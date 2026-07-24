package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-api-template/common/config"
	"go-api-template/common/logger"
	_ "go-api-template/common/metrics" // 触发 init()，注册 build_info 指标
	"go-api-template/handler/middle"
	"go-api-template/handler/routers"
	"go-api-template/internal/store/db"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginautostoplight "github.com/hargeek/gin-auto-stoplight-doc"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

func RunServer() {
	// 启动业务 HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", config.AppConfig.ServerConfig.HttpPort),
		Handler: mainEngine(),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("gin server listen failed", "error", err)
		}
	}()
	logger.Info("gin server started", "address", srv.Addr)

	// 启动 metrics server
	metricsSrv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", config.AppConfig.ServerConfig.MetricPort),
		Handler: promhttp.Handler(),
	}
	go func() {
		if err := metricsSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("metrics server listen failed", "error", err)
		}
	}()
	logger.Info("metrics server started", "address", metricsSrv.Addr, "path", "/metrics")

	// 优雅退出：等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("gin server shutdown error", "error", err)
	}
	if err := metricsSrv.Shutdown(ctx); err != nil {
		logger.Error("metrics server shutdown error", "error", err)
	}
	if err := db.Close(); err != nil {
		logger.Error("db shutdown error", "error", err)
	}
	// flush 未导出的 Span（必须在最后，确保所有 Span 都已记录完毕）
	if TelemetryShutdown != nil {
		if err := TelemetryShutdown(ctx); err != nil {
			logger.Error("telemetry shutdown error", "error", err)
		}
	}
	logger.Info("all servers exited")
}

func mainEngine() *gin.Engine {
	if viper.GetBool("debug") {
		gin.SetMode(gin.DebugMode) // Set Gin to Debug mode
	} else {
		gin.SetMode(gin.ReleaseMode) // Set Gin to Release mode
	}

	r := gin.New()
	// Recovery middleware, to recover from any panics and write a 500 if there was one.
	r.Use(gin.Recovery())

	if viper.GetBool("debug") {
		pprof.Register(r) // Automatically add routers for net/http/pprof only if config enables it
	}

	r.Use(middle.Trace())   // 链路追踪（最先注册，确保后续中间件和 Handler 都在同一 Span 下）
	r.Use(middle.Metrics()) // Prometheus 指标采集
	r.Use(middle.Logger())  // 访问日志
	r.Use(middle.Cors())    // 跨域

	// UseRawPath = true, 保留原始路径
	r.UseRawPath = true
	// swagger
	routers.Router.RegisterSwagger(r)
	// redoc
	ginautostoplight.Register(r)
	// init router
	routers.Router.InitApiRouter(r)

	return r
}
