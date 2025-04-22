package main

import "github.com/spf13/viper"

const (
	ConfigDir  = "./config"
	ConfigName = "api"
	ConfigType = "yaml"
)

func ReadInConfig() error {
	viper.AddConfigPath(ConfigDir)
	viper.SetConfigType(ConfigType)
	viper.SetConfigName(ConfigName)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func Key(prefix, key string) string {
	return prefix + "." + key
}
