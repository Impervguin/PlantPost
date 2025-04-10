package registry

import (
	"PlantSite/internal/models/plant"
	"sync"
)

type PlantSpecificationRegistry struct {
	fromDbMap     map[string]FromDB
	fromDomainMap map[string]FromDomain
	lock          sync.RWMutex
}

func (m *PlantSpecificationRegistry) FromDB(category string, json JsonB) (PlantSpecification, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if fromDb, ok := m.fromDbMap[category]; ok {
		return fromDb(json)
	}
	return nil, ErrCategoryNotFound
}

func (m *PlantSpecificationRegistry) Register(category string, fromDB FromDB, fromDomain FromDomain) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.fromDbMap[category] = fromDB
	m.fromDomainMap[category] = fromDomain
	return nil
}

func (m *PlantSpecificationRegistry) FromDomain(category string, spec plant.PlantSpecification) (PlantSpecification, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if fromDomain, ok := m.fromDomainMap[category]; ok {
		return fromDomain(spec)
	}
	return nil, ErrCategoryNotFound
}

func NewPlantSpecificationRegistry() *PlantSpecificationRegistry {
	return &PlantSpecificationRegistry{
		fromDbMap:     make(map[string]FromDB),
		fromDomainMap: make(map[string]FromDomain),
		lock:          sync.RWMutex{},
	}
}
