package adapter

import (
	"fmt"
)

// WeatherAdapterImpl 是 WeatherAdapter 的一个简单实现
// 实际项目中可调用第三方API，这里仅做演示

type WeatherAdapterImpl struct {
	ApiKey string
}

func NewWeatherAdapterImpl(apiKey string) *WeatherAdapterImpl {
	return &WeatherAdapterImpl{ApiKey: apiKey}
}

func (w *WeatherAdapterImpl) GetWeather(city string) (string, error) {
	return fmt.Sprintf("%s: 晴，25°C", city), nil
}
