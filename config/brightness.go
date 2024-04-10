package config

import (
	"fmt"

	"github.com/spf13/cobra"
	ite8291 "github.com/v4n6/ite8291r3tool/pkg"
)

const brightnessDefault = 25

type BrightnessProp struct {
	Brightness uint8
}

func (c *BrightnessProp) BrightnessVal() (uint8, error) {
	err := validateMaxUint8Value("brightness", c.Brightness, ite8291.BrightnessMaxValue)
	if err != nil {
		return 0, err
	}

	return c.Brightness, nil
}

// addBrightnessFlag add brightness flag the provided cmd
func AddBrightnessFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Uint8VarP(&Config.Brightness, "brightness", "b", brightnessDefault,
		fmt.Sprintf("Brightness of the keyboard backlight [0-%d]", ite8291.BrightnessMaxValue))
}
