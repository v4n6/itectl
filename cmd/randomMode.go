package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/params"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// randomModeDescription - random-mode command description.
const randomModeDescription = "Set keyboard backlight to 'random' mode."

// newRandomModeCmd creates, initializes and returns command
// to set keyboard backlight to 'random' mode.
func newRandomModeCmd(v *viper.Viper, call ite8291Ctl) *cobra.Command {

	var randomModeCmd = &cobra.Command{
		Use:           "random-mode",
		Short:         randomModeDescription,
		Long:          randomModeDescription,
		Args:          cobra.NoArgs,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return call(cmd, func(ctl *ite8291.Controller) error {
				return ctl.SetRandomMode(params.Speed(v), params.Brightness(v),
					params.ColorNum(v), params.Reactive(v), params.Save(v))
			})
		},
	}

	params.AddSpeed(randomModeCmd, v)
	params.AddBrightness(randomModeCmd, v)
	params.AddColorNum(randomModeCmd, v)
	params.AddReactive(randomModeCmd, v)
	params.AddSave(randomModeCmd, v)
	params.AddReset(randomModeCmd, v)

	return randomModeCmd
}
