package middle

import (
	"github.com/gin-gonic/gin"
	"go-api-template/common/logger"
	"log/slog"
	"time"
)

// Logger 日志中间件，记录请求日志
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()

		clientIP := c.ClientIP()
		clientUserName := c.GetString("api_username")
		statusCode := c.Writer.Status()
		latencyTime := endTime.Sub(startTime).String()
		timestamp := endTime.Unix()
		reqMethod := c.Request.Method
		reqPath := c.Request.URL.Path
		queryString := c.Request.URL.RawQuery
		userAgent := c.GetHeader("User-Agent")

		// 使用 slog 记录日志
		logger.BaseLogger().Info(
			"request log",
			slog.String("client_ip", clientIP),
			slog.String("user_name", clientUserName),
			slog.Int("status", statusCode),
			slog.String("latency", latencyTime),
			slog.Int64("timestamp", timestamp),
			slog.String("method", reqMethod),
			slog.String("path", reqPath),
			slog.String("query", queryString),
			slog.String("user_agent", userAgent),
		)
	}
}
