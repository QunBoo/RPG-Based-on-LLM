package api

import (
	"FantasticLife/server"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
)

type Router struct {
	fx.In
	CorsMiddleware gin.HandlerFunc `name:"cors"`
	HostPort       string          `name:"hostPort"`
	GptBotServer   server.BOT
}

func (r *Router) Handler() http.Handler {
	engine := gin.New()
	engine.Use(gin.Recovery(), r.CorsMiddleware)
	{
		apiV1 := engine.Group("api/v1")
		{
			apiV1.GET("/", func(c *gin.Context) {

				c.JSON(200, gin.H{
					"message": "Hello, World!",
				})
			})
			apiV1.POST("bot-chat", r.GptBotServer.BOTChat)
		}

	}
	return engine
}

func NewHttpServer(router Router) *http.Server {
	return &http.Server{
		Addr:    router.HostPort,
		Handler: router.Handler(),
	}
}
