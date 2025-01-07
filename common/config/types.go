package config

type Config struct {
	EnvConfig      string         `mapstructure:"env"`      // 环境变量
	ServerConfig   ServerConfig   `mapstructure:"server"`   // 服务配置
	DataBaseConfig DataBaseConfig `mapstructure:"database"` // 数据库配置
}
