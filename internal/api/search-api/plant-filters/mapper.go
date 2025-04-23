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

func ParsePlantFilter(ftype string, params map[string]interface{}) (PlantFilter, error) {
	switch ftype {
	case PlantNameFilterID:
		var f PlantNameFilter
		if err := f.Bind(params); err != nil {
			return nil, err
		}
		return &f, nil
	case PlantCategoryFilterID:
		var f PlantCategoryFilter
		if err := f.Bind(params); err != nil {
			return nil, err
		}
		return &f, nil
	case PlantLatinNameFilterID:
		var f PlantLatinNameFilter
		if err := f.Bind(params); err != nil {
			return nil, err
		}
		return &f, nil
	case PlantHeightFilterID:
		var f PlantHeightFilter
		if err := f.Bind(params); err != nil {
			return nil, err
		}
		return &f, nil
	case PlantDiameterFilterID:
		var f PlantDiameterFilter
		if err := f.Bind(params); err != nil {
			return nil, err
		}
		return &f, nil
	case PlantSoilAcidityFilterID:
		var f PlantSoilAcidityFilter
		if err := f.Bind(params); err != nil {
			return nil, err
		}
		return &f, nil
	case PlantSoilMoistureFilterID:
		var f PlantSoilMoistureFilter
		if err := f.Bind(params); err != nil {
			return nil, err
		}
		return &f, nil
	case PlantLightRelationFilterID:
		var f PlantLightRelationFilter
		if err := f.Bind(params); err != nil {
			return nil, err
		}
		return &f, nil
	case PlantSoilTypeFilterID:
		var f PlantSoilTypeFilter
		if err := f.Bind(params); err != nil {
			return nil, err
		}
		return &f, nil
	case PlantWinterHardinessFilterID:
		var f PlantWinterHardinessFilter
		if err := f.Bind(params); err != nil {
			return nil, err
		}
		return &f, nil
	case PlantFloweringPeriodFilterID:
		var f PlantFloweringPeriodFilter
		if err := f.Bind(params); err != nil {
			return nil, err
		}
		return &f, nil
	default:
		return nil, ErrInvalidFilterType
	}
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
