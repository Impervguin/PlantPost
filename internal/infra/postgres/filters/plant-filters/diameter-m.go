package plantfilters

import (
	registry "PlantSite/internal/infra/postgres/filters/registry"
	pgconsts "PlantSite/internal/infra/postgres/pg-consts"
	"PlantSite/internal/models/search"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPlantFilter(search.PlantDiameterFilterID, PlantDiameterMFilterFactory)
}

var _ registry.PlantFilterFactory = PlantDiameterMFilterFactory

func PlantDiameterMFilterFactory(f search.PlantFilter) (registry.PostgresPlantFilter, error) {
	pf, ok := f.(*search.PlantDiameterFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	// Between {min} and {max}
	filt := squirrel.And{
		squirrel.GtOrEq{pgconsts.JsonBDiameterMKey: pf.Min},
		squirrel.LtOrEq{pgconsts.JsonBDiameterMKey: pf.Max},
	}

	return filt, nil
}
