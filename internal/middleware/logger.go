package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware 结构化日志中间件
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 计算延迟
		latency := time.Since(start)

		// 获取状态码
		status := c.Writer.Status()

		// 构建日志字段
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", raw),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		// 添加错误信息
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		// 根据状态码选择日志级别
		switch {
		case status >= 500:
			logger.Error("HTTP Request", fields...)
		case status >= 400:
			logger.Warn("HTTP Request", fields...)
		default:
			logger.Info("HTTP Request", fields...)
		}
	}
}
