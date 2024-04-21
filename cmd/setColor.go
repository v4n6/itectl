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

// newSetColorCmd creates, initializes and returns command
// to set keyboard backlight predefined color.
func newSetColorCmd(v *viper.Viper, call ite8291r3Ctl) *cobra.Command {

	var colorNum func() byte
	var color func() *ite8291.Color

	// setColorCmd represents the set-color command
	var setColorCmd = &cobra.Command{
		Use:   "set-color",
		Short: "Set keyboard backlight predefined color.",
		Long: `Set a predefined color of the keyboard backlight to the provided RGB value.
  The predefined color must be specified by its number provided by --color-num options.
  The RGB value must be provided by one or more of its atoms using --red, --green, --blue options.`,

		RunE: func(cmd *cobra.Command, args []string) error {
			return call(func(ctl *ite8291.Controller) error {
				return ctl.SetColor(colorNum(), color())
			})
		},
	}

	colorNum = params.AddAssignableColorNum(setColorCmd, v)
	color = params.AddColor(setColorCmd, v)

	return setColorCmd
}
