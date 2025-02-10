package articlefeed

import (
	"context"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
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
	connection      *amqp.Connection
	channel         *amqp.Channel
	errorChan       chan error
	logger          *slog.Logger
	service         articleGetter
	profileManager  profileGetter
	followerCentral followersGetter
	feedAppender    feedAppender
	queueName       string
	queue           amqp.Queue
}

func NewArticleFeedWorker(connection *amqp.Connection, articleGetter articleGetter, profileGetter profileGetter, followersGetter followersGetter, feedAppender feedAppender, queueName string, errorChan chan error, logger *slog.Logger) (*ArticleFeedWorker, error) {
	articleQueuePublisher := &ArticleFeedWorker{
		connection:      connection,
		queueName:       queueName,
		errorChan:       errorChan,
		logger:          logger,
		service:         articleGetter,
		profileManager:  profileGetter,
		followerCentral: followersGetter,
		feedAppender:    feedAppender,
	}
	if err := articleQueuePublisher.setupChannel(); err != nil {
		return nil, err
	}
	if err := articleQueuePublisher.setupQueue(); err != nil {
		return nil, err
	}
	return articleQueuePublisher, nil
}

func (w *ArticleFeedWorker) setupChannel() error {
	if w.channel != nil {
		return nil
	}

	channel, err := w.connection.Channel()
	if err != nil {
		return err
	}

	w.channel = channel
	return nil
}

func (w *ArticleFeedWorker) setupQueue() error {
	queue, err := w.channel.QueueDeclare(w.queueName, true, false, false, false, nil)
	if err != nil {
		w.channel = nil
		return err
	}
	w.queue = queue
	return nil
}

func (w *ArticleFeedWorker) Consume() {
	messageDeliveryChan, err := w.channel.Consume(w.queue.Name, "", false, false, false, false, nil)
	if err != nil {
		w.errorChan <- err
		return
	}
	for messageDelivery := range messageDeliveryChan {
		ctx := context.Background()
		articleID := string(messageDelivery.Body)
		w.logger.Debug("Received a message", "message", articleID)

		article, err := w.service.GetArticleByID(ctx, articleID)
		if err != nil {
			w.logger.Error("Failed to find article", "articleID", articleID, "error", err)
			w.ack(messageDelivery)
			continue
		}
		w.logger.Debug("Found Article", "article", article)

		author, err := w.profileManager.GetUserByID(ctx, *article.Author)
		if err != nil {
			w.logger.Error("Failed to find article author", "authorID", *article.Author, "error", err)
			w.ack(messageDelivery)
			continue
		}
		w.logger.Debug("Found Author", "author", author)

		followers, err := w.followerCentral.GetFollowers(ctx, author.ID.Hex())
		if err != nil {
			w.logger.Error("Failed to author's followers", "authorID", *article.Author, "error", err)
			w.ack(messageDelivery)
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
			continue
		}

		w.ack(messageDelivery)
	}
}

func (w *ArticleFeedWorker) ack(delivery amqp.Delivery) {
	if err := delivery.Ack(false); err != nil {
		w.logger.Error("Failed to ACK message", "error", err)
	}
}
