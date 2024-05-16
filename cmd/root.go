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
	"io"
	"os"
	"slices"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/params"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// Execute invokes the application.
func Execute() {

	_ = ExecuteCmd(os.Args[1:], os.Stdout, os.Stderr, findIteDevice, params.ReadConfig)
}

// ExecuteCmd invokes the command provided by args or sets keyboard backlight to a configured mode.
// output, errOut provide corresponding output and error streams.
// v specifies viper to use. find function is used to look up a supported ite8291 device.
// readConfig function is used to retrieve configuration either from configuration file provided
// by corresponding flag or from default global and/or user configuration files.
func ExecuteCmd(args []string, output, errOut io.Writer,
	find findDevice, readConf readConfig) (err error) {

	cobra.EnableTraverseRunHooks = true

	v := viper.New()
	rootCmd := newRootCmd(v, find) // root command
	rootCmd.InitDefaultCompletionCmd()
	rootCmd.InitDefaultHelpCmd()
	rootCmd.InitDefaultHelpFlag()

	cmd, flags, err := rootCmd.Traverse(args) // get sub-command
	if err != nil {
		cobra.CheckErr(err)
	}

	cfgFile, err := params.ConfigFile(rootCmd, flags)
	if err != nil {
		cobra.CheckErr(err)
	}

	if err = readConf(rootCmd, v, cfgFile); err != nil {
		return err
	}

	if cmd.Use == rootCmd.Use &&
		!cmd.Flag("help").Changed &&
		!slices.ContainsFunc(flags, func(s string) bool {
			return s == cobra.ShellCompRequestCmd ||
				s == cobra.ShellCompNoDescRequestCmd
		}) {

		// no sub-command provided
		defaultMode := params.DefaultMode(v) // configured default mode
		if len(defaultMode) > 0 {
			// insert default mode command
			args = slices.Insert(args, 0, fmt.Sprintf("%s-mode", defaultMode))
		}
	}

	rootCmd.SetArgs(args)
	rootCmd.SetOut(output)
	rootCmd.SetErr(errOut)
	return rootCmd.Execute()
}

// readConfig type provides a function to retrieve and merge viper configuration.
type readConfig func(cmd *cobra.Command, v *viper.Viper, cfgFile string) error

// ite8291Call type provides a function that calls a sub command method on provided controller.
type ite8291Call func(ctl *ite8291.Controller) error

// ite8291Ctl type defines a function that provides ite8291r3 controller and calls given f with it.
type ite8291Ctl func(f ite8291Call) error

// findDevice type provides a function that looks up a supported ite8291r3 device based on the given parameters.
// It returns pointer to found device or occurred error.
type findDevice func(useDevice bool, bus, address int,
	pollInterval, timeout time.Duration) (dev ite8291.Device, err error)

// findIteDevice looks up a supported ite8291r3 device.
//
// useDevice specifies whether an ite8291r3 device identified by bus and address must be used.
// timeout specifies maximum duration to wait till a supported device can be found.
// If timeout is 0 or negative, it doesn't wait and return the corresponding error immediately.
// pollInterval specifies duration to wait between consequent search attempts.
func findIteDevice(useDevice bool, bus, address int,
	pollInterval, pollTimeout time.Duration) (dev ite8291.Device, err error) {

	devChecker := ite8291.CheckDeviceByVendorProduct
	if useDevice {
		devChecker = ite8291.NewCheckDeviceByBusAddress(bus, address)
	}

	return ite8291.FindDevice(pollInterval, pollTimeout, devChecker)
}

// newRootCmd creates, initializes and returns root command.
// v is a viper instance used by commands instead of static one.
// find is a findDevice function used to obtain ite8291r3 device instance.
// readConf is a function used to retrieve viper configuration.
func newRootCmd(v *viper.Viper, find findDevice) *cobra.Command {

	var rootCmd = &cobra.Command{
		Use:               "itectl",
		Short:             "Control ITE 8291 keyboard backlight",
		Long:              `Control ITE 8291 keyboard backlight`,
		TraverseChildren:  true,
		CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
	}

	// config flag
	params.AddConfigFlag(rootCmd)

	// poll related properties
	pollInterval, pollTimeout := params.AddPoll(rootCmd, v)
	// device related properties
	useDevive, deviceBus, deviceAddress := params.AddDevice(rootCmd, v)

	// ite8291Ctl
	exec := func(f ite8291Call) error {

		dev, err := find(useDevive(), deviceBus(), deviceAddress(), pollInterval(), pollTimeout())
		if err != nil {
			return err
		}

		ctl := ite8291.NewController(dev)
		defer ctl.Close()

		return f(ctl)
	}

	// build commands hierarchy
	rootCmd.AddCommand(newOffModeCmd(v, exec))
	rootCmd.AddCommand(newAuroraModeCmd(v, exec))
	rootCmd.AddCommand(newBreathModeCmd(v, exec))
	rootCmd.AddCommand(newFireworksModeCmd(v, exec))
	rootCmd.AddCommand(newMarqueeModeCmd(v, exec))
	rootCmd.AddCommand(newRainbowModeCmd(v, exec))
	rootCmd.AddCommand(newRaindropModeCmd(v, exec))
	rootCmd.AddCommand(newRandomModeCmd(v, exec))
	rootCmd.AddCommand(newRippleModeCmd(v, exec))
	rootCmd.AddCommand(newSingleColorModeCmd(v, exec))
	rootCmd.AddCommand(newWaveModeCmd(v, exec))
	rootCmd.AddCommand(newBrightnessCmd(v, exec))
	rootCmd.AddCommand(newSetBrightnessCmd(v, exec))
	rootCmd.AddCommand(newFirmwareVersionCmd(v, exec))
	rootCmd.AddCommand(newStateCmd(v, exec))
	rootCmd.AddCommand(newSetColorCmd(v, exec))

	return rootCmd
}
