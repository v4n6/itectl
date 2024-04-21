package params

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/ite8291r3tool/pkg/ite8291"
)

// speedDefault is default speed value
const speedDefault = 5

// speedProp is name of speed flag and config property.
const speedProp = "speed"

// AddSpeed adds speed flag to the provided cmd.
// It also adds hook to bind it to the corresponding viper config property
// and to validate speed value.
// It returns function to retrieve current speed value.
func AddSpeed(cmd *cobra.Command, v *viper.Viper) (speed func() byte) {

	cmd.PersistentFlags().Uint8P(speedProp, "s", speedDefault, fmt.Sprintf("Speed of the keyboard backlight mode [0-%d]",
		ite8291.SpeedMaxValue))
	bindAndValidate(cmd, v, speedProp, speedProp, func() error {
		return validateMaxUint8Value(speedProp, byte(v.GetUint(speedProp)), ite8291.SpeedMaxValue)
	})

	return func() byte { return byte(ite8291.SpeedMaxValue - v.GetUint(speedProp)) }
}
