package routers

import (
	"github.com/gin-gonic/gin"
	"go-api-template/handler/controller"
)

func (r *router) InitAuxiliaryRouter(api *gin.RouterGroup) {
	api.GET("/health", controller.Auxiliary.GetHealthy)
	api.GET("/delayed-health", controller.Auxiliary.GetDelayedHealthy)
	api.GET("/echo-get", controller.Auxiliary.EchoAnyGet)
	api.POST("/echo-post", controller.Auxiliary.EchoAnyPost)
}
