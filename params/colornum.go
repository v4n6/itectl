package params

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/ite8291r3tool/pkg/ite8291"
)

const (
	// colorNumDefault is color number property default value
	colorNumDefault = ite8291.ColorRandom
	// assignableColorNumDefault is assignable color number property default value
	assignableColorNumDefault = ite8291.AssignableColorNumMinValue
)

const (
	// colorNumProp is name of color number configuration property.
	colorNumProp = "colorNum"
	// colorNumFlag is name of color number flag.
	colorNumFlag = "color-num"
)

// AddColorNum adds color number flag to the provided cmd.
// It also adds hook to bind it to the corresponding viper config property.
// It returns functions to retrieve current color number property value.
func AddColorNum(cmd *cobra.Command, v *viper.Viper) (colorNum func() byte) {

	cmd.PersistentFlags().Uint8P(colorNumFlag, "c", colorNumDefault,
		fmt.Sprintf("Number of the predfined color of keyboard backlight [%d-%d]",
			ite8291.ColorNumMinValue, ite8291.ColorNumMaxValue))
	bindAndValidate(cmd, v, colorNumFlag, colorNumProp, func() error {
		return validateMaxUint8Value(colorNumFlag, byte(v.GetUint(colorNumProp)),
			ite8291.ColorNumMaxValue)
	})

	return func() byte { return byte(v.GetUint(colorNumProp)) }
}

// AddAssignableColorNum adds assignable color number flag to the provided cmd.
// It also adds hook to bind it to the corresponding viper config property.
// It returns functions to retrieve current assignable color number property value.
func AddAssignableColorNum(cmd *cobra.Command, v *viper.Viper) (assignableColorNum func() byte) {

	cmd.PersistentFlags().Uint8P(colorNumFlag, "c", assignableColorNumDefault,
		fmt.Sprintf("Number of the predfined color of keyboard backlight [%d-%d]",
			ite8291.AssignableColorNumMinValue, ite8291.AssignableColorNumMaxValue))
	bindAndValidate(cmd, v, colorNumFlag, colorNumProp, func() error {
		return validateMinMaxUint8Value(colorNumFlag, byte(v.GetUint(colorNumProp)),
			ite8291.AssignableColorNumMinValue, ite8291.AssignableColorNumMaxValue)
	})
	return func() byte { return byte(v.GetUint(colorNumProp)) }
}
