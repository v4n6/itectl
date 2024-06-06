package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// firmwareVersionDescription - firmware-version command description.
const firmwareVersionDescription = "Retrieve and print firmware version of the keyboard backlight controller."

// newFirmwareVersionCmd creates, initializes and returns command
// to retrieve and print keyboard backlight controller firmware version.
func newFirmwareVersionCmd(call ite8291Ctl) *cobra.Command {

	return &cobra.Command{
		Use:           "firmware-version",
		Short:         firmwareVersionDescription,
		Long:          firmwareVersionDescription,
		Args:          cobra.NoArgs,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return call(cmd, func(ctl *ite8291.Controller) error {
				ver, err := ctl.FirmwareVersion()
				if err != nil {
					return err
				}

				fmt.Fprintln(cmd.OutOrStdout(), ver)
				return nil
			})
		},
	}
}
