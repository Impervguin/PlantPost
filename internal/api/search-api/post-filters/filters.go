package postfilters

import (
	"PlantSite/internal/models/search"
	"fmt"

	"github.com/google/uuid"
)

type PostFilter interface {
	ToDomain() (search.PostFilter, error)
	Bind(params map[string]interface{}) error
	Type() string
}

type PostTitleFilter struct {
	Title string `json:"title" form:"title" binding:"required"`
}

func (f *PostTitleFilter) ToDomain() (search.PostFilter, error) {
	return search.NewPostTitleFilter(f.Title), nil
}

func (f *PostTitleFilter) Bind(params map[string]interface{}) error {
	title, ok := params["title"]
	if !ok {
		return fmt.Errorf("title not found")
	}
	f.Title, ok = title.(string)
	if !ok {
		return fmt.Errorf("title is not a string")
	}
	return nil
}

func (f *PostTitleFilter) Type() string {
	return PostTitleFilterID
}

type PostTitleContainsFilter struct {
	Part string `json:"part" form:"part" binding:"required"`
}

func (f *PostTitleContainsFilter) ToDomain() (search.PostFilter, error) {
	return search.NewPostTitleContainsFilter(f.Part), nil
}

func (f *PostTitleContainsFilter) Bind(params map[string]interface{}) error {
	part, ok := params["part"]
	if !ok {
		return fmt.Errorf("part not found")
	}
	f.Part, ok = part.(string)
	if !ok {
		return fmt.Errorf("part is not a string")
	}
	return nil
}

func (f *PostTitleContainsFilter) Type() string {
	return PostTitleContainsFilterID
}

type PostTagFilter struct {
	Tags []string `json:"tags" form:"tags" binding:"required"`
}

func (f *PostTagFilter) ToDomain() (search.PostFilter, error) {
	return search.NewPostTagFilter(f.Tags), nil
}

func (f *PostTagFilter) Bind(params map[string]interface{}) error {
	tags, ok := params["tags"]
	if !ok {
		return fmt.Errorf("tags not found")
	}
	tagsList, ok := tags.([]interface{})
	if !ok {
		return fmt.Errorf("tags is not a list")
	}
	f.Tags = make([]string, 0, len(tagsList))
	for _, tag := range tagsList {
		tagStr, ok := tag.(string)
		if !ok {
			return fmt.Errorf("tag is not a string")
		}
		f.Tags = append(f.Tags, tagStr)
	}
	return nil
}

func (f *PostTagFilter) Type() string {
	return PostTagFilterID
}

type PostAuthorFilter struct {
	AuthorID string `json:"author_id" form:"author_id" binding:"required"`
}

func (f *PostAuthorFilter) ToDomain() (search.PostFilter, error) {
	id, err := uuid.Parse(f.AuthorID)
	if err != nil {
		return nil, err
	}
	return search.NewPostAuthorFilter(id), nil
}

func (f *PostAuthorFilter) Bind(params map[string]interface{}) error {
	authorID, ok := params["author_id"]
	if !ok {
		return fmt.Errorf("author_id not found")
	}
	authorIDStr, ok := authorID.(string)
	if !ok {
		return fmt.Errorf("author_id is not a string")
	}
	f.AuthorID = authorIDStr
	return nil
}

func (f *PostAuthorFilter) Type() string {
	return PostAuthorFilterID
}
