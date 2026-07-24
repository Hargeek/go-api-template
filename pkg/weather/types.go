package weather

import (
	"fmt"
	"strconv"
)

// currentWeatherResponse 只映射外部协议中当前业务需要的字段，避免协议 DTO 泄漏到 Adapter 之外。
type currentWeatherResponse struct {
	CurrentCondition []currentCondition `json:"current_condition"`
}

type currentCondition struct {
	TemperatureCelsius string               `json:"temp_C"`
	WeatherDescription []weatherDescription `json:"weatherDesc"`
}

type weatherDescription struct {
	Value string `json:"value"`
}

func (c currentCondition) temperatureCelsius() (float64, error) {
	temperature, err := strconv.ParseFloat(c.TemperatureCelsius, 64)
	if err != nil {
		return 0, fmt.Errorf("parse weather temperature %q: %w", c.TemperatureCelsius, err)
	}
	return temperature, nil
}
