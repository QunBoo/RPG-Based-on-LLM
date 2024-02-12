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
	apiEndpoint := fmt.Sprintf("https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/plugin/dswdc1hehb9y4bi8/?access_token=%s", accessToken)

	// 定义一个结构体来匹配新的 JSON 结构
	payloadStruct := struct {
		Query   string   `json:"query"`
		Plugins []string `json:"plugins"`
		Verbose bool     `json:"verbose"`
	}{
		Query:   "高等数字通信这门课学分多少？",
		Plugins: []string{"uuid-zhishiku"},
		Verbose: true,
	}

	// 将结构体序列化为 JSON
	payloadBytes, err := json.Marshal(payloadStruct)
	if err != nil {
		// 处理错误
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// 使用序列化后的 JSON 作为请求体
	resp1, err := http.Post(apiEndpoint, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		// 处理错误
		fmt.Println("Error making HTTP request:", err)
		return
	}

	// 确保响应体被关闭
	defer resp1.Body.Close()
	if err != nil {
		fmt.Println("Failed to send message:")
		return
	}
	defer resp1.Body.Close()
	body, err = ioutil.ReadAll(resp1.Body)
	if err != nil {
		fmt.Println("Failed to read API response:", err)
		return
	}
	fmt.Println("API response:", string(body))
}
