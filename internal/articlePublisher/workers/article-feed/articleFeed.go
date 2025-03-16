package articlefeed

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/articlePublisher/models"
	followerCentralModels "github.com/ravilock/goduit/internal/followerCentral/models"
	profileManagerModels "github.com/ravilock/goduit/internal/profileManager/models"
)

type articleGetter interface {
	GetArticleByID(ctx context.Context, articleID string) (*models.Article, error)
}

type profileGetter interface {
	GetUserByID(ctx context.Context, ID string) (*profileManagerModels.User, error)
}

type followersGetter interface {
	GetFollowers(ctx context.Context, followed string) ([]*followerCentralModels.Follower, error)
}

type feedAppender interface {
	AppendArticleToUserFeeds(ctx context.Context, article *models.Article, author *profileManagerModels.User, userIDs []string) error
}

type ArticleFeedHandler struct {
	logger          *slog.Logger
	service         articleGetter
	profileManager  profileGetter
	followerCentral followersGetter
	feedAppender    feedAppender
}

func NewArticleFeedHandler(articleGetter articleGetter, profileGetter profileGetter, followersGetter followersGetter, feedAppender feedAppender, logger *slog.Logger) *ArticleFeedHandler {
	logger = logger.With("emitter", "article-feed-worker")
	return &ArticleFeedHandler{
		logger:          logger,
		service:         articleGetter,
		profileManager:  profileGetter,
		followerCentral: followersGetter,
		feedAppender:    feedAppender,
	}
}

func (w *ArticleFeedHandler) Handle(message app.Message) {
	ctx := context.Background()
	articleID := string(message.Data())
	w.logger.Debug("Received a message", "message", articleID)

	article, err := w.service.GetArticleByID(ctx, articleID)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.ArticleNotFoundErrorCode:
				w.logger.Info("Article not found", "articleID", articleID)
				w.success(message)
				return
			}
		}
		w.logger.Error("Failed to find article", "articleID", articleID, "error", err)
		w.failure(message)
		return
	}
	w.logger.Debug("Found Article", "article", article)

	author, err := w.profileManager.GetUserByID(ctx, *article.Author)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.UserNotFoundErrorCode:
				w.logger.Info("Author not found", "authorID", *article.Author)
				w.success(message)
				return
			}
		}
		w.logger.Error("Failed to find article author", "authorID", *article.Author, "error", err)
		w.failure(message)
		return
	}
	w.logger.Debug("Found Author", "author", author)

	followers, err := w.followerCentral.GetFollowers(ctx, author.ID.Hex())
	if err != nil {
		w.logger.Error("Failed to author's followers", "authorID", *article.Author, "error", err)
		w.failure(message)
		return
	}

	if len(followers) == 0 {
		w.logger.Debug("No Followers Found", "author", author)
		w.success(message)
		return
	}
	w.logger.Debug("Found Followers", "followers", followers)

	followerIDs := make([]string, 0, len(followers))
	for _, follower := range followers {
		followerIDs = append(followerIDs, *follower.Follower)
	}

	// TODO: Only do appending of article for active users (last 30 days)
	if err := w.feedAppender.AppendArticleToUserFeeds(ctx, article, author, followerIDs); err != nil {
		w.logger.Error("Failed to write feed for users", "error", err)
		w.failure(message)
		return
	}

	w.success(message)
	w.logger.Debug("Successfully appended article to user feeds")
}

func (w *ArticleFeedHandler) success(message app.Message) {
	if err := message.Success(); err != nil {
		w.logger.Error("Failed to ACK message", "error", err)
	}
}

func (w *ArticleFeedHandler) failure(message app.Message) {
	if err := message.Failure(); err != nil {
		w.logger.Error("Failed to ACK message", "error", err)
	}
}
