package poststorage

import (
	"PlantSite/internal/infra/sqdb"
	"PlantSite/internal/models/post"
	"PlantSite/internal/models/post/parser"
	postget "PlantSite/internal/repositories/postgres/post-get"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type PostgresPostRepository struct {
	db          sqdb.SquirrelDatabase
	plantGetter parser.PlantGetter
}

func NewPostgresPostRepository(ctx context.Context, db sqdb.SquirrelDatabase, plantGetter parser.PlantGetter) (*PostgresPostRepository, error) {
	return &PostgresPostRepository{db: db, plantGetter: plantGetter}, nil
}

func (repo *PostgresPostRepository) Create(ctx context.Context, pst *post.Post) (*post.Post, error) {
	err := repo.db.Transaction(ctx, func(tx sqdb.SquirrelQuirier) error {
		var content post.Content
		var contentWithPlant *post.ContentWithPlant
		content = pst.Content()
		if post.CheckContentWithPlant(&content) {
			plantParser, err := parser.GetParser(&content, repo.plantGetter)
			if err != nil {
				return fmt.Errorf("PostgresPostRepository.Create can't get plant parser: %w", err)
			}
			contentWithPlant, err = post.NewContentWithPlant(pst.Content().Text, post.ContentFormat(pst.Content().ContentType), plantParser)
			if err != nil {
				return fmt.Errorf("PostgresPostRepository.Create can't create plant content: %w", err)
			}
			content = contentWithPlant.Content
		}

		_, err := tx.Insert(ctx, squirrel.Insert("post").
			Columns("id", "title", "body", "content_type", "author_id", "created_at", "updated_at").
			Values(pst.ID(), pst.Title(), content.Text, content.ContentType, pst.AuthorID(), pst.CreatedAt(), pst.UpdatedAt()),
		)
		if err != nil {
			return err
		}
		if pst.Photos().Len() > 0 {
			query := squirrel.Insert("post_photo").
				Columns("id", "post_id", "file_id", "place_number")

			for _, photo := range pst.Photos().List() {
				query = query.Values(photo.ID(), pst.ID(), photo.FileID(), photo.PlaceNumber())
			}
			_, err = tx.Insert(ctx, query)
			if err != nil {
				return err
			}
		}

		if len(pst.Tags()) > 0 {
			query := squirrel.Insert("post_tag").
				Columns("id", "post_id", "tag")

			for _, tag := range pst.Tags() {
				query = query.Values(uuid.New(), pst.ID(), tag)
			}
			_, err = tx.Insert(ctx, query)
			if err != nil {
				return err
			}
		}

		if contentWithPlant != nil {

			if len(contentWithPlant.PlantIDs()) > 0 {
				query := squirrel.Insert("plant_post").Columns("id", "plant_id", "post_id")
				for _, plantID := range contentWithPlant.PlantIDs() {
					query = query.Values(uuid.New(), plantID, pst.ID())
				}
				_, err = tx.Insert(ctx, query)
				if err != nil {
					return fmt.Errorf("PostgresPostRepository.Create can't insert plant posts: %w", err)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("PostgresPostRepository.Create failed %w", err)
	}

	return pst, nil
}

type Post struct {
	ID          uuid.UUID
	Title       string
	Body        string
	ContentType string
	AuthorID    uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PostPhoto struct {
	ID          uuid.UUID
	PlaceNumber int
	PhotoID     uuid.UUID
}

func (repo *PostgresPostRepository) Get(ctx context.Context, postID uuid.UUID) (*post.Post, error) {
	get := postget.NewPostgresPostGet(repo.db)
	pst, err := get.Get(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("PostgresPostRepository.Get failed %w", err)
	}
	return pst, nil
}

func (repo *PostgresPostRepository) Update(ctx context.Context, id uuid.UUID, updateFn func(*post.Post) (*post.Post, error)) (*post.Post, error) {
	pst, err := repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	updatedPst, err := updateFn(pst)
	if err != nil {
		return nil, err
	}

	err = repo.db.Transaction(ctx, func(tx sqdb.SquirrelQuirier) error {
		var content post.Content
		var contentWithPlant *post.ContentWithPlant
		content = pst.Content()
		if post.CheckContentWithPlant(&content) {
			plantParser, err := parser.GetParser(&content, repo.plantGetter)
			if err != nil {
				return fmt.Errorf("PostgresPostRepository.Create can't get plant parser: %w", err)
			}
			contentWithPlant, err = post.NewContentWithPlant(pst.Content().Text, post.ContentFormat(pst.Content().ContentType), plantParser)
			if err != nil {
				return fmt.Errorf("PostgresPostRepository.Create can't create plant content: %w", err)
			}
			content = contentWithPlant.Content
		}

		_, err := repo.db.Update(ctx, squirrel.Update("post").
			Set("title", updatedPst.Title()).
			Set("body", content.Text).
			Set("content_type", content.ContentType).
			Set("author_id", updatedPst.AuthorID()).
			Set("updated_at", updatedPst.UpdatedAt()).
			Where(squirrel.Eq{"id": id}))
		if err != nil {
			return err
		}
		_, err = tx.Delete(ctx, squirrel.Delete("post_photo").
			Where(squirrel.Eq{"post_id": id}))
		if err != nil && !errors.Is(err, sqdb.ErrNoRows) {
			return err
		}
		for _, photo := range updatedPst.Photos().List() {
			_, err = tx.Insert(ctx, squirrel.Insert("post_photo").
				Columns("id", "post_id", "file_id", "place_number").
				Values(photo.ID(), id, photo.FileID(), photo.PlaceNumber()))
			if err != nil {
				return err
			}
		}
		_, err = tx.Delete(ctx, squirrel.Delete("post_tag").
			Where(squirrel.Eq{"post_id": id}))
		if err != nil && !errors.Is(err, sqdb.ErrNoRows) {
			return err
		}
		for _, tag := range updatedPst.Tags() {
			_, err = tx.Insert(ctx, squirrel.Insert("post_tag").
				Columns("id", "post_id", "tag").
				Values(uuid.New(), id, tag))
			if err != nil {
				return err
			}
		}

		_, err = tx.Delete(ctx, squirrel.Delete("plant_post").
			Where(squirrel.Eq{"post_id": id}))
		if err != nil && !errors.Is(err, sqdb.ErrNoRows) {
			return fmt.Errorf("PostgresPostRepository.Update can't delete post plants: %w", err)
		}

		if contentWithPlant != nil {

			if len(contentWithPlant.PlantIDs()) > 0 {
				query := squirrel.Insert("plant_post").Columns("id", "plant_id", "post_id")
				for _, plantID := range contentWithPlant.PlantIDs() {
					query = query.Values(uuid.New(), plantID, pst.ID())
				}
				_, err = tx.Insert(ctx, query)
				if err != nil {
					return fmt.Errorf("PostgresPostRepository.Create can't insert plant posts: %w", err)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("PostgresPostRepository.Update failed %w", err)
	}
	return updatedPst, nil
}

func (repo *PostgresPostRepository) Delete(ctx context.Context, postID uuid.UUID) error {
	err := repo.db.Transaction(ctx, func(tx sqdb.SquirrelQuirier) error {
		_, err := tx.Delete(ctx, squirrel.Delete("post_photo").
			Where(squirrel.Eq{"post_id": postID}))
		if err != nil && !errors.Is(err, sqdb.ErrNoRows) {
			return err
		}
		_, err = tx.Delete(ctx, squirrel.Delete("post_tag").
			Where(squirrel.Eq{"post_id": postID}))

		if err != nil && !errors.Is(err, sqdb.ErrNoRows) {
			return err
		}

		_, err = tx.Delete(ctx, squirrel.Delete("plant_post").
			Where(squirrel.Eq{"post_id": postID}))
		if err != nil && !errors.Is(err, sqdb.ErrNoRows) {
			return fmt.Errorf("PostgresPostRepository.Delete can't delete post plants: %w", err)
		}

		_, err = tx.Delete(ctx, squirrel.Delete("post").
			Where(squirrel.Eq{"id": postID}))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("PostgresPostRepository.Delete failed %w", err)
	}
	return nil
}

type PostRow struct {
	ID          uuid.UUID
	Title       string
	Body        string
	ContentType string
	AuthorID    uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (repo *PostgresPostRepository) ListAuthorPosts(ctx context.Context, authorID uuid.UUID) ([]*post.Post, error) {
	postRows, err := repo.fetchPostsByAuthor(ctx, authorID)
	if err != nil {
		return nil, fmt.Errorf("PostgresPostRepository.ListAuthorPosts failed %w", err)
	}
	posts := make([]*post.Post, 0)
	for _, postRow := range postRows {
		photos, err := repo.fetchPostPhotos(ctx, postRow.ID)
		if err != nil {
			return nil, fmt.Errorf("PostgresPostRepository.ListAuthorPosts failed %w", err)
		}
		tags, err := repo.fetchPostTags(ctx, postRow.ID)
		if err != nil {
			return nil, fmt.Errorf("PostgresPostRepository.ListAuthorPosts failed %w", err)
		}
		content, err := post.NewContent(postRow.Body, post.ContentFormat(postRow.ContentType))
		if err != nil {
			return nil, fmt.Errorf("PostgresPostRepository.ListAuthorPosts failed %w", err)
		}
		p, err := post.CreatePost(postRow.ID, postRow.Title, *content, tags, postRow.AuthorID, *photos, postRow.CreatedAt, postRow.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("PostgresPostRepository.ListAuthorPosts failed %w", err)
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (repo *PostgresPostRepository) fetchPostsByAuthor(ctx context.Context, authorID uuid.UUID) ([]*PostRow, error) {
	rows, err := repo.db.Query(ctx, squirrel.Select("id", "title", "body", "content_type", "author_id", "created_at", "updated_at").
		From("post").
		Where(squirrel.Eq{"author_id": authorID}),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := make([]*PostRow, 0)
	for rows.Next() {
		var pst PostRow
		err := rows.Scan(&pst.ID, &pst.Title, &pst.Body, &pst.ContentType, &pst.AuthorID, &pst.CreatedAt, &pst.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &pst)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return posts, nil
}

func (repo *PostgresPostRepository) fetchPostPhotos(ctx context.Context, postID uuid.UUID) (*post.PostPhotos, error) {
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

func (repo *PostgresPostRepository) fetchPostTags(ctx context.Context, postID uuid.UUID) ([]string, error) {
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
