//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/miles0wu/meme-coin-api/internal/repository"
	"github.com/miles0wu/meme-coin-api/internal/repository/dao"
	"github.com/miles0wu/meme-coin-api/internal/service"
	"github.com/miles0wu/meme-coin-api/internal/web"
	"github.com/miles0wu/meme-coin-api/ioc"
)

func InitApp() *App {
	wire.Build(
		ioc.InitLogger,
		ioc.InitDB,
		dao.NewGormCoinDAO,
		repository.NewCachedCoinRepository,
		service.NewCoinService,
		web.NewCoinHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
		wire.Struct(new(App), "*"),
	)
	return &App{}
}
