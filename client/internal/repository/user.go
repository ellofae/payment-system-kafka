package repository

import (
	"context"
	"errors"

	"github.com/ellofae/payment-system-kafka/client/internal/domain"
	"github.com/ellofae/payment-system-kafka/client/internal/dto"
	response_errors "github.com/ellofae/payment-system-kafka/client/internal/errors"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	logger  hclog.Logger
	storage *Storage
}

func NewUserRepository(storage *Storage) domain.IUserRepository {
	return &UserRepository{
		logger:  logger.GetLogger(),
		storage: storage,
	}
}

const get_user_transactions string = `SELECT * FROM transactions WHERE user_id = $1`

func (r *UserRepository) GetUserTransaction(ctx context.Context, userID int) ([]*dto.TransactionDisplayData, error) {
	conn, err := r.storage.GetPgConnPool().Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var transactions []*dto.TransactionDisplayData

	rows, err := tx.Query(ctx, get_user_transactions, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, response_errors.ErrNoRecordFound
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction dto.TransactionDisplayData
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.TransactionID, &transaction.CardNumber, &transaction.Amount); err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

const get_current_user_data string = `SELECT users.id, users.first_name, users.last_name,
credentials.email, credentials.register_date FROM users
FULL JOIN credentials
ON users.credential_id = credentials.id
WHERE users.id = $1`

func (r *UserRepository) GetUserData(ctx context.Context, userId int) (*dto.UserIntermediateData, error) {
	intermediateData := dto.UserIntermediateData{}

	conn, err := r.storage.GetPgConnPool().Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, get_current_user_data, userId).Scan(&intermediateData.ID, &intermediateData.FirstName, &intermediateData.LastName, &intermediateData.Email, &intermediateData.RegisterDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, response_errors.ErrNoRecordFound
		}
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &intermediateData, nil
}
