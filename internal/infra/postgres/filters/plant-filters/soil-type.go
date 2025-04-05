package plantfilters

import (
	registry "PlantSite/internal/infra/postgres/filters/registry"
	pgconsts "PlantSite/internal/infra/postgres/pg-consts"
	"PlantSite/internal/models/search"

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
		pgconsts.JsonBSoilTypeKey: pf.PossibleSoilTypes,
	}

	return filt, nil
}
