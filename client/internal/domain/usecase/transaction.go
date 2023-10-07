package usecase

import (
	"context"

	"github.com/ellofae/payment-system-kafka/client/internal/domain"
	"github.com/ellofae/payment-system-kafka/client/internal/dto"
	"github.com/ellofae/payment-system-kafka/client/internal/utils"
	"github.com/ellofae/payment-system-kafka/pkg/encryption"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/hashicorp/go-hclog"
)

type TransactionUsecase struct {
	logger hclog.Logger
	repo   domain.ITransactionRepository
}

func NewTransactionUsecase(repo domain.ITransactionRepository) domain.ITransactionUsecase {
	return &TransactionUsecase{
		logger: logger.GetLogger(),
		repo:   repo,
	}
}

func (u *TransactionUsecase) ValidateDTOStruct(model interface{}) error {
	validate := utils.NewValidator()

	if err := validate.Struct(model); err != nil {
		validation_errors := utils.ValidatorErrors(err)
		for _, errValidation := range validation_errors {
			u.logger.Error("Validation error occured", "error", errValidation)
		}

		return err
	}

	return nil
}

func (u *TransactionUsecase) PlaceTransaction(req *dto.TransactionData) error {
	if err := u.ValidateDTOStruct(req); err != nil {
		return err
	}

	req.CardNumber = encryption.EncryptData([]byte(req.CardNumber))

	return nil
}

func (u *TransactionUsecase) AttachTransaction(ctx context.Context, req *dto.TransactionData) error {
	_, err := u.repo.AttachTrasaction(ctx, req.UserID, req.TransactionID)
	if err != nil {
		u.logger.Error("Unable to store user's transaction", "error", err.Error())
		return err
	}

	return nil
}
