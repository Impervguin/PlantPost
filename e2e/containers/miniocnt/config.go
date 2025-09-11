package miniocnt

import "github.com/spf13/viper"

type MinioConfig struct {
	// From config
	User     string   `mapstructure:"user"`
	Password string   `mapstructure:"password"`
	Buckets  []string `mapstructure:"buckets"`
	Image    string   `mapstructure:"image"`
	Host     string   `mapstructure:"host"`
	Port     uint16   `mapstructure:"port"`

	// From docker setup
	OuterPort *uint16
	OuterHost *string
}

const (
	minioKey = "minio"
)

func GetConfig(configPath string) (*MinioConfig, error) {
	v := viper.New()
	v.SetConfigName("minio")
	v.AddConfigPath(configPath)
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config MinioConfig
	err = v.UnmarshalKey(minioKey, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
