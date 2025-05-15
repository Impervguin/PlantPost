package plantsquery

import (
	"PlantSite/internal/models/search"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	registry = newQueryPlantFilterRegistry()
	once     = sync.Once{}
)

func registryInit() {
	once.Do(func() {
		registry.register(PlantNameFilterParam, parsePlantNameFilterfunc)
		registry.register(PlantCategoryFilterParam, parsePlantCategoryFilterfunc)
		registry.register(PlantLatinNameFilterParam, parsePlantLatinNameFilterfunc)
		registry.register(PlantHeightFilterParam, parsePlantHeightFilterfunc)
		registry.register(PlantDiameterFilterParam, parsePlantDiameterFilterfunc)
		registry.register(PlantSoilAcidityFilterParam, parsePlantSoilAcidityFilterfunc)
		registry.register(PlantSoilMoistureFilterParam, parsePlantSoilMoistureFilterfunc)
		registry.register(PlantLightRelationFilterParam, parsePlantLightRelationFilterfunc)
		registry.register(PlantSoilTypeFilterParam, parseSoilTypeFilterfunc)
		registry.register(PlantWinterHardinessFilterParam, parsePlantWinterHardinessFilterfunc)
		registry.register(PlantFloweringPeriodFilterParam, parsePlantFloweringPeriodFilterfunc)
	})
}

func ParseQueryPlantFilter(filterType PlantFilterParam, queryValue string) (search.PlantFilter, error) {
	registryInit()
	return registry.parse(filterType, queryValue)
}

func ParseGinQueryPlantSearch(c *gin.Context) (*search.PlantSearch, error) {
	registryInit()
	params := c.Request.URL.Query()
	srch := search.NewPlantSearch()
	for filterType, query := range params {
		if len(query) == 0 {
			continue
		}
		for _, q := range query {
			filter, err := registry.parse(PlantFilterParam(filterType), q)
			if err != nil {
				return &search.PlantSearch{}, fmt.Errorf("can't parse filter in search: %w", err)
			}
			srch.AddFilter(filter)
		}
	}
	return srch, nil
}
