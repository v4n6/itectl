package ite8291

import "github.com/gotmc/libusb/v2"

// RequestType constants
const (
	sendControlRequestType = byte(libusb.HostToDevice) | byte(libusb.Class) | byte(libusb.InterfaceRecipient)
	rcvControlRequestType  = byte(libusb.DeviceToHost) | byte(libusb.Class) | byte(libusb.InterfaceRecipient)
)

// Commands
const (
	getEffectCommand          = 0x88
	setEffectCommand          = 0x8
	setBrightnessCommand      = 0x9
	setColorCommand           = 0x14
	setRowIndexCommand        = 0x16
	getFirmwareVersionCommand = 0x80
)

// ite8291 control
const (
	setEffectControl = 0x2
	setOffControl    = 0x1
)

// state indicators
const (
	stateOff = 0x1
)

// effect constants
const (
	breathingEffect = 0x2
	waveEffect      = 0x3
	randomEffect    = 0x4
	rainbowEffect   = 0x5
	rippleEffect    = 0x6
	marqueeEffect   = 0x9
	raindropEffect  = 0xA
	auroraEffect    = 0xE
	fireworksEffect = 0x11
	userEffect      = 0x33
)

type Direction byte

// direction
const (
	DirectionNone  Direction = 0x0
	DirectionRight Direction = 0x1
	DirectionLeft  Direction = 0x2
	DirectionUp    Direction = 0x3
	DirectionDown  Direction = 0x4
)

// keyboard size
const (
	RowsNumber    = 6
	ColumnsNumber = 21
)

// ite8291 keyboard buffer props
const (
	rowBufferLength = 3*ColumnsNumber + 2

	rowBlueOffset  = 1
	rowGreenOffset = rowBlueOffset + ColumnsNumber
	rowRedOffset   = rowGreenOffset + ColumnsNumber
)

// prop boundries
const (
	BrightnessMaxValue = 50

	SpeedMaxValue = 10

	ColorMaxValue = 0x8
)

// special colors
const (
	ColorNone   = 0x0
	ColorRandom = 0x8
)

type ControlTransfer interface {
	ControlTransfer(requestType byte, request byte, value uint16, index uint16, data []byte, length int,
		timeout int,
	) (int, error)
}

func controlSend(ctl ControlTransfer, data []byte) error {

	_, err := ctl.ControlTransfer(sendControlRequestType,
		0x009, // bRequest (HID set_report)
		0x300, // wValue (HID feature)
		0x001, // wIndex
		data,
		len(data),
		0)

	return err
}

func controlReceive(ctl ControlTransfer, data []byte) error {
	_, err := ctl.ControlTransfer(rcvControlRequestType,
		0x001, // bRequest (HID set_report)
		0x300, // wValue (HID feature)
		0x001, // wIndex
		data,
		len(data),
		0)

	return err
}

func bool2Byte(b bool) byte {
	if b {

		return 1
	} else {

		return 0
	}
}

func setEffect(ctl ControlTransfer, cntrl, effect, speed, brightness, colorNum, reactive_or_direction byte, save bool) error {

	return controlSend(ctl, []byte{setEffectCommand, cntrl, effect, speed, brightness, colorNum, reactive_or_direction, bool2Byte(save)})
}

func setEffectWithReactive(ctl ControlTransfer, cntrl, effect, speed, brightness, colorNum byte, reactive, save bool) error {
	return setEffect(ctl, cntrl, effect, speed, brightness, colorNum, bool2Byte(reactive), save)
}

func SetOff(ctl ControlTransfer) error {

	return setEffect(ctl, setOffControl, 0, 0, 0, 0, 0, false)
}

func State(ctl ControlTransfer) (state bool, err error) {
	err = controlSend(ctl, []byte{getEffectCommand})
	if err != nil {
		return false, err
	}

	out := []byte{8, 0, 0, 0, 0, 0, 0, 0}
	err = controlReceive(ctl, out)
	if err != nil {
		return false, err
	}

	return out[1] != stateOff, nil
}

func SetBrightness(ctl ControlTransfer, brightness byte) error {
	return controlSend(ctl, []byte{setBrightnessCommand, setEffectControl, brightness})
}

func GetBrightness(ctl ControlTransfer) (brightness byte, err error) {
	err = controlSend(ctl, []byte{getEffectCommand})
	if err != nil {
		return 0, err
	}

	out := []byte{8, 0, 0, 0, 0, 0, 0, 0}
	err = controlReceive(ctl, out)
	if err != nil {
		return 0, err
	}

	return out[4], nil
}

func SetAuroraMode(ctl ControlTransfer, speed, brightness, colorNum byte, reactive, save bool) error {

	return setEffectWithReactive(ctl, setEffectControl, auroraEffect, speed, brightness, colorNum, reactive, save)
}

func SetBreathingMode(ctl ControlTransfer, speed, brightness, colorNum byte, save bool) error {

	return setEffectWithReactive(ctl, setEffectControl, breathingEffect, speed, brightness, colorNum, false, save)
}

func SetFireworksMode(ctl ControlTransfer, speed, brightness, colorNum byte, reactive, save bool) error {

	return setEffectWithReactive(ctl, setEffectControl, fireworksEffect, speed, brightness, colorNum, reactive, save)
}

func SetMarqueeMode(ctl ControlTransfer, speed, brightness byte, save bool) error {

	return setEffectWithReactive(ctl, setEffectControl, marqueeEffect, speed, brightness, 0, false, save)
}

func SetRainbowMode(ctl ControlTransfer, brightness byte, save bool) error {

	return setEffectWithReactive(ctl, setEffectControl, rainbowEffect, 0, brightness, 0, false, save)
}

func SetRainDropMode(ctl ControlTransfer, speed, brightness, colorNum byte, save bool) error {

	return setEffectWithReactive(ctl, setEffectControl, raindropEffect, speed, brightness, colorNum, false, save)
}

func SetRandomMode(ctl ControlTransfer, speed, brightness, colorNum byte, reactive, save bool) error {

	return setEffectWithReactive(ctl, setEffectControl, randomEffect, speed, brightness, colorNum, reactive, save)
}

func SetRippleMode(ctl ControlTransfer, speed, brightness, colorNum byte, reactive, save bool) error {

	return setEffectWithReactive(ctl, setEffectControl, rippleEffect, speed, brightness, colorNum, reactive, save)
}

func SetWaveMode(ctl ControlTransfer, speed, brightness byte, direction Direction, save bool) error {

	return setEffect(ctl, setEffectControl, waveEffect, speed, brightness, 0, byte(direction), save)
}

func setUserMode(ctl ControlTransfer, brightness byte, save bool) error {

	return setEffectWithReactive(ctl, setEffectControl, userEffect, 0, brightness, 0, false, save)
}

func setRowIndex(ctl ControlTransfer, idx byte) error {

	return controlSend(ctl, []byte{setRowIndexCommand, 0, idx})
}

func SetSingleColorMode(dev *libusb.Device, h *libusb.DeviceHandle, brightness, red, green, blue byte, save bool) error {

	err := setUserMode(h, brightness, save)
	if err != nil {
		return err
	}

	writer, err := GetBulkWriteFunc(dev, h)
	if err != nil {
		return err
	}

	rowBuffer := make([]byte, rowBufferLength)
	for i := range RowsNumber {
		err = setRowIndex(h, byte(i))
		if err != nil {
			return err
		}

		for j := range ColumnsNumber {
			rowBuffer[j+rowBlueOffset] = blue
			rowBuffer[j+rowGreenOffset] = green
			rowBuffer[j+rowRedOffset] = red
		}

		_, err := writer(rowBuffer)
		if err != nil {
			return err
		}
	}

	return nil
}

func SetColor(ctl ControlTransfer, colorNum, red, green, blue byte) error {

	return controlSend(ctl, []byte{setColorCommand, 0, colorNum, red, green, blue})
}
