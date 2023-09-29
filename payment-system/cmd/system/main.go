package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ellofae/payment-system-kafka/payment-system/internal/data"
	"github.com/ellofae/payment-system-kafka/payment-system/internal/processor"
	"github.com/ellofae/payment-system-kafka/payment-system/internal/producer"
	"github.com/ellofae/payment-system-kafka/payment-system/pkg/logger"
)

const topic string = "Payment System"

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

	c, err := processor.InitializeConsumer()
	if err != nil {
		os.Exit(1)
	}

	transactionProducer := producer.NewTransactionProducer(p, topic)
	go func() {
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
	}()

	err = processor.NewTransactionConsumer(c, topic)
	if err != nil {
		os.Exit(1)
	}

	go func() {
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
			default:
				fmt.Printf("default behavior")
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	signal := <-sigChan
	log.Info("Signal has been caught", "signal", signal)

	os.Exit(0)
}
