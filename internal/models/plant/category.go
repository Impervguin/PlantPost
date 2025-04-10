package plant

import (
	"github.com/google/uuid"
)

type ParamType string

const (
	ParameterTypeNumber ParamType = "number"
	ParameterTypeFloat  ParamType = "float"
	ParameterTypeString ParamType = "string"
	ParameterTypeSelect ParamType = "select"
)

type PlantParam struct {
	Name    string
	Type    ParamType
	Options []string
}

type PlantCategory struct {
	Name        string
	MainPhotoID uuid.UUID
	Params      []PlantParam
}
