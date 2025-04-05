package search

import (
	"PlantSite/internal/models/post"
	"slices"
	"strings"

	"github.com/google/uuid"
)

type PostFilter interface {
	Filter(p *post.Post) bool
	Identifier() string
}

type PostTitleFilter struct {
	Title string
}

var _ PostFilter = &PostTitleFilter{}

func (p *PostTitleFilter) Identifier() string {
	return PostTitleFilterID
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

var _ PostFilter = &PostTitleContainsFilter{}

func (p *PostTitleContainsFilter) Identifier() string {
	return PostTitleContainsFilterID
}

func NewPostTitleContainsFilter(part string) *PostTitleContainsFilter {
	return &PostTitleContainsFilter{Part: part}
}

func (p *PostTitleContainsFilter) Filter(post *post.Post) bool {
	return strings.Contains(strings.ToLower(post.Title()), strings.ToLower(p.Part))
}

type PostTagFilter struct {
	Tags []string
}

var _ PostFilter = &PostTagFilter{}

func (p *PostTagFilter) Identifier() string {
	return PostTagFilterID
}

func NewPostTagFilter(tags []string) *PostTagFilter {
	return &PostTagFilter{Tags: tags}
}

func (p *PostTagFilter) Filter(post *post.Post) bool {
	for _, tag := range p.Tags {
		if slices.Contains(post.Tags(), tag) {
			return true
		}
	}
	return false
}

type PostAuthorFilter struct {
	AuthorID uuid.UUID
}

var _ PostFilter = &PostAuthorFilter{}

func (p *PostAuthorFilter) Identifier() string {
	return PostAuthorFilterID
}

func NewPostAuthorFilter(authorID uuid.UUID) *PostAuthorFilter {
	return &PostAuthorFilter{AuthorID: authorID}
}

func (p *PostAuthorFilter) Filter(post *post.Post) bool {
	return post.AuthorID() == p.AuthorID
}
