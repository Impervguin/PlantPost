package filters

import (
	"PlantSite/internal/models/search"
	"fmt"
	"sync"
)

type PlantFilterFactory func(search.PlantFilter) (PostgresPlantFilter, error)
type PostFilterFactory func(search.PostFilter) (PostgresPostFilter, error)

type FilterRegistry struct {
	plant map[string]PlantFilterFactory
	post  map[string]PostFilterFactory
	mut   sync.RWMutex
}

func NewFilterRegistry() *FilterRegistry {
	return &FilterRegistry{
		plant: make(map[string]PlantFilterFactory),
		post:  make(map[string]PostFilterFactory),
		mut:   sync.RWMutex{},
	}
}

func (r *FilterRegistry) RegisterPlantFilter(id string, factory PlantFilterFactory) error {
	r.mut.Lock()
	r.plant[id] = factory
	r.mut.Unlock()
	return nil
}

func (r *FilterRegistry) RegisterPostFilter(id string, factory PostFilterFactory) error {
	r.mut.Lock()
	r.post[id] = factory
	r.mut.Unlock()
	return nil
}

func (r *FilterRegistry) MapPlantFilter(filter search.PlantFilter) (PostgresPlantFilter, error) {
	r.mut.RLock()
	defer r.mut.RUnlock()
	factory, ok := r.plant[filter.Identifier()]
	if !ok {
		return nil, fmt.Errorf("unknown plant filter type: %s", filter.Identifier())
	}
	return factory(filter)
}

func (r *FilterRegistry) MapPostFilter(filter search.PostFilter) (PostgresPostFilter, error) {
	r.mut.RLock()
	defer r.mut.RUnlock()
	factory, ok := r.post[filter.Identifier()]
	if !ok {
		return nil, fmt.Errorf("unknown post filter type: %s", filter.Identifier())
	}
	return factory(filter)
}
