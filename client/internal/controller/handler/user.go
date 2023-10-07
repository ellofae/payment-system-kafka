package handler

import (
	"github.com/ellofae/payment-system-kafka/client/internal/controller"
	"github.com/ellofae/payment-system-kafka/client/internal/domain"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
)

type UserHandler struct {
	logger      hclog.Logger
	userUsecase domain.IUserUsecase
}

func NewUserHandler(userUsecase domain.IUserUsecase) controller.IHandler {
	return &UserHandler{
		logger:      logger.GetLogger(),
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) Register(r *gin.Engine) {
	//	userGroup := r.Group("/user")

}
