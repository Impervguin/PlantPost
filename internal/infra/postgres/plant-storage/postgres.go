package plantstorage

import (
	"PlantSite/internal/infra/sqpgx"
	"PlantSite/internal/models/plant"
	"context"
	"fmt"
	"time"

	plantget "PlantSite/internal/infra/postgres/plant-get"
	specificationmapper "PlantSite/internal/infra/postgres/specification-mapper"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

// type PlantRepository interface {
// 	Create(ctx context.Context, plant *Plant) (*Plant, error)
// 	Update(ctx context.Context, plantID uuid.UUID, updateFn func(*Plant) (*Plant, error)) (*Plant, error)
// 	Delete(ctx context.Context, plantID uuid.UUID) error
// 	Get(ctx context.Context, plantID uuid.UUID) (*Plant, error)
// }

// type PlantCategoryRepository interface {
// 	GetCategories(ctx context.Context) ([]PlantCategory, error)
// 	GetCategory(ctx context.Context, name string) (*PlantCategory, error)
// }

type PostgresPlantRepository struct {
	db *sqpgx.SquirrelPgx
}

func configString() string {
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

func NewPostgresPlantRepository(ctx context.Context) (*PostgresPlantRepository, error) {
	connStr := configString()
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	return &PostgresPlantRepository{db: sqpgx.NewSquirrelPgx(pool)}, err
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
	return get.Get(ctx, plantID)
}

func (repo *PostgresPlantRepository) Create(ctx context.Context, plnt *plant.Plant) (*plant.Plant, error) {
	tx, err := repo.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	tmpSpec, err := specificationmapper.SpecificationFromDomain(plnt.GetCategory(), plnt.GetSpecification())
	if err != nil {
		return nil, err
	}
	specJson, err := tmpSpec.ToJsonB()
	if err != nil {
		return nil, err
	}
	_, err = tx.Insert(ctx, squirrel.Insert("plant").
		Columns("id", "name", "latin_name", "description", "main_photo_id", "category", "created_at", "updated_at", "specification").
		Values(plnt.ID(), plnt.GetName(), plnt.GetLatinName(), plnt.GetDescription(), plnt.MainPhotoID(), plnt.GetCategory(), plnt.CreatedAt(), plnt.UpdatedAt(), specJson),
	)

	if err != nil {
		return nil, err
	}
	query := squirrel.Insert("plant_photo").
		Columns("id", "plant_id", "description")

	err = plnt.GetPhotos().Iterate(func(e plant.PlantPhoto) error {
		query = query.Values(e.ID(), plnt.ID(), e.Description())
		return nil
	})
	if err != nil {
		return nil, err
	}
	_, err = tx.Insert(ctx, query)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)

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

	tx, err := repo.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Update(ctx, squirrel.Update("plant").
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
		return nil, err
	}

	_, err = tx.Delete(ctx, squirrel.Delete("plant_photo").
		Where(squirrel.Eq{"plant_id": plantID}))
	if err != nil {
		return nil, err
	}
	query := squirrel.Insert("plant_photo").
		Columns("id", "plant_id", "description")
	err = plnt.GetPhotos().Iterate(func(e plant.PlantPhoto) error {
		query = query.Values(e.ID(), plnt.ID(), e.Description())
		return nil
	})
	if err != nil {
		return nil, err
	}
	_, err = tx.Insert(ctx, query)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)
	return plnt, err
}

func (repo *PostgresPlantRepository) Delete(ctx context.Context, plantID uuid.UUID) error {
	tx, err := repo.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Delete(ctx, squirrel.Delete("plant_photo").Where(squirrel.Eq{"plant_id": plantID}))
	if err != nil {
		return err
	}

	_, err = tx.Delete(ctx, squirrel.Delete("plant").Where(squirrel.Eq{"id": plantID}))
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
