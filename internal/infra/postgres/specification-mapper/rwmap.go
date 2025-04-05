package plantstorage

import (
	"PlantSite/internal/models/plant"
	"sync"
)

type RWMapPlantSpecificationMapper struct {
	fromDbMap     map[string]FromDb
	fromDomainMap map[string]FromDomain
	lock          sync.RWMutex
}

func (m *RWMapPlantSpecificationMapper) FromDB(category string, json JsonB) (PlantSpecification, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if fromDb, ok := m.fromDbMap[category]; ok {
		return fromDb(json)
	}
	return nil, ErrCategoryNotFound
}

func (m *RWMapPlantSpecificationMapper) Register(category string, fromDB FromDb, fromDomain FromDomain) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.fromDbMap[category] = fromDB
	m.fromDomainMap[category] = fromDomain
}

func (m *RWMapPlantSpecificationMapper) FromDomain(category string, spec plant.PlantSpecification) (PlantSpecification, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if fromDomain, ok := m.fromDomainMap[category]; ok {
		return fromDomain(spec)
	}
	return nil, ErrCategoryNotFound
}

func RWMapPlantSpecificationMapperFactory() PlantSpecificationMapper {
	return &RWMapPlantSpecificationMapper{
		fromDbMap:     make(map[string]FromDb),
		fromDomainMap: make(map[string]FromDomain),
		lock:          sync.RWMutex{},
	}
}
