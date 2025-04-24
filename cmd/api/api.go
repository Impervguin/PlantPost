package main

import "github.com/spf13/viper"

const (
	ApiPrefix       = "api"
	ApiUrlPrefixKey = "urlprefix"
	ApiPortKey      = "port"
)

func GetApiUrlPrefix() string {
	if err := ReadInConfig(); err != nil {
		panic(err)
	}
	return viper.GetString(Key(ApiPrefix, ApiUrlPrefixKey))
}

func GetApiPort() uint {
	if err := ReadInConfig(); err != nil {
		panic(err)
	}
	return viper.GetUint(Key(ApiPrefix, ApiPortKey))
}
