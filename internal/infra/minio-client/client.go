package minioclient

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	minio.Client
	bucket string
}

func NewMinioClient(conf *MinioConfig) (*MinioClient, error) {
	client, err := minio.New(
		conf.Endpoint,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				conf.User,
				conf.Password,
				"",
			),
			Secure: false,
		},
	)
	if err != nil {
		return nil, err
	}
	return &MinioClient{Client: *client, bucket: conf.Bucket}, nil
}

func (client *MinioClient) GetBucket() string {
	return client.bucket
}
