package plantfilters

import (
	registry "PlantSite/internal/infra/postgres/filters/registry"
	pgconsts "PlantSite/internal/infra/postgres/pg-consts"
	"PlantSite/internal/models/search"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPlantFilter(search.PlantWinterHardinessFilterID, PlantWinterHardinessFilterFactory)
}

var _ registry.PlantFilterFactory = PlantWinterHardinessFilterFactory

func PlantWinterHardinessFilterFactory(f search.PlantFilter) (registry.PostgresPlantFilter, error) {
	pf, ok := f.(*search.PlantHardinessFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	// Between {min} and {max}
	filt := squirrel.And{
		squirrel.GtOrEq{pgconsts.JsonBWinterHardinessKey: pf.Min},
		squirrel.LtOrEq{pgconsts.JsonBWinterHardinessKey: pf.Max},
	}

	return filt, nil
}
