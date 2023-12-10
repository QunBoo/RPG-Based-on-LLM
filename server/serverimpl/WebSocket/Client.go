package WebSocket

import (
	"FantasticLife/utils"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"runtime/debug"

	"github.com/gorilla/websocket"
)

const (
	// 用户连接超时时间
	heartbeatExpirationTime = 6 * 60
)

// 用户登录
type login struct {
	AppId  uint32
	UserId string
	Client *Client
}

// GetKey 获取 key
func (l *login) GetKey() (key string) {
	key = GetUserKey(l.AppId, l.UserId)

	return
}

// 用户连接
type Client struct {
	Addr              string          // 客户端地址
	Socket            *websocket.Conn // 用户连接
	Send              chan []byte     // 待发送的数据
	AppId             uint32          // 登录的平台Id app/web/ios
	UserId            string          // 用户Id，用户登录以后才有
	FirstTime         uint64          // 首次连接事件
	HeartbeatTime     uint64          // 用户上次心跳时间
	LoginTime         uint64          // 登录时间 登录以后才有
	ClientManagerHook *ClientManager
	logger            *zap.Logger
}

// 初始化
func NewClient(addr string, socket *websocket.Conn, firstTime uint64, clientManagerHook *ClientManager, logger *zap.Logger) (client *Client) {
	client = &Client{
		Addr:              addr,
		Socket:            socket,
		Send:              make(chan []byte, 100),
		FirstTime:         firstTime,
		HeartbeatTime:     firstTime,
		ClientManagerHook: clientManagerHook,
		logger:            logger,
	}

	return
}

// GetKey 获取 key
func (c *Client) GetKey() (key string) {
	key = GetUserKey(c.AppId, c.UserId)

	return
}

// 读取客户端数据
func (c *Client) read() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("write stop", string(debug.Stack()), r)
		}
	}()

	defer func() {
		//fmt.Println("读取客户端数据 关闭send", c)
		c.logger.Info("读取客户端数据 关闭send", zap.String("addr", c.Addr))
		close(c.Send)
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			//fmt.Println("读取客户端数据 错误", c.Addr, err)
			c.logger.Error("读取客户端数据 错误", zap.String("addr", c.Addr), zap.Error(err))
			return
		}
		// 处理程序
		//fmt.Println("[Cli::read()]读取客户端数据 处理:", string(message))
		c.logger.Info("读取客户端数据", zap.String("addr", c.Addr), zap.String("message", string(message)))

		c.ProcessData(message)
	}
}

// 向客户端写数据
func (c *Client) write() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("write stop", string(debug.Stack()), r)

		}
	}()

	defer func() {
		//clientManager.Unregister <- c
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// 发送数据错误 关闭连接
				//fmt.Println("Client发送数据 关闭连接", c.Addr, "ok", ok)
				c.logger.Info("Client发送数据 关闭连接", zap.String("addr", c.Addr), zap.Bool("ok", ok))

				return
			}
			//fmt.Printf("[Cli::write()]Client发送数据%s\n", message)
			c.logger.Info("[Cli::write()] Client发送数据", zap.String("addr", c.Addr), zap.String("message", string(message)))
			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// SendMsg 发送数据
func (c *Client) SendMsg(msg []byte) {

	if c == nil {

		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("SendMsg stop:", r, string(debug.Stack()))
		}
	}()
	//fmt.Printf("[Cli::SendMsg]:%s\n", msg)
	c.logger.Info("[Cli::SendMsg] SendMsg", zap.String("addr", c.Addr), zap.String("message", string(msg)))
	c.Send <- msg
}

// close 关闭客户端连接
func (c *Client) close() {
	close(c.Send)
}

// 用户登录
func (c *Client) Login(appId uint32, userId string, loginTime uint64) {
	c.AppId = appId
	c.UserId = userId
	c.LoginTime = loginTime
	// 登录成功=心跳一次
	c.Heartbeat(loginTime)
}

// 用户心跳
func (c *Client) Heartbeat(currentTime uint64) {
	c.HeartbeatTime = currentTime

	return
}

// 心跳超时
func (c *Client) IsHeartbeatTimeout(currentTime uint64) (timeout bool) {
	if c.HeartbeatTime+heartbeatExpirationTime <= currentTime {
		timeout = true
	}

	return
}

// 是否登录了
func (c *Client) IsLogin() (isLogin bool) {

	// 用户登录了
	if c.UserId != "" {
		isLogin = true

		return
	}

	return
}

func (c *Client) ProcessData(message []byte) {

	fmt.Println("处理数据", c.Addr, string(message))

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("处理数据 stop", r)
		}
	}()

	request := &Request{}

	err := json.Unmarshal(message, request)
	if err != nil {
		fmt.Println("处理数据 json Unmarshal", err)
		c.SendMsg([]byte("数据不合法"))

		return
	}

	requestData, err := json.Marshal(request.Data)
	if err != nil {
		fmt.Println("处理数据 json Marshal", err)
		c.SendMsg([]byte("处理数据失败"))

		return
	}

	seq := request.Seq
	cmd := request.Cmd

	var (
		code uint32
		msg  string
		data interface{}
	)

	// request
	fmt.Println("接入模块服务请求", cmd, c.Addr)

	// 在句柄注册的map中寻找句柄，并调用句柄，得到响应Response
	clientManagerHook := c.ClientManagerHook
	if value, ok := clientManagerHook.GetHandlers(cmd); ok {
		code, msg, data = value(c, seq, requestData)
	} else {
		code = utils.RoutingNotExist
		fmt.Println("处理数据 路由不存在", c.Addr, "cmd", cmd)
	}

	msg = utils.GetErrorMessage(code, msg)

	responseHead := NewResponseHead(seq, cmd, code, msg, data)

	headByte, err := json.Marshal(responseHead)
	if err != nil {
		fmt.Println("处理数据 json Marshal", err)

		return
	}
	// fmt.Printf("《《《《headByte: %s\n", headByte)
	c.SendMsg(headByte)

	//fmt.Println("acc_response send", c.Addr, c.AppId, c.UserId, "cmd", cmd, "code", code)
	c.logger.Info("acc_response", zap.String("addr", c.Addr), zap.Uint32("appId", c.AppId), zap.String("userId", c.UserId), zap.String("cmd", cmd), zap.Uint32("code", code))

	return
}
