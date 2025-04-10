package plantfilters

import (
	registry "PlantSite/internal/infra/filters/registry"
	"PlantSite/internal/models/search"
	"fmt"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPlantFilter(search.PlantNameFilterID, PlantNameFilterFactory)
}

var _ registry.PlantFilterFactory = PlantNameFilterFactory

func PlantNameFilterFactory(f search.PlantFilter) (registry.PostgresPlantFilter, error) {
	plantNameFilter, ok := f.(*search.PlantNameFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	// ILIKE %{name}%
	filt := squirrel.ILike{
		"name": fmt.Sprintf("%%%s%%", plantNameFilter.Name),
	}

	return filt, nil
}
