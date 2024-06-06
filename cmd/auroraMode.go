package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/params"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// auroraModeDescription - aurora-mode command description.
const auroraModeDescription = "Set keyboard backlight to 'aurora' mode."

// newAuroraModeCmd creates, initializes and returns command
// to set keyboard backlight to 'aurora' mode.
func newAuroraModeCmd(v *viper.Viper, call ite8291Ctl) *cobra.Command {

	auroraModeCmd := &cobra.Command{
		Use:           "aurora-mode",
		Short:         auroraModeDescription,
		Long:          auroraModeDescription,
		Args:          cobra.NoArgs,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return call(cmd, func(ctl *ite8291.Controller) error {
				return ctl.SetAuroraMode(params.Speed(v), params.Brightness(v),
					params.ColorNum(v), params.Reactive(v), params.Save(v))
			})
		},
	}

	params.AddSpeed(auroraModeCmd, v)
	params.AddBrightness(auroraModeCmd, v)
	params.AddColorNum(auroraModeCmd, v)
	params.AddReactive(auroraModeCmd, v)
	params.AddSave(auroraModeCmd, v)
	params.AddReset(auroraModeCmd, v)

	return auroraModeCmd
}
