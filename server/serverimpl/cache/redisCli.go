package cache

import (
	"FantasticLife/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	ServersHashKey       = "acc:hash:servers" // 全部的服务器
	ServersHashCacheTime = 2 * 60 * 60        // key过期时间
	ServersHashTimeout   = 3 * 60             // 超时时间
)

func NewRedisCli(logger *zap.Logger, config *config.Config) (redisClient *redis.Client) {
	redisConf := config.Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:         redisConf.Addr,
		Password:     redisConf.Password,
		DB:           redisConf.DB,
		PoolSize:     redisConf.PoolSize,
		MinIdleConns: redisConf.MinIdleConns,
	})

	pong, err := redisClient.Ping(context.Background()).Result()
	//fmt.Println("初始化redis:", pong, err)
	logger.Info("初始化redis", zap.String("pong", pong), zap.Error(err))
	if err != nil {
		panic(fmt.Sprintf("初始化redis失败:%v", err))
	}
	return redisClient
	// Output: PONG <nil>
}

func GetServersHashKey() (key string) {
	key = fmt.Sprintf("%s", ServersHashKey)
	return key
}
