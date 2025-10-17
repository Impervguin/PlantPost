package search

import (
	"PlantSite/internal/models/album"
	"PlantSite/internal/models/plant"
	"strings"

	"slices"

	"github.com/google/uuid"
)

type PlantFilter interface {
	Filter(p *plant.Plant) bool
	Identifier() string
}

type ExactPlantNameFilter struct {
	Name string
}

type PlantNameFilter struct {
	Name string
}

var _ PlantFilter = &ExactPlantNameFilter{}

func NewExactPlantNameFilter(name string) *ExactPlantNameFilter {
	return &ExactPlantNameFilter{Name: name}
}

func (p *ExactPlantNameFilter) Identifier() string {
	return ExactPlantNameFilterID
}

func (p *ExactPlantNameFilter) Filter(plant *plant.Plant) bool {
	return plant.GetName() == p.Name
}

type PlantIDsFilter struct {
	IDs []uuid.UUID
}

var _ PlantFilter = &PlantIDsFilter{}

func NewPlantIDsFilter(ids []uuid.UUID) *PlantIDsFilter {
	return &PlantIDsFilter{IDs: ids}
}

func (p *PlantIDsFilter) Identifier() string {
	return PlantIDsFilterID
}

func (p *PlantIDsFilter) Filter(plant *plant.Plant) bool {
	return slices.Contains(p.IDs, plant.ID())
}

var _ PlantFilter = &PlantNameFilter{}

func NewPlantNameFilter(name string) *PlantNameFilter {
	return &PlantNameFilter{Name: name}
}

func (p *PlantNameFilter) Identifier() string {
	return PlantNameFilterID
}

func (p *PlantNameFilter) Filter(plant *plant.Plant) bool {
	return strings.Contains(plant.GetName(), p.Name)
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

func NewPlantHeightFilter(minVal, maxVal float64) *PlantHeightFilter {
	return &PlantHeightFilter{Min: minVal, Max: maxVal}
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

func NewPlantDiameterFilter(minVal, maxVal float64) *PlantDiameterFilter {
	return &PlantDiameterFilter{Min: minVal, Max: maxVal}
}

func (p *PlantDiameterFilter) Identifier() string {
	return PlantDiameterFilterID
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

func NewSoilAcidityFilter(minVal, maxVal plant.SoilAcidity) *PlantSoilAcidityFilter {
	return &PlantSoilAcidityFilter{Min: minVal, Max: maxVal}
}

func (p *PlantSoilAcidityFilter) Identifier() string {
	return PlantSoilAcidityFilterID
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

func NewSoilMoistureFilter(possibleMoistures []plant.SoilMoisture) *PlantSoilMoistureFilter {
	return &PlantSoilMoistureFilter{PossibleMoistures: possibleMoistures}
}

func (p *PlantSoilMoistureFilter) Identifier() string {
	return PlantSoilMoistureFilterID
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

func NewLightRelationFilter(possibleRelations []plant.LightRelation) *PlantLightRelationFilter {
	return &PlantLightRelationFilter{PossibleRelations: possibleRelations}
}

func (p *PlantLightRelationFilter) Identifier() string {
	return PlantLightRelationFilterID
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

func NewWinterHardinessFilter(minVal, maxVal plant.WinterHardiness) *PlantHardinessFilter {
	return &PlantHardinessFilter{Min: minVal, Max: maxVal}
}

func (f *PlantHardinessFilter) Identifier() string {
	return PlantWinterHardinessFilterID
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

func NewSoilTypeFilter(possibleSoilTypes []plant.Soil) *PlantSoilTypeFilter {
	return &PlantSoilTypeFilter{PossibleSoilTypes: possibleSoilTypes}
}

func (p *PlantSoilTypeFilter) Identifier() string {
	return PlantSoilTypeFilterID
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

func NewFloweringPeriodFilter(arr []plant.FloweringPeriod) *PlantFloweringPeriodFilter {
	return &PlantFloweringPeriodFilter{
		PossibleFloweringPeriods: arr,
	}
}

func (p *PlantFloweringPeriodFilter) Identifier() string {
	return PlantFloweringPeriodFilterID
}

func (p *PlantFloweringPeriodFilter) Filter(pl *plant.Plant) bool {
	spec := pl.GetSpecification()
	impl, ok := spec.(*plant.DeciduousSpecification)
	if !ok {
		return false
	}
	return slices.Contains(p.PossibleFloweringPeriods, impl.GetFloweringPeriod())
}

type PlantAlbumFilter struct {
	AlbumID uuid.UUID
	Albums  []*album.Album
}

func NewPlantAlbumFilter(albmID uuid.UUID, albms []*album.Album) *PlantAlbumFilter {
	if albms == nil {
		albms = make([]*album.Album, 0)
	}
	return &PlantAlbumFilter{AlbumID: albmID, Albums: albms}
}

var _ PlantFilter = &PlantAlbumFilter{}

func (p *PlantAlbumFilter) Identifier() string {
	return PlantAlbumFilterID
}

func (p *PlantAlbumFilter) Filter(pl *plant.Plant) bool {
	for _, albm := range p.Albums {
		if albm.ID() == p.AlbumID && slices.Contains(albm.PlantIDs(), pl.ID()) {
			return true
		}
	}
	return false
}
