package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ellofae/payment-system-kafka/client/internal/controller/handler"
	"github.com/ellofae/payment-system-kafka/client/internal/controller/middleware"
	"github.com/ellofae/payment-system-kafka/client/internal/domain/usecase"
	"github.com/ellofae/payment-system-kafka/client/internal/repository"
	"github.com/ellofae/payment-system-kafka/config"
	"github.com/ellofae/payment-system-kafka/pkg/encryption"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/ellofae/payment-system-kafka/pkg/postgres"
	"github.com/gin-gonic/gin"
)

func Run() {
	log := logger.GetLogger()
	cfg := config.ParseConfig(config.ConfigureViper())
	ctx := context.Background()

	encryption.InitializeEncryptionKey(cfg)
	repository.InitSessionStorage(cfg)

	connPool := postgres.OpenPoolConnection(ctx, cfg)
	if err := connPool.Ping(ctx); err != nil {
		log.Error("Unable to ping the database connection", "error", err.Error())
		os.Exit(1)
	}

	postgres.RunMigrationsUp(ctx, cfg)

	middleware.InitJWTSecretKey(cfg)

	storage := repository.NewStorage(connPool)
	router := initRouter(storage)

	srv := initHTTPServer(router, cfg)

	startServer(ctx, srv)

}

func initRouter(storage *repository.Storage) *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("client/web/templates/*.html")

	r.Static("/assets", "./client/web/assets")

	authenticationRepository := repository.NewAuthenticationRepository(storage)
	authenticationUsecase := usecase.NewAuthenticationUsecase(authenticationRepository)
	authenticationHandler := handler.NewAuthenticationHandler(authenticationUsecase)

	transactionRepository := repository.NewTransactionRepository(storage)
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepository)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	userRepository := repository.NewUserRepository(storage)
	userUsecase := usecase.NewUserUsecase(userRepository)
	userHandler := handler.NewUserHandler(userUsecase)

	informationHandler := handler.NewInformationHandler()

	authenticationHandler.Register(r)
	informationHandler.Register(r)
	transactionHandler.Register(r)
	userHandler.Register(r)

	return r
}

func initHTTPServer(router *gin.Engine, cfg *config.Config) http.Server {
	readTimeoutSecondsCount, _ := strconv.Atoi(cfg.ClientServer.ReadTimeout)
	writeTimeoutSecondsCount, _ := strconv.Atoi(cfg.ClientServer.WriteTimeout)
	idleTimeoutSecondsCount, _ := strconv.Atoi(cfg.ClientServer.IdleTimeout)

	bindAddr := cfg.ClientServer.BindAddr

	srv := http.Server{
		Addr:         bindAddr,
		Handler:      router,
		ReadTimeout:  time.Duration(readTimeoutSecondsCount) * time.Second,
		WriteTimeout: time.Duration(writeTimeoutSecondsCount) * time.Second,
		IdleTimeout:  time.Duration(idleTimeoutSecondsCount) * time.Second,
	}

	return srv
}

func startServer(ctx context.Context, srv http.Server) {
	log := logger.GetLogger()

	go func() {
		log.Info("Starting server...")
		err := srv.ListenAndServe()
		if err != nil {
			log.Error("Server was stopped", "error", err.Error())
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	signal := <-sigChan
	log.Info(fmt.Sprintf("Signal has been caught: %v", signal))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
}
