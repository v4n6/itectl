package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/params"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// setBrightnessDescription - set-brightness command description.
const setBrightnessDescription = "Set keyboard backlight brightness."

// newSetBrightnessCmd creates, initializes and returns command
// to set keyboard backlight brightness.
func newSetBrightnessCmd(v *viper.Viper, call ite8291Ctl) *cobra.Command {

	var setBrightnessCmd = &cobra.Command{
		Use:           "set-brightness",
		Short:         setBrightnessDescription,
		Long:          setBrightnessDescription,
		Args:          cobra.NoArgs,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return call(cmd, func(ctl *ite8291.Controller) error {
				return ctl.SetBrightness(params.Brightness(v))
			})
		},
	}

	params.AddBrightness(setBrightnessCmd, v)

	if err := setBrightnessCmd.MarkPersistentFlagRequired("brightness"); err != nil {
		panic(err)
	}

	return setBrightnessCmd
}
