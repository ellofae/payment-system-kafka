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
)

type (
	IAuthenticationRepository interface {
		GetUserCredByEmail(context.Context, string) (*dto.CredentialDTO, error)
		GetRoleByTitle(context.Context, string) (int, error)
		SignUp(context.Context, *dto.UserCreationForm, string) (int, error)
		GetUserIDByEmail(context.Context, string) (int, error)
	}
)
