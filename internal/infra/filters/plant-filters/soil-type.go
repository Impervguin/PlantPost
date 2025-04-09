package plantfilters

import (
	registry "PlantSite/internal/infra/filters/registry"
	pgconsts "PlantSite/internal/infra/pg-consts"
	"PlantSite/internal/models/search"
	"fmt"

	"github.com/Masterminds/squirrel"
)

// const (
// 	PlantSoilTypeFilterID        = "PlantSoilTypeFilter"
// 	PlantWinterHardinessFilterID = "PlantWinterHardinessFilter"
// 	PlantFloweringPeriodFilterID = "PlantFloweringPeriodFilter"
// )

func init() {
	registry.RegisterPlantFilter(search.PlantSoilTypeFilterID, PlantSoilTypeFilterFactory)
}

var _ registry.PlantFilterFactory = PlantSoilTypeFilterFactory

func PlantSoilTypeFilterFactory(f search.PlantFilter) (registry.PostgresPlantFilter, error) {
	pf, ok := f.(*search.PlantSoilTypeFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	// in {list}
	filt := squirrel.Eq{
		fmt.Sprintf("specification->>'%s'", pgconsts.JsonBSoilTypeKey): pf.PossibleSoilTypes,
	}

	return filt, nil
}
