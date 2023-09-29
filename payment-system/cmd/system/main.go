package main

import (
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ellofae/payment-system-kafka/payment-system/internal/data"
	"github.com/ellofae/payment-system-kafka/payment-system/internal/producer"
	"github.com/ellofae/payment-system-kafka/payment-system/pkg/logger"
)

const topic string = "Payment System"

func InitializeProducer() (*kafka.Producer, error) {
	log := logger.GetLogger()

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "payment_producer",
		"acks":              "all",
	})

	if err != nil {
		log.Error("Unable to start a payment producer", "error", err)
		return nil, err
	}

	return producer, nil
}

func main() {
	p, err := InitializeProducer()
	if err != nil {
		os.Exit(1)
	}

	transactionProducer := producer.NewTransactionProducer(p, topic)
	for i := 1; i <= 10000; i++ {
		data := &data.TransactionData{
			UserID:        1,
			TransactionID: i,
			CardNumber:    "xxx-1024-5213", // encription suposed to be done here
			Description:   "transaction description",
			Amount:        100.0,
		}

		if err := transactionProducer.ProcessTransaction(data); err != nil {
			os.Exit(1)
		}
		time.Sleep(time.Second * 3)
	}
}
