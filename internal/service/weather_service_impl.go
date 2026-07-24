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
	weatherAdapter adapter.WeatherAdapter
}

func NewWeatherServiceImpl(weatherAdapter adapter.WeatherAdapter) *WeatherServiceImpl {
	return &WeatherServiceImpl{weatherAdapter: weatherAdapter}
}

func (s *WeatherServiceImpl) QueryWeather(ctx context.Context, city string) (string, *errort.ApiError) {
	result, err := s.weatherAdapter.GetWeather(ctx, city)
	if err != nil {
		metrics.WeatherQueryTotal.WithLabelValues(city, "fail").Inc()
		logger.ErrorContext(ctx, "weather query failed",
			"city", city,
			"error", err,
		)
		return "", errort.NewApiError(errort.GeneralError, fmt.Errorf(errort.MsgWeatherQueryFailed, err))
	}
	metrics.WeatherQueryTotal.WithLabelValues(city, "success").Inc()
	return result, nil
}
