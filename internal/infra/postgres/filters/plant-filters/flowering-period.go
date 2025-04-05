package plantfilters

import (
	registry "PlantSite/internal/infra/postgres/filters/registry"
	pgconsts "PlantSite/internal/infra/postgres/pg-consts"
	"PlantSite/internal/models/search"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPlantFilter(search.PlantFloweringPeriodFilterID, PlantFloweringPeriodFilterFactory)
}

var _ registry.PlantFilterFactory = PlantFloweringPeriodFilterFactory

func PlantFloweringPeriodFilterFactory(f search.PlantFilter) (registry.PostgresPlantFilter, error) {
	pf, ok := f.(*search.PlantFloweringPeriodFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	// in {list}
	filt := squirrel.Eq{
		pgconsts.JsonBFloweringPeriodKey: pf.PossibleFloweringPeriods,
	}

	return filt, nil
}
