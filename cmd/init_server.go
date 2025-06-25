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
	helloService := &service.HelloServiceImpl{}
	controller.Hello.Service = helloService
	// init weather service
	weatherAdapter := &adapter.WeatherAdapterImpl{ApiKey: "demo-key"}
	weatherService := &service.WeatherServiceImpl{Adapter: weatherAdapter}
	controller.Weather.Service = weatherService
}
