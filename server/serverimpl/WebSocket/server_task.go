package WebSocket

import (
	"FantasticLife/server/serverimpl/task"
	"FantasticLife/utils"
	"fmt"
	"go.uber.org/zap"
	"runtime/debug"
	"time"
)

// 服务器定时注册初始化，定时进行服务器注册，将服务器信息输入到redis中，以hashmap形式存储
func (manager *ClientManager) ServerInit() {
	task.Timer(2*time.Second, 60*time.Second, manager.server, "", manager.serverDefer, "")
}

// TODO 服务器注册函数，调用函数获取本机"ip:port"并设置到redis里
func (manager *ClientManager) server(param interface{}) (result bool) {
	result = true

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("服务注册 stop", r, string(debug.Stack()))
		}
	}()

	server := NewServer(utils.GetServerIp(), "8080")
	currentTime := uint64(time.Now().Unix())
	fmt.Println("定时任务，服务注册", param, server, currentTime)

	err := manager.SetServerInfo(server, currentTime)
	if err != nil {
		manager.logger.Error("定时任务，服务注册错误", zap.Error(err), zap.Any("server", server), zap.Uint64("currentTime", currentTime))
		return false
	}

	return
}

// TODO Redis中服务器下线
func (manager *ClientManager) serverDefer(param interface{}) (result bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("服务下线 stop", r, string(debug.Stack()))
		}
	}()

	fmt.Println("服务下线", param)

	server := NewServer(utils.GetServerIp(), "8080")
	err := manager.DelServerInfo(server)
	if err != nil {
		manager.logger.Error("服务下线错误", zap.Error(err), zap.Any("server", server), zap.Any("param", param))
		return false
	}

	return
}
