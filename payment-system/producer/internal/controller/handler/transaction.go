package handler

import (
	"net/http"
	"sync"

	"github.com/ellofae/payment-system-kafka/payment-system/producer/internal/controller"
	"github.com/ellofae/payment-system-kafka/payment-system/producer/internal/domain"
	"github.com/ellofae/payment-system-kafka/payment-system/producer/internal/domain/entity"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
)

type TransactionHandler struct {
	logger             hclog.Logger
	transactionUsecase domain.ITransactionUsecase
}

func NewTransactionHandler(transactionUsecase domain.ITransactionUsecase) controller.IHandler {
	return &TransactionHandler{
		logger:             logger.GetLogger(),
		transactionUsecase: transactionUsecase,
	}
}

func (th *TransactionHandler) Register(r *gin.Engine) {
	transactionGroup := r.Group("/transaction")

	transactionGroup.POST("/place", th.handlePlaceTransaction)
}

func (th *TransactionHandler) handlePlaceTransaction(c *gin.Context) {
	transactionData := &entity.TransactionData{}

	err := c.ShouldBindJSON(transactionData)
	if err != nil {
		th.logger.Error("Unable to bind JSON transaction data", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to bind JSON transaction data",
			"error":   err.Error(),
		})
		return
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		err = th.transactionUsecase.PlaceTransaction(c.Request.Context(), transactionData)
		if err != nil {
			th.logger.Error("Internal error occured", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Some internal error occured",
				"error":   err.Error(),
			})
			return
		}

	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction has been produced",
	})

	wg.Wait()
}
