package miniocnt

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Migrate(ctx context.Context, db *MinioConfig) error {
	clnt, err := minio.New(fmt.Sprintf("%s:%d", *db.OuterHost, *db.OuterPort), &minio.Options{
		Creds:  credentials.NewStaticV4(db.User, db.Password, ""),
		Secure: false,
	})
	if err != nil {
		return err
	}
	for _, bucket := range db.Buckets {
		if err := clnt.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func MigrateDown(ctx context.Context, db *MinioConfig) error {
	clnt, err := minio.New(fmt.Sprintf("%s:%d", *db.OuterHost, *db.OuterPort), &minio.Options{
		Creds:  credentials.NewStaticV4(db.User, db.Password, ""),
		Secure: false,
	})
	if err != nil {
		return err
	}
	for _, bucket := range db.Buckets {
		if err := clnt.RemoveBucketWithOptions(ctx, bucket, minio.RemoveBucketOptions{ForceDelete: true}); err != nil {
			return err
		}

	}
	return nil
}
