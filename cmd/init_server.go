package cmd

import (
	"go-api-template/common/config"
	"go-api-template/common/logger"
)

func init() {
	config.LoadConfig() // 初始化配置
	logger.InitLogger() // 初始化日志logger
}
