package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ellofae/payment-system-kafka/config"
	"github.com/ellofae/payment-system-kafka/payment-system/producer/internal/controller/handler"
	"github.com/ellofae/payment-system-kafka/payment-system/producer/internal/domain/usecase"
	"github.com/ellofae/payment-system-kafka/payment-system/producer/internal/producing"
	"github.com/ellofae/payment-system-kafka/pkg/encryption"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/gin-gonic/gin"
)

const topic string = "purchases"

func Run() {
	cfg := config.ParseConfig(config.ConfigureViper())
	ctx := context.Background()

	encryption.InitializeEncryptionKey(cfg)

	p, err := InitializeProducer(cfg)
	if err != nil {
		os.Exit(1)
	}

	producer := producing.NewTransactionProducer(p, topic)
	router := InitRouter(producer)

	srv := InitHTTPServer(router, cfg)
	StartServer(ctx, srv)
}

func InitializeProducer(cfg *config.Config) (*kafka.Producer, error) {
	log := logger.GetLogger()

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", cfg.Kafka.BootstrapServersHost, cfg.Kafka.BootstrapServersPort),
		"client.id":         cfg.Kafka.ProducerID,
		"acks":              cfg.Kafka.Acks,
	})

	if err != nil {
		log.Error("Unable to start a transaction producer", "error", err)
		return nil, err
	}

	return producer, nil
}

func InitRouter(producer *producing.TransactionProducer) *gin.Engine {
	r := gin.Default()

	transactionUsecase := usecase.NewTransactionUsecase(producer)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	transactionHandler.Register(r)

	return r
}

func InitHTTPServer(router *gin.Engine, cfg *config.Config) http.Server {
	readTimeoutSecondsCount, _ := strconv.Atoi(cfg.ProducerServer.ReadTimeout)
	writeTimeoutSecondsCount, _ := strconv.Atoi(cfg.ProducerServer.WriteTimeout)
	idleTimeoutSecondsCount, _ := strconv.Atoi(cfg.ProducerServer.IdleTimeout)

	bindAddr := cfg.ProducerServer.BindAddr

	srv := http.Server{
		Addr:         bindAddr,
		Handler:      router,
		ReadTimeout:  time.Duration(readTimeoutSecondsCount) * time.Second,
		WriteTimeout: time.Duration(writeTimeoutSecondsCount) * time.Second,
		IdleTimeout:  time.Duration(idleTimeoutSecondsCount) * time.Second,
	}

	return srv
}

func StartServer(ctx context.Context, srv http.Server) {
	log := logger.GetLogger()

	go func() {
		log.Info("Starting server...")
		err := srv.ListenAndServe()
		if err != nil {
			log.Info("Server was stopped")
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	signal := <-sigChan
	log.Info(fmt.Sprintf("Signal has been caught: %v", signal))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
}
