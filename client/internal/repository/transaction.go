package repository

import (
	"context"

	"github.com/ellofae/payment-system-kafka/client/internal/domain"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/hashicorp/go-hclog"
)

type TransactionRepository struct {
	logger  hclog.Logger
	storage *Storage
}

func NewTransactionRepository(storage *Storage) domain.ITransactionRepository {
	return &TransactionRepository{
		logger:  logger.GetLogger(),
		storage: storage,
	}
}

const attach_transaction string = `INSERT INTO transactions(user_id, transaction_id, card_number, amount) 
VALUES($1, $2, $3, $4)
RETURNING id`

func (r *TransactionRepository) AttachTrasaction(ctx context.Context, userID int, transactionID string, cardNumber string, amount float64) (int, error) {
	conn, err := r.storage.GetPgConnPool().Acquire(ctx)
	if err != nil {
		return -1, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback(ctx)

	var receivedID int
	err = tx.QueryRow(ctx, attach_transaction, userID, transactionID, cardNumber, amount).Scan(&receivedID)
	if err != nil {
		return -1, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return -1, err
	}

	return receivedID, nil
}
