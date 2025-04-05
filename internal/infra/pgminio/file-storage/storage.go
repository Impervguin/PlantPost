package filestorage

import (
	"PlantSite/internal/infra/sqpgx"
	"PlantSite/internal/models"
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

type PgMinioStorage struct {
	minioClient *minio.Client
	bucketName  string
	db          *sqpgx.SquirrelPgx
}

func pgConfigString() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s port=%d host=%s pool_max_conns=%d pool_max_conn_lifetime=%s",
		viper.GetString(ConfigPostgresUserKey),
		viper.GetString(ConfigPostgresPasswordKey),
		viper.GetString(ConfigPostgresDbNameKey),
		viper.GetInt(ConfigPostgresPortKey),
		viper.GetString(ConfigPostgresHostKey),
		viper.GetInt(ConfigMaxConnectionsKey),
		viper.GetString(ConfigMaxConnectionLifetimeKey),
	)
}

func NewPgMinioStorage(ctx context.Context, bucketName string) (*PgMinioStorage, error) {
	minioClient, err := minio.New(
		viper.GetString(ConfigMinioEndpoint),
		&minio.Options{
			Creds: credentials.NewStaticV4(
				viper.GetString(ConfigMinioUser),
				viper.GetString(ConfigMinioPassword),
				"",
			),
			Secure: false,
		},
	)
	if err != nil {
		return nil, err
	}

	connStr := pgConfigString()
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	return &PgMinioStorage{
		minioClient: minioClient,
		bucketName:  bucketName,
		db:          sqpgx.NewSquirrelPgx(pool),
	}, nil
}

func (storage *PgMinioStorage) UploadFile(ctx context.Context, data models.FileData) (*models.File, error) {
	f, err := models.NewFile(data.Name)
	if err != nil {
		return nil, err
	}

	info, err := storage.minioClient.PutObject(ctx,
		storage.bucketName,
		f.URL,
		data.Reader,
		-1,
		minio.PutObjectOptions{ContentType: data.ContentType},
	)
	f.URL = info.Key
	if err != nil {
		return nil, err
	}
	_, err = storage.db.Insert(ctx, squirrel.Insert("file").
		Columns("id", "name", "url", "created_at").
		Values(f.ID, f.Name, f.URL, f.CreatedAt))
	if err != nil {
		return nil, err
	}
	f.URL = fmt.Sprintf("%s/%s", storage.bucketName, f.URL)

	return f, nil
}

func (storage *PgMinioStorage) Download(ctx context.Context, fileID uuid.UUID) (*models.FileData, error) {
	f, err := storage.Get(ctx, fileID)
	if err != nil {
		return nil, err
	}
	attrs, err := storage.minioClient.StatObject(ctx, storage.bucketName, f.URL, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	data, err := storage.minioClient.GetObject(ctx, storage.bucketName, f.URL, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer data.Close()

	body, err := io.ReadAll(data)
	if err != nil {
		return nil, err
	}
	return &models.FileData{
		Name:        f.Name,
		Reader:      bytes.NewReader(body),
		ContentType: attrs.ContentType,
	}, nil
}

func (storage *PgMinioStorage) DeleteFile(ctx context.Context, fileID uuid.UUID) error {
	f, err := storage.Get(ctx, fileID)
	if err != nil {
		return err
	}
	_, err = storage.db.Delete(ctx, squirrel.Delete("file").Where(squirrel.Eq{"id": fileID}))

	errminio := storage.minioClient.RemoveObject(ctx, storage.bucketName, f.URL, minio.RemoveObjectOptions{})
	if errminio != nil || err != nil {
		return fmt.Errorf("error deleting file from database or minio: %w %w", err, errminio)
	}
	return nil
}

func (storage *PgMinioStorage) Update(ctx context.Context, fileID uuid.UUID, data *models.FileData) (*models.File, error) {
	f, err := storage.Get(ctx, fileID)
	if err != nil {
		return nil, err
	}

	f.Name = data.Name
	f.CreatedAt = time.Now()

	info, err := storage.minioClient.PutObject(ctx, storage.bucketName, f.URL, data.Reader, -1, minio.PutObjectOptions{})
	if err != nil {
		return nil, err
	}
	f.URL = info.Key
	_, err = storage.db.Update(ctx, squirrel.Update("file").
		Set("name", f.Name).
		Set("url", f.URL).
		Set("created_at", f.CreatedAt).
		Where(squirrel.Eq{"id": fileID}))
	if err != nil {
		return nil, err
	}
	f.URL = fmt.Sprintf("%s/%s", storage.bucketName, f.URL)
	return f, nil
}

func (storage *PgMinioStorage) Get(ctx context.Context, fileID uuid.UUID) (*models.File, error) {
	var f models.File
	err := storage.db.QueryRow(ctx, squirrel.Select("id", "name", "url", "created_at").
		From("file").
		Where(squirrel.Eq{"id": fileID})).Scan(&f.ID, &f.Name, &f.URL, &f.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, models.ErrFileNotFound
	} else if err != nil {
		return nil, err
	}
	return &f, nil

}
