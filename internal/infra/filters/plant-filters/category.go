package plantfilters

import (
	"PlantSite/internal/models/search"

	registry "PlantSite/internal/infra/filters/registry"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPlantFilter(search.PlantCategoryFilterID, PlantCategoryFilterFactory)
}

var _ registry.PlantFilterFactory = PlantCategoryFilterFactory

func PlantCategoryFilterFactory(f search.PlantFilter) (registry.PostgresPlantFilter, error) {
	pf, ok := f.(*search.PlantCategoryFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	filt := squirrel.Eq{
		"category": pf.Category,
	}

	return filt, nil
}
