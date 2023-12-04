package serverimpl

import (
	"FantasticLife/config"
	"FantasticLife/server"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type LLMEntity struct {
	Id     string `json:"id"`
	Conn   server.LLMTransceiver
	logger *zap.Logger
}

type BaiChuanConn struct {
	Key      string
	EndPoint string
	logger   *zap.Logger
}
type RequestData struct {
	Messages struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}
type ResponseData struct {
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// TODO: 1. 实现对话接口
func (co *BaiChuanConn) SpeakToLLM(c *gin.Context, messageMapSlice []map[string]string) (respMessage string) {
	url := co.EndPoint
	api_key := co.Key
	reqBody := map[string]interface{}{
		"model":    "Baichuan2",
		"messages": messageMapSlice,
		"stream":   false,
	}
	jsonData, err := json.Marshal(reqBody)
	//fmt.Println(url, api_key, string(jsonData))
	client := &http.Client{}                                            // 创建客户端
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData)) // 创建请求

	req.Header.Add("Content-Type", "application/json") // 添加请求头
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api_key))
	res, err := client.Do(req) // 发送请求
	if err != nil {
		panic(err)
		co.logger.Error("发送请求失败", zap.Error(err))
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
			co.logger.Error("关闭请求失败", zap.Error(err))
		}
	}(res.Body) // 关闭请求
	// 返回消息
	//fmt.Println(res.StatusCode)
	body, err := io.ReadAll(res.Body)
	if res.StatusCode == 200 {
		co.logger.Info("响应SpeakToBot:", zap.String("body", string(body)))
	} else {
		co.logger.Warn("响应SpeakToBot:", zap.String("body", string(body)))
	}
	var respData ResponseData
	err = json.Unmarshal(body, &respData)
	if err != nil {
		panic(err)
		co.logger.Error("解析响应失败", zap.Error(err))
	}
	c.JSON(http.StatusOK, gin.H{
		"message": respData.Choices[0].Message.Content,
	})
	respMessage = respData.Choices[0].Message.Content
	return respMessage
}

func (b *LLMEntity) SpeakToBot(c *gin.Context, messageMapSlice []map[string]string) (respMessage string) {
	//发送消息
	respMessage = b.Conn.SpeakToLLM(c, messageMapSlice)
	// 返回消息
	c.JSON(http.StatusOK, gin.H{
		"message": "SpeakToBot Success!",
	})
	return respMessage
}

func (b *LLMEntity) SpeakToBot_server(c *gin.Context) {
	//消息处理
	var requestData RequestData
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	messageMap := map[string]string{
		"role":    requestData.Messages.Role,
		"content": requestData.Messages.Content,
	}
	messageMapSlice := []map[string]string{messageMap}
	//发送消息
	respMessage := b.Conn.SpeakToLLM(c, messageMapSlice)
	// 返回消息
	c.JSON(http.StatusOK, gin.H{
		"message": respMessage,
	})
}

func NewLLMBOT(pConn server.LLMTransceiver, zapLogger *zap.Logger) (server.LLMBOT, error) {
	LLMe := LLMEntity{
		Id:     "Default",
		Conn:   pConn,
		logger: zapLogger,
	}
	return &LLMe, nil
}
func NewLLMTransceiver(config *config.Config, zapLogger *zap.Logger) server.LLMTransceiver {
	gptLArk := config.GptLark
	LLMName := "BaiChuan"
	if LLMName == "BaiChuan" {
		bcConn := BaiChuanConn{
			Key:      gptLArk.Key,
			EndPoint: gptLArk.EndPoint,
			logger:   zapLogger,
		}
		return &bcConn
	}
	return nil
}
