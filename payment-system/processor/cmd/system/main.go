package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ellofae/payment-system-kafka/config"
	"github.com/ellofae/payment-system-kafka/payment-system/processor/internal/domain/entity"
	"github.com/ellofae/payment-system-kafka/payment-system/processor/internal/processing"
	"github.com/ellofae/payment-system-kafka/payment-system/processor/internal/repository"
	"github.com/ellofae/payment-system-kafka/pkg/encryption"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/ellofae/payment-system-kafka/pkg/redis"
)

const topic string = "purchases"
const MIN_COMMIT_COUNT = 5

func main() {
	log := logger.GetLogger()
	cfg := config.ParseConfig(config.ConfigureViper())
	ctx := context.Background()

	redisClient := redis.OpenRedisConnection(ctx, cfg)
	redisRepository := repository.NewStroage(redisClient.Client)

	encryption.InitializeEncryptionKey(cfg)

	c, err := processing.InitializeConsumer(cfg)
	if err != nil {
		os.Exit(1)
	}
	defer c.Close()

	err = processing.NewTransactionConsumer(c, topic)
	if err != nil {
		os.Exit(1)
	}

	log.Info("Starting processing transactions..")

	var wg sync.WaitGroup
	var intrSignal bool

	msg_count := 0
	for {
		intrSignal = false
		eventValue := c.Poll(100)

		switch ev := eventValue.(type) {
		case *kafka.Message:
			msg_count += 1
			if msg_count%MIN_COMMIT_COUNT == 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()

					_, err := c.Commit()
					if err != nil {
						log.Warn("Unable to commit the offset", "error", err.Error())
						intrSignal = true
					}
				}()
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				transactionData := &entity.TransactionData{}

				decoder := json.NewDecoder(bytes.NewReader(ev.Value))
				err := decoder.Decode(transactionData)
				if err != nil {
					log.Error("Unable to decode transaction data", "error", err.Error())
					return
				}

				if err := redisRepository.Insert(ctx, transactionData); err != nil {
					log.Error("Failed to store consumed transaction", "error", err.Error())
					return
				}
				fmt.Printf("processed transaction: %v\n", transactionData)
			}()
		case kafka.PartitionEOF:
			log.Warn("reached", "event value", ev)
		case kafka.Error:
			log.Error("kafka error has occured", "error", ev)
			intrSignal = true
		}

		if intrSignal {
			break
		}
	}

	log.Error("Interrupt signal has been caught, something went wrong..")

	wg.Wait()
	os.Exit(1)
}
