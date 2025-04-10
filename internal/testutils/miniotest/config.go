package miniotest

import "github.com/spf13/viper"

type MinioConfig struct {
	User         string
	Password     string
	Bucket       string
	Port         uint16
	Image        string
	MigrationDir string
}

const TestEnvPrefix = "test"

const (
	ConfigTestMinioUserKey         = "minio_user"
	ConfigTestMinioPasswordKey     = "minio_password"
	ConfigTestMinioBucketKey       = "minio_dbucket"
	ConfigTestMinioPortKey         = "minio_port"
	ConfigTestMinioImageKey        = "minio_image"
	ConfigTestMinioMigrationDirKey = "minio_migration_dir"
)

const (
	DefaultTestMinioUser                = "test"
	DefaultTestMinioPassword            = "testtesttest"
	DefaultTestMinioPort         uint16 = 9000
	DefaultTestMinioImage               = "minio/minio:latest"
	DefaultTestMinioMigrationDir        = "migrations/minio"
	DefaultTestMinioBucket              = "test"
)

func GetConfig() *MinioConfig {
	viper.SetEnvPrefix(TestEnvPrefix)
	viper.AutomaticEnv()
	viper.SetDefault(ConfigTestMinioUserKey, DefaultTestMinioUser)
	viper.SetDefault(ConfigTestMinioPasswordKey, DefaultTestMinioPassword)
	viper.SetDefault(ConfigTestMinioBucketKey, DefaultTestMinioBucket)
	viper.SetDefault(ConfigTestMinioPortKey, DefaultTestMinioPort)
	viper.SetDefault(ConfigTestMinioImageKey, DefaultTestMinioImage)
	viper.SetDefault(ConfigTestMinioMigrationDirKey, DefaultTestMinioMigrationDir)

	return &MinioConfig{
		User:         viper.GetString(ConfigTestMinioUserKey),
		Password:     viper.GetString(ConfigTestMinioPasswordKey),
		Bucket:       viper.GetString(ConfigTestMinioBucketKey),
		Port:         viper.GetUint16(ConfigTestMinioPortKey),
		Image:        viper.GetString(ConfigTestMinioImageKey),
		MigrationDir: viper.GetString(ConfigTestMinioMigrationDirKey),
	}
}
