package main

import (
	"time"

	"github.com/spf13/viper"
)

const (
	AuthPrefix           = "auth"
	SessionExpireTimeKey = "session_expire_time"
)

func GetSessionExpireTime() time.Duration {
	if err := ReadInConfig(); err != nil {
		panic(err)
	}
	return viper.GetDuration(Key(AuthPrefix, SessionExpireTimeKey))
}
