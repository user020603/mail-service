package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"thanhnt208/mail-service/api/routes"
	"thanhnt208/mail-service/config"
	"thanhnt208/mail-service/external/client"
	"thanhnt208/mail-service/internal/delivery/rest"
	"thanhnt208/mail-service/internal/service"
	"thanhnt208/mail-service/pkg/logger"
	"time"
)

func main() {
	cfg := config.LoadConfig()
	if cfg == nil {
		panic("Failed to load configuration")
	}

	logger, err := logger.NewLogger(cfg.LogLevel, cfg.LogFile)
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	grpcConn, err := client.StartGrpcClient()
	if err != nil {
		logger.Error("Failed to start gRPC client", "error", err)
		panic("Failed to start gRPC client: " + err.Error())
	}

	getContainerInfoClient := client.NewGetContainerInfoClient(grpcConn, logger)
	if getContainerInfoClient == nil {
		logger.Error("Failed to create GetContainerInfoClient")
		panic("Failed to create GetContainerInfoClient")
	}

	mailService := service.NewMailService(getContainerInfoClient, logger)
	mailHandler := rest.NewMailHandler(mailService)

	r := routes.SetupMailRoutes(mailHandler)

	port := cfg.RestPort
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		logger.Info("Starting REST server", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start REST server", "error", err)
			panic("Failed to start REST server: " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server gracefully...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Fatal("REST server forced to shutdown:", "error", err)
	}

	logger.Info("Server exited gracefully")
}
