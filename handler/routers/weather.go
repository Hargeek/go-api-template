package routers

import (
	"go-api-template/handler/controller"

	"github.com/gin-gonic/gin"
)

func (r *router) InitWeatherRouter(api *gin.RouterGroup) {
	api.GET("/weather", controller.Weather.QueryWeather)
}
