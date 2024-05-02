package params

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// BrightnessDefault default brightness value
const BrightnessDefault = 25

const (
	// BrightnessProp is name of brightness flag and config property.
	BrightnessProp = "brightness"
	// BrightnessShortFlag is name of brightness short flag
	BrightnessShortFlag = "b"
)

// AddBrightness adds brightness flag to the provided cmd.
// It also adds hook to bind the flag to the corresponding viper config property
// and to validate brightness value.
// It returns function to retrieve current brightness value.
func AddBrightness(cmd *cobra.Command, v *viper.Viper) (brightness func() byte) {

	cmd.PersistentFlags().Uint8P(BrightnessProp, BrightnessShortFlag, BrightnessDefault,
		fmt.Sprintf("Brightness of the keyboard backlight; min value 0, max value %d. %s", ite8291.BrightnessMaxValue, configurationWarning))
	bindAndValidate(cmd, v, BrightnessProp, BrightnessProp, func() error {
		return validateMaxUint8Value(fmt.Sprintf("-%s, --%s", BrightnessShortFlag, BrightnessProp),
			byte(v.GetUint(BrightnessProp)), ite8291.BrightnessMaxValue)
	})

	return func() byte { return byte(v.GetUint(BrightnessProp)) }
}
