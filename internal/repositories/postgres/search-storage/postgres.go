package searchstorage

import (
	"PlantSite/internal/infra/filters"
	specificationmapper "PlantSite/internal/infra/specification-mapper"
	"PlantSite/internal/infra/sqdb"
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/plant"
	"PlantSite/internal/models/post"
	"PlantSite/internal/models/search"
	plantget "PlantSite/internal/repositories/postgres/plant-get"
	postget "PlantSite/internal/repositories/postgres/post-get"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type PostgresSearchRepository struct {
	db sqdb.SquirrelDatabase
}

func NewPostgresSearchRepository(ctx context.Context, db sqdb.SquirrelDatabase) (*PostgresSearchRepository, error) {
	return &PostgresSearchRepository{db: db}, nil
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
	ID          uuid.UUID
	Title       string
	Body        string
	AuthorID    uuid.UUID
	Tags        []string
	ContentType string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PostPhoto struct {
	ID          uuid.UUID
	PhotoID     uuid.UUID
	PlaceNumber int
}

func (repo *PostgresSearchRepository) SearchPosts(ctx context.Context, srch *search.PostSearch) ([]*post.Post, error) {
	whereClause, err := filters.NewPostgresPostSearch()
	if err != nil {
		return nil, fmt.Errorf("PostgresSearchRepository.SearchPosts failed %w", err)
	}

	srch.Iterate(func(pf search.PostFilter) error {
		filt, err := filters.MapPostFilter(pf)
		if err != nil {
			return fmt.Errorf("PostgresSearchRepository.SearchPosts failed %w", err)
		}
		return whereClause.AddFilter(filt)
	})

	rows, err := repo.db.Query(ctx,
		squirrel.Select("id", "title", "body", "author_id", "content_type", "updated_at", "created_at").
			From("post").
			Where(whereClause),
	)
	if errors.Is(err, sqdb.ErrNoRows) {
		return nil, post.ErrPostNotFound
	} else if err != nil {
		return nil, fmt.Errorf("PostgresSearchRepository.SearchPosts failed %w", err)
	}
	psts := make([]Post, 0)
	for rows.Next() {
		var pst Post
		err := rows.Scan(&pst.ID, &pst.Title, &pst.Body, &pst.AuthorID, &pst.ContentType, &pst.UpdatedAt, &pst.CreatedAt)
		if err != nil {
			return nil, err
		}
		psts = append(psts, pst)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	posts := make([]*post.Post, 0)
	for _, pst := range psts {
		photos, err := repo.fetchPostPhotos(ctx, pst.ID)
		if err != nil {
			return nil, fmt.Errorf("PostgresSearchRepository.SearchPosts failed %w", err)
		}
		tags, err := repo.fetchPostTags(ctx, pst.ID)
		if err != nil {
			return nil, fmt.Errorf("PostgresSearchRepository.SearchPosts failed %w", err)
		}
		content, err := post.NewContent(pst.Body, post.ContentFormat(pst.ContentType))
		if err != nil {
			return nil, fmt.Errorf("PostgresSearchRepository.SearchPosts failed %w", err)
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
			return nil, fmt.Errorf("PostgresSearchRepository.SearchPosts failed %w", err)
		}
		posts = append(posts, newPst)
	}
	return posts, nil
}

func (repo *PostgresSearchRepository) fetchPostPhotos(ctx context.Context, postID uuid.UUID) (*post.PostPhotos, error) {
	rows, err := repo.db.Query(ctx, squirrel.Select("id", "place_number", "file_id").
		From("post_photo").
		Where(squirrel.Eq{"post_id": postID}),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	photos := post.NewPostPhotos()
	for rows.Next() {
		var tmpPhoto PostPhoto
		err := rows.Scan(&tmpPhoto.ID, &tmpPhoto.PlaceNumber, &tmpPhoto.PhotoID)
		if err != nil {
			return nil, err
		}
		photo, err := post.CreatePostPhoto(tmpPhoto.ID, tmpPhoto.PhotoID, tmpPhoto.PlaceNumber)
		if err != nil {
			return nil, err
		}
		err = photos.Add(photo)
		if err != nil {
			return nil, err
		}
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return photos, nil
}

func (repo *PostgresSearchRepository) fetchPostTags(ctx context.Context, postID uuid.UUID) ([]string, error) {
	rows, err := repo.db.Query(ctx, squirrel.Select("tag").
		From("post_tag").
		Where(squirrel.Eq{"post_id": postID}),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tags := make([]string, 0)
	for rows.Next() {
		var tmpTag string
		err := rows.Scan(&tmpTag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tmpTag)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return tags, nil
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
		return nil, fmt.Errorf("PostgresSearchRepository.SearchPlants failed %w", err)
	}

	err = srch.Iterate(func(pf search.PlantFilter) error {
		filt, err := filters.MapPlantFilter(pf)
		if err != nil {
			return err
		}
		return whereClause.AddFilter(filt)
	})

	if err != nil {
		return nil, fmt.Errorf("PostgresSearchRepository.SearchPlants failed %w", err)
	}

	rows, err := repo.db.Query(ctx,
		squirrel.Select("id", "name", "latin_name", "description", "main_photo_id", "category", "updated_at", "created_at", "specification").
			From("plant").
			Where(whereClause),
	)
	if err != nil {
		return nil, fmt.Errorf("PostgresSearchRepository.SearchPlants failed %w", err)
	}
	defer rows.Close()
	plants := make([]Plant, 0)
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
		plants = append(plants, plnt)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("PostgresSearchRepository.SearchPlants failed %w", rows.Err())
	}

	truePlants := make([]*plant.Plant, 0)
	for _, plnt := range plants {
		photos, err := repo.fetchPlantPhotos(ctx, plnt.ID)
		if err != nil {
			return nil, fmt.Errorf("PostgresSearchRepository.SearchPlants failed %w", err)
		}
		plantSpec, err := plnt.Specification.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("PostgresSearchRepository.SearchPlants failed %w", err)
		}
		truePlant, err := plant.CreatePlant(
			plnt.ID,
			plnt.Name,
			plnt.LatinName,
			plnt.Description,
			plnt.MainPhotoID,
			*photos,
			plnt.Category,
			plantSpec,
			plnt.CreatedAt,
			plnt.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("PostgresSearchRepository.SearchPlants failed %w", err)
		}
		truePlants = append(truePlants, truePlant)
	}

	return truePlants, nil
}

func (repo *PostgresSearchRepository) fetchPlantPhotos(ctx context.Context, plantID uuid.UUID) (*plant.PlantPhotos, error) {
	rows, err := repo.db.Query(ctx, squirrel.Select("id", "file_id", "description").
		From("plant_photo").
		Where(squirrel.Eq{"plant_id": plantID}),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	photos := plant.NewPlantPhotos()
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
		err = photos.Add(photo)
		if err != nil {
			return nil, err
		}
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return photos, nil
}

func (s *PostgresSearchRepository) GetPostAuthors(ctx context.Context) ([]*auth.Author, error) {
	rows, err := s.db.Query(ctx, squirrel.Select("app_user.id", "app_user.username", "app_user.email", "app_user.password_hash", "app_user.created_at", "author.has_rights", "author.grant_at", "author.revoke_at").
		From("author").Join("app_user ON author.id = app_user.id").
		Where(squirrel.Expr("EXISTS (SELECT 1 FROM post WHERE post.author_id = author.id)")))
	if errors.Is(err, sqdb.ErrNoRows) {
		return []*auth.Author{}, nil
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()
	authors := make([]*auth.Author, 0)
	for rows.Next() {
		var memberID uuid.UUID
		var username, email string
		var passwordHash []byte
		var hasRights bool
		var grantAt, revokeAt, createdAt time.Time
		err := rows.Scan(&memberID, &username, &email, &passwordHash, &createdAt, &hasRights, &grantAt, &revokeAt)
		if err != nil {
			return nil, err
		}
		member, err := auth.CreateMember(memberID, username, email, passwordHash, createdAt)
		if err != nil {
			return nil, err
		}
		author, err := auth.CreateAuthor(*member, grantAt, hasRights, revokeAt)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return authors, nil
}

func (s *PostgresSearchRepository) GetPostTags(ctx context.Context) ([]string, error) {
	rows, err := s.db.Query(ctx, squirrel.Select("DISTINCT tag").
		From("post_tag").
		Where(squirrel.Expr("EXISTS (SELECT 1 FROM post WHERE post.id = post_tag.post_id)")))
	if errors.Is(err, sqdb.ErrNoRows) {
		return []string{}, nil
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()
	tags := make([]string, 0)
	for rows.Next() {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return tags, nil
}
