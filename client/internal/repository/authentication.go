package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ellofae/payment-system-kafka/client/internal/domain"
	"github.com/ellofae/payment-system-kafka/client/internal/domain/entity"
	"github.com/ellofae/payment-system-kafka/client/internal/dto"
	response_errors "github.com/ellofae/payment-system-kafka/client/internal/errors"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v5"
)

type AuthenticationRepository struct {
	logger  hclog.Logger
	storage *Storage
}

func NewAuthenticationRepository(storage *Storage) domain.IAuthenticationRepository {
	return &AuthenticationRepository{
		logger:  logger.GetLogger(),
		storage: storage,
	}
}

const get_user_id_by_email_query = `SELECT users.id FROM users
LEFT JOIN credentials ON users.credential_id = credentials.id
WHERE credentials.email = $1`

func (r *AuthenticationRepository) GetUserIDByEmail(ctx context.Context, email string) (int, error) {
	var id int

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

	err = tx.QueryRow(ctx, get_user_id_by_email_query, email).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, response_errors.ErrNoRecordFound
		}
		return -1, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return -1, err
	}

	return id, nil
}

const get_user_by_email_query = `SELECT * FROM credentials WHERE email = $1`

func (r *AuthenticationRepository) GetUserCredByEmail(ctx context.Context, email string) (*dto.CredentialDTO, error) {
	entity := &entity.Credential{}

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

	err = tx.QueryRow(ctx, get_user_by_email_query, email).Scan(&entity.ID, &entity.Email, &entity.Password, &entity.RegisterDate)
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

	dto_model := &dto.CredentialDTO{
		ID:       entity.ID,
		Email:    entity.Email,
		Password: entity.Password,
	}

	return dto_model, nil
}

const credential_creation_query = `INSERT INTO credentials(email, password_hash, register_date) 
VALUES($1, $2, $3)
RETURNING id
`

const user_creation_query = `INSERT INTO users(first_name, last_name, credential_id) 
VALUES($1, $2, $3)
RETURNING id
`

func (r *AuthenticationRepository) SignUp(ctx context.Context, user_dto *dto.UserCreationForm, hashed_password string) (int, error) {
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

	var credentialID int

	err = tx.QueryRow(ctx, credential_creation_query, user_dto.Email, hashed_password, time.Now()).Scan(&credentialID)
	if err != nil {
		return -1, err
	}

	var userID int

	err = tx.QueryRow(ctx, user_creation_query, user_dto.FirstName, user_dto.LastName, credentialID).Scan(&userID)
	if err != nil {
		return -1, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return -1, err
	}

	return userID, nil
}
