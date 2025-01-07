package config

type DataBaseConfig struct {
	Host     string `mapstructure:"host"`          // 数据库地址
	Port     int64  `mapstructure:"port"`          // 数据库端口
	Database string `mapstructure:"db_name"`       // 数据库名
	Username string `mapstructure:"db_user"`       // 数据库用户名
	Password string `mapstructure:"db_password"`   // 数据库密码
	LogMode  bool   `mapstructure:"log_mode"`      // 是否开启日志模式
	MaxIdle  int    `mapstructure:"max_idle_conn"` // 最大空闲连接
	MaxOpen  int    `mapstructure:"max_open_conn"` // 最大连接数
	MaxLife  int    `mapstructure:"max_life_time"` // 最大生存时间
}
