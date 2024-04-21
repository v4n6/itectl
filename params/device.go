package params

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// deviceBusDefault is default value of device usb bus
	deviceBusDefault = 0
	// deviceBusDefault is default value of device usb address
	deviceAddressDefault = 0
)

const (
	// deviceBusProp is name of the device bus configuration property
	deviceBusProp = "device.bus"
	// deviceBusFlag is name of the device bus flag
	deviceBusFlag = "device-bus"

	// deviceAddressProp is name of the device address configuration property
	deviceAddressProp = "device.address"
	// deviceAddressFlag is name of the device address flag
	deviceAddressFlag = "device-address"
)

// AddDevice adds device related flags to the provided cmd.
// It also adds hook to bind them to the corresponding viper config properties
// and to validate device properties. It ensures that either
// both device bus and device address properties are not set, or both are set.
// It returns functions to retrieve current device, deviceBus and deviceAddress values.
func AddDevice(cmd *cobra.Command, v *viper.Viper) (useDevice func() bool, deviceBus, deviceAddress func() int) {

	var device bool

	cmd.PersistentFlags().Int(deviceBusFlag, deviceBusDefault, "Bus number of the keyboard backlight device.")
	bindAndValidate(cmd, v, deviceBusFlag, deviceBusProp, nil)
	cmd.PersistentFlags().Int(deviceAddressFlag, deviceAddressDefault, "Address of the keyboard backlight device.")
	bindAndValidate(cmd, v, deviceAddressFlag, deviceAddressProp, nil)

	addValidationHook(cmd, func() error {

		if !(v.IsSet(deviceBusProp) || cmd.Flag(deviceBusFlag).Changed) &&
			!(v.IsSet(deviceAddressProp) || cmd.Flag(deviceAddressFlag).Changed) {

			return nil
		}

		if v.IsSet(deviceBusProp) || cmd.Flag(deviceBusFlag).Changed {
			if !(v.IsSet(deviceAddressProp) || cmd.Flag(deviceAddressFlag).Changed) {
				return fmt.Errorf("%w: missing device adddress (either configured or specified explicitly via --%s flag)",
					InvalidOptionValueError, deviceAddressFlag)
			}
		} else if v.IsSet(deviceAddressProp) || cmd.Flag(deviceAddressFlag).Changed {
			return fmt.Errorf("%w: missing device bus number (either configured or specified explicitly via --%s flag)",
				InvalidOptionValueError, deviceBusFlag)
		}

		device = true
		return nil
	})

	return func() bool { return device },
		func() int { return v.GetInt(deviceBusProp) },
		func() int { return v.GetInt(deviceAddressProp) }
}
