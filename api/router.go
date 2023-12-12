package api

import (
	"FantasticLife/config"
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
			apiV1.GET("userList", r.ChatSession.GetUserList)
			apiV1.POST("ChatSessionSendMessageAll", r.ChatSession.ChatSessionSendMessageAll)
			apiV1.POST("completionsTest", r.LLMBotServer.SpeakToBot_server)
			apiV1.POST("ChatInit", r.ChatSession.InitSession)
			apiV1.POST("completions", r.ChatSession.SendMessageToBot)
		}

	}
	return engine
}

func NewHttpServer(router Router, config *config.Config) *http.Server {
	//app := config.App
	//httpPort := app.HttpPort
	return &http.Server{
		Addr:    router.HostPort,
		Handler: router.Handler(),
	}
}
