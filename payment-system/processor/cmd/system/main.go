package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ellofae/payment-system-kafka/payment-system/data"
	"github.com/ellofae/payment-system-kafka/payment-system/processor/internal/processing"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
)

const topic string = "purchases"

func main() {
	log := logger.GetLogger()

	c, err := processing.InitializeConsumer()
	if err != nil {
		os.Exit(1)
	}

	err = processing.NewTransactionConsumer(c, topic)
	if err != nil {
		os.Exit(1)
	}

	log.Info("Starting processing transactions..")
	for {
		eventValue := c.Poll(100)

		switch eventValue := eventValue.(type) {
		case *kafka.Message:
			transactionData := &data.TransactionData{}

			decoder := json.NewDecoder(bytes.NewReader(eventValue.Value))
			err := decoder.Decode(transactionData)
			if err != nil {
				log.Error("Unable to decode transaction data", "error", err.Error())
				os.Exit(1)
			}

			fmt.Printf("processed transaction: %v\n", transactionData)
		case kafka.Error:
			fmt.Printf("kafka error: %v\n", eventValue)
		}
	}
}
