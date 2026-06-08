package telemetry

import (
	"os"

	"go-api-template/common/types"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// buildResource 构造资源属性，优先级：环境变量 > 代码内置值
//
// 自动读取的环境变量：
//
//	OTEL_SERVICE_NAME        覆盖服务名（优先级高于代码中的 types.ServiceName）
//	OTEL_RESOURCE_ATTRIBUTES 追加自定义属性，格式 key1=val1,key2=val2
func buildResource() *resource.Resource {
	// 解析 OTEL_SERVICE_NAME：有则用，无则降级到 types.ServiceName
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = types.ServiceName
	}

	// Default() 作为 base，提供 SDK 运行时属性（telemetry.sdk.*）以及 OTEL_RESOURCE_ATTRIBUTES 中的自定义键
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("branch", types.Branch),
			attribute.String("revision", types.Revision),
			attribute.String("go.version", types.GoVersion),
			attribute.String("build.user", types.BuildUser),
		),
	)
	if err != nil {
		return resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("branch", types.Branch),
			attribute.String("revision", types.Revision),
			attribute.String("go.version", types.GoVersion),
			attribute.String("build.user", types.BuildUser),
		)
	}
	return r
}
