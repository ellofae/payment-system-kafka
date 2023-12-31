package domain

import (
	"context"

	"github.com/ellofae/payment-system-kafka/client/internal/dto"
)

type (
	IAuthenticationUsecase interface {
		ValidateDTOStruct(interface{}) error
		SignUp(context.Context, *dto.UserCreationForm) (int, error)
		SignIn(context.Context, *dto.UserLoginForm) (string, error)
	}

	ITransactionUsecase interface {
		ValidateDTOStruct(interface{}) error
		PlaceTransaction(*dto.TransactionData) error
		AttachTransaction(context.Context, *dto.TransactionData) error
	}

	IUserUsecase interface {
		GetUserTransaction(context.Context, int) ([]*dto.TransactionDisplayData, error)
		GetUserData(context.Context, int) (*dto.UserData, error)
	}
)

type (
	IAuthenticationRepository interface {
		GetUserCredByEmail(context.Context, string) (*dto.CredentialDTO, error)
		SignUp(context.Context, *dto.UserCreationForm, string) (int, error)
		GetUserIDByEmail(context.Context, string) (int, error)
	}

	ITransactionRepository interface {
		AttachTrasaction(context.Context, int, string, string, float64) (int, error)
	}

	IUserRepository interface {
		GetUserTransaction(context.Context, int) ([]*dto.TransactionDisplayData, error)
		GetUserData(context.Context, int) (*dto.UserIntermediateData, error)
	}
)
