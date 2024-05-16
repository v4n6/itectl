package cmd

import (
	"bytes"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/params"
	"github.com/v4n6/itectl/pkg/ite8291"
)

var _ = Describe("cmd package", func() {

	var dev *deviceStubT

	var findDevCall *findDeviceCallT
	var readConfigCall *readConfigCallT

	var cmdArgs []string
	var cmdErr error
	var cmdOut, cmdErrOut *gbytes.Buffer

	var subCmd string
	var useCmdImplicit bool

	JustBeforeEach(func() {

		if useCmdImplicit {
			// configure subCmd as configured mode
			mode, found := strings.CutSuffix(subCmd, "-mode")
			Ω(found).Should(BeTrue())

			Ω(readConfigCall.v.MergeConfigMap(map[string]any{ // add viper mode config
				"mode": mode,
			})).Should(Succeed())

		} else {
			// use subCmd explicitly as 1st arg
			cmdArgs = slices.Insert(cmdArgs, 0, subCmd)
		}

		// execute the command
		cmdErr = ExecuteCmd(cmdArgs, cmdOut, cmdErrOut,
			newFindDevice(dev, findDevCall), // find device function
			newReadConfig(readConfigCall),
		)
		// close out & err streams after command execution
		Ω(cmdErrOut.Close()).Should(Succeed())
		Ω(cmdOut.Close()).Should(Succeed())
	})

	BeforeEach(func() {
		// init
		dev = &deviceStubT{}
		findDevCall = &findDeviceCallT{}
		readConfigCall = &readConfigCallT{v: viper.New()}

		cmdOut, cmdErrOut = gbytes.NewBuffer(), gbytes.NewBuffer()

		subCmd, useCmdImplicit = "", false

		// add named colors to viper config
		Ω(readConfigCall.v.MergeConfigMap(map[string]any{
			"namedColors": namedColorsConfig,
		})).Should(Succeed())
	})

	Describe("happy cases", func() {

		JustBeforeEach(func() {
			Ω(cmdErr).Should(Succeed()) // assert command executed without error
		})

		DescribeTableSubtree("with configuration",
			func(defaults *defaultsT) {

				var defs *defaultsT

				BeforeEach(func() {

					if defaults == nil {
						defs = zeroDefaults() // use only programm defaults
						return
					}

					defs = defaults // use specified defaults
					newConf := map[string]any{
						params.BrightnessProp: defs.brightness,
						params.SpeedProp:      defs.speed,
						params.ColorNumProp:   defs.colorNum,
						params.DirectionProp:  defs.direction,

						params.SingleColorProp: defs.singleColor,
						params.ResetProp:       defs.reset,
						params.ReactiveProp:    defs.reactive,
						params.SaveProp:        defs.save,

						params.PredefinedColorProp: predefinedColorsConfig(defs.predefinedColors),

						"poll": map[string]any{
							"interval": defs.pollInterval,
							"timeout":  defs.pollTimeout,
						},
					}

					if defs.deviceBus >= 0 && defs.deviceAddress >= 0 {
						// configure device
						newConf["device"] = map[string]any{
							"bus":     defs.deviceBus,
							"address": defs.deviceAddress,
						}
					}

					Ω(readConfigCall.v.MergeConfigMap(newConf)).Should(Succeed()) // merge with viper config
				})

				DescribeTableSubtree("state",
					func(ex *execT) {

						BeforeEach(func() {
							subCmd, cmdArgs = "state", ex.genArgs(nil, defs, r)
						})

						DescribeTableSubtree("when device is",
							func(state int, msg string) {

								BeforeEach(func() {
									dev.ctlChangedData = [][]byte{nil, []byte{8, byte(state), 0, 0, 0, 0, 0, 0}}
								})

								It("correctly calls usb device and prints correct state", func() {
									assertOnlyControlCall(dev, []*ctlArgsT{&ctlArgsT{
										requestType: 0x21, request: 9, value: 0x300, index: 1, data: []byte{0x88}, length: 1, timeout: 0,
									}, &ctlArgsT{
										requestType: 0xA1, request: 1, value: 0x300, index: 1, data: []byte{8, 0, 0, 0, 0, 0, 0, 0}, length: 8, timeout: 0,
									}})

									assertCommandOutput(cmdOut, cmdErrOut, msg)

									assertFindDeviceCall(findDevCall, ex.Device(), ex.DeviceBus(), ex.DeviceAddress(),
										ex.PollInterval(), ex.PollTimeout())

									assertReadConfigCall(readConfigCall, ex.ConfigFile())
								})
							},
							Entry("when device is ON", 2, "On"),
							Entry("when device is OFF", 1, "Off"),
						)
					},
					newExecs().
						ConfigFile(configFileAll...).
						PollInterval(pollIntervalAll...).
						PollTimeout(pollTimeoutAll...).entries(),
					newExecs().
						RequiredDeviceBus(deviceBusAll...).
						RequiredDeviceAddress(deviceAddressAll...).entries(),
				)

				DescribeTableSubtree("set-brightness",

					func(ex *execT) {

						BeforeEach(func() {
							subCmd, cmdArgs = "set-brightness", ex.genArgs(nil, defs, r)
						})

						It("correctly calls usb device", func() {
							assertOnlyControlCall(dev, []*ctlArgsT{&ctlArgsT{
								requestType: 0x21, request: 9, value: 0x300, index: 1,
								data: []byte{0x9, 0x2, ex.Brightness()}, length: 3, timeout: 0,
							}})

							assertFindDeviceCall(findDevCall, ex.Device(), ex.DeviceBus(), ex.DeviceAddress(),
								ex.PollInterval(), ex.PollTimeout())

							assertReadConfigCall(readConfigCall, ex.ConfigFile())
						})
					},
					newExecs().
						ConfigFile(configFileAll...).
						PollInterval(pollIntervalAll...).
						PollTimeout(pollTimeoutAll...).
						RequiredBrightness(brightnessAll...).entries(),
					newExecs().
						RequiredDeviceBus(deviceBusAll...).
						RequiredDeviceAddress(deviceAddressAll...).
						RequiredBrightness(brightnessAll...).entries(),
				)

				DescribeTableSubtree("brightness",
					func(ex *execT) {
						var brightness int

						BeforeEach(func() {
							brightness, subCmd = r.Intn(ite8291.BrightnessMaxValue), "brightness"
							cmdArgs = ex.genArgs(nil, defs, r)

							dev.ctlChangedData = [][]byte{nil, []byte{8, 0, 0, 0, byte(brightness), 0, 0, 0}}
						})

						It("correctly calls usb device and prints obtained brightness", func() {

							assertOnlyControlCall(dev, []*ctlArgsT{&ctlArgsT{
								requestType: 0x21, request: 9, value: 0x300, index: 1, data: []byte{0x88}, length: 1, timeout: 0,
							}, &ctlArgsT{
								requestType: 0xA1, request: 1, value: 0x300, index: 1, data: []byte{8, 0, 0, 0, 0, 0, 0, 0}, length: 8, timeout: 0,
							}})

							assertCommandOutput(cmdOut, cmdErrOut, strconv.Itoa(brightness))

							assertFindDeviceCall(findDevCall, ex.Device(), ex.DeviceBus(), ex.DeviceAddress(),
								ex.PollInterval(), ex.PollTimeout())

							assertReadConfigCall(readConfigCall, ex.ConfigFile())
						})
					},
					newExecs().
						ConfigFile(configFileAll...).
						PollInterval(pollIntervalAll...).
						PollTimeout(pollTimeoutAll...).entries(),
					newExecs().
						RequiredDeviceBus(deviceBusAll...).
						RequiredDeviceAddress(deviceAddressAll...).entries(),
				)

				DescribeTableSubtree("firmware-version",
					func(ex *execT) {

						var version []byte

						BeforeEach(func() {
							subCmd, cmdArgs = "firmware-version", ex.genArgs(nil, defs, r)

							version = []byte{byte(r.Intn(255)), byte(r.Intn(255)), byte(r.Intn(255)), byte(r.Intn(255))}

							dev.ctlChangedData = [][]byte{nil, []byte{8, version[0], version[1], version[2], version[3], 0, 0, 0}}
						})

						It("correctly calls usb device and prints obtained firmware version", func() {

							assertOnlyControlCall(dev, []*ctlArgsT{&ctlArgsT{
								requestType: 0x21, request: 9, value: 0x300, index: 1, data: []byte{0x80}, length: 1, timeout: 0,
							}, &ctlArgsT{
								requestType: 0xA1, request: 1, value: 0x300, index: 1, data: []byte{8, 0, 0, 0, 0, 0, 0, 0}, length: 8, timeout: 0,
							}})

							assertCommandOutput(cmdOut, cmdErrOut, fmt.Sprintf("%d.%d.%d.%d", version[0], version[1], version[2], version[3]))

							assertFindDeviceCall(findDevCall, ex.Device(), ex.DeviceBus(), ex.DeviceAddress(),
								ex.PollInterval(), ex.PollTimeout())

							assertReadConfigCall(readConfigCall, ex.ConfigFile())
						})
					},
					newExecs().
						ConfigFile(configFileAll...).
						PollInterval(pollIntervalAll...).
						PollTimeout(pollTimeoutAll...).entries(),
					newExecs().
						RequiredDeviceBus(deviceBusAll...).
						RequiredDeviceAddress(deviceAddressAll...).entries(),
				)

				DescribeTableSubtree("set-color",
					func(ex *execT) {

						BeforeEach(func() {
							subCmd, cmdArgs = "set-color", ex.genArgs(nil, defs, r)
						})

						It("correctly calls usb device", func() {

							col := ex.Color()

							assertOnlyControlCall(dev, []*ctlArgsT{&ctlArgsT{
								requestType: 0x21, request: 9, value: 0x300, index: 1,
								data: []byte{0x14, 0x0, ex.CustomColorNum(), col.Red, col.Green, col.Blue}, length: 6, timeout: 0,
							}})

							assertFindDeviceCall(findDevCall, ex.Device(), ex.DeviceBus(), ex.DeviceAddress(),
								ex.PollInterval(), ex.PollTimeout())

							assertReadConfigCall(readConfigCall, ex.ConfigFile())
						})
					},
					newExecs().
						ConfigFile(configFileAll...).
						PollInterval(pollIntervalAll...).
						PollTimeout(pollTimeoutAll...).
						CustomColorNum(customColorNumAll...).
						RequiredRed(colorAll...).
						Green(colorAll...).
						Blue(colorAll...).entries(),
					newExecs().
						RequiredDeviceBus(deviceBusAll...).
						RequiredDeviceAddress(deviceAddressAll...).
						CustomColorNum(randomLabel).
						RequiredRGB(rgbAll...).entries(),
					newExecs().
						CustomColorNum(randomLabel).
						RequiredColorName(colorNameAll...).entries(),
				)

				DescribeTableSubtree("called implicitly",
					func(implicitCmd bool) {

						BeforeEach(func() {
							useCmdImplicit = implicitCmd
						})

						DescribeTableSubtree("aurora-mode",
							func(ex *execT) {

								BeforeEach(func() {
									subCmd, cmdArgs = "aurora-mode", ex.genArgs(nil, defs, r)
								})

								It("correctly calls usb device", func() {
									assertOnlyControlCallWithReset(dev, ex.PredefinedColors(), ex.Reset(),
										[]*ctlArgsT{&ctlArgsT{
											requestType: 0x21, request: 9, value: 0x300, index: 1,
											data: []byte{0x8, 0x2, 0xe, ex.Speed(), ex.Brightness(), ex.ColorNum(),
												ex.Reactive(), ex.Save()}, length: 8, timeout: 0,
										}})

									assertFindDeviceCall(findDevCall, ex.Device(), ex.DeviceBus(), ex.DeviceAddress(),
										ex.PollInterval(), ex.PollTimeout())

									assertReadConfigCall(readConfigCall, ex.ConfigFile())
								})
							},
							newExecs().
								ConfigFile(configFileAll...).
								PollInterval(pollIntervalAll...).
								PollTimeout(pollTimeoutAll...).
								Brightness(brightnessAll...).
								Speed(speedAll...).
								ColorNum(colorNumAll...).
								Reset(resetAll...).
								Reactive(reactiveAll...).
								Save(saveAll...).entries(),
							newExecs().
								RequiredDeviceBus(deviceBusAll...).
								RequiredDeviceAddress(deviceAddressAll...).entries(),
						)

						DescribeTableSubtree("breath-mode",
							func(ex *execT) {

								BeforeEach(func() {
									subCmd, cmdArgs = "breath-mode", ex.genArgs(nil, defs, r)
								})

								It("correctly calls usb device", func() {
									assertOnlyControlCallWithReset(dev, ex.PredefinedColors(), ex.Reset(),
										[]*ctlArgsT{&ctlArgsT{
											requestType: 0x21, request: 9, value: 0x300, index: 1,
											data: []byte{0x8, 0x2, 0x2, ex.Speed(), ex.Brightness(), ex.ColorNum(), 0,
												ex.Save()}, length: 8, timeout: 0,
										}})

									assertFindDeviceCall(findDevCall, ex.Device(), ex.DeviceBus(), ex.DeviceAddress(),
										ex.PollInterval(), ex.PollTimeout())

									assertReadConfigCall(readConfigCall, ex.ConfigFile())
								})
							},
							newExecs().
								ConfigFile(configFileAll...).
								PollInterval(pollIntervalAll...).
								PollTimeout(pollTimeoutAll...).
								Brightness(brightnessAll...).
								Speed(speedAll...).
								ColorNum(colorNumAll...).
								Reset(resetAll...).
								Save(saveAll...).entries(),
							newExecs().
								ConfigFile(configFileAll...).
								RequiredDeviceBus(deviceBusAll...).
								RequiredDeviceAddress(deviceAddressAll...).entries(),
						)

						DescribeTableSubtree("fireworks-mode",
							func(ex *execT) {

								BeforeEach(func() {
									subCmd, cmdArgs = "fireworks-mode", ex.genArgs(nil, defs, r)
								})

								It("correctly calls usb device", func() {
									assertOnlyControlCallWithReset(dev, ex.PredefinedColors(), ex.Reset(),
										[]*ctlArgsT{&ctlArgsT{
											requestType: 0x21, request: 9, value: 0x300, index: 1,
											data: []byte{0x8, 0x2, 0x11, ex.Speed(), ex.Brightness(), ex.ColorNum(),
												ex.Reactive(), ex.Save()}, length: 8, timeout: 0,
										}})

									assertFindDeviceCall(findDevCall, ex.Device(), ex.DeviceBus(), ex.DeviceAddress(),
										ex.PollInterval(), ex.PollTimeout())

									assertReadConfigCall(readConfigCall, ex.ConfigFile())
								})
							},
							newExecs().
								ConfigFile(configFileAll...).
								PollInterval(pollIntervalAll...).
								PollTimeout(pollTimeoutAll...).
								Brightness(brightnessAll...).
								Speed(speedAll...).
								ColorNum(colorNumAll...).
								Reset(resetAll...).
								Reactive(reactiveAll...).
								Save(saveAll...).entries(),
							newExecs().
								RequiredDeviceBus(deviceBusAll...).
								RequiredDeviceAddress(deviceAddressAll...).entries(),
						)

						DescribeTableSubtree("marquee-mode",
							func(ex *execT) {

								BeforeEach(func() {
									subCmd, cmdArgs = "marquee-mode", ex.genArgs(nil, defs, r)
								})

								It("correctly calls usb device", func() {
									assertOnlyControlCallWithReset(dev, ex.PredefinedColors(), ex.Reset(),
										[]*ctlArgsT{&ctlArgsT{
											requestType: 0x21, request: 9, value: 0x300, index: 1,
											data: []byte{0x8, 0x2, 0x9, ex.Speed(), ex.Brightness(), 0, 0,
												ex.Save()}, length: 8, timeout: 0,
										}})

									assertFindDeviceCall(findDevCall, ex.Device(), ex.DeviceBus(), ex.DeviceAddress(),
										ex.PollInterval(), ex.PollTimeout())

									assertReadConfigCall(readConfigCall, ex.ConfigFile())
								})
							},
							newExecs().
								ConfigFile(configFileAll...).
								PollInterval(pollIntervalAll...).
								PollTimeout(pollTimeoutAll...).
								Brightness(brightnessAll...).
								Speed(speedAll...).
								Reset(resetAll...).
								Save(saveAll...).entries(),
							newExecs().
								RequiredDeviceBus(deviceBusAll...).
								RequiredDeviceAddress(deviceAddressAll...).entries(),
						)

						DescribeTableSubtree("off-mode",
							func(ex *execT) {

								BeforeEach(func() {
									subCmd, cmdArgs = "off-mode", ex.genArgs(nil, defs, r)
								})

								It("correctly calls usb device", func() {
									assertOnlyControlCallWithReset(dev, ex.PredefinedColors(), ex.Reset(),
										[]*ctlArgsT{&ctlArgsT{
											requestType: 0x21, request: 9, value: 0x300, index: 1,
											data: []byte{0x8, 0x1, 0x0, 0, 0, 0, 0, 0}, length: 8, timeout: 0,
										}})

									assertFindDeviceCall(findDevCall, ex.Device(), ex.DeviceBus(), ex.DeviceAddress(),
										ex.PollInterval(), ex.PollTimeout())

									assertReadConfigCall(readConfigCall, ex.ConfigFile())
								})
							},
							newExecs().
								ConfigFile(configFileAll...).
								PollInterval(pollIntervalAll...).
								PollTimeout(pollTimeoutAll...).
								Reset(resetAll...).entries(),
							newExecs().
								RequiredDeviceBus(deviceBusAll...).
								RequiredDeviceAddress(deviceAddressAll...).entries(),
						)

						DescribeTableSubtree("rainbow-mode",
							func(ex *execT) {

								BeforeEach(func() {
									subCmd, cmdArgs = "rainbow-mode", ex.genArgs(nil, defs, r)
								})

								It("correctly calls usb device", func() {
									assertOnlyControlCallWithReset(dev, ex.PredefinedColors(), ex.Reset(),
										[]*ctlArgsT{&ctlArgsT{
											requestType: 0x21, request: 9, value: 0x300, index: 1,
											data: []byte{0x8, 0x2, 0x5, 0, ex.Brightness(), 0, 0,
												ex.Save()}, length: 8, timeout: 0,
										}})

									assertFindDeviceCall(findDevCall, ex.Device(), ex.DeviceBus(), ex.DeviceAddress(),
										ex.PollInterval(), ex.PollTimeout())

									assertReadConfigCall(readConfigCall, ex.ConfigFile())
								})
							},
							newExecs().
								ConfigFile(configFileAll...).
								PollInterval(pollIntervalAll...).
								PollTimeout(pollTimeoutAll...).
								Brightness(brightnessAll...).
								Reset(resetAll...).
								Save(saveAll...).entries(),
							newExecs().
								RequiredDeviceBus(deviceBusAll...).
								RequiredDeviceAddress(deviceAddressAll...).entries(),
						)

						DescribeTableSubtree("raindrop-mode",
							func(ex *execT) {

								BeforeEach(func() {
									subCmd, cmdArgs = "raindrop-mode", ex.genArgs(nil, defs, r)
								})

								It("correctly calls usb device", func() {
									assertOnlyControlCallWithReset(dev, ex.PredefinedColors(), ex.Reset(),
										[]*ctlArgsT{&ctlArgsT{
											requestType: 0x21, request: 9, value: 0x300, index: 1,
											data: []byte{0x8, 0x2, 0xA, ex.Speed(), ex.Brightness(), ex.ColorNum(), 0,
												ex.Save()}, length: 8, timeout: 0,
										}})

									assertFindDeviceCall(findDevCall, ex.Device(), ex.DeviceBus(), ex.DeviceAddress(),
										ex.PollInterval(), ex.PollTimeout())

									assertReadConfigCall(readConfigCall, ex.ConfigFile())
								})
							},
							newExecs().
								ConfigFile(configFileAll...).
								PollInterval(pollIntervalAll...).
								PollTimeout(pollTimeoutAll...).
								Brightness(brightnessAll...).
								Speed(speedAll...).
								ColorNum(colorNumAll...).
								Reset(resetAll...).
								Save(saveAll...).entries(),
							newExecs().
								RequiredDeviceBus(deviceBusAll...).
								RequiredDeviceAddress(deviceAddressAll...).entries(),
						)

						DescribeTableSubtree("random-mode",
							func(ex *execT) {

								BeforeEach(func() {
									subCmd, cmdArgs = "random-mode", ex.genArgs(nil, defs, r)
								})

								It("correctly calls usb device", func() {
									assertOnlyControlCallWithReset(dev, ex.PredefinedColors(), ex.Reset(),
										[]*ctlArgsT{&ctlArgsT{
											requestType: 0x21, request: 9, value: 0x300, index: 1,
											data: []byte{0x8, 0x2, 0x4, ex.Speed(), ex.Brightness(), ex.ColorNum(),
												ex.Reactive(), ex.Save()}, length: 8, timeout: 0,
										}})

									assertFindDeviceCall(findDevCall, ex.Device(), ex.DeviceBus(), ex.DeviceAddress(),
										ex.PollInterval(), ex.PollTimeout())

									assertReadConfigCall(readConfigCall, ex.ConfigFile())
								})
							},
							newExecs().
								ConfigFile(configFileAll...).
								PollInterval(pollIntervalAll...).
								PollTimeout(pollTimeoutAll...).
								Brightness(brightnessAll...).
								Speed(speedAll...).
								ColorNum(colorNumAll...).
								Reactive(reactiveAll...).
								Reset(resetAll...).
								Save(saveAll...).entries(),
							newExecs().
								RequiredDeviceBus(deviceBusAll...).
								RequiredDeviceAddress(deviceAddressAll...).entries(),
						)

						DescribeTableSubtree("ripple-mode",
							func(ex *execT) {

								BeforeEach(func() {
									subCmd, cmdArgs = "ripple-mode", ex.genArgs(nil, defs, r)
								})

								It("correctly calls usb device", func() {
									assertOnlyControlCallWithReset(dev, ex.PredefinedColors(), ex.Reset(),
										[]*ctlArgsT{&ctlArgsT{
											requestType: 0x21, request: 9, value: 0x300, index: 1,
											data: []byte{0x8, 0x2, 0x6, ex.Speed(), ex.Brightness(), ex.ColorNum(),
												ex.Reactive(), ex.Save()}, length: 8, timeout: 0,
										}})

									assertFindDeviceCall(findDevCall, ex.Device(), ex.DeviceBus(), ex.DeviceAddress(),
										ex.PollInterval(), ex.PollTimeout())

									assertReadConfigCall(readConfigCall, ex.ConfigFile())
								})
							},
							newExecs().
								ConfigFile(configFileAll...).
								PollInterval(pollIntervalAll...).
								PollTimeout(pollTimeoutAll...).
								Brightness(brightnessAll...).
								Speed(speedAll...).
								ColorNum(colorNumAll...).
								Reactive(reactiveAll...).
								Reset(resetAll...).
								Save(saveAll...).entries(),
							newExecs().
								RequiredDeviceBus(deviceBusAll...).
								RequiredDeviceAddress(deviceAddressAll...).entries(),
						)

						DescribeTableSubtree("wave-mode",
							func(ex *execT) {

								BeforeEach(func() {
									subCmd, cmdArgs = "wave-mode", ex.genArgs(nil, defs, r)
								})

								It("correctly calls usb device", func() {
									assertOnlyControlCallWithReset(dev, ex.PredefinedColors(), ex.Reset(),
										[]*ctlArgsT{&ctlArgsT{
											requestType: 0x21, request: 9, value: 0x300, index: 1,
											data: []byte{0x8, 0x2, 0x3, ex.Speed(), ex.Brightness(), 0, ex.Direction(),
												ex.Save()}, length: 8, timeout: 0,
										}})

									assertFindDeviceCall(findDevCall, ex.Device(), ex.DeviceBus(), ex.DeviceAddress(),
										ex.PollInterval(), ex.PollTimeout())

									assertReadConfigCall(readConfigCall, ex.ConfigFile())
								})
							},
							newExecs().
								ConfigFile(configFileAll...).
								PollInterval(pollIntervalAll...).
								PollTimeout(pollTimeoutAll...).
								Brightness(brightnessAll...).
								Speed(speedAll...).
								Direction(directionAll...).
								Reset(resetAll...).
								Save(saveAll...).entries(),
							newExecs().
								ConfigFile(configFileAll...).
								RequiredDeviceBus(deviceBusAll...).
								RequiredDeviceAddress(deviceAddressAll...).entries(),
						)

						DescribeTableSubtree("single-color-mode",
							func(ex *execT) {

								BeforeEach(func() {
									subCmd, cmdArgs = "single-color-mode", ex.genArgs(nil, defs, r)
								})

								It("correctly calls usb device", func() {

									expectedCallArgs := []*ctlArgsT{
										&ctlArgsT{
											requestType: 0x21, request: 9, value: 0x300, index: 1,
											data: []byte{0x8, 0x2, 0x33, 0, ex.Brightness(), 0, 0,
												ex.Save()}, length: 8, timeout: 0,
										}}

									for i := range 6 {
										expectedCallArgs = append(expectedCallArgs,
											&ctlArgsT{
												requestType: 0x21, request: 9, value: 0x300, index: 1,
												data: []byte{0x16, 0x0, byte(i)}, length: 3, timeout: 0,
											})
									}

									expectedCallArgs = assertControlCall(dev, ex.PredefinedColors(), ex.Reset(), expectedCallArgs)
									assertGetBulkWriteCall(dev, 1)
									assertCloseCallAfter(dev, len(expectedCallArgs), 1, ite8291.RowsNumber)

									col := ex.SingleModeColor()
									z, r, g, b := []byte{0}, []byte{col.Red}, []byte{col.Green}, []byte{col.Blue}
									expectedRow := slices.Concat(z, bytes.Repeat(b, ite8291.ColumnsNumber),
										bytes.Repeat(g, ite8291.ColumnsNumber),
										bytes.Repeat(r, ite8291.ColumnsNumber), z)

									Ω(dev.bulkBuffer.Contents()).Should(Equal(bytes.Repeat(expectedRow, ite8291.RowsNumber)))

									assertFindDeviceCall(findDevCall, ex.Device(), ex.DeviceBus(), ex.DeviceAddress(),
										ex.PollInterval(), ex.PollTimeout())

									assertReadConfigCall(readConfigCall, ex.ConfigFile())
								})
							},
							newExecs().
								ConfigFile(configFileAll...).
								PollInterval(pollIntervalAll...).
								PollTimeout(pollTimeoutAll...).
								Brightness(brightnessAll...).
								Reset(resetAll...).
								Save(saveAll...).
								Red(colorAll...).
								Green(colorAll...).
								Blue(colorAll...).entries(),
							newExecs().
								RequiredRGB(rgbAll...).entries(),
							newExecs().
								RequiredColorName(colorNameAll...).entries(),
							newExecs().
								RequiredDeviceBus(deviceBusAll...).
								RequiredDeviceAddress(deviceAddressAll...).entries(),
						)
					},
					Entry("NO", false),
					Entry("YES", true),
				)
			},
			// different configurations DFC
			Entry("no", nil),

			Entry("case 1", &defaultsT{
				brightness:       40,
				speed:            2,
				colorNum:         3,
				direction:        "doWn",
				reactive:         true,
				save:             true,
				reset:            true,
				singleColor:      colorNameAll[2],
				predefinedColors: []string{"#eA2", colorNameAll[0], "0xd1c2a3", "0XBBCCAA", "#1A1B1C", "2A3B44", colorNameAll[2]},
				pollInterval:     500 * time.Millisecond,
				pollTimeout:      800 * time.Millisecond,
				deviceBus:        10,
				deviceAddress:    22}),

			Entry("case 2", &defaultsT{
				brightness:       10,
				speed:            8,
				colorNum:         2,
				direction:        "LEFt",
				reactive:         false,
				save:             true,
				reset:            true,
				singleColor:      "#115522",
				predefinedColors: []string{"", "291", "", "0XBBCCAA", "", colorNameAll[1]},
				pollInterval:     1 * time.Second,
				pollTimeout:      10 * time.Minute,
				deviceBus:        -1,
				deviceAddress:    -1}))
	})

	Describe("validation erros", func() {

		JustBeforeEach(func() {
			Ω(cmdErr).Should(HaveOccurred()) // command must return an error
		})

		DescribeTableSubtree("brightness",
			func(ex *execT) {

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {
							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)
						})

						It("fails and correctly repots the error", func() {
							assertValidationError(dev, cmdErr, cmdErrOut, ex.brightness.value, brightnessFlagAll...)
						})
					},

					EntryDescription("%q"),

					Entry(nil, "aurora-mode"),
					Entry(nil, "breath-mode"),
					Entry(nil, "fireworks-mode"),
					Entry(nil, "marquee-mode"),
					Entry(nil, "rainbow-mode"),
					Entry(nil, "raindrop-mode"),
					Entry(nil, "random-mode"),
					Entry(nil, "ripple-mode"),
					Entry(nil, "single-color-mode"),
					Entry(nil, "wave-mode"),
				)
			},
			newExecs().RequiredBrightness("-1", "51", "abcd", "1efg", "hjk2").entries(),
		)

		DescribeTableSubtree("speed",
			func(ex *execT) {

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {
							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)
						})

						It("fails and correctly repots the error", func() {
							assertValidationError(dev, cmdErr, cmdErrOut, ex.speed.value, speedFlagAll...)
						})
					},

					EntryDescription("%q"),

					Entry(nil, "aurora-mode"),
					Entry(nil, "breath-mode"),
					Entry(nil, "fireworks-mode"),
					Entry(nil, "marquee-mode"),
					Entry(nil, "raindrop-mode"),
					Entry(nil, "random-mode"),
					Entry(nil, "ripple-mode"),
					Entry(nil, "wave-mode"),
				)
			},
			newExecs().RequiredSpeed("-1", "11", "abcd", "2Speed", "speed1").entries(),
		)

		DescribeTableSubtree("color-num",
			func(ex *execT) {

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {
							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)
						})

						It("fails and correctly repots the error", func() {
							assertValidationError(dev, cmdErr, cmdErrOut, ex.colorNum.value, colorNumFlagAll...)
						})
					},

					EntryDescription("%q"),

					Entry(nil, "aurora-mode"),
					Entry(nil, "breath-mode"),
					Entry(nil, "fireworks-mode"),
					Entry(nil, "raindrop-mode"),
					Entry(nil, "random-mode"),
					Entry(nil, "ripple-mode"),
				)
			},
			newExecs().RequiredColorNum("-1", "11", "abcd", "2Speed", "speed1").entries(),
		)

		DescribeTableSubtree("custom color-num",
			func(ex *execT) {

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {
							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)
						})

						It("fails and correctly repots the error", func() {
							assertValidationError(dev, cmdErr, cmdErrOut, ex.customColorNum.value, colorNumFlagAll...)
						})
					},

					EntryDescription("%q"),

					Entry(nil, "set-color"),
				)
			},
			newExecs().
				RequiredRGB("#FFFFFF").
				CustomColorNum("0", "8", "abcd", "2Speed", "speed1").entries(),
		)

		DescribeTableSubtree("direction",
			func(ex *execT) {

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {
							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)
						})

						It("fails and correctly repots the error", func() {
							assertValidationError(dev, cmdErr, cmdErrOut, ex.direction.value, directionFlagAll...)
						})
					},

					EntryDescription("%q"),

					Entry(nil, "wave-mode"),
				)
			},
			newExecs().RequiredDirection("unknown", "11", "zz").entries(),
		)

		DescribeTableSubtree("reactive",
			func(ex *execT) {

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {
							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)
						})

						It("fails and correctly repots the error", func() {
							assertValidationError(dev, cmdErr, cmdErrOut, ex.reactive.value, reactiveFlagAll...)
						})
					},

					EntryDescription("%q"),

					Entry(nil, "aurora-mode"),
					Entry(nil, "fireworks-mode"),
					Entry(nil, "random-mode"),
					Entry(nil, "ripple-mode"),
				)
			},
			newExecs().RequiredReactive("true1", "nottrue", "2false", "nofalse", "maybe").entries(),
		)

		DescribeTableSubtree("save",
			func(ex *execT) {

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {
							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)
						})

						It("fails and correctly repots the error", func() {
							assertValidationError(dev, cmdErr, cmdErrOut, ex.save.value, saveFlagAll...)
						})

					},

					EntryDescription("%q"),

					Entry(nil, "aurora-mode"),
					Entry(nil, "breath-mode"),
					Entry(nil, "fireworks-mode"),
					Entry(nil, "marquee-mode"),
					Entry(nil, "rainbow-mode"),
					Entry(nil, "raindrop-mode"),
					Entry(nil, "random-mode"),
					Entry(nil, "ripple-mode"),
					Entry(nil, "single-color-mode"),
					Entry(nil, "wave-mode"),
				)
			},
			newExecs().RequiredSave("true1", "nottrue", "2false", "nofalse", "maybe").entries(),
		)

		DescribeTableSubtree("reset",
			func(ex *execT) {

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {
							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)
						})

						It("fails and correctly repots the error", func() {
							assertValidationError(dev, cmdErr, cmdErrOut, ex.reset.value, resetFlagAll...)
						})

					},

					EntryDescription("%q"),

					Entry(nil, "aurora-mode"),
					Entry(nil, "breath-mode"),
					Entry(nil, "fireworks-mode"),
					Entry(nil, "marquee-mode"),
					Entry(nil, "off-mode"),
					Entry(nil, "rainbow-mode"),
					Entry(nil, "raindrop-mode"),
					Entry(nil, "random-mode"),
					Entry(nil, "ripple-mode"),
					Entry(nil, "single-color-mode"),
					Entry(nil, "wave-mode"),
				)
			},
			newExecs().RequiredReset("true1", "nottrue", "2false", "nofalse", "maybe").entries(),
		)

		DescribeTableSubtree("poll-interval",
			func(ex *execT) {

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {
							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)
						})

						It("fails and correctly repots the error", func() {
							assertValidationError(dev, cmdErr, cmdErrOut, ex.pollInterval.value, pollIntervalFlagAll...)
						})
					},

					EntryDescription("%q"),

					Entry(nil, "brightness"),
					Entry(nil, "firmware-version"),
					Entry(nil, "state"),

					Entry(nil, "aurora-mode"),
					Entry(nil, "breath-mode"),
					Entry(nil, "fireworks-mode"),
					Entry(nil, "marquee-mode"),
					Entry(nil, "off-mode"),
					Entry(nil, "rainbow-mode"),
					Entry(nil, "raindrop-mode"),
					Entry(nil, "random-mode"),
					Entry(nil, "ripple-mode"),
					Entry(nil, "single-color-mode"),
					Entry(nil, "wave-mode"),
				)
			},
			newExecs().RequiredPollInterval("-1µs", "0s", "11c", "as", "2Sm", "22").entries(),
		)

		DescribeTableSubtree("poll-timeout",
			func(ex *execT) {

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {
							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)
						})

						It("fails and correctly repots the error", func() {
							assertValidationError(dev, cmdErr, cmdErrOut, ex.pollTimeout.value, pollTimeoutFlagAll...)
						})
					},

					EntryDescription("%q"),

					Entry(nil, "brightness"),
					Entry(nil, "firmware-version"),
					Entry(nil, "state"),

					Entry(nil, "aurora-mode"),
					Entry(nil, "breath-mode"),
					Entry(nil, "fireworks-mode"),
					Entry(nil, "marquee-mode"),
					Entry(nil, "off-mode"),
					Entry(nil, "rainbow-mode"),
					Entry(nil, "raindrop-mode"),
					Entry(nil, "random-mode"),
					Entry(nil, "ripple-mode"),
					Entry(nil, "single-color-mode"),
					Entry(nil, "wave-mode"),
				)
			},
			newExecs().RequiredPollTimeout("-1", "11c", "as", "2Sm", "22").entries(),
		)

		DescribeTableSubtree("poll-interval > poll-timeout",
			func(ex *execT) {

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {
							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)
						})

						It("fails and correctly repots the error", func() {
							assertValidationError(dev, cmdErr, cmdErrOut, ex.pollInterval.value, ex.pollTimeout.value, "interval", "timeout")
						})
					},

					EntryDescription("%q"),

					Entry(nil, "brightness"),
					Entry(nil, "firmware-version"),
					Entry(nil, "state"),

					Entry(nil, "aurora-mode"),
					Entry(nil, "breath-mode"),
					Entry(nil, "fireworks-mode"),
					Entry(nil, "marquee-mode"),
					Entry(nil, "off-mode"),
					Entry(nil, "rainbow-mode"),
					Entry(nil, "raindrop-mode"),
					Entry(nil, "random-mode"),
					Entry(nil, "ripple-mode"),
					Entry(nil, "single-color-mode"),
					Entry(nil, "wave-mode"),
				)
			},
			newExecs().
				RequiredPollInterval("10s").
				RequiredPollTimeout("9s").entries(),
		)

		DescribeTableSubtree("device-bus",
			func(ex *execT) {

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {
							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)
						})

						It("fails and correctly repots the error", func() {
							assertValidationError(dev, cmdErr, cmdErrOut, ex.deviceBus.value, "device", "bus")
						})
					},

					EntryDescription("%q"),

					Entry(nil, "brightness"),
					Entry(nil, "firmware-version"),
					Entry(nil, "state"),

					Entry(nil, "aurora-mode"),
					Entry(nil, "breath-mode"),
					Entry(nil, "fireworks-mode"),
					Entry(nil, "marquee-mode"),
					Entry(nil, "off-mode"),
					Entry(nil, "rainbow-mode"),
					Entry(nil, "raindrop-mode"),
					Entry(nil, "random-mode"),
					Entry(nil, "ripple-mode"),
					Entry(nil, "single-color-mode"),
					Entry(nil, "wave-mode"),
				)
			},
			newExecs().
				DeviceBus("bus", "11bus", "as2").
				RequiredDeviceAddress(deviceAddressAll...).entries(),
		)

		DescribeTableSubtree("device-address",
			func(ex *execT) {

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {
							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)
						})

						It("fails and correctly repots the error", func() {
							assertValidationError(dev, cmdErr, cmdErrOut, ex.deviceAddress.value, "device", "address")
						})
					},

					EntryDescription("%q"),

					Entry(nil, "brightness"),
					Entry(nil, "firmware-version"),
					Entry(nil, "state"),

					Entry(nil, "aurora-mode"),
					Entry(nil, "breath-mode"),
					Entry(nil, "fireworks-mode"),
					Entry(nil, "marquee-mode"),
					Entry(nil, "off-mode"),
					Entry(nil, "rainbow-mode"),
					Entry(nil, "raindrop-mode"),
					Entry(nil, "random-mode"),
					Entry(nil, "ripple-mode"),
					Entry(nil, "single-color-mode"),
					Entry(nil, "wave-mode"),
				)
			},
			newExecs().
				RequiredDeviceBus(deviceBusAll...).
				DeviceAddress("1a", "x010", "address", "a1").entries(),
		)

		DescribeTableSubtree("red",
			func(ex *execT) {

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {
							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)

							if command == "set-color" {
								cmdArgs = append(cmdArgs, "-c", "1")
							}
						})

						It("fails and correctly repots the error", func() {
							assertValidationError(dev, cmdErr, cmdErrOut, ex.red.value, redFlagAll...)
						})
					},

					EntryDescription("%q"),

					Entry(nil, "set-color"),
					Entry(nil, "single-color-mode"),
				)
			},
			newExecs().RequiredRed("-1", "1end", "start2", "NaN", "1.1", "1.0", "0.0", "256").entries(),
		)

		DescribeTableSubtree("green",
			func(ex *execT) {

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {
							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)

							if command == "set-color" {
								cmdArgs = append(cmdArgs, "-c", "1")
							}
						})

						It("fails and correctly repots the error", func() {
							assertValidationError(dev, cmdErr, cmdErrOut, ex.green.value, greenFlagAll...)
						})

					},
					EntryDescription("%q"),

					Entry(nil, "set-color"),
					Entry(nil, "single-color-mode"),
				)
			},
			newExecs().RequiredGreen("-1", "1end", "start2", "NaN", "1.1", "1.0", "0.0", "256").entries(),
		)

		DescribeTableSubtree("blue",
			func(ex *execT) {

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {
							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)

							if command == "set-color" {
								cmdArgs = append(cmdArgs, "-c", "1")
							}
						})

						It("fails and correctly repots the error", func() {
							assertValidationError(dev, cmdErr, cmdErrOut, ex.blue.value, blueFlagAll...)
						})

					},

					EntryDescription("%q"),

					Entry(nil, "set-color"),
					Entry(nil, "single-color-mode"),
				)
			},
			newExecs().RequiredBlue("-1", "1end", "start2", "NaN", "1.1", "1.0", "0.0", "256").entries(),
		)

		DescribeTableSubtree("rgb",
			func(ex *execT) {

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {
							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)

							if command == "set-color" {
								cmdArgs = append(cmdArgs, "-c", "1")
							}
						})

						It("fails and correctly repots the error", func() {
							assertValidationError(dev, cmdErr, cmdErrOut, ex.rgb.value, rgbFlagAll...)
						})
					},

					EntryDescription("%q"),

					Entry(nil, "set-color"),
					Entry(nil, "single-color-mode"),
				)
			},
			newExecs().RequiredRGB("1234", "12", "#78", "#5678", "0x1234", "0xab",
				"0X1234567", "0xabcde", "#1a2b3c4", "0xa1b2c3d", "#bcde").entries(),
		)

		DescribeTableSubtree("color-name",
			func(ex *execT) {

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {
							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)

							if command == "set-color" {
								cmdArgs = append(cmdArgs, "-c", "1")
							}

							Ω(readConfigCall.v.MergeConfigMap(map[string]any{
								"namedColors": map[string]any{
									"existing-color1": "1234",
									"existing-color2": "12",
									"existing-color3": "#78",
									"existing-color4": "#5678",
									"existing-color6": "0x1234",
									"existing-color7": "0X1234567",
									"existing-color8": "0xabcde",
									"existing-color9": "#1a2b3c4",
									"existing-colora": "0xa1b2c3d",
									"existing-colorb": "#bcde",
									"existing-colorc": "0xab",
								},
							})).Should(Succeed())
						})

						It("fails and correctly repots the error", func() {
							assertValidationError(dev, cmdErr, cmdErrOut, ex.colorName.value, colorNameFlagAll...)
						})
					},

					EntryDescription("%q"),

					Entry(nil, "set-color"),
					Entry(nil, "single-color-mode"),
				)
			},
			newExecs().RequiredColorName("NONAME-Color1", "NNONAME-Color", "1234", "12",
				"existing-color1", "existing-color2", "existing-color3",
				"existing-color4", "existing-color5", "existing-color6",
				"existing-color7", "existing-color8", "existing-color9",
				"existing-colora", "existing-colorb", "existing-colorc").entries(),
		)

		DescribeTableSubtree("predefinedColors",
			func(ex *execT) {

				var colorNum int

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {
							colorNum = r.Intn(7) + 1

							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)

							newConf := map[string]any{
								params.PredefinedColorProp: map[string]any{
									fmt.Sprintf("color%d", colorNum): "#S1FF55",
								},
							}

							Ω(readConfigCall.v.MergeConfigMap(newConf)).Should(Succeed()) // merge with viper config

						})

						It("fails and correctly repots the error", func() {
							Ω(cmdErr).Should(MatchError(SatisfyAll(
								ContainSubstring("predefined"),
								ContainSubstring("color"),
								ContainSubstring(strconv.Itoa(colorNum)))))
						})
					},

					EntryDescription("%q"),

					Entry(nil, "aurora-mode"),
					Entry(nil, "breath-mode"),
					Entry(nil, "fireworks-mode"),
					Entry(nil, "marquee-mode"),
					Entry(nil, "rainbow-mode"),
					Entry(nil, "raindrop-mode"),
					Entry(nil, "random-mode"),
					Entry(nil, "ripple-mode"),
					Entry(nil, "single-color-mode"),
					Entry(nil, "wave-mode"),
				)
			},
			newExecs().RequiredReset("true").entries(),
		)

		DescribeTableSubtree("singleModeColors",
			func(ex *execT) {

				DescribeTableSubtree("with command",
					func(command string) {

						BeforeEach(func() {

							subCmd, cmdArgs = command, ex.genArgs(nil, zeroDefaults(), r)

							newConf := map[string]any{
								"singleModeColor": "#FF55",
							}

							Ω(readConfigCall.v.MergeConfigMap(newConf)).Should(Succeed()) // merge with viper config

						})

						It("fails and correctly repots the error", func() {
							Ω(cmdErr).Should(MatchError(SatisfyAll(
								ContainSubstring("single"),
								ContainSubstring("mode"),
								ContainSubstring("color"),
								ContainSubstring("#FF55"))))
						})
					},

					EntryDescription("%q"),

					Entry(nil, "single-color-mode"),
				)
			},
			newExecs().entries(),
		)
	})

	Describe("device errors", func() {

		JustBeforeEach(func() {
			Ω(cmdErr).Should(HaveOccurred()) // assert command returns an error
		})

		BeforeEach(func() {
			uuid.SetRand(r) // seed uuid generation
		})

		DescribeTableSubtree("in ControlTransfer method",
			func(callNum int, callCmd string, callArgs ...string) {

				BeforeEach(func() {
					dev.ctlRtnError = fmt.Errorf("control transfer error %s", uuid.New())
					dev.ctlRtnErrorAtCall = callNum - 1
					subCmd, cmdArgs = callCmd, callArgs
				})

				It("fails and correctly repots the error", func() {
					Ω(cmdErr).Should(MatchError(dev.ctlRtnError))
				})
			},

			func(n int, c string, p ...string) string {
				sb := &strings.Builder{}
				fmt.Fprintf(sb, "on %d method invocation in cmd: %s", n, c)
				for _, s := range p {
					sb.WriteRune(' ')
					sb.WriteString(s)
				}

				return sb.String()
			},

			Entry(nil, 1, "aurora-mode"),
			Entry(nil, 1, "breath-mode"),
			Entry(nil, 1, "fireworks-mode"),
			Entry(nil, 1, "brightness"),
			Entry(nil, 1, "firmware-version"),
			Entry(nil, 1, "marquee-mode"),
			Entry(nil, 1, "off-mode"),
			Entry(nil, 1, "rainbow-mode"),
			Entry(nil, 1, "raindrop-mode"),
			Entry(nil, 1, "random-mode"),
			Entry(nil, 1, "ripple-mode"),
			Entry(nil, 1, "set-brightness", "-b", "20"),
			Entry(nil, 1, "set-color", "-c", "2", "--rgb", "#11AAFF"),
			Entry(nil, 1, "single-color-mode"),
			Entry(nil, 1, "state"),
			Entry(nil, 1, "wave-mode"),

			Entry(nil, 1, "aurora-mode", "--reset"),
			Entry(nil, 1, "breath-mode", "--reset"),
			Entry(nil, 1, "fireworks-mode", "--reset"),
			Entry(nil, 1, "marquee-mode", "--reset"),
			Entry(nil, 1, "off-mode", "--reset"),
			Entry(nil, 1, "rainbow-mode", "--reset"),
			Entry(nil, 1, "raindrop-mode", "--reset"),
			Entry(nil, 1, "random-mode", "--reset"),
			Entry(nil, 1, "ripple-mode", "--reset"),
			Entry(nil, 1, "single-color-mode", "--reset"),
			Entry(nil, 1, "wave-mode", "--reset"),

			Entry(nil, 2, "brightness"),
			Entry(nil, 2, "firmware-version"),
			Entry(nil, 2, "single-color-mode"),
			Entry(nil, 2, "state"),

			Entry(nil, 8, "aurora-mode", "--reset"),
			Entry(nil, 8, "breath-mode", "--reset"),
			Entry(nil, 8, "fireworks-mode", "--reset"),
			Entry(nil, 8, "marquee-mode", "--reset"),
			Entry(nil, 8, "off-mode", "--reset"),
			Entry(nil, 8, "rainbow-mode", "--reset"),
			Entry(nil, 8, "raindrop-mode", "--reset"),
			Entry(nil, 8, "random-mode", "--reset"),
			Entry(nil, 8, "ripple-mode", "--reset"),
			Entry(nil, 8, "single-color-mode", "--reset"),
			Entry(nil, 8, "wave-mode", "--reset"),
		)

		DescribeTableSubtree("in GetBulkWrite method",
			func(callNum int, callCmd string, callArgs ...string) {

				BeforeEach(func() {
					dev.getBulkWriteRtnError = fmt.Errorf("get bulk write error %s", uuid.New())
					dev.getBulkWriteRtnErrorAtCall = callNum
					subCmd, cmdArgs = callCmd, callArgs
				})

				It("fails and correctly repots the error", func() {
					Ω(cmdErr).Should(MatchError(dev.getBulkWriteRtnError))
				})
			},

			func(n int, c string, p ...string) string {
				sb := &strings.Builder{}
				fmt.Fprintf(sb, "on %d method invocation in cmd: %s", n, c)
				for _, s := range p {
					sb.WriteRune(' ')
					sb.WriteString(s)
				}

				return sb.String()
			},

			Entry(nil, 1, "single-color-mode"),
		)

		DescribeTableSubtree("in bulk Write method",
			func(callNum int, callCmd string, callArgs ...string) {

				BeforeEach(func() {
					dev.bulkWriteRtnError = fmt.Errorf("get bulk write error %s", uuid.New())
					dev.bulkWriteRtnErrorAtCall = callNum
					subCmd, cmdArgs = callCmd, callArgs
				})

				It("fails and correctly repots the error", func() {
					Ω(cmdErr).Should(MatchError(dev.bulkWriteRtnError))
				})
			},

			func(n int, c string, p ...string) string {
				sb := &strings.Builder{}
				fmt.Fprintf(sb, "on %d method invocation in cmd: %s", n, c)
				for _, s := range p {
					sb.WriteRune(' ')
					sb.WriteString(s)
				}

				return sb.String()
			},

			Entry(nil, 1, "single-color-mode"),
		)

		DescribeTableSubtree("in find device function",
			func(callCmd string, callArgs ...string) {

				BeforeEach(func() {
					findDevCall.rtnError = fmt.Errorf("find device error %s", uuid.New())
					subCmd, cmdArgs = callCmd, callArgs
				})

				It("fails and correctly repots the error", func() {
					Ω(cmdErr).Should(MatchError(findDevCall.rtnError))
				})
			},

			func(c string, p ...string) string {
				sb := &strings.Builder{}
				fmt.Fprintf(sb, "in cmd: %s", c)
				for _, s := range p {
					sb.WriteRune(' ')
					sb.WriteString(s)
				}

				return sb.String()
			},

			Entry(nil, "aurora-mode"),
			Entry(nil, "breath-mode"),
			Entry(nil, "fireworks-mode"),
			Entry(nil, "brightness"),
			Entry(nil, "firmware-version"),
			Entry(nil, "marquee-mode"),
			Entry(nil, "off-mode"),
			Entry(nil, "rainbow-mode"),
			Entry(nil, "raindrop-mode"),
			Entry(nil, "random-mode"),
			Entry(nil, "ripple-mode"),
			Entry(nil, "set-brightness", "-b", "20"),
			Entry(nil, "set-color", "-c", "2", "--rgb", "#11AAFF"),
			Entry(nil, "single-color-mode"),
			Entry(nil, "state"),
			Entry(nil, "wave-mode"),
		)

		DescribeTableSubtree("in read config function",
			func(callCmd string, callArgs ...string) {

				BeforeEach(func() {
					readConfigCall.rtnError = fmt.Errorf("read config error %s", uuid.New())
					subCmd, cmdArgs = callCmd, callArgs
				})

				It("fails and correctly repots the error", func() {
					Ω(cmdErr).Should(MatchError(readConfigCall.rtnError))
				})
			},

			func(c string, p ...string) string {
				sb := &strings.Builder{}
				fmt.Fprintf(sb, "in cmd: %s", c)
				for _, s := range p {
					sb.WriteRune(' ')
					sb.WriteString(s)
				}

				return sb.String()
			},

			Entry(nil, "aurora-mode"),
			Entry(nil, "breath-mode"),
			Entry(nil, "fireworks-mode"),
			Entry(nil, "brightness"),
			Entry(nil, "firmware-version"),
			Entry(nil, "marquee-mode"),
			Entry(nil, "off-mode"),
			Entry(nil, "rainbow-mode"),
			Entry(nil, "raindrop-mode"),
			Entry(nil, "random-mode"),
			Entry(nil, "ripple-mode"),
			Entry(nil, "set-brightness", "-b", "20"),
			Entry(nil, "set-color", "-c", "2", "--rgb", "#11AAFF"),
			Entry(nil, "single-color-mode"),
			Entry(nil, "state"),
			Entry(nil, "wave-mode"),
		)
	})
})
