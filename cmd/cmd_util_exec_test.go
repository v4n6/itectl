package cmd

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/v4n6/itectl/params"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// randomLabel is label to mark random value in tests description
const randomLabel = "<rnd>"

// true/false string representations
const (
	trueStr  = "true"
	falseStr = "false"
)

// namedColors configured named colors
var namedColors = map[string]string{
	"NONAME-Color": "ABC",
	"colour_cyAN":  "#112233",
	"123":          "0XddEEFf",
}

// named colors viper config
var namedColorsConfig = map[string]any{}

func init() {
	for n, c := range namedColors {
		namedColorsConfig[n] = c
	}
}

func init() {
	for k := range namedColors {
		colorNameAll = append(colorNameAll, k)
	}
}

// all possible flags
var (
	brightnessFlagAll     = []string{"-" + params.BrightnessShortFlag, "--" + params.BrightnessProp}
	speedFlagAll          = []string{"-" + params.SpeedShortFlag, "--" + params.SpeedProp}
	colorNumFlagAll       = []string{"-" + params.ColorNumShortFlag, "--" + params.ColorNumFlag}
	customColorNumFlagAll = []string{"-" + params.ColorNumShortFlag, "--" + params.ColorNumFlag}
	directionFlagAll      = []string{"-" + params.DirectionShortFlag, "--" + params.DirectionProp}
	reactiveFlagAll       = []string{"--" + params.ReactiveProp}

	saveFlagAll = []string{"--" + params.SaveProp}

	redFlagAll   = []string{"--" + params.ColorRedFlag}
	greenFlagAll = []string{"--" + params.ColorGreenFlag}
	blueFlagAll  = []string{"--" + params.ColorBlueFlag}

	colorNameFlagAll = []string{"--" + params.ColorNameFlag}
	rgbFlagAll       = []string{"--" + params.ColorRGBFlag}

	resetFlagAll = []string{"--" + params.ResetProp}

	pollIntervalFlagAll = []string{"--" + params.PollIntervalFlag}
	pollTimeoutFlagAll  = []string{"--" + params.PollTimeoutFlag}

	deviceBusFlagAll     = []string{"--" + params.DeviceBusFlag}
	deviceAddressFlagAll = []string{"--" + params.DeviceAddressFlag}

	configFileFlagAll = []string{"--" + params.ConfigFileFlag}
)

// all flag values to test
var (
	brightnessAll     = []string{"0", randomLabel, "50"}
	speedAll          = []string{"0", randomLabel, "10"}
	customColorNumAll = []string{"1", randomLabel, "7"}
	colorNumAll       = []string{"0", randomLabel, "8"}
	directionAll      = []string{"None", "right", "LEFT", "uP", "dOWn"}
	reactiveAll       = []string{"", trueStr, falseStr}

	colorAll     = []string{"0", randomLabel, "255"}
	rgbAll       = []string{"123567", "#a1b2c3", "0x1a1b1c", "#X2a2B2c", "#ABC", "def"}
	colorNameAll []string

	saveAll = []string{"", trueStr, falseStr}

	resetAll = []string{"", trueStr, falseStr}

	pollIntervalAll = []string{"1µs", "1s", "1m"}
	pollTimeoutAll  = []string{"20m", "10m", "1h"}

	deviceBusAll     = []string{"0", randomLabel, "10000"}
	deviceAddressAll = []string{"0", randomLabel, "20000"}

	configFileAll = []string{"", randomLabel}
)

// predefinedColorsConfig converts given colors to viper config.
func predefinedColorsConfig(colors []string) map[string]any {
	config := map[string]any{}
	for i, c := range colors {
		config[fmt.Sprintf("color%d", i+1)] = c
	}

	return config
}

// defaultsT type provides flags default/configured values to use when flag is not set.
type defaultsT struct {
	brightness byte
	speed      byte
	colorNum   byte
	direction  string
	reactive   bool
	save       bool

	reset bool

	singleColor string

	predefinedColors []string

	pollInterval time.Duration
	pollTimeout  time.Duration

	deviceBus     int
	deviceAddress int
}

// zeroDefaults creates and returns default flags values initialized with default flag values
func zeroDefaults() *defaultsT {
	return &defaultsT{
		brightness: params.BrightnessDefault,
		speed:      params.SpeedDefault,
		colorNum:   params.ColorNumDefault,
		direction:  params.DirectionDefault,
		reactive:   params.ReactiveDefault,

		save: params.ReactiveDefault,

		reset: params.ResetDefault,

		singleColor:      params.SingleColorDefault,
		predefinedColors: params.PredefinedColorsDefault,

		pollInterval: params.PollIntervalDefault,
		pollTimeout:  params.PollTimeoutDefault,

		deviceBus:     -1,
		deviceAddress: -1,
	}
}

// flag type represents cmdline flag
type flag struct {
	name  string
	value string
}

// execT type provides values to use for test cases generation and
// asserting commands executions.
type execT struct {
	conf *defaultsT

	brightness     *flag
	speed          *flag
	colorNum       *flag
	customColorNum *flag
	direction      *flag
	reactive       *flag

	red   *flag
	green *flag
	blue  *flag

	rgb       *flag
	colorName *flag

	save *flag

	reset *flag

	pollInterval *flag
	pollTimeout  *flag

	deviceBus     *flag
	deviceAddress *flag

	configFile *flag
}

// genArgs generates cmd line arguments based on execT values
func (e *execT) genArgs(args []string, conf *defaultsT, r *rand.Rand) []string {
	e.conf = conf

	args = e.genFlagArgs(args, e.brightness, func() string {
		return strconv.Itoa(r.Intn(ite8291.BrightnessMaxValue))
	})

	args = e.genFlagArgs(args, e.speed, func() string {
		return strconv.Itoa(r.Intn(ite8291.SpeedMaxValue))
	})

	args = e.genFlagArgs(args, e.colorNum, func() string {
		return strconv.Itoa(ite8291.ColorNumMinValue + r.Intn(ite8291.ColorNumMaxValue-ite8291.ColorNumMinValue))
	})

	args = e.genFlagArgs(args, e.customColorNum, func() string {
		return strconv.Itoa(ite8291.CustomColorNumMinValue + r.Intn(ite8291.CustomColorNumMaxValue-ite8291.CustomColorNumMinValue))
	})

	args = e.genFlagArgs(args, e.direction, nil)

	args = e.genFlagArgs(args, e.red, func() string {
		return strconv.Itoa(r.Intn(255))
	})

	args = e.genFlagArgs(args, e.green, func() string {
		return strconv.Itoa(r.Intn(255))
	})

	args = e.genFlagArgs(args, e.blue, func() string {
		return strconv.Itoa(r.Intn(255))
	})

	args = e.genFlagArgs(args, e.rgb, nil)

	args = e.genFlagArgs(args, e.colorName, nil)

	args = e.genBoolFlagArgs(args, e.reactive)
	args = e.genBoolFlagArgs(args, e.save)
	args = e.genBoolFlagArgs(args, e.reset)

	args = e.genFlagArgs(args, e.pollInterval, nil)
	args = e.genFlagArgs(args, e.pollTimeout, nil)

	args = e.genFlagArgs(args, e.deviceBus, func() string {
		return strconv.Itoa(r.Intn(1000))
	})
	args = e.genFlagArgs(args, e.deviceAddress, func() string {
		return strconv.Itoa(r.Intn(1000))
	})

	args = e.genFlagArgs(args, e.configFile, func() string {
		return fmt.Sprintf("config file %s", uuid.New())
	})

	return args
}

// genFlagArgs generates cmd line arguments to represent the given flag.
// If flag value is set to randomLabel, the value will be replaced with the one
// generated by randVal function.
func (e *execT) genFlagArgs(args []string, flag *flag, randVal func() string) []string {
	if flag == nil || len(flag.name) == 0 {
		return args
	}

	if flag.value == randomLabel {
		flag.value = randVal()
	}

	return append(args, flag.name, flag.value)
}

// genBoolFlagArgs generates cmd line arguments to represent the given boolean flag.
// If flag name is empty, no arg will be generated.
// If flag value is empty, only flag without the value will be generated
// If name and value are set, arg in a form "--name=value" will be generated.
func (e *execT) genBoolFlagArgs(args []string, flag *flag) []string {

	if flag == nil || len(flag.name) == 0 {
		return args
	}

	if len(flag.value) == 0 {
		return append(args, flag.name)
	}

	return append(args, fmt.Sprintf("%s=%s", flag.name, flag.value))
}

// getIntValue returns value of the given flag as int.
// If flag name is empty, specified defValue is returned.
// Flag's value is converted from string to int.
func getIntValue(flag *flag, defValue int) int {

	if flag == nil || len(flag.name) == 0 {
		return defValue
	}

	val, err := strconv.Atoi(flag.value)
	Ω(err).Should(Succeed())

	return val
}

// getByteValue returns value of the given flag as byte.
// If flag name is empty, the specified defValue is returned.
// Flag's value is converted from string to byte.
func getByteValue(flag *flag, defValue byte) byte {
	return byte(getIntValue(flag, int(defValue)))
}

// getDurationValue returns value of the given flag as time.Duration.
// If flag name is empty, the specified defValue is returned.
// Flag's value is converted from string to time.Duration.
func getDurationValue(flag *flag, defValue time.Duration) time.Duration {

	if flag == nil || len(flag.name) == 0 {
		return defValue
	}

	duration, err := time.ParseDuration(flag.value)
	Ω(err).Should(Succeed())
	return duration
}

// getByteBool returns value of the given boolean flag as byte.
// If flag name is empty, the specified defValue is used.
// If flag value is emty, true value is used.
// true is converted to 1, false to 0.
func getByteBool(flag *flag, defValue bool) byte {

	var b bool

	switch {
	case flag == nil || len(flag.name) == 0:
		b = defValue
	case len(flag.value) == 0:
		b = true
	default:
		var err error
		b, err = strconv.ParseBool(flag.value)
		Ω(err).Should(Succeed())
	}

	if b {
		return 1
	} else {
		return 0
	}
}

func (e *execT) SetBrightness(flag *flag) {
	e.brightness = flag
}

func (e *execT) Brightness() byte {
	return getByteValue(e.brightness, e.conf.brightness)
}

func (e *execT) SetSpeed(flag *flag) {
	e.speed = flag
}

func (e *execT) Speed() byte {

	if e.speed == nil || len(e.speed.name) == 0 {
		return ite8291.SpeedMaxValue - e.conf.speed
	}

	val, err := strconv.Atoi(e.speed.value)
	Ω(err).Should(Succeed())

	return ite8291.SpeedMaxValue - byte(val)
}

func (e *execT) SetColorNum(flag *flag) {
	e.colorNum = flag
}

func (e *execT) ColorNum() byte {
	return getByteValue(e.colorNum, e.conf.colorNum)
}

func (e *execT) SetCustomColorNum(flag *flag) {
	e.customColorNum = flag
}

func (e *execT) CustomColorNum() byte {
	return getByteValue(e.customColorNum, 0)
}

func (e *execT) SetDirection(flag *flag) {
	e.direction = flag
}

func (e *execT) Direction() byte {

	var dir string
	if e.direction != nil && len(e.direction.name) > 0 {
		dir = e.direction.value
	} else {
		dir = e.conf.direction
	}

	d, err := params.ParseDirectionName(dir)
	Ω(err).Should(Succeed())

	return byte(d)
}

func (e *execT) SetRed(flag *flag) {
	e.red = flag
}

func (e *execT) Red() byte {
	return getByteValue(e.red, 0)
}

func (e *execT) SetGreen(flag *flag) {
	e.green = flag
}

func (e *execT) Green() byte {
	return getByteValue(e.green, 0)
}

func (e *execT) SetBlue(flag *flag) {
	e.blue = flag
}

func (e *execT) Blue() byte {
	return getByteValue(e.blue, 0)
}

func (e *execT) SetRGB(flag *flag) {
	e.rgb = flag
}

func (e *execT) SetColorName(flag *flag) {
	e.colorName = flag
}

func (e *execT) getNamedOrRGBColor() *ite8291.Color {

	if e.colorName != nil && len(e.colorName.name) > 0 {
		s, ok := namedColors[e.colorName.value]
		Ω(ok).Should(BeTrue())

		col, err := ite8291.ParseColor(s)
		Ω(err).Should(Succeed())

		return col
	}

	if e.rgb != nil && len(e.rgb.name) > 0 {
		col, err := ite8291.ParseColor(e.rgb.value)
		Ω(err).Should(Succeed())

		return col
	}

	return nil
}

func (e *execT) Color() *ite8291.Color {

	if col := e.getNamedOrRGBColor(); col != nil {
		return col
	}

	Ω((e.red != nil && len(e.red.name) > 0) ||
		(e.green != nil && len(e.green.name) > 0) ||
		(e.blue != nil && len(e.blue.name) > 0)).Should(BeTrue())

	return &ite8291.Color{Red: e.Red(), Green: e.Green(), Blue: e.Blue()}
}

func (e *execT) SingleModeColor() *ite8291.Color {

	if col := e.getNamedOrRGBColor(); col != nil {
		return col
	}

	if (e.red != nil && len(e.red.name) > 0) ||
		(e.green != nil && len(e.green.name) > 0) ||
		(e.blue != nil && len(e.blue.name) > 0) {

		return &ite8291.Color{Red: e.Red(), Green: e.Green(), Blue: e.Blue()}
	}

	s, ok := namedColors[e.conf.singleColor]
	if !ok {
		s = e.conf.singleColor
	}

	col, err := ite8291.ParseColor(s)
	Ω(err).Should(Succeed())
	return col
}

func (e *execT) SetReset(flag *flag) {
	e.reset = flag
}

func (e *execT) Reset() bool {
	return getByteBool(e.reset, e.conf.reset) != 0
}

func (e *execT) SetPollInterval(flag *flag) {
	e.pollInterval = flag
}

func (e *execT) PollInterval() time.Duration {
	return getDurationValue(e.pollInterval, e.conf.pollInterval)
}

func (e *execT) SetPollTimeout(flag *flag) {
	e.pollTimeout = flag
}

func (e *execT) PollTimeout() time.Duration {
	return getDurationValue(e.pollTimeout, e.conf.pollTimeout)
}

func (e *execT) SetDeviceBus(flag *flag) {
	e.deviceBus = flag
}

func (e *execT) DeviceBus() int {
	return getIntValue(e.deviceBus, e.conf.deviceBus)
}

func (e *execT) SetDeviceAddress(flag *flag) {
	e.deviceAddress = flag
}

func (e *execT) DeviceAddress() int {
	return getIntValue(e.deviceAddress, e.conf.deviceAddress)
}

func (e *execT) Device() bool {
	return e.DeviceBus() >= 0 && e.DeviceAddress() >= 0
}

func (e *execT) PredefinedColors() []string {

	l, cl := ite8291.CustomColorNumMaxValue-ite8291.CustomColorNumMinValue+1, len(e.conf.predefinedColors)
	colors := make([]string, l)

	for i := range l {

		var color string
		if i < cl {
			// try configured predefined colors
			color = e.conf.predefinedColors[i]
		}

		if len(color) == 0 {
			// ith configured predefined color not found -> take default
			color = params.PredefinedColorsDefault[i]
		}

		colors[i] = color
	}

	return colors
}

func (e *execT) SetReactive(flag *flag) {
	e.reactive = flag
}

func (e *execT) Reactive() byte {
	return getByteBool(e.reactive, e.conf.reactive)
}

func (e *execT) SetSave(flag *flag) {
	e.save = flag
}

func (e *execT) Save() byte {
	return getByteBool(e.save, e.conf.save)
}

func (e *execT) SetConfigFile(flag *flag) {
	e.configFile = flag
}

func (e *execT) ConfigFile() string {
	if e.configFile == nil {
		return ""
	}
	return e.configFile.value
}

// String returns test case label
func (e *execT) String() string {
	b := &strings.Builder{}
	b.WriteString("`flags:")

	// add flag with value
	addFlag := func(flag *flag) {
		if flag != nil && len(flag.name) > 0 {
			b.WriteRune(' ')
			b.WriteString(flag.name)
			b.WriteRune(' ')
			b.WriteString(flag.value)
		}
	}

	// add boolean flag
	addBool := func(flag *flag) {
		if flag != nil && len(flag.name) > 0 {
			b.WriteRune(' ')
			b.WriteString(flag.name)

			if len(flag.value) > 0 {
				b.WriteRune('=')
				b.WriteString(flag.value)
			}
		}
	}

	addFlag(e.brightness)
	addFlag(e.speed)
	addFlag(e.colorNum)
	addFlag(e.customColorNum)
	addFlag(e.direction)

	addFlag(e.red)
	addFlag(e.green)
	addFlag(e.blue)
	addFlag(e.rgb)
	addFlag(e.colorName)

	addBool(e.reactive)
	addBool(e.save)
	addBool(e.reset)
	addFlag(e.pollInterval)
	addFlag(e.pollTimeout)
	addFlag(e.deviceBus)
	addFlag(e.deviceAddress)

	addFlag(e.configFile)

	b.WriteRune('`')

	return b.String()
}

// execsT type represents collection of command invocations
type execsT []*execT

// newExecs creates collection od command invocations with one initial invocation without any flags.
func newExecs() execsT {
	return execsT([]*execT{&execT{}})
}

// addFlagWithValue adds test cases with given flags and values.
// update function is used to set invocation flag and is usually one of execT setters.
// useDefaultVal specifies whether an invocation without the given flag will be added.
func (exs execsT) addFlagWithValue(useDefaultVal bool, flags []string, update func(src *execT, flag *flag),
	values ...string) execsT {

	f, v, le, lv, lf, l := 0, 0, len(exs), len(values)-1, len(flags)-1, len(values)

	if l < lf+1 {
		l = lf + 1 // use max length
	}

	if useDefaultVal {
		l++ // add another test case
	}

	if l < le {
		l = le // use max length
	}

	for i := range l {
		if i == le {
			// need another test case
			ex := &execT{}
			*ex = *exs[i-1] // copy prev test case

			exs = append(exs, ex)
			le++
		}

		if useDefaultVal {
			// add test case wit empty flag
			useDefaultVal = false
			update(exs[i], &flag{})
			continue
		}

		// add execution with new or last flag and value
		update(exs[i], &flag{name: flags[f], value: values[v]})

		if f < lf {
			f++ // take next flag name
		}
		if v < lv {
			v++ // take next flag value
		}
	}

	return exs
}

func (exs execsT) Brightness(values ...string) execsT {
	return exs.addFlagWithValue(true, brightnessFlagAll, (*execT).SetBrightness, values...)
}

func (exs execsT) RequiredBrightness(values ...string) execsT {
	return exs.addFlagWithValue(false, brightnessFlagAll, (*execT).SetBrightness, values...)
}

func (exs execsT) Speed(values ...string) execsT {
	return exs.addFlagWithValue(true, speedFlagAll, (*execT).SetSpeed, values...)
}

func (exs execsT) RequiredSpeed(values ...string) execsT {
	return exs.addFlagWithValue(false, speedFlagAll, (*execT).SetSpeed, values...)
}

func (exs execsT) ColorNum(values ...string) execsT {
	return exs.addFlagWithValue(true, colorNumFlagAll, (*execT).SetColorNum,
		values...)
}

func (exs execsT) RequiredColorNum(values ...string) execsT {
	return exs.addFlagWithValue(false, colorNumFlagAll, (*execT).SetColorNum,
		values...)
}

func (exs execsT) CustomColorNum(values ...string) execsT {
	return exs.addFlagWithValue(false, customColorNumFlagAll, (*execT).SetCustomColorNum,
		values...)
}

func (exs execsT) Direction(values ...string) execsT {
	return exs.addFlagWithValue(true, directionFlagAll, (*execT).SetDirection,
		values...)
}

func (exs execsT) RequiredDirection(values ...string) execsT {
	return exs.addFlagWithValue(false, directionFlagAll, (*execT).SetDirection,
		values...)
}

func (exs execsT) Red(values ...string) execsT {
	return exs.addFlagWithValue(true, redFlagAll, (*execT).SetRed, values...)
}

func (exs execsT) RequiredRed(values ...string) execsT {
	return exs.addFlagWithValue(false, redFlagAll, (*execT).SetRed, values...)
}

func (exs execsT) Green(values ...string) execsT {
	return exs.addFlagWithValue(true, greenFlagAll, (*execT).SetGreen, values...)
}

func (exs execsT) RequiredGreen(values ...string) execsT {
	return exs.addFlagWithValue(false, greenFlagAll, (*execT).SetGreen, values...)
}

func (exs execsT) Blue(values ...string) execsT {
	return exs.addFlagWithValue(true, blueFlagAll, (*execT).SetBlue, values...)
}

func (exs execsT) RequiredBlue(values ...string) execsT {
	return exs.addFlagWithValue(false, blueFlagAll, (*execT).SetBlue, values...)
}

func (exs execsT) RequiredRGB(values ...string) execsT {
	return exs.addFlagWithValue(false, rgbFlagAll, (*execT).SetRGB, values...)
}

func (exs execsT) ColorName(values ...string) execsT {
	return exs.addFlagWithValue(true, colorNameFlagAll, (*execT).SetColorName, values...)
}

func (exs execsT) RequiredColorName(values ...string) execsT {
	return exs.addFlagWithValue(false, colorNameFlagAll, (*execT).SetColorName, values...)
}

func (exs execsT) Reactive(values ...string) execsT {
	return exs.addFlagWithValue(true, reactiveFlagAll, (*execT).SetReactive, values...)
}

func (exs execsT) RequiredReactive(values ...string) execsT {
	return exs.addFlagWithValue(false, reactiveFlagAll, (*execT).SetReactive, values...)
}

func (exs execsT) Save(values ...string) execsT {
	return exs.addFlagWithValue(true, saveFlagAll, (*execT).SetSave, values...)
}

func (exs execsT) RequiredSave(values ...string) execsT {
	return exs.addFlagWithValue(false, saveFlagAll, (*execT).SetSave, values...)
}

func (exs execsT) Reset(values ...string) execsT {
	return exs.addFlagWithValue(true, resetFlagAll, (*execT).SetReset, values...)
}

func (exs execsT) RequiredReset(values ...string) execsT {
	return exs.addFlagWithValue(false, resetFlagAll, (*execT).SetReset, values...)
}

func (exs execsT) PollInterval(values ...string) execsT {
	return exs.addFlagWithValue(true, pollIntervalFlagAll, (*execT).SetPollInterval, values...)
}

func (exs execsT) RequiredPollInterval(values ...string) execsT {
	return exs.addFlagWithValue(false, pollIntervalFlagAll, (*execT).SetPollInterval, values...)
}

func (exs execsT) PollTimeout(values ...string) execsT {
	return exs.addFlagWithValue(true, pollTimeoutFlagAll, (*execT).SetPollTimeout, values...)
}

func (exs execsT) RequiredPollTimeout(values ...string) execsT {
	return exs.addFlagWithValue(false, pollTimeoutFlagAll, (*execT).SetPollTimeout, values...)
}

func (exs execsT) DeviceBus(values ...string) execsT {
	return exs.addFlagWithValue(true, deviceBusFlagAll, (*execT).SetDeviceBus, values...)
}

func (exs execsT) RequiredDeviceBus(values ...string) execsT {
	return exs.addFlagWithValue(false, deviceBusFlagAll, (*execT).SetDeviceBus, values...)
}

func (exs execsT) DeviceAddress(values ...string) execsT {
	return exs.addFlagWithValue(true, deviceAddressFlagAll, (*execT).SetDeviceAddress, values...)
}

func (exs execsT) RequiredDeviceAddress(values ...string) execsT {
	return exs.addFlagWithValue(false, deviceAddressFlagAll, (*execT).SetDeviceAddress, values...)
}

func (exs execsT) ConfigFile(values ...string) execsT {
	return exs.addFlagWithValue(true, configFileFlagAll, (*execT).SetConfigFile, values...)
}

func (exs execsT) RequiredConfigFile(values ...string) execsT {
	return exs.addFlagWithValue(false, configFileFlagAll, (*execT).SetConfigFile, values...)
}

// entries returns test cases as ginkgo test entries
func (exs execsT) entries() []TableEntry {

	entries := []TableEntry{}
	for _, in := range exs {
		entries = append(entries, Entry(in.String(), in))
	}

	return entries
}
