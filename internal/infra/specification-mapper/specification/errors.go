package plantstorage

import (
	pgconsts "PlantSite/internal/infra/pg-consts"
	"fmt"
)

var (
	ErrJSONBMissingHeightM         = fmt.Errorf("missing %s value", pgconsts.JSONBHeightMKey)
	ErrJSONBMissingDiameterM       = fmt.Errorf("missing %s value", pgconsts.JSONBDiameterMKey)
	ErrJSONBMissingFloweringPeriod = fmt.Errorf("missing %s value", pgconsts.JSONBFloweringPeriodKey)
	ErrJSONBMissingSoilAcidity     = fmt.Errorf("missing %s value", pgconsts.JSONBSoilAcidityKey)
	ErrJSONBMissingSoilMoisture    = fmt.Errorf("missing %s value", pgconsts.JSONBSoilMoistureKey)
	ErrJSONBMissingLightRelation   = fmt.Errorf("missing %s value", pgconsts.JSONBLightRelationKey)
	ErrJSONBMissingSoilType        = fmt.Errorf("missing %s value", pgconsts.JSONBSoilTypeKey)
	ErrJSONBMissingWinterHardiness = fmt.Errorf("missing %s value", pgconsts.JSONBWinterHardinessKey)
)

var (
	ErrJSONBFormatHeightM         = fmt.Errorf("invalid %s value", pgconsts.JSONBHeightMKey)
	ErrJSONBFormatDiameterM       = fmt.Errorf("invalid %s value", pgconsts.JSONBDiameterMKey)
	ErrJSONBFormatFloweringPeriod = fmt.Errorf("invalid %s value", pgconsts.JSONBFloweringPeriodKey)
	ErrJSONBFormatSoilAcidity     = fmt.Errorf("invalid %s value", pgconsts.JSONBSoilAcidityKey)
	ErrJSONBFormatSoilMoisture    = fmt.Errorf("invalid %s value", pgconsts.JSONBSoilMoistureKey)
	ErrJSONBFormatLightRelation   = fmt.Errorf("invalid %s value", pgconsts.JSONBLightRelationKey)
	ErrJSONBFormatSoilType        = fmt.Errorf("invalid %s value", pgconsts.JSONBSoilTypeKey)
	ErrJSONBFormatWinterHardiness = fmt.Errorf("invalid %s value", pgconsts.JSONBWinterHardinessKey)
)
