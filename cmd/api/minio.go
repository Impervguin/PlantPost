package main

import (
	minioclient "PlantSite/internal/infra/minio-client"

	"github.com/spf13/viper"
)

const (
	MinioPrefix       = "minio"
	PlantBucketPrefix = "plant"
	PostBucketPrefix  = "post"
	EndpointKey       = "endpoint"
	MinioUserKey      = "user"
	MinioPasswordKey  = "password"
	BucketKey         = "bucket"
)

func getMinioConfig(bucketPrefix string) *minioclient.MinioConfig {
	if err := ReadInConfig(); err != nil {
		panic(err)
	}
	pref := Key(MinioPrefix, bucketPrefix)
	return &minioclient.MinioConfig{
		Endpoint: viper.GetString(Key(pref, EndpointKey)),
		User:     viper.GetString(Key(pref, MinioUserKey)),
		Password: viper.GetString(Key(pref, MinioPasswordKey)),
		Bucket:   viper.GetString(Key(pref, BucketKey)),
	}
}

func GetPlantMinioConfig() *minioclient.MinioConfig {
	return getMinioConfig(PlantBucketPrefix)
}

func GetPostMinioConfig() *minioclient.MinioConfig {
	return getMinioConfig(PostBucketPrefix)
}
