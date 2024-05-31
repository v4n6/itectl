package params

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// SpeedDefault is default speed value.
const SpeedDefault = 5

const (
	// SpeedProp is name of speed flag and config property.
	SpeedProp = "speed"
	// SpeedShortFlag is name of speed short flag.
	SpeedShortFlag = "s"
)

// AddSpeed adds speed flag to the provided cmd.
// It also adds hook to bind it to the corresponding viper configuration property
// and to validate its value.
func AddSpeed(cmd *cobra.Command, v *viper.Viper) {

	cmd.PersistentFlags().Uint8P(SpeedProp, SpeedShortFlag, SpeedDefault,
		fmt.Sprintf("Speed of the keyboard backlight effect; min value 0, max value %d. %s",
			ite8291.SpeedMaxValue, configurationWarning))
	bindAndValidate(cmd, v, SpeedProp, SpeedProp, func() error {
		return validateMaxUint8Value(fmt.Sprintf("-%s, --%s", SpeedShortFlag, SpeedProp),
			byte(v.GetUint(SpeedProp)), ite8291.SpeedMaxValue)
	})
}

// Speed returns speed property value.
func Speed(v *viper.Viper) byte {
	return byte(ite8291.SpeedMaxValue - v.GetUint(SpeedProp))
}
