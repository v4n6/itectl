package params

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/pkg/ite8291"
)

const (
	// ColorNumDefault is color number property default value
	ColorNumDefault = ite8291.ColorRandom
)

const (
	// ColorNumProp is name of color number configuration property.
	ColorNumProp = "colorNum"
	// ColorNumFlag is name of color number flag.
	ColorNumFlag = "color-num"
	// ColorNumFlag is name of color number short flag.
	ColorNumShortFlag = "c"
)

// AddColorNum adds color number flag to the provided cmd.
// It also adds hook to bind it to the corresponding viper config property
// and to validate its value.
// It returns function to retrieve current color number property value.
func AddColorNum(cmd *cobra.Command, v *viper.Viper) (colorNum func() byte) {

	cmd.PersistentFlags().Uint8P(ColorNumFlag, ColorNumShortFlag, ColorNumDefault,
		fmt.Sprintf(
			"Number of the predfined color of keyboard backlight to use; min value %d, max value %d, 0 means no color, 1-7 customizable color, 8 random color. %s",
			ite8291.ColorNumMinValue, ite8291.ColorNumMaxValue, configurationWarning))
	bindAndValidate(cmd, v, ColorNumFlag, ColorNumProp, func() error {
		return validateMaxUint8Value(fmt.Sprintf("-%s, --%s", ColorNumShortFlag, ColorNumFlag),
			byte(v.GetUint(ColorNumProp)), ite8291.ColorNumMaxValue)
	})

	return func() byte { return byte(v.GetUint(ColorNumProp)) }
}

// AddCustomColorNum adds customizable color number flag to the provided cmd.
// It also adds hook to validate its value.
// It returns function to retrieve current color number flag value.
func AddCustomColorNum(cmd *cobra.Command, v *viper.Viper) (assignableColorNum func() byte) {

	var customColorNum byte

	cmd.PersistentFlags().Uint8VarP(&customColorNum, ColorNumFlag, ColorNumShortFlag, 0,
		fmt.Sprintf("Number of the predfined color of keyboard backlight to set; min value %d, max value %d.",
			ite8291.CustomColorNumMinValue, ite8291.CustomColorNumMaxValue))
	addValidationHook(cmd, func() error {
		if !cmd.Flag(ColorNumFlag).Changed {
			return nil
		}

		return validateMinMaxUint8Value(fmt.Sprintf("-%s, --%s", ColorNumShortFlag, ColorNumFlag), customColorNum,
			ite8291.CustomColorNumMinValue, ite8291.CustomColorNumMaxValue)
	})

	if err := cmd.MarkPersistentFlagRequired(ColorNumFlag); err != nil {
		panic(err)
	}

	return func() byte { return customColorNum }
}
