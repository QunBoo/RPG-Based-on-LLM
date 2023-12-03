package servicesimpl

import (
	"FantasticLife/server"
	"FantasticLife/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type ChatSessionService struct {
	ChatSessionList map[string]*ChatSession
	logger          *zap.Logger
}
type ChatSession struct {
	ChatSessionId string
	BotInter      server.BOT
}

// TODO: 管理ChatSession和Service的关系，以及talkfunction接口
func (s *ChatSessionService) SendMessageToBot(c *gin.Context) {
	var input struct {
		Messages string `json:"messages"`
	}
	// 从请求中读取 JSON 数据
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 修改 JSON 数据格式
	modifiedMessage := map[string]string{
		"role":    "user",
		"content": input.Messages,
	}
	fmt.Println(modifiedMessage)
	//s.BotInter.SpeakToBot(c, modifiedMessage)
}
func (s *ChatSessionService) InitSession(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func NewChatSessionService(zapLogger *zap.Logger) services.TalkFunction {
	CSService := ChatSessionService{
		logger: zapLogger,
	}
	//TkFunc := services.TalkFunction(&CSService)
	return &CSService
}
