package plant

import "fmt"

// Хвойные
const ConiferousCategory = "coniferous"

type ConiferousSpecification struct {
	heightM   float64
	diameterM float64

	soilAcidity     SoilAcidity
	soilMoisture    SoilMoisture
	lightRelation   LightRelation
	soilType        Soil
	winterHardiness WinterHardiness
}

func NewConiferousSpecification(heightM, diameterM float64,
	soilAcidity SoilAcidity,
	soilMoisture SoilMoisture,
	lightRelation LightRelation,
	soilType Soil,
	winterHardiness WinterHardiness) (*ConiferousSpecification, error) {
	spec := &ConiferousSpecification{
		heightM:         heightM,
		diameterM:       diameterM,
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

func (c *ConiferousSpecification) Validate() error {
	if c.heightM <= 0 || c.diameterM <= 0 {
		return fmt.Errorf("height_m and diameter_m should be greater than 0")
	}

	if c.soilAcidity <= 0 {
		return fmt.Errorf("soil_acidity should be greater than 0")
	}

	if err := c.soilMoisture.Validate(); err != nil {
		return err
	}

	if err := c.lightRelation.Validate(); err != nil {
		return err
	}

	if err := c.soilType.Validate(); err != nil {
		return err
	}

	if err := c.winterHardiness.Validate(); err != nil {
		return err
	}

	return nil
}

func (c *ConiferousSpecification) GetHeightM() float64 {
	return c.heightM
}

func (c *ConiferousSpecification) GetDiameterM() float64 {
	return c.diameterM
}

func (c *ConiferousSpecification) GetSoilAcidity() SoilAcidity {
	return c.soilAcidity
}

func (c *ConiferousSpecification) GetSoilMoisture() SoilMoisture {
	return c.soilMoisture
}

func (c *ConiferousSpecification) GetLightRelation() LightRelation {
	return c.lightRelation
}

func (c *ConiferousSpecification) GetSoilType() Soil {
	return c.soilType
}

func (c *ConiferousSpecification) GetWinterHardiness() WinterHardiness {
	return c.winterHardiness
}
