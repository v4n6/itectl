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

// newRaindropModeCmd creates, initializes and returns command
// to set keyboard backlight to raindrop mode.
func newRaindropModeCmd(v *viper.Viper, call ite8291r3Ctl) *cobra.Command {

	var speed func() byte
	var brightness func() byte
	var colorNum func() byte
	var save func() bool

	// raindropModeCmd represents the raindrop-mode command
	var raindropModeCmd = &cobra.Command{
		Use:   "raindrop-mode",
		Short: "Set keyboard backlight to 'raindrop' mode.",
		Long: `Set keyboard backlight to 'raindrop' mode.
	Brightness of the mode can be provided by --brightness (-b) option.
  Speed of the mode's animation can be provided by --speed (-s) option.
  The predefined color used by the mode can be specified by its number provided by --color-num options.
  Color number '0' indicates black (none) color. Color number '8' indicates random color.
  If --save is provided the mode is saved in controller.`,

		RunE: func(cmd *cobra.Command, args []string) error {
			return call(func(ctl *ite8291.Controller) error {
				return ctl.SetRaindropMode(speed(), brightness(), colorNum(), save())
			})
		},
	}

	speed = params.AddSpeed(raindropModeCmd, v)
	brightness = params.AddBrightness(raindropModeCmd, v)
	colorNum = params.AddColorNum(raindropModeCmd, v)
	save = params.AddSave(raindropModeCmd, v)

	return raindropModeCmd
}
