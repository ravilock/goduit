package producers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type FeedProducer struct {
	Channel     *amqp.Channel
	ChannelName string
}

func NewFeedProducer(channel *amqp.Channel, channelName string) (*FeedProducer, error) {
	_, err := channel.QueueDeclare(channelName, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return &FeedProducer{channel, channelName}, nil
}

func (p *FeedProducer) Publish(ctx context.Context, article *models.Article) error {
	// TODO: Handle error
	articleJson, err := json.Marshal(article)
	if err != nil {
		return err
	}
	return p.Channel.PublishWithContext(ctx, "", p.ChannelName, false, false, amqp.Publishing{
		ContentType: "application/json",
		// CorrelationId: "", // TODO: Should get request id?
		MessageId: uuid.NewString(),
		Timestamp: time.Now(),
		Body:      articleJson,
	})
}
