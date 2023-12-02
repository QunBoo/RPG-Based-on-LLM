package tWXYY

import (
	"fmt"
	"io"
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
