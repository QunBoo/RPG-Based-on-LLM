package tWXYY

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetAccessToken(client_id, client_secret string) {
	url := "https://aip.baidubce.com/oauth/2.0/token?client_id=" + client_id + "&client_secret=" + client_secret + "&grant_type=client_credentials"
	payload := strings.NewReader(``)
	client := &http.Client{}                          // 创建客户端
	req, err := http.NewRequest("POST", url, payload) // 创建请求
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json") // 添加请求头
	req.Header.Add("Accept", "application/json")
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
	fmt.Println(string(body))
}

func SpeakToWXYY(apiKey, secretKey string) {
	//apiKey := "your_API_Key"
	//secretKey := "your_Secret_Key"
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
	message := []byte(`{  
   "messages": [  
    {"role":"user","content":"请你作为一个文字类型的RPG游戏的讲述者，这个游戏的剧情中会根据我的选择而不同，你将作为这个故事的讲述者，并为我讲述故事并生成选项以供我选择我要做什么，我将会说出我的决策，以供你生成下一段故事。我希望以《哈利波特》的世界观作为故事背景，我想要扮演一个和哈利波特同岁的魔法师，与哈利波特、赫敏、马尔福同一年入学，将与他们共同冒险一次只讲述一段故事，每当你描绘完一段故事后，请你停止生成，并等待我的选择。现在，请为我讲述第一段故事"}  
   ]  
 }`)
	resp, err = http.Post(apiEndpoint, "application/json", bytes.NewBuffer(message))
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
	fmt.Println("API response:", string(body))
}
