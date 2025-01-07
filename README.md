# Go API Template

![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.22-%23007d9c)
![Gin Version](https://img.shields.io/badge/Gin-%3E%3D1.10-green)
[![GoDoc](https://godoc.org/github.com/hargeek/go-api-template?status.svg)](https://pkg.go.dev/github.com/hargeek/go-api-template)
[![Contributors](https://img.shields.io/github/contributors/hargeek/go-api-template)](https://github.com/hargeek/go-api-template/graphs/contributors)
[![License](https://img.shields.io/github/license/hargeek/go-api-template)](./LICENSE)

用于快速构建Go API后端工程的模板项目

## 项目结构

```shell
.
├── cmd                     # 项目启动入口
│   ├── init_server.go      # 初始化逻辑
│   └── run_server.go       # 启动逻辑
├── common                  # 公共代码
│   ├── config              # 配置加载
│   ├── error               # 自定义错误码
│   ├── logger              # 日志组件
│   └── types               # 类型定义
├── config                  # 配置文件
│   └── conf.yaml.example   # 示例配置文件
├── controller              # 控制器层，处理请求和响应
│   └── auxiliary.go        # 示例控制器
├── docs                    # 文档
│   ├── docs.go             # Swagger 文档生成
│   ├── swagger.json        # Swagger JSON 文件
│   └── swagger.yaml        # Swagger YAML 文件
├── go.mod                  # Go 模块定义
├── go.sum                  # Go 模块校验
├── main.go                 # 程序入口
├── middle                  # 中间件
│   ├── cors.go             # 跨域中间件
│   └── logger.go           # 访问日志中间件
├── resource                # 静态资源
│   ├── embed.go            # 静态资源嵌入
│   ├── error_code          # 错误码说明文档（Markdown 格式）
│   └── static_resource     # 其他静态资源
├── routers                 # 路由定义
├── service                 # 服务层，处理业务逻辑
├── store                   # 存储层
│   ├── dao                 # 数据访问层（CRUD 操作）
│   ├── db                  # 数据库连接
│   └── model               # 数据模型
└── utils                   # 工具代码，常用工具函数
```

## 开箱即用
- [x] 中间件: 访问日志、跨域
- [x] 使用`Go` 1.21 以上版本支持的`slog`日志库
- [x] 集成`Swagger`文档和`Redoc`文档
- [x] 自定义错误码，以及独立的错误码说明文档，最终会在接口文档中展示
- [x] 通过`yaml`配置文件加载配置，支持通过配置文件启用`debug`模式，快速开启`pprof`性能监控，方便性能调优
- [x] `Docker`部署, 支持`buildx`, 使用多阶段构建，以及在`docker build`时调用`Makefile`，注入服务版本信息等
- [x] 提供常用的`Makefile`命令，如`make build`, `make run`, `make buildx`, `make test`, `make generate-error`等，方便开发时直接使用，更多`make`命令请执行`make`或`make help`查看

## 快速开始

### clone项目

```bash
git clone https://github.com/Hargeek/go-api-template.git
```

### 重命名项目

```bash
mv go-api-template your-project-name
cd your-project-name
sed -i '' 's/go-api-template/your-project-name/g' $(grep go-api-template -rl .)
```

### 修改配置文件

```bash
cd your-project-name
cp config/conf.yaml.example config/conf.yaml
# 准备数据库，修改配置文件
```

### 安装依赖

```bash
go mod tidy
```

### 安装swag命令及make指令

- 建议安装指定版本的swag命令（不同版本的swag命令生成的接口文档格式可能不一样）

```bash
$ go install github.com/swaggo/swag/cmd/swag@v1.8.4
```

- make请参考[make](https://www.gnu.org/software/make/)

### 运行项目

```bash
make run
```