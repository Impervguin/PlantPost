package poststorage

import (
	postget "PlantSite/internal/infra/postgres/post-get"
	"PlantSite/internal/infra/sqpgx"
	"PlantSite/internal/models/post"
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

type PostgresPostRepository struct {
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

func NewPostgresPostRepository(ctx context.Context) (*PostgresPostRepository, error) {
	connStr := configString()
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	return &PostgresPostRepository{db: sqpgx.NewSquirrelPgx(pool)}, err
}

func (repo *PostgresPostRepository) Create(ctx context.Context, pst *post.Post) (*post.Post, error) {
	tx, err := repo.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Insert(ctx, squirrel.Insert("post").
		Columns("id", "title", "body", "author_id", "created_at", "updated_at").
		Values(pst.ID(), pst.Title(), pst.Content().Text, pst.AuthorID(), pst.CreatedAt(), pst.UpdatedAt()),
	)
	if err != nil {
		return nil, err
	}

	query := squirrel.Insert("post_photo").
		Columns("id", "post_id", "photo_id", "place_number")

	for _, photo := range pst.Photos().List() {
		query = query.Values(photo.ID(), pst.ID(), photo.FileID(), photo.PlaceNumber())
	}
	_, err = tx.Insert(ctx, query)
	if err != nil {
		return nil, err
	}

	query = squirrel.Insert("post_tag").
		Columns("id", "post_id", "tag")

	for _, tag := range pst.Tags() {
		query = query.Values(uuid.New(), pst.ID(), tag)
	}
	_, err = tx.Insert(ctx, query)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)
	return pst, err
}

type Post struct {
	ID        uuid.UUID
	Title     string
	Body      string
	AuthorID  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostPhoto struct {
	ID          uuid.UUID
	PlaceNumber int
	PhotoID     uuid.UUID
}

func (repo *PostgresPostRepository) Get(ctx context.Context, postID uuid.UUID) (*post.Post, error) {
	get := postget.NewPostgresPostGet(repo.db)
	return get.Get(ctx, postID)
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

	tx, err := repo.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	_, err = repo.db.Update(ctx, squirrel.Update("post").
		Set("title", updatedPst.Title()).
		Set("body", updatedPst.Content().Text).
		Set("author_id", updatedPst.AuthorID()).
		Set("updated_at", updatedPst.UpdatedAt()).
		Where(squirrel.Eq{"id": id}))
	if err != nil {
		return nil, err
	}
	_, err = tx.Delete(ctx, squirrel.Delete("post_photo").
		Where(squirrel.Eq{"post_id": id}))
	if err != nil {
		return nil, err
	}
	for _, photo := range updatedPst.Photos().List() {
		_, err = tx.Insert(ctx, squirrel.Insert("post_photo").
			Columns("id", "post_id", "photo_id", "place_number").
			Values(photo.ID(), id, photo.FileID(), photo.PlaceNumber()))
		if err != nil {
			return nil, err
		}
	}
	_, err = tx.Delete(ctx, squirrel.Delete("post_tag").
		Where(squirrel.Eq{"post_id": id}))
	if err != nil {
		return nil, err
	}
	for _, tag := range updatedPst.Tags() {
		_, err = tx.Insert(ctx, squirrel.Insert("post_tag").
			Columns("id", "post_id", "tag").
			Values(uuid.New(), id, tag))
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit(ctx)
	return updatedPst, err
}

func (repo *PostgresPostRepository) Delete(ctx context.Context, postID uuid.UUID) error {
	tx, err := repo.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Delete(ctx, squirrel.Delete("post_photo").
		Where(squirrel.Eq{"post_id": postID}))
	if err != nil {
		return err
	}
	_, err = tx.Delete(ctx, squirrel.Delete("post_tag").
		Where(squirrel.Eq{"post_id": postID}))
	if err != nil {
		return err
	}
	_, err = tx.Delete(ctx, squirrel.Delete("post").
		Where(squirrel.Eq{"id": postID}))
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (repo *PostgresPostRepository) ListAuthorPosts(ctx context.Context, authorID uuid.UUID) ([]*post.Post, error) {
	rows, err := repo.db.Query(ctx, squirrel.Select("id", "title", "body", "author_id", "created_at", "updated_at").
		From("post").
		Where(squirrel.Eq{"author_id": authorID}))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := make([]*post.Post, 0)
	for rows.Next() {
		var pst Post
		err := rows.Scan(&pst.ID, &pst.Title, &pst.Body, &pst.AuthorID, &pst.CreatedAt, &pst.UpdatedAt)
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
			err = photos.Add(photo)
			if err != nil {
				return nil, err
			}
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
		post, err := post.CreatePost(pst.ID, pst.Title, *content, tags, pst.AuthorID, *photos, pst.CreatedAt, pst.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
