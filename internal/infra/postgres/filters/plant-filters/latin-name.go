package plantfilters

import (
	registry "PlantSite/internal/infra/postgres/filters/registry"
	"PlantSite/internal/models/search"
	"fmt"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPlantFilter(search.PlantLatinNameFilterID, PlantLatinNameFilterFactory)
}

var _ registry.PlantFilterFactory = PlantLatinNameFilterFactory

func PlantLatinNameFilterFactory(f search.PlantFilter) (registry.PostgresPlantFilter, error) {
	pf, ok := f.(*search.PlantLatinNameFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	// ILIKE %{latin_name}%
	filt := squirrel.ILike{
		"latin_name": fmt.Sprintf("%%%s%%", pf.LatinName),
	}

	return filt, nil
}
