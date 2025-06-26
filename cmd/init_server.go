package cmd

import (
	"go-api-template/common/config"
	"go-api-template/common/logger"
	"go-api-template/handler/controller"
	"go-api-template/internal/adapter"
	"go-api-template/internal/service"
)

func init() {
	config.LoadConfig() // 初始化配置
	logger.InitLogger() // 初始化日志logger

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
