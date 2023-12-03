package api

import (
	"FantasticLife/server"
	"FantasticLife/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
)

type Router struct {
	fx.In
	CorsMiddleware gin.HandlerFunc `name:"cors"`
	HostPort       string          `name:"hostPort"`
	GptBotServer   server.BOT
	ChatSession    services.TalkFunction
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
			apiV1.POST("completionsTest", r.GptBotServer.SpeakToBot_server)
			apiV1.POST("init", r.GptBotServer.InitBot)
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
