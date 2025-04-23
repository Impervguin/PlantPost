package mapper

import (
	"PlantSite/internal/models/plant"
	"errors"
	"fmt"
)

var ErrInvalidCategory = errors.New("invalid category (must be coniferous or deciduous)")

type ConiferousSpecification struct {
	HeightM         float64 `json:"height_m" form:"height_m" binding:"required"`
	DiameterM       float64 `json:"diameter_m" form:"diameter_m" binding:"required"`
	SoilAcidity     int     `json:"soil_acidity" form:"soil_acidity" binding:"required"`
	SoilMoisture    string  `json:"soil_moisture" form:"soil_moisture" binding:"required"`
	LightRelation   string  `json:"light_relation" form:"light_relation" binding:"required"`
	SoilType        string  `json:"soil_type" form:"soil_type" binding:"required"`
	WinterHardiness int     `json:"winter_hardiness" form:"winter_hardiness" binding:"required"`
}

func (c *ConiferousSpecification) Category() string {
	return plant.ConiferousCategory
}

func (c *ConiferousSpecification) ToDomain() (plant.PlantSpecification, error) {
	return plant.NewConiferousSpecification(c.HeightM, c.DiameterM,
		plant.SoilAcidity(c.SoilAcidity),
		plant.SoilMoisture(c.SoilMoisture),
		plant.LightRelation(c.LightRelation),
		plant.Soil(c.SoilType),
		plant.WinterHardiness(c.WinterHardiness))
}

type DeciduousSpecification struct {
	HeightM   float64 `json:"height_m" form:"height_m" binding:"required"`
	DiameterM float64 `json:"diameter_m" form:"diameter_m" binding:"required"`

	FloweringPeriod string `json:"flowering_period" form:"flowering_period" binding:"required"`
	SoilAcidity     int    `json:"soil_acidity" form:"soil_acidity" binding:"required"`
	SoilMoisture    string `json:"soil_moisture" form:"soil_moisture" binding:"required"`
	LightRelation   string `json:"light_relation" form:"light_relation" binding:"required"`
	SoilType        string `json:"soil_type" form:"soil_type" binding:"required"`
	WinterHardiness int    `json:"winter_hardiness" form:"winter_hardiness" binding:"required"`
}

func (d *DeciduousSpecification) Category() string {
	return plant.DeciduousCategory
}

func (d *DeciduousSpecification) ToDomain() (plant.PlantSpecification, error) {
	return plant.NewDeciduousSpecification(d.HeightM, d.DiameterM,
		plant.FloweringPeriod(d.FloweringPeriod),
		plant.SoilAcidity(d.SoilAcidity),
		plant.SoilMoisture(d.SoilMoisture),
		plant.LightRelation(d.LightRelation),
		plant.Soil(d.SoilType),
		plant.WinterHardiness(d.WinterHardiness))
}

type PlantSpecification interface {
	Category() string
	ToDomain() (plant.PlantSpecification, error)
}

func MapConiferousSpecification(specification plant.PlantSpecification) (*ConiferousSpecification, error) {
	spec, ok := specification.(*plant.ConiferousSpecification)
	if !ok {
		return nil, fmt.Errorf("invalid specification type: %T", specification)
	}
	return &ConiferousSpecification{
		HeightM:         spec.GetHeightM(),
		DiameterM:       spec.GetDiameterM(),
		SoilAcidity:     int(spec.GetSoilAcidity()),
		SoilMoisture:    string(spec.GetSoilMoisture()),
		LightRelation:   string(spec.GetLightRelation()),
		SoilType:        string(spec.GetSoilType()),
		WinterHardiness: int(spec.GetWinterHardiness()),
	}, nil
}

func MapDeciduousSpecification(specification plant.PlantSpecification) (*DeciduousSpecification, error) {
	spec, ok := specification.(*plant.DeciduousSpecification)
	if !ok {
		return nil, fmt.Errorf("invalid specification type: %T", specification)
	}
	return &DeciduousSpecification{
		HeightM:         spec.GetHeightM(),
		DiameterM:       spec.GetDiameterM(),
		FloweringPeriod: string(spec.GetFloweringPeriod()),
		SoilAcidity:     int(spec.GetSoilAcidity()),
		SoilMoisture:    string(spec.GetSoilMoisture()),
		LightRelation:   string(spec.GetLightRelation()),
		SoilType:        string(spec.GetSoilType()),
		WinterHardiness: int(spec.GetWinterHardiness()),
	}, nil
}

func MapSpecification(specification plant.PlantSpecification) (PlantSpecification, error) {
	switch specification.Category() {
	case plant.ConiferousCategory:
		return MapConiferousSpecification(specification)
	case plant.DeciduousCategory:
		return MapDeciduousSpecification(specification)
	default:
		return nil, ErrInvalidCategory
	}
}
