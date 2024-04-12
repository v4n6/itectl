package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	ite8291 "github.com/v4n6/ite8291r3tool/pkg"
)

const speedDefault = 5

const speedProp = "speed"

// Speed returns either specified, configured or default value of the speed flag.
func Speed() uint8 {
	return byte(ite8291.SpeedMaxValue - viper.GetUint(speedProp))
}

// AddSpeedFlag adds speed flag to the provided cmd and binds it to the corresponding viper config property.
// It also adds hook to validate speed value.
func AddSpeedFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Uint8P(speedProp, "s", speedDefault, fmt.Sprintf("Speed of the keyboard backlight mode [0-%d]",
		ite8291.SpeedMaxValue))

	bindAndValidate(cmd, speedProp, speedProp, func() error {
		return validateMaxUint8Value(speedProp, byte(viper.GetUint(speedProp)), ite8291.SpeedMaxValue)
	})
}
