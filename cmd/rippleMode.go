package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/params"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// rippleModeDescription - ripple-mode command description.
const rippleModeDescription = "Set keyboard backlight to 'ripple' mode."

// newRippleModeCmd creates, initializes and returns command
// to set keyboard backlight to 'ripple' mode.
func newRippleModeCmd(v *viper.Viper, call ite8291Ctl) *cobra.Command {

	var rippleModeCmd = &cobra.Command{
		Use:           "ripple-mode",
		Short:         rippleModeDescription,
		Long:          rippleModeDescription,
		Args:          cobra.NoArgs,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return call(cmd, func(ctl *ite8291.Controller) error {
				return ctl.SetRippleMode(params.Speed(v), params.Brightness(v),
					params.ColorNum(v), params.Reactive(v), params.Save(v))
			})
		},
	}

	params.AddSpeed(rippleModeCmd, v)
	params.AddBrightness(rippleModeCmd, v)
	params.AddColorNum(rippleModeCmd, v)
	params.AddReactive(rippleModeCmd, v)
	params.AddSave(rippleModeCmd, v)
	params.AddReset(rippleModeCmd, v)

	return rippleModeCmd
}
