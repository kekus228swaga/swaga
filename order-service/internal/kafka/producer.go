package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/kekus228swaga/orderflow/order-service/internal/publisher"
)

type Producer struct {
	asyncProducer sarama.AsyncProducer
	topic         string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	// Горутина для обработки успехов/ошибок (обязательно для AsyncProducer)
	go func() {
		for range producer.Successes() {
			// log.Println("✅ Kafka message sent")
		}
	}()
	go func() {
		for err := range producer.Errors() {
			log.Printf("❌ Kafka async error: %v", err)
		}
	}()

	return &Producer{asyncProducer: producer, topic: topic}, nil
}

func (p *Producer) Send(event publisher.OrderEvent) {
	body, err := json.Marshal(event)
	if err != nil {
		log.Printf("❌ JSON marshal error: %v", err)
		return
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(body),
	}

	p.asyncProducer.Input() <- msg
}
