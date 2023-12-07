package serverimpl

import (
	"FantasticLife/config"
	"FantasticLife/server"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
type WXConn struct {
	Key       string
	AppSecret string
	EndPoint  string
	logger    *zap.Logger
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

type ResponseDataForWX struct {
	Result string `json:"result"`
}

func (co WXConn) SpeakToLLM(c *gin.Context, messageMapSlice []map[string]string) (respMessage string) {
	apiKey := co.Key
	secretKey := co.AppSecret
	tokenURL := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s", apiKey, secretKey)
	resp, err := http.Get(tokenURL)
	if err != nil {
		fmt.Println("Failed to get access token:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read access token response:", err)
		return
	}
	type AccessToken struct {
		AccessToken string `json:"access_token"`
	}
	var accessToken1 AccessToken
	err = json.Unmarshal(body, &accessToken1)
	if err != nil {
		fmt.Println("Failed to unmarshal access token:", err)
		return
	}
	accessToken := accessToken1.AccessToken

	// 步骤二：调用文心一言API发送信息
	apiEndpoint := fmt.Sprintf("https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions_pro?access_token=%s", accessToken)

	// 创建一个新的结构体，用于生成期望的 JSON 格式
	type Messages struct {
		Messages []map[string]string `json:"messages"`
	}
	// 将 messageMapSlice 放入 Messages 结构体
	messages := Messages{Messages: messageMapSlice}
	// MarshalIndent 用于生成格式化的 JSON 字符串
	jsonData, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// 将 JSON 字符串转换为 []byte 类型
	byteSlice := []byte(jsonData)
	resp, err = http.Post(apiEndpoint, "application/json", bytes.NewBuffer(byteSlice))
	if err != nil {
		fmt.Println("Failed to send message:")
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read API response:", err)
		return
	}
	//fmt.Println("API response:", string(body))
	var respData ResponseDataForWX
	err = json.Unmarshal(body, &respData)
	c.JSON(http.StatusOK, gin.H{
		"message": respData.Result,
	})
	return respData.Result
}

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
	if b.Conn == nil {
		b.logger.Error("conn == nil")
		return "conn == nil"
	}
	//发送消息
	respMessage = b.Conn.SpeakToLLM(c, messageMapSlice)
	// 返回消息
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
	LLMName := gptLArk.LLMName
	if LLMName == "BaiChuan" {
		bcConn := BaiChuanConn{
			Key:      gptLArk.Key,
			EndPoint: gptLArk.EndPoint,
			logger:   zapLogger,
		}
		return &bcConn
	} else if LLMName == "WXYY" {
		WXConn := WXConn{
			EndPoint:  gptLArk.EndPoint,
			Key:       gptLArk.Key,
			AppSecret: gptLArk.AppSecret,
			logger:    zapLogger,
		}
		return &WXConn
	}
	return nil
}
