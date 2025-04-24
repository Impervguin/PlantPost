package main

import (
	"PlantSite/internal/infra/sqpgx"
	"context"

	"github.com/spf13/viper"
)

const (
	SqpgxPrefix        = "database"
	HostKey            = "host"
	UserKey            = "user"
	PasswordKey        = "password"
	DatabaseKey        = "name"
	PortKey            = "port"
	MaxConnsKey        = "max_connections"
	MaxConnLifeTimeKey = "max_conn_life_time"
)

func GetSqpgxConfig() *sqpgx.SqpgxConfig {
	if err := ReadInConfig(); err != nil {
		panic(err)
	}
	conf, err := sqpgx.NewSqpgxConfig(
		viper.GetString(Key(SqpgxPrefix, UserKey)),
		viper.GetString(Key(SqpgxPrefix, PasswordKey)),
		viper.GetString(Key(SqpgxPrefix, DatabaseKey)),
		viper.GetString(Key(SqpgxPrefix, HostKey)),
		viper.GetInt(Key(SqpgxPrefix, PortKey)),
	)
	if err != nil {
		panic(err)
	}
	conf.MaxConnections = viper.GetInt(Key(SqpgxPrefix, MaxConnsKey))
	conf.MaxConnectionsLifetime = viper.GetDuration(Key(SqpgxPrefix, MaxConnLifeTimeKey))
	return conf
}

func GetSqpgx(ctx context.Context) *sqpgx.SquirrelPgx {
	conf := GetSqpgxConfig()
	sqpgx, err := sqpgx.NewSquirrelPgx(ctx, conf)
	if err != nil {
		panic(err)
	}
	return sqpgx
}
