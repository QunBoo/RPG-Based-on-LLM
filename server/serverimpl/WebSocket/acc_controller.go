package WebSocket

import (
	"FantasticLife/utils"
	"encoding/json"
	"fmt"
	"time"
)

// ping
func PingController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = utils.OK
	fmt.Println("webSocket_request ping接口", client.Addr, seq, string(message))

	data = "pong"

	return
}

func LoginController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = utils.OK
	fmt.Println("webSocket_request login接口", client.Addr, seq, string(message))
	currentTime := uint64(time.Now().Unix())

	request := &Login{}
	if err := json.Unmarshal(message, request); err != nil {
		code = utils.ParameterIllegal
		fmt.Println("用户登录 解析数据失败", seq, err)
		//client.logger.Error("用户登录 解析数据失败",)

		return
	}

	fmt.Println("webSocket_request 用户登录", seq, "ServiceToken", request.ServiceToken)

	// 本项目只是演示，所以直接过去客户端传入的用户ID
	if request.UserId == "" || len(request.UserId) >= 20 {
		code = utils.UnauthorizedUserId
		fmt.Println("用户登录 非法的用户", seq, request.UserId)

		return
	}

	if !client.ClientManagerHook.InAppIds(request.AppId) {
		code = utils.Unauthorized
		fmt.Println("用户登录 不支持的平台", seq, request.AppId)

		return
	}

	if client.IsLogin() {
		fmt.Println("用户登录 用户已经登录", client.AppId, client.UserId, seq)
		code = utils.OperationFailure

		return
	}

	client.Login(request.AppId, request.UserId, currentTime)

	//// TODO 存储用户的登录数据到redis
	//userOnline := models.UserLogin(serverIp, serverPort, request.AppId, request.UserId, client.Addr, currentTime)
	//err := cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	//if err != nil {
	//	code = utils.ServerError
	//	fmt.Println("用户登录 SetUserOnlineInfo", seq, err)
	//
	//	return
	//}

	// 用户登录
	login := &login{
		AppId:  request.AppId,
		UserId: request.UserId,
		Client: client,
	}
	clientManagerHook := client.ClientManagerHook
	clientManagerHook.Login <- login

	fmt.Println("用户登录 成功", seq, client.Addr, request.UserId)
	return
}

func HeartbeatController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = utils.OK
	fmt.Println("webSocket_request hb接口", client.Addr, seq, string(message))
	currentTime := uint64(time.Now().Unix())

	request := &HeartBeat{}
	if err := json.Unmarshal(message, request); err != nil {
		code = utils.ParameterIllegal
		fmt.Println("心跳接口 解析数据失败", seq, err)

		return
	}

	fmt.Println("webSocket_request 心跳接口", client.AppId, client.UserId)

	if !client.IsLogin() {
		fmt.Println("心跳接口 用户未登录", client.AppId, client.UserId, seq)
		code = utils.NotLoggedIn

		return
	}
	client.Heartbeat(currentTime)
	// TODO 获取redis用户的登录数据
	//userOnline, err := cache.GetUserOnlineInfo(client.GetKey())
	//if err != nil {
	//	if err == redis.Nil {
	//		code = utils.NotLoggedIn
	//		fmt.Println("心跳接口 用户未登录", seq, client.AppId, client.UserId)
	//
	//		return
	//	} else {
	//		code = utils.ServerError
	//		fmt.Println("心跳接口 GetUserOnlineInfo", seq, client.AppId, client.UserId, err)
	//
	//		return
	//	}
	//}

	//userOnline.Heartbeat(currentTime)
	//err = cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	//if err != nil {
	//	code = utils.ServerError
	//	fmt.Println("心跳接口 SetUserOnlineInfo", seq, client.AppId, client.UserId, err)
	//
	//	return
	//}
	data = "hb"
	return
}

/************************  请求数据  **************************/
// 通用请求数据格式
type Request struct {
	Seq  string      `json:"seq"`            // 消息的唯一Id
	Cmd  string      `json:"cmd"`            // 请求命令字
	Data interface{} `json:"data,omitempty"` // 数据 json
}

// 登录请求数据
type Login struct {
	ServiceToken string `json:"serviceToken"` // 验证用户是否登录
	AppId        uint32 `json:"appId,omitempty"`
	UserId       string `json:"userId,omitempty"`
}

// 心跳请求数据
type HeartBeat struct {
	UserId string `json:"userId,omitempty"`
}
