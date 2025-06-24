package routers

import (
	"github.com/gin-gonic/gin"
	"go-api-template/handler/controller"
)

func (r *router) InitHelloRouter(api *gin.RouterGroup) {
	api.GET("/hello", controller.Hello.HelloController)
}
