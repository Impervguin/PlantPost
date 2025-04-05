package plantfilters

import (
	registry "PlantSite/internal/infra/postgres/filters/registry"
	pgconsts "PlantSite/internal/infra/postgres/pg-consts"
	"PlantSite/internal/models/search"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPlantFilter(search.PlantSoilMoistureFilterID, PlantSoilMoistureFilterFactory)
}

var _ registry.PlantFilterFactory = PlantSoilMoistureFilterFactory

func PlantSoilMoistureFilterFactory(f search.PlantFilter) (registry.PostgresPlantFilter, error) {
	pf, ok := f.(*search.PlantSoilMoistureFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	// in {list}
	filt := squirrel.Eq{
		pgconsts.JsonBSoilMoistureKey: pf.PossibleMoistures,
	}

	return filt, nil
}
