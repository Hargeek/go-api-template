package cmd

import (
	"context"
	"fmt"
	"os"

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

	// init weather service
	//weatherAdapter := &adapter.WeatherAdapterImpl{"demo-key"} // 使用字面量初始化
	//apiKey := config.AppConfig.Weather.Apikey
	weatherAdapter := adapter.NewWeatherAdapterImpl("demo-key") // 使用构造函数初始化
	weatherService := service.NewWeatherServiceImpl(weatherAdapter)
	controller.Weather = controller.NewWeatherController(weatherService)

	// init task service
	taskService := service.NewTaskServiceImpl(d)
	controller.Task = controller.NewTaskController(taskService)
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
	_, _ = cyan.Println(systemLogo)
	logger.Info(fmt.Sprintf("GoVersion: %s, Branch: %s, Revision: %s, BuildDate: %s, BuildUser: %s", types.GoVersion, types.Branch, types.Revision, types.BuildDate, types.BuildUser))
}
