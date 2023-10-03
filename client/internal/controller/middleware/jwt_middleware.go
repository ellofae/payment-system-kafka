package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	response_errors "github.com/ellofae/payment-system-kafka/client/internal/errors"
	"github.com/ellofae/payment-system-kafka/config"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	authorizationHeader = "Authorization"
)

type AccessTokenData struct {
	Expiry   int64
	IssuedAt int64
	UserID   string
	State    string
}

var jwtSecretKey string

func InitJWTSecretKey(cfg *config.Config) {
	jwtSecretKey = cfg.Authentication.JWTSecretKey
}

func GenerateAccessToken(user_id int) (string, error) {
	logger := logger.GetLogger()

	current_time := time.Now()
	jwt_token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id":   strconv.Itoa(user_id),
			"issued_at": current_time.Unix(),
			"expiry":    current_time.Add(time.Minute * 30).Unix(),
			"state":     "access_token",
		})

	access_token, err := jwt_token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		logger.Error("Unable to generate an access token", "error", err.Error())
		return "", errors.New("unable to generate an access token")
	}

	return access_token, nil
}

func ParseToken(tokenString string) (*AccessTokenData, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("unable to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expiry := claims["expiry"].(float64)
		issued_at := claims["issued_at"].(float64)

		return &AccessTokenData{
			Expiry:   int64(expiry),
			IssuedAt: int64(issued_at),
			UserID:   claims["user_id"].(string),
			State:    claims["state"].(string),
		}, nil
	} else {
		return nil, err
	}
}

func AuthenticateUser(c *gin.Context) {
	header_data := c.GetHeader(authorizationHeader)

	if header_data == "" {
		response_errors.NewErrorResponse(c, http.StatusUnauthorized, "Authorization header is empty")
		return
	}

	jwtString := strings.Split(header_data, "Bearer ")

	if len(jwtString) < 2 {
		response_errors.NewErrorResponse(c, http.StatusInternalServerError, "Must provide Authorization header with format `Bearer {token}`")
		return
	}

	token_claims, err := ParseToken(jwtString[1])
	if err != nil {
		response_errors.NewErrorResponse(c, http.StatusInternalServerError, "Incorrect access token provided")
		return
	}

	expiry := token_claims.Expiry
	if expiry < time.Now().Unix() {
		response_errors.NewErrorResponse(c, http.StatusInternalServerError, "Token expired")
		return
	}
	c.Set("user_id", token_claims.UserID)
	c.Next()
}
