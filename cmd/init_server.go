package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"go-api-template/common/config"
	"go-api-template/common/logger"
	"go-api-template/common/types"
	"go-api-template/handler/controller"
	"go-api-template/internal/adapter"
	"go-api-template/internal/service"
)

func init() {
	config.LoadConfig()   // 初始化配置
	logger.InitLogger()   // 初始化日志logger
	showInfoDisplayLogo() // 显示logo

	// init hello service
	helloService := service.NewHelloServiceImpl()
	controller.Hello = controller.NewHelloController(helloService)

	// init weather service
	//weatherAdapter := &adapter.WeatherAdapterImpl{"demo-key"} // 使用字面量初始化
	//apiKey := config.AppConfig.Weather.Apikey
	weatherAdapter := adapter.NewWeatherAdapterImpl("demo-key") // 使用构造函数初始化
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
	_, _ = cyan.Println(systemLogo)
	logger.Info(fmt.Sprintf("GoVersion: %s, Branch: %s, Revision: %s, BuildDate: %s, BuildUser: %s", types.GoVersion, types.Branch, types.Revision, types.BuildDate, types.BuildUser))
}
