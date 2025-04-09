package albumstorage

import (
	"PlantSite/internal/infra/sqdb"
	"PlantSite/internal/models/album"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type PostgresAlbumRepository struct {
	db sqdb.SquirrelDatabase
}

func NewPostgresAlbumRepository(ctx context.Context, db sqdb.SquirrelDatabase) (*PostgresAlbumRepository, error) {
	return &PostgresAlbumRepository{db: db}, nil
}

type Album struct {
	ID          uuid.UUID
	Name        string
	Description string
	OwnerID     uuid.UUID
	PlantIDs    uuid.UUIDs
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (repo *PostgresAlbumRepository) Create(ctx context.Context, alb *album.Album) (*album.Album, error) {
	err := repo.db.Transaction(ctx, func(tx sqdb.SquirrelQuirier) error {
		_, err := tx.Insert(ctx, squirrel.Insert("album").
			Columns("id", "name", "description", "owner_id", "created_at", "updated_at").
			Values(alb.ID(), alb.Name(), alb.Description(), alb.GetOwnerID(), alb.CreatedAt(), alb.UpdatedAt()),
		)

		if err != nil {
			return fmt.Errorf("PostgresAlbumRepository.Create failed %w", err)
		}

		plantIds := alb.PlantIDs()
		if len(plantIds) > 0 {
			query := squirrel.Insert("plant_album").
				Columns("id", "album_id", "plant_id")

			for _, id := range plantIds {
				query = query.Values(uuid.New(), alb.ID(), id)
			}

			_, err = tx.Insert(ctx, query)
			if err != nil {
				return fmt.Errorf("PostgresAlbumRepository.Create failed %w", err)
			}
		}

		return nil
	})

	return alb, err
}

func (repo *PostgresAlbumRepository) Get(ctx context.Context, id uuid.UUID) (*album.Album, error) {
	var tmpAlbum Album
	row, err := repo.db.QueryRow(ctx, squirrel.Select("id", "name", "description", "owner_id", "created_at", "updated_at").
		From("album").
		Where(squirrel.Eq{"id": id}),
	)
	if errors.Is(err, sqdb.ErrNoRows) {
		return nil, fmt.Errorf("PostgresAlbumRepository.Get failed %w", album.ErrAlbumNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("PostgresAlbumRepository.Get failed %w", err)
	}
	err = row.Scan(&tmpAlbum.ID, &tmpAlbum.Name, &tmpAlbum.Description, &tmpAlbum.OwnerID, &tmpAlbum.CreatedAt, &tmpAlbum.UpdatedAt)

	if errors.Is(err, sqdb.ErrNoRows) {
		return nil, fmt.Errorf("PostgresAlbumRepository.Get failed %w", album.ErrAlbumNotFound)
	} else if err != nil {
		return nil, err
	}

	plantIDs, err := repo.fetchPlantIDs(ctx, tmpAlbum.ID)
	if err != nil {
		return nil, fmt.Errorf("PostgresAlbumRepository.Get failed %w", err)
	}

	alb, err := album.CreateAlbum(
		tmpAlbum.ID,
		tmpAlbum.Name,
		tmpAlbum.Description,
		plantIDs,
		tmpAlbum.OwnerID,
		tmpAlbum.CreatedAt,
		tmpAlbum.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("PostgresAlbumRepository.Get failed %w", err)
	}
	return alb, nil
}

func (repo *PostgresAlbumRepository) Update(ctx context.Context, id uuid.UUID, updateFn func(*album.Album) (*album.Album, error)) (*album.Album, error) {
	alb, err := repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("PostgresAlbumRepository.Update failed %w", err)
	}
	alb, err = updateFn(alb)
	if err != nil {
		return nil, fmt.Errorf("PostgresAlbumRepository.Update failed %w", err)
	}

	err = repo.db.Transaction(ctx, func(tx sqdb.SquirrelQuirier) error {
		_, err := tx.Update(ctx, squirrel.Update("album").
			Set("name", alb.Name()).
			Set("description", alb.Description()).
			Set("owner_id", alb.GetOwnerID()).
			Set("updated_at", alb.UpdatedAt()).
			Where(squirrel.Eq{"id": alb.ID()}),
		)
		if err != nil {
			return err
		}
		_, err = tx.Delete(ctx, squirrel.Delete("plant_album").
			Where(squirrel.Eq{"album_id": id}),
		)

		if err != nil {
			return fmt.Errorf("PostgresAlbumRepository.Update failed %w", err)
		}
		plantIds := alb.PlantIDs()
		if len(plantIds) > 0 {
			query := squirrel.Insert("plant_album").
				Columns("id", "album_id", "plant_id")

			for _, id := range plantIds {
				query = query.Values(uuid.New(), alb.ID(), id)
			}

			_, err = tx.Insert(ctx, query)
			if err != nil {
				return fmt.Errorf("PostgresAlbumRepository.Update plants failed %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("PostgresAlbumRepository.Update failed %w", err)
	}

	return alb, nil
}

func (repo *PostgresAlbumRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := repo.db.Transaction(ctx, func(tx sqdb.SquirrelQuirier) error {

		_, err := tx.Delete(ctx, squirrel.Delete("plant_album").
			Where(squirrel.Eq{"album_id": id}),
		)

		if err != nil {
			return fmt.Errorf("PostgresAlbumRepository.Delete failed %w", err)
		}

		_, err = tx.Delete(ctx, squirrel.Delete("album").
			Where(squirrel.Eq{"id": id}),
		)
		if err != nil {
			return fmt.Errorf("PostgresAlbumRepository.Delete failed %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("PostgresAlbumRepository.Delete failed %w", err)
	}

	return nil
}

type AlbumRow struct {
	ID          uuid.UUID
	Name        string
	Description string
	OwnerID     uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (repo *PostgresAlbumRepository) List(ctx context.Context, ownerID uuid.UUID) ([]*album.Album, error) {
	albums := make([]*album.Album, 0)
	albs, err := repo.fetchAlbumsByOwner(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("PostgresAlbumRepository.List failed %w", err)
	}
	for _, alb := range albs {
		var plantIDs uuid.UUIDs
		plantIDs, err = repo.fetchPlantIDs(ctx, alb.ID)
		if err != nil {
			return nil, fmt.Errorf("PostgresAlbumRepository.List failed %w", err)
		}
		alb, err := album.CreateAlbum(
			alb.ID,
			alb.Name,
			alb.Description,
			plantIDs,
			alb.OwnerID,
			alb.CreatedAt,
			alb.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("PostgresAlbumRepository.List failed %w", err)
		}
		albums = append(albums, alb)
	}
	return albums, nil
}

func (repo *PostgresAlbumRepository) fetchAlbumsByOwner(ctx context.Context, ownerID uuid.UUID) ([]*AlbumRow, error) {
	var albums []*AlbumRow
	rows, err := repo.db.Query(ctx, squirrel.Select("id", "name", "description", "owner_id", "created_at", "updated_at").
		From("album").
		Where(squirrel.Eq{"owner_id": ownerID}),
	)
	if errors.Is(err, sqdb.ErrNoRows) {
		return albums, nil
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tmpAlbum AlbumRow
		err := rows.Scan(&tmpAlbum.ID, &tmpAlbum.Name, &tmpAlbum.Description, &tmpAlbum.OwnerID, &tmpAlbum.CreatedAt, &tmpAlbum.UpdatedAt)
		if err != nil {
			return nil, err
		}
		albums = append(albums, &tmpAlbum)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return albums, nil
}

func (repo *PostgresAlbumRepository) fetchPlantIDs(ctx context.Context, albumID uuid.UUID) (uuid.UUIDs, error) {
	plantIDs := make(uuid.UUIDs, 0)
	rows, err := repo.db.Query(ctx, squirrel.Select("plant_id").
		From("plant_album").
		Where(squirrel.Eq{"album_id": albumID}),
	)
	if errors.Is(err, sqdb.ErrNoRows) {
		return plantIDs, nil
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var plantID uuid.UUID
		err := rows.Scan(&plantID)
		if err != nil {
			return nil, err
		}
		plantIDs = append(plantIDs, plantID)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return plantIDs, nil
}
