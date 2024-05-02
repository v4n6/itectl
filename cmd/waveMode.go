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

// waveModeDescription - wave-mode command description
const waveModeDescription = "Set keyboard backlight to 'wave' mode."

// newWaveModeCmd creates, initializes and returns command
// to set keyboard backlight to 'wave' mode.
func newWaveModeCmd(v *viper.Viper, call ite8291Ctl) *cobra.Command {

	var speed func() byte
	var brightness func() byte
	var direction func() ite8291.Direction
	var save func() bool
	var optionallyResetColors ite8291Call

	var waveModeCmd = &cobra.Command{
		Use:   "wave-mode",
		Short: waveModeDescription,
		Long:  waveModeDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return call(func(ctl *ite8291.Controller) error {
				if err := optionallyResetColors(ctl); err != nil {
					return err
				}
				return ctl.SetWaveMode(speed(), brightness(), direction(), save())
			})
		},
	}

	speed = params.AddSpeed(waveModeCmd, v)
	brightness = params.AddBrightness(waveModeCmd, v)
	direction = params.AddDirection(waveModeCmd, v)
	save = params.AddSave(waveModeCmd, v)
	optionallyResetColors = params.AddReset(waveModeCmd, v)

	return waveModeCmd
}
