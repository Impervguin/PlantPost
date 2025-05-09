package plantsquery

import (
	"PlantSite/internal/models/plant"
	"PlantSite/internal/models/search"
	"fmt"
	"strconv"
	"strings"
)

func parsePlantNameFilterfunc(queryValue string) (search.PlantFilter, error) {
	filt := search.NewPlantNameFilter(queryValue)
	if filt == nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantNameFilterParam, queryValue)
	}
	return filt, nil
}

func parsePlantLatinNameFilterfunc(queryValue string) (search.PlantFilter, error) {
	filt := search.NewPlantLatinNameFilter(queryValue)
	if filt == nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantLatinNameFilterParam, queryValue)
	}
	return filt, nil
}

func parsePlantCategoryFilterfunc(queryValue string) (search.PlantFilter, error) {
	filt := search.NewPlantCategoryFilter(queryValue)
	if filt == nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantCategoryFilterParam, queryValue)
	}
	return filt, nil
}

func parsePlantHeightFilterfunc(queryValue string) (search.PlantFilter, error) {
	// {min}-{max} format
	parts := strings.Split(queryValue, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantHeightFilterParam, queryValue)
	}
	min, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantHeightFilterParam, queryValue)
	}
	max, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantHeightFilterParam, queryValue)
	}
	filt := search.NewPlantHeightFilter(min, max)
	if filt == nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantHeightFilterParam, queryValue)
	}
	return filt, nil
}

func parsePlantDiameterFilterfunc(queryValue string) (search.PlantFilter, error) {
	// {min}-{max} format
	parts := strings.Split(queryValue, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantDiameterFilterParam, queryValue)
	}
	min, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantDiameterFilterParam, queryValue)
	}
	max, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantDiameterFilterParam, queryValue)
	}
	filt := search.NewPlantDiameterFilter(min, max)
	if filt == nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantDiameterFilterParam, queryValue)
	}
	return filt, nil
}

func parsePlantSoilAcidityFilterfunc(queryValue string) (search.PlantFilter, error) {
	// {min}-{max} format
	parts := strings.Split(queryValue, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantSoilAcidityFilterParam, queryValue)
	}
	min, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantSoilAcidityFilterParam, queryValue)
	}
	max, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantSoilAcidityFilterParam, queryValue)
	}
	filt := search.NewSoilAcidityFilter(plant.SoilAcidity(min), plant.SoilAcidity(max))
	if filt == nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantSoilAcidityFilterParam, queryValue)
	}
	return filt, nil
}

func parsePlantSoilMoistureFilterfunc(queryValue string) (search.PlantFilter, error) {
	// var1,var2,... format
	vars := strings.Split(queryValue, ",")
	if len(vars) == 0 {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantSoilMoistureFilterParam, queryValue)
	}
	possibleMoistures := make([]plant.SoilMoisture, 0, len(vars))
	for _, moisture := range vars {
		moisture = strings.TrimSpace(moisture)
		possibleMoistures = append(possibleMoistures, plant.SoilMoisture(moisture))
	}
	filt := search.NewSoilMoistureFilter(possibleMoistures)
	if filt == nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantSoilMoistureFilterParam, queryValue)
	}
	return filt, nil
}

func parsePlantLightRelationFilterfunc(queryValue string) (search.PlantFilter, error) {
	// var1,var2,... format
	vars := strings.Split(queryValue, ",")
	if len(vars) == 0 {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantLightRelationFilterParam, queryValue)
	}
	possibleRelations := make([]plant.LightRelation, 0, len(vars))
	for _, relation := range vars {
		relation = strings.TrimSpace(relation)
		possibleRelations = append(possibleRelations, plant.LightRelation(relation))
	}
	filt := search.NewLightRelationFilter(possibleRelations)
	if filt == nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantLightRelationFilterParam, queryValue)
	}
	return filt, nil
}

func parseSoilTypeFilterfunc(queryValue string) (search.PlantFilter, error) {
	// var1,var2,... format
	vars := strings.Split(queryValue, ",")
	if len(vars) == 0 {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantSoilTypeFilterParam, queryValue)
	}
	possibleSoilTypes := make([]plant.Soil, 0, len(vars))
	for _, soilType := range vars {
		soilType = strings.TrimSpace(soilType)
		possibleSoilTypes = append(possibleSoilTypes, plant.Soil(soilType))
	}
	filt := search.NewSoilTypeFilter(possibleSoilTypes)
	if filt == nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantSoilTypeFilterParam, queryValue)
	}
	return filt, nil
}

func parsePlantWinterHardinessFilterfunc(queryValue string) (search.PlantFilter, error) {
	// {min}-{max} format
	parts := strings.Split(queryValue, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantWinterHardinessFilterParam, queryValue)
	}
	min, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantWinterHardinessFilterParam, queryValue)
	}
	max, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantWinterHardinessFilterParam, queryValue)
	}
	filt := search.NewWinterHardinessFilter(plant.WinterHardiness(min), plant.WinterHardiness(max))
	if filt == nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantWinterHardinessFilterParam, queryValue)
	}
	return filt, nil
}

func parsePlantFloweringPeriodFilterfunc(queryValue string) (search.PlantFilter, error) {
	// var1,var2,... format
	vars := strings.Split(queryValue, ",")
	if len(vars) == 0 {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantFloweringPeriodFilterParam, queryValue)
	}
	possibleFloweringPeriods := make([]plant.FloweringPeriod, 0, len(vars))
	for _, floweringPeriod := range vars {
		floweringPeriod = strings.TrimSpace(floweringPeriod)
		possibleFloweringPeriods = append(possibleFloweringPeriods, plant.FloweringPeriod(floweringPeriod))
	}
	filt := search.NewFloweringPeriodFilter(possibleFloweringPeriods)
	if filt == nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrParsingFailed, PlantFloweringPeriodFilterParam, queryValue)
	}
	return filt, nil
}
