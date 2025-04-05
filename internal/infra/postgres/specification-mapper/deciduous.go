package plantstorage

import (
	pgconsts "PlantSite/internal/infra/postgres/pg-consts"
	"PlantSite/internal/models/plant"
	"fmt"
)

func init() {
	Register(plant.DeciduousCategory, DeciduousFromJsonB, DeciduousFromDomain)
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

func (spec *DeciduousSpecification) ToJsonB() (JsonB, error) {
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

func DeciduousFromJsonB(JsonB JsonB) (PlantSpecification, error) {
	var decSpec DeciduousSpecification
	if val, ok := JsonB[pgconsts.JsonBHeightMKey]; ok {
		decSpec.HeightM = val.(float64)
	} else {
		return nil, ErrJsonBMissingHeightM
	}

	if val, ok := JsonB[pgconsts.JsonBDiameterMKey]; ok {
		decSpec.DiameterM = val.(float64)
	} else {
		return nil, ErrJsonBMissingDiameterM
	}

	if val, ok := JsonB[pgconsts.JsonBFloweringPeriodKey]; ok {
		decSpec.FloweringPeriod = val.(plant.FloweringPeriod)
	} else {
		return nil, ErrJsonBMissingFloweringPeriod
	}

	if val, ok := JsonB[pgconsts.JsonBSoilAcidityKey]; ok {
		decSpec.SoilAcidity = val.(plant.SoilAcidity)
	} else {
		return nil, ErrJsonBMissingSoilAcidity
	}

	if val, ok := JsonB[pgconsts.JsonBSoilMoistureKey]; ok {
		decSpec.SoilMoisture = val.(plant.SoilMoisture)
	} else {
		return nil, ErrJsonBMissingSoilMoisture
	}

	if val, ok := JsonB[pgconsts.JsonBLightRelationKey]; ok {
		decSpec.LightRelation = val.(plant.LightRelation)
	} else {
		return nil, ErrJsonBMissingLightRelation
	}

	if val, ok := JsonB[pgconsts.JsonBSoilTypeKey]; ok {
		decSpec.SoilType = val.(plant.Soil)
	} else {
		return nil, ErrJsonBMissingSoilType
	}

	if val, ok := JsonB[pgconsts.JsonBWinterHardinessKey]; ok {
		decSpec.WinterHardiness = val.(plant.WinterHardiness)
	} else {
		return nil, ErrJsonBMissingWinterHardiness
	}
	return &decSpec, nil
}

func DeciduousFromDomain(plSpec plant.PlantSpecification) (PlantSpecification, error) {
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
