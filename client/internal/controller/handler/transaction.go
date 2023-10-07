package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ellofae/payment-system-kafka/client/internal/controller"
	"github.com/ellofae/payment-system-kafka/client/internal/controller/middleware"
	"github.com/ellofae/payment-system-kafka/client/internal/domain"
	"github.com/ellofae/payment-system-kafka/client/internal/dto"
	response_errors "github.com/ellofae/payment-system-kafka/client/internal/errors"
	"github.com/ellofae/payment-system-kafka/client/internal/repository"
	"github.com/ellofae/payment-system-kafka/client/internal/utils"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
)

type TransactionHandler struct {
	logger  hclog.Logger
	usecase domain.ITransactionUsecase
}

func NewTransactionHandler(usecase domain.ITransactionUsecase) controller.IHandler {
	return &TransactionHandler{
		logger:  logger.GetLogger(),
		usecase: usecase,
	}
}

func (h *TransactionHandler) Register(r *gin.Engine) {
	transactionGroup := r.Group("/transaction")

	transactionGroup.GET("/commit", middleware.AuthenticateUser, h.handleTransactionCommit)
	transactionGroup.POST("/send", middleware.AuthenticateUser, h.handleTransactionSending)
}

func (h *TransactionHandler) handleTransactionCommit(c *gin.Context) {
	c.HTML(http.StatusOK, "commit.html", nil)
}

func (h *TransactionHandler) handleTransactionSending(c *gin.Context) {
	store := repository.SessionStorage()
	session, err := store.Get(c.Request, "session")
	if err != nil {
		response_errors.NewHTTPResponse(c, http.StatusInternalServerError, "Unable to get the session", err)
		return
	}

	userID, ok := session.Values["user_id"]
	if !ok {
		response_errors.NewErrorResponse(c, http.StatusUnauthorized, "No user ID was recieved")
		return
	}

	userIDNumeric, _ := strconv.Atoi(userID.(string))
	transactionRequest := &dto.TransactionData{
		UserID:        userIDNumeric,
		TransactionID: utils.GenerateUniqueRandomString(14),
	}

	err = c.ShouldBind(transactionRequest)
	if err != nil {
		response_errors.NewHTTPErrorResposne(c, http.StatusBadRequest, "Incorrect provided data", err)
		return
	}

	if err := h.usecase.PlaceTransaction(transactionRequest); err != nil {
		response_errors.NewHTTPErrorResposne(c, http.StatusInternalServerError, "Field validation error occured", err)
		return
	}

	jsonData, err := json.Marshal(transactionRequest)
	if err != nil {
		response_errors.NewHTTPErrorResposne(c, http.StatusInternalServerError, "Unable to marshall data", err)
		return
	}

	response, err := http.Post("http://localhost:8000/transaction/place", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		response_errors.NewHTTPErrorResposne(c, http.StatusInternalServerError, "Unable to send JSON data to transaction server", err)
		return
	}
	defer response.Body.Close()

	if err := h.usecase.AttachTransaction(c.Request.Context(), transactionRequest); err != nil {
		response_errors.NewHTTPErrorResposne(c, http.StatusInternalServerError, "Unable to attach transaction to a user", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "transaction is placed",
	})
}
