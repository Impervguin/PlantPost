package plantstorage

import (
	pgconsts "PlantSite/internal/infra/postgres/pg-consts"
	"PlantSite/internal/models/plant"
	"fmt"
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
	Register(plant.DeciduousCategory, DeciduousFromJsonB, DeciduousFromDomain)
}

func (spec *ConiferousSpecification) ToJsonB() (JsonB, error) {
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

func ConiferousFromJsonB(JsonB JsonB) (PlantSpecification, error) {
	var conSpec ConiferousSpecification
	if val, ok := JsonB[pgconsts.JsonBHeightMKey]; ok {
		conSpec.HeightM = val.(float64)
	} else {
		return nil, ErrJsonBMissingHeightM
	}

	if val, ok := JsonB[pgconsts.JsonBDiameterMKey]; ok {
		conSpec.DiameterM = val.(float64)
	} else {
		return nil, ErrJsonBMissingDiameterM
	}

	if val, ok := JsonB[pgconsts.JsonBSoilAcidityKey]; ok {
		conSpec.SoilAcidity = val.(plant.SoilAcidity)
	} else {
		return nil, ErrJsonBMissingSoilAcidity
	}

	if val, ok := JsonB[pgconsts.JsonBSoilMoistureKey]; ok {
		conSpec.SoilMoisture = val.(plant.SoilMoisture)
	} else {
		return nil, ErrJsonBMissingSoilMoisture
	}

	if val, ok := JsonB[pgconsts.JsonBLightRelationKey]; ok {
		conSpec.LightRelation = val.(plant.LightRelation)
	} else {
		return nil, ErrJsonBMissingLightRelation
	}

	if val, ok := JsonB[pgconsts.JsonBSoilTypeKey]; ok {
		conSpec.SoilType = val.(plant.Soil)
	} else {
		return nil, ErrJsonBMissingSoilType
	}

	if val, ok := JsonB[pgconsts.JsonBWinterHardinessKey]; ok {
		conSpec.WinterHardiness = val.(plant.WinterHardiness)
	} else {
		return nil, ErrJsonBMissingWinterHardiness
	}
	return &conSpec, nil
}

func ConiferousFromDomain(plSpec plant.PlantSpecification) (PlantSpecification, error) {
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
