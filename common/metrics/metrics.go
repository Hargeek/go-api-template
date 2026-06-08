// Package metrics 提供 Prometheus 指标注册与采集。
//
// 使用方式：
//   - HTTP 框架层指标由 handler/middle/metrics.go 中间件自动采集，无需手动调用
//   - 业务指标在 service 层关键路径手动调用，例如：
//       metrics.WeatherQueryTotal.WithLabelValues(city, "success").Inc()
//
// 新增业务指标步骤：
//  1. 在本文件 var 块中用 promauto.NewXxxVec 声明并注册
//  2. 在对应 service 层调用 .WithLabelValues(...).Inc() / .Observe()
//  3. 运行服务后访问 http://localhost:<metric_port>/metrics 验证指标出现
package metrics

import (
	"go-api-template/common/types"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// -------------------------------------------------------------------------
	// HTTP 层指标（由 handler/middle/metrics.go 自动采集，无需手动调用）
	// -------------------------------------------------------------------------

	// HttpRequestsTotal 每个路由的请求总量。
	// 标签 path 使用路由模板（如 /api/v1/tasks/:id），避免路径参数导致高基数。
	// 常用查询：
	//   rate(http_requests_total[1m])           — 每分钟 QPS
	//   sum by (status)(http_requests_total)    — 按状态码汇总
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests, labeled by method, path template and status code.",
		},
		[]string{"method", "path", "status"},
	)

	// HttpRequestDuration 每个路由的请求延迟分布（秒）。
	// 使用默认 Bucket（.005 .01 .025 .05 .1 .25 .5 1 2.5 5 10）。
	// 常用查询：
	//   histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))  — P99 延迟
	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds, labeled by method, path template and status code.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	// -------------------------------------------------------------------------
	// 示例业务指标（演示自定义 Counter 注册方式，实际项目中按需替换）
	// -------------------------------------------------------------------------

	// WeatherQueryTotal 天气查询总量，按城市和结果（success/fail）分维度。
	// 在 WeatherServiceImpl.QueryWeather 返回后调用：
	//   metrics.WeatherQueryTotal.WithLabelValues(city, "success").Inc()
	WeatherQueryTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "weather_query_total",
			Help: "Total number of weather queries, labeled by city and result (success/fail).",
		},
		[]string{"city", "result"},
	)

	// TaskOperationsTotal Task CRUD 操作量，按操作类型和结果分维度
	// 在 TaskServiceImpl 各方法返回前调用：
	//   metrics.TaskOperationsTotal.WithLabelValues("create", "success").Inc()
	TaskOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "task_operations_total",
			Help: "Total number of task CRUD operations, labeled by operation and result (success/fail).",
		},
		[]string{"operation", "result"},
	)
)

func init() {
	// 构造 app_build_info metric
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "app_build_info",
			Help: "Application build info. Value is always 1.",
			ConstLabels: prometheus.Labels{
				"branch":     types.Branch,
				"revision":   types.Revision,
				"go_version": types.GoVersion,
				"build_user": types.BuildUser,
			},
		},
		func() float64 { return 1 },
	))
}
