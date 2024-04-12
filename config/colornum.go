package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	ite8291 "github.com/v4n6/ite8291r3tool/pkg"
)

const (
	colorNumDefault = ite8291.ColorRandom

	assignableColorNumDefault = 1
)

const colorNumFlag = "color-num"

const colorNumProp = "colorNum"

var assignableColorNum uint8

// ColorNum returns either specified, configured or default value of the color-num flag.
func ColorNum() uint8 {
	return byte(viper.GetUint(colorNumProp))
}

// AssignableColorNum returns value of the color-num flag indicating color number to assign a color to.
func AssignableColorNum() uint8 {
	return assignableColorNum
}

// AddColorNumFlag adds color-num flag to the provided cmd and binds it to the corresponding viper config property.
// It also adds hook to validate color number value.
func AddColorNumFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Uint8P(colorNumFlag, "c", colorNumDefault,
		fmt.Sprintf("Number of the predfined color of keyboard backlight [%d-%d]",
			ite8291.ColorNone, ite8291.ColorMaxValue))

	bindAndValidate(cmd, colorNumFlag, colorNumProp, func() error {
		return validateMaxUint8Value(colorNumFlag, ColorNum(), ite8291.ColorMaxValue)
	})
}

// AddAssignableColorNumFlag adds color-num flag to the provided cmd.
// It also adds hook to validate color number value.
func AddAssignableColorNumFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Uint8VarP(&assignableColorNum, colorNumFlag, "c", assignableColorNumDefault,
		fmt.Sprintf("Number of the predfined color of keyboard backlight [%d-%d]",
			ite8291.ColorNone+1, ite8291.ColorMaxValue-1))

	addValidationHook(cmd, func() error {
		return validateMinMaxUint8Value(colorNumFlag, ColorNum(), ite8291.ColorNone+1, ite8291.ColorMaxValue-1)
	})

}
