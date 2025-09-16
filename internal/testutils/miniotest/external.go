package miniotest

import (
	"context"
	"fmt"
)

func NewTestExternalMinio(ctx context.Context, config *MinioConfig) (*MinioCredentials, error) {
	if !config.External {
		return nil, fmt.Errorf("minio is not configured")
	}
	creds := NewMinioCredentials(config.User, config.Password, config.Bucket, config.Host, config.Port)
	return creds, nil
}
