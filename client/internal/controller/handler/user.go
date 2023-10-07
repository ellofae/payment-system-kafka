package handler

import (
	"net/http"
	"strconv"

	"github.com/ellofae/payment-system-kafka/client/internal/controller"
	"github.com/ellofae/payment-system-kafka/client/internal/domain"
	response_errors "github.com/ellofae/payment-system-kafka/client/internal/errors"
	"github.com/ellofae/payment-system-kafka/client/internal/repository"
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
	userGroup := r.Group("/user")

	userGroup.GET("/", h.handleUserRecords)
}

func (h *UserHandler) handleUserRecords(c *gin.Context) {
	store := repository.SessionStorage()
	session, err := store.Get(c.Request, "session")
	if err != nil {
		response_errors.NewHTTPResponse(c, http.StatusInternalServerError, "Unable to get the session", err)
		return
	}

	userID, ok := session.Values["user_id"]
	if !ok {
		response_errors.NewHTTPErrorResposne(c, http.StatusUnauthorized, "No user ID was recieved", err)
		return
	}

	userIDNumeric, _ := strconv.Atoi(userID.(string))

	transactions, err := h.userUsecase.GetUserTransaction(c.Request.Context(), userIDNumeric)
	if err != nil {
		if err != response_errors.ErrNoRecordFound {
			response_errors.NewHTTPErrorResposne(c, http.StatusInternalServerError, "Failed to get users' transactions", err)
			return
		}
	}

	userData, err := h.userUsecase.GetUserData(c.Request.Context(), userIDNumeric)
	if err != nil {
		response_errors.NewHTTPErrorResposne(c, http.StatusInternalServerError, "Failed to get users' data", err)
		return
	}

	c.HTML(http.StatusOK, "records.html", gin.H{
		"user":         userData,
		"transactions": transactions,
	})
}
