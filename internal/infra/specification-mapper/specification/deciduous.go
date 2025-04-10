package plantstorage

import (
	pgconsts "PlantSite/internal/infra/pg-consts"
	registry "PlantSite/internal/infra/specification-mapper/plant-registry"
	"PlantSite/internal/models/plant"
	"fmt"
	"math"
)

func init() {
	registry.Register(plant.DeciduousCategory, DeciduousFromJsonB, DeciduousFromDomain)
}

var _ registry.PlantSpecification = &DeciduousSpecification{}

type DeciduousSpecification struct {
	HeightM         float64
	DiameterM       float64
	FloweringPeriod plant.FloweringPeriod
	SoilAcidity     plant.SoilAcidity
	SoilMoisture    plant.SoilMoisture
	LightRelation   plant.LightRelation
	SoilType        plant.Soil
	WinterHardiness plant.WinterHardiness
}

func (spec *DeciduousSpecification) ToJsonB() (registry.JsonB, error) {
	return map[string]interface{}{
		pgconsts.JsonBHeightMKey:         spec.HeightM,
		pgconsts.JsonBDiameterMKey:       spec.DiameterM,
		pgconsts.JsonBFloweringPeriodKey: spec.FloweringPeriod,
		pgconsts.JsonBSoilAcidityKey:     spec.SoilAcidity,
		pgconsts.JsonBSoilMoistureKey:    spec.SoilMoisture,
		pgconsts.JsonBLightRelationKey:   spec.LightRelation,
		pgconsts.JsonBSoilTypeKey:        spec.SoilType,
		pgconsts.JsonBWinterHardinessKey: spec.WinterHardiness,
	}, nil
}

func (spec *DeciduousSpecification) ToDomain() (plant.PlantSpecification, error) {
	return plant.NewDeciduousSpecification(
		spec.HeightM,
		spec.DiameterM,
		spec.FloweringPeriod,
		spec.SoilAcidity,
		spec.SoilMoisture,
		spec.LightRelation,
		spec.SoilType,
		spec.WinterHardiness,
	)
}

func DeciduousFromJsonB(JsonB registry.JsonB) (registry.PlantSpecification, error) {
	var decSpec DeciduousSpecification
	if val, ok := JsonB[pgconsts.JsonBHeightMKey]; ok {
		switch fact := val.(type) {
		case float64:
			decSpec.HeightM = fact
		case int:
			decSpec.HeightM = float64(fact)
		default:
			return nil, ErrJsonBFormatHeightM
		}
	} else {
		return nil, ErrJsonBMissingHeightM
	}

	if val, ok := JsonB[pgconsts.JsonBDiameterMKey]; ok {
		switch fact := val.(type) {
		case float64:
			decSpec.DiameterM = fact
		case int:
			decSpec.DiameterM = float64(fact)
		default:
			return nil, ErrJsonBFormatDiameterM
		}
	} else {
		return nil, ErrJsonBMissingDiameterM
	}

	if val, ok := JsonB[pgconsts.JsonBSoilAcidityKey]; ok {
		switch fact := val.(type) {
		case float64:
			decSpec.SoilAcidity = plant.SoilAcidity(math.Round(fact))
		case int:
			decSpec.SoilAcidity = plant.SoilAcidity(fact)
		default:
			return nil, ErrJsonBFormatSoilAcidity
		}
	} else {
		return nil, ErrJsonBMissingSoilAcidity
	}

	if val, ok := JsonB[pgconsts.JsonBSoilMoistureKey]; ok {
		switch fact := val.(type) {
		case string:
			decSpec.SoilMoisture = plant.SoilMoisture(fact)
		default:
			return nil, ErrJsonBFormatSoilMoisture
		}
	} else {
		return nil, ErrJsonBMissingSoilMoisture
	}

	if val, ok := JsonB[pgconsts.JsonBLightRelationKey]; ok {
		switch fact := val.(type) {
		case string:
			decSpec.LightRelation = plant.LightRelation(fact)
		default:
			return nil, ErrJsonBFormatLightRelation
		}
	} else {
		return nil, ErrJsonBMissingLightRelation
	}

	if val, ok := JsonB[pgconsts.JsonBSoilTypeKey]; ok {
		switch fact := val.(type) {
		case string:
			decSpec.SoilType = plant.Soil(fact)
		default:
			return nil, ErrJsonBFormatSoilType
		}
	} else {
		return nil, ErrJsonBMissingSoilType
	}

	if val, ok := JsonB[pgconsts.JsonBWinterHardinessKey]; ok {
		switch fact := val.(type) {
		case float64:
			decSpec.WinterHardiness = plant.WinterHardiness(math.Round(fact))
		case int:
			decSpec.WinterHardiness = plant.WinterHardiness(fact)
		default:
			return nil, ErrJsonBFormatWinterHardiness
		}
	} else {
		return nil, ErrJsonBMissingWinterHardiness
	}
	if val, ok := JsonB[pgconsts.JsonBFloweringPeriodKey]; ok {
		switch fact := val.(type) {
		case string:
			decSpec.FloweringPeriod = plant.FloweringPeriod(fact)
		default:
			return nil, ErrJsonBFormatFloweringPeriod
		}
	} else {
		return nil, ErrJsonBMissingFloweringPeriod
	}
	return &decSpec, nil
}

func DeciduousFromDomain(plSpec plant.PlantSpecification) (registry.PlantSpecification, error) {
	decSpec, ok := plSpec.(*plant.DeciduousSpecification)
	if !ok {
		return nil, fmt.Errorf("invalid plant specification type")
	}
	return &DeciduousSpecification{
		HeightM:         decSpec.GetHeightM(),
		DiameterM:       decSpec.GetDiameterM(),
		FloweringPeriod: decSpec.GetFloweringPeriod(),
		SoilAcidity:     decSpec.GetSoilAcidity(),
		SoilMoisture:    decSpec.GetSoilMoisture(),
		LightRelation:   decSpec.GetLightRelation(),
		SoilType:        decSpec.GetSoilType(),
		WinterHardiness: decSpec.GetWinterHardiness(),
	}, nil
}
