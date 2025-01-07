package db

import (
	"database/sql"
	"fmt"
	"go-api-template/common/config"
	logger2 "go-api-template/common/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	isInit bool
	GORM   *gorm.DB
	DB     *sql.DB
	err    error
)

// Init db的初始化函数，与数据库建立连接
func Init() {
	if isInit {
		return
	}

	// 构建 PostgreSQL 连接字符串
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		config.AppConfig.DataBaseConfig.Host,
		config.AppConfig.DataBaseConfig.Username,
		config.AppConfig.DataBaseConfig.Password,
		config.AppConfig.DataBaseConfig.Database,
		config.AppConfig.DataBaseConfig.Port,
	)

	// 配置 GORM 的日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      false,         // Disable color
		},
	)

	// 连接 PostgreSQL 数据库
	GORM, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("connecting database failed: " + err.Error())
	}

	// 获取底层的 *sql.DB 对象以设置连接池
	sqlDB, err := GORM.DB()
	if err != nil {
		panic("get database connection failed: " + err.Error())
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(config.AppConfig.DataBaseConfig.MaxIdle)
	sqlDB.SetMaxOpenConns(config.AppConfig.DataBaseConfig.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(config.AppConfig.DataBaseConfig.MaxLife) * time.Second)

	// 标记初始化完成
	isInit = true

	// 记录日志
	logger2.Info(fmt.Sprintf("%s database connected success",
		config.AppConfig.DataBaseConfig.Host+":"+
			fmt.Sprintf("%d", config.AppConfig.DataBaseConfig.Port)+"/"+
			config.AppConfig.DataBaseConfig.Database))

	// 将 *sql.DB 赋值给全局变量
	DB = sqlDB
}

// Close db的关闭函数
func Close() error {
	logger2.Info("closing db connection...")
	return DB.Close()
}
