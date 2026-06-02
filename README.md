# Go API Template

![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.25-%23007d9c)
![Gin Version](https://img.shields.io/badge/Gin-%3E%3D1.12-green)
[![GoDoc](https://godoc.org/github.com/hargeek/go-api-template?status.svg)](https://pkg.go.dev/github.com/hargeek/go-api-template)
[![Contributors](https://img.shields.io/github/contributors/hargeek/go-api-template)](https://github.com/hargeek/go-api-template/graphs/contributors)
[![License](https://img.shields.io/github/license/hargeek/go-api-template)](./LICENSE)

用于快速构建 Go REST API 后端工程的生产就绪模板，提供完整的分层架构、配置管理、日志、接口文档、错误码体系和 Docker 部署支持。

## 技术栈

| 组件 | 库 | 说明 |
|------|-----|------|
| HTTP 框架 | `gin-gonic/gin` | 路由、中间件、参数绑定 |
| 数据库 ORM | `gorm.io/gorm` + PostgreSQL 驱动 | 连接池管理、迁移 |
| 配置管理 | `spf13/viper` | YAML 配置 + 环境变量覆盖 |
| 日志 | 标准库 `log/slog` | JSON 结构化日志，支持多输出目标 |
| API 文档 | `swaggo/swag` + Redoc | 自动生成 Swagger / Stoplight 文档 |
| 性能分析 | `gin-contrib/pprof` | debug 模式下暴露 pprof 端点 |
| 错误处理 | 自定义 `ErrCode` + stringer | 六位分层错误码，自动生成字符串映射 |

## 项目结构

```
.
├── main.go                     # 程序入口（Swagger 注解定义在此）
├── cmd/
│   ├── init_server.go          # 依赖初始化与手动依赖注入
│   └── run_server.go           # Gin 引擎构建、启动与优雅退出
├── handler/                    # 表现层
│   ├── routers/                # 路由注册（按模块拆分）
│   ├── controller/             # 请求处理、参数校验、响应组装
│   └── middle/                 # 中间件（CORS、访问日志）
├── internal/                   # 业务核心（包外不可直接引用）
│   ├── service/                # 业务逻辑接口 + 实现
│   ├── adapter/                # 外部服务适配器（接口隔离第三方依赖）
│   ├── static/                 # embed 静态资源（错误码文档等）
│   └── store/
│       ├── db/                 # 数据库连接管理（GORM + 连接池）
│       ├── dao/                # 数据访问对象（预留，按需添加）
│       └── model/              # 数据模型（预留，按需添加）
├── common/                     # 公共基础库
│   ├── config/                 # 配置结构定义与加载（含零值校验）
│   ├── logger/                 # 日志初始化与封装
│   ├── error/                  # 错误码定义（go generate 生成字符串映射）
│   └── types/                  # 通用类型、统一响应结构、构建信息变量
├── config/
│   ├── conf.yaml               # 运行时配置文件（git ignored）
│   └── conf.yaml.example       # 配置模板
├── pkg/utils/                  # 工具函数（预留扩展）
├── scripts/                    # 运维/辅助脚本（预留扩展）
├── docs/                       # swag 自动生成的 Swagger 文档
├── Makefile                    # 常用开发命令
└── Dockerfile                  # 多阶段构建镜像
```

## 开箱即用

- **分层架构**：Controller → Service → Adapter 经典三层分离，各层依赖接口而非实现，便于 Mock 测试和切换实现
- **统一响应**：所有接口返回 `{code, msg, data}` 结构，通过 `ApiResponse` 统一收口
- **错误码体系**：六位数字编码（前三位模块、后三位错误），通过 `go generate` 自动生成 `String()` 方法
- **配置零值校验**：启动时反射检查所有必填配置项，任何字段为空即 `panic`，防止配置缺漏在运行时才暴露
- **结构化日志**：基于 `slog` 的 JSON 日志，支持同时输出到 stdout 和文件，访问日志中间件记录完整请求信息
- **接口文档**：集成 Swagger UI（`/swagger/index.html`）和 Stoplight Redoc（`/doc`）双文档界面
- **性能调试**：`debug: true` 时自动注册 pprof 路由，生产环境自动关闭
- **优雅退出**：监听 `SIGINT/SIGTERM`，5 秒超时内完成在途请求后关闭 Server 和数据库连接
- **构建信息注入**：Makefile 通过 `-ldflags` 将 Branch/Revision/BuildDate/BuildUser 注入二进制，健康检查接口可直接返回
- **Docker 部署**：多阶段构建，最终镜像基于 Alpine，以非 root 用户运行，支持 `buildx` 交叉编译
- **热重载开发**：集成 `air`，代码改动后自动重启

## API 端点

所有接口挂载在 `/api/v1` 前缀下：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/health` | 健康检查，返回服务版本和构建信息 |
| GET | `/api/v1/delayed-health?delay_sec=N` | 延迟 N 秒响应，用于超时测试 |
| GET | `/api/v1/echo-get` | 回显请求信息（IP、Header、查询参数） |
| POST | `/api/v1/echo-post` | 回显请求信息（IP、Header、Body） |
| GET | `/api/v1/hello` | Hello World 示例接口 |
| GET | `/api/v1/weather?city=城市名` | 查询天气示例接口（演示 Adapter 模式） |

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/Hargeek/go-api-template.git
```

### 2. 重命名为你的项目

```bash
mv go-api-template your-project-name
cd your-project-name
sed -i '' 's/go-api-template/your-project-name/g' $(grep go-api-template -rl .)
```

### 3. 准备配置文件

```bash
cp config/conf.yaml.example config/conf.yaml
# 按实际情况修改数据库地址、端口、账号密码等
```

配置文件说明（`config/conf.yaml`）：

```yaml
debug: false          # true 时开启 pprof 和 gin debug 日志
env: local            # 环境标识，健康检查接口会返回此值
server:
  http_port: 8080
database:
  host: localhost
  port: 5432
  db_name: go_api_template
  db_user: go_api_template
  db_password: 123456
  log_mode: false       # true 时打印慢查询日志（阈值 3s）
  max_idle_conn: 10
  max_open_conn: 100
  max_life_time: 300    # 连接最大存活时间（秒）
logging:
  level: info           # debug / info / warn / error
  output:
    - stdout            # 标准输出
    - http.log          # 同时写入文件
```

### 4. 安装依赖

```bash
go mod tidy
```

### 5. 安装 swag（生成接口文档）

```bash
go install github.com/swaggo/swag/cmd/swag@v1.16.4
```

### 6. 运行项目

```bash
make run          # 直接运行
make run-air      # air 热重载（需先安装 air）
make local        # 使用本地专用配置运行
```

## 常用 Make 命令

```bash
make help              # 查看所有命令说明
make run               # 运行默认环境
make run-air           # air 热重载开发
make local             # 使用 conf-local.yaml 运行
make build             # 编译为 bin/app
make buildx            # CI/CD 用精简编译（注入版本信息，输出 /bin/server）
make swagger           # 格式化并重新生成 Swagger 文档
make test              # 运行所有测试并生成覆盖率报告
make generate-error    # 通过 stringer 重新生成错误码字符串映射
make vet               # go vet 静态检查
make fieldalignment    # 检查并修复结构体内存对齐
make install-hook      # 安装 git hooks（.githooks/）
make clean             # 清理编译产物
```

## Docker 部署

```bash
# 构建镜像
docker build -t go-api-template:latest .

# 运行容器（挂载配置文件）
docker run -d \
  -p 8080:8080 \
  -v /path/to/conf.yaml:/config/conf.yaml \
  go-api-template:latest
```

## 扩展指南

### 添加新业务模块

以添加 `user` 模块为例：

1. `internal/store/model/` — 添加数据模型 `user.go`
2. `internal/store/dao/` — 添加 DAO 操作
3. `internal/service/` — 定义 `UserService` 接口和实现
4. `handler/controller/` — 添加 `user.go` 控制器
5. `handler/routers/` — 注册路由，在 `init.go` 的 `InitApiRouter` 中调用
6. `cmd/init_server.go` — 在 `init()` 中完成依赖装配

### 替换天气 Adapter

`internal/adapter/weather_adapter_impl.go` 中的实现是演示用的 Mock。接入真实第三方 API 时：

1. 新建一个实现 `WeatherAdapter` 接口的结构体
2. 在 `cmd/init_server.go` 中将 API Key 改为从配置读取后传入
3. 原有业务逻辑和测试代码无需改动

### 添加新错误码

在 `common/error/const.go` 中按规则添加常量后，运行：

```bash
make generate-error
```

错误码编码规则：六位数字，前三位为模块号，后三位为错误序号（如 `101001` = 系统模块第 1 个参数错误）。

## License

[Apache 2.0](./LICENSE)
