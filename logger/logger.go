package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Logger() gin.HandlerFunc {
	logger, _ := zap.NewProduction(
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("Request processed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.RequestURI),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Since(start)),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
