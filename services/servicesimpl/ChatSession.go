package servicesimpl

import (
	"FantasticLife/server"
	"FantasticLife/server/serverimpl/WebSocket"
	"FantasticLife/services"
	"FantasticLife/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type ChatSessionServiceImpl struct {
	ChatSessionList map[string]*ChatSession
	ClientManager   *WebSocket.ClientManager
	logger          *zap.Logger
	Producer        sarama.SyncProducer
	ConsumerGroup   sarama.ConsumerGroup
}
type ChatSession struct {
	ChatSessionId string
	ChatHistory   []map[string]string
	LLMBOTInter   server.LLMBOT
}

func (s *ChatSessionServiceImpl) ChatSendMessageMQ(c *gin.Context) {
	var input struct {
		SessionId string `json:"sessionId"`
		Messages  string `json:"messages"`
	}
	// 从请求中读取 JSON 数据
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//如果不存在SessionId，则输出错误，返回，打印当前的ChatSessionList
	if _, ok := s.ChatSessionList[input.SessionId]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SessionId not exist!", "SessionId": input.SessionId})
		s.logger.Info("SendMessageToBot", zap.Any("ChatSessionList", s.ChatSessionList))
		return
	}

	//	将Message发送到MQ
	msg := &sarama.ProducerMessage{
		Topic: "chat",
		Value: sarama.StringEncoder(input.Messages),
	}
	partition, offset, err := s.Producer.SendMessage(msg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.logger.Info("ChatSendMessageMQ", zap.Any("log", fmt.Sprintf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", "chat", partition, offset)))

}
func (s *ChatSessionServiceImpl) ChatGetMQMessage() {
	ctx := context.Background()
	handler := ExampleConsumerGroupHandler{
		CSImplP: s,
	}
	group := s.ConsumerGroup
	// 消费者组循环，确保在消费者出错时可以重新加入
	for {
		if err := group.Consume(ctx, []string{"chat"}, handler); err != nil {
			log.Printf("Error from consumer: %v", err)
		}
	}

}

// 和Bot的交互功能
func (s *ChatSessionServiceImpl) SendMessageToBot(c *gin.Context) {
	var input struct {
		SessionId string `json:"sessionId"`
		Messages  string `json:"messages"`
	}
	// 从请求中读取 JSON 数据
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//如果不存在SessionId，则输出错误，返回，打印当前的ChatSessionList
	if _, ok := s.ChatSessionList[input.SessionId]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SessionId not exist!", "SessionId": input.SessionId})
		s.logger.Info("SendMessageToBot", zap.Any("ChatSessionList", s.ChatSessionList))
		return
	}
	TempChatSessionP := s.ChatSessionList[input.SessionId]
	TempChatSessionP.ChatHistory = append(TempChatSessionP.ChatHistory, map[string]string{
		"role":    "user",
		"content": input.Messages,
	})
	//s.logger.Info("SendMessageToBot0", zap.Any("ChatHistory", TempChatSessionP.ChatHistory))
	respMessage := TempChatSessionP.LLMBOTInter.SpeakToBot(c, TempChatSessionP.ChatHistory)
	TempChatSessionP.ChatHistory = append(TempChatSessionP.ChatHistory, map[string]string{
		"role":    "assistant",
		"content": respMessage,
	})
	s.logger.Info("SendMessageToBot", zap.Any("ChatHistory", TempChatSessionP.ChatHistory))
	c.JSON(http.StatusOK, gin.H{
		"message": respMessage,
	})
	//通过调用sendMessageAll函数，将消息发送给所有用户
	userId := "小助手"
	msgId := "msgId"
	_, err := s.ClientManager.SendUserMessageAll(101, userId, msgId, utils.MessageCmdMsg, respMessage)
	if err != nil {
		return
	}
}
func (s *ChatSessionServiceImpl) InitSession(c *gin.Context) {
	var input struct {
		SessionId string `json:"sessionId"`
	}
	// 从请求中读取 JSON 数据
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//如果不存在SessionId，则输出错误，返回，打印当前的ChatSessionList
	if _, ok := s.ChatSessionList[input.SessionId]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SessionId not exist!", "SessionId": input.SessionId})
		s.logger.Info("InitSession", zap.Any("ChatSessionList", s.ChatSessionList))
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
	Response(c, utils.OK, "", data)
}

func Response(c *gin.Context, code uint32, msg string, data map[string]interface{}) {
	message := utils.ResponseMsg(code, msg, data)

	// 允许跨域
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Origin", "*") // 这是允许访问所有域
	c.Header("Access-Control-Allow-Methods",
		"POST, GET, OPTIONS, PUT, DELETE,UPDATE") // 服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
	c.Header("Access-Control-Allow-Headers",
		"Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
	c.Header("Access-Control-Expose-Headers",
		"Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
	c.Header("Access-Control-Allow-Credentials",
		"true") //  跨域请求是否需要带cookie信息 默认设置为true
	c.Set("content-type",
		"application/json") // 设置返回格式是json

	c.JSON(http.StatusOK, message)

	return
}

// 登录
func (s *ChatSessionServiceImpl) Login(c *gin.Context) {
	// 需要username, password
	var userInfo struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	// 从请求中读取 JSON 数据
	if err := c.BindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 生成password_hash
	passwordHash := utils.GeneratePasswordHash(userInfo.Password)
	// 查询数据库
	dbConn := s.ClientManager.MysqlCli
	var userDB struct {
		Username     string `json:"username"`
		PasswordHash string `json:"password_hash"`
	}
	result := dbConn.Table("users").Where("username = ?", userInfo.Username).First(&userDB)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}
	if userDB.PasswordHash != passwordHash {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password Error!"})
		return
	}
	// 生成token
	token, err := utils.GenerateJWT(userInfo.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login Success!",
		"token":   token,
	})
}

// 注册
func (s *ChatSessionServiceImpl) SignUp(c *gin.Context) {
	// 需要username, password, email
	var userInfo struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	// 从请求中读取 JSON 数据
	if err := c.BindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//// 获取当前数据库user_id的最大值
	//dbConn := s.ClientManager.MysqlCli
	//var maxUserId uint32
	//dbConn.Table("user").Select("max(user_id)").Scan(&maxUserId)
	//
	//// 生成新的user_id
	//curUserId := maxUserId + 1
	// 生成password_hash
	passwordHash := utils.GeneratePasswordHash(userInfo.Password)
	type UserDB struct {
		Username     string `json:"username"`
		PasswordHash string `json:"password_hash"`
		Email        string `json:"email"`
	}
	//	// 插入数据库
	dbConn := s.ClientManager.MysqlCli
	result := dbConn.Table("users").Create(&UserDB{Username: userInfo.Username, PasswordHash: passwordHash, Email: userInfo.Email})
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "SignUp Success!",
	})
	return

}

func (s *ChatSessionServiceImpl) ChatSessionSendMessageAll(c *gin.Context) {
	var wg sync.WaitGroup
	var gofunc func(msg string)
	gofunc = func(msg string) {
		defer wg.Done()
		SessionId := "Default"
		Messages := msg
		var inputData struct {
			SessionId string `json:"sessionId"`
			Messages  string `json:"Messages"`
		}
		inputData.SessionId = SessionId
		inputData.Messages = Messages
		data, _ := json.Marshal(inputData)

		// 发送请求, 注意当部署到服务器上时，需要将ip地址改为服务器的外网ip地址
		resp, err := http.Post("http://127.0.0.1:8080/api/v1/completions", "application/json", bytes.NewBuffer(data))
		if err != nil {
			fmt.Println("Error:", err)
		}
		defer resp.Body.Close()

		// 读取响应
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error:", err)
		}

		fmt.Println("<<<< resp From MQ", string(body))
	}
	// 获取参数
	appIdStr := c.PostForm("appId")
	userId := c.PostForm("userId")
	msgId := c.PostForm("msgId")
	message := c.PostForm("message")
	appIdUint64, _ := strconv.ParseInt(appIdStr, 10, 32)
	appId := uint32(appIdUint64)

	fmt.Println("http_request 给全体用户发送消息", appIdStr, userId, msgId, message)

	data := make(map[string]interface{})

	//如果message中包括"@小助手"，则自动调用SendMessageToBot函数
	if strings.Contains(message, "@小助手") {
		wg.Add(1)
		go gofunc(message)
	}

	sendResults, err := s.ClientManager.SendUserMessageAll(appId, userId, msgId, utils.MessageCmdMsg, message)
	if err != nil {
		s.logger.Error("发送消息报错", zap.Error(err))
	}

	data["sendResults"] = sendResults
	//wg.Wait()

	Response(c, utils.OK, "", data)
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
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // 确保消息被写入所有副本后才认为是成功的
	config.Producer.Retry.Max = 5                    // 最大重试次数
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{"8.141.81.87:9092"}, config)
	if err != nil {
		zapLogger.Fatal("Failed to start Sarama producer:", zap.Error(err))

	}
	configCon := sarama.NewConfig()
	configCon.Version = sarama.V2_0_0_0 // 确保版本兼容
	configCon.Consumer.Return.Errors = true
	configCon.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	configCon.Consumer.Offsets.Initial = sarama.OffsetOldest // 从最早的消息开始消费

	ConsumerGroup, err := sarama.NewConsumerGroup([]string{"8.141.81.87:9092"}, "your_consumer_group_id", config)
	if err != nil {
		zapLogger.Fatal("Error creating consumer group:", zap.Error(err))

	}
	CSService := ChatSessionServiceImpl{
		ChatSessionList: SessionList,
		ClientManager:   ClientManager,
		logger:          zapLogger,
		Producer:        producer,
		ConsumerGroup:   ConsumerGroup,
	}

	go CSService.ChatGetMQMessage() //开启一个协程，监听MQ消息

	return &CSService
}

// ExampleConsumerGroupHandler represents a Sarama consumer group consumer
type ExampleConsumerGroupHandler struct {
	CSImplP *ChatSessionServiceImpl
}

func (ExampleConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ExampleConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h ExampleConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)
		sess.MarkMessage(msg, "") // 标记消息为已处理

		SessionId := "Default"
		Messages := string(msg.Value)
		var inputData struct {
			SessionId string `json:"sessionId"`
			Messages  string `json:"Messages"`
		}
		inputData.SessionId = SessionId
		inputData.Messages = Messages
		data, _ := json.Marshal(inputData)

		// 发送请求, 注意当部署到服务器上时，需要将ip地址改为服务器的外网ip地址
		resp, err := http.Post("http://127.0.0.1:8080/api/v1/completions", "application/json", bytes.NewBuffer(data))
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		defer resp.Body.Close()

		// 读取响应
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		fmt.Println("<<<< resp From MQ", string(body))
	}
	return nil
}
