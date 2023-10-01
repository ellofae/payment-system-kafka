package processing

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ellofae/payment-system-kafka/config"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
)

func InitializeConsumer(cfg *config.Config) (*kafka.Consumer, error) {
	log := logger.GetLogger()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", cfg.Kafka.BootstrapServersPort, cfg.Kafka.BootstrapServersHost),
		"group.id":          cfg.Kafka.GroupID,
		"auto.offset.reset": cfg.Kafka.AutoOffsetReset,
	})

	if err != nil {
		log.Error("Unable to start a transaction producer", "error", err.Error())
		return nil, err
	}

	return consumer, nil
}

func NewTransactionConsumer(consumer *kafka.Consumer, topic string) error {
	log := logger.GetLogger()

	if err := consumer.Subscribe(topic, nil); err != nil {
		log.Error("Unable to subscribe a transaction consumer", "error", err.Error())
		return err
	}
	return nil
}
