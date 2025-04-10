package plantfilters

import (
	registry "PlantSite/internal/infra/filters/registry"
	pgconsts "PlantSite/internal/infra/pg-consts"
	"PlantSite/internal/models/search"
	"fmt"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPlantFilter(search.PlantWinterHardinessFilterID, PlantWinterHardinessFilterFactory)
}

var _ registry.PlantFilterFactory = PlantWinterHardinessFilterFactory

func PlantWinterHardinessFilterFactory(f search.PlantFilter) (registry.PostgresPlantFilter, error) {
	pf, ok := f.(*search.PlantHardinessFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	// Between {min} and {max}
	filt := squirrel.And{
		squirrel.GtOrEq{fmt.Sprintf("specification->'%s'", pgconsts.JsonBWinterHardinessKey): pf.Min},
		squirrel.LtOrEq{fmt.Sprintf("specification->'%s'", pgconsts.JsonBWinterHardinessKey): pf.Max},
	}

	return filt, nil
}
