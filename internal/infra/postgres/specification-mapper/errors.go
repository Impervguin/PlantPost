package plantstorage

import (
	pgconsts "PlantSite/internal/infra/postgres/pg-consts"
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
