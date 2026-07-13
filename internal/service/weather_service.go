package service

import (
	"context"

	errort "go-api-template/common/error"
)

// WeatherService 定义天气服务接口
type WeatherService interface {
	// QueryWeather 查询指定城市的当前天气描述
	QueryWeather(ctx context.Context, city string) (string, *errort.ApiError)
}
