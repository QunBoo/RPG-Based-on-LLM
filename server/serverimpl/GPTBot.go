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
	chatMap map[string]string
}

func (b *GptBot) BOTChat(c *gin.Context) {
	url := b.conn.EndPoint
	api_key := b.conn.Key
	fmt.Println("api_key:", api_key, "url:", url)
	data := map[string]interface{}{
		"model": "Baichuan2",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": "你是谁？",
			},
		},
		"stream": false,
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
	defer res.Body.Close() // 关闭请求

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if res.StatusCode == 200 {
		fmt.Println("请求成功！")
		fmt.Println("响应body:", string(body))
		fmt.Println("请求成功，X-BC-Request-Id:", res.Header.Get("X-BC-Request-Id"))
	} else {
		fmt.Println("请求失败，状态码:", res.StatusCode)
		fmt.Println("请求失败，body:", string(body))
		fmt.Println("请求失败，X-BC-Request-Id:", res.Header.Get("X-BC-Request-Id"))
	}
}
func (b *GptBot) BOTRemember(c *gin.Context) {

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
