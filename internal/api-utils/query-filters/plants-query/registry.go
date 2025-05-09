package plantsquery

import (
	"PlantSite/internal/models/search"
	"fmt"
	"sync"
)

type QueryPlantFilterParser func(string) (search.PlantFilter, error)

type PlantFilterParam string

const (
	PlantNameFilterParam            PlantFilterParam = "name"
	PlantCategoryFilterParam        PlantFilterParam = "category"
	PlantLatinNameFilterParam       PlantFilterParam = "latin_name"
	PlantHeightFilterParam          PlantFilterParam = "height"
	PlantDiameterFilterParam        PlantFilterParam = "diameter"
	PlantSoilAcidityFilterParam     PlantFilterParam = "soil_acidity"
	PlantSoilMoistureFilterParam    PlantFilterParam = "soil_moisture"
	PlantLightRelationFilterParam   PlantFilterParam = "light_relation"
	PlantSoilTypeFilterParam        PlantFilterParam = "soil_type"
	PlantWinterHardinessFilterParam PlantFilterParam = "winter_hardiness"
	PlantFloweringPeriodFilterParam PlantFilterParam = "flowering_period"
)

type QueryPlantFilterRegistry struct {
	parsers map[PlantFilterParam]QueryPlantFilterParser
	mut     sync.RWMutex
}

func newQueryPlantFilterRegistry() *QueryPlantFilterRegistry {
	return &QueryPlantFilterRegistry{
		parsers: make(map[PlantFilterParam]QueryPlantFilterParser),
		mut:     sync.RWMutex{},
	}
}

func (r *QueryPlantFilterRegistry) register(name PlantFilterParam, parser QueryPlantFilterParser) {
	r.mut.Lock()
	defer r.mut.Unlock()
	r.parsers[name] = parser
}

func (r *QueryPlantFilterRegistry) parse(name PlantFilterParam, queryValue string) (search.PlantFilter, error) {
	r.mut.RLock()
	defer r.mut.RUnlock()
	parser, ok := r.parsers[name]
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrParserNotFound, name)
	}
	return parser(queryValue)
}
