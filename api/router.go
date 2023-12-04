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
	ZapLogger      gin.HandlerFunc `name:"zaplogger"`
	HostPort       string          `name:"hostPort"`
	LLMBotServer   server.LLMBOT
	ChatSession    services.ChatSessionService
}

func (r *Router) Handler() http.Handler {
	engine := gin.New()
	engine.Use(gin.Recovery(), r.CorsMiddleware)
	engine.Use(gin.Recovery(), r.ZapLogger)
	{
		apiV1 := engine.Group("api/v1")
		{
			apiV1.GET("/", func(c *gin.Context) {

				c.JSON(200, gin.H{
					"message": "Hello, World!",
				})
			})
			apiV1.POST("completionsTest", r.LLMBotServer.SpeakToBot_server)
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
