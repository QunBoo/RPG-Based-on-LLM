package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func ZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		t := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		latency := time.Since(t)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		// 记录日志
		logger.Info("request",
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("clientIP", clientIP),
			zap.String("method", method),
			zap.String("path", path),
		)
	}
}
