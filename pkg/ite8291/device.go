package ite8291

import (
	"errors"
	"fmt"
	"time"

	"github.com/gotmc/libusb/v2"
)

// ite8291r3 request type
const (
	SendControlRequestType    = byte(libusb.HostToDevice) | byte(libusb.Class) | byte(libusb.InterfaceRecipient)
	ReceiveControlRequestType = byte(libusb.DeviceToHost) | byte(libusb.Class) | byte(libusb.InterfaceRecipient)
)

// vendorID - vendor id of ite8291r3 usb device
const vendorID = 0x048D

// endpointOutDirection - out endpoint direction
const endpointOutDirection libusb.EndpointDirection = 0

// productIDs - supported product ids of ite8291r3 devices
var productIDs = map[uint16]bool{0x6004: true, 0x6006: true, 0xCE00: true}

// NoDeviceFoundError - error indicating that no ite8291r3 device found
var NoDeviceFoundError = errors.New("no ite8291r3 device found")

// UnsupportedDeviceError - error indicating that a device with given bus and address is not an ite8291r3 device
var UnsupportedDeviceError = errors.New("is not ite8291 device")

// NoOutEndpointFoundError - error indicating that no out endpoint found at ite8291r3 device
var NoOutEndpointFoundError = errors.New("no output endpoint found")

// USBDevice type provides ite8291r3 usb device
type USBDevice struct {
	*libusb.DeviceHandle

	ctx *libusb.Context
	dev *libusb.Device
}

// Close USBDevice. It closes underlying libusb.Contexst and libusb.DeviceHandle.
func (d *USBDevice) Close() error {

	err := d.DeviceHandle.Close()

	if err := d.ctx.Close(); err != nil {
		return err
	}

	return err
}

// WriteFunc type provides function intended to write bulk data to a USB device.
type WriteFunc func(p []byte) (n int, err error)

// GetBulkWrite returns WriteFunc intended to write bulk data to the ite8291r3 device.
// It is used to set color(s) of all/specific key(s) of keyboard backlight.
func (d *USBDevice) GetBulkWrite() (WriteFunc, error) {

	cfg, err := d.dev.ActiveConfigDescriptor()
	if err != nil {
		return nil, err
	}

	for _, iface := range cfg.SupportedInterfaces {
		for _, desc := range iface.InterfaceDescriptors {
			for _, ep := range desc.EndpointDescriptors {

				if ep.Direction() == endpointOutDirection {

					return func(p []byte) (n int, err error) {
						return d.BulkTransferOut(ep.EndpointAddress, p, 0)
					}, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("%w", NoOutEndpointFoundError)
}

// CheckDevice type provides function to check wether a given dev
// is a supported ite8291r3 device.
type CheckDevice func(dev *libusb.Device) (ok bool, err error)

// CheckDeviceByVendorProduct checks wether vendor and product ids
// of given dev are that of supported ite8291r3 devices.
func CheckDeviceByVendorProduct(dev *libusb.Device) (found bool, err error) {
	d, err := dev.DeviceDescriptor()
	if err != nil {
		return false, err
	}

	return d.VendorID == vendorID && productIDs[d.ProductID], nil
}

// NewCheckDeviceByBusAddress returns CheckDevice function
// that checks wether given dev has specified bus and address,
// and its vendor and product ids are that of supported ite8291r3 devices.
//
// It returns instance of UnsupportedDeviceError if a device with
// correct bus and address is not a supported ite8291r3 device.
func NewCheckDeviceByBusAddress(bus, address int) CheckDevice {
	return func(dev *libusb.Device) (bool, error) {
		b, _ := dev.BusNumber()
		a, _ := dev.DeviceAddress()
		if bus != b || address != a {
			return false, nil
		}

		ok, err := CheckDeviceByVendorProduct(dev)
		if err != nil {
			return false, err
		}

		if ok {
			return true, nil
		}

		return false, fmt.Errorf("device (bus=%d,address=%d) %w",
			bus, address, UnsupportedDeviceError)

	}
}

// targetInterfaceNumber - interface number of ite8291r3 device to check and detach, if necessary, kernel driver.
const targetInterfaceNumber = 1

// LookupDevice traverses all usb devices and returns first found supported ite8291r3 device.
// It uses given check function to decide whether a device is supported.
// It returns pointer to found usbDevice or occured error.
// LookupDevice returns instance of NoDeviceFoundError if no supported device found.
func LookupDevice(ctx *libusb.Context, check CheckDevice) (usbDevice *USBDevice, err error) {

	devs, err := ctx.DeviceList()
	if err != nil {
		return nil, err
	}

	for _, dev := range devs {

		var found bool
		found, err = check(dev)
		if errors.Is(err, UnsupportedDeviceError) {
			return nil, err // device found but it isn't an ite8291 device
		}

		if !(err == nil && found) {
			continue // either not found or errored in checker -> save error and continue
		}

		// open handler
		h, err := dev.Open()
		if err != nil {
			return nil, err
		}

		usbDevice = &USBDevice{DeviceHandle: h, dev: dev, ctx: ctx}

		// detach kernel driver
		kern, err := usbDevice.KernelDriverActive(targetInterfaceNumber)
		if err != nil {
			_ = h.Close() // close handle on error
			return nil, err
		}

		if kern {
			if err = usbDevice.DetachKernelDriver(targetInterfaceNumber); err != nil {
				_ = h.Close() // close handle on error
				return nil, err
			}
		}

		return usbDevice, nil
	}

	if err != nil {
		// report not found error together with saved error
		return nil, fmt.Errorf("%w: %w", NoDeviceFoundError, err)
	}

	return nil, fmt.Errorf("%w", NoDeviceFoundError)
}

// FindDevice searches for a supported ite8291r3 device. check function decides whether a device is a supported one.
//
// If no device was found it repeats the search after pollInterval duration. If no device was found in timeout
// duration, FindDevice returns instance of NoDeviceFoundError. If timeout is 0 or negative no further searches
// are done and the error is returned immediately.
func FindDevice(pollInterval, timeout time.Duration, check CheckDevice) (usbDevice *USBDevice, err error) {

	ctx, err := libusb.NewContext()
	if err != nil {
		return nil, err
	}

	if timeout <= 0 {
		if usbDevice, err = LookupDevice(ctx, check); err != nil {
			_ = ctx.Close()
		}
		return usbDevice, err
	}

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	exit := make(chan bool)
	go func() {
		time.Sleep(timeout)
		exit <- true
	}()

	for {

		if usbDevice, err = LookupDevice(ctx, check); err == nil {
			return usbDevice, err
		}

		if err != nil && !errors.Is(err, NoDeviceFoundError) {
			_ = ctx.Close()
			return nil, err
		}

		select {

		case <-exit:
			_ = ctx.Close()
			return nil, err

		case <-ticker.C:
		}
	}
}
