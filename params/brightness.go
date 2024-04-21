package params

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/ite8291r3tool/pkg/ite8291"
)

// brightnessDefault default brightness value
const brightnessDefault = 25

// brightnessProp is name of brightness flag and config property.
const brightnessProp = "brightness"

// AddBrightness adds brightness flag to the provided cmd.
// It also adds hook to bind the flag to the corresponding viper config property
// and to validate brightness value.
// It returns function to retrieve current brightness value.
func AddBrightness(cmd *cobra.Command, v *viper.Viper) (brightness func() byte) {

	cmd.PersistentFlags().Uint8P(brightnessProp, "b", brightnessDefault,
		fmt.Sprintf("Brightness of the keyboard backlight [0-%d]", ite8291.BrightnessMaxValue))
	bindAndValidate(cmd, v, brightnessProp, brightnessProp, func() error {
		return validateMaxUint8Value(brightnessProp, byte(v.GetUint(brightnessProp)), ite8291.BrightnessMaxValue)
	})

	return func() byte { return byte(v.GetUint(brightnessProp)) }
}
