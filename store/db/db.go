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

func Init() {
	if isInit {
		return
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.AppConfig.DataBaseConfig.Host,
		config.AppConfig.DataBaseConfig.Username,
		config.AppConfig.DataBaseConfig.Password,
		config.AppConfig.DataBaseConfig.Database,
		config.AppConfig.DataBaseConfig.Port,
	)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      false,         // Disable color
		},
	)
	GORM, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("connecting database failed: " + err.Error())
	}
	sqlDB, err := GORM.DB()
	if err != nil {
		panic("get database connection failed: " + err.Error())
	}
	sqlDB.SetMaxIdleConns(config.AppConfig.DataBaseConfig.MaxIdle)
	sqlDB.SetMaxOpenConns(config.AppConfig.DataBaseConfig.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(config.AppConfig.DataBaseConfig.MaxLife) * time.Second)

	isInit = true
	logger2.Info(fmt.Sprintf("%s database connected success",
		config.AppConfig.DataBaseConfig.Host+":"+
			fmt.Sprintf("%d", config.AppConfig.DataBaseConfig.Port)+"/"+
			config.AppConfig.DataBaseConfig.Database))

	DB = sqlDB
}

func Close() error {
	logger2.Info("closing db connection...")
	return DB.Close()
}
