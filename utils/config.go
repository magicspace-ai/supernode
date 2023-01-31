package utils

import (
	"github.com/spf13/viper"
)

var isInitialized = false
var vconfig *viper.Viper= nil

func initConfig() {

	v := viper.New()
	v.SetConfigName("main")    // name of config file (without extension)
	v.SetConfigType("toml")      // REQUIRED if the config file does not have the extension in the name
	v.AddConfigPath("./config/") // path to look for the config file in

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			HandleError(err, "Config file was not found at config/main.toml", true)
		} else {
			HandleError(err, "Error on loading config/main.toml", true)
		}
	}

	vconfig = v
}

func GetConfigs() *viper.Viper {

	if vconfig == nil {
		initConfig()
	}

	return vconfig
}

/*
 * getConfig by key
 */
func GetConfig(key string, fallback interface{}) interface{} {
	value := GetConfigs().Get(key)
	if value == nil {
		return fallback
	}
	return value
}
