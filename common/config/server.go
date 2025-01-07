package config

type ServerConfig struct {
	HttpPort int `mapstructure:"http_port"` // http api端口
}
