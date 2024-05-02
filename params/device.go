package params

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// DeviceBusDefault is default value of device usb bus
	DeviceBusDefault = 0
	// deviceBusDefault is default value of device usb address
	DeviceAddressDefault = 0
)

const (
	// deviceBusProp is name of the device bus configuration property
	deviceBusProp = "device.bus"
	// DeviceBusFlag is name of the device bus flag
	DeviceBusFlag = "device-bus"

	// deviceAddressProp is name of the device address configuration property
	deviceAddressProp = "device.address"
	// DeviceAddressFlag is name of the device address flag
	DeviceAddressFlag = "device-address"
)

// AddDevice adds device related flags to the provided cmd.
// It also adds hook to bind them to the corresponding viper config properties
// and to validate device properties. It ensures that either
// both device bus and device address properties are not set, or both are set.
// It returns functions to retrieve current useDevice, deviceBus and deviceAddress values.
func AddDevice(cmd *cobra.Command, v *viper.Viper) (useDevice func() bool, deviceBus, deviceAddress func() int) {

	var useDev bool

	cmd.PersistentFlags().Int(DeviceBusFlag, DeviceBusDefault,
		fmt.Sprintf("Bus number of the keyboard backlight device. %s", configurationWarning))
	bindAndValidate(cmd, v, DeviceBusFlag, deviceBusProp, nil)
	cmd.PersistentFlags().Int(DeviceAddressFlag, DeviceAddressDefault,
		fmt.Sprintf("Address of the keyboard backlight device. %s", configurationWarning))
	bindAndValidate(cmd, v, DeviceAddressFlag, deviceAddressProp, nil)

	addValidationHook(cmd, func() error {

		if !(v.IsSet(deviceBusProp) || cmd.Flag(DeviceBusFlag).Changed) &&
			!(v.IsSet(deviceAddressProp) || cmd.Flag(DeviceAddressFlag).Changed) {
			useDev = false
			return nil // device is not set
		}

		if v.IsSet(deviceBusProp) || cmd.Flag(DeviceBusFlag).Changed {
			if !(v.IsSet(deviceAddressProp) || cmd.Flag(DeviceAddressFlag).Changed) {
				return fmt.Errorf("%w device address missing for \"--%s\" (either configured or specified explicitly)",
					InvalidOptionValueError, DeviceAddressFlag)
			}
		} else if v.IsSet(deviceAddressProp) || cmd.Flag(DeviceAddressFlag).Changed {
			return fmt.Errorf("%w device bus number missing for \"--%s\" (either configured or specified explicitly)",
				InvalidOptionValueError, DeviceBusFlag)
		}

		useDev = true
		return nil
	})

	return func() bool { return useDev },
		func() int { return v.GetInt(deviceBusProp) },
		func() int { return v.GetInt(deviceAddressProp) }
}
