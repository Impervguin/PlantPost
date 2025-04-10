package tests

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetEnvPrefix(TestEnvPrefix)
	viper.SetDefault(TestWorkingDirKey, DefaultWorkingDir)
	viper.AutomaticEnv()
}

const TestEnvPrefix = "test"

const TestWorkingDirKey = "twd"

const DefaultWorkingDir = "."

func GetTestWorkingDir() string {
	return viper.GetString(TestWorkingDirKey)
}
