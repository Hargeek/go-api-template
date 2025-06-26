package controller

import (
	errort "go-api-template/common/error"
	res "go-api-template/common/types/response"
	"go-api-template/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Weather *WeatherController

type WeatherController struct {
	Service service.WeatherService
}

func NewWeatherController(s service.WeatherService) *WeatherController {
	return &WeatherController{Service: s}
}

// QueryWeather 查询天气接口
// @Accept      json
// @Produce     json
// @Summary     查询天气
// @Description 查询指定城市天气
// @Tags        Weather API
// @Param       city query    string true "城市名"
// @Success     200  {object} res.CommonApiResponseData
// @Router      /api/v1/weather [get]
func (w *WeatherController) QueryWeather(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		res.ApiResponse(c, http.StatusBadRequest, errort.GeneralError, "city参数不能为空", nil)
		return
	}
	result, err := w.Service.QueryWeather(city)
	if err != nil {
		res.ApiResponse(c, http.StatusInternalServerError, errort.GeneralError, "查询天气失败", gin.H{
			"error": err.Error(),
		})
		return
	}
	res.ApiResponse(c, http.StatusOK, errort.NoError, "查询天气成功", gin.H{
		"data": result,
	})
}
