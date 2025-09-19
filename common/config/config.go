package config

import (
	"encoding/json"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
	"sync"
)

const configDir = "config"

var (
	currentStaticCompiledAbsFilename string
	once                             sync.Once
	AppConfig                        *Config
)

var inited = false

func InitViper() {
	if !inited {
		inited = true
		LoadConfig()
	}
}

func LoadConfig() {
	loadConfigFile()
	//loadEnvConfig()
	checkZeroValue(reflect.ValueOf(AppConfig).Elem(), "")
}

func loadConfigFile() {
	once.Do(func() {
		viper.SetConfigName("conf")                                                             // name of config file (without extension)
		viper.AddConfigPath("/" + configDir)                                                    // first try load config from "/config"
		viper.AddConfigPath(configDir)                                                          // second try load config from "config/"
		viper.AddConfigPath(path.Join(currentStaticCompiledAbsFilename, "..", "..", configDir)) // !!!Important, this line config path is compile-time path
		viper.AutomaticEnv()                                                                    // read in environment variables that match
		if err := viper.ReadInConfig(); err != nil {                                            // Handle errors reading the config file
			log.Panic("fail to load config file", err)
		}
		{ // setup default value
			viper.SetDefault("testConfigKey", "testConfigItem")
		}
		AppConfig = new(Config)
		if err := viper.Unmarshal(AppConfig); err != nil {
			log.Panic("common: load config, unmarshal config to struct failed!", err)
		}
		if viper.GetBool("debug") {
			if configJson, err := json.MarshalIndent(AppConfig, "", "  "); err != nil {
				panic(err)
			} else {
				log.Printf("common: load config, marshal config to json: \n%+v", string(configJson))
			}
		}
		checkZeroValue(reflect.ValueOf(AppConfig).Elem(), "")
	})
}

func checkZeroValue(v reflect.Value, parentFieldName string) {
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := typ.Field(i).Name
		// 如果parentFieldName不为空，就将其添加到前面，创建一个层级化的字段名
		fullFieldName := fieldName
		if parentFieldName != "" {
			fullFieldName = parentFieldName + "." + fieldName
		}

		// 检查字段是否未初始化,即零值
		if field.Kind() == reflect.Struct && field.IsZero() {
			log.Panicf("common: check zero value, field key is missing or has zero value: %s", fullFieldName)
		}

		switch field.Kind() {
		case reflect.String:
			if field.Interface().(string) == "" {
				log.Panicf("common: check zero value, field key without value: %s", fullFieldName)
			}
		case reflect.Int, reflect.Int32, reflect.Int64:
			if field.Int() == 0 {
				log.Panicf("common: check zero value, field key without value: %s", fullFieldName)
			}
		case reflect.Bool:
			continue
		case reflect.Struct:
			checkZeroValue(field, fullFieldName) // 对嵌套结构体递归
		case reflect.Slice, reflect.Array: // 检查切片或数组元素是否为零值
			if field.Len() == 0 {
				log.Panicf("common: check zero value, field key without value: %s", fullFieldName)
			}
		default:
			panic("unhandled default case")
		}
	}
}

// loadEnvConfig 从环境变量中加载配置
func loadEnvConfig() {}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsBool 获取环境变量并转换为布尔值
func getEnvAsBool(key string, defaultValue bool) bool {
	value := strings.ToLower(os.Getenv(key))
	switch value {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		return defaultValue
	}
}
