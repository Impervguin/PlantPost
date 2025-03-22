package plant

import "fmt"

// func init() {
// 	globalPlantRegistry.Register("coniferous", NewConiferousSpecification)
// }

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

// func NewConiferousSpecification(data map[string]interface{}) (PlantSpecification, error) {
// 	heightM, ok := data["height_m"].(float64)
// 	if !ok {
// 		return nil, fmt.Errorf("height_m is not a float64")
// 	}

// 	diameterM, ok := data["diameter_m"].(float64)
// 	if !ok {
// 		return nil, fmt.Errorf("diameter_m is not a float64")
// 	}
// 	soilAcidity, ok := data["soil_acidity"].(int)
// 	if !ok {
// 		return nil, fmt.Errorf("soil_acidity is not an int")
// 	}

// 	soilMoisture, ok := data["soil_moisture"].(SoilMoisture)
// 	if !ok {
// 		return nil, fmt.Errorf("soil_moisture is not a SoilMoisture")
// 	}
// 	if err := soilMoisture.Validate(); err != nil {
// 		return nil, err
// 	}

// 	LightRelation, ok := data["light_sensivity"].(LightRelation)
// 	if !ok {
// 		return nil, fmt.Errorf("light_sensivity is not a LightRelation")
// 	}
// 	if err := LightRelation.Validate(); err != nil {
// 		return nil, err
// 	}

// 	soilType, ok := data["soil_type"].(Soil)
// 	if !ok {
// 		return nil, fmt.Errorf("soil_type is not a Soil")
// 	}

// 	winterHardiness, ok := data["winter_hardiness"].(WinterHardiness)
// 	if !ok {
// 		return nil, fmt.Errorf("winter_hardiness is not a WinterHardiness")
// 	}
// 	spec := &ConiferousSpecification{
// 		HeightM:         heightM,
// 		DiameterM:       diameterM,
// 		SoilAcidity:     soilAcidity,
// 		SoilMoisture:    soilMoisture,
// 		LightRelation:  LightRelation,
// 		SoilType:        soilType,
// 		WinterHardiness: winterHardiness,
// 	}
// 	if err := spec.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return spec, nil
// }

func NewConiferousSpecification(heightM, diameterM float64,
	soilAcidity SoilAcidity,
	soilMoisture SoilMoisture,
	LightRelation LightRelation,
	soilType Soil,
	winterHardiness WinterHardiness) (*ConiferousSpecification, error) {
	spec := &ConiferousSpecification{
		heightM:         heightM,
		diameterM:       diameterM,
		soilAcidity:     soilAcidity,
		soilMoisture:    soilMoisture,
		lightRelation:   LightRelation,
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
