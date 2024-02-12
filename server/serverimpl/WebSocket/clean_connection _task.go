package WebSocket

import (
	"FantasticLife/server/serverimpl/task"
	"fmt"
	"runtime/debug"
	"time"
)

func (manager *ClientManager) CleanConnectionInit() {
	task.Timer(3*time.Second, 30*time.Second, manager.cleanConnection, "", nil, nil)
}

// TODO 清理超时连接
func (manager *ClientManager) cleanConnection(param interface{}) (result bool) {
	result = true

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ClearTimeoutConnections stop", r, string(debug.Stack()))
		}
	}()

	fmt.Println("定时任务: 清理超时连接", param)

	manager.ClearTimeoutConnections()

	return
}
