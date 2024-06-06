package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/params"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// setColorDescription - set-color command description.
const setColorDescription = "Set keyboard backlight customizable predefined color."

// newSetColorCmd creates, initializes and returns command to set
// keyboard backlight customizable predefined color.
func newSetColorCmd(v *viper.Viper, call ite8291Ctl) *cobra.Command {

	var colorNum func() byte
	var color func() *ite8291.Color

	var setColorCmd = &cobra.Command{
		Use:   "set-color",
		Short: setColorDescription,
		Long: fmt.Sprintf(`Set a customizable predefined color of the keyboard backlight to the specified value.

The number of the predefined color to set must be specified by "(-%s,--%s)" flag.
The color value can be by given by a name "(--%s)" of the color cobfigured via %q configuration property.
e.g. %[4]s:
       azure: "#007FFF"

It can also be specified by RGB string "(--%s)" directly in a one of the following formats %q.
The color can also be provided by a combination of (--%s, --%s, --%s) flags.`,
			params.ColorNumShortFlag, params.ColorNumFlag,
			params.ColorNameFlag, params.NamedColorsProp,
			params.ColorRGBFlag, ite8291.SupportedColorStringFormats,
			params.ColorRedFlag, params.ColorGreenFlag, params.ColorBlueFlag),
		Args:          cobra.NoArgs,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return call(cmd, func(ctl *ite8291.Controller) error {
				return ctl.SetColor(colorNum(), color())
			})
		},
	}

	colorNum = params.AddCustomColorNum(setColorCmd)
	color = params.AddColor(setColorCmd, v)

	return setColorCmd
}
