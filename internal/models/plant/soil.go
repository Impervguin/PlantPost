package plant

import "fmt"

type Soil string

const (
	LightSoil  Soil = "light"
	MediumSoil Soil = "medium"
	HeavySoil  Soil = "heavy"
)

func (s *Soil) Validate() error {
	switch *s {
	case LightSoil, MediumSoil, HeavySoil:
		return nil
	default:
		return fmt.Errorf("invalid soil type: %s", *s)
	}
}
