package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// brightnessDescription - brightness command description.
const brightnessDescription = "Retrieve and print current brightness of the keyboard backlight."

// newBrightnessCmd creates, initializes and returns command
// to retrieve and print keyboard backlight brightness.
func newBrightnessCmd(call ite8291Ctl) *cobra.Command {

	return &cobra.Command{
		Use:           "brightness",
		Short:         brightnessDescription,
		Long:          brightnessDescription,
		Args:          cobra.NoArgs,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return call(cmd, func(ctl *ite8291.Controller) error {
				brightness, err := ctl.Brightness()
				if err != nil {
					return err
				}

				fmt.Fprintf(cmd.OutOrStdout(), "%d\n", brightness)
				return nil
			})
		},
	}
}
