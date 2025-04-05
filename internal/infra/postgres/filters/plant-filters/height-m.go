package plantfilters

import (
	registry "PlantSite/internal/infra/postgres/filters/registry"
	pgconsts "PlantSite/internal/infra/postgres/pg-consts"
	"PlantSite/internal/models/search"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPlantFilter(search.PlantHeightFilterID, PlantHeightMFilterFactory)
}

var _ registry.PlantFilterFactory = PlantHeightMFilterFactory

func PlantHeightMFilterFactory(f search.PlantFilter) (registry.PostgresPlantFilter, error) {
	pf, ok := f.(*search.PlantHeightFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	// Between {min} and {max}
	filt := squirrel.And{
		squirrel.GtOrEq{pgconsts.JsonBHeightMKey: pf.Min},
		squirrel.LtOrEq{pgconsts.JsonBHeightMKey: pf.Max},
	}

	return filt, nil
}
