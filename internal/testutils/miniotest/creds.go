package miniotest

import "fmt"

type MinioCredentials struct {
	User     string
	Password string
	Bucket   string
	Host     string
	Port     uint16
}

func NewMinioCredentials(user, password, bucket, host string, port uint16) *MinioCredentials {
	return &MinioCredentials{
		User:     user,
		Password: password,
		Bucket:   bucket,
		Host:     host,
		Port:     port,
	}
}

func (c *MinioCredentials) GetEndpoint() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
