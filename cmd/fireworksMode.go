package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/params"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// fireworksModeDescription - fireworks-mode command description.
const fireworksModeDescription = "Set keyboard backlight to 'fireworks' mode."

// newFireworksModeCmd creates, initializes and returns command
// to set keyboard backlight to 'fireworks' mode.
func newFireworksModeCmd(v *viper.Viper, call ite8291Ctl) *cobra.Command {

	fireworksModeCmd := &cobra.Command{
		Use:           "fireworks-mode",
		Short:         fireworksModeDescription,
		Long:          fireworksModeDescription,
		Args:          cobra.NoArgs,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return call(cmd, func(ctl *ite8291.Controller) error {
				return ctl.SetFireworksMode(params.Speed(v), params.Brightness(v),
					params.ColorNum(v), params.Reactive(v), params.Save(v))
			})
		},
	}

	params.AddSpeed(fireworksModeCmd, v)
	params.AddBrightness(fireworksModeCmd, v)
	params.AddColorNum(fireworksModeCmd, v)
	params.AddReactive(fireworksModeCmd, v)
	params.AddSave(fireworksModeCmd, v)
	params.AddReset(fireworksModeCmd, v)

	return fireworksModeCmd
}
