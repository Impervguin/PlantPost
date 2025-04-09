package pgtest

import "github.com/spf13/viper"

type PostgresConfig struct {
	User         string
	Password     string
	Database     string
	Port         uint16
	Image        string
	MigrationDir string
}

const (
	ConfigTestPostgresUserKey         = "postgres_user"
	ConfigTestPostgresPasswordKey     = "postgres_password"
	ConfigTestPostgresDatabaseKey     = "postgres_database"
	ConfigTestPostgresPortKey         = "postgres_port"
	ConfigTestPostgresImageKey        = "postgres_image"
	ConfigTestPostgresMigrationDirKey = "postgres_migration_dir"
)

const (
	DefaultTestPostgresUser                = "test"
	DefaultTestPostgresPassword            = "test"
	DefaultTestPostgresDatabase            = "test"
	DefaultTestPostgresPort         uint16 = 5432
	DefaultTestPostgresImage               = "postgres:latest"
	DefaultTestPostgresMigrationDir        = "migrations/postgr"
)
const TestEnvPrefix = "test"

func GetConfig() *PostgresConfig {

	viper.SetDefault(ConfigTestPostgresUserKey, DefaultTestPostgresUser)
	viper.SetDefault(ConfigTestPostgresPasswordKey, DefaultTestPostgresPassword)
	viper.SetDefault(ConfigTestPostgresDatabaseKey, DefaultTestPostgresDatabase)
	viper.SetDefault(ConfigTestPostgresPortKey, DefaultTestPostgresPort)
	viper.SetDefault(ConfigTestPostgresImageKey, DefaultTestPostgresImage)
	viper.SetDefault(ConfigTestPostgresMigrationDirKey, DefaultTestPostgresMigrationDir)
	viper.SetEnvPrefix(TestEnvPrefix)
	viper.AutomaticEnv()

	return &PostgresConfig{
		User:         viper.GetString(ConfigTestPostgresUserKey),
		Password:     viper.GetString(ConfigTestPostgresPasswordKey),
		Database:     viper.GetString(ConfigTestPostgresDatabaseKey),
		Port:         viper.GetUint16(ConfigTestPostgresPortKey),
		Image:        viper.GetString(ConfigTestPostgresImageKey),
		MigrationDir: viper.GetString(ConfigTestPostgresMigrationDirKey),
	}
}
