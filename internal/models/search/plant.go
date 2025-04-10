package search

import (
	"PlantSite/internal/models/plant"

	"slices"
)

type PlantFilter interface {
	Filter(p *plant.Plant) bool
	Identifier() string
}

type PlantNameFilter struct {
	Name string
}

var _ PlantFilter = &PlantNameFilter{}

func NewPlantNameFilter(name string) *PlantNameFilter {
	return &PlantNameFilter{Name: name}
}

func (p *PlantNameFilter) Identifier() string {
	return PlantNameFilterID
}

func (p *PlantNameFilter) Filter(plant *plant.Plant) bool {
	return plant.GetName() == p.Name
}

type PlantCategoryFilter struct {
	Category string
}

var _ PlantFilter = &PlantCategoryFilter{}

func NewPlantCategoryFilter(category string) *PlantCategoryFilter {
	return &PlantCategoryFilter{Category: category}
}

func (p *PlantCategoryFilter) Filter(plant *plant.Plant) bool {
	return plant.GetCategory() == p.Category
}

func (p *PlantCategoryFilter) Identifier() string {
	return PlantCategoryFilterID
}

type PlantLatinNameFilter struct {
	LatinName string
}

var _ PlantFilter = &PlantLatinNameFilter{}

func NewPlantLatinNameFilter(latinName string) *PlantLatinNameFilter {
	return &PlantLatinNameFilter{LatinName: latinName}
}

func (p *PlantLatinNameFilter) Filter(plant *plant.Plant) bool {
	return plant.GetLatinName() == p.LatinName
}

func (p *PlantLatinNameFilter) Identifier() string {
	return PlantLatinNameFilterID
}

type PlantHeightFilter struct {
	Min, Max float64
}

var _ PlantFilter = &PlantHeightFilter{}

func NewPlantHeightFilter(min, max float64) *PlantHeightFilter {
	return &PlantHeightFilter{Min: min, Max: max}
}

func (p *PlantHeightFilter) Identifier() string {
	return PlantHeightFilterID
}

func (p *PlantHeightFilter) Filter(pl *plant.Plant) bool {
	spec := pl.GetSpecification()
	switch impl := spec.(type) {
	case *plant.ConiferousSpecification:
		return impl.GetHeightM() > p.Min && impl.GetHeightM() < p.Max
	case *plant.DeciduousSpecification:
		return impl.GetHeightM() > p.Min && impl.GetHeightM() < p.Max
	}
	return false
}

type PlantDiameterFilter struct {
	Min, Max float64
}

var _ PlantFilter = &PlantDiameterFilter{}

func (p *PlantDiameterFilter) Identifier() string {
	return PlantDiameterFilterID
}

func NewPlantDiameterFilter(min, max float64) *PlantDiameterFilter {
	return &PlantDiameterFilter{Min: min, Max: max}
}

func (p *PlantDiameterFilter) Filter(pl *plant.Plant) bool {
	spec := pl.GetSpecification()
	switch impl := spec.(type) {
	case *plant.ConiferousSpecification:
		return impl.GetDiameterM() > p.Min && impl.GetDiameterM() < p.Max
	case *plant.DeciduousSpecification:
		return impl.GetDiameterM() > p.Min && impl.GetDiameterM() < p.Max
	}
	return false
}

type PlantSoilAcidityFilter struct {
	Min, Max plant.SoilAcidity
}

var _ PlantFilter = &PlantSoilAcidityFilter{}

func (p *PlantSoilAcidityFilter) Identifier() string {
	return PlantSoilAcidityFilterID
}

func NewSoilAcidityFilter(min, max plant.SoilAcidity) *PlantSoilAcidityFilter {
	return &PlantSoilAcidityFilter{Min: min, Max: max}
}

func (p *PlantSoilAcidityFilter) Filter(pl *plant.Plant) bool {
	spec := pl.GetSpecification()
	switch impl := spec.(type) {
	case *plant.ConiferousSpecification:
		return impl.GetSoilAcidity() >= p.Min && impl.GetSoilAcidity() <= p.Max
	case *plant.DeciduousSpecification:
		return impl.GetSoilAcidity() >= p.Min && impl.GetSoilAcidity() <= p.Max
	}
	return false
}

type PlantSoilMoistureFilter struct {
	PossibleMoistures []plant.SoilMoisture
}

var _ PlantFilter = &PlantSoilMoistureFilter{}

func (p *PlantSoilMoistureFilter) Identifier() string {
	return PlantSoilMoistureFilterID
}

func NewSoilMoistureFilter(possibleMoistures []plant.SoilMoisture) *PlantSoilMoistureFilter {
	return &PlantSoilMoistureFilter{PossibleMoistures: possibleMoistures}
}

func (p *PlantSoilMoistureFilter) Filter(pl *plant.Plant) bool {
	spec := pl.GetSpecification()
	switch impl := spec.(type) {
	case *plant.ConiferousSpecification:
		return slices.Contains(p.PossibleMoistures, impl.GetSoilMoisture())
	case *plant.DeciduousSpecification:
		return slices.Contains(p.PossibleMoistures, impl.GetSoilMoisture())
	}
	return false
}

type PlantLightRelationFilter struct {
	PossibleRelations []plant.LightRelation
}

var _ PlantFilter = &PlantLightRelationFilter{}

func (p *PlantLightRelationFilter) Identifier() string {
	return PlantLightRelationFilterID
}

func NewLightRelationFilter(possibleRelations []plant.LightRelation) *PlantLightRelationFilter {
	return &PlantLightRelationFilter{PossibleRelations: possibleRelations}
}

func (p *PlantLightRelationFilter) Filter(pl *plant.Plant) bool {
	spec := pl.GetSpecification()
	switch impl := spec.(type) {
	case *plant.ConiferousSpecification:
		return slices.Contains(p.PossibleRelations, impl.GetLightRelation())
	case *plant.DeciduousSpecification:
		return slices.Contains(p.PossibleRelations, impl.GetLightRelation())
	}
	return false
}

type PlantHardinessFilter struct {
	Min, Max plant.WinterHardiness
}

var _ PlantFilter = &PlantHardinessFilter{}

func (f *PlantHardinessFilter) Identifier() string {
	return PlantWinterHardinessFilterID
}

func NewWinterHardinessFilter(min, max plant.WinterHardiness) *PlantHardinessFilter {
	return &PlantHardinessFilter{Min: min, Max: max}
}

func (p *PlantHardinessFilter) Filter(pl *plant.Plant) bool {
	spec := pl.GetSpecification()
	switch impl := spec.(type) {
	case *plant.ConiferousSpecification:
		return impl.GetWinterHardiness() > p.Min && impl.GetWinterHardiness() < p.Max
	case *plant.DeciduousSpecification:
		return impl.GetWinterHardiness() > p.Min && impl.GetWinterHardiness() < p.Max
	}
	return false
}

type PlantSoilTypeFilter struct {
	PossibleSoilTypes []plant.Soil
}

var _ PlantFilter = &PlantSoilTypeFilter{}

func (p *PlantSoilTypeFilter) Identifier() string {
	return PlantSoilTypeFilterID
}

func NewSoilTypeFilter(possibleSoilTypes []plant.Soil) *PlantSoilTypeFilter {
	return &PlantSoilTypeFilter{PossibleSoilTypes: possibleSoilTypes}
}

func (p *PlantSoilTypeFilter) Filter(pl *plant.Plant) bool {
	spec := pl.GetSpecification()
	switch impl := spec.(type) {
	case *plant.ConiferousSpecification:
		return slices.Contains(p.PossibleSoilTypes, impl.GetSoilType())
	case *plant.DeciduousSpecification:
		return slices.Contains(p.PossibleSoilTypes, impl.GetSoilType())
	}
	return false
}

type PlantFloweringPeriodFilter struct {
	PossibleFloweringPeriods []plant.FloweringPeriod
}

var _ PlantFilter = &PlantFloweringPeriodFilter{}

func (p *PlantFloweringPeriodFilter) Identifier() string {
	return PlantFloweringPeriodFilterID
}

func NewFloweringPeriodFilter(arr []plant.FloweringPeriod) *PlantFloweringPeriodFilter {
	return &PlantFloweringPeriodFilter{
		PossibleFloweringPeriods: arr,
	}
}

func (p *PlantFloweringPeriodFilter) Filter(pl *plant.Plant) bool {
	spec := pl.GetSpecification()
	switch impl := spec.(type) {
	case *plant.DeciduousSpecification:
		return slices.Contains(p.PossibleFloweringPeriods, impl.GetFloweringPeriod())
	}
	return false
}
