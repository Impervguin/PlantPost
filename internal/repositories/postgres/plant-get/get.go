package plantget

import (
	"PlantSite/internal/infra/sqdb"
	"PlantSite/internal/models/plant"
	"context"
	"fmt"
	"time"

	specificationmapper "PlantSite/internal/infra/specification-mapper"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

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

type PostgresPlantGet struct {
	db sqdb.SquirrelDatabase
}

func NewPostgresPlantGet(db sqdb.SquirrelDatabase) *PostgresPlantGet {
	return &PostgresPlantGet{db: db}
}

func (g *PostgresPlantGet) Get(ctx context.Context, plantID uuid.UUID) (*plant.Plant, error) {
	var tmpPlant Plant
	var tmpSpec specificationmapper.JsonB
	photos := plant.NewPlantPhotos()

	row, err := g.db.QueryRow(ctx,
		squirrel.Select("id", "name", "latin_name", "description", "main_photo_id", "category", "created_at", "updated_at", "specification").
			From("plant").
			Where(squirrel.Eq{"id": plantID}),
	)
	if err != nil {
		return nil, fmt.Errorf("PostgresPlantGet failed %w", err)
	}

	err = row.Scan(&tmpPlant.ID, &tmpPlant.Name, &tmpPlant.LatinName, &tmpPlant.Description, &tmpPlant.MainPhotoID, &tmpPlant.Category, &tmpPlant.CreatedAt, &tmpPlant.UpdatedAt, &tmpSpec)
	if err == sqdb.ErrNoRows {
		return nil, plant.ErrPlantNotFound
	} else if err != nil {
		return nil, err
	}

	tmpPlant.Specification, err = specificationmapper.SpecificationFromDB(tmpPlant.Category, tmpSpec)
	if err != nil {
		return nil, err
	}

	rows, err := g.db.Query(ctx,
		squirrel.Select("id", "file_id", "description").
			From("plant_photo").
			Where(squirrel.Eq{"plant_id": plantID}),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tmpPhoto PlantPhoto
		err := rows.Scan(&tmpPhoto.ID, &tmpPhoto.PhotoID, &tmpPhoto.Description)
		if err != nil {
			return nil, err
		}
		photo, err := plant.CreatePlantPhoto(tmpPhoto.ID, tmpPhoto.PhotoID, tmpPhoto.Description)
		if err != nil {
			return nil, err
		}
		photos.Add(photo)
	}

	plntSpec, err := tmpPlant.Specification.ToDomain()
	if err != nil {
		return nil, err
	}

	plnt, err := plant.CreatePlant(
		tmpPlant.ID,
		tmpPlant.Name,
		tmpPlant.LatinName,
		tmpPlant.Description,
		tmpPlant.MainPhotoID,
		*photos,
		tmpPlant.Category,
		plntSpec,
		tmpPlant.CreatedAt,
		tmpPlant.UpdatedAt,
	)
	return plnt, err
}
