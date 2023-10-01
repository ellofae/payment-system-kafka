package main

import (
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ellofae/payment-system-kafka/payment-system/data"
	"github.com/ellofae/payment-system-kafka/payment-system/producer/internal/producing"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
)

const topic string = "purchases"

func InitializeProducer() (*kafka.Producer, error) {
	log := logger.GetLogger()

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "transaction_producer",
		"acks":              "all",
	})

	if err != nil {
		log.Error("Unable to start a transaction producer", "error", err)
		return nil, err
	}

	return producer, nil
}

func main() {
	log := logger.GetLogger()

	p, err := InitializeProducer()
	if err != nil {
		os.Exit(1)
	}

	transactionProducer := producing.NewTransactionProducer(p, topic)

	log.Info("Starting producing transactions..")
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
