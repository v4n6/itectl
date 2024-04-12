package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	ite8291 "github.com/v4n6/ite8291r3tool/pkg"
)

const directionDefault string = "none"

const directionProp = "direction"

var directions = map[string]ite8291.Direction{
	"none":  ite8291.DirectionNone,
	"right": ite8291.DirectionRight,
	"left":  ite8291.DirectionLeft,
	"up":    ite8291.DirectionUp,
	"down":  ite8291.DirectionDown,
}

var directionNames = []string{"none", "right", "left", "up", "down"}

// Direction returns either specified, configured or default value of the direction flag.
func Direction() ite8291.Direction {

	return directions[viper.GetString(directionProp)]
}

// AddDirectionFlag adds direction flag to the provided cmd and binds it to the corresponding viper config property.
// It also adds hook to validate direction value.
func AddDirectionFlag(cmd *cobra.Command) {

	cmd.PersistentFlags().StringP("direction", "d", directionDefault, fmt.Sprintf("Direction of the keyboard backlight effect %q",
		directionNames))

	bindAndValidate(cmd, directionProp, directionProp, func() error {

		d, found := directions[viper.GetString(directionProp)]
		if !found {
			return fmt.Errorf("%w direction; expected on of %q was %q", InvalidOptionValueError, directionNames, d)
		}

		return nil
	})
}
