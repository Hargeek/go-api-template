package middle

import (
	"go-api-template/common/types"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Trace 为每个 HTTP 请求自动生成 Span，挂载在全局 TracerProvider 下
// 必须注册在所有业务中间件之前，确保后续中间件和 Handler 的日志都能关联到同一条链路
func Trace() gin.HandlerFunc {
	return otelgin.Middleware(types.ServiceName)
}
