package plantfilters

import (
	registry "PlantSite/internal/infra/filters/registry"
	pgconsts "PlantSite/internal/infra/pg-consts"
	"PlantSite/internal/models/search"
	"fmt"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPlantFilter(search.PlantHeightFilterID, PlantHeightMFilterFactory)
}

var _ registry.PlantFilterFactory = PlantHeightMFilterFactory

func PlantHeightMFilterFactory(f search.PlantFilter) (registry.PostgresPlantFilter, error) {
	pf, ok := f.(*search.PlantHeightFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	// Between {min} and {max}
	filt := squirrel.And{
		squirrel.GtOrEq{fmt.Sprintf("specification->'%s'", pgconsts.JsonBHeightMKey): pf.Min},
		squirrel.LtOrEq{fmt.Sprintf("specification->'%s'", pgconsts.JsonBHeightMKey): pf.Max},
	}

	return filt, nil
}
