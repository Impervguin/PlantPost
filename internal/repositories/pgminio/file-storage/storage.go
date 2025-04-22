package filestorage

import (
	minioclient "PlantSite/internal/infra/minio-client"
	"PlantSite/internal/infra/sqdb"
	"PlantSite/internal/models"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type PgMinioStorage struct {
	minioClient *minioclient.MinioClient
	bucketName  string
	db          sqdb.SquirrelDatabase
}

func NewPgMinioStorage(ctx context.Context, db sqdb.SquirrelDatabase, minioCl *minioclient.MinioClient) (*PgMinioStorage, error) {
	return &PgMinioStorage{
		minioClient: minioCl,
		bucketName:  minioCl.GetBucket(),
		db:          db,
	}, nil
}

func (storage *PgMinioStorage) Upload(ctx context.Context, data *models.FileData) (*models.File, error) {
	if data.Reader == nil {
		return nil, errors.New("data.Reader should not be nil")
	}
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
	f, err := storage.get(ctx, fileID)
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

func (storage *PgMinioStorage) Delete(ctx context.Context, fileID uuid.UUID) error {
	f, err := storage.get(ctx, fileID)
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
	f, err := storage.get(ctx, fileID)
	if err != nil {
		return nil, err
	}

	f.Name = data.Name
	if data.Reader != nil {
		info, err := storage.minioClient.PutObject(ctx, storage.bucketName, f.URL, data.Reader, -1, minio.PutObjectOptions{})
		if err != nil {
			return nil, err
		}
		f.URL = info.Key
		f.CreatedAt = time.Now()
	}
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
	f, err := storage.get(ctx, fileID)
	if err != nil {
		return nil, err
	}

	f.URL = fmt.Sprintf("%s/%s", storage.bucketName, f.URL)
	return f, nil
}

func (storage *PgMinioStorage) get(ctx context.Context, fileID uuid.UUID) (*models.File, error) {
	var f models.File
	row, err := storage.db.QueryRow(ctx, squirrel.Select("id", "name", "url", "created_at").
		From("file").
		Where(squirrel.Eq{"id": fileID}))

	if errors.Is(err, sqdb.ErrNoRows) {
		return nil, models.ErrFileNotFound
	} else if err != nil {
		return nil, err
	}

	err = row.Scan(&f.ID, &f.Name, &f.URL, &f.CreatedAt)

	if err == sqdb.ErrNoRows {
		return nil, models.ErrFileNotFound
	} else if err != nil {
		return nil, err
	}
	return &f, nil

}
