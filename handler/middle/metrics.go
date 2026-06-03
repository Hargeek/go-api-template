package middle

import (
	"strconv"
	"time"

	"go-api-template/common/metrics"

	"github.com/gin-gonic/gin"
)

// Metrics 采集每次 HTTP 请求的总量和延迟，标签使用路由模板（如 /api/v1/tasks/:id）避免高基数
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.FullPath() // 路由模板，非实际 URL，避免 ID 导致高基数
		if path == "" {
			path = "unknown" // 未匹配路由（如 404）
		}
		method := c.Request.Method
		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()

		metrics.HttpRequestsTotal.WithLabelValues(method, path, status).Inc()
		metrics.HttpRequestDuration.WithLabelValues(method, path, status).Observe(duration)
	}
}
