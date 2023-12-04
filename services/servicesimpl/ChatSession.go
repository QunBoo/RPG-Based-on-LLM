package servicesimpl

import (
	"FantasticLife/server"
	"FantasticLife/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type ChatSessionServiceImpl struct {
	ChatSessionList map[string]*ChatSession
	logger          *zap.Logger
}
type ChatSession struct {
	ChatSessionId string
	LLMBOTInter   *server.LLMBOT
}

// TODO: 管理ChatSession和Service的关系，以及talkfunction接口
func (s *ChatSessionServiceImpl) SendMessageToBot(c *gin.Context) {
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
func (s *ChatSessionServiceImpl) InitSession(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func NewChatSession(llmbot server.LLMBOT) *ChatSession {
	return &ChatSession{
		ChatSessionId: "Default",
		LLMBOTInter:   &llmbot,
	}
}

// TODO：初始化ChatSession，研究fx怎么用
func NewChatSessionService(zapLogger *zap.Logger, defaultSesstion *ChatSession) services.ChatSessionService {
	SessionList := make(map[string]*ChatSession)
	SessionList["Default"] = defaultSesstion
	CSService := ChatSessionServiceImpl{
		ChatSessionList: make(map[string]*ChatSession),
		logger:          zapLogger,
	}
	//TkFunc := services.ChatSessionService(&CSService)
	return &CSService
}
