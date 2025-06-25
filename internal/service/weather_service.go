package service

// WeatherService 定义天气服务接口

type WeatherService interface {
	QueryWeather(city string) (string, error)
}
