package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// App
// @title MemeCoins
// @version 0.1.0
// @description Meme coins api swagger
//
// @contact.name miles0wu
// @contact.email miles4w701@gmail.com
//
// @license.name GPL-3.0
// @license.url https://spdx.org/licenses/GPL-3.0-only.html
//
// @BasePath /
type App struct {
	server *gin.Engine
}

func main() {
	initViper()
	initLogger()
	app := InitApp()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.server,
	}

	// start server
	go func() {
		zap.L().Info("Server starting", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("listen: ", zap.Error(err))
		}
	}()

	// listening signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutting down server...")

	// shutdown timeout 5 secs
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server forced to shutdown", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}

func initViper() {
	cfile := pflag.String("config", "config/config.yaml", "config file")
	pflag.Parse()

	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cfile)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
