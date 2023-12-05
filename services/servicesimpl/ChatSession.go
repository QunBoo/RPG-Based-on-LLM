package servicesimpl

import (
	"FantasticLife/server"
	"FantasticLife/services"
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
	ChatHistory   []map[string]string
	LLMBOTInter   server.LLMBOT
}

// TODO: 管理ChatSession和Service的关系，以及talkfunction接口
func (s *ChatSessionServiceImpl) SendMessageToBot(c *gin.Context) {
	var input struct {
		SessionId string `json:"session_id"`
		Messages  string `json:"messages"`
	}
	// 从请求中读取 JSON 数据
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	TempChatSessionP := s.ChatSessionList[input.SessionId]
	TempChatSessionP.ChatHistory = append(TempChatSessionP.ChatHistory, map[string]string{
		"role":    "user",
		"content": input.Messages,
	})
	respMessage := TempChatSessionP.LLMBOTInter.SpeakToBot(c, TempChatSessionP.ChatHistory)
	TempChatSessionP.ChatHistory = append(TempChatSessionP.ChatHistory, map[string]string{
		"role":    "assistant",
		"content": respMessage,
	})
	s.logger.Info("SendMessageToBot", zap.Any("ChatHistory", TempChatSessionP.ChatHistory))
	c.JSON(http.StatusOK, gin.H{
		"message": respMessage,
	})
	//s.BotInter.SpeakToBot(c, modifiedMessage)
}
func (s *ChatSessionServiceImpl) InitSession(c *gin.Context) {
	var input struct {
		SessionId string `json:"session_id"`
	}
	// 从请求中读取 JSON 数据
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 初始化，Session置零
	TempChatSessionP := s.ChatSessionList[input.SessionId]
	TempChatSessionP.ChatHistory = nil
	c.JSON(http.StatusOK, gin.H{
		"message": "InitSession Success!",
	})
}

func NewChatSession(llmbot server.LLMBOT) *ChatSession {
	return &ChatSession{
		ChatSessionId: "Default",
		ChatHistory:   nil,
		LLMBOTInter:   llmbot,
	}
}

// TODO：初始化ChatSession，研究fx怎么用
func NewChatSessionService(zapLogger *zap.Logger, defaultSesstion *ChatSession) services.ChatSessionService {
	SessionList := make(map[string]*ChatSession)
	SessionList["Default"] = defaultSesstion
	CSService := ChatSessionServiceImpl{
		ChatSessionList: SessionList,
		logger:          zapLogger,
	}
	//TkFunc := services.ChatSessionService(&CSService)
	return &CSService
}
