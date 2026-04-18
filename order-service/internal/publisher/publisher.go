package publisher

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type OrderEvent struct {
	OrderID     int64  `json:"order_id"`
	UserID      int64  `json:"user_id"`
	ProductName string `json:"product_name"`
}

type Publisher struct {
	Channel *amqp091.Channel
	queue   string
}

func NewPublisher(amqpURL, queueName string) (*Publisher, error) {
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Объявляем очередь (если её нет, создаст; если есть, пропустит)
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &Publisher{Channel: ch, queue: queueName}, nil
}

func (p *Publisher) Publish(ctx context.Context, event OrderEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.Channel.PublishWithContext(ctx, "", p.queue, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		return err
	}

	log.Printf("📤 Published to RabbitMQ: %+v", event)
	return nil
}
