package service

import (
	"context"
	"fmt"

	errort "go-api-template/common/error"
	"go-api-template/common/logger"
	// profile:mtl:start
	"go-api-template/common/metrics"
	// profile:mtl:end
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
		// profile:mtl:start
		metrics.WeatherQueryTotal.WithLabelValues(city, "fail").Inc()
		// profile:mtl:end
		logger.ErrorContext(ctx, "weather query failed",
			"city", city,
			"error", err,
		)
		return "", errort.NewApiError(errort.GeneralError, fmt.Errorf(errort.MsgWeatherQueryFailed, err))
	}
	// profile:mtl:start
	metrics.WeatherQueryTotal.WithLabelValues(city, "success").Inc()
	// profile:mtl:end
	return result, nil
}
