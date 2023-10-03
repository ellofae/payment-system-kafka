package controller

import "github.com/gin-gonic/gin"

type IHandler interface {
	Register(*gin.Engine)
}
