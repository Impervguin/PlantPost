package plantstorage

import (
	"PlantSite/internal/infra/sqdb"
	"PlantSite/internal/models/plant"
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

type PostgresPlantCategoryRepository struct {
	db sqdb.SquirrelDatabase
}

func NewPostgresPlantCategoryRepository(db sqdb.SquirrelDatabase) (*PostgresPlantCategoryRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("nil db")
	}
	return &PostgresPlantCategoryRepository{db: db}, nil
}

func (r *PostgresPlantCategoryRepository) GetCategory(ctx context.Context, name string) (*plant.PlantCategory, error) {
	var category plant.PlantCategory
	row, err := r.db.QueryRow(ctx, squirrel.Select("name", "photo_id").From("plant_category").Where(squirrel.Eq{"name": name}))
	if err != nil {
		return nil, err
	}
	err = row.Scan(&category.Name, &category.MainPhotoID)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *PostgresPlantCategoryRepository) GetCategories(ctx context.Context) ([]plant.PlantCategory, error) {
	var categories []plant.PlantCategory
	rows, err := r.db.Query(ctx, squirrel.Select("name", "photo_id").From("plant_category"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var category plant.PlantCategory
		err = rows.Scan(&category.Name, &category.MainPhotoID)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
