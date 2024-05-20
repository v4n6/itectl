package cmd

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// assertOnlyControlCall asserts that controller ControlTransfer was called with given args.
// It asserts that bulk write methods were not called.
// It asserts that Close was called correctly.
func assertOnlyControlCall(device *deviceStubT, args []*ctlArgsT) {
	GinkgoHelper()

	assertOnlyControlCallWithReset(device, nil, false, args)
}

// assertOnlyControlCallWithReset asserts that controller ControlTransfer was called with given args.
// It asserts that controller was requested to reset preconfigured colors before the main controller call.
// It asserts that bulk write methods were not called.
// It asserts that Close was called correctly.
func assertOnlyControlCallWithReset(device *deviceStubT, predefinedColors []string, doReset bool, args []*ctlArgsT) {
	GinkgoHelper()

	args = assertControlCall(device, predefinedColors, doReset, args)
	assertCloseCallAfter(device, len(args), 0, 0)
	assertGetBulkWriteCall(device, 0)
	assertBulkWriteCall(device, 0)
}

// assertControlCall asserts that controller ControlTransfer was called with given args.
// doReset decides whether reset predefined colors calls shoud be asserted.
// predefinedColors specifies predefined color values.
func assertControlCall(dev *deviceStubT, predefinedColors []string, doReset bool, args []*ctlArgsT) []*ctlArgsT {
	GinkgoHelper()

	if doReset {

		l := ite8291.CustomColorNumMaxValue - ite8291.CustomColorNumMinValue + 1

		ctlResetArgs := make([]*ctlArgsT, l)

		for i := range l {
			color := predefinedColors[i]

			s, ok := namedColors[color]
			if !ok {
				s = color
			}

			val, err := ite8291.ParseColor(s)
			Ω(err).Should(Succeed())

			ctlResetArgs[i] = &ctlArgsT{
				requestType: 0x21, request: 9, value: 0x300, index: 1,
				data: []byte{0x14, 0x0, byte(i + 1), val.Red, val.Green, val.Blue}, length: 6, timeout: 0,
			}
		}

		args = append(ctlResetArgs, args...)
	}

	Ω(dev.ctlCallNum).Should(Equal(len(args)), "invalid number of control calls")

	for i, arg := range args {
		Ω(dev.ctlArgs[i]).Should(Equal(arg), "invalid control call %d args", i+1)
	}

	return args
}

// assertGetBulkWriteCall asserts that GetBulkWrite method was called correctly
func assertGetBulkWriteCall(dev *deviceStubT, getBulkWriteCallNum int) {
	GinkgoHelper()
	Ω(dev.getBulkWriteCallNum).Should(Equal(getBulkWriteCallNum), "GetBulkWrite must be called %d times", getBulkWriteCallNum)
}

// assertBulkWriteCall asserts that Write function returned by GetBulkWrite was called correctly
func assertBulkWriteCall(dev *deviceStubT, bulkWriteCallNum int) {
	GinkgoHelper()
	Ω(dev.bulkWriteCallNum).Should(Equal(bulkWriteCallNum), "GetBulkWrite must be called %d times", bulkWriteCallNum)
}

// assertCommandOutput asserts that given msg was written to stdout and nothing was written to stderr
func assertCommandOutput(out, err *gbytes.Buffer, msg string) {
	Ω(err.Contents()).Should(BeEmpty())
	Ω(out.Contents()).Should(ContainSubstring(msg))
}

// assertCloseCallAfter asserts Close was called after given number of ControlTransfer, GetBulkWrite, and Write calls.
func assertCloseCallAfter(dev *deviceStubT, ctlCallNum, getBulkWriterCallNum, bulkWriterCallNum int) {
	Ω(dev.closeCallNum).Should(Equal(1), "invalid number of close calls")
	Ω(dev.closePreCtlCallNum).Should(Equal(ctlCallNum), "invalid number of control calls preceding close call")
	Ω(dev.closePreGetBulkWriteCallNum).Should(Equal(getBulkWriterCallNum), "invalid number of get bulk write calls preceding close call")
	Ω(dev.closePreBulkWriteCallNum).Should(Equal(bulkWriterCallNum), "invalid number of bulk write calls preceding close call")
}

// assertDeviceNotCalled asserts no controller method was called.
func assertDeviceNotCalled(dev *deviceStubT) {
	Ω(dev.ctlCallNum).Should(Equal(0), "control must not be called")
	Ω(dev.getBulkWriteCallNum).Should(Equal(0), "get bulk write call must not be called")
	Ω(dev.bulkWriteCallNum).Should(Equal(0), "bulk write call must not be called")
	Ω(dev.closeCallNum).Should(Equal(0), "close must not be called")
}

// assertValidationError asserts the error returned by command and messages written to stderr.
func assertValidationError(dev *deviceStubT, cmdErr error, err *gbytes.Buffer, value string, flags ...string) {
	assertDeviceNotCalled(dev)

	matchers := []OmegaMatcher{}
	if len(value) > 0 {
		matchers = append(matchers, ContainSubstring("%q", value))
	}

	for _, flag := range flags {
		matchers = append(matchers, ContainSubstring(flag))
	}

	Ω(cmdErr).Should(MatchError(SatisfyAll(matchers...)))
	Ω(err.Contents()).Should(BeEmpty())
}

// assertFindDeviceCall asserts that findDevice function was called only one with given arguments.
func assertFindDeviceCall(call *findDeviceCallT, useDevice bool, devBus, devAddress int,
	pollInterval, pollTimeout time.Duration) {

	Ω(call.callNum).Should(Equal(1))

	Ω(call.pollInterval).Should(Equal(pollInterval))
	Ω(call.pollTimeout).Should(Equal(pollTimeout))

	Ω(call.useDevice).Should(Equal(useDevice))
	if useDevice {
		Ω(call.devBus).Should(Equal(devBus))
		Ω(call.devAddress).Should(Equal(devAddress))
	}
}

// assertReadConfigCall asserts that readConfig function was called only one with given arguments.
func assertReadConfigCall(call *readConfigCallT, cfgFile string) {

	Ω(call.callNum).Should(Equal(1))

	Ω(call.cfgFile).Should(Equal(cfgFile))
}
