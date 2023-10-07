package usecase

import (
	"context"
	"time"

	"github.com/ellofae/payment-system-kafka/client/internal/domain"
	"github.com/ellofae/payment-system-kafka/client/internal/dto"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/hashicorp/go-hclog"
)

type UserUsecase struct {
	logger hclog.Logger
	repo   domain.IUserRepository
}

func NewUserUsecase(repo domain.IUserRepository) domain.IUserUsecase {
	return &UserUsecase{
		logger: logger.GetLogger(),
		repo:   repo,
	}
}

func (u *UserUsecase) GetUserTransaction(ctx context.Context, userID int) ([]*dto.TransactionDisplayData, error) {
	transaction, err := u.repo.GetUserTransaction(ctx, userID)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (u *UserUsecase) GetUserData(ctx context.Context, userID int) (*dto.UserData, error) {
	intermediateData, err := u.repo.GetUserData(ctx, userID)
	if err != nil {
		return nil, err
	}

	formatedDate := intermediateData.RegisterDate.Format(time.RFC1123)

	userData := &dto.UserData{
		ID:           intermediateData.ID,
		FirstName:    intermediateData.FirstName,
		LastName:     intermediateData.LastName,
		Email:        intermediateData.Email,
		RegisterDate: formatedDate,
	}

	return userData, nil
}
