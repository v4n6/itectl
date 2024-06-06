package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/params"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// rainbowModeDescription - rainbow-mode command description.
const rainbowModeDescription = "Set keyboard backlight to 'rainbow' mode."

// newRainbowModeCmd creates, initializes and returns command
// to set keyboard backlight to 'rainbow' mode.
func newRainbowModeCmd(v *viper.Viper, call ite8291Ctl) *cobra.Command {

	var rainbowModeCmd = &cobra.Command{
		Use:           "rainbow-mode",
		Short:         rainbowModeDescription,
		Long:          rainbowModeDescription,
		Args:          cobra.NoArgs,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return call(cmd, func(ctl *ite8291.Controller) error {
				return ctl.SetRainbowMode(params.Brightness(v), params.Save(v))
			})
		},
	}

	params.AddBrightness(rainbowModeCmd, v)
	params.AddSave(rainbowModeCmd, v)
	params.AddReset(rainbowModeCmd, v)

	return rainbowModeCmd
}
