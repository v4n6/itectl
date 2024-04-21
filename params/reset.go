package params

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/ite8291r3tool/pkg/ite8291"
)

// resetDefault is a default value of "reset" property.
const resetDefault = false

// resetFlag is name of "reset" flag.
const resetFlag = "reset"

const (
	// predefinedColorPropTemplate is template of the name of predefined color property.
	predefinedColorPropTemplate = "predefinedColors.color_%d"
	// predefinedColorsNumber is a number of assignable predefined colors.
	predefinedColorsNumber = ite8291.AssignableColorNumMaxValue - ite8291.AssignableColorNumMinValue + 1
)

// default values of predefined colors.
var predefinedColorsDefault []string = []string{
	"#FFFFFF",
	"#FF0000",
	"#FFFF00",
	"#00FF00",
	"#0000FF",
	"#00FFFF",
	"#FF00FF",
}

// AddReset adds "reset" flag to the provided cmd.
// It also adds hook to validate configured predefined colors if "reset" property is true.
// It returns functions to retrieve current "reset" and predefinedColors values.
func AddReset(cmd *cobra.Command, v *viper.Viper) (reset func() bool, predefinedColors func() []*ite8291.Color) {

	r := resetDefault

	colors := make([]*ite8291.Color, predefinedColorsNumber)

	cmd.PersistentFlags().BoolVar(&r, resetFlag, resetDefault, "Reset the controller predefined colors to the initailly configured state.")

	addValidationHook(cmd, func() (err error) {

		if r {

			for i := range predefinedColorsNumber {
				n := i + 1

				val := v.GetString(fmt.Sprintf(predefinedColorPropTemplate, n))
				if val == "" {
					val = predefinedColorsDefault[i]
				}

				if colors[i], err = colorNameToColor(val, v); err == nil {
					continue
				}

				if colors[i], err = ite8291.ParseColor(val); err != nil {
					return fmt.Errorf("%w of color #%d: %w", InvalidOptionValueError, n, err)
				}
			}
		}

		return nil
	})

	return func() bool { return r },
		func() []*ite8291.Color { return colors }
}
