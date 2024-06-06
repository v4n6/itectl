package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/params"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// offModeDescription - off-mode command description.
const offModeDescription = "Turn the keyboard backlight off."

// newOffModeCmd creates, initializes and returns command
// to set keyboard backlight off.
func newOffModeCmd(v *viper.Viper, call ite8291Ctl) *cobra.Command {

	offModeCmd := &cobra.Command{
		Use:           "off-mode",
		Short:         offModeDescription,
		Long:          offModeDescription,
		Args:          cobra.NoArgs,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return call(cmd, func(ctl *ite8291.Controller) error {
				return ctl.SetOffMode()
			})
		},
	}

	params.AddReset(offModeCmd, v)

	return offModeCmd
}
