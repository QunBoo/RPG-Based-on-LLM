package servicesimpl

import (
	"FantasticLife/server"
	"FantasticLife/server/serverimpl/WebSocket"
	"FantasticLife/services"
	"FantasticLife/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type ChatSessionServiceImpl struct {
	ChatSessionList map[string]*ChatSession
	ClientManager   *WebSocket.ClientManager
	logger          *zap.Logger
}
type ChatSession struct {
	ChatSessionId string
	ChatHistory   []map[string]string
	LLMBOTInter   server.LLMBOT
}

// 和Bot的交互功能
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
func (s *ChatSessionServiceImpl) GetUserList(c *gin.Context) {
	appIdStr := c.Query("appId")
	appIdUint64, _ := strconv.ParseInt(appIdStr, 10, 32)
	appId := uint32(appIdUint64)

	//fmt.Println("http_request 查看全部在线用户", appId)
	s.logger.Info("http_request 查看全部在线用户", zap.Uint32("appId", appId))

	data := make(map[string]interface{})

	//userList := WebSocket.ClientManager.GetUserList(appId)
	userList := s.ClientManager.GetUserList(appId)
	data["userList"] = userList
	data["userCount"] = len(userList)
	c.JSON(http.StatusOK, gin.H{
		"message": data,
	})
}
func (s *ChatSessionServiceImpl) ChatSessionSendMessageAll(c *gin.Context) {
	// 获取参数
	appIdStr := c.PostForm("appId")
	userId := c.PostForm("userId")
	msgId := c.PostForm("msgId")
	message := c.PostForm("message")
	appIdUint64, _ := strconv.ParseInt(appIdStr, 10, 32)
	appId := uint32(appIdUint64)

	fmt.Println("http_request 给全体用户发送消息", appIdStr, userId, msgId, message)

	data := make(map[string]interface{})
	//if cache.SeqDuplicates(msgId) {
	//	fmt.Println("给用户发送消息 重复提交:", msgId)
	//	controllers.Response(c, common.OK, "", data)
	//
	//	return
	//}

	sendResults, err := s.ClientManager.SendUserMessageAll(appId, userId, msgId, utils.MessageCmdMsg, message)
	if err != nil {
		s.logger.Error("发送消息报错", zap.Error(err))
	}

	data["sendResults"] = sendResults

	c.JSON(utils.OK, gin.H{
		"message": data,
	})
}

func NewChatSession(llmbot server.LLMBOT) *ChatSession {
	return &ChatSession{
		ChatSessionId: "Default",
		ChatHistory:   nil,
		LLMBOTInter:   llmbot,
	}
}

func NewChatSessionService(zapLogger *zap.Logger, defaultSesstion *ChatSession, ClientManager *WebSocket.ClientManager) services.ChatSessionService {
	SessionList := make(map[string]*ChatSession)
	SessionList["Default"] = defaultSesstion
	CSService := ChatSessionServiceImpl{
		ChatSessionList: SessionList,
		ClientManager:   ClientManager,
		logger:          zapLogger,
	}

	return &CSService
}
