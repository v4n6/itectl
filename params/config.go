package params

import (
	"errors"
	"fmt"
	"slices"

	"github.com/adrg/xdg"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// viper configuration file properties
const (
	ConfigName = "ite8291r3"

	ConfigType = "yaml"
)

// defaultModeProp is a viper config property name for default mode
const defaultModeProp = "mode"

// ConfigFileFlag is name of viper config file flag
const ConfigFileFlag = "config"

// InvalidOptionValueError is an error indicating that a provided option has an invalid value.
var InvalidOptionValueError = errors.New("invalid option value")

// DefaultMode returns name of configured default mode without "-mode" suffix
func DefaultMode(v *viper.Viper) string {
	return v.GetString(defaultModeProp)
}

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

// bindAndValidate sets PersistentPreRunE hook of the given cmd.
// New hook calls the original PersistentPreRunE hook if it existed before.
// It binds a flag specified via flagName with viper config property specified via configProp.
// It also validates given property using provided via validate arg function, if given.
func bindAndValidate(cmd *cobra.Command, v *viper.Viper, flagName, configProp string, validate func() error) {

	origin := cmd.PersistentPreRunE

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {

		if origin != nil {
			if err := origin(cmd, args); err != nil {
				return err
			}
		}

		if err := v.BindPFlag(configProp, cmd.Flag(flagName)); err != nil {
			return err
		}

		if validate != nil {
			return validate()
		}

		return nil
	}
}

// addValidationHook sets PersistentPreRunE hook of the given cmd.
// New hook calls the original one if it existed before.
// It validates command using provided via validate arg function.
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

// ReadConfig initializes viper instance specified via v.
// It reads configuration from specific config file if provided by the user
// via flag and reports error if it doesn't exist or read error occurred.
// If no flag provided, it reads and merges configuration first from global
// config files from XDG directories and then from config file in user XDG
// directory. In this case errors are printed as warnings and no error is reported.
// No warning is printed if no default config file(s) found.
func ReadConfig(cmd *cobra.Command, v *viper.Viper, args []string) error {

	_ = cmd.Flags().Parse(args) // ignore errors

	if cfgFlag := cmd.Flag(ConfigFileFlag); cfgFlag.Changed {
		// Use config cfgFile provided via flag.
		cfgFile, err := cast.ToStringE(cfgFlag.Value)
		if err != nil {
			return err
		}

		if err := readInConfig(v, cfgFile, ""); err != nil {
			return fmt.Errorf("error reading configuration file %q: %w", cfgFile, err)
		}

		return nil
	}

	// merges all default config files
	for _, dir := range slices.Concat(xdg.ConfigDirs, []string{xdg.ConfigHome}) {

		if err := mergeConfig(cmd, v, dir); err != nil {
			// print and ignore errors
			fmt.Fprintf(cmd.ErrOrStderr(), "Warning: %v\n", err)
		}
	}

	return nil
}

// mergeConfig merges configuration found in specified dir into provided via v param viper instance
func mergeConfig(cmd *cobra.Command, v *viper.Viper, dir string) error {

	v_ := viper.New() // local viper instance

	if err := readInConfig(v_, "", dir); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// return error if config file was found and error occured
			return err
		}
	}

	return v.MergeConfigMap(v_.AllSettings())
}

// readInConfig reads configuration into v from config file specified via cfgFile, if provided.
// It reads it from file in dir if cfgFile is set to default value.
func readInConfig(v *viper.Viper, cfgFile, dir string) error {

	v.SetConfigType(ConfigType)

	if len(cfgFile) > 0 {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigName(ConfigName)
		v.AddConfigPath(dir)
	}

	return v.ReadInConfig()
}
