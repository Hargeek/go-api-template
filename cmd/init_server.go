package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"go-api-template/common/config"
	"go-api-template/common/logger"
	"go-api-template/common/types"
	"go-api-template/handler/controller"
	"go-api-template/internal/adapter"
	"go-api-template/internal/service"
	"go-api-template/internal/store/dao"
	"go-api-template/internal/store/db"
	"go-api-template/internal/store/db/migrate"
	"go-api-template/pkg/telemetry"
	"go-api-template/pkg/weather"

	"github.com/fatih/color"
)

// TelemetryShutdown 由 RunServer() 在优雅退出时调用，确保未导出的 Span 全部 flush
var TelemetryShutdown func(context.Context) error

func init() {
	config.LoadConfig() // 初始化配置

	// 注册 TraceProvider/LoggerProvider
	shutdown, err := telemetry.Setup(context.Background())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "telemetry setup failed: %v\n", err)
		os.Exit(1)
	}
	TelemetryShutdown = shutdown

	logger.InitLogger()   // 初始化 logger，挂载 otelslog bridge
	showInfoDisplayLogo() // 显示 logo

	migrate.AutoMigrate() // 自动迁移数据库表结构（触发 GORM 初始化 + OTEL Plugin 注册）

	// init dao
	d := dao.NewDao(db.GetGORM())

	// init hello service
	helloService := service.NewHelloServiceImpl()
	controller.Hello = controller.NewHelloController(helloService)

	initWeather()

	// init task service
	taskService := service.NewTaskServiceImpl(d)
	controller.Task = controller.NewTaskController(taskService)
}

// initWeather 创建天气客户端并完成 Adapter、Service 和 Controller 装配。
func initWeather() {
	weatherClient, err := weather.NewClient(weather.Config{
		BaseURL: config.AppConfig.WeatherConfig.BaseURL,
		Timeout: time.Duration(config.AppConfig.WeatherConfig.TimeoutSeconds) * time.Second,
	})
	if err != nil {
		logger.Error("create weather client failed", "error", err)
		os.Exit(1)
	}
	weatherAdapter := adapter.NewWeatherAdapter(weatherClient)
	weatherService := service.NewWeatherServiceImpl(weatherAdapter)
	controller.Weather = controller.NewWeatherController(weatherService)
}

const systemLogo = `
 _____ ____    ____  ____  _    _____ _____ _      ____  _     ____ _____ _____
/  __//  _ \  /  _ \/  __\/ \  /__ __Y  __// \__/|/  __\/ \   /  _ Y__ __Y  __/
| |  _| / \|  | / \||  \/|| |    / \ |  \  | |\/|||  \/|| |   | / \| / \ |  \  
| |_//| \_/|  | |-|||  __/| |    | | |  /_ | |  |||  __/| |_/\| |-|| | | |  /_ 
\____\\____/  \_/ \|\_/   \_/    \_/ \____\\_/  \|\_/   \____/\_/ \| \_/ \____\
`

func showInfoDisplayLogo() {
	cyan := color.New(color.FgCyan, color.Bold)
	_, _ = cyan.Print(systemLogo)
	logger.Info("build information",
		"go_version", types.GoVersion,
		"branch", types.Branch,
		"revision", types.Revision,
		"build_date", types.BuildDate,
		"build_user", types.BuildUser,
	)
}
