package post

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	id      uuid.UUID
	title   string
	content Content
	tags    []string
	photos  PostPhotos

	authorID  uuid.UUID
	updatedAt time.Time
	createdAt time.Time
}

func (p *Post) Validate() error {
	if p.id == uuid.Nil {
		return fmt.Errorf("post ID cannot be empty")
	}
	if p.title == "" {
		return fmt.Errorf("post title cannot be empty")
	}
	if err := p.content.Validate(); err != nil {
		return fmt.Errorf("post content validation failed: %v", err)
	}
	if p.authorID == uuid.Nil {
		return fmt.Errorf("author ID cannot be empty")
	}
	if p.tags == nil {
		return fmt.Errorf("tags cannot be nil")
	}
	if err := p.photos.Validate(); err != nil {
		return fmt.Errorf("post photo validation failed: %v", err)
	}

	for _, tag := range p.tags {
		if tag == "" {
			return fmt.Errorf("empty tag")
		}
	}

	if p.updatedAt.Before(p.createdAt) {
		return fmt.Errorf("post update date cannot be before creation date")
	}

	if p.updatedAt.After(time.Now()) {
		return fmt.Errorf("post update date cannot be in the future")
	}

	return nil
}

func CreatePost(
	id uuid.UUID,
	title string,
	content Content,
	tags []string,
	authorID uuid.UUID,
	photos PostPhotos,
	createdAt time.Time,
	updatedAt time.Time,
) (*Post, error) {
	post := &Post{
		id:        id,
		title:     title,
		content:   content,
		tags:      tags,
		authorID:  authorID,
		createdAt: createdAt,
		photos:    photos,
		updatedAt: updatedAt,
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
	photos *PostPhotos,
) (*Post, error) {
	t := time.Now()
	return CreatePost(uuid.New(), title, content, tags, authorID, *photos, t, t)
}

func (p Post) ID() uuid.UUID {
	return p.id
}

func (p Post) AuthorID() uuid.UUID {
	return p.authorID
}

func (p Post) Title() string {
	return p.title
}

func (p Post) Content() Content {
	return p.content
}

func (p Post) Tags() []string {
	return p.tags
}

func (p Post) Photos() PostPhotos {
	return p.photos
}

func (p Post) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p Post) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Post) UpdateContent(content Content) error {
	if err := content.Validate(); err != nil {
		return err
	}
	p.content = content
	p.updatedAt = time.Now()
	return nil
}

func (p *Post) UpdateTitle(title string) error {
	p.title = title
	p.updatedAt = time.Now()
	return nil
}

func (p *Post) UpdateTags(tags []string) error {
	p.tags = tags
	p.updatedAt = time.Now()
	return nil
}

func (p *Post) AddTag(tag string) error {
	if len(p.tags) >= MaximumTagCount {
		return fmt.Errorf("maximum number of tags exceeded")
	}
	if slices.Contains(p.tags, tag) {
		return fmt.Errorf("tag already exists: %s", tag)
	}
	p.tags = append(p.tags, tag)
	return nil
}

func (p *Post) RemoveTag(tag string) error {
	index := slices.Index(p.tags, tag)
	if index == -1 {
		return fmt.Errorf("tag not found: %s", tag)
	}
	if index == 0 {
		p.tags = p.tags[index+1:]
		return nil
	}
	p.tags = append(p.tags[:index-1], p.tags[index+1:]...)
	return nil
}

func (p *Post) AddPhoto(photo *PostPhoto) error {
	return p.photos.Add(photo)
}

type PostRepository interface {
	Create(ctx context.Context, post *Post) (*Post, error)
	Update(ctx context.Context, id uuid.UUID, updateFn func(*Post) (*Post, error)) (*Post, error)
	Delete(ctx context.Context, postID uuid.UUID) error
	Get(ctx context.Context, postID uuid.UUID) (*Post, error)
}
