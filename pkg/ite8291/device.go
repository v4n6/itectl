package ite8291

import (
	"errors"
	"fmt"
	"time"

	"github.com/gotmc/libusb/v2"
)

// ite8291r3 request types
const (
	sendControlRequestType = byte(libusb.HostToDevice) | byte(libusb.Class) | byte(libusb.InterfaceRecipient)
	rcvControlRequestType  = byte(libusb.DeviceToHost) | byte(libusb.Class) | byte(libusb.InterfaceRecipient)
)

// vendor id of ite8291r3 usb device
const vendorID = 0x048D

// product ids of ite8291r3 usb device
var (
	productIDs = map[uint16]bool{0x6004: true, 0x6006: true, 0xCE00: true}
)

// NoDeviceFoundError is an error indicating that no ite8291r3 device can be found
var NoDeviceFoundError = errors.New("no ite8291r3 device found")

// UnsupportedDeviceError is an error indicating that a device with given bus and address is not an ite8291r3 device
var UnsupportedDeviceError = errors.New("is not ite8291 device")

// NoOutEndpointFoundError is an error indicating that no out endpoint can be found at ite8291r3 device
var NoOutEndpointFoundError = errors.New("no output endpoint found")

// USBDevice represents ite8291r3 usb device
type USBDevice struct {
	ctx *libusb.Context
	*libusb.DeviceHandle
	dev *libusb.Device
}

// Close ...
func (d *USBDevice) Close() error {
	err1 := d.DeviceHandle.Close()
	err2 := d.ctx.Close()

	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}

	return nil
}

// WriteFunc is a write function that can be used to set keys colors.
type WriteFunc func(p []byte) (n int, err error)

const endpointOutDirection libusb.EndpointDirection = 0

// GetBulkWrite returns write function that can be used to set keys colors.
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

type Context interface {
	DeviceList() ([]*libusb.Device, error)
}

type DeviceCheckerFunc func(dev *libusb.Device) (bool, error)

func VendorProductDeviceCheckerFunc(dev *libusb.Device) (bool, error) {
	if d, err := dev.DeviceDescriptor(); err == nil &&
		d.VendorID == vendorID && productIDs[d.ProductID] {

		return true, nil
	}

	return false, nil
}

func NewAddressDeviceCheckerFunc(bus, address int) DeviceCheckerFunc {
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

type FindDeviceFunc func(ctx *libusb.Context, checker DeviceCheckerFunc) (usbDevice *USBDevice, err error)

const targetInterfaceNumber = 1

func FindDeviceWithoutPollingFunc(ctx *libusb.Context, checker DeviceCheckerFunc) (usbDevice *USBDevice, err error) {

	devs, err := ctx.DeviceList()
	if err != nil {
		return nil, err
	}

	for _, dev := range devs {
		match, err := checker(dev)
		if err != nil {
			_ = ctx.Close()
			return nil, err
		}
		if match {
			h, err := dev.Open()
			if err != nil {
				_ = ctx.Close()
				return nil, err
			}

			usbDevice = &USBDevice{DeviceHandle: h, dev: dev, ctx: ctx}

			kern, err := usbDevice.KernelDriverActive(targetInterfaceNumber)
			if err != nil {
				_ = usbDevice.Close()
				return nil, err
			}

			if kern {
				if err = usbDevice.DetachKernelDriver(targetInterfaceNumber); err != nil {
					_ = usbDevice.Close()
					return nil, err
				}
			}

			return usbDevice, nil
		}
	}

	return nil, fmt.Errorf("%w", NoDeviceFoundError)
}

func NewFindDeviceWithPollingFunc(pollInterval, timeout time.Duration) FindDeviceFunc {

	return func(ctx *libusb.Context, checker DeviceCheckerFunc) (usbDevice *USBDevice, err error) {

		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()

		exit := make(chan bool)
		go func() {
			time.Sleep(timeout)
			exit <- true
		}()

		for {

			dev, err := FindDeviceWithoutPollingFunc(ctx, checker)
			if err != nil && !errors.Is(err, NoDeviceFoundError) {
				return nil, err
			}

			if err == nil {
				return dev, nil
			}

			select {

			case <-exit:
				return nil, fmt.Errorf("%w", NoDeviceFoundError)

			case <-ticker.C:
			}
		}
	}
}

func GetDevice(finder FindDeviceFunc, checker DeviceCheckerFunc) (usbDevice *USBDevice, err error) {

	ctx, err := libusb.NewContext()
	if err != nil {
		return nil, err
	}

	dev, err := finder(ctx, checker)
	if err != nil {
		if errors.Is(err, NoDeviceFoundError) {
			_ = ctx.Close()
		}

		return nil, err
	}

	return dev, nil
}

func GetDeviceByVendorProduct() (usbDevice *USBDevice, err error) {

	return GetDevice(FindDeviceWithoutPollingFunc, VendorProductDeviceCheckerFunc)
}

func GetDeviceByAddress(bus, address int) (usbDevice *USBDevice, err error) {

	return GetDevice(FindDeviceWithoutPollingFunc, NewAddressDeviceCheckerFunc(bus, address))
}

func GetDeviceWithPolling(pollInterval, timeout time.Duration) (usbDevice *USBDevice, err error) {

	return GetDevice(NewFindDeviceWithPollingFunc(pollInterval, timeout), VendorProductDeviceCheckerFunc)
}

func GetDeviceByAddressWithPolling(bus, address int, pollInterval, timeout time.Duration) (usbDevice *USBDevice, err error) {

	return GetDevice(NewFindDeviceWithPollingFunc(pollInterval, timeout), NewAddressDeviceCheckerFunc(bus, address))
}

type DeviceDescriptor struct {
	*libusb.Descriptor
	BusNumber     int
	DeviceAddress int
	PortNumber    int
}

func ListDevices() (devices []*DeviceDescriptor, err error) {

	ctx, err := libusb.NewContext()
	if err != nil {
		return nil, err
	}
	defer ctx.Close()

	devs, err := ctx.DeviceList()
	if err != nil {
		return nil, err
	}

	devices = []*DeviceDescriptor{}
	for _, dev := range devs {
		found, err := VendorProductDeviceCheckerFunc(dev)
		if err != nil {
			return nil, err
		}
		if found {
			d, err := dev.DeviceDescriptor()
			if err != nil {
				return nil, err
			}

			busNumber, err := dev.BusNumber()
			if err != nil {
				return nil, err
			}

			deviceAddress, err := dev.DeviceAddress()
			if err != nil {
				return nil, err
			}

			portNumber, err := dev.PortNumber()
			if err != nil {
				return nil, err
			}

			devices = append(devices, &DeviceDescriptor{Descriptor: d,
				BusNumber: busNumber, DeviceAddress: deviceAddress, PortNumber: portNumber})
		}
	}

	return devices, nil
}
