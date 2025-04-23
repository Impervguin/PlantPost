package plantfilters

import (
	"PlantSite/internal/models/plant"
	"PlantSite/internal/models/search"
	"fmt"
)

type PlantFilter interface {
	ToDomain() (search.PlantFilter, error)
	Bind(params map[string]interface{}) error
	Type() string
}

type PlantNameFilter struct {
	Name string `json:"name" form:"name" binding:"required"`
}

func (f *PlantNameFilter) ToDomain() (search.PlantFilter, error) {
	return search.NewPlantNameFilter(f.Name), nil
}

func (f *PlantNameFilter) Bind(params map[string]interface{}) error {
	name, ok := params["name"]
	if !ok {
		return fmt.Errorf("name not found in params")
	}
	nameStr, ok := name.(string)
	if !ok {
		return fmt.Errorf("name is not a string")
	}
	f.Name = nameStr
	return nil
}

func (f *PlantNameFilter) Type() string {
	return PlantNameFilterID
}

type PlantLatinNameFilter struct {
	LatinName string `json:"latin_name" form:"latin_name" binding:"required"`
}

func (f *PlantLatinNameFilter) ToDomain() (search.PlantFilter, error) {
	return search.NewPlantLatinNameFilter(f.LatinName), nil
}

func (f *PlantLatinNameFilter) Bind(params map[string]interface{}) error {
	latinName, ok := params["latin_name"]
	if !ok {
		return fmt.Errorf("latin_name not found in params")
	}
	latinNameStr, ok := latinName.(string)
	if !ok {
		return fmt.Errorf("latin_name is not a string")
	}
	f.LatinName = latinNameStr
	return nil
}

func (f *PlantLatinNameFilter) Type() string {
	return PlantLatinNameFilterID
}

type PlantCategoryFilter struct {
	Category string `json:"category" form:"category" binding:"required"`
}

func (f *PlantCategoryFilter) ToDomain() (search.PlantFilter, error) {
	return search.NewPlantCategoryFilter(f.Category), nil
}

func (f *PlantCategoryFilter) Bind(params map[string]interface{}) error {
	category, ok := params["category"]
	if !ok {
		return fmt.Errorf("category not found in params")
	}
	categoryStr, ok := category.(string)
	if !ok {
		return fmt.Errorf("category is not a string")
	}
	f.Category = categoryStr
	return nil
}

func (f *PlantCategoryFilter) Type() string {
	return PlantCategoryFilterID
}

type PlantHeightFilter struct {
	Min float64 `json:"min" form:"min" binding:"required"`
	Max float64 `json:"max" form:"max" binding:"required"`
}

func (f *PlantHeightFilter) ToDomain() (search.PlantFilter, error) {
	return search.NewPlantHeightFilter(f.Min, f.Max), nil
}

func (f *PlantHeightFilter) Bind(params map[string]interface{}) error {
	min, ok := params["min"]
	if !ok {
		return fmt.Errorf("min not found in params")
	}
	minFloat, ok := min.(float64)
	if !ok {
		return fmt.Errorf("min is not a float64")
	}
	f.Min = minFloat

	max, ok := params["max"]
	if !ok {
		return fmt.Errorf("max not found in params")
	}
	maxFloat, ok := max.(float64)
	if !ok {
		return fmt.Errorf("max is not a float64")
	}
	f.Max = maxFloat
	return nil
}

func (f *PlantHeightFilter) Type() string {
	return PlantHeightFilterID
}

type PlantDiameterFilter struct {
	Min float64 `json:"min" form:"min" binding:"required"`
	Max float64 `json:"max" form:"max" binding:"required"`
}

func (f *PlantDiameterFilter) ToDomain() (search.PlantFilter, error) {
	return search.NewPlantDiameterFilter(f.Min, f.Max), nil
}

func (f *PlantDiameterFilter) Bind(params map[string]interface{}) error {
	min, ok := params["min"]
	if !ok {
		return fmt.Errorf("min not found in params")
	}
	minFloat, ok := min.(float64)
	if !ok {
		return fmt.Errorf("min is not a float64")
	}
	f.Min = minFloat

	max, ok := params["max"]
	if !ok {
		return fmt.Errorf("max not found in params")
	}
	maxFloat, ok := max.(float64)
	if !ok {
		return fmt.Errorf("max is not a float64")
	}
	f.Max = maxFloat
	return nil
}

func (f *PlantDiameterFilter) Type() string {
	return PlantDiameterFilterID
}

type PlantSoilTypeFilter struct {
	PossibleSoilTypes []plant.Soil `json:"soil_types" form:"soil_types" binding:"required"`
}

func (f *PlantSoilTypeFilter) ToDomain() (search.PlantFilter, error) {
	return search.NewSoilTypeFilter(f.PossibleSoilTypes), nil
}

func (f *PlantSoilTypeFilter) Bind(params map[string]interface{}) error {
	possibleSoilTypes, ok := params["soil_types"]
	if !ok {
		return fmt.Errorf("soil_types not found in params")
	}
	possibleSoilTypesList, ok := possibleSoilTypes.([]interface{})
	if !ok {
		return fmt.Errorf("soil_types is not a list")
	}
	f.PossibleSoilTypes = make([]plant.Soil, 0, len(possibleSoilTypesList))
	for _, soilType := range possibleSoilTypesList {
		soilType, ok := soilType.(string)
		if !ok {
			return fmt.Errorf("soil_type is not a string")
		}
		f.PossibleSoilTypes = append(f.PossibleSoilTypes, plant.Soil(soilType))
	}
	return nil
}

func (f *PlantSoilTypeFilter) Type() string {
	return PlantSoilTypeFilterID
}

type PlantSoilAcidityFilter struct {
	Min plant.SoilAcidity `json:"min" form:"min" binding:"required"`
	Max plant.SoilAcidity `json:"max" form:"max" binding:"required"`
}

func (f *PlantSoilAcidityFilter) ToDomain() (search.PlantFilter, error) {
	return search.NewSoilAcidityFilter(f.Min, f.Max), nil
}

func (f *PlantSoilAcidityFilter) Bind(params map[string]interface{}) error {
	min, ok := params["min"]
	if !ok {
		return fmt.Errorf("min not found in params")
	}
	minInt, ok := min.(int)
	if !ok {
		return fmt.Errorf("min is not a int")
	}
	f.Min = plant.SoilAcidity(minInt)

	max, ok := params["max"]
	if !ok {
		return fmt.Errorf("max not found in params")
	}
	maxInt, ok := max.(int)
	if !ok {
		return fmt.Errorf("max is not a int")
	}
	f.Max = plant.SoilAcidity(maxInt)
	return nil
}

func (f *PlantSoilAcidityFilter) Type() string {
	return PlantSoilAcidityFilterID
}

type PlantSoilMoistureFilter struct {
	PossibleMoistures []plant.SoilMoisture `json:"moistures" form:"moistures" binding:"required"`
}

func (f *PlantSoilMoistureFilter) ToDomain() (search.PlantFilter, error) {
	return search.NewSoilMoistureFilter(f.PossibleMoistures), nil
}

func (f *PlantSoilMoistureFilter) Bind(params map[string]interface{}) error {
	possibleMoistures, ok := params["moistures"]
	if !ok {
		return fmt.Errorf("moistures not found in params")
	}
	possibleMoisturesList, ok := possibleMoistures.([]interface{})
	if !ok {
		return fmt.Errorf("moistures is not a list")
	}
	f.PossibleMoistures = make([]plant.SoilMoisture, 0, len(possibleMoisturesList))
	for _, moisture := range possibleMoisturesList {
		moisture, ok := moisture.(string)
		if !ok {
			return fmt.Errorf("moisture is not a string")
		}
		f.PossibleMoistures = append(f.PossibleMoistures, plant.SoilMoisture(moisture))
	}
	return nil
}

func (f *PlantSoilMoistureFilter) Type() string {
	return PlantSoilMoistureFilterID
}

type PlantLightRelationFilter struct {
	PossibleRelations []plant.LightRelation `json:"light_relations" form:"light_relations" binding:"required"`
}

func (f *PlantLightRelationFilter) ToDomain() (search.PlantFilter, error) {
	return search.NewLightRelationFilter(f.PossibleRelations), nil
}

func (f *PlantLightRelationFilter) Bind(params map[string]interface{}) error {
	possibleRelations, ok := params["light_relations"]
	if !ok {
		return fmt.Errorf("relations not found in params")
	}
	possibleRelationsList, ok := possibleRelations.([]interface{})
	if !ok {
		return fmt.Errorf("relations is not a list")
	}
	f.PossibleRelations = make([]plant.LightRelation, 0, len(possibleRelationsList))
	for _, relation := range possibleRelationsList {
		relation, ok := relation.(string)
		if !ok {
			return fmt.Errorf("relation is not a string")
		}
		f.PossibleRelations = append(f.PossibleRelations, plant.LightRelation(relation))
	}
	return nil
}

func (f *PlantLightRelationFilter) Type() string {
	return PlantLightRelationFilterID
}

type PlantWinterHardinessFilter struct {
	Min plant.WinterHardiness `json:"min" form:"min" binding:"required"`
	Max plant.WinterHardiness `json:"max" form:"max" binding:"required"`
}

func (f *PlantWinterHardinessFilter) ToDomain() (search.PlantFilter, error) {
	return search.NewWinterHardinessFilter(f.Min, f.Max), nil
}

func (f *PlantWinterHardinessFilter) Bind(params map[string]interface{}) error {
	min, ok := params["min"]
	if !ok {
		return fmt.Errorf("min not found in params")
	}
	minInt, ok := min.(int)
	if !ok {
		return fmt.Errorf("min is not a int")
	}
	f.Min = plant.WinterHardiness(minInt)

	max, ok := params["max"]
	if !ok {
		return fmt.Errorf("max not found in params")
	}
	maxInt, ok := max.(int)
	if !ok {
		return fmt.Errorf("max is not a int")
	}
	f.Max = plant.WinterHardiness(maxInt)
	return nil
}

func (f *PlantWinterHardinessFilter) Type() string {
	return PlantWinterHardinessFilterID
}

type PlantFloweringPeriodFilter struct {
	PossibleFloweringPeriods []plant.FloweringPeriod `json:"flowering_periods" form:"flowering_periods" binding:"required"`
}

func (f *PlantFloweringPeriodFilter) ToDomain() (search.PlantFilter, error) {
	return search.NewFloweringPeriodFilter(f.PossibleFloweringPeriods), nil
}

func (f *PlantFloweringPeriodFilter) Bind(params map[string]interface{}) error {
	possibleFloweringPeriods, ok := params["flowering_periods"]
	if !ok {
		return fmt.Errorf("flowering_periods not found in params")
	}
	possibleFloweringPeriodsList, ok := possibleFloweringPeriods.([]interface{})
	if !ok {
		return fmt.Errorf("flowering_periods is not a list")
	}
	f.PossibleFloweringPeriods = make([]plant.FloweringPeriod, 0, len(possibleFloweringPeriodsList))
	for _, floweringPeriod := range possibleFloweringPeriodsList {
		floweringPeriod, ok := floweringPeriod.(string)
		if !ok {
			return fmt.Errorf("flowering_period is not a string")
		}
		f.PossibleFloweringPeriods = append(f.PossibleFloweringPeriods, plant.FloweringPeriod(floweringPeriod))
	}
	return nil
}

func (f *PlantFloweringPeriodFilter) Type() string {
	return PlantFloweringPeriodFilterID
}
