package plant

import "fmt"

type SoilAcidity int

func (a SoilAcidity) Validate() error {
	if a <= 0 {
		return fmt.Errorf("invalid soil acidity")
	}
	return nil
}
