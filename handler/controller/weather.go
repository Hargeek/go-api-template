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
//
//	@Accept			json
//	@Produce		json
//	@Summary		查询天气
//	@Description	查询指定城市天气
//	@Tags			Weather API
//	@Param			city	query		string	true	"城市名"
//	@Success		200		{object}	res.CommonApiResponseData
//	@Router			/api/v1/weather [get]
func (w *WeatherController) QueryWeather(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		res.ApiResponse(c, http.StatusBadRequest, errort.ParamMissing, "city参数不能为空", nil)
		return
	}
	result, apiErr := w.Service.QueryWeather(c.Request.Context(), city)
	if apiErr != nil {
		res.ApiResponse(c, http.StatusInternalServerError, apiErr.Code, apiErr.Msg, nil)
		return
	}
	res.ApiResponse(c, http.StatusOK, errort.NoError, "查询天气成功", gin.H{
		"data": result,
	})
}
