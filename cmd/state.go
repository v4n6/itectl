package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/v4n6/itectl/pkg/ite8291"
)

const (
	isOnMessage  = "On"
	isOffMessage = "Off"
)

// newStateCmd creates, initializes and returns command to get and
// print keyboard backlight state.
func newStateCmd(call ite8291Ctl) *cobra.Command {

	// stateCmd represents the state command
	var stateCmd = &cobra.Command{
		Use:           "state",
		Short:         "Get the current state of the keyboard backlight.",
		Long:          `Print state of the keyboard backlight as either 'on' or 'off'.`,
		Args:          cobra.NoArgs,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return call(cmd, func(ctl *ite8291.Controller) error {
				isOn, err := ctl.State()
				if err != nil {
					return err
				}

				if isOn {
					fmt.Fprintln(cmd.OutOrStdout(), isOnMessage)
				} else {
					fmt.Fprintln(cmd.OutOrStdout(), isOffMessage)
				}

				return nil
			})
		},
	}

	return stateCmd
}
