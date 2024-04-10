package ite8291

import (
	"errors"
	"fmt"
	"time"

	"github.com/gotmc/libusb/v2"
)

const vendorID = 0x048D

var (
	productIDs = map[uint16]bool{0x6004: true, 0x6006: true, 0xCE00: true}
)

const targetInterfaceNumber = 1

const endpointOutDirection libusb.EndpointDirection = 0

var NoDeviceFoundError = errors.New("no ite8291 device found")

var UnsupportedDeviceError = errors.New("is not ite8291 device")

var NoOutEndpointFoundError = errors.New("no output endpoint found")

type Context interface {
	DeviceList() ([]*libusb.Device, error)
}

type deviceChecker func(dev *libusb.Device) (bool, error)

func vendorProductDeviceChecker(dev *libusb.Device) (bool, error) {
	if d, err := dev.DeviceDescriptor(); err == nil &&
		d.VendorID == vendorID && productIDs[d.ProductID] {

		return true, nil
	}

	return false, nil
}

func getAddressDeviceChecker(bus, address int) deviceChecker {
	return func(dev *libusb.Device) (bool, error) {
		b, _ := dev.BusNumber()
		a, _ := dev.DeviceAddress()
		if bus == b && address == a {

			d, err := dev.DeviceDescriptor()
			if err != nil {
				return false, err
			}

			if d.VendorID != vendorID || !productIDs[d.ProductID] {
				return false, fmt.Errorf("device (bus=%d,address=%d) %w", bus, address, UnsupportedDeviceError)
			}

			return true, nil
		}

		return false, nil
	}
}

type findDevice func(ctx Context, checker deviceChecker) (device *libusb.Device, handle *libusb.DeviceHandle, err error)

func findDeviceImpl(ctx Context, checker deviceChecker) (device *libusb.Device, handle *libusb.DeviceHandle, err error) {

	devs, err := ctx.DeviceList()
	if err != nil {
		return nil, nil, err
	}

	for _, dev := range devs {
		match, err := checker(dev)
		if err != nil {
			return nil, nil, err
		}
		if match {
			h, err := dev.Open()
			if err != nil {
				return nil, nil, err
			}

			kern, err := h.KernelDriverActive(targetInterfaceNumber)
			if err != nil {
				return nil, nil, err
			}

			if kern {
				err = h.DetachKernelDriver(targetInterfaceNumber)
				if err != nil {
					return nil, nil, err
				}
			}

			return dev, h, nil
		}
	}

	return nil, nil, fmt.Errorf("%w", NoDeviceFoundError)
}

func getFindDeviceWithPollingImpl(pollInterval, timeout time.Duration) findDevice {

	return func(ctx Context, checker deviceChecker) (device *libusb.Device, handle *libusb.DeviceHandle, err error) {

		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()

		exit := make(chan bool)
		go func() {
			time.Sleep(timeout)
			exit <- true
		}()

		for {

			dev, h, err := findDeviceImpl(ctx, checker)
			if err != nil && !errors.Is(err, NoDeviceFoundError) {
				return nil, nil, err
			}

			if err == nil {
				return dev, h, nil
			}

			select {

			case <-exit:
				return nil, nil, fmt.Errorf("%w", NoDeviceFoundError)

			case <-ticker.C:
			}
		}
	}
}

func getDevice(finder findDevice, checker deviceChecker) (device *libusb.Device, handle *libusb.DeviceHandle, done func(), err error) {

	ctx, err := libusb.NewContext()
	if err != nil {
		return nil, nil, nil, err
	}

	dev, h, err := finder(ctx, checker)
	if err != nil {
		_ = ctx.Close()
		return nil, nil, nil, err
	}

	return dev, h, func() {
		h.Close()
		ctx.Close()
	}, nil
}

func GetDevice() (device *libusb.Device, handle *libusb.DeviceHandle,
	done func(), err error) {

	return getDevice(findDeviceImpl, vendorProductDeviceChecker)
}

func GetDeviceByAddress(bus, address int) (device *libusb.Device, handle *libusb.DeviceHandle,
	done func(), err error) {

	return getDevice(findDeviceImpl, getAddressDeviceChecker(bus, address))
}

func GetDeviceWithPolling(pollInterval, timeout time.Duration) (device *libusb.Device, handle *libusb.DeviceHandle,
	done func(), err error) {

	return getDevice(getFindDeviceWithPollingImpl(pollInterval, timeout), vendorProductDeviceChecker)
}

func GetDeviceByAddressWithPolling(bus, address int, pollInterval, timeout time.Duration) (device *libusb.Device, handle *libusb.DeviceHandle,
	done func(), err error) {

	return getDevice(getFindDeviceWithPollingImpl(pollInterval, timeout), getAddressDeviceChecker(bus, address))
}

type write func(p []byte) (n int, err error)

func GetBulkWriteFunc(dev *libusb.Device, handle *libusb.DeviceHandle) (write, error) {

	cfg, err := dev.ActiveConfigDescriptor()
	if err != nil {
		return nil, err
	}

	for _, iface := range cfg.SupportedInterfaces {
		for _, desc := range iface.InterfaceDescriptors {
			for _, ep := range desc.EndpointDescriptors {
				if ep.Direction() == endpointOutDirection {

					return func(p []byte) (n int, err error) {
						return handle.BulkTransferOut(ep.EndpointAddress, p, 0)
					}, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("%w", NoOutEndpointFoundError)
}
