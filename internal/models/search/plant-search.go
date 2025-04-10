package search

import "PlantSite/internal/models/plant"

type PlantSearch struct {
	filters []PlantFilter
}

func NewPlantSearch() *PlantSearch {
	return &PlantSearch{
		filters: make([]PlantFilter, 0),
	}
}

func (s *PlantSearch) AddFilter(filter PlantFilter) {
	s.filters = append(s.filters, filter)
}

func (s *PlantSearch) Filter(pl *plant.Plant) bool {
	for _, f := range s.filters {
		if !f.Filter(pl) {
			return false
		}
	}
	return true
}

func (s *PlantSearch) Iterate(fn func(PlantFilter) error) error {
	for _, f := range s.filters {
		if err := fn(f); err != nil {
			return err
		}
	}
	return nil
}
