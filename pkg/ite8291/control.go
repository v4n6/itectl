package ite8291

import (
	"fmt"
)

// ite8291r3 controller commands.
const (
	GetEffectCommand          = 0x88
	SetEffectCommand          = 0x8
	SetBrightnessCommand      = 0x9
	SetColorCommand           = 0x14
	SetRowIndexCommand        = 0x16
	GetFirmwareVersionCommand = 0x80
)

// ite8291r3 effect operations.
const (
	SetEffectOp = 0x2
	SetOffOp    = 0x1
)

// ite8291r3 states.
const (
	OffState = 0x1
)

// ite8291r3 effect types.
const (
	BreathingEffect = 0x2
	WaveEffect      = 0x3
	RandomEffect    = 0x4
	RainbowEffect   = 0x5
	RippleEffect    = 0x6
	MarqueeEffect   = 0x9
	RaindropEffect  = 0xA
	AuroraEffect    = 0xE
	FireworksEffect = 0x11
	UserEffect      = 0x33
)

// Direction represents direction of ite8291r3 effects.
type Direction byte

// ite8291r3 effect directions.
const (
	DirectionNone  Direction = 0x0
	DirectionRight Direction = 0x1
	DirectionLeft  Direction = 0x2
	DirectionUp    Direction = 0x3
	DirectionDown  Direction = 0x4
)

// ite8291r3 controller parameters boundaries.
const (
	BrightnessMaxValue = 50

	SpeedMaxValue = 10

	ColorNumMinValue = 0x0
	ColorNumMaxValue = 0x8

	CustomColorNumMinValue = 0x1
	CustomColorNumMaxValue = 0x7
)

// special color numbers.
const (
	ColorNone   = 0x0
	ColorRandom = 0x8
)

// ite8291r3 keyboard dimension.
const (
	RowsNumber    = 6
	ColumnsNumber = 21
)

// auxiliary ite8291r3 keyboard buffer constants.
const (
	rowBufferLength = 3*ColumnsNumber + 2

	rowBlueOffset  = 1
	rowGreenOffset = rowBlueOffset + ColumnsNumber
	rowRedOffset   = rowGreenOffset + ColumnsNumber
)

// Device interface abstracts ite8291r3 usb device.
type Device interface {

	// ControlTransfer method of *libusb.DeviceHandle. It is used to set
	// effects and their attributes and get/set global ite8291r3
	// properties.
	ControlTransfer(requestType byte, request byte, value uint16, index uint16, data []byte, length int,
		timeout int) (int, error)

	// GetBulkWrite returns write function that can be used to set keys
	// colors.
	GetBulkWrite() (WriteFunc, error)

	// Close closes the device.
	Close() error
}

// Controller provides ite8291r3 controller functionality.
type Controller struct {
	dev Device
}

// NewController creates a new controller backed by provided ite8291r3
// usb device.
func NewController(d Device) *Controller {

	return &Controller{dev: d}
}

// Close cleans up underlying ite8291r3 usb device.
func (c *Controller) Close() error {
	return c.dev.Close()
}

// ControlSend sends data to ite8291r3 controller.
func (c *Controller) ControlSend(data []byte) error {

	_, err := c.dev.ControlTransfer(SendControlRequestType,
		0x009, // bRequest (HID set_report)
		0x300, // wValue (HID feature)
		0x001, // wIndex
		data,
		len(data),
		0)

	return err
}

// ControlSend receives data from ite8291r3 controller.
func (c *Controller) controlReceive(data []byte) error {
	_, err := c.dev.ControlTransfer(ReceiveControlRequestType,
		0x001, // bRequest (HID set_report)
		0x300, // wValue (HID feature)
		0x001, // wIndex
		data,
		len(data),
		0)

	return err
}

// SetEffect sets ite8291r3 effect and its attributes.
func (c *Controller) SetEffect(cntrl, effect, speed, brightness, colorNum,
	reactOrDiv byte, save bool) error {

	return c.ControlSend([]byte{SetEffectCommand, cntrl, effect, speed, brightness, colorNum,
		reactOrDiv, bool2Byte(save)})
}

// SetEffect sets ite8291r3 reactive effect and its attributes.
func (c *Controller) setEffectWithReactive(effect, speed, brightness, colorNum byte, reactive, save bool) error {
	return c.SetEffect(SetEffectOp, effect, speed, brightness, colorNum, bool2Byte(reactive), save)
}

// SetOffMode switches ite8291r3 keyboard backlight off.
func (c *Controller) SetOffMode() error {

	return c.SetEffect(SetOffOp, 0, 0, 0, 0, 0, false)
}

// State retrieves ite8291r3 keyboard backlight state: whether it's On
// (true) or Off (false).
func (c *Controller) State() (state bool, err error) {

	if err := c.ControlSend([]byte{GetEffectCommand}); err != nil {
		return false, err
	}

	out := []byte{8, 0, 0, 0, 0, 0, 0, 0}
	if err := c.controlReceive(out); err != nil {
		return false, err
	}

	return out[1] != OffState, nil
}

// SetBrightness sets brightness of ite8291r3 keyboard backlight. The
// maximum value is specified by BrightnessMaxValue.
func (c *Controller) SetBrightness(brightness byte) error {

	return c.ControlSend([]byte{SetBrightnessCommand, SetEffectOp, brightness})
}

// Brightness returns brightness of ite8291r3 keyboard backlight. The
// maximum value is specified by BrightnessMaxValue.
func (c *Controller) Brightness() (brightness byte, err error) {

	if err = c.ControlSend([]byte{GetEffectCommand}); err != nil {
		return 0, err
	}

	out := []byte{8, 0, 0, 0, 0, 0, 0, 0}
	if err := c.controlReceive(out); err != nil {
		return 0, err
	}

	return out[4], nil
}

// SetAuroraMode sets ite8291r3 keyboard backlight controller to
// 'aurora' effect.
func (c *Controller) SetAuroraMode(speed, brightness, colorNum byte,
	reactive, save bool) error {

	return c.setEffectWithReactive(AuroraEffect, speed, brightness,
		colorNum, reactive, save)
}

// SetBreathingMode sets ite8291r3 keyboard backlight controller to
// 'breathing' effect.
func (c *Controller) SetBreathingMode(speed, brightness, colorNum byte, save bool) error {

	return c.setEffectWithReactive(BreathingEffect, speed, brightness, colorNum, false, save)
}

// SetFireworksMode sets ite8291r3 keyboard backlight to 'fireworks'
// effect.
func (c *Controller) SetFireworksMode(speed, brightness, colorNum byte, reactive, save bool) error {

	return c.setEffectWithReactive(FireworksEffect, speed, brightness, colorNum, reactive, save)
}

// SetMarqueeMode sets ite8291r3 keyboard backlight to 'marquee'
// effect.
func (c *Controller) SetMarqueeMode(speed, brightness byte, save bool) error {

	return c.setEffectWithReactive(MarqueeEffect, speed, brightness, 0, false, save)
}

// SetRainbowMode sets ite8291r3 keyboard backlight to 'rainbow' effect.
func (c *Controller) SetRainbowMode(brightness byte, save bool) error {

	return c.setEffectWithReactive(RainbowEffect, 0, brightness, 0, false, save)
}

// SetRaindropMode sets ite8291r3 keyboard backlight to 'raindrop' effect.
func (c *Controller) SetRaindropMode(speed, brightness, colorNum byte, save bool) error {

	return c.setEffectWithReactive(RaindropEffect, speed, brightness, colorNum, false, save)
}

// SetRandomMode sets ite8291r3 keyboard backlight to 'random' effect.
func (c *Controller) SetRandomMode(speed, brightness, colorNum byte, reactive, save bool) error {

	return c.setEffectWithReactive(RandomEffect, speed, brightness, colorNum, reactive, save)
}

// SetRippleMode sets ite8291r3 keyboard backlight to 'ripple' effect.
func (c *Controller) SetRippleMode(speed, brightness, colorNum byte, reactive, save bool) error {

	return c.setEffectWithReactive(RippleEffect, speed, brightness, colorNum, reactive, save)
}

// SetWaveMode sets ite8291r3 keyboard backlight to 'wave' effect.
func (c *Controller) SetWaveMode(speed, brightness byte, direction Direction, save bool) error {

	return c.SetEffect(SetEffectOp, WaveEffect, speed, brightness, 0, byte(direction), save)
}

// SetUserMode sets ite8291r3 keyboard backlight to 'user' effect.
func (c *Controller) setUserMode(brightness byte, save bool) error {

	return c.setEffectWithReactive(UserEffect, 0, brightness, 0, false, save)
}

// setRowIndex sets current keyboard row of 'user' effect to the
// specified value.
func (c *Controller) setRowIndex(idx byte) error {

	return c.ControlSend([]byte{SetRowIndexCommand, 0, idx})
}

// SetSingleColorMode sets color of all keyboard backlight key to the specified color.
func (c *Controller) SetSingleColorMode(brightness byte, color *Color, save bool) error {

	if err := c.setUserMode(brightness, save); err != nil {
		return err
	}

	write, err := c.dev.GetBulkWrite()
	if err != nil {
		return err
	}

	rowBuffer := make([]byte, rowBufferLength)
	for i := range RowsNumber {
		if err := c.setRowIndex(byte(i)); err != nil {
			return err
		}

		for j := range ColumnsNumber {
			rowBuffer[j+rowBlueOffset] = color.Blue
			rowBuffer[j+rowGreenOffset] = color.Green
			rowBuffer[j+rowRedOffset] = color.Red
		}

		if _, err := write(rowBuffer); err != nil {
			return err
		}
	}

	return nil
}

// SetColor sets predefined color specified by its colorNum to the
// given color.
func (c *Controller) SetColor(colorNum byte, color *Color) error {

	return c.ControlSend([]byte{SetColorCommand, 0, colorNum, color.Red, color.Green, color.Blue})
}

// SetColors sets predefined colors to the provided colors.
func (c *Controller) SetColors(colors []*Color) error {

	for i, col := range colors[:CustomColorNumMaxValue-CustomColorNumMinValue+1] {

		if err := c.ControlSend([]byte{SetColorCommand, 0, byte(i + 1), col.Red, col.Green, col.Blue}); err != nil {
			return err
		}
	}

	return nil
}

// FirmwareVersion returns firmware version of ite8291r3 controller
// as string.
func (c *Controller) FirmwareVersion() (string, error) {

	if err := c.ControlSend([]byte{GetFirmwareVersionCommand}); err != nil {
		return "", err
	}

	out := []byte{8, 0, 0, 0, 0, 0, 0, 0}
	err := c.controlReceive(out)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d.%d.%d.%d", out[1], out[2], out[3], out[4]), nil
}

// bool2Byte converts provided boolean to byte. False is converted to
// 0, true to 1.
func bool2Byte(b bool) byte {
	if b {
		return 1
	} else {
		return 0
	}
}
