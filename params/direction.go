package params

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/ite8291r3tool/pkg/ite8291"
	"strings"
)

// directionDefault is default value of direction configuration property.
const directionDefault string = "none"

// directionProp is name of direction flag and configuration property.
const directionProp = "direction"

var directions = map[string]ite8291.Direction{
	"none":  ite8291.DirectionNone,
	"right": ite8291.DirectionRight,
	"left":  ite8291.DirectionLeft,
	"up":    ite8291.DirectionUp,
	"down":  ite8291.DirectionDown,
}

var directionNames = []string{"none", "right", "left", "up", "down"}

// AddDirection adds direction flag to the provided cmd.
// It also adds hook to bind the flag to the corresponding viper config property
// and to validate the direction value.
// It returns function to retrieve current direction value.
func AddDirection(cmd *cobra.Command, v *viper.Viper) (direction func() ite8291.Direction) {

	var d ite8291.Direction

	cmd.PersistentFlags().StringP("direction", "d", directionDefault,
		fmt.Sprintf("Direction of the keyboard backlight effect %q", directionNames))
	bindAndValidate(cmd, v, directionProp, directionProp, func() error {

		found, val := false, strings.ToLower(v.GetString(directionProp))
		if d, found = directions[val]; !found {
			return fmt.Errorf("%w direction; expected on of %q was %q", InvalidOptionValueError, directionNames, val)
		}

		return nil
	})

	return func() ite8291.Direction { return d }
}
