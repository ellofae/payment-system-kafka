package producer

import (
	"bytes"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ellofae/payment-system-kafka/payment-system/internal/data"
	"github.com/ellofae/payment-system-kafka/payment-system/pkg/logger"
)

type TransactionProducer struct {
	producer          *kafka.Producer
	topic             string
	transactionStatus chan kafka.Event
}

func NewTransactionProducer(p *kafka.Producer, topic string) *TransactionProducer {
	return &TransactionProducer{
		producer:          p,
		topic:             topic,
		transactionStatus: make(chan kafka.Event, 10000),
	}
}

func (tp *TransactionProducer) ProcessTransaction(transactionData *data.TransactionData) error {
	log := logger.GetLogger()

	buffer := &bytes.Buffer{}
	err := json.NewEncoder(buffer).Encode(transactionData)
	if err != nil {
		log.Error("Unable to encode transaction data", "error", err.Error())
		return err
	}

	dataBytes := buffer.Bytes()

	err = tp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &tp.topic, Partition: kafka.PartitionAny},
		Value:          dataBytes,
	}, tp.transactionStatus)

	if err != nil {
		log.Error("Unable to send transaction data to a consumer", "error", err.Error())
		return err
	}

	<-tp.transactionStatus
	log.Info("Placed transaction on the queue, transaction ID: %d", transactionData.TransactionID)

	return nil
}
