package miniotest

import "github.com/spf13/viper"

type MinioConfig struct {
	User         string
	Password     string
	Bucket       string
	Host         string
	Port         uint16
	Image        string
	MigrationDir string
	External     bool
}

const TestEnvPrefix = "test"

const (
	ConfigTestMinioUserKey         = "minio_user"
	ConfigTestMinioPasswordKey     = "minio_password"
	ConfigTestMinioBucketKey       = "minio_dbucket"
	ConfigTestMinioHostKey         = "minio_host"
	ConfigTestMinioPortKey         = "minio_port"
	ConfigTestMinioImageKey        = "minio_image"
	ConfigTestMinioMigrationDirKey = "minio_migration_dir"
	ConfigTestMinioExternalKey     = "minio_external"
)

const (
	DefaultTestMinioUser                = "test"
	DefaultTestMinioPassword            = "testtesttest"
	DefaultTestMinioPort         uint16 = 9000
	DefaultTestMinioHost                = "localhost"
	DefaultTestMinioImage               = "minio/minio:latest"
	DefaultTestMinioMigrationDir        = "migrations/minio"
	DefaultTestMinioBucket              = "test"
	DefaultTestMinioExternal            = false
)

func GetConfig() *MinioConfig {
	viper.SetEnvPrefix(TestEnvPrefix)
	viper.AutomaticEnv()
	viper.SetDefault(ConfigTestMinioUserKey, DefaultTestMinioUser)
	viper.SetDefault(ConfigTestMinioPasswordKey, DefaultTestMinioPassword)
	viper.SetDefault(ConfigTestMinioBucketKey, DefaultTestMinioBucket)
	viper.SetDefault(ConfigTestMinioHostKey, DefaultTestMinioHost)
	viper.SetDefault(ConfigTestMinioPortKey, DefaultTestMinioPort)
	viper.SetDefault(ConfigTestMinioImageKey, DefaultTestMinioImage)
	viper.SetDefault(ConfigTestMinioMigrationDirKey, DefaultTestMinioMigrationDir)
	viper.SetDefault(ConfigTestMinioExternalKey, DefaultTestMinioExternal)

	return &MinioConfig{
		User:         viper.GetString(ConfigTestMinioUserKey),
		Password:     viper.GetString(ConfigTestMinioPasswordKey),
		Bucket:       viper.GetString(ConfigTestMinioBucketKey),
		Host:         viper.GetString(ConfigTestMinioHostKey),
		Port:         viper.GetUint16(ConfigTestMinioPortKey),
		Image:        viper.GetString(ConfigTestMinioImageKey),
		MigrationDir: viper.GetString(ConfigTestMinioMigrationDirKey),
		External:     viper.GetBool(ConfigTestMinioExternalKey),
	}
}
