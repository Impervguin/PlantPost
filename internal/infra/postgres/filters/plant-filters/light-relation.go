package plantfilters

import (
	registry "PlantSite/internal/infra/postgres/filters/registry"
	pgconsts "PlantSite/internal/infra/postgres/pg-consts"
	"PlantSite/internal/models/search"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPlantFilter(search.PlantLightRelationFilterID, PlantLightRelationFilterFactory)
}

var _ registry.PlantFilterFactory = PlantLightRelationFilterFactory

func PlantLightRelationFilterFactory(f search.PlantFilter) (registry.PostgresPlantFilter, error) {
	pf, ok := f.(*search.PlantLightRelationFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	// in {list}
	filt := squirrel.Eq{
		pgconsts.JsonBLightRelationKey: pf.PossibleRelations,
	}

	return filt, nil
}
