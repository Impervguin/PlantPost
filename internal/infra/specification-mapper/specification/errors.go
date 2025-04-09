package plantstorage

import (
	pgconsts "PlantSite/internal/infra/pg-consts"
	"fmt"
)

var (
	ErrJsonBMissingHeightM         = fmt.Errorf("missing %s value", pgconsts.JsonBHeightMKey)
	ErrJsonBMissingDiameterM       = fmt.Errorf("missing %s value", pgconsts.JsonBDiameterMKey)
	ErrJsonBMissingFloweringPeriod = fmt.Errorf("missing %s value", pgconsts.JsonBFloweringPeriodKey)
	ErrJsonBMissingSoilAcidity     = fmt.Errorf("missing %s value", pgconsts.JsonBSoilAcidityKey)
	ErrJsonBMissingSoilMoisture    = fmt.Errorf("missing %s value", pgconsts.JsonBSoilMoistureKey)
	ErrJsonBMissingLightRelation   = fmt.Errorf("missing %s value", pgconsts.JsonBLightRelationKey)
	ErrJsonBMissingSoilType        = fmt.Errorf("missing %s value", pgconsts.JsonBSoilTypeKey)
	ErrJsonBMissingWinterHardiness = fmt.Errorf("missing %s value", pgconsts.JsonBWinterHardinessKey)
)

var (
	ErrJsonBFormatHeightM         = fmt.Errorf("invalid %s value", pgconsts.JsonBHeightMKey)
	ErrJsonBFormatDiameterM       = fmt.Errorf("invalid %s value", pgconsts.JsonBDiameterMKey)
	ErrJsonBFormatFloweringPeriod = fmt.Errorf("invalid %s value", pgconsts.JsonBFloweringPeriodKey)
	ErrJsonBFormatSoilAcidity     = fmt.Errorf("invalid %s value", pgconsts.JsonBSoilAcidityKey)
	ErrJsonBFormatSoilMoisture    = fmt.Errorf("invalid %s value", pgconsts.JsonBSoilMoistureKey)
	ErrJsonBFormatLightRelation   = fmt.Errorf("invalid %s value", pgconsts.JsonBLightRelationKey)
	ErrJsonBFormatSoilType        = fmt.Errorf("invalid %s value", pgconsts.JsonBSoilTypeKey)
	ErrJsonBFormatWinterHardiness = fmt.Errorf("invalid %s value", pgconsts.JsonBWinterHardinessKey)
)
