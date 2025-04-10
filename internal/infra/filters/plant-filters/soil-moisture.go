package plantfilters

import (
	registry "PlantSite/internal/infra/filters/registry"
	pgconsts "PlantSite/internal/infra/pg-consts"
	"PlantSite/internal/models/search"
	"fmt"

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
		fmt.Sprintf("specification->>'%s'", pgconsts.JsonBSoilMoistureKey): pf.PossibleMoistures,
	}

	return filt, nil
}
