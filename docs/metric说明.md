# Metrics 说明

本项目集成 Prometheus 指标采集，HTTP 层指标自动采集，业务指标按需添加。

---

## 访问方式

指标通过独立端口暴露，不经过业务中间件：

```
http://localhost:8081/metrics
```

端口在 `config/conf.yaml` 中配置：

```yaml
server:
  http_port: 8080
  metric_port: 8081
```

---

## 已有指标

### 自动采集（无需代码）

| 指标                              | 类型           | 说明                                        |
|---------------------------------|--------------|-------------------------------------------|
| `http_requests_total`           | Counter      | 请求总量，按 `method` / `path` / `status` 分维度   |
| `http_request_duration_seconds` | Histogram    | 请求延迟分布，可计算 P99                            |
| `app_build_info`                | Gauge（值恒为 1） | 携带 `branch`、`revision`、`go_version` 等构建信息 |
| Go 运行时指标                        | 多种           | 内存、GC、goroutine 数量，由 SDK 自动提供             |

### 示例业务指标

| 指标                      | 采集位置                              | 说明                                   |
|-------------------------|-----------------------------------|--------------------------------------|
| `weather_query_total`   | `WeatherServiceImpl.QueryWeather` | 按城市和结果（success/fail）统计，演示 Counter 用法 |
| `task_operations_total` | `TaskServiceImpl` 各方法             | 按操作类型和结果统计，演示多维度 Counter 用法          |

---

## 添加新业务指标

**第一步**：在 `common/metrics/metrics.go` 注册

```go
var MyCounter = promauto.NewCounterVec(prometheus.CounterOpts{
    Name: "my_operation_total",
    Help: "描述这个指标",
    },
    []string{"label1", "label2"}
)
```

**第二步**：在 Service 层调用

```go
metrics.MyCounter.WithLabelValues("val1", "val2").Inc()
```

**第三步**：验证

启动服务后访问 `http://localhost:8081/metrics`，搜索指标名称。

---

## 指标类型速查

| 场景            | 类型        | 方法                              |
|---------------|-----------|---------------------------------|
| 统计次数（请求量、操作数） | Counter   | `.Inc()` / `.Add(n)`            |
| 当前状态（并发数、连接池） | Gauge     | `.Set(v)` / `.Inc()` / `.Dec()` |
| 延迟/大小分布       | Histogram | `.Observe(seconds)`             |

