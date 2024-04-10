package config

import (
	"fmt"

	"github.com/spf13/cobra"
	ite8291 "github.com/v4n6/ite8291r3tool/pkg"
)

const (
	colorNumDefault           = ite8291.ColorRandom
	assignableColorNumDefault = 1
	redDefault                = 0
	greenDefault              = 0
	blueDefault               = 0
)

type ColorsProp struct {
	ColorNum uint8
	ColorRed,
	ColorGreen,
	ColorBlue uint8
}

func (c *ColorsProp) ColorNumVal() (uint8, error) {
	err := validateMaxUint8Value("color-num", c.ColorNum, ite8291.ColorMaxValue)
	if err != nil {
		return 0, err
	}

	return c.ColorNum, nil
}

func (c *ColorsProp) AssignableColorNumVal() (uint8, error) {
	err := validateMinMaxUint8Value("color-num", c.ColorNum, ite8291.ColorNone+1, ite8291.ColorMaxValue-1)
	if err != nil {
		return 0, err
	}

	return c.ColorNum, nil
}

func (c *ColorsProp) ColorRedVal() uint8 {
	return c.ColorRed
}

func (c *ColorsProp) ColorGreenVal() uint8 {
	return c.ColorGreen
}

func (c *ColorsProp) ColorBlueVal() uint8 {
	return c.ColorBlue
}

func AddColorNumFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Uint8VarP(&Config.ColorNum, "color-num", "c", colorNumDefault,
		fmt.Sprintf("Number of the predfined color of keyboard backlight [%d-%d]",
			ite8291.ColorNone, ite8291.ColorMaxValue))
}

func AddAssignableColorNumFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Uint8VarP(&Config.ColorNum, "color-num", "c", assignableColorNumDefault,
		fmt.Sprintf("Number of the predfined color of keyboard backlight [%d-%d]",
			ite8291.ColorNone+1, ite8291.ColorMaxValue-1))
}

func AddRedColorFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Uint8Var(&Config.ColorRed, "red", redDefault, "Red atom of RGB color.")
}

func AddGreenColorFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Uint8Var(&Config.ColorGreen, "green", greenDefault, "Green atom of RGB color.")
}

func AddBlueColorFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Uint8Var(&Config.ColorBlue, "blue", blueDefault, "Blue atom of RGB color.")
}
