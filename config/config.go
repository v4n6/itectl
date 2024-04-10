package config

import (
	"errors"
	"fmt"
)

type Configuration struct {
	PollProp
	SpeedProp
	BrightnessProp
	ColorsProp
	DirectionProp
	ReactiveProp
	SaveProp
}

var Config = &Configuration{}

// InvalidOptionValueError is an error indicating that a provided option has an invalid value.
var InvalidOptionValueError = errors.New("invalid option value")

// validateMinMaxUint8Value validates a uint8 option to be in a range provided by valMin and valMax.
func validateMinMaxUint8Value(name string, val, valMin, valMax byte) error {

	if valMax < val || val < valMin {
		return fmt.Errorf("%w %q; expected [%d,%d] was %d", InvalidOptionValueError, name, valMin, valMax, val)
	}

	return nil
}

// validateMaxUint8Value validates a uint8 option to be less than given valMax.
func validateMaxUint8Value(name string, val, valMax byte) error {

	if valMax < val {
		return fmt.Errorf("%w %q; expected [%d,%d] was %d", InvalidOptionValueError, name, 0, valMax, val)
	}

	return nil
}
