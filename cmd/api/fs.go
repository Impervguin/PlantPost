package main

import "github.com/spf13/viper"

const (
	FSPrefix            = "fs"
	FSPlantBucketPrefix = "plant"
	FSPostBucketPrefix  = "post"
	FSRootKey           = "root"
	FSBucketKey         = "bucket"
)

func GetFsBucket(bucketPrefix string) string {
	if err := ReadInConfig(); err != nil {
		panic(err)
	}
	pref := Key(FSPrefix, bucketPrefix)
	return viper.GetString(Key(pref, FSBucketKey))
}

func GetFsRoot() string {
	if err := ReadInConfig(); err != nil {
		panic(err)
	}
	return viper.GetString(Key(FSPrefix, FSRootKey))
}
