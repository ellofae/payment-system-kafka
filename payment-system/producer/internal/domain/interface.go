package domain

import (
	"context"

	"github.com/ellofae/payment-system-kafka/payment-system/producer/internal/domain/entity"
)

type ITransactionUsecase interface {
	PlaceTransaction(context.Context, *entity.TransactionData) error
}
