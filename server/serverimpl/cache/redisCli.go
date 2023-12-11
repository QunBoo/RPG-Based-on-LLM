package cache

import (
	"FantasticLife/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisCli redis.Client

const (
	serversHashKey       = "acc:hash:servers" // 全部的服务器
	serversHashCacheTime = 2 * 60 * 60        // key过期时间
	serversHashTimeout   = 3 * 60             // 超时时间
)

func NewRedisCli(logger *zap.Logger, config *config.Config) (redisClient *RedisCli) {
	redisConf := config.Redis
	redisClient = (*RedisCli)(redis.NewClient(&redis.Options{
		Addr:         redisConf.Addr,
		Password:     redisConf.Password,
		DB:           redisConf.DB,
		PoolSize:     redisConf.PoolSize,
		MinIdleConns: redisConf.MinIdleConns,
	}))

	pong, err := redisClient.Ping(context.Background()).Result()
	//fmt.Println("初始化redis:", pong, err)
	logger.Info("初始化redis", zap.String("pong", pong), zap.Error(err))
	if err != nil {
		panic(fmt.Sprintf("初始化redis失败:%v", err))
	}
	return redisClient
	// Output: PONG <nil>
}

func (cli *RedisCli) getServersHashKey() (key string) {
	key = fmt.Sprintf("%s", serversHashKey)
	return key
}
