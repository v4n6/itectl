package params

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// ResetDefault is a default value of "reset" property.
const ResetDefault = false

// ResetProp is name of reset flag and configuration property.
const ResetProp = "reset"

const (
	// PredefinedColorProp name of predefined colors configuration property.
	PredefinedColorProp = "predefinedColors"
	// predefinedColorPropTemplate is template of the name of predefined color property.
	predefinedColorPropTemplate = "%s.color%d"
	// predefinedColorsNumber is a number of customizable predefined colors.
	predefinedColorsNumber = ite8291.CustomColorNumMaxValue - ite8291.CustomColorNumMinValue + 1
)

// default values of predefined colors.
var PredefinedColorsDefault []string = []string{
	"#FFFFFF",
	"#FF0000",
	"#FFFF00",
	"#00FF00",
	"#0000FF",
	"#00FFFF",
	"#FF00FF",
}

// AddReset adds "reset" flag to the provided cmd.
// It also adds hook to validate configured predefined colors,
// if "reset" property is set to true.
// AddReset returns function to reset all customizable predefined colors
// to their corresponding configured/default values, if "reset" is set to true or do nothing otherwise.
func AddReset(cmd *cobra.Command, v *viper.Viper) (optionallyResetColors func(ctl *ite8291.Controller) error) {

	colors := make([]*ite8291.Color, predefinedColorsNumber)

	cmd.PersistentFlags().Bool(ResetProp, ResetDefault,
		fmt.Sprintf("Reset the controller customizable predefined colors to their corresponding configured/default values. %s",
			configurationWarning))
	bindAndValidate(cmd, v, ResetProp, ResetProp, func() (err error) {

		if v.GetBool(ResetProp) {

			for i := range predefinedColorsNumber {
				n := i + 1

				val := v.GetString(fmt.Sprintf(predefinedColorPropTemplate, PredefinedColorProp, n))
				if val == "" {
					val = PredefinedColorsDefault[i]
				}
				// try as color name
				if colors[i], err = colorNameToColor(val, v); err == nil {
					continue
				}
				// it isn't color name -> try as rgb
				if colors[i], err = ite8291.ParseColor(val); err != nil {
					return fmt.Errorf("%w for predefined color #%d: %w",
						InvalidOptionValueError, n, err)
				}
			}
		}

		return nil
	})

	return func(ctl *ite8291.Controller) error {
		if v.GetBool(ResetProp) {
			return ctl.SetColors(colors)
		}
		return nil
	}
}
