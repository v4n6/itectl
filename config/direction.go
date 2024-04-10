package config

import (
	"fmt"

	"github.com/spf13/cobra"
	ite8291 "github.com/v4n6/ite8291r3tool/pkg"
)

const directionDefault string = "none"

var directions = map[string]ite8291.Direction{
	"none":  ite8291.DirectionNone,
	"right": ite8291.DirectionRight,
	"left":  ite8291.DirectionLeft,
	"up":    ite8291.DirectionUp,
	"down":  ite8291.DirectionDown,
}

var directionNames = []string{"none", "right", "left", "up", "down"}

type DirectionProp struct {
	Direction string
}

func (c *DirectionProp) DirectionVal() (ite8291.Direction, error) {

	direction, found := directions[c.Direction]
	if !found {
		return 0, fmt.Errorf("%w direction; expected on of %q was %q", InvalidOptionValueError, directionNames, c.Direction)
	}

	return direction, nil
}

func AddDirectionNameFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&Config.Direction, "direction", "d", directionDefault, fmt.Sprintf("Direction of the keyboard backlight effect %q",
		directionNames))
}
