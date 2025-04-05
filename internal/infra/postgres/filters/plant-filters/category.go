package plantfilters

import (
	"PlantSite/internal/models/search"
	"fmt"

	registry "PlantSite/internal/infra/postgres/filters/registry"

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

	fmt.Println(filt.ToSql())

	return filt, nil
}
