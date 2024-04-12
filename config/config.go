package config

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// InvalidOptionValueError is an error indicating that a provided option has an invalid value.
var InvalidOptionValueError = errors.New("invalid option value")

// validateMinMaxUint8Value validates a uint8 option to be in a range provided by valMin and valMax.
func validateMinMaxUint8Value(name string, val, valMin, valMax uint8) error {

	if valMax < val || val < valMin {
		return fmt.Errorf("%w %q; expected [%d,%d] was %d", InvalidOptionValueError, name, valMin, valMax, val)
	}

	return nil
}

// validateMaxUint8Value validates a uint8 option to be less than given valMax.
func validateMaxUint8Value(name string, val, valMax uint8) error {

	if valMax < val {
		return fmt.Errorf("%w %q; expected [%d,%d] was %d", InvalidOptionValueError, name, 0, valMax, val)
	}

	return nil
}

// bindAndValidate replaces PersistentPreRunE hook with a new one.
// New hook calls the old one if existed before.
// It binds a cobra flag specified via flagName with viper config property specified via configProp.
// It also validates given property using specified via validate arg function, if provided.
func bindAndValidate(cmd *cobra.Command, flagName, configProp string, validate func() error) {

	origin := cmd.PersistentPreRunE

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {

		if origin != nil {
			if err := origin(cmd, args); err != nil {
				return err
			}
		}

		err := viper.BindPFlag(configProp, cmd.Flag(flagName))
		if err != nil {
			return err
		}

		if validate != nil {
			return validate()
		}

		return nil
	}
}

// addValidationHook replaces PersistentPreRunE hook with a new one.
// New hook calls the old one if existed before.
// It validates given property using specified via validate arg function.
func addValidationHook(cmd *cobra.Command, validate func() error) {

	origin := cmd.PersistentPreRunE

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {

		if origin != nil {
			if err := origin(cmd, args); err != nil {
				return err
			}
		}

		return validate()
	}
}
