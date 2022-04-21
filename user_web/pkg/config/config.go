package config

import (
	"github.com/spf13/cast"
	viperLib "github.com/spf13/viper"
)

var Config *viperLib.Viper

// SetupConfig 初始化配置
func SetupConfig() {
	Config = viperLib.New()
	Config.SetConfigType("env")
	Config.AddConfigPath(".")
	Config.AutomaticEnv()
	Config.SetConfigName(".env")
	if err := Config.ReadInConfig(); err != nil {
		panic(err)
	}
	Config.WatchConfig()
}

func internalGet(path string) interface{} {
	if !Config.IsSet(path) {
		return nil
	}
	return Config.Get(path)
}

func GetString(path string) string {
	return cast.ToString(internalGet(path))
}

func GetInt(path string) int {
	return cast.ToInt(internalGet(path))
}

func GetBool(path string) bool {
	return cast.ToBool(internalGet(path))
}

func GetFloat64(path string) float64 {
	return cast.ToFloat64(internalGet(path))
}

func Get(path string) interface{} {
	return internalGet(path)
}
