package config

import (
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"runtime"
)

const configDir = "config"

type Config struct {
	GptLark struct {
		Key       string
		EndPoint  string
		AppSecret string
	}
}

func GetConfig() *Config {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	viper.SetConfigName("config")                // 配置文件名称（无扩展名）
	viper.SetConfigType("yaml")                  // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(filepath.Join(basepath)) // 查找配置文件所在的路径

	viper.AutomaticEnv() // 读取匹配的环境变量

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		log.Fatalf("unable to decode into struct, %s", err)
	}
	return config
}
