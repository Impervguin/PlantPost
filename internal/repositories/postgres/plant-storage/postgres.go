package plantstorage

import (
	"PlantSite/internal/infra/sqdb"
	"PlantSite/internal/models/plant"
	"context"
	"errors"
	"fmt"
	"time"

	specificationmapper "PlantSite/internal/infra/specification-mapper"
	plantget "PlantSite/internal/repositories/postgres/plant-get"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type PostgresPlantRepository struct {
	db sqdb.SquirrelDatabase
}

func NewPostgresPlantRepository(ctx context.Context, db sqdb.SquirrelDatabase) (*PostgresPlantRepository, error) {
	return &PostgresPlantRepository{db: db}, nil
}

type Plant struct {
	ID            uuid.UUID
	Name          string
	LatinName     string
	Description   string
	MainPhotoID   uuid.UUID
	Category      string
	Specification specificationmapper.PlantSpecification
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type PlantPhoto struct {
	ID          uuid.UUID
	PhotoID     uuid.UUID
	Description string
}

func (repo *PostgresPlantRepository) Get(ctx context.Context, plantID uuid.UUID) (*plant.Plant, error) {
	get := plantget.NewPostgresPlantGet(repo.db)
	plnt, err := get.Get(ctx, plantID)
	if err != nil {
		return nil, fmt.Errorf("PostgresPlantRepository.Get failed %w", err)
	}
	return plnt, nil
}

func (repo *PostgresPlantRepository) Create(ctx context.Context, plnt *plant.Plant) (*plant.Plant, error) {
	err := repo.db.Transaction(ctx, func(tx sqdb.SquirrelQuirier) error {
		tmpSpec, err := specificationmapper.SpecificationFromDomain(plnt.GetCategory(), plnt.GetSpecification())
		if err != nil {
			return err
		}
		specJson, err := tmpSpec.ToJsonB()
		if err != nil {
			return err
		}
		_, err = tx.Insert(ctx, squirrel.Insert("plant").
			Columns("id", "name", "latin_name", "description", "main_photo_id", "category", "created_at", "updated_at", "specification").
			Values(plnt.ID(), plnt.GetName(), plnt.GetLatinName(), plnt.GetDescription(), plnt.MainPhotoID(), plnt.GetCategory(), plnt.CreatedAt(), plnt.UpdatedAt(), specJson),
		)

		if err != nil {
			return err
		}
		if plnt.GetPhotos().Len() > 0 {
			query := squirrel.Insert("plant_photo").
				Columns("id", "plant_id", "file_id", "description")

			err = plnt.GetPhotos().Iterate(func(e plant.PlantPhoto) error {
				query = query.Values(e.ID(), plnt.ID(), e.FileID(), e.Description())
				return nil
			})
			if err != nil {
				return err
			}
			_, err = tx.Insert(ctx, query)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("PostgresPlantRepository.Create failed %w", err)
	}
	return plnt, err
}

func (repo *PostgresPlantRepository) Update(ctx context.Context, plantID uuid.UUID, updateFn func(*plant.Plant) (*plant.Plant, error)) (*plant.Plant, error) {
	plnt, err := repo.Get(ctx, plantID)
	if err != nil {
		return nil, err
	}
	plnt, err = updateFn(plnt)
	if err != nil {
		return nil, err
	}
	tmpSpec, err := specificationmapper.SpecificationFromDomain(plnt.GetCategory(), plnt.GetSpecification())
	if err != nil {
		return nil, err
	}
	specJson, err := tmpSpec.ToJsonB()
	if err != nil {
		return nil, err
	}

	err = repo.db.Transaction(ctx, func(tx sqdb.SquirrelQuirier) error {
		_, err := tx.Update(ctx, squirrel.Update("plant").
			Set("name", plnt.GetName()).
			Set("latin_name", plnt.GetLatinName()).
			Set("description", plnt.GetDescription()).
			Set("main_photo_id", plnt.MainPhotoID()).
			Set("category", plnt.GetCategory()).
			Set("created_at", plnt.CreatedAt()).
			Set("updated_at", plnt.UpdatedAt()).
			Set("specification", specJson).
			Where(squirrel.Eq{"id": plantID}),
		)
		if err != nil {
			return err
		}

		_, err = tx.Delete(ctx, squirrel.Delete("plant_photo").
			Where(squirrel.Eq{"plant_id": plantID}))
		if err != nil && !errors.Is(err, sqdb.ErrNoRows) {
			return err
		}
		if plnt.GetPhotos().Len() > 0 {
			query := squirrel.Insert("plant_photo").
				Columns("id", "plant_id", "file_id", "description")
			err = plnt.GetPhotos().Iterate(func(e plant.PlantPhoto) error {
				query = query.Values(e.ID(), plnt.ID(), e.FileID(), e.Description())
				return nil
			})
			if err != nil {
				return err
			}
			_, err = tx.Insert(ctx, query)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("PostgresPlantRepository.Update failed %w", err)
	}
	return plnt, nil
}

func (repo *PostgresPlantRepository) Delete(ctx context.Context, plantID uuid.UUID) error {
	err := repo.db.Transaction(ctx, func(tx sqdb.SquirrelQuirier) error {
		_, err := tx.Delete(ctx, squirrel.Delete("plant_photo").Where(squirrel.Eq{"plant_id": plantID}))
		if err != nil && !errors.Is(err, sqdb.ErrNoRows) {
			return err
		}
		_, err = tx.Delete(ctx, squirrel.Delete("plant").Where(squirrel.Eq{"id": plantID}))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("PostgresPlantRepository.Delete failed %w", err)
	}
	return nil
}
