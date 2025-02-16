package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/miles0wu/meme-coin-api/api/docs"
	"github.com/miles0wu/meme-coin-api/internal/web"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, coinHdl *web.CoinHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)

	coinHdl.RegisterRoutes(server)
	server.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return server
}

func InitGinMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			AllowCredentials: true,
			AllowHeaders:     []string{"Content-Type"},
			AllowOriginFunc: func(origin string) bool {
				return strings.HasPrefix(origin, "http://localhost")
			},
			MaxAge: 12 * time.Hour,
		}),
	}
}
