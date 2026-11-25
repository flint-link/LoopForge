package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryMiddleware 异常恢复中间件
func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Error("Panic recovered",
			zap.Any("error", recovered),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("client_ip", c.ClientIP()),
		)

		c.AbortWithStatusJSON(500, gin.H{
			"code":    500,
			"message": "Internal Server Error",
			"data":    nil,
		})
	})
}
