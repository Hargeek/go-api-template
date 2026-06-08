// Package telemetry 管理 OpenTelemetry 的生命周期（Trace Provider）
//
// 使用方式：
//
//  1. 在 cmd/init_server.go 的 init() 中调用 Setup()，获取 shutdown 函数
//  2. 在 cmd/run_server.go 的优雅退出处调用 shutdown，确保 buffer 中的 Span 全部 flush
//
// 环境变量：
//
//	OTEL_SERVICE_NAME              覆盖服务名（默认使用 types.ServiceName）
//	OTEL_RESOURCE_ATTRIBUTES       追加资源属性，格式 key1=val1,key2=val2
//	OTEL_EXPORTER_ENABLED          总开关，true 时才处理子开关，默认 false
//	OTEL_EXPORTER_ENABLED_TRACES   Trace 开关，需总开关为 true，默认 false
//	OTEL_EXPORTER_ENABLED_LOGS     Log 开关（Log 接入后生效），需总开关为 true
//	OTEL_EXPORTER_OTLP_ENDPOINT    OTLP Collector 地址；未设置时降级为 stdout 输出
package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// Setup 初始化所有 OTEL Provider，返回 shutdown 函数。
// shutdown 必须在服务退出前调用，确保未导出的 Span 全部 flush。
func Setup(ctx context.Context) (shutdown func(context.Context) error, err error) {
	tp, err := initTraceProvider(ctx)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return func(ctx context.Context) error {
		return tp.Shutdown(ctx)
	}, nil
}
