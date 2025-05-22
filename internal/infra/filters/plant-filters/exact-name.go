package plantfilters

import (
	registry "PlantSite/internal/infra/filters/registry"
	"PlantSite/internal/models/search"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPlantFilter(search.ExactPlantNameFilterID, ExactPlantNameFilterFactory)
}

var _ registry.PlantFilterFactory = ExactPlantNameFilterFactory

func ExactPlantNameFilterFactory(f search.PlantFilter) (registry.PostgresPlantFilter, error) {
	exactPlantNameFilter, ok := f.(*search.ExactPlantNameFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	// ILIKE %{name}%
	filt := squirrel.Eq{
		"name": exactPlantNameFilter.Name,
	}

	return filt, nil
}
