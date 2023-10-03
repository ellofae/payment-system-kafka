package usecase

import (
	"context"

	"github.com/ellofae/payment-system-kafka/client/internal/controller/middleware"
	"github.com/ellofae/payment-system-kafka/client/internal/domain"
	"github.com/ellofae/payment-system-kafka/client/internal/dto"
	response_errors "github.com/ellofae/payment-system-kafka/client/internal/errors"
	"github.com/ellofae/payment-system-kafka/client/internal/utils"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/hashicorp/go-hclog"
)

type AuthenticationUsecase struct {
	logger hclog.Logger
	repo   domain.IAuthenticationRepository
}

func NewAuthenticationUsecase(repo domain.IAuthenticationRepository) *AuthenticationUsecase {
	return &AuthenticationUsecase{
		logger: logger.GetLogger(),
		repo:   repo,
	}
}

func (u *AuthenticationUsecase) ValidateDTOStruct(model interface{}) error {
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

func (u *AuthenticationUsecase) SignUp(ctx context.Context, user_dto *dto.UserCreationForm) (int, error) {
	if err := u.ValidateDTOStruct(user_dto); err != nil {
		return -1, err
	}

	_, err := u.repo.GetUserCredByEmail(ctx, user_dto.Email)
	if err == nil {
		return -1, response_errors.ErrAlreadyExists
	} else if err != response_errors.ErrNoRecordFound {
		return -1, err
	}

	hashed_password, err := utils.HashPassword(user_dto.Password)
	if err != nil {
		return -1, err
	}

	return u.repo.SignUp(ctx, user_dto, hashed_password)
}

func (u *AuthenticationUsecase) SignIn(ctx context.Context, user_dto *dto.UserLoginForm) (string, error) {
	if err := u.ValidateDTOStruct(user_dto); err != nil {
		return "", err
	}

	cred, err := u.repo.GetUserCredByEmail(ctx, user_dto.Email)
	if err != nil {
		if err == response_errors.ErrNoRecordFound {
			return "", response_errors.ErrIncorrectEmail
		}
		return "", err
	}

	if !utils.CheckPasswordHash(user_dto.Password, cred.Password) {
		return "", response_errors.ErrIncorrectPassword
	}

	user_id, err := u.repo.GetUserIDByEmail(ctx, user_dto.Email)
	if err != nil {
		return "", err
	}

	accessToken, err := middleware.GenerateAccessToken(user_id)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
