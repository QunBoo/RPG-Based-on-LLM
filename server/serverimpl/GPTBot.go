package serverimpl

import (
	"FantasticLife/server"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type GptConn struct {
	Key       string
	EndPoint  string
	AppSecret string
}
type GptBot struct {
	conn    *GptConn
	chatMap []map[string]string
}
type RequestData struct {
	Messages struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}
type BotResponseData struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func (b *GptBot) SpeakToBot(c *gin.Context, messageMap map[string]string) {
	//消息处理
	//var requestData RequestData
	//if err := c.BindJSON(&requestData); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//messageMap := map[string]string{
	//	"role":    requestData.Messages.Role,
	//	"content": requestData.Messages.Content,
	//}
	b.chatMap = append(b.chatMap, messageMap)
	//发送消息
	url := b.conn.EndPoint
	api_key := b.conn.Key
	data := map[string]interface{}{
		"model":    "Baichuan2",
		"messages": b.chatMap,
		"stream":   false,
	}
	jsonData, err := json.Marshal(data)
	client := &http.Client{}                                            // 创建客户端
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData)) // 创建请求
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json") // 添加请求头
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api_key))
	res, err := client.Do(req) // 发送请求
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(res.Body) // 关闭请求

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if res.StatusCode == 200 {
		fmt.Println("请求成功！")
		fmt.Println("响应body:", string(body))
	} else {
		fmt.Println("请求失败，状态码:", res.StatusCode)
		fmt.Println("请求失败，body:", string(body))
	}
	// 解析 JSON
	var botResp BotResponseData
	err = json.Unmarshal(body, &botResp)
	if err != nil {
		fmt.Println("Error parsing JSON: ", err)
		return
	}

	// 将消息添加到 chatMap
	for _, choice := range botResp.Choices {
		//b.Messages[choice.Message.Role] = choice.Message.Content
		messageMap = map[string]string{
			"role":    choice.Message.Role,
			"content": choice.Message.Content,
		}
		b.chatMap = append(b.chatMap, messageMap)
	}
	// 返回消息
	c.JSON(http.StatusOK, gin.H{
		"message": botResp.Choices[0].Message.Content,
	})
}
func (b *GptBot) SpeakToBot_server(c *gin.Context) {
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
	b.chatMap = append(b.chatMap, messageMap)
	//发送消息
	url := b.conn.EndPoint
	api_key := b.conn.Key
	data := map[string]interface{}{
		"model":    "Baichuan2",
		"messages": b.chatMap,
		"stream":   false,
	}
	jsonData, err := json.Marshal(data)
	client := &http.Client{}                                            // 创建客户端
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData)) // 创建请求
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json") // 添加请求头
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api_key))
	res, err := client.Do(req) // 发送请求
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(res.Body) // 关闭请求

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if res.StatusCode == 200 {
		fmt.Println("请求成功！")
		fmt.Println("响应body:", string(body))
	} else {
		fmt.Println("请求失败，状态码:", res.StatusCode)
		fmt.Println("请求失败，body:", string(body))
	}
	// 解析 JSON
	var botResp BotResponseData
	err = json.Unmarshal(body, &botResp)
	if err != nil {
		fmt.Println("Error parsing JSON: ", err)
		return
	}

	// 将消息添加到 chatMap
	for _, choice := range botResp.Choices {
		//b.Messages[choice.Message.Role] = choice.Message.Content
		messageMap = map[string]string{
			"role":    choice.Message.Role,
			"content": choice.Message.Content,
		}
		b.chatMap = append(b.chatMap, messageMap)
	}
	// 返回消息
	c.JSON(http.StatusOK, gin.H{
		"message": botResp.Choices[0].Message.Content,
	})
}
func (b *GptBot) InitBot(c *gin.Context) {
	//初始化GptBot,清空chatMap
	b.chatMap = nil
	c.JSON(http.StatusOK, gin.H{
		"message": "InitBot Success!",
	})
}

//func NewGptConn(Key, EndPoint, Appsecret string) *GptConn {
//	return &GptConn{
//		Key:       Key,
//		EndPoint:  EndPoint,
//		AppSecret: Appsecret,
//	}
//}

func NewGptBot(pConn *GptConn) (server.BOT, error) {
	return &GptBot{
		conn: pConn,
	}, nil
}
