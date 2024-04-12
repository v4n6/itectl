package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	ite8291 "github.com/v4n6/ite8291r3tool/pkg"
)

const brightnessDefault = 25

const brightnessProp = "brightness"

// Brightness returns either specified, configured or default value of the brightness flag.
func Brightness() uint8 {
	return byte(viper.GetUint(brightnessProp))
}

// AddBrightnessFlag adds brightness flag to the provided cmd and binds it to the corresponding viper config property.
// It also adds hook to validate brightness value.
func AddBrightnessFlag(cmd *cobra.Command) {

	cmd.PersistentFlags().Uint8P(brightnessProp, "b", brightnessDefault,
		fmt.Sprintf("Brightness of the keyboard backlight [0-%d]", ite8291.BrightnessMaxValue))

	bindAndValidate(cmd, brightnessProp, brightnessProp, func() error {
		return validateMaxUint8Value(brightnessProp, Brightness(), ite8291.BrightnessMaxValue)
	})
}
