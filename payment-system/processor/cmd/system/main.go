package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ellofae/payment-system-kafka/config"
	"github.com/ellofae/payment-system-kafka/internal/encryption"
	"github.com/ellofae/payment-system-kafka/payment-system/processor/internal/domain/entity"
	"github.com/ellofae/payment-system-kafka/payment-system/processor/internal/processing"
	"github.com/ellofae/payment-system-kafka/pkg/logger"

	"github.com/pkg/profile"
)

const topic string = "purchases"
const MIN_COMMIT_COUNT = 5

func main() {
	log := logger.GetLogger()
	cfg := config.ParseConfig(config.ConfigureViper())

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
	defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
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

				// transactionData.CardNumber, err = encryption.DecryptData(transactionData.CardNumber)
				// if err != nil {
				// 	os.Exit(1)
				// }

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
