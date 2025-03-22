package plant

import "fmt"

type SoilMoisture string

const (
	DryMoisture    SoilMoisture = "dry"
	LowMoisture    SoilMoisture = "low"
	MediumMoisture SoilMoisture = "medium"
	HighMoisture   SoilMoisture = "high"
)

func (s *SoilMoisture) Validate() error {
	switch *s {
	case DryMoisture, LowMoisture, MediumMoisture, HighMoisture:
		return nil
	default:
		return fmt.Errorf("invalid soil moisture: %s", *s)
	}
}
