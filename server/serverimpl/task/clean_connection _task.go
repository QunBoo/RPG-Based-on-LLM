package task

import (
	"fmt"
	"runtime/debug"
	"time"
)

// TODO 定时清理超时连接，调用此函数初始化定时器来定时调用函数清理超时连接
func Init() {
	Timer(3*time.Second, 30*time.Second, cleanConnection, "", nil, nil)
}

// TODO 清理超时连接
func cleanConnection(param interface{}) (result bool) {
	result = true

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ClearTimeoutConnections stop", r, string(debug.Stack()))
		}
	}()

	fmt.Println("定时任务，清理超时连接", param)

	//WebSocket.ClearTimeoutConnections()

	return
}
