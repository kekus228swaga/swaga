package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type AnalyticsConsumer struct {
	consumerGroup sarama.ConsumerGroup
	topic         string
}

func New(brokers []string, groupID, topic string) (*AnalyticsConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // Читать с самого начала

	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &AnalyticsConsumer{consumerGroup: client, topic: topic}, nil
}

func (ac *AnalyticsConsumer) Start(ctx context.Context) {
	handler := &ConsumerHandler{}

	for {
		// Consume блокирует поток и слушает
		if err := ac.consumerGroup.Consume(ctx, []string{ac.topic}, handler); err != nil {
			log.Printf("❌ Error from consumer: %s", err.Error())
		}
		if ctx.Err() != nil {
			return
		}
	}
}

// ConsumerHandler реализует интерфейс sarama.ConsumerGroupHandler
type ConsumerHandler struct{}

func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		fmt.Printf("📊 Analytics: Received message: %s\n", string(message.Value))
		session.MarkMessage(message, "") // Подтверждаем прочтение
	}
	return nil
}
