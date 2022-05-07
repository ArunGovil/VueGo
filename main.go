//go:generate go run -tags=dev assets_generate.go

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArunGovil/VueGo/assets"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

var Version = "Welcome"

func main() {
	var (
		err               error
		nodeContextCancel context.CancelFunc
	)

	logger := logrus.New().WithField("who", "Example")

	httpServer := echo.New()
	httpServer.Use(middleware.CORS())

	httpServer.GET("/*", echo.WrapHandler(http.FileServer(assets.Assets)))
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

	if Version == "Welcome" {
		_, nodeContextCancel = context.WithCancel(context.Background())
	}

	go func() {
		logger.Info("Starting Node development server...")
		var cmd *exec.Cmd
		var err error

		if cmd, err = StartClientApp(); err != nil {
			logger.WithError(err).Fatal("Error starting Node server..")
		}

		cmd.Wait()
		logger.Info("Stopping Node development server...")
	}()

	<-quit

	if Version == "Welcome" {
		nodeContextCancel()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = httpServer.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("There was an error shutting down the server")
	}

	logger.Info("Application stopped")

}

func StartClientApp() (*exec.Cmd, error) {
	var err error

	cmd := exec.Command("yarn", "serve")
	cmd.Dir = "./frontend"
	cmd.Stdout = os.Stdout

	if err = cmd.Start(); err != nil {
		return cmd, fmt.Errorf("Error starting Yarn: %W", err)
	}
	return cmd, nil
}
