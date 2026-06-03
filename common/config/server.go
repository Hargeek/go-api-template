package config

type ServerConfig struct {
	HttpPort   int `mapstructure:"http_port"`   // HTTP API 端口
	MetricPort int `mapstructure:"metric_port"` //  metrics 端口
}
