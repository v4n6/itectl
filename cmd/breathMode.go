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

// breathModeDescription - breath-mode command description.
const breathModeDescription = "Set keyboard backlight to 'breathing' mode."

// newBreathModeCmd creates, initializes and returns command
// to set keyboard backlight to 'breathing' mode.
func newBreathModeCmd(v *viper.Viper, call ite8291Ctl) *cobra.Command {

	breathModeCmd := &cobra.Command{
		Use:           "breath-mode",
		Short:         breathModeDescription,
		Long:          breathModeDescription,
		Args:          cobra.NoArgs,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return call(cmd, func(ctl *ite8291.Controller) error {
				return ctl.SetBreathingMode(params.Speed(v), params.Brightness(v),
					params.ColorNum(v), params.Save(v))
			})
		},
	}

	params.AddSpeed(breathModeCmd, v)
	params.AddBrightness(breathModeCmd, v)
	params.AddColorNum(breathModeCmd, v)
	params.AddSave(breathModeCmd, v)
	params.AddReset(breathModeCmd, v)

	return breathModeCmd
}
