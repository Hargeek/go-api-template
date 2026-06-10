# Log 说明

本项目集成 `log/slog` + `otelslog bridge`，本地输出保持 JSON 格式，可按需通过环境变量开启 OTEL 导出，将日志送入 Collector →
Loki，在 Grafana 中实现日志与 Trace 的关联跳转。

---

## 本地输出（默认）

不设置任何环境变量时，日志输出到 stdout（以及配置文件中指定的文件）。

请求日志（由 gin 中间件生成）：

```json
{
  "time": "2026-06-10 16:49:13.107",
  "level": "INFO",
  "msg": "request log",
  "client_ip": "::1",
  "user_name": "",
  "status": 200,
  "latency": "1.690917ms",
  "timestamp": 1781081353,
  "method": "GET",
  "path": "/api/v1/tasks",
  "query": "",
  "user_agent": "curl/8.7.1",
  "trace_id": "e097076fe9bd9d1d3976f703c417f5b4",
  "span_id": "6ffdb393a2c9db18"
}
```

`trace_id` / `span_id` 字段仅在 Trace 启用（span context 有效）时附加。

---

## 环境变量

| 变量                            | 说明                                        |
|-------------------------------|-------------------------------------------|
| `OTEL_EXPORTER_ENABLED`       | 总开关，`true` 时才处理子开关，默认不启用                  |
| `OTEL_EXPORTER_ENABLED_LOGS`  | Log 子开关，需总开关为 `true`                      |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Collector 地址，与 Trace 共用；不设置时 OTEL Log 不采集 |

---

## 常用配置场景

```bash
# 不启用（默认）：仅本地 slog stdout + 文件
go run main.go

# 接入本地 Loki（需先 make log-up 或 make obs-up）
OTEL_EXPORTER_ENABLED=true \
OTEL_EXPORTER_ENABLED_LOGS=true \
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317 \
go run main.go
```

> `OTEL_EXPORTER_ENABLED_LOGS=true` 但不设置 `OTEL_EXPORTER_OTLP_ENDPOINT` 时，OTEL Log 不采集。
> 本地调试直接看 slog 的 stdout 输出即可，无需重复输出。

---

## 本地联调

```bash
make log-up    # 启动 Collector + Loki + Grafana（仅日志）
make log-down  # 停止日志环境
make obs-up    # 同时启动 Trace（Jaeger）+ Log（Loki）完整环境
make obs-down  # 停止全部
```

| 服务         | 地址                    |
|------------|-----------------------|
| Grafana UI | http://localhost:3000 |
| Loki API   | http://localhost:3100 |
| OTLP gRPC  | localhost:4317        |

Grafana 启动后自动加载 Loki 和 Jaeger 数据源。在 Grafana → Explore → Loki 中查询日志，含 trace context 的日志行会出现
**Jaeger** 跳转按钮，点击直接查看对应链路详情。

> **注意**：Log → Trace 跳转需要同时启用 Trace，建议使用 `make obs-up` 启动完整环境。  
> 仅 `make log-up`（无 Jaeger）时 Jaeger 按钮不可用。

![Grafana Log → Trace 跳转演示](images/grafana-loki-jaeger.gif)

---

## Log → Trace 跳转原理

### 数据流

```
应用（slog + otelslog bridge）
  └─ OTLP gRPC ──→ OTEL Collector
       └─ transform/logs processor  ←─ 追加 trace_id 到 body
            └─ OTLP HTTP ──→ Loki
                 └─ Grafana derivedFields（regex）──→ Jaeger
```

### 关键点

1. **本地 slog 输出**：`trace_id` 作为 JSON 字段附加到日志（通过 `withTrace()`），格式不变。

2. **Collector transform**：`transform/logs` processor 在 body 末尾追加 ` trace_id=<hex>`，仅对含 `trace_id` 属性的日志生效（无 trace context 的日志不受影响）。

3. **Loki 存储**：log body = `"request log trace_id=e097076fe9bd9d1d3976f703c417f5b4"`

4. **Grafana regex**：derivedFields 使用 `matcherRegex: 'trace_id=(\w+)'` 从 body 提取，任何含 trace_id 的日志都可跳转，不依赖日志的具体格式。

---

## 日志写入规范

所有涉及请求上下文的日志，使用带 context 的函数：

```go
// 正确：自动附加 trace_id / span_id（通过 otelslog bridge 从 ctx 提取）
logger.InfoContext(ctx, "task created", slog.Uint64("id", uint64(task.ID)))

// 不推荐：丢失 trace 关联
logger.Info("task created", slog.Uint64("id", uint64(task.ID)))
```

启动/关闭等无请求上下文的场景使用不带 context 的函数：

```go
logger.Info("server starting", slog.Int("port", cfg.HttpPort))
```
