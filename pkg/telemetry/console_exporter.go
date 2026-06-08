package telemetry

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// consoleExporter 是开发模式下的轻量 Span 输出器，每个 Span 只打一行。
//
// 输出格式：
//
//	[trace] <traceID前8位> | <span名称>                  | <耗时ms> | <关键属性>
//
// 示例：
//
//	[trace] 4bf92f35 | GET /api/v1/tasks                  |  12.34ms | http.status=200
//	[trace] 4bf92f35 | db.query SELECT                    |   3.12ms | db.operation=SELECT
type consoleExporter struct {
	writer io.Writer
}

func newConsoleExporter() *consoleExporter {
	return &consoleExporter{writer: os.Stdout}
}

func (e *consoleExporter) ExportSpans(_ context.Context, spans []sdktrace.ReadOnlySpan) error {
	for _, span := range spans {
		traceID := span.SpanContext().TraceID().String()
		if len(traceID) > 8 {
			traceID = traceID[:8]
		}

		duration := span.EndTime().Sub(span.StartTime())
		ms := float64(duration.Microseconds()) / 1000

		// 只摘取有助于快速定位的关键属性
		var parts []string
		for _, attr := range span.Attributes() {
			key := string(attr.Key)
			switch key {
			case "http.status_code", "http.method", "http.route",
				"db.operation", "db.system",
				"rpc.method", "rpc.service":
				parts = append(parts, fmt.Sprintf("%s=%v", key, attr.Value.AsInterface()))
			case "db.statement":
				// SQL 语句截断，避免过长
				stmt := fmt.Sprintf("%v", attr.Value.AsInterface())
				if len(stmt) > 60 {
					stmt = stmt[:60] + "..."
				}
				parts = append(parts, "sql="+stmt)
			}
		}

		// 错误状态单独标注
		if span.Status().Code.String() == "Error" {
			parts = append(parts, "error="+span.Status().Description)
		}

		attrStr := ""
		if len(parts) > 0 {
			attrStr = " | " + strings.Join(parts, " ")
		}

		fmt.Fprintf(e.writer, "[trace] %s | %-42s | %7.2fms%s\n",
			traceID,
			span.Name(),
			ms,
			attrStr,
		)
	}
	return nil
}

func (e *consoleExporter) Shutdown(_ context.Context) error { return nil }
