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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/ite8291r3tool/params"
	"github.com/v4n6/ite8291r3tool/pkg/ite8291"
)

// newSingleColorModeCmd creates, initializes and returns command
// to set keyboard backlight to single color mode.
func newSingleColorModeCmd(v *viper.Viper, call ite8291r3Ctl) *cobra.Command {

	var brightness func() byte
	var color func() *ite8291.Color
	var save func() bool

	// singleColorModeCmd represents the single-color-mode command
	var singleColorModeCmd = &cobra.Command{
		Use:   "single-color-mode",
		Short: "Set keyboard backlight to single color",
		Long: `Set keyboard backlight to the given color.
  The color is specified by its RGB atoms using --red, --green, --blue options.
  Brightness can be provided by --brightness (-b) option.
  If --save is provided the corresponding colors are saved in controller.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return call(func(ctl *ite8291.Controller) error {
				return ctl.SetSingleColorMode(brightness(), color(), save())
			})
		},
	}

	brightness = params.AddBrightness(singleColorModeCmd, v)
	color = params.AddSingleColor(singleColorModeCmd, v)
	save = params.AddSave(singleColorModeCmd, v)

	return singleColorModeCmd
}
