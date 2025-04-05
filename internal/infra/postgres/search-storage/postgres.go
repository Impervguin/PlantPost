package searchstorage

import (
	"PlantSite/internal/infra/postgres/filters"
	plantget "PlantSite/internal/infra/postgres/plant-get"
	postget "PlantSite/internal/infra/postgres/post-get"
	specificationmapper "PlantSite/internal/infra/postgres/specification-mapper"
	"PlantSite/internal/infra/sqpgx"
	"PlantSite/internal/models/plant"
	"PlantSite/internal/models/post"
	"PlantSite/internal/models/search"
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

type PostgresSearchRepository struct {
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

func NewPostgresSearchRepository(ctx context.Context) (*PostgresSearchRepository, error) {
	connStr := configString()
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	return &PostgresSearchRepository{db: sqpgx.NewSquirrelPgx(pool)}, err
}

func (repo *PostgresSearchRepository) GetPostByID(ctx context.Context, id uuid.UUID) (*post.Post, error) {
	get := postget.NewPostgresPostGet(repo.db)
	return get.Get(ctx, id)
}

func (repo *PostgresSearchRepository) GetPlantByID(ctx context.Context, id uuid.UUID) (*plant.Plant, error) {
	get := plantget.NewPostgresPlantGet(repo.db)
	return get.Get(ctx, id)
}

type Post struct {
	ID        uuid.UUID
	Title     string
	Body      string
	AuthorID  uuid.UUID
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostPhoto struct {
	ID          uuid.UUID
	PhotoID     uuid.UUID
	PlaceNumber int
}

func (repo *PostgresSearchRepository) SearchPosts(ctx context.Context, srch *search.PostSearch) ([]*post.Post, error) {
	whereClause, err := filters.NewPostgresPlantSearch()
	if err != nil {
		return nil, err
	}

	srch.Iterate(func(pf search.PostFilter) error {
		filt, err := filters.MapPostFilter(pf)
		if err != nil {
			return err
		}
		return whereClause.AddFilter(filt)
	})

	rows, err := repo.db.Query(ctx,
		squirrel.Select("id", "title", "body", "author_id", "updated_at", "created_at").
			From("post").
			Where(whereClause),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := make([]*post.Post, 0)
	for rows.Next() {
		var pst Post
		err := rows.Scan(&pst.ID, &pst.Title, &pst.Body, &pst.AuthorID, &pst.UpdatedAt, &pst.CreatedAt)
		if err != nil {
			return nil, err
		}

		photos := post.NewPostPhotos()
		subRows, err := repo.db.Query(ctx, squirrel.Select("id", "place_number", "photo_id").
			From("post_photo").
			Where(squirrel.Eq{"post_id": pst.ID}),
		)
		if err != nil {
			return nil, err
		}
		defer subRows.Close()
		for subRows.Next() {
			var tmpPhoto PostPhoto
			err := subRows.Scan(&tmpPhoto.ID, &tmpPhoto.PlaceNumber, &tmpPhoto.PhotoID)
			if err != nil {
				return nil, err
			}
			photo, err := post.CreatePostPhoto(tmpPhoto.ID, tmpPhoto.PhotoID, tmpPhoto.PlaceNumber)
			if err != nil {
				return nil, err
			}
			photos.Add(photo)
		}

		tags := make([]string, 0)
		subRows, err = repo.db.Query(ctx, squirrel.Select("tag").
			From("post_tag").
			Where(squirrel.Eq{"post_id": pst.ID}),
		)
		if err != nil {
			return nil, err
		}
		defer subRows.Close()
		for subRows.Next() {
			var tmpTag string
			err := subRows.Scan(&tmpTag)
			if err != nil {
				return nil, err
			}
			tags = append(tags, tmpTag)
		}
		content, err := post.NewContent(pst.Body, post.ContentTypePlainText)
		if err != nil {
			return nil, err
		}
		pst.Tags = tags
		newPst, err := post.CreatePost(
			pst.ID,
			pst.Title,
			*content,
			pst.Tags,
			pst.AuthorID,
			*photos,
			pst.CreatedAt,
			pst.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, newPst)
	}
	return posts, nil
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

func (repo *PostgresSearchRepository) SearchPlants(ctx context.Context, srch *search.PlantSearch) ([]*plant.Plant, error) {
	whereClause, err := filters.NewPostgresPlantSearch()
	if err != nil {
		return nil, err
	}

	srch.Iterate(func(pf search.PlantFilter) error {
		filt, err := filters.MapPlantFilter(pf)
		if err != nil {
			return err
		}
		return whereClause.AddFilter(filt)
	})

	rows, err := repo.db.Query(ctx,
		squirrel.Select("id", "name", "latin_name", "description", "main_photo_id", "category", "updated_at", "created_at", "specification").
			From("plant").
			Where(whereClause),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	plants := make([]*plant.Plant, 0)
	for rows.Next() {
		var plnt Plant
		var tmpSpec specificationmapper.JsonB
		err := rows.Scan(&plnt.ID, &plnt.Name, &plnt.LatinName, &plnt.Description, &plnt.MainPhotoID, &plnt.Category, &plnt.UpdatedAt, &plnt.CreatedAt, &tmpSpec)
		if err != nil {
			return nil, err
		}

		plnt.Specification, err = specificationmapper.SpecificationFromDB(plnt.Category, tmpSpec)
		if err != nil {
			return nil, err
		}

		photos := plant.NewPlantPhotos()
		subRows, err := repo.db.Query(ctx, squirrel.Select("id", "photo_id", "description").
			From("plant_photo").
			Where(squirrel.Eq{"plant_id": plnt.ID}),
		)
		if err != nil {
			return nil, err
		}
		defer subRows.Close()
		for subRows.Next() {
			var tmpPhoto PlantPhoto
			err := subRows.Scan(&tmpPhoto.ID, &tmpPhoto.PhotoID, &tmpPhoto.Description)
			if err != nil {
				return nil, err
			}
			photo, err := plant.CreatePlantPhoto(tmpPhoto.ID, tmpPhoto.PhotoID, tmpPhoto.Description)
			if err != nil {
				return nil, err
			}
			photos.Add(photo)
		}

		plntSpec, err := plnt.Specification.ToDomain()
		if err != nil {
			return nil, err
		}
		tmpPlant, err := plant.CreatePlant(
			plnt.ID,
			plnt.Name,
			plnt.LatinName,
			plnt.Description,
			plnt.MainPhotoID,
			*photos,
			plnt.Category,
			plntSpec,
			plnt.CreatedAt,
			plnt.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		plants = append(plants, tmpPlant)
	}
	return plants, nil
}
