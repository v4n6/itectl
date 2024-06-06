package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/params"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// singleColorModeDescription - single-color-mode command description.
const singleColorModeDescription = "Set keyboard backlight to 'single color' mode"

// newSingleColorModeCmd creates, initializes and returns command to
// set keyboard backlight to 'single color' mode.
func newSingleColorModeCmd(v *viper.Viper, call ite8291Ctl) *cobra.Command {

	var color func() *ite8291.Color

	var singleColorModeCmd = &cobra.Command{
		Use:   "single-color-mode",
		Short: singleColorModeDescription,
		Long: fmt.Sprintf(`Set keyboard backlight to 'single color' mode.

In the 'single color' mode keyboard backlight is set to the same color.
The color can be by given by a name "(--%s)" of the color cobfigured via %q configuration property.
e.g. %[2]s:
       azure: "#007FFF"

It can also be specified by RGB string "(--%s)" directly in a one of the following formats %q.
The color can also be provided by a combination of (--%s, --%s, --%s) flags.

If color is not provided directly via flag(s), the value specified by %q configuration property will be used.
It can be set to either a color name or an rgb directly.`,
			params.ColorNameFlag, params.NamedColorsProp,
			params.ColorRGBFlag, ite8291.SupportedColorStringFormats,
			params.ColorRedFlag, params.ColorGreenFlag, params.ColorBlueFlag,
			params.SingleColorProp),
		Args:          cobra.NoArgs,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return call(cmd, func(ctl *ite8291.Controller) error {
				return ctl.SetSingleColorMode(params.Brightness(v), color(),
					params.Save(v))
			})
		},
	}

	color = params.AddSingleModeColor(singleColorModeCmd, v)
	params.AddBrightness(singleColorModeCmd, v)
	params.AddSave(singleColorModeCmd, v)
	params.AddReset(singleColorModeCmd, v)

	return singleColorModeCmd
}
