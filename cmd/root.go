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
	"os"
	"slices"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/ite8291r3tool/params"
	"github.com/v4n6/ite8291r3tool/pkg/ite8291"
)

// Execute executes the application
func Execute() {

	cobra.EnableTraverseRunHooks = true
	args := os.Args[1:]

	v := viper.New()

	rootCmd := newRootCmd(v, findIte8291r3) // root command
	cmd, _, err := rootCmd.Find(args)       // get sub-command
	if err != nil {
		cobra.CheckErr(err)
	}

	// read provided config file or read & merge default configuration files
	if err := params.ReadConfig(cmd, v, args); err != nil {
		cobra.CheckErr(err)
	}

	v.AutomaticEnv() // read in environment variables that match

	if cmd.Use == rootCmd.Use {
		// no sub-command provided
		defaultMode := params.DefaultMode(v) // configured default mode
		if len(defaultMode) > 0 {
			// insert default mode command
			rootCmd.SetArgs(slices.Insert(args, 0, fmt.Sprintf("%s-mode", defaultMode)))
		}
	}

	if err := rootCmd.Execute(); err != nil {
		cobra.CheckErr(err)
	}
}

// ite8291r3Call is a function that calls a method on provided controller
type ite8291r3Call func(ctl *ite8291.Controller) error

// ite8291r3Ctl is a function that provides ite8291r3 controller and calls given f on it
type ite8291r3Ctl func(f ite8291r3Call) error

// findDevice is a function that finds ite8291r3 device based om specified parameters
// useDevice specifies whether a specific device identified by bus and address must be used
// poll specifies whether function must poll for device presence
// in given pollInterval intervals with given pollTimeout timeout
type findDevice func(useDevice bool, bus, address int, poll bool, pollInterval, pollTimeout time.Duration) (dev ite8291.Device, err error)

// findIte8291r3 finds implements findDevice function.
// It finds ite8291r3 device based om specified parameters
// useDevice specifies whether a specific device identified by bus and address must be used
// poll specifies whether function must poll for device presence
// in given pollInterval intervals with given pollTimeout timeout
func findIte8291r3(useDevice bool, bus, address int,
	poll bool, pollInterval, pollTimeout time.Duration) (dev ite8291.Device, err error) {

	devChecker := ite8291.VendorProductDeviceCheckerFunc
	if useDevice {
		devChecker = ite8291.NewAddressDeviceCheckerFunc(bus, address)
	}

	devFinder := ite8291.FindDeviceWithoutPollingFunc
	if poll {
		devFinder = ite8291.NewFindDeviceWithPollingFunc(pollInterval, pollTimeout)
	}

	return ite8291.GetDevice(devFinder, devChecker)
}

// newRootCmd creates and returns initialized root command.
// v is a viper instance used by commands instead of static one.
// find is findDevice function used to obtain ite8291r3 device instance.
func newRootCmd(v *viper.Viper, find findDevice) *cobra.Command {

	var cfgFile string // viper config file provided by user

	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:              "ite8291r3ctl",
		Short:            "Control ite8291r3 keyboard backlight",
		Long:             `Control ite8291r3 keyboard backlight`,
		TraverseChildren: true,
	}

	// add --config flag to override viper config file discovery
	rootCmd.PersistentFlags().StringVar(&cfgFile, params.ConfigFileFlag, "",
		fmt.Sprintf("Configuration file to use instead of xdg config files [\"/etc/%[1]s.%[2]s\",\"~/.config/%[1]s.%[2]s\"]",
			params.ConfigName, params.ConfigType))
	// add poll relative parameters
	poll, pollInterval, pollTimeout := params.AddPoll(rootCmd, v)
	// add reset parameters
	reset, predefinedColors := params.AddReset(rootCmd, v)
	// add deice relative parameters
	useDevive, deviceBus, deviceAddress := params.AddDevice(rootCmd, v)

	// function that obtains ite8291r3 controller and calls given f on it
	exec := func(f ite8291r3Call) error {

		dev, err := find(useDevive(), deviceBus(), deviceAddress(), poll(), pollInterval(), pollTimeout())
		if err != nil {
			return err
		}

		ctl := ite8291.NewController(dev)
		defer ctl.Close()

		if reset() {
			if err := ctl.SetColors(predefinedColors()); err != nil {
				return err
			}
		}

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
	rootCmd.AddCommand(newGetBrightnessCmd(v, exec))
	rootCmd.AddCommand(newSetBrightnessCmd(v, exec))
	rootCmd.AddCommand(newGetFirmwareVersionCmd(v, exec))
	rootCmd.AddCommand(newStateCmd(v, exec))
	rootCmd.AddCommand(newSetColorCmd(v, exec))
	rootCmd.AddCommand(newListDevicesCmd(v, exec))

	return rootCmd
}
