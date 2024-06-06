package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/params"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// raindropModeDescription - raindrop-mode command description.
const raindropModeDescription = "Set keyboard backlight to 'raindrop' mode."

// newRaindropModeCmd creates, initializes and returns command
// to set keyboard backlight to 'raindrop' mode.
func newRaindropModeCmd(v *viper.Viper, call ite8291Ctl) *cobra.Command {

	var raindropModeCmd = &cobra.Command{
		Use:           "raindrop-mode",
		Short:         raindropModeDescription,
		Long:          raindropModeDescription,
		Args:          cobra.NoArgs,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return call(cmd, func(ctl *ite8291.Controller) error {
				return ctl.SetRaindropMode(params.Speed(v), params.Brightness(v),
					params.ColorNum(v), params.Save(v))
			})
		},
	}

	params.AddSpeed(raindropModeCmd, v)
	params.AddBrightness(raindropModeCmd, v)
	params.AddColorNum(raindropModeCmd, v)
	params.AddSave(raindropModeCmd, v)
	params.AddReset(raindropModeCmd, v)

	return raindropModeCmd
}
