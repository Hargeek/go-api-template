package db

import (
	"database/sql"
	"fmt"
	"go-api-template/common/config"
	"go-api-template/common/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

var (
	GORM   *gorm.DB
	DB     *sql.DB
	err    error
	dbOnce sync.Once
)

func newDBWithConfig() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.AppConfig.DataBaseConfig.Host,
		config.AppConfig.DataBaseConfig.Username,
		config.AppConfig.DataBaseConfig.Password,
		config.AppConfig.DataBaseConfig.Database,
		config.AppConfig.DataBaseConfig.Port,
	)
	logLevel := gormLog.Silent
	if config.AppConfig.DataBaseConfig.LogMode {
		logLevel = gormLog.Warn
	}
	newLogger := gormLog.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormLog.Config{
			SlowThreshold: 3 * time.Second, // Slow SQL threshold
			LogLevel:      logLevel,        // Log level
			Colorful:      false,           // Disable color
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

	logger.Info(fmt.Sprintf("%s database connected success",
		config.AppConfig.DataBaseConfig.Host+":"+
			fmt.Sprintf("%d", config.AppConfig.DataBaseConfig.Port)+"/"+
			config.AppConfig.DataBaseConfig.Database))

	DB = sqlDB
}

func GetGORM() *gorm.DB {
	dbOnce.Do(func() {
		newDBWithConfig()
	})
	return GORM
}

func Close() error {
	logger.Info("closing db connection...")
	return DB.Close()
}
