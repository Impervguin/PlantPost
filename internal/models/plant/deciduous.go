package plant

import "fmt"

// func init() {
// 	globalPlantRegistry.Register("deciduous", NewDeciduousSpecification)
// }

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

// func NewDeciduousSpecification(data map[string]interface{}) (PlantSpecification, error) {
// 	heightM, ok := data["height_m"].(float64)
// 	if !ok {
// 		return nil, fmt.Errorf("height_m is not a float64")
// 	}

// 	diameterM, ok := data["diameter_m"].(float64)
// 	if !ok {
// 		return nil, fmt.Errorf("diameter_m is not a float64")
// 	}

// 	floweringPeriod, ok := data["flowering_period"].(FloweringPeriod)
// 	if !ok {
// 		return nil, fmt.Errorf("flowering_period is not a FloweringPeriod")
// 	}

// 	soilAcidity, ok := data["soil_acidity"].(int)
// 	if !ok {
// 		return nil, fmt.Errorf("soil_acidity is not an int")
// 	}

// 	soilMoisture, ok := data["soil_moisture"].(SoilMoisture)
// 	if !ok {
// 		return nil, fmt.Errorf("soil_moisture is not a SoilMoisture")
// 	}

// 	LightRelation, ok := data["light_sensivity"].(LightRelation)
// 	if !ok {
// 		return nil, fmt.Errorf("light_sensivity is not a LightRelation")
// 	}

// 	soilType, ok := data["soil_type"].(Soil)
// 	if !ok {
// 		return nil, fmt.Errorf("soil_type is not a Soil")
// 	}

// 	winterHardiness, ok := data["winter_hardiness"].(WinterHardiness)
// 	if !ok {
// 		return nil, fmt.Errorf("winter_hardiness is not a WinterHardiness")
// 	}

// 	spec := &DeciduousSpecification{
// 		HeightM:         heightM,
// 		DiameterM:       diameterM,
// 		FloweringPeriod: floweringPeriod,
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
	if d.soilAcidity <= 0 {
		return fmt.Errorf("soil_acidity should be greater than 0")
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
