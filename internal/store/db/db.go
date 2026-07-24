package db

import (
	"database/sql"
	"log"
	"os"
	"sync"
	"time"

	gormSqlite "github.com/glebarez/sqlite"

	"go-api-template/common/config"
	"go-api-template/common/logger"

	// "gorm.io/driver/postgres" // 切换 PostgreSQL 时取消注释
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
	// profile:mtl:start
	"gorm.io/plugin/opentelemetry/tracing"
	// profile:mtl:end
)

var (
	GORM   *gorm.DB
	DB     *sql.DB
	dbOnce sync.Once
)

func newDBWithConfig() {
	sqlite := config.AppConfig.SQLiteConfig
	// pg := config.AppConfig.PostgresConfig // 切换 PostgreSQL 时取消注释

	logLevel := gormLog.Silent
	if sqlite.LogMode {
		logLevel = gormLog.Warn
	}
	gormLogger := gormLog.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormLog.Config{
			SlowThreshold: 3 * time.Second,
			LogLevel:      logLevel,
			Colorful:      false,
		},
	)

	var err error
	GORM, err = gorm.Open(gormSqlite.Open(sqlite.Path), &gorm.Config{Logger: gormLogger})
	if err != nil {
		panic("connecting sqlite failed: " + err.Error())
	}
	logger.Info("database connected", "driver", "sqlite", "path", sqlite.Path)

	// 切换 PostgreSQL 时替换上方 sqlite 连接块为：
	// dsn := fmt.Sprintf(
	// 	"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable TimeZone=Asia/Shanghai",
	// 	pg.Host, pg.Port, pg.Database, pg.Username, pg.Password,
	// )
	// GORM, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
	// if err != nil {
	// 	panic("connecting postgres failed: " + err.Error())
	// }
	// logger.Info(fmt.Sprintf("database connected: driver=postgres host=%s:%d/%s", pg.Host, pg.Port, pg.Database))

	// profile:mtl:start
	// 注册 GORM OTEL Plugin：每条 SQL 自动生成子 Span，挂载在当前请求链路下。
	// WithoutMetrics() 避免与 Prometheus 指标重复注册
	if err := GORM.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		panic("register gorm otel plugin failed: " + err.Error())
	}
	// profile:mtl:end

	sqlDB, err := GORM.DB()
	if err != nil {
		panic("get database connection failed: " + err.Error())
	}
	sqlDB.SetMaxIdleConns(sqlite.MaxIdle)
	sqlDB.SetMaxOpenConns(sqlite.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(sqlite.MaxLife) * time.Second)
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
	if DB == nil {
		return nil
	}
	return DB.Close()
}
