package plant

import "fmt"

type FloweringPeriod string

const (
	January   FloweringPeriod = "january"
	February  FloweringPeriod = "february"
	March     FloweringPeriod = "march"
	April     FloweringPeriod = "april"
	May       FloweringPeriod = "may"
	June      FloweringPeriod = "june"
	July      FloweringPeriod = "july"
	August    FloweringPeriod = "august"
	September FloweringPeriod = "september"
	October   FloweringPeriod = "october"
	November  FloweringPeriod = "november"
	December  FloweringPeriod = "december"
	Winter    FloweringPeriod = "winter"
	Spring    FloweringPeriod = "spring"
	Summer    FloweringPeriod = "summer"
	Autumn    FloweringPeriod = "autumn"
)

func (fp FloweringPeriod) Validate() error {
	switch fp {
	case January, February, March, April, May, June, July, August, September, October, November, December, Winter, Spring, Summer, Autumn:
		return nil
	default:
		return fmt.Errorf("invalid flowering period: %s", fp)
	}
}
