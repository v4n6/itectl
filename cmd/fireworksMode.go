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

// fireworksModeCmd represents the fireworks-mode command
var fireworksModeCmd = &cobra.Command{
	Use:   "fireworks-mode",
	Short: "Set keyboard backlight to 'fireworks' mode.",
	Long: `Set keyboard backlight to 'fireworks' mode.
  Brightness of the mode can be provided by --brightness (-b) option.
  Speed of the mode's animation can be provided by --speed (-s) option.
  The predefined color used by the mode can be specified by its number provided by --color-num options.
  Color number '0' indicates black (none) color. Color number '8' indicates random color.
  If --reactive is provided the backlight reacts to user input.
  If --save is provided the mode is saved in controller.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		return executeCommand(func(dev *libusb.Device, h *libusb.DeviceHandle) error {
			return ite8291.SetFireworksMode(h, config.Speed(), config.Brightness(), config.ColorNum(),
				config.Reactive(), config.Save())
		})
	},
}

func init() {
	rootCmd.AddCommand(fireworksModeCmd)

	config.AddSpeedFlag(fireworksModeCmd)
	config.AddBrightnessFlag(fireworksModeCmd)
	config.AddColorNumFlag(fireworksModeCmd)

	config.AddReactiveFlag(fireworksModeCmd)
	config.AddSaveFlag(fireworksModeCmd)
}
