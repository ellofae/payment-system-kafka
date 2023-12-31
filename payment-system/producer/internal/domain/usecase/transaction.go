package usecase

import (
	"context"
	"time"

	"github.com/ellofae/payment-system-kafka/payment-system/producer/internal/domain"
	"github.com/ellofae/payment-system-kafka/payment-system/producer/internal/domain/entity"
	"github.com/ellofae/payment-system-kafka/payment-system/producer/internal/producing"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/hashicorp/go-hclog"
)

type TransactionUsecase struct {
	logger   hclog.Logger
	producer *producing.TransactionProducer
}

func NewTransactionUsecase(producer *producing.TransactionProducer) domain.ITransactionUsecase {
	return &TransactionUsecase{
		logger:   logger.GetLogger(),
		producer: producer,
	}
}

func (tu *TransactionUsecase) PlaceTransaction(ctx context.Context, data *entity.TransactionData) error {
	if err := tu.producer.ProduceTransaction(data); err != nil {
		return err
	}
	time.Sleep(time.Second * 3)

	return nil
}
