package plant

import "fmt"

type WinterHardiness int

func (w WinterHardiness) Validate() error {
	if w < 1 || w > 11 {
		return fmt.Errorf("invalid winter hardiness: %d", w)
	}
	return nil
}
