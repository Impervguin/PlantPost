package appcnt

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
	DockerFilePath     string `mapstructure:"dockerfile"`
	BuildContext       string `mapstructure:"build_context"`
	AppConfigPath      string `mapstructure:"app_config"`
	AdminUser          string `mapstructure:"admin_user"`
	AdminPassword      string `mapstructure:"admin_password"`
	Host               string `mapstructure:"host"`
	Port               uint16 `mapstructure:"port"`
	ExternalDataSource bool   `mapstructure:"external_data_source"` // If Databases is external

	// From docker setup
	OuterPort *uint16
	OuterHost *string

	// From other containers
	DbHost     *string
	DbPort     *uint16
	DbName     *string
	DbUser     *string
	DbPassword *string

	MinioHost        *string
	MinioPort        *uint16
	MinioUser        *string
	MinioPassword    *string
	MinioPostBucket  *string
	MinioPlantBucket *string
}

const (
	appKey = "app"
)

func GetConfig(configPath string) (*AppConfig, error) {
	v := viper.New()
	v.SetConfigName("app")
	v.AddConfigPath(configPath)
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config AppConfig
	err = v.UnmarshalKey(appKey, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *AppConfig) GetApiYamlReader(ApiPort uint16) io.Reader {
	return strings.NewReader(
		fmt.Sprintf(
			`
database:
  name: %s
  user: %s
  password: %s
  host: %s
  port: %d
  max_connections: 10
  max_conn_life_time: 10m 
hasher:
  hash_cost: 31

api:
  urlprefix: api
  port: %d
  static: ./internal/view/static/
  media: /media


minio:
  post:
    endpoint: %s:%d
    user: %s
    password: %s
    bucket: %s
  plant:
    endpoint: %s:%d
    user: %s
    password: %s
    bucket: %s

admins:
  - login: %s
    password: %s

auth:
  session_expire_time: 1h

log:
  console_level: info
  file_level: info
  dir: /logs
  type: dev
  file_factory: every_day
  file_type: json`,
			*c.DbName,
			*c.DbUser,
			*c.DbPassword,
			*c.DbHost,
			*c.DbPort,
			ApiPort,
			*c.MinioHost,
			*c.MinioPort,
			*c.MinioUser,
			*c.MinioPassword,
			*c.MinioPostBucket,
			*c.MinioHost,
			*c.MinioPort,
			*c.MinioUser,
			*c.MinioPassword,
			*c.MinioPlantBucket,
			c.AdminUser,
			c.AdminPassword,
		),
	)
}
