package task

import (
	"fmt"
	"runtime/debug"
	"time"
)

// 服务器定时注册初始化，定时进行服务器注册，将服务器信息输入到redis中，以hashmap形式存储
func ServerInit() {
	Timer(2*time.Second, 60*time.Second, server, "", serverDefer, "")
}

// TODO 服务器注册函数，调用函数获取本机"ip:port"并设置到redis里
func server(param interface{}) (result bool) {
	result = true

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("服务注册 stop", r, string(debug.Stack()))
		}
	}()

	//server := websocket.GetServer()
	//currentTime := uint64(time.Now().Unix())
	//fmt.Println("定时任务，服务注册", param, server, currentTime)
	//
	//cache.SetServerInfo(server, currentTime)

	return
}

// TODO Redis中服务器下线
func serverDefer(param interface{}) (result bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("服务下线 stop", r, string(debug.Stack()))
		}
	}()

	fmt.Println("服务下线", param)

	//server := websocket.GetServer()
	//cache.DelServerInfo(server)

	return
}
