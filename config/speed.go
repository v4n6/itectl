package config

import (
	"fmt"

	"github.com/spf13/cobra"
	ite8291 "github.com/v4n6/ite8291r3tool/pkg"
)

const speedDefault = 5

type SpeedProp struct {
	Speed uint8
}

func (c *SpeedProp) SpeedVal() (uint8, error) {
	err := validateMaxUint8Value("speed", c.Speed, ite8291.SpeedMaxValue)
	if err != nil {
		return 0, err
	}

	return ite8291.SpeedMaxValue - c.Speed, nil
}

// addSpeedFlag adds spped flag to the provided cmd
func AddSpeedFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Uint8VarP(&Config.Speed, "speed", "s", speedDefault, fmt.Sprintf("Speed of the keyboard backlight mode [0-%d]",
		ite8291.SpeedMaxValue))
}
