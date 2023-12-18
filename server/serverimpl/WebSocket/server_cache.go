package WebSocket

import (
	"FantasticLife/server/serverimpl/cache"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

func getServersHashKey() (key string) {
	key = fmt.Sprintf("%s", cache.ServersHashKey)

	return
}

// 设置服务器信息
func (manager *ClientManager) SetServerInfo(server *Server, currentTime uint64) (err error) {
	key := getServersHashKey()
	value := fmt.Sprintf("%d", currentTime)
	redisClient := manager.RedisCli
	number, err := redisClient.Do(context.Background(), "hSet", key, server.String(), value).Int()
	if err != nil {
		fmt.Println("SetServerInfo", key, number, err)
		return
	}
	redisClient.Do(context.Background(), "Expire", key, cache.ServersHashCacheTime)
	return
}

// 下线服务器信息
func (manager *ClientManager) DelServerInfo(server *Server) (err error) {
	key := getServersHashKey()
	redisClient := manager.RedisCli
	number, err := redisClient.Do(context.Background(), "hDel", key, server.String()).Int()
	if err != nil {
		fmt.Println("DelServerInfo", key, number, err)

		return
	}

	if number != 1 {

		return
	}

	redisClient.Do(context.Background(), "Expire", key, cache.ServersHashCacheTime)

	return
}

// 根据Redis中存储的所有服务器信息查找服务器，返回服务器列表，以models.Server存储
func (manager *ClientManager) GetServerAll(currentTime uint64) (servers []*Server, err error) {

	servers = make([]*Server, 0)
	//得到一个在Redis服务器中存储服务器信息的哈希键
	key := getServersHashKey()

	redisClient := manager.RedisCli

	val, err := redisClient.Do(context.Background(), "hGetAll", key).Result()
	if err != nil {
		fmt.Println("Redis hGetAll", key, err)
		return
	}

	valByte, _ := json.Marshal(val)
	fmt.Println("GetServerAll", key, string(valByte))

	serverMap, err := redisClient.HGetAll(context.Background(), key).Result()
	if err != nil {
		fmt.Println("SetServerInfo", key, err)
		return
	}

	for key, value := range serverMap {
		valueUint64, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			fmt.Println("GetServerAll", key, err)

			return nil, err
		}

		// 超时
		if valueUint64+cache.ServersHashTimeout <= currentTime {
			continue
		}

		server, err := StringToServer(key)
		if err != nil {
			fmt.Println("GetServerAll", key, err)

			return nil, err
		}

		servers = append(servers, server)
	}

	return
}
