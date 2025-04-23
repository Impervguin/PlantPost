package request

import (
	"PlantSite/internal/models/plant"
	"fmt"
)

type PlantSpecification interface {
	GetCategory() string
	ToDomain() (plant.PlantSpecification, error)
}

type ConiferousSpecification struct {
	HeightM   float64
	DiameterM float64

	SoilAcidity     int
	SoilMoisture    string
	LightRelation   string
	SoilType        string
	WinterHardiness int
}

func (c *ConiferousSpecification) GetCategory() string {
	return plant.ConiferousCategory
}

func (c *ConiferousSpecification) ToDomain() (plant.PlantSpecification, error) {
	spec, err := plant.NewConiferousSpecification(c.HeightM, c.DiameterM,
		plant.SoilAcidity(c.SoilAcidity),
		plant.SoilMoisture(c.SoilMoisture),
		plant.LightRelation(c.LightRelation),
		plant.Soil(c.SoilType),
		plant.WinterHardiness(c.WinterHardiness))
	if err != nil {
		return nil, fmt.Errorf("can't create coniferous specification: %w", err)
	}
	return spec, nil
}

type DeciduousSpecification struct {
	HeightM   float64
	DiameterM float64

	FloweringPeriod string
	SoilAcidity     int
	SoilMoisture    string
	LightRelation   string
	SoilType        string
	WinterHardiness int
}

func (d *DeciduousSpecification) GetCategory() string {
	return plant.DeciduousCategory
}

func (d *DeciduousSpecification) ToDomain() (plant.PlantSpecification, error) {
	spec, err := plant.NewDeciduousSpecification(d.HeightM, d.DiameterM,
		plant.FloweringPeriod(d.FloweringPeriod),
		plant.SoilAcidity(d.SoilAcidity),
		plant.SoilMoisture(d.SoilMoisture),
		plant.LightRelation(d.LightRelation),
		plant.Soil(d.SoilType),
		plant.WinterHardiness(d.WinterHardiness))
	if err != nil {
		return nil, fmt.Errorf("can't create deciduous specification: %w", err)
	}
	return spec, nil
}
