package serverimpl

import (
	"FantasticLife/config"
	"FantasticLife/utils"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type WebSocketManager struct {
	WebSocketPort string
	RpcPort       string
}

// 启动websocket服务器程序，设定serverPort和webSocketPort
func StartWebSocket(config *config.Config, clientManager *ClientManager) {

	serverIp := utils.GetServerIp()
	app := config.App
	webSocketPort := app.WebSocketPort
	rpcPort := app.RpcPort

	serverPort := rpcPort //9001

	http.HandleFunc("/acc", wsPage)

	// 添加处理程序
	go clientManager.start()
	fmt.Println("WebSocket 启动程序成功", serverIp, serverPort)
	http.ListenAndServe(":"+webSocketPort, nil)

}

// 升级协议，在tpl文件中访问对应的链路实现函数调用
func wsPage(w http.ResponseWriter, req *http.Request) {

	// 升级协议
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		fmt.Println("升级协议", "ua:", r.Header["User-Agent"], "referer:", r.Header["Referer"])

		return true
	}}).Upgrade(w, req, nil)
	if err != nil {
		http.NotFound(w, req)

		return
	}

	fmt.Println("webSocket 建立连接:", conn.RemoteAddr().String())

	currentTime := uint64(time.Now().Unix())
	client := NewClient(conn.RemoteAddr().String(), conn, currentTime)

	go client.read()
	go client.write()

	//TODO 用户连接事件
	//clientManager.Register <- client
}