package articlefeed

import (
	"context"
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

type ArticleFeedWorker struct {
	logger          *slog.Logger
	queueConsumer   app.Consumer
	service         articleGetter
	profileManager  profileGetter
	followerCentral followersGetter
	feedAppender    feedAppender
}

func NewArticleFeedWorker(articleWriteQueueConsumer app.Consumer, articleGetter articleGetter, profileGetter profileGetter, followersGetter followersGetter, feedAppender feedAppender, logger *slog.Logger) *ArticleFeedWorker {
	return &ArticleFeedWorker{
		logger:          logger,
		queueConsumer:   articleWriteQueueConsumer,
		service:         articleGetter,
		profileManager:  profileGetter,
		followerCentral: followersGetter,
		feedAppender:    feedAppender,
	}
}

func (w *ArticleFeedWorker) Consume() {
	go w.queueConsumer.StartConsumer()
	messageQueue := w.queueConsumer.Consume()
	for message := range messageQueue {
		ctx := context.Background()
		articleID := string(message.Data())
		w.logger.Debug("Received a message", "message", articleID)

		article, err := w.service.GetArticleByID(ctx, articleID)
		if err != nil {
			w.logger.Error("Failed to find article", "articleID", articleID, "error", err)
			w.success(message)
			continue
		}
		w.logger.Debug("Found Article", "article", article)

		author, err := w.profileManager.GetUserByID(ctx, *article.Author)
		if err != nil {
			w.logger.Error("Failed to find article author", "authorID", *article.Author, "error", err)
			w.success(message)
			continue
		}
		w.logger.Debug("Found Author", "author", author)

		followers, err := w.followerCentral.GetFollowers(ctx, author.ID.Hex())
		if err != nil {
			w.logger.Error("Failed to author's followers", "authorID", *article.Author, "error", err)
			w.success(message)
			continue
		}
		w.logger.Debug("Found Followers", "followers", followers)

		followerIDs := make([]string, 0, len(followers))
		for _, follower := range followers {
			followerIDs = append(followerIDs, *follower.Follower)
		}

		// TODO: Only do appending of article for active users (last 30 days)
		if err := w.feedAppender.AppendArticleToUserFeeds(ctx, article, author, followerIDs); err != nil {
			w.logger.Error("Failed to write feed for users", "error", err)
			w.success(message)
			continue
		}

		w.success(message)
	}
}

func (w *ArticleFeedWorker) success(message app.Message) {
	if err := message.Success(); err != nil {
		w.logger.Error("Failed to ACK message", "error", err)
	}
}
