package main

import "github.com/spf13/viper"

const (
	ApiPrefix          = "api"
	ApiUrlPrefixKey    = "urlprefix"
	ApiPortKey         = "port"
	ApiStaticKey       = "static"
	ApiMediaKey        = "media"
	ApiMediaStorageKey = "media-storage"
)

const (
	MediaStorageFs    = "fs"
	MediaStorageMinio = "minio"
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

func GetStaticPath() string {
	if err := ReadInConfig(); err != nil {
		panic(err)
	}
	return viper.GetString(Key(ApiPrefix, ApiStaticKey))
}

func GetMediaPath() string {
	if err := ReadInConfig(); err != nil {
		panic(err)
	}
	return viper.GetString(Key(ApiPrefix, ApiMediaKey))
}

func GetMediaStorage() string {
	if err := ReadInConfig(); err != nil {
		panic(err)
	}
	mediaStorage := viper.GetString(Key(ApiPrefix, ApiMediaStorageKey))
	switch mediaStorage {
	case MediaStorageFs:
		return MediaStorageFs
	case MediaStorageMinio:
		return MediaStorageMinio
	default:
		panic("unknown media storage")
	}
}
