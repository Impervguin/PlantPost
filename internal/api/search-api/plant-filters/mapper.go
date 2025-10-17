package plantfilters

import (
	"PlantSite/internal/models/search"
	"errors"
)

const (
	PlantNameFilterID            = "name"
	PlantCategoryFilterID        = "category"
	PlantLatinNameFilterID       = "latin_name"
	PlantHeightFilterID          = "height"
	PlantDiameterFilterID        = "diameter"
	PlantSoilAcidityFilterID     = "soil_acidity"
	PlantSoilMoistureFilterID    = "soil_moisture"
	PlantLightRelationFilterID   = "light_relation"
	PlantSoilTypeFilterID        = "soil_type"
	PlantWinterHardinessFilterID = "winter_hardiness"
	PlantFloweringPeriodFilterID = "flowering_period"
)

var (
	ErrInvalidFilterType = errors.New("invalid plant filter type")
)

var (
	filterConstructors = map[string]func() PlantFilter{
		PlantNameFilterID:            func() PlantFilter { return &PlantNameFilter{} },
		PlantCategoryFilterID:        func() PlantFilter { return &PlantCategoryFilter{} },
		PlantLatinNameFilterID:       func() PlantFilter { return &PlantLatinNameFilter{} },
		PlantHeightFilterID:          func() PlantFilter { return &PlantHeightFilter{} },
		PlantDiameterFilterID:        func() PlantFilter { return &PlantDiameterFilter{} },
		PlantSoilAcidityFilterID:     func() PlantFilter { return &PlantSoilAcidityFilter{} },
		PlantSoilMoistureFilterID:    func() PlantFilter { return &PlantSoilMoistureFilter{} },
		PlantLightRelationFilterID:   func() PlantFilter { return &PlantLightRelationFilter{} },
		PlantSoilTypeFilterID:        func() PlantFilter { return &PlantSoilTypeFilter{} },
		PlantWinterHardinessFilterID: func() PlantFilter { return &PlantWinterHardinessFilter{} },
		PlantFloweringPeriodFilterID: func() PlantFilter { return &PlantFloweringPeriodFilter{} },
	}
)

func ParsePlantFilter(ftype string, params map[string]interface{}) (PlantFilter, error) {
	if _, ok := filterConstructors[ftype]; !ok {
		return nil, ErrInvalidFilterType
	}
	filter := filterConstructors[ftype]()
	err := filter.Bind(params)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

func MapPlantFilters(filters []PlantFilter) ([]search.PlantFilter, error) {
	domainFilters := make([]search.PlantFilter, 0, len(filters))
	for _, f := range filters {
		domainFilter, err := f.ToDomain()
		if err != nil {
			return nil, err
		}
		domainFilters = append(domainFilters, domainFilter)
	}
	return domainFilters, nil
}
