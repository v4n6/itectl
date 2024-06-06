package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/params"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// waveModeDescription - wave-mode command description.
const waveModeDescription = "Set keyboard backlight to 'wave' mode."

// newWaveModeCmd creates, initializes and returns command to set
// keyboard backlight to 'wave' mode.
func newWaveModeCmd(v *viper.Viper, call ite8291Ctl) *cobra.Command {

	var direction func() ite8291.Direction

	var waveModeCmd = &cobra.Command{
		Use:           "wave-mode",
		Short:         waveModeDescription,
		Long:          waveModeDescription,
		Args:          cobra.NoArgs,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return call(cmd, func(ctl *ite8291.Controller) error {
				return ctl.SetWaveMode(params.Speed(v), params.Brightness(v),
					direction(), params.Save(v))
			})
		},
	}

	params.AddSpeed(waveModeCmd, v)
	params.AddBrightness(waveModeCmd, v)
	direction = params.AddDirection(waveModeCmd, v)
	params.AddSave(waveModeCmd, v)
	params.AddReset(waveModeCmd, v)

	return waveModeCmd
}
