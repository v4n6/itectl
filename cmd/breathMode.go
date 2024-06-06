package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/params"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// breathModeDescription - breath-mode command description.
const breathModeDescription = "Set keyboard backlight to 'breathing' mode."

// newBreathModeCmd creates, initializes and returns command
// to set keyboard backlight to 'breathing' mode.
func newBreathModeCmd(v *viper.Viper, call ite8291Ctl) *cobra.Command {

	breathModeCmd := &cobra.Command{
		Use:           "breath-mode",
		Short:         breathModeDescription,
		Long:          breathModeDescription,
		Args:          cobra.NoArgs,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return call(cmd, func(ctl *ite8291.Controller) error {
				return ctl.SetBreathingMode(params.Speed(v), params.Brightness(v),
					params.ColorNum(v), params.Save(v))
			})
		},
	}

	params.AddSpeed(breathModeCmd, v)
	params.AddBrightness(breathModeCmd, v)
	params.AddColorNum(breathModeCmd, v)
	params.AddSave(breathModeCmd, v)
	params.AddReset(breathModeCmd, v)

	return breathModeCmd
}
