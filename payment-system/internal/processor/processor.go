package processor

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ellofae/payment-system-kafka/payment-system/pkg/logger"
)

func InitializeConsumer() (*kafka.Consumer, error) {
	log := logger.GetLogger()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "transaction_consumer",
		"auto.offset.reset": "smallest",
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
