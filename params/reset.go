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
// It also adds hook to bind it to the corresponding viper config property.
func AddReset(cmd *cobra.Command, v *viper.Viper) {

	cmd.PersistentFlags().Bool(ResetProp, ResetDefault,
		fmt.Sprintf(
			"Reset the controller customizable predefined colors to their corresponding configured/default values. %s",
			configurationWarning))
	bindAndValidate(cmd, v, ResetProp, ResetProp, nil)
}

// Reset returns value of reset property
func Reset(v *viper.Viper) bool {
	return v.GetBool(ResetProp)
}

// PredefinedColor returns color of i-th predefined color
func PredefinedColor(v *viper.Viper, i int) (color *ite8291.Color, err error) {

	val := v.GetString(fmt.Sprintf(predefinedColorPropTemplate, PredefinedColorProp, i))
	if val == "" {
		val = PredefinedColorsDefault[i-1]
	}

	// try as color name
	if color, err = colorNameToColor(val, v); err == nil {
		return color, nil
	}
	// it isn't color name -> try as rgb
	if color, err = ite8291.ParseColor(val); err != nil {
		return nil, fmt.Errorf("%w for predefined color #%d: %w",
			InvalidOptionValueError, i, err)
	}

	return color, nil
}
