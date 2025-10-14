package pgcnt

import "github.com/spf13/viper"

type PostgresConfig struct {
	// From config
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	Image        string `mapstructure:"image"`
	MigrationDir string `mapstructure:"migration_dir"`
	Host         string `mapstructure:"host"`
	Port         uint16 `mapstructure:"port"`

	// From docker setup
	OuterPort *uint16
	OuterHost *string
}

const (
	dbKey = "database"
)

func GetConfig(configPath string) (*PostgresConfig, error) {
	v := viper.New()
	v.SetConfigName("db")
	v.AddConfigPath(configPath)
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config PostgresConfig
	err = v.UnmarshalKey(dbKey, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
