package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

var Version = "Welcome to VueGo!"

func main() {
	var (
		err error
	)

	logger := logrus.New().WithField("who", "Example")
	httpServer := echo.New()
	httpServer.Use(middleware.CORS())

	httpServer.GET("/api/version", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, Version)
	})

	go func() {
		var err error

		logger.WithFields(logrus.Fields{
			"serverVersion": Version,
		}).Infof("Starting application")

		err = httpServer.Start("0.0.0.0:8080")

		if err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Unable to start application")
		} else {
			logger.Info("Shutting down the server...")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = httpServer.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("There was an error shutting down the server")
	}

	logger.Info("Application stopped")
}
