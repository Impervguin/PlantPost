package filestorage

import (
	filedir "PlantSite/internal/infra/os/file-dir"
	"PlantSite/internal/infra/sqdb"
	"PlantSite/internal/models"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type PgOsFileStorage struct {
	bucketName string
	fcl        *filedir.FileClient
	db         sqdb.SquirrelDatabase
}

var _ models.FileRepository = (*PgOsFileStorage)(nil)

func NewPgOsFileStorage(bucketName string, fcl *filedir.FileClient, db sqdb.SquirrelDatabase) *PgOsFileStorage {
	fcl.Mkdir(bucketName)
	return &PgOsFileStorage{
		bucketName: bucketName,
		fcl:        fcl,
		db:         db,
	}
}

func (storage *PgOsFileStorage) filePath(name string) string {
	return fmt.Sprintf("%s/%s", storage.bucketName, name)
}

func (storage *PgOsFileStorage) get(ctx context.Context, fileID uuid.UUID) (*models.File, error) {
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

func (storage *PgOsFileStorage) Upload(ctx context.Context, data *models.FileData) (*models.File, error) {
	if data.Reader == nil {
		return nil, models.ErrFileNotFound
	}

	f, err := models.NewFile(data.Name)
	if err != nil {
		return nil, err
	}
	path := storage.filePath(f.URL)

	err = storage.fcl.Put(path, data.Reader)
	if err != nil {
		return nil, err
	}

	_, err = storage.db.Insert(ctx, squirrel.Insert("file").
		Columns("id", "name", "url", "created_at").
		Values(f.ID, f.Name, f.URL, f.CreatedAt))
	if err != nil {
		return nil, err
	}
	f.URL = storage.filePath(f.URL)

	return f, nil
}

func (storage *PgOsFileStorage) Download(ctx context.Context, fileID uuid.UUID) (*models.FileData, error) {
	f, err := storage.get(ctx, fileID)
	if err != nil {
		return nil, err
	}
	attrs, err := storage.fcl.Stat(storage.filePath(f.URL))
	if err != nil {
		return nil, err
	}

	data, err := storage.fcl.Get(storage.filePath(f.URL))
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

func (storage *PgOsFileStorage) Delete(ctx context.Context, fileID uuid.UUID) error {
	f, err := storage.get(ctx, fileID)
	if err != nil {
		return err
	}
	_, err = storage.db.Delete(ctx, squirrel.Delete("file").Where(squirrel.Eq{"id": fileID}))

	errfs := storage.fcl.Delete(storage.filePath(f.URL))
	if errors.Is(errfs, filedir.ErrFileNotFound) {
		errfs = nil
	}
	if errfs != nil || err != nil {
		return fmt.Errorf("error deleting file from database or minio: %w %w", err, errfs)
	}
	return nil
}

func (storage *PgOsFileStorage) Update(ctx context.Context, fileID uuid.UUID, data *models.FileData) (*models.File, error) {
	f, err := storage.get(ctx, fileID)
	if err != nil {
		return nil, err
	}

	f.Name = data.Name
	if data.Reader != nil {
		path := storage.filePath(f.URL)
		err = storage.fcl.Put(path, data.Reader)
		if err != nil {
			return nil, err
		}

		stat, err := storage.fcl.Stat(path)
		if err != nil {
			return nil, err
		}
		f.CreatedAt = stat.ModTime
	}
	_, err = storage.db.Update(ctx, squirrel.Update("file").
		Set("name", f.Name).
		Set("url", f.URL).
		Set("created_at", f.CreatedAt).
		Where(squirrel.Eq{"id": fileID}))
	if err != nil {
		return nil, err
	}
	f.URL = storage.filePath(f.URL)
	return f, nil
}

func (storage *PgOsFileStorage) Get(ctx context.Context, fileID uuid.UUID) (*models.File, error) {
	f, err := storage.get(ctx, fileID)
	if err != nil {
		return nil, err
	}

	f.URL = storage.filePath(f.URL)
	return f, nil
}
