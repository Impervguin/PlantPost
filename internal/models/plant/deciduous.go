package plant

import "fmt"

// Лиственные
const DeciduousCategory = "deciduous"

type DeciduousSpecification struct {
	heightM   float64
	diameterM float64

	floweringPeriod FloweringPeriod
	soilAcidity     SoilAcidity
	soilMoisture    SoilMoisture
	lightRelation   LightRelation
	soilType        Soil
	winterHardiness WinterHardiness
}

func NewDeciduousSpecification(heightM, diameterM float64,
	floweringPeriod FloweringPeriod,
	soilAcidity SoilAcidity,
	soilMoisture SoilMoisture,
	lightRelation LightRelation,
	soilType Soil,
	winterHardiness WinterHardiness) (*DeciduousSpecification, error) {
	spec := &DeciduousSpecification{
		heightM:         heightM,
		diameterM:       diameterM,
		floweringPeriod: floweringPeriod,
		soilAcidity:     soilAcidity,
		soilMoisture:    soilMoisture,
		lightRelation:   lightRelation,
		soilType:        soilType,
		winterHardiness: winterHardiness,
	}
	if err := spec.Validate(); err != nil {
		return nil, err
	}
	return spec, nil
}

func (d *DeciduousSpecification) Validate() error {
	if d.heightM <= 0 || d.diameterM <= 0 {
		return fmt.Errorf("height_m and diameter_m should be greater than 0")
	}
	if err := d.soilAcidity.Validate(); err != nil {
		return err
	}
	if err := d.floweringPeriod.Validate(); err != nil {
		return err
	}
	if err := d.soilMoisture.Validate(); err != nil {
		return err
	}
	if err := d.lightRelation.Validate(); err != nil {
		return err
	}
	if err := d.soilType.Validate(); err != nil {
		return err
	}
	if err := d.winterHardiness.Validate(); err != nil {
		return err
	}
	return nil
}

func (d DeciduousSpecification) GetHeightM() float64 {
	return d.heightM
}

func (d DeciduousSpecification) GetDiameterM() float64 {
	return d.diameterM
}

func (d DeciduousSpecification) GetFloweringPeriod() FloweringPeriod {
	return d.floweringPeriod
}

func (d DeciduousSpecification) GetSoilAcidity() SoilAcidity {
	return d.soilAcidity
}

func (d DeciduousSpecification) GetSoilMoisture() SoilMoisture {
	return d.soilMoisture
}

func (d DeciduousSpecification) GetLightRelation() LightRelation {
	return d.lightRelation
}

func (d DeciduousSpecification) GetSoilType() Soil {
	return d.soilType
}

func (d DeciduousSpecification) GetWinterHardiness() WinterHardiness {
	return d.winterHardiness
}
