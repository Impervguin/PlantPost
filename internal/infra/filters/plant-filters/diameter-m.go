package plantfilters

import (
	registry "PlantSite/internal/infra/filters/registry"
	pgconsts "PlantSite/internal/infra/pg-consts"
	"PlantSite/internal/models/search"
	"fmt"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPlantFilter(search.PlantDiameterFilterID, PlantDiameterMFilterFactory)
}

var _ registry.PlantFilterFactory = PlantDiameterMFilterFactory

func PlantDiameterMFilterFactory(f search.PlantFilter) (registry.PostgresPlantFilter, error) {
	pf, ok := f.(*search.PlantDiameterFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	// Between {min} and {max}
	filt := squirrel.And{
		squirrel.GtOrEq{fmt.Sprintf("specification->'%s'", pgconsts.JsonBDiameterMKey): pf.Min},
		squirrel.LtOrEq{fmt.Sprintf("specification->'%s'", pgconsts.JsonBDiameterMKey): pf.Max},
	}

	return filt, nil
}
