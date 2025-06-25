package service

import "go-api-template/internal/adapter"

type WeatherServiceImpl struct {
	Adapter adapter.WeatherAdapter
}

func (s *WeatherServiceImpl) QueryWeather(city string) (string, error) {
	return s.Adapter.GetWeather(city)
}
