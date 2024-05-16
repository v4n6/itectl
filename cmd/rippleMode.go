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
	"github.com/v4n6/itectl/params"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// rippleModeDescription - ripple-mode command description
const rippleModeDescription = "Set keyboard backlight to 'ripple' mode."

// newRippleModeCmd creates, initializes and returns command
// to set keyboard backlight to 'ripple' mode.
func newRippleModeCmd(v *viper.Viper, call ite8291Ctl) *cobra.Command {

	var speed func() byte
	var brightness func() byte
	var colorNum func() byte
	var reactive func() bool
	var save func() bool
	var optionallyResetColors ite8291Call

	var rippleModeCmd = &cobra.Command{
		Use:   "ripple-mode",
		Short: rippleModeDescription,
		Long:  rippleModeDescription,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return call(func(ctl *ite8291.Controller) error {
				if err := optionallyResetColors(ctl); err != nil {
					return err
				}
				return ctl.SetRippleMode(speed(), brightness(), colorNum(), reactive(), save())
			})
		},
	}

	speed = params.AddSpeed(rippleModeCmd, v)
	brightness = params.AddBrightness(rippleModeCmd, v)
	colorNum = params.AddColorNum(rippleModeCmd, v)
	reactive = params.AddReactive(rippleModeCmd, v)
	save = params.AddSave(rippleModeCmd, v)
	optionallyResetColors = params.AddReset(rippleModeCmd, v)

	return rippleModeCmd
}
