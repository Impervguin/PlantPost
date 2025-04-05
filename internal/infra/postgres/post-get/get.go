package postget

import (
	"PlantSite/internal/infra/sqpgx"
	"PlantSite/internal/models/post"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostgresPostGet struct {
	db *sqpgx.SquirrelPgx
}

func NewPostgresPostGet(db *sqpgx.SquirrelPgx) *PostgresPostGet {
	return &PostgresPostGet{db: db}
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

func (g *PostgresPostGet) Get(ctx context.Context, postID uuid.UUID) (*post.Post, error) {
	var pst Post
	err := g.db.QueryRow(ctx, squirrel.Select("id", "title", "body", "author_id", "created_at", "updated_at").
		From("post").
		Where(squirrel.Eq{"id": postID}),
	).Scan(&pst.ID, &pst.Title, &pst.Body, &pst.AuthorID, &pst.CreatedAt, &pst.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("post not found")
	} else if err != nil {
		return nil, err
	}
	photos := post.NewPostPhotos()
	rows, err := g.db.Query(ctx, squirrel.Select("id", "place_number", "photo_id").
		From("post_photo").
		Where(squirrel.Eq{"post_id": postID}),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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
	tags := make([]string, 0)
	rows, err = g.db.Query(ctx, squirrel.Select("tag").
		From("post_tag").
		Where(squirrel.Eq{"post_id": postID}),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tmpTag string
		err := rows.Scan(&tmpTag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tmpTag)
	}

	content, err := post.NewContent(pst.Body, post.ContentTypePlainText)
	if err != nil {
		return nil, err
	}

	post, err := post.CreatePost(
		pst.ID,
		pst.Title,
		*content,
		tags,
		pst.AuthorID,
		*photos,
		pst.CreatedAt,
		pst.UpdatedAt,
	)

	return post, err

}
