package routers

import (
	"github.com/gin-gonic/gin"
)

var Router router

type router struct{}

func (r *router) InitApiRouter(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		r.InitAuxiliaryRouter(api)
		r.InitHelloRouter(api)
		r.InitWeatherRouter(api)
	}
}
