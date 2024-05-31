package cmd

import (
	"slices"
	"time"

	"github.com/onsi/gomega/gbytes"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/pkg/ite8291"

	. "github.com/onsi/gomega"
)

// ctlArgsT type provides data to collect and assert ControlTransfer arguments.
type ctlArgsT struct {
	requestType byte
	request     byte
	value       uint16
	index       uint16
	data        []byte
	length      int
	timeout     int
}

// deviceStubT type provides stub to use instead of ite8291 controller to collect and assert its calls.
type deviceStubT struct {
	ctlArgs           []*ctlArgsT
	ctlCallNum        int
	ctlChangedData    [][]byte
	ctlRtnError       error
	ctlRtnErrorAtCall int

	closeCallNum                int
	closePreCtlCallNum          int
	closePreGetBulkWriteCallNum int
	closePreBulkWriteCallNum    int

	getBulkWriteCallNum        int
	getBulkWriteRtnError       error
	getBulkWriteRtnErrorAtCall int

	bulkWriteCallNum        int
	bulkWriteRtnError       error
	bulkWriteRtnErrorAtCall int

	bulkBuffer *gbytes.Buffer
}

// ControlTransfer intercepts controller ControlTransfer calls.
func (d *deviceStubT) ControlTransfer(requestType byte, request byte, value uint16, index uint16,
	data []byte, length int, timeout int) (int, error) {

	// return a given d.ctlRtnError at d.ctlRtnErrorAtCall
	if d.ctlRtnError != nil && d.ctlRtnErrorAtCall == d.ctlCallNum {
		return 0, d.ctlRtnError
	}

	// collect call args
	d.ctlArgs = append(d.ctlArgs, &ctlArgsT{
		requestType: requestType,
		request:     request,
		value:       value,
		index:       index,
		data:        slices.Clone(data),
		length:      length,
		timeout:     timeout,
	})

	// change data if requested by d.ctlChangedData
	if len(d.ctlChangedData) > d.ctlCallNum && len(d.ctlChangedData[d.ctlCallNum]) > 0 {
		copy(data, d.ctlChangedData[d.ctlCallNum])
	}
	d.ctlCallNum++ // increase call counter

	return len(data), nil
}

// GetBulkWrite intercepts controller GetBulkWrite calls.
func (d *deviceStubT) GetBulkWrite() (ite8291.WriteFunc, error) {
	d.getBulkWriteCallNum++ // increase call counter
	// init bulk buffer if it doesn't exist
	if d.bulkBuffer == nil {
		d.bulkBuffer = gbytes.NewBuffer()
	}

	// return given d.getBulkWriteRtnError at d.getBulkWriteRtnErrorAtCall call
	if d.getBulkWriteRtnError != nil && d.getBulkWriteRtnErrorAtCall == d.getBulkWriteCallNum {
		return nil, d.getBulkWriteRtnError
	}

	return func(data []byte) (n int, err error) {
		d.bulkWriteCallNum++              // increase call counter
		n, err = d.bulkBuffer.Write(data) // write data

		// return given d.getBulkWriteRtnError at d.getBulkWriteRtnErrorAtCall call
		if d.bulkWriteRtnError != nil && d.bulkWriteRtnErrorAtCall == d.bulkWriteCallNum {
			return 0, d.bulkWriteRtnError
		}

		return n, err
	}, nil
}

// Close intercepts controller Close calls.
func (d *deviceStubT) Close() error {
	// collect other methods call counters
	d.closePreCtlCallNum = d.ctlCallNum
	d.closePreGetBulkWriteCallNum = d.getBulkWriteCallNum
	d.closePreBulkWriteCallNum = d.bulkWriteCallNum

	d.closeCallNum++ // increase call counters

	// close bulk buffer if it exists
	if d.bulkBuffer != nil {
		Ω(d.bulkBuffer.Close()).Should(Succeed())
	}

	return nil
}

type findDeviceCallT struct {
	callNum                   int
	useDevice                 bool
	devBus, devAddress        int
	pollInterval, pollTimeout time.Duration
	rtnError                  error
}

// newFindDevice function used to return given dev controller stub as controller to collect and assert findDevice calls.
func newFindDevice(dev *deviceStubT, call *findDeviceCallT) func(bool, int, int,
	time.Duration, time.Duration) (ite8291.Device, error) {
	return func(useDevice bool, devBus, devAddress int, pollInterval, pollTimeout time.Duration) (ite8291.Device, error) {

		call.callNum++
		call.useDevice = useDevice
		call.devBus = devBus
		call.devAddress = devAddress
		call.pollInterval = pollInterval
		call.pollTimeout = pollTimeout

		// return call.rtnError if requested
		if call.rtnError != nil {
			return nil, call.rtnError
		}

		return dev, nil
	}
}

type readConfigCallT struct {
	callNum  int
	cfgFile  string
	v        *viper.Viper
	rtnError error
}

// newFindDevice function used to return given dev controller stub as controller to collect and assert findDevice calls.
func newReadConfig(call *readConfigCallT) func(*cobra.Command, *viper.Viper, string) error {

	return func(cmd *cobra.Command, v *viper.Viper, cfgFile string) error {
		Ω(v.MergeConfigMap(call.v.AllSettings())).Should(Succeed())

		call.callNum++ // increase call counter

		call.cfgFile = cfgFile

		// return call.rtnError if requested
		if call.rtnError != nil {
			return call.rtnError
		}

		return nil
	}
}
