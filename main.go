package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

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
	app.server.Run(":8080")
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
