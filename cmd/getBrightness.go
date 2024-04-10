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

	"github.com/gotmc/libusb/v2"
	"github.com/spf13/cobra"
	ite8291 "github.com/v4n6/ite8291r3tool/pkg"
)

// getBrightnessCmd represents the get-brightness command
var getBrightnessCmd = &cobra.Command{
	Use:   "get-brightness",
	Short: "Get current brightness of the keyboard backlight.",
	Long:  `Print current brightness of the keyboard backlight.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		return executeCommand(func(dev *libusb.Device, h *libusb.DeviceHandle) error {

			brightness, err := ite8291.GetBrightness(h)
			if err != nil {
				return err
			}

			fmt.Printf("%d", brightness)
			return nil
		})
	},
}

func init() {
	rootCmd.AddCommand(getBrightnessCmd)
}
