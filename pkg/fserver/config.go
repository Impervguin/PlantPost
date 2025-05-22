package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type FileServerConfig struct {
	Port int
	Root string
}

const (
	FSPrefix  = "fileserver"
	FSPortKey = "port"
	FSRootKey = "root"
)

func Key(prefix string, key string) string {
	return prefix + "." + key
}

func GetFileServerConfig(configFile string) *FileServerConfig {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	fmt.Println(viper.AllKeys())
	return &FileServerConfig{
		Port: viper.GetInt(Key(FSPrefix, FSPortKey)),
		Root: viper.GetString(Key(FSPrefix, FSRootKey)),
	}
}
