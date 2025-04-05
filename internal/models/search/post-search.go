package search

import "PlantSite/internal/models/post"

type PostSearch struct {
	filters []PostFilter
}

func NewPostSearch() *PostSearch {
	return &PostSearch{
		filters: make([]PostFilter, 0),
	}
}

func (s *PostSearch) AddFilter(filter PostFilter) {
	s.filters = append(s.filters, filter)
}

func (s *PostSearch) Filter(post *post.Post) bool {
	for _, f := range s.filters {
		if !f.Filter(post) {
			return false
		}
	}
	return true
}

func (s *PostSearch) Iterate(fn func(PostFilter) error) error {
	for _, f := range s.filters {
		if err := fn(f); err != nil {
			return err
		}
	}
	return nil
}
