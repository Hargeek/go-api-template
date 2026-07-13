package service

import (
	"context"
	"fmt"

	errort "go-api-template/common/error"
	"go-api-template/common/logger"
	"go-api-template/common/metrics"
	"go-api-template/internal/adapter"
)

type WeatherServiceImpl struct {
	Adapter adapter.WeatherAdapter
}

func NewWeatherServiceImpl(adapter adapter.WeatherAdapter) *WeatherServiceImpl {
	return &WeatherServiceImpl{Adapter: adapter}
}

func (s *WeatherServiceImpl) QueryWeather(ctx context.Context, city string) (string, *errort.ApiError) {
	result, err := s.Adapter.GetWeather(ctx, city)
	if err != nil {
		metrics.WeatherQueryTotal.WithLabelValues(city, "fail").Inc()
		logger.Error(fmt.Sprintf(errort.MsgWeatherQueryFailed, err))
		return "", errort.NewApiError(errort.GeneralError, fmt.Errorf(errort.MsgWeatherQueryFailed, err))
	}
	metrics.WeatherQueryTotal.WithLabelValues(city, "success").Inc()
	return result, nil
}
