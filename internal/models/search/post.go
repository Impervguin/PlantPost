package search

import (
	"PlantSite/internal/models/post"
	"slices"
	"strings"

	"github.com/google/uuid"
)

type PostFilter interface {
	Filter(p *post.Post) bool
}

type PostTitleFilter struct {
	Title string
}

func NewPostTitleFilter(title string) *PostTitleFilter {
	return &PostTitleFilter{Title: title}
}

func (p *PostTitleFilter) Filter(post *post.Post) bool {
	return post.Title() == p.Title
}

type PostTitleContainsFilter struct {
	Part string
}

func NewPostTitleContainsFilter(part string) *PostTitleContainsFilter {
	return &PostTitleContainsFilter{Part: part}
}

func (p *PostTitleContainsFilter) Filter(post *post.Post) bool {
	return strings.Contains(strings.ToLower(post.Title()), strings.ToLower(p.Part))
}

type PostTagFilter struct {
	Tag string
}

func NewPostTagFilter(tag string) *PostTagFilter {
	return &PostTagFilter{Tag: tag}
}

func (p *PostTagFilter) Filter(post *post.Post) bool {
	return slices.Contains(post.Tags(), p.Tag)
}

type PostAuthorFilter struct {
	AuthorID uuid.UUID
}

func NewPostAuthorFilter(authorID uuid.UUID) *PostAuthorFilter {
	return &PostAuthorFilter{AuthorID: authorID}
}

func (p *PostAuthorFilter) Filter(post *post.Post) bool {
	return post.AuthorID() == p.AuthorID
}
