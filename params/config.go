package params

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// viper configuration file properties
const (
	// ConfigName is name of configuration file(s) used by viper
	ConfigName = "itectl"
	// ConfigType is type of configuration file(s) used by viper
	ConfigType = "yaml"
)

const (
	// defaultModeProp is viper config property name for implicit mode
	defaultModeProp = "mode"
	// ConfigFileFlag is name of the viper config file flag
	ConfigFileFlag = "config"
	// configurationWarning is warning about possibility that value can be retrieved from viper configuration
	configurationWarning = "The default value can be overwritten by global and/or user configuration!"
)

// InvalidOptionValueError is an error indicating that a property has an invalid value.
var InvalidOptionValueError = errors.New("invalid option value")

// DefaultMode returns name of configured implicit mode without "-mode" suffix
func DefaultMode(v *viper.Viper) string {
	return v.GetString(defaultModeProp)
}

// validateMinMaxUint8Value validates a uint8 option to be in a range provided by valMin and valMax.
func validateMinMaxUint8Value(name string, val, valMin, valMax uint8) error {

	if valMax < val || val < valMin {
		return fmt.Errorf("%w \"%d\" for %q; expected [%d,%d]", InvalidOptionValueError, val, name, valMin, valMax)
	}

	return nil
}

// validateMaxUint8Value validates a uint8 option to be less than given valMax.
func validateMaxUint8Value(name string, val, valMax uint8) error {

	if valMax < val {
		return fmt.Errorf("%w \"%d\" for %q; expected [%d,%d]", InvalidOptionValueError, val, name, 0, valMax)
	}

	return nil
}

// bindAndValidate sets PersistentPreRunE hook of the given cmd.
// New hook calls the original PersistentPreRunE hook first if it existed before.
// It binds a flag specified via flagName with viper config property specified via configProp.
// It also validates given property using provided via validate parameter function, if given.
func bindAndValidate(cmd *cobra.Command, v *viper.Viper, flagName, configProp string, validate func() error) {

	origin := cmd.PersistentPreRunE

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {

		if origin != nil {
			if err := origin(cmd, args); err != nil {
				return err
			}
		}

		if err := v.BindPFlag(configProp, cmd.Flag(flagName)); err != nil {
			panic(err)
		}

		if validate != nil {
			return validate()
		}

		return nil
	}
}

// addValidationHook sets PersistentPreRunE hook of the given cmd.
// New hook calls the original PersistentPreRunE hook first if it existed before.
// It validates properties using provided via validate parameter function.
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

// AddConfigFlag adds config flag to the given cmd.
// config flag identifies configuration file to use instead of default.
func AddConfigFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(ConfigFileFlag, "",
		fmt.Sprintf("Configuration file to use instead of xdg config files [\"/etc/%[1]s.%[2]s\",\"~/.config/%[1]s.%[2]s\"]",
			ConfigName, ConfigType))
}

// ReadConfig adds retrieved configuration to viper instance specified via v.
// If cfgFile is provided, ReadConfig reads configuration from that file and merges it with given v viper.
// It returns error if specified file does not exist or io error occurred.
// If cfgFile is empty, ReadConfig retrieves configuration from default global/user XDG config files.
// In this case all occurred errors are printed as warnings and no error is returned.
// No warning is printed if default config file(s) not found.
func ReadConfig(cmd *cobra.Command, v *viper.Viper, cfgFile string) error {

	if len(cfgFile) > 0 {

		if err := readInConfig(v, cfgFile, ""); err != nil {
			return fmt.Errorf("error reading configuration file %q: %w", cfgFile, err)
		}

		return nil
	}

	// merges all default config files
	for _, dir := range slices.Concat(xdg.ConfigDirs, []string{xdg.ConfigHome}) {

		if err := mergeConfig(cmd, v, dir); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Warning: %v\n", err) // print and ignore errors
		}
	}

	v.SetEnvPrefix("ITECTL")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv() // read in environment variables that match

	return nil
}

// mergeConfig merges configuration found in specified dir into provided via v param viper instance
func mergeConfig(cmd *cobra.Command, v *viper.Viper, dir string) error {

	v_ := viper.New() // local viper instance

	if err := readInConfig(v_, "", dir); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err // return error if config file was found and error occured
		}
	}

	return v.MergeConfigMap(v_.AllSettings())
}

// readInConfig retrieves properties into viper instance specified via v
// from config file specified via cfgFile, if provided.
// It retrieves them from file in dir if cfgFile is not provided.
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

// ConfigFile returns value of config flag.
func ConfigFile(cmd *cobra.Command, flags []string) (cfgFile string, err error) {
	fs := pflag.NewFlagSet("config", pflag.ContinueOnError)
	fs.ParseErrorsWhitelist.UnknownFlags = true
	fs.AddFlagSet(cmd.Flags())
	fs.AddFlagSet(cmd.PersistentFlags())

	_ = fs.Parse(flags) // ignore errors

	return fs.GetString(ConfigFileFlag)
}
