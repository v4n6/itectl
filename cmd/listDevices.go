/*
Copyright © 2024 Sergey Morozov

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/ite8291r3tool/pkg/ite8291"
)

const deviceOutputTemplate = "Bus:{{printf \"%03d\" .BusNumber}} Device:{{printf \"%03d\" .DeviceAddress}} Port:{{printf \"%03d\" .PortNumber}} Vendor:{{printf \"%04x\" .VendorID}} Product:{{printf \"%04x\" .ProductID}} Rev:{{ .DeviceReleaseNumber}}"

// newListDevicesCmd creates, initializes and returns command
// to list all available supported devices.
func newListDevicesCmd(v *viper.Viper, call ite8291r3Ctl) *cobra.Command {

	// listDevicesCmd represents the list-devices command
	return &cobra.Command{
		Use:   "list-devices",
		Short: "List supported devices",
		Long:  `List supported devices`,
		RunE: func(cmd *cobra.Command, args []string) error {

			descs, err := ite8291.ListDevices()
			if err != nil {
				return err
			}

			tmpl, err := template.New("device").Parse(deviceOutputTemplate)
			if err != nil {
				return err
			}

			if len(descs) == 0 {
				fmt.Println("No supported devices found")
				return nil
			}

			for _, d := range descs {

				err = tmpl.Execute(os.Stdout, d)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}
