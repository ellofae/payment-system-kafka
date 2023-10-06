package handler

import (
	"fmt"
	"net/http"

	"github.com/ellofae/payment-system-kafka/client/internal/controller"
	"github.com/ellofae/payment-system-kafka/client/internal/controller/middleware"
	"github.com/ellofae/payment-system-kafka/client/internal/domain"
	"github.com/ellofae/payment-system-kafka/client/internal/dto"
	response_errors "github.com/ellofae/payment-system-kafka/client/internal/errors"
	"github.com/ellofae/payment-system-kafka/client/internal/repository"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
)

type AuthenticationHandler struct {
	logger      hclog.Logger
	authUsecase domain.IAuthenticationUsecase
}

func NewAuthenticationHandler(authUsecase domain.IAuthenticationUsecase) controller.IHandler {
	return &AuthenticationHandler{
		logger:      logger.GetLogger(),
		authUsecase: authUsecase,
	}
}

func (h *AuthenticationHandler) Register(r *gin.Engine) {
	authGroup := r.Group("/auth")

	authGroup.GET("/signup", h.handleRegistrationRendering)
	authGroup.GET("/signin", h.handleUserLoginRendering)

	authGroup.GET("/logout", middleware.AuthenticateUser, h.handleUserLogout)

	authGroup.POST("/signup", h.handleUserRegistration)
	authGroup.POST("/signin", h.handleUserLogin)
}

func (h *AuthenticationHandler) handleRegistrationRendering(c *gin.Context) {
	c.HTML(http.StatusOK, "registration.html", gin.H{
		"title": "register",
	})
}

func (h *AuthenticationHandler) handleUserLoginRendering(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "login",
	})
}

func (h *AuthenticationHandler) handleUserRegistration(c *gin.Context) {
	req := &dto.UserCreationForm{}

	err := c.ShouldBind(req)
	if err != nil {
		response_errors.NewHTTPResponse(c, http.StatusBadRequest, "Incorrect provided data", err)
		return
	}

	user_id, err := h.authUsecase.SignUp(c.Request.Context(), req)
	if err != nil {
		if err == response_errors.ErrNoRecordFound {
			response_errors.NewHTTPResponse(c, http.StatusNotFound, "No such record has been found", err)
			return
		}

		if err == response_errors.ErrAlreadyExists {
			response_errors.NewHTTPResponse(c, http.StatusBadRequest, "Provided email is already registered", err)
			return
		}

		response_errors.NewHTTPErrorResposne(c, http.StatusInternalServerError, "An error occured during user registration", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User has been registered",
		"user_id": user_id,
	})
}

func (h *AuthenticationHandler) handleUserLogin(c *gin.Context) {
	req := &dto.UserLoginForm{}

	err := c.ShouldBind(req)
	if err != nil {
		response_errors.NewHTTPResponse(c, http.StatusBadRequest, "Incorrect provided data", err)
		return
	}

	access_token, err := h.authUsecase.SignIn(c.Request.Context(), req)
	if err != nil {
		if err == response_errors.ErrIncorrectEmail {
			response_errors.NewHTTPResponse(c, http.StatusBadRequest, "Incorrect email was provided", err)
			return
		}

		if err == response_errors.ErrIncorrectPassword {
			response_errors.NewHTTPResponse(c, http.StatusBadRequest, "Incorrect password was provided", err)
			return
		}

		response_errors.NewHTTPErrorResposne(c, http.StatusInternalServerError, "An error occured while signing in", err)
		return
	}

	store := repository.SessionStorage()
	session, err := store.Get(c.Request, "session")
	if err != nil {
		response_errors.NewHTTPResponse(c, http.StatusInternalServerError, "Unable to get the session", err)
		return
	}

	session.Values["access_token"] = fmt.Sprintf("%s %s", "Bearer", access_token)
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		response_errors.NewHTTPResponse(c, http.StatusInternalServerError, "Unable to save session data", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": access_token,
	})
}

func (h *AuthenticationHandler) handleUserLogout(c *gin.Context) {
	store := repository.SessionStorage()

	session, err := store.Get(c.Request, "session")
	if err != nil {
		response_errors.NewHTTPResponse(c, http.StatusInternalServerError, "Unable to get the session", err)
		return
	}

	delete(session.Values, "access_token")
	if err = session.Save(c.Request, c.Writer); err != nil {
		response_errors.NewHTTPResponse(c, http.StatusInternalServerError, "Unable to save session data", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "You have logged out",
	})
}
