package WebSocket

import (
	"FantasticLife/config"
	"FantasticLife/utils"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type WebSocketManager struct {
	WebSocketPort string
	RpcPort       string
}

// 启动websocket服务器程序，设定serverPort和webSocketPort
func StartWebSocket(config *config.Config, clientManager *ClientManager, logger *zap.Logger) {

	serverIp := utils.GetServerIp()
	app := config.App
	webSocketPort := app.WebSocketPort
	rpcPort := app.RpcPort

	serverPort := rpcPort //9001
	wsPage := func(w http.ResponseWriter, req *http.Request) {
		// 从acc接收到http请求后，升级协议
		conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
			//fmt.Println("升级协议", "ua:", r.Header["User-Agent"], "referer:", r.Header["Referer"])
			logger.Info("升级协议", zap.Any("ua", r.Header["User-Agent"]), zap.Any("referer", r.Header["Referer"]))
			return true
		}}).Upgrade(w, req, nil)
		if err != nil {
			http.NotFound(w, req)

			return
		}

		//fmt.Println("webSocket 建立连接:", conn.RemoteAddr().String())
		logger.Info("webSocket 建立连接:", zap.String("Addr", conn.RemoteAddr().String()))

		currentTime := uint64(time.Now().Unix())
		client := NewClient(conn.RemoteAddr().String(), conn, currentTime, clientManager, logger)

		go client.read()
		go client.write()
		// 用户连接事件
		clientManager.RegisterChan <- client
	}

	http.HandleFunc("/acc", wsPage)

	// 启动管道事务的处理，包括建立连接、用户登录、断开连接、广播事件
	go clientManager.start()
	//fmt.Println("WebSocket 启动程序成功", "serverIp:", serverIp, "rpcPort:", serverPort)
	logger.Info("WebSocket 启动程序成功", zap.String("serverIp", serverIp), zap.String("rpcPort", serverPort))
	err := http.ListenAndServe(":"+webSocketPort, nil)
	if err != nil {
		logger.Error("WebSocket 启动程序失败", zap.Error(err))
	}

}
