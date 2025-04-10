package plantfilters

import (
	registry "PlantSite/internal/infra/filters/registry"
	pgconsts "PlantSite/internal/infra/pg-consts"
	"PlantSite/internal/models/search"
	"fmt"

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
		fmt.Sprintf("specification->>'%s'", pgconsts.JsonBLightRelationKey): pf.PossibleRelations,
	}

	return filt, nil
}
