package plant

import "fmt"

type LightRelation string

const (
	Shadow     LightRelation = "shadow"
	HalfShadow LightRelation = "halfshadow"
	Light      LightRelation = "light"
)

func (l *LightRelation) Validate() error {
	switch *l {
	case Shadow, HalfShadow, Light:
		return nil
	default:
		return fmt.Errorf("invalid light sensitivity: %s", *l)
	}
}
