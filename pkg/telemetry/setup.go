// Package telemetry 管理 OpenTelemetry 的生命周期（Trace Provider + Log Provider）
//
// 使用方式：
//
//  1. 在 cmd/init_server.go 的 init() 中最先调用 Setup()，获取 shutdown 函数
//  2. Setup() 返回后再调用 logger.InitLogger()，确保 otelslog bridge 能拿到 LoggerProvider
//  3. 在 cmd/run_server.go 的优雅退出处调用 shutdown，flush 所有未导出的 Span 和 Log
//
// 环境变量：
//
//	OTEL_SERVICE_NAME              覆盖服务名（默认使用 types.ServiceName）
//	OTEL_RESOURCE_ATTRIBUTES       追加资源属性，格式 key1=val1,key2=val2
//	OTEL_EXPORTER_ENABLED          总开关，true 时才处理子开关，默认 false
//	OTEL_EXPORTER_ENABLED_TRACES   Trace 开关，需总开关为 true，默认 false
//	OTEL_EXPORTER_ENABLED_LOGS     Log 开关，需总开关为 true，默认 false
//	OTEL_EXPORTER_OTLP_ENDPOINT    OTLP Collector 地址；Log 未设置时不采集（Trace 未设置时降级 stdout）
package telemetry

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// Setup 初始化所有 OTEL Provider（Trace + Log），返回 shutdown 函数。
// shutdown 必须在服务退出前调用，确保未导出的 Span / Log 全部 flush。
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

	lp, err := initLogProvider(ctx)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) error {
		return errors.Join(tp.Shutdown(ctx), lp.Shutdown(ctx))
	}, nil
}
