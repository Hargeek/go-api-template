package telemetry

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// initLogProvider 根据环境变量决定 Log 行为：
//
//	OTEL_EXPORTER_ENABLED=false（或未设置）      → 禁用，不导出
//	OTEL_EXPORTER_ENABLED=true
//	  OTEL_EXPORTER_ENABLED_LOGS=true
//	    OTEL_EXPORTER_OTLP_ENDPOINT 未设置       → noop，不导出
//	    OTEL_EXPORTER_OTLP_ENDPOINT=<addr>       → OTLP gRPC 导出
//	  OTEL_EXPORTER_ENABLED_LOGS 非 true         → 禁用
func initLogProvider(ctx context.Context) (*sdklog.LoggerProvider, error) {
	if os.Getenv("OTEL_EXPORTER_ENABLED") != "true" ||
		os.Getenv("OTEL_EXPORTER_ENABLED_LOGS") != "true" {
		return sdklog.NewLoggerProvider(), nil
	}

	var (
		exporter sdklog.Exporter
		err      error
	)

	if endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); endpoint != "" {
		exporter, err = otlploggrpc.New(ctx,
			otlploggrpc.WithEndpoint(endpoint),
			otlploggrpc.WithInsecure(),
			otlploggrpc.WithRetry(otlploggrpc.RetryConfig{
				Enabled:         true,
				InitialInterval: 500 * time.Millisecond,
				MaxInterval:     5 * time.Second,
				MaxElapsedTime:  30 * time.Second,
			}),
		)
	} else {
		// 未配置 Collector 地址时不导出，避免与 slog 本地输出重复
		//exporter, err = stdoutlog.New()
		return sdklog.NewLoggerProvider(), nil
	}
	if err != nil {
		return nil, err
	}

	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
		sdklog.WithResource(buildResource()),
	)
	global.SetLoggerProvider(lp)
	return lp, nil
}
