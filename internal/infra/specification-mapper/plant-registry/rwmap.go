package registry

import (
	"PlantSite/internal/models/plant"
	"sync"
)

type PlantSpecificationRegistry struct {
	fromDBMap     map[string]FromDB
	fromDomainMap map[string]FromDomain
	lock          sync.RWMutex
}

func NewPlantSpecificationRegistry() *PlantSpecificationRegistry {
	return &PlantSpecificationRegistry{
		fromDBMap:     make(map[string]FromDB),
		fromDomainMap: make(map[string]FromDomain),
		lock:          sync.RWMutex{},
	}
}

func (m *PlantSpecificationRegistry) FromDB(category string, json JSONB) (PlantSpecification, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if fromDB, ok := m.fromDBMap[category]; ok {
		return fromDB(json)
	}
	return nil, ErrCategoryNotFound
}

func (m *PlantSpecificationRegistry) Register(category string, fromDB FromDB, fromDomain FromDomain) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.fromDBMap[category] = fromDB
	m.fromDomainMap[category] = fromDomain
	return nil
}

func (m *PlantSpecificationRegistry) FromDomain(
	category string,
	spec plant.PlantSpecification,
) (PlantSpecification, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if fromDomain, ok := m.fromDomainMap[category]; ok {
		return fromDomain(spec)
	}
	return nil, ErrCategoryNotFound
}
