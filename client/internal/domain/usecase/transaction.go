package usecase

import (
	"github.com/ellofae/payment-system-kafka/client/internal/domain"
	"github.com/ellofae/payment-system-kafka/client/internal/dto"
	"github.com/ellofae/payment-system-kafka/client/internal/utils"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/hashicorp/go-hclog"
)

type TransactionUsecase struct {
	logger hclog.Logger
}

func NewTransactionUsecase() domain.ITransactionUsecase {
	return &TransactionUsecase{
		logger: logger.GetLogger(),
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

	return nil
}
