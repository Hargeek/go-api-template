package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	errort "go-api-template/common/error"
	"go-api-template/common/types"
	"go-api-template/store/db"
	"io"
	"net/http"
	"strings"
	"time"
)

var Auxiliary auxiliary

type auxiliary struct{}

// GetHealthy 获取健康检查状态
//
// @Accept      json
// @Produce     json
// @Summary     健康检查接口
// @Description 健康检查接口
// @Tags        Auxiliary API
// @Success     200    {object} types.CommonApiResponse
// @Router      /api/v1/health [get]
func (*auxiliary) GetHealthy(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": errort.NoError,
		"data": gin.H{
			"service_name": types.ServiceName,
			"branch":       types.Branch,
			"env":          viper.GetString("env"),
			"revision":     types.Revision,
			"build_date":   types.BuildDate,
			"build_user":   types.BuildUser,
			"go_version":   types.GoVersion,
			"db_status":    db.DB.Stats(),
		},
		"msg": "It is healthy!",
	})
}

// GetDelayedHealthy 延迟响应测试接口
//
// @Accept      json
// @Produce     json
// @Summary     延迟响应测试接口
// @Description 延迟响应测试接口
// @Tags        Auxiliary API
// @Param       delay_sec query    int true "delay time(second)"
// @Success     200       {object} types.CommonApiResponse
// @Router      /api/v1/delayed-health [get]
func (*auxiliary) GetDelayedHealthy(ctx *gin.Context) {
	params := new(struct {
		DelaySec int `form:"delay_sec" binding:"required"`
	})
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	time.Sleep(time.Duration(params.DelaySec) * time.Second)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "It is " + ctx.Query("delay_sec") + " second delayed healthy!",
		"data": nil,
		"code": errort.NoError,
	})
}

// EchoAnyGet 回显请求信息(get)
//
// @Accept      json
// @Produce     json
// @Summary     回显请求信息(get)
// @Description 回显请求信息(get)
// @Tags        Auxiliary API
// @Success     200 {object} types.CommonApiResponse
// @Router      /api/v1/echo-get [get]
func (*auxiliary) EchoAnyGet(ctx *gin.Context) {
	headers := make(map[string]string)
	for key, value := range ctx.Request.Header {
		headers[key] = strings.Join(value, ",")
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": errort.NoError,
		"data": gin.H{
			"client_ip":       ctx.ClientIP(),
			"remote_addr":     ctx.Request.RemoteAddr,
			"request_uri":     ctx.Request.RequestURI,
			"request_path":    ctx.Request.URL.Path,
			"x-forwarded-for": ctx.GetHeader("X-Forwarded-For"),
			"request_headers": headers,
		},
		"msg": "It is echo get!",
	})
}

// EchoAnyPost 回显请求信息(post)
//
// @Accept      json
// @Produce     json
// @Summary     回显请求信息(post)
// @Description 回显请求信息(post)
// @Tags        Auxiliary API
// @Param       params body     interface{} true "Request Body""
// @Success     200 {object} types.CommonApiResponse
// @Router      /api/v1/echo-post [post]
func (*auxiliary) EchoAnyPost(ctx *gin.Context) {
	headers := make(map[string]string)
	for key, value := range ctx.Request.Header {
		headers[key] = strings.Join(value, ",")
	}
	bodyBytes, _ := io.ReadAll(ctx.Request.Body)
	defer ctx.Request.Body.Close()
	ctx.JSON(http.StatusOK, gin.H{
		"code": errort.NoError,
		"data": gin.H{
			"client_ip":        ctx.ClientIP(),
			"remote_addr":      ctx.Request.RemoteAddr,
			"request_uri":      ctx.Request.RequestURI,
			"request_path":     ctx.Request.URL.Path,
			"x-forwarded-for":  ctx.GetHeader("X-Forwarded-For"),
			"request_headers":  headers,
			"request_body_str": string(bodyBytes),
		},
		"msg": "It is echo post!",
	})
}
