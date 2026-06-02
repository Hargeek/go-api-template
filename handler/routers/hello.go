package routers

import (
	"go-api-template/handler/controller"

	"github.com/gin-gonic/gin"
)

func (r *router) InitHelloRouter(api *gin.RouterGroup) {
	api.GET("/hello", controller.Hello.HelloController)
}
