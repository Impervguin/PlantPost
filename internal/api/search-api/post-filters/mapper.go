package postfilters

import (
	"PlantSite/internal/models/search"
	"errors"
)

var (
	ErrInvalidFilterType = errors.New("invalid filter type")
)

const (
	PostTitleFilterID         = "title"
	PostTitleContainsFilterID = "title_contains"
	PostTagFilterID           = "tag"
	PostAuthorFilterID        = "author"
)

func ParsePostFilter(ftype string, params map[string]interface{}) (PostFilter, error) {
	switch ftype {
	case PostTitleFilterID:
		var f PostTitleFilter
		if err := f.Bind(params); err != nil {
			return nil, err
		}
		return &f, nil
	case PostTitleContainsFilterID:
		var f PostTitleContainsFilter
		if err := f.Bind(params); err != nil {
			return nil, err
		}
		return &f, nil
	case PostTagFilterID:
		var f PostTagFilter
		if err := f.Bind(params); err != nil {
			return nil, err
		}
		return &f, nil
	case PostAuthorFilterID:
		var f PostAuthorFilter
		if err := f.Bind(params); err != nil {
			return nil, err
		}
		return &f, nil
	default:
		return nil, ErrInvalidFilterType
	}
}

func MapPostFilters(filters []PostFilter) ([]search.PostFilter, error) {
	domainFilters := make([]search.PostFilter, 0, len(filters))
	for _, f := range filters {
		domainFilter, err := f.ToDomain()
		if err != nil {
			return nil, err
		}
		domainFilters = append(domainFilters, domainFilter)
	}
	return domainFilters, nil
}
