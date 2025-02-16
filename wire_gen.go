// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/miles0wu/meme-coin-api/internal/repository"
	"github.com/miles0wu/meme-coin-api/internal/repository/dao"
	"github.com/miles0wu/meme-coin-api/internal/service"
	"github.com/miles0wu/meme-coin-api/internal/web"
	"github.com/miles0wu/meme-coin-api/ioc"
)

// Injectors from wire.go:

func InitApp() *App {
	v := ioc.InitGinMiddlewares()
	logger := ioc.InitLogger()
	db := ioc.InitDB(logger)
	coinDAO := dao.NewGormCoinDAO(db, logger)
	coinRepository := repository.NewCachedCoinRepository(coinDAO)
	coinService := service.NewCoinService(coinRepository)
	coinHandler := web.NewCoinHandler(coinService, logger)
	engine := ioc.InitWebServer(v, coinHandler)
	app := &App{
		server: engine,
	}
	return app
}
