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
	registry.Register(plant.ConiferousCategory, ConiferousFromJSONB, ConiferousFromDomain)
}

var _ registry.PlantSpecification = &ConiferousSpecification{}

var coniferousSpecDecoders = map[string]func(v any, conSpec *ConiferousSpecification) error{
	pgconsts.JSONBHeightMKey: func(v any, conSpec *ConiferousSpecification) error {
		switch fact := v.(type) {
		case float64:
			conSpec.HeightM = fact
		case int:
			conSpec.HeightM = float64(fact)
		default:
			return ErrJSONBFormatHeightM
		}
		return nil
	},
	pgconsts.JSONBDiameterMKey: func(v any, conSpec *ConiferousSpecification) error {
		switch fact := v.(type) {
		case float64:
			conSpec.DiameterM = fact
		case int:
			conSpec.DiameterM = float64(fact)
		default:
			return ErrJSONBFormatDiameterM
		}
		return nil
	},
	pgconsts.JSONBSoilAcidityKey: func(v any, conSpec *ConiferousSpecification) error {
		switch fact := v.(type) {
		case float64:
			conSpec.SoilAcidity = plant.SoilAcidity(math.Round(fact))
		case int:
			conSpec.SoilAcidity = plant.SoilAcidity(fact)
		default:
			return ErrJSONBFormatSoilAcidity
		}
		return nil
	},
	pgconsts.JSONBSoilMoistureKey: func(v any, conSpec *ConiferousSpecification) error {
		switch fact := v.(type) {
		case string:
			conSpec.SoilMoisture = plant.SoilMoisture(fact)
		default:
			return ErrJSONBFormatSoilMoisture
		}
		return nil
	},
	pgconsts.JSONBLightRelationKey: func(v any, conSpec *ConiferousSpecification) error {
		switch fact := v.(type) {
		case string:
			conSpec.LightRelation = plant.LightRelation(fact)
		default:
			return ErrJSONBFormatLightRelation
		}
		return nil
	},
	pgconsts.JSONBSoilTypeKey: func(v any, conSpec *ConiferousSpecification) error {
		switch fact := v.(type) {
		case string:
			conSpec.SoilType = plant.Soil(fact)
		default:
			return ErrJSONBFormatSoilType
		}
		return nil
	},
	pgconsts.JSONBWinterHardinessKey: func(v any, conSpec *ConiferousSpecification) error {
		switch fact := v.(type) {
		case float64:
			conSpec.WinterHardiness = plant.WinterHardiness(math.Round(fact))
		case int:
			conSpec.WinterHardiness = plant.WinterHardiness(fact)
		default:
			return ErrJSONBFormatWinterHardiness
		}
		return nil
	},
}

func (spec *ConiferousSpecification) ToJSONB() (registry.JSONB, error) {
	return map[string]interface{}{
		pgconsts.JSONBHeightMKey:         spec.HeightM,
		pgconsts.JSONBDiameterMKey:       spec.DiameterM,
		pgconsts.JSONBSoilAcidityKey:     spec.SoilAcidity,
		pgconsts.JSONBSoilMoistureKey:    spec.SoilMoisture,
		pgconsts.JSONBLightRelationKey:   spec.LightRelation,
		pgconsts.JSONBSoilTypeKey:        spec.SoilType,
		pgconsts.JSONBWinterHardinessKey: spec.WinterHardiness,
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

func ConiferousFromJSONB(jsonB registry.JSONB) (registry.PlantSpecification, error) {
	var conSpec ConiferousSpecification
	var mustBeKeys = map[string]error{
		pgconsts.JSONBHeightMKey:         ErrJSONBMissingHeightM,
		pgconsts.JSONBDiameterMKey:       ErrJSONBMissingDiameterM,
		pgconsts.JSONBSoilAcidityKey:     ErrJSONBMissingSoilAcidity,
		pgconsts.JSONBSoilMoistureKey:    ErrJSONBMissingSoilMoisture,
		pgconsts.JSONBLightRelationKey:   ErrJSONBMissingLightRelation,
		pgconsts.JSONBSoilTypeKey:        ErrJSONBMissingSoilType,
		pgconsts.JSONBWinterHardinessKey: ErrJSONBMissingWinterHardiness,
	}

	err := checkJSONBKeys(jsonB, mustBeKeys)
	if err != nil {
		return nil, err
	}

	for key, decoder := range coniferousSpecDecoders {
		if err := decoder(jsonB[key], &conSpec); err != nil {
			return nil, err
		}
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
