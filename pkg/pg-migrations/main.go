package main

import (
	"PgMigrations/migrator"
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/viper"
)

const ConfigDir = "./config"
const ConfigName = "pgmigr"
const ConfigType = "yaml"

const (
	MigrationDirKey = "migration_dir"
	DatabaseKey     = "database"
	UserKey         = "user"
	PasswordKey     = "password"
	HostKey         = "host"
	PortKey         = "port"
)

func initMigrator() (*migrator.Migrator, error) {
	viper.AddConfigPath(ConfigDir)
	viper.SetConfigType(ConfigType)
	viper.SetConfigName(ConfigName)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := migrator.Config{
		MigrationDir: viper.GetString(MigrationDirKey),
		Database:     viper.GetString(DatabaseKey),
		User:         viper.GetString(UserKey),
		Password:     viper.GetString(PasswordKey),
		Host:         viper.GetString(HostKey),
		Port:         viper.GetInt(PortKey),
	}
	migrator, err := migrator.NewMigrator(&config)
	if err != nil {
		return nil, err
	}
	return migrator, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: pg-migrations <up|down>")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "up":
		migrator, err := initMigrator()
		if err != nil {
			panic(err)
		}
		err = migrator.Up()
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No change")
			os.Exit(0)
		} else if err != nil {
			panic(err)
		}
	case "down":
		migrator, err := initMigrator()
		if err != nil {
			panic(err)
		}
		err = migrator.Down()
		if err != nil {
			panic(err)
		}
	default:
		fmt.Println("Invalid command")
	}
}
