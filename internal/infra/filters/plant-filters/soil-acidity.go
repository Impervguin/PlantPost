package plantfilters

import (
	registry "PlantSite/internal/infra/filters/registry"
	pgconsts "PlantSite/internal/infra/pg-consts"
	"PlantSite/internal/models/search"
	"fmt"

	"github.com/Masterminds/squirrel"
)

func init() {
	registry.RegisterPlantFilter(search.PlantSoilAcidityFilterID, PlantSoilAcidityFilterFactory)
}

var _ registry.PlantFilterFactory = PlantSoilAcidityFilterFactory

func PlantSoilAcidityFilterFactory(f search.PlantFilter) (registry.PostgresPlantFilter, error) {
	pf, ok := f.(*search.PlantSoilAcidityFilter)
	if !ok {
		return nil, registry.ErrInvalidFilterType
	}

	// Between {min} and {max}
	filt := squirrel.And{
		squirrel.GtOrEq{fmt.Sprintf("specification->'%s'", pgconsts.JsonBSoilAcidityKey): pf.Min},
		squirrel.LtOrEq{fmt.Sprintf("specification->'%s'", pgconsts.JsonBSoilAcidityKey): pf.Max},
	}

	return filt, nil
}
