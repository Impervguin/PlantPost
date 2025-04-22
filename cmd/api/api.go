package main

import "github.com/spf13/viper"

const (
	ApiPrefix       = "api"
	ApiUrlPrefixKey = "urlprefix"
)

func GetApiUrlPrefix() string {
	if err := ReadInConfig(); err != nil {
		panic(err)
	}
	return viper.GetString(Key(ApiPrefix, ApiUrlPrefixKey))
}
