package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/params"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// marqueeModeDescription - marquee-mode command description.
const marqueeModeDescription = "Set keyboard backlight to 'marquee' mode."

// newMarqueeModeCmd creates, initializes and returns command
// to set keyboard backlight to 'marquee' mode.
func newMarqueeModeCmd(v *viper.Viper, call ite8291Ctl) *cobra.Command {

	var marqueeModeCmd = &cobra.Command{
		Use:           "marquee-mode",
		Short:         marqueeModeDescription,
		Long:          marqueeModeDescription,
		Args:          cobra.NoArgs,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return call(cmd, func(ctl *ite8291.Controller) error {
				return ctl.SetMarqueeMode(params.Speed(v), params.Brightness(v),
					params.Save(v))
			})
		},
	}

	params.AddSpeed(marqueeModeCmd, v)
	params.AddBrightness(marqueeModeCmd, v)
	params.AddSave(marqueeModeCmd, v)
	params.AddReset(marqueeModeCmd, v)

	return marqueeModeCmd
}
