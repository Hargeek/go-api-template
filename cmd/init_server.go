package cmd

import (
	"go-api-template/common/config"
	"go-api-template/common/logger"
	"go-api-template/store/db"
)

func init() {
	config.LoadConfig() // 初始化配置
	logger.InitLogger() // 初始化日志logger
	db.Init()           // 初始化数据库连接
}
