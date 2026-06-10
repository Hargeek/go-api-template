package telemetry

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// initTraceProvider 根据环境变量决定 Trace 行为：
//
//	OTEL_EXPORTER_ENABLED=false（或未设置）      → 禁用，不输出
//	OTEL_EXPORTER_ENABLED=true
//	  OTEL_EXPORTER_ENABLED_TRACES=true
//	    OTEL_EXPORTER_OTLP_ENDPOINT 未设置       → stdout compact JSON
//	    OTEL_EXPORTER_OTLP_ENDPOINT=<addr>       → OTLP gRPC 导出
//	  OTEL_EXPORTER_ENABLED_TRACES 非 true       → 禁用
func initTraceProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	if os.Getenv("OTEL_EXPORTER_ENABLED") != "true" ||
		os.Getenv("OTEL_EXPORTER_ENABLED_TRACES") != "true" {
		return sdktrace.NewTracerProvider(), nil
	}

	var (
		exporter sdktrace.SpanExporter
		err      error
	)

	if endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); endpoint != "" {
		// passthrough:///host:port 格式，例如：
		//   OTEL_EXPORTER_OTLP_ENDPOINT=passthrough:///10.0.0.1:4317
		exporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
				Enabled:         true,
				InitialInterval: 500 * time.Millisecond,
				MaxInterval:     5 * time.Second,
				MaxElapsedTime:  30 * time.Second,
			}),
		)
	} else {
		// 启用但未配置 Collector 地址，降级到 stdout 方便本地调试
		// 如需更简洁的单行格式，可替换为：exporter = newConsoleExporter()
		exporter, err = stdouttrace.New()
	}
	if err != nil {
		return nil, err
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(buildResource()),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	), nil
}
