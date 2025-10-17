package miniotest

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var ErrBucketDoesNotExist = fmt.Errorf("bucket does not exist")

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

func MigrateDown(ctx context.Context, creds *MinioCredentials) error {
	clnt, err := minio.New(creds.GetEndpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(creds.User, creds.Password, ""),
		Secure: false,
	})
	if err != nil {
		return err
	}
	return clnt.RemoveBucketWithOptions(ctx, creds.Bucket, minio.RemoveBucketOptions{ForceDelete: true})
}

func CheckBucketExists(ctx context.Context, creds *MinioCredentials) error {
	clnt, err := minio.New(creds.GetEndpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(creds.User, creds.Password, ""),
		Secure: false,
	})
	if err != nil {
		return err
	}
	f, err := clnt.BucketExists(ctx, creds.Bucket)
	if err != nil {
		return err
	}
	if !f {
		return ErrBucketDoesNotExist
	}
	return nil
}

func CleanUpBucket(ctx context.Context, creds *MinioCredentials) error {
	clnt, err := minio.New(creds.GetEndpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(creds.User, creds.Password, ""),
		Secure: false,
	})
	if err != nil {
		return err
	}
	clnt.RemoveObjects(
		ctx,
		creds.Bucket,
		clnt.ListObjects(ctx, creds.Bucket, minio.ListObjectsOptions{}),
		minio.RemoveObjectsOptions{},
	)
	return nil
}
