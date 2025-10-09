package cmd

import (
	"context"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginautostoplight "github.com/hargeek/gin-auto-stoplight-doc"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go-api-template/common/config"
	"go-api-template/common/logger"
	"go-api-template/handler/middle"
	"go-api-template/handler/routers"
	"go-api-template/internal/store/db"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func RunServer() {
	// 启动gin server
	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", strconv.Itoa(config.AppConfig.ServerConfig.HttpPort)),
		Handler: mainEngine(),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("listen: %s\n", err)
		}
	}()
	logger.Info(fmt.Sprintf("gin server is running on %s", fmt.Sprintf("0.0.0.0:%s", strconv.Itoa(config.AppConfig.ServerConfig.HttpPort))))
	// graceful shutdown
	// Wait for the interrupt signal, shut down all servers gracefully
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("gin server shutdown...")
	// set ctx timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// shutdown gin server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("gin server shutdown error:", err)
	}
	logger.Info("gin server exiting...")
	// 关闭db
	if err := db.Close(); err != nil {
		logger.Fatal("db shutdown error", err)
	}
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

	r.Use(middle.Logger()) // logger middleware
	r.Use(middle.Cors())   // cors middleware

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
