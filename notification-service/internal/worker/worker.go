package worker

import (
	"context"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type NotificationWorker struct {
	Conn    *amqp091.Connection
	Channel *amqp091.Channel
	queue   string
}

func New(ctx context.Context, amqpURL, queueName string) (*NotificationWorker, error) {
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &NotificationWorker{Conn: conn, Channel: ch, queue: queueName}, nil
}

func (w *NotificationWorker) Start(ctx context.Context) error {
	msgs, err := w.Channel.Consume(w.queue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	log.Println("📥 Notification worker started, waiting for messages...")

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return nil
			}
			w.process(msg)
			err := msg.Ack(false)
			if err != nil {
				return err
			} // Подтверждаем обработку
		}
	}
}

func (w *NotificationWorker) process(msg amqp091.Delivery) {
	log.Printf("📧 Processing notification: %s", string(msg.Body))

	// Имитация тяжелой работы (отправка email, генерация PDF)
	time.Sleep(2 * time.Second)

	log.Printf("✅ Notification sent for message: %s", string(msg.Body))
}
