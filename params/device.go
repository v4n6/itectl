package params

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// device related properties default values.
const (
	// DeviceBusDefault - default value of device usb bus.
	DeviceBusDefault = 0
	// deviceBusDefault - default value of device usb address.
	DeviceAddressDefault = 0
)

// device related properties and flags names.
const (
	// deviceBusProp - name of the device bus configuration property.
	deviceBusProp = "device.bus"
	// DeviceBusFlag - name of the device bus flag.
	DeviceBusFlag = "device-bus"

	// deviceAddressProp - name of the device address configuration property.
	deviceAddressProp = "device.address"
	// DeviceAddressFlag - name of the device address flag.
	DeviceAddressFlag = "device-address"
)

// AddDevice adds device related flags to the provided cmd. It also
// adds hook to bind them to the corresponding viper config
// properties.
func AddDevice(cmd *cobra.Command, v *viper.Viper) {

	cmd.PersistentFlags().Uint(DeviceBusFlag, DeviceBusDefault,
		"Bus number of the keyboard backlight device. "+configurationWarning)
	bindAndValidate(cmd, v, DeviceBusFlag, deviceBusProp, nil)

	cmd.PersistentFlags().Uint(DeviceAddressFlag, DeviceAddressDefault,
		"Address of the keyboard backlight device. "+configurationWarning)
	bindAndValidate(cmd, v, DeviceAddressFlag, deviceAddressProp, nil)
}

// Device returns device related property values: useDevice - whether
// a specific device identified by device bus and number should be
// used; deviceBus, deviceAddress - device bus and number properties
// respectively. It also validates device bus and number. It ensures
// that either both are positive, or non-positive.
func Device(v *viper.Viper) (useDevice bool,
	deviceBus, deviceAddress int, err error) {

	deviceBus, deviceAddress = v.GetInt(deviceBusProp), v.GetInt(deviceAddressProp)

	if deviceBus > 0 && deviceAddress > 0 {
		return true, deviceBus, deviceAddress, nil // device is set
	}

	if deviceBus == 0 {
		if deviceAddress > 0 {
			return false, 0, 0,
				fmt.Errorf("%w device bus number missing for \"--%s\" (either configured or specified explicitly)",
					ErrInvalidOptVal, DeviceBusFlag)
		}

		return false, 0, 0, nil // device is not set
	}

	return false, 0, 0,
		fmt.Errorf("%w device address missing for \"--%s\" (either configured or specified explicitly)",
			ErrInvalidOptVal, DeviceAddressFlag)
}
