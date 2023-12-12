package utils

import (
	"go.uber.org/zap"
)

// NewZapLogger returns a zap logger instance.
func NewZapLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction() // 或者使用 zap.NewDevelopment() 用于开发环境
	if err != nil {
		//log.Fatalf("can't initialize zap logger: %v", err)
		logger.Error("can't initialize zap logger: %v", zap.Error(err))
		return nil, err
	}
	return logger, nil
}
