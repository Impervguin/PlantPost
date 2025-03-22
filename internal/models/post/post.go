package post

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID
	Title     string
	Content   Content
	Tags      []string
	AuthorID  uuid.UUID
	CreatedAt time.Time
	Photos    []PostPhoto
}

type PostPhoto struct {
	ID          uuid.UUID
	PlaceNumber int
	FileID      uuid.UUID
}

func (p *Post) Validate() error {
	if p.ID == uuid.Nil {
		return fmt.Errorf("post ID cannot be empty")
	}
	if p.Title == "" {
		return fmt.Errorf("post title cannot be empty")
	}
	if err := p.Content.Validate(); err != nil {
		return fmt.Errorf("post content validation failed: %v", err)
	}
	if p.AuthorID == uuid.Nil {
		return fmt.Errorf("author ID cannot be empty")
	}
	if len(p.Photos) > MaximumPhotoPerPostCount {
		return fmt.Errorf("maximum number of photos exceeded")
	}
	if p.Tags == nil {
		return fmt.Errorf("tags cannot be nil")
	}
	if p.Photos == nil {
		return fmt.Errorf("photos cannot be nil")
	}
	photoSet := make(map[int]struct{})
	for _, photo := range p.Photos {
		if photo.PlaceNumber < 1 || photo.PlaceNumber > 10 {
			return fmt.Errorf("invalid photo place number: %d", photo.PlaceNumber)
		}
		if _, exists := photoSet[photo.PlaceNumber]; exists {
			return fmt.Errorf("duplicate photo place number: %d", photo.PlaceNumber)
		}
		photoSet[photo.PlaceNumber] = struct{}{}
		if photo.FileID == uuid.Nil {
			return fmt.Errorf("photo file ID cannot be empty")
		}
	}

	for _, tag := range p.Tags {
		if tag == "" {
			return fmt.Errorf("empty tag")
		}
	}

	if p.CreatedAt.After(time.Now()) {
		return fmt.Errorf("post creation date cannot be in the future")
	}
	return nil
}

func CreatePost(
	id uuid.UUID,
	title string,
	content Content,
	tags []string,
	authorID uuid.UUID,
	photos []PostPhoto,
	createdAt time.Time,
) (*Post, error) {
	post := &Post{
		ID:        id,
		Title:     title,
		Content:   content,
		Tags:      tags,
		AuthorID:  authorID,
		CreatedAt: createdAt,
		Photos:    photos,
	}
	if err := post.Validate(); err != nil {
		return nil, err
	}
	return post, nil
}

func NewPost(
	title string,
	content Content,
	tags []string,
	authorID uuid.UUID,
	photos []PostPhoto,
) (*Post, error) {
	return CreatePost(uuid.New(), title, content, tags, authorID, photos, time.Now())
}

func CreatePostPhoto(id uuid.UUID, fileID uuid.UUID, placeNumber int) (*PostPhoto, error) {
	postPhoto := &PostPhoto{
		ID:          id,
		PlaceNumber: placeNumber,
		FileID:      fileID,
	}
	return postPhoto, nil
}

func NewPostPhoto(fileID uuid.UUID, placeNumber int) (*PostPhoto, error) {
	return CreatePostPhoto(uuid.New(), fileID, placeNumber)
}

type PostRepository interface {
	Create(ctx context.Context, post *Post) (*Post, error)
	Update(ctx context.Context, post *Post, updateFn func(*Post) (*Post, error)) (*Post, error)
	Delete(ctx context.Context, postID uuid.UUID) error
	Get(ctx context.Context, postID uuid.UUID) (*Post, error)
}
