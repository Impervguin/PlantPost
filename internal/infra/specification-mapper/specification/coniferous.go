package plantstorage

import (
	pgconsts "PlantSite/internal/infra/pg-consts"
	registry "PlantSite/internal/infra/specification-mapper/plant-registry"
	"PlantSite/internal/models/plant"
	"fmt"
	"math"
)

type ConiferousSpecification struct {
	HeightM         float64
	DiameterM       float64
	SoilAcidity     plant.SoilAcidity
	SoilMoisture    plant.SoilMoisture
	LightRelation   plant.LightRelation
	SoilType        plant.Soil
	WinterHardiness plant.WinterHardiness
}

func init() {
	registry.Register(plant.ConiferousCategory, ConiferousFromJsonB, ConiferousFromDomain)
}

var _ registry.PlantSpecification = &ConiferousSpecification{}

func (spec *ConiferousSpecification) ToJsonB() (registry.JsonB, error) {
	return map[string]interface{}{
		pgconsts.JsonBHeightMKey:         spec.HeightM,
		pgconsts.JsonBDiameterMKey:       spec.DiameterM,
		pgconsts.JsonBSoilAcidityKey:     spec.SoilAcidity,
		pgconsts.JsonBSoilMoistureKey:    spec.SoilMoisture,
		pgconsts.JsonBLightRelationKey:   spec.LightRelation,
		pgconsts.JsonBSoilTypeKey:        spec.SoilType,
		pgconsts.JsonBWinterHardinessKey: spec.WinterHardiness,
	}, nil
}

func (spec *ConiferousSpecification) ToDomain() (plant.PlantSpecification, error) {
	return plant.NewConiferousSpecification(
		spec.HeightM,
		spec.DiameterM,
		spec.SoilAcidity,
		spec.SoilMoisture,
		spec.LightRelation,
		spec.SoilType,
		spec.WinterHardiness,
	)
}

func ConiferousFromJsonB(JsonB registry.JsonB) (registry.PlantSpecification, error) {
	var conSpec ConiferousSpecification
	if val, ok := JsonB[pgconsts.JsonBHeightMKey]; ok {
		switch fact := val.(type) {
		case float64:
			conSpec.HeightM = fact
		case int:
			conSpec.HeightM = float64(fact)
		default:
			return nil, ErrJsonBFormatHeightM
		}
	} else {
		return nil, ErrJsonBMissingHeightM
	}

	if val, ok := JsonB[pgconsts.JsonBDiameterMKey]; ok {
		switch fact := val.(type) {
		case float64:
			conSpec.DiameterM = fact
		case int:
			conSpec.DiameterM = float64(fact)
		default:
			return nil, ErrJsonBFormatDiameterM
		}
	} else {
		return nil, ErrJsonBMissingDiameterM
	}

	if val, ok := JsonB[pgconsts.JsonBSoilAcidityKey]; ok {
		switch fact := val.(type) {
		case float64:
			conSpec.SoilAcidity = plant.SoilAcidity(math.Round(fact))
		case int:
			conSpec.SoilAcidity = plant.SoilAcidity(fact)
		default:
			return nil, ErrJsonBFormatSoilAcidity
		}
	} else {
		return nil, ErrJsonBMissingSoilAcidity
	}

	if val, ok := JsonB[pgconsts.JsonBSoilMoistureKey]; ok {
		switch fact := val.(type) {
		case string:
			conSpec.SoilMoisture = plant.SoilMoisture(fact)
		default:
			return nil, ErrJsonBFormatSoilMoisture
		}
	} else {
		return nil, ErrJsonBMissingSoilMoisture
	}

	if val, ok := JsonB[pgconsts.JsonBLightRelationKey]; ok {
		switch fact := val.(type) {
		case string:
			conSpec.LightRelation = plant.LightRelation(fact)
		default:
			return nil, ErrJsonBFormatLightRelation
		}
	} else {
		return nil, ErrJsonBMissingLightRelation
	}

	if val, ok := JsonB[pgconsts.JsonBSoilTypeKey]; ok {
		switch fact := val.(type) {
		case string:
			conSpec.SoilType = plant.Soil(fact)
		default:
			return nil, ErrJsonBFormatSoilType
		}
	} else {
		return nil, ErrJsonBMissingSoilType
	}

	if val, ok := JsonB[pgconsts.JsonBWinterHardinessKey]; ok {
		switch fact := val.(type) {
		case float64:
			conSpec.WinterHardiness = plant.WinterHardiness(math.Round(fact))
		case int:
			conSpec.WinterHardiness = plant.WinterHardiness(fact)
		default:
			return nil, ErrJsonBFormatWinterHardiness
		}
	} else {
		return nil, ErrJsonBMissingWinterHardiness
	}
	return &conSpec, nil
}

func ConiferousFromDomain(plSpec plant.PlantSpecification) (registry.PlantSpecification, error) {
	conSpec, ok := plSpec.(*plant.ConiferousSpecification)
	if !ok {
		return nil, fmt.Errorf("invalid plant specification type")
	}
	return &ConiferousSpecification{
		HeightM:         conSpec.GetHeightM(),
		DiameterM:       conSpec.GetDiameterM(),
		SoilAcidity:     conSpec.GetSoilAcidity(),
		SoilMoisture:    conSpec.GetSoilMoisture(),
		LightRelation:   conSpec.GetLightRelation(),
		SoilType:        conSpec.GetSoilType(),
		WinterHardiness: conSpec.GetWinterHardiness(),
	}, nil
}
