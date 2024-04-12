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

	"github.com/adrg/xdg"
	"github.com/gotmc/libusb/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/ite8291r3tool/config"
	ite8291 "github.com/v4n6/ite8291r3tool/pkg"
)

const configName = "ite8291r3tool"

const configType = "yaml"

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ite8291r3tool",
	Short: "A brief description of your application",
	Long:  `ROOOT1`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	cobra.EnableTraverseRunHooks = true

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Configuration file to use instead of xdg config files.")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	config.AddPollFlags(rootCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		for _, dir := range xdg.ConfigDirs {
			mergeConfig(dir)
		}
		mergeConfig(xdg.ConfigHome)
	}

	viper.AutomaticEnv() // read in environment variables that match
}

// mergeConfig ...
func mergeConfig(dir string) {
	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType(configType)

	v.AddConfigPath(dir)
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "Error reading configuration: %v\n", err)
		}
	}

	err = viper.GetViper().MergeConfigMap(v.AllSettings())
	if err != nil {
		// should never happen
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

type ite8291Command func(dev *libusb.Device, h *libusb.DeviceHandle) error

func findDevice() (*libusb.Device, *libusb.DeviceHandle, func(), error) {

	if config.Poll() {

		return ite8291.GetDeviceWithPolling(config.PollInterval(), config.PollTimeout())
	}

	return ite8291.GetDevice()
}

// executeCommand executes given cmd providing obtained instance of device
func executeCommand(cmd ite8291Command) error {

	dev, h, done, err := findDevice()
	if err != nil {
		return err
	}
	defer done()

	return cmd(dev, h)
}
