package logger

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func TestWithTrace(t *testing.T) {
	t.Run("append trace fields for valid span context", func(t *testing.T) {
		spanContext := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			SpanID:  trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8},
		})
		ctx := trace.ContextWithSpanContext(context.Background(), spanContext)

		args := withTrace(ctx, []interface{}{"operation", "query"})

		require.Len(t, args, 4)
		traceAttr, ok := args[2].(slog.Attr)
		require.True(t, ok)
		spanAttr, ok := args[3].(slog.Attr)
		require.True(t, ok)
		assert.Equal(t, "trace_id", traceAttr.Key)
		assert.Equal(t, spanContext.TraceID().String(), traceAttr.Value.String())
		assert.Equal(t, "span_id", spanAttr.Key)
		assert.Equal(t, spanContext.SpanID().String(), spanAttr.Value.String())
	})

	t.Run("keep fields unchanged without valid span context", func(t *testing.T) {
		args := []interface{}{"operation", "query"}

		got := withTrace(context.Background(), args)

		assert.Equal(t, args, got)
	})
}
