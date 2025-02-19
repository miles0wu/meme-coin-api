//go:build wireinject

//go:generate wire

package main

import (
	"github.com/google/wire"
	"github.com/miles0wu/meme-coin-api/internal/repository"
	"github.com/miles0wu/meme-coin-api/internal/repository/cache"
	"github.com/miles0wu/meme-coin-api/internal/repository/dao"
	"github.com/miles0wu/meme-coin-api/internal/service"
	"github.com/miles0wu/meme-coin-api/internal/web"
	"github.com/miles0wu/meme-coin-api/ioc"
)

var thirdPartySet = wire.NewSet(
	ioc.InitLogger,
	ioc.InitDB,
	ioc.InitRedis,
)

func InitApp() *App {
	wire.Build(
		thirdPartySet,
		dao.NewGormCoinDAO,
		cache.NewRedisCoinCache,
		repository.NewCachedCoinRepository,
		service.NewCoinService,
		web.NewCoinHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
		wire.Struct(new(App), "*"),
	)
	return &App{}
}
