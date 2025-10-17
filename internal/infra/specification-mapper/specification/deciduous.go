package plantstorage

import (
	pgconsts "PlantSite/internal/infra/pg-consts"
	registry "PlantSite/internal/infra/specification-mapper/plant-registry"
	"PlantSite/internal/models/plant"
	"fmt"
	"math"
)

func init() {
	registry.Register(plant.DeciduousCategory, DeciduousFromJSONB, DeciduousFromDomain)
}

var _ registry.PlantSpecification = &DeciduousSpecification{}

var decidousSpecDecoders = map[string]func(v any, decSpec *DeciduousSpecification) error{
	pgconsts.JSONBHeightMKey: func(v any, decSpec *DeciduousSpecification) error {
		switch fact := v.(type) {
		case float64:
			decSpec.HeightM = fact
		case int:
			decSpec.HeightM = float64(fact)
		default:
			return ErrJSONBFormatHeightM
		}
		return nil
	},
	pgconsts.JSONBDiameterMKey: func(v any, decSpec *DeciduousSpecification) error {
		switch fact := v.(type) {
		case float64:
			decSpec.DiameterM = fact
		case int:
			decSpec.DiameterM = float64(fact)
		default:
			return ErrJSONBFormatDiameterM
		}
		return nil
	},
	pgconsts.JSONBFloweringPeriodKey: func(v any, decSpec *DeciduousSpecification) error {
		switch fact := v.(type) {
		case string:
			decSpec.FloweringPeriod = plant.FloweringPeriod(fact)
		default:
			return ErrJSONBFormatFloweringPeriod
		}
		return nil
	},
	pgconsts.JSONBSoilAcidityKey: func(v any, decSpec *DeciduousSpecification) error {
		switch fact := v.(type) {
		case float64:
			decSpec.SoilAcidity = plant.SoilAcidity(math.Round(fact))
		case int:
			decSpec.SoilAcidity = plant.SoilAcidity(fact)
		default:
			return ErrJSONBFormatSoilAcidity
		}
		return nil
	},
	pgconsts.JSONBSoilMoistureKey: func(v any, decSpec *DeciduousSpecification) error {
		switch fact := v.(type) {
		case string:
			decSpec.SoilMoisture = plant.SoilMoisture(fact)
		default:
			return ErrJSONBFormatSoilMoisture
		}
		return nil
	},
	pgconsts.JSONBLightRelationKey: func(v any, decSpec *DeciduousSpecification) error {
		switch fact := v.(type) {
		case string:
			decSpec.LightRelation = plant.LightRelation(fact)
		default:
			return ErrJSONBFormatLightRelation
		}
		return nil
	},
	pgconsts.JSONBSoilTypeKey: func(v any, decSpec *DeciduousSpecification) error {
		switch fact := v.(type) {
		case string:
			decSpec.SoilType = plant.Soil(fact)
		default:
			return ErrJSONBFormatSoilType
		}
		return nil
	},
	pgconsts.JSONBWinterHardinessKey: func(v any, decSpec *DeciduousSpecification) error {
		switch fact := v.(type) {
		case float64:
			decSpec.WinterHardiness = plant.WinterHardiness(math.Round(fact))
		case int:
			decSpec.WinterHardiness = plant.WinterHardiness(fact)
		default:
			return ErrJSONBFormatWinterHardiness
		}
		return nil
	},
}

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

func (spec *DeciduousSpecification) ToJSONB() (registry.JSONB, error) {
	return map[string]interface{}{
		pgconsts.JSONBHeightMKey:         spec.HeightM,
		pgconsts.JSONBDiameterMKey:       spec.DiameterM,
		pgconsts.JSONBFloweringPeriodKey: spec.FloweringPeriod,
		pgconsts.JSONBSoilAcidityKey:     spec.SoilAcidity,
		pgconsts.JSONBSoilMoistureKey:    spec.SoilMoisture,
		pgconsts.JSONBLightRelationKey:   spec.LightRelation,
		pgconsts.JSONBSoilTypeKey:        spec.SoilType,
		pgconsts.JSONBWinterHardinessKey: spec.WinterHardiness,
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

func DeciduousFromJSONB(jsonB registry.JSONB) (registry.PlantSpecification, error) {
	var decSpec DeciduousSpecification
	var mustBeKeys = map[string]error{
		pgconsts.JSONBHeightMKey:         ErrJSONBMissingHeightM,
		pgconsts.JSONBDiameterMKey:       ErrJSONBMissingDiameterM,
		pgconsts.JSONBFloweringPeriodKey: ErrJSONBMissingFloweringPeriod,
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

	for key, decoder := range decidousSpecDecoders {
		if err := decoder(jsonB[key], &decSpec); err != nil {
			return nil, err
		}
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
