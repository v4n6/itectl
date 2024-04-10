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
	"github.com/gotmc/libusb/v2"
	"github.com/spf13/cobra"
	"github.com/v4n6/ite8291r3tool/config"
	ite8291 "github.com/v4n6/ite8291r3tool/pkg"
)

// setBrightnessCmd represents the set-brightness command
var setBrightnessCmd = &cobra.Command{
	Use:   "set-brightness",
	Short: "Change keyboard backlight brightness.",
	Long: `Set brightness of the keyboard backlight to the provided value
  Brightness of the mode must be specified by --brightness (-b) option.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		brightness, err := config.Config.BrightnessVal()
		if err != nil {
			return err
		}

		return executeCommand(func(dev *libusb.Device, h *libusb.DeviceHandle) error {
			return ite8291.SetBrightness(h, brightness)
		})
	},
}

func init() {
	rootCmd.AddCommand(setBrightnessCmd)

	config.AddBrightnessFlag(setBrightnessCmd)
	_ = setBrightnessCmd.MarkPersistentFlagRequired("brightness")
}
