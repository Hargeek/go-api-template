package routers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-api-template/common/logger"
	"go-api-template/docs"
	swagDoc "go-api-template/docs"
	"go-api-template/static"
)

func (r *router) RegisterSwagger(e *gin.Engine) {
	// programmatically set swagger info description
	content, err := static.GetErrorCodeEmbed().ReadFile("error_code/error_code.md")
	if err != nil {
		logger.Panic("error reading error code file: ", err)
		return
	}
	docs.SwaggerInfo.Description = string(content)
	e.GET("/api/v1/swagger/*any", func(c *gin.Context) {
		protocol := "http"
		if forwardedProto := c.Request.Header.Get("X-Forwarded-Proto"); forwardedProto == "https" {
			protocol = "https"
		} else if c.Request.TLS != nil {
			protocol = "https"
		} else {
			protocol = "http"
		}
		host := c.Request.Host
		swagDoc.SwaggerInfo.Host = host
		swagDoc.SwaggerInfo.Schemes = []string{protocol}
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})
}
