package plantfilters

import (
	registry "PlantSite/internal/infra/filters/registry"
	"PlantSite/internal/models/search"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPlantFilter(search.PlantIDsFilterID, PlantIDsFilterFactory)
}

var _ registry.PlantFilterFactory = PlantIDsFilterFactory

func PlantIDsFilterFactory(f search.PlantFilter) (registry.PostgresPlantFilter, error) {
	pf, ok := f.(*search.PlantIDsFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	// in {list}
	filt := squirrel.Eq{
		"id": pf.IDs,
	}

	return filt, nil
}
