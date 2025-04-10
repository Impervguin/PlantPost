package miniotest

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Migrate(ctx context.Context, creds *MinioCredentials) error {
	clnt, err := minio.New(creds.GetEndpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(creds.User, creds.Password, ""),
		Secure: false,
	})
	if err != nil {
		return err
	}

	err = clnt.MakeBucket(ctx, creds.Bucket, minio.MakeBucketOptions{})
	if err != nil {
		return err
	}

	return nil
}
