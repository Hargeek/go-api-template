package config

// WeatherConfig 定义天气示例使用的外部服务配置。
type WeatherConfig struct {
	BaseURL        string `mapstructure:"base_url"`        // 天气服务基础地址，例如 https://wttr.in
	TimeoutSeconds int    `mapstructure:"timeout_seconds"` // 单次 HTTP 请求超时时间，单位秒
}
