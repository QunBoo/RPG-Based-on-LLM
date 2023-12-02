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
	data := map[string]interface{}{
		"model": "Baichuan2",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": "请你为我提供一个文字类型的RPG游戏，游戏过程中会根据我的选择推进故事的发展，\n你将作为这个故事的讲述者，并为我生成描绘对应的场景以供我选择我要做什么，我将会说出我的决策，以供你生成下一段故事\n我希望以《哈利波特》的世界观作为故事背景，我想要扮演一个和哈利波特同岁的魔法师，与哈利波特、马尔福同一年入学，根据我的选择的不同我可能会成为哈利波特的重要伙伴，或是黑魔法师中的一员，甚至代替哈利波特成为救世主。\n一次只描绘一个场景，每当你描绘完一个场景后，请你停止生成，并等待我的选择",
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
