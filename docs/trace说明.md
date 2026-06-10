# Trace 说明

本项目集成 OpenTelemetry 链路追踪，HTTP 请求和 SQL 查询自动生成 Span，通过环境变量控制是否启用及导出目标。

---

## 自动采集范围

无需修改业务代码，以下层自动产生 Span：

| 层       | 实现               | Span 内容            |
|---------|------------------|--------------------|
| HTTP 请求 | `otelgin` 中间件    | 方法、路由模板、状态码、客户端 IP |
| SQL 查询  | GORM OTEL Plugin | 操作类型、SQL 语句、数据库类型  |

请求 `GET /api/v1/tasks/1` 产生的链路示例：

```
GET /api/v1/tasks/:id  (12ms)
  └── select tasks      (3ms)   ← SQL 子 Span
```

---

## 环境变量

| 变量                             | 说明                              |
|--------------------------------|---------------------------------|
| `OTEL_EXPORTER_ENABLED`        | 总开关，`true` 时处理子开关，默认不启用         |
| `OTEL_EXPORTER_ENABLED_TRACES` | Trace 开关，需总开关为 `true`           |
| `OTEL_EXPORTER_OTLP_ENDPOINT`  | Collector 地址；不设置时输出到 stdout     |
| `OTEL_SERVICE_NAME`            | 覆盖服务名                           |
| `OTEL_RESOURCE_ATTRIBUTES`     | 追加资源属性，格式 `key1=val1,key2=val2` |

---

## 常用配置场景

```bash
# 不启用（默认）
go run main.go

# 启用，输出到 stdout 查看 Span 结构
OTEL_EXPORTER_ENABLED=true OTEL_EXPORTER_ENABLED_TRACES=true go run main.go

# 接入本地 Jaeger（需先 make trace-up）
OTEL_EXPORTER_ENABLED=true \
OTEL_EXPORTER_ENABLED_TRACES=true \
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317 \
go run main.go
```

---

## 本地 Jaeger 联调

```bash
make trace-up    # 启动 OTEL Collector + Jaeger
make trace-down  # 停止
```

启动后访问 `http://localhost:16686` 查看链路。

---

## Context 透传（重要）

**SQL Span 能挂在 HTTP Span 下，依赖 context 全链路透传。**

新增 DAO 方法时必须接受并传递 `context.Context`：

```go
// 正确：透传 context，SQL Span 会挂在当前请求链路下
func (d *TaskDAO) List(ctx context.Context) ([]model.Task, error) {
return d.db.WithContext(ctx).Find(&tasks).Error
}

// 错误：不传 context，SQL Span 变成孤立 Trace
func (d *TaskDAO) List() ([]model.Task, error) {
return d.db.Find(&tasks).Error
}
```

调用链路：`Controller` 用 `c.Request.Context()` → `Service` 透传 `ctx` → `DAO` 使用 `db.WithContext(ctx)`。

---

## 访问日志中的 trace_id

Trace 启用后，访问日志会自动附加 `trace_id` 和 `span_id`：

```json
{
  "msg": "request log",
  "trace_id": "4bf92f35...",
  "span_id": "00f067aa...",
  "method": "GET",
  "path": "/api/v1/tasks",
  "status": 200
}
```

这通过 `logger.InfoContext(c.Request.Context(), ...)` 实现，`common/logger` 包内部自动从 context 提取。

启动完整可观测性环境（`make obs-up`）并同时开启 Log 导出后，Grafana Loki 中的日志行会出现 **Jaeger** 跳转按钮，点击直接定位到对应链路。详见 [log说明.md](log说明.md)。

---

## 添加业务 Span

自动采集覆盖 HTTP 和 SQL 两层，复杂业务可在 Service 层手动补充：

```go
func (s *TaskServiceImpl) Create(ctx context.Context, ...) (*model.Task, error) {
    ctx, span := otel.Tracer("task").Start(ctx, "TaskService.Create")
    defer span.End()
    
    if err := s.dao.Create(ctx, task); err != nil {
    span.RecordError(err)
    span.SetStatus(codes.Error, err.Error())
    return nil, err
    }
    return task, nil
}
```

