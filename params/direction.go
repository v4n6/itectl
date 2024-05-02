package params

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/pkg/ite8291"
	"strings"
)

// DirectionDefault is default value of direction configuration property.
const DirectionDefault string = "none"

const (
	// DirectionProp is name of direction flag and configuration property.
	DirectionProp = "direction"
	// DirectionShortFlag is name of direction short flag.
	DirectionShortFlag = "d"
)

var directions = map[string]ite8291.Direction{
	"none":  ite8291.DirectionNone,
	"right": ite8291.DirectionRight,
	"left":  ite8291.DirectionLeft,
	"up":    ite8291.DirectionUp,
	"down":  ite8291.DirectionDown,
}

var directionNames = []string{"none", "right", "left", "up", "down"}

// ParseDirectionName parses given direction name
// to the corresponding ite8291.Direction value.
// It reports InvalidOptionValueError if direction name is not a valid name.
// Direction name is case insensitive.
func ParseDirectionName(name string) (ite8291.Direction, error) {

	if dir, found := directions[strings.ToLower(name)]; found {
		return dir, nil
	}
	return 0, fmt.Errorf("%w %q for %q; expected one of %q",
		InvalidOptionValueError, name,
		fmt.Sprintf("-%s, --%s", DirectionShortFlag, DirectionProp), directionNames)
}

// AddDirection adds direction flag to the provided cmd.
// It also adds hook to bind the flag to the corresponding viper config property
// and to validate the direction value.
// AddDirection returns function to retrieve current direction value.
func AddDirection(cmd *cobra.Command, v *viper.Viper) (direction func() ite8291.Direction) {

	var dir ite8291.Direction

	cmd.PersistentFlags().StringP(DirectionProp, DirectionShortFlag, DirectionDefault,
		fmt.Sprintf("Direction of the keyboard backlight effect %q. %s",
			directionNames, configurationWarning))
	bindAndValidate(cmd, v, DirectionProp, DirectionProp, func() (err error) {

		if dir, err = ParseDirectionName(v.GetString(DirectionProp)); err != nil {
			return err
		}

		return nil
	})

	return func() ite8291.Direction { return dir }
}
