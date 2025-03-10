package publishers

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type ArticleQueuePublisher struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
	queue      amqp.Queue
}

func NewArticleQueuePublisher(connection *amqp.Connection, queueName string) (*ArticleQueuePublisher, error) {
	articleQueuePublisher := &ArticleQueuePublisher{connection: connection, queueName: queueName}
	if err := articleQueuePublisher.setupChannel(); err != nil {
		return nil, err
	}
	if err := articleQueuePublisher.setupQueue(); err != nil {
		return nil, err
	}
	return articleQueuePublisher, nil
}

func (p *ArticleQueuePublisher) setupChannel() error {
	if p.channel != nil {
		return nil
	}

	channel, err := p.connection.Channel()
	if err != nil {
		return err
	}

	p.channel = channel
	return nil
}

func (p *ArticleQueuePublisher) setupQueue() error {
	queue, err := p.channel.QueueDeclare(p.queueName, true, false, false, false, nil)
	if err != nil {
		p.channel = nil
		return err
	}
	p.queue = queue
	return nil
}

func (p *ArticleQueuePublisher) PublishArticle(ctx context.Context, article *models.Article) error {
	if err := p.channel.Publish("", p.queue.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(article.ID.Hex()),
	}); err != nil {
		return err
	}
	return nil
}
