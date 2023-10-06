package handler

import (
	"net/http"

	"github.com/ellofae/payment-system-kafka/client/internal/controller"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
)

type InformationHandler struct {
	logger hclog.Logger
}

func NewInformationHandler() controller.IHandler {
	return &InformationHandler{
		logger: logger.GetLogger(),
	}
}

func (h *InformationHandler) Register(r *gin.Engine) {
	r.GET("/index", h.handleIndexPage)
}

func (h *InformationHandler) handleIndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
