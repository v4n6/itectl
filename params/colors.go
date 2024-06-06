package params

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/itectl/pkg/ite8291"
)

// red, green, blue properties default values.
const (
	// RedDefault - default value of red flag.
	RedDefault = 0
	// GreenDefault - default value of green flag.
	GreenDefault = 0
	// BlueDefault - default value of blue flag.
	BlueDefault = 0
)

// SingleColorDefault - default value of single color property.
const SingleColorDefault = "#FFFFFF"

// color properties and flags names.
const (
	// ColorRedFlag - name of the red color flag.
	ColorRedFlag = "red"
	// ColorGreenFlag - name of the green color flag.
	ColorGreenFlag = "green"
	// ColorBlueFlag - name of the blue color flag.
	ColorBlueFlag = "blue"

	// ColorNameFlag - name of the color name flag.
	ColorNameFlag = "color-name"
	// ColorRGBFlag - name of the RGB flag.
	ColorRGBFlag = "rgb"

	// SingleColorProp - name of the single color configuration property.
	SingleColorProp = "singleModeColor"

	// NamedColorsProp - name of the named colors configuration property.
	NamedColorsProp = "namedColors"
)

// colorNameToColor converts given color name to the corresponding
// instance of ite8291.Color. It returns ErrInvalidOptVal if color
// name was not configured or its value is not a valid color.
func colorNameToColor(name string, v *viper.Viper) (color *ite8291.Color, err error) {

	if val := v.GetString(fmt.Sprintf("%s.%s", NamedColorsProp, name)); len(val) > 0 {

		color, err = ite8291.ParseColor(val)
		if err != nil {
			return nil, fmt.Errorf("%w %q for %q: %w", ErrInvalidOptVal, name,
				"--"+ColorNameFlag, err)
		}

		return color, nil
	}

	return nil, fmt.Errorf("%w %q for %q is an unknown color name", ErrInvalidOptVal, name,
		"--"+ColorNameFlag)
}

// addColorFlags adds color related flags to the provided cmd. It also
// adds hook to validate their values. The 'required' parameter
// specifies whether color must be specified explicitly.
// addColorFlags returns pointers to red, green, blue and target color
// values.
func addColorFlags(cmd *cobra.Command, v *viper.Viper, required bool) (red, green, blue *byte, color **ite8291.Color) {

	var r, g, b byte
	var name, rgb string
	var col *ite8291.Color

	cmd.PersistentFlags().Uint8Var(&r, ColorRedFlag, RedDefault, "Red part of RGB color.")
	cmd.PersistentFlags().Uint8Var(&g, ColorGreenFlag, GreenDefault, "Green part of RGB color.")
	cmd.PersistentFlags().Uint8Var(&b, ColorBlueFlag, BlueDefault, "Blue part of RGB color.")

	cmd.PersistentFlags().StringVar(&name, ColorNameFlag, "",
		fmt.Sprintf("Name of the color to use. One of color names configured via %q property in configuration file(s).",
			NamedColorsProp))
	cmd.PersistentFlags().StringVar(&rgb, ColorRGBFlag, "",
		fmt.Sprintf("Color as RGB value in a one of the following formats %q",
			ite8291.SupportedColorStringFormats))

	cmd.MarkFlagsMutuallyExclusive(ColorRGBFlag, ColorNameFlag, ColorRedFlag)
	cmd.MarkFlagsMutuallyExclusive(ColorRGBFlag, ColorNameFlag, ColorGreenFlag)
	cmd.MarkFlagsMutuallyExclusive(ColorRGBFlag, ColorNameFlag, ColorBlueFlag)
	if required {
		cmd.MarkFlagsOneRequired(ColorNameFlag, ColorRGBFlag,
			ColorRedFlag, ColorGreenFlag, ColorBlueFlag)
	}

	addValidationHook(cmd, func() (err error) {

		if len(name) > 0 {

			if col, err = colorNameToColor(name, v); err != nil {
				return err
			}

			return nil
		}

		if len(rgb) > 0 {

			if col, err = ite8291.ParseColor(rgb); err != nil {
				return fmt.Errorf("%w %q for %q: %w", ErrInvalidOptVal, rgb,
					"--"+ColorRGBFlag, err)
			}
		}

		return nil
	})

	return &r, &g, &b, &col
}

// AddColor adds color related flags to the provided cmd. It also adds
// hook to validate their values. It ensures that the color is
// specified explicitly via flags. AddColor returns function to
// retrieve current color value.
func AddColor(cmd *cobra.Command, v *viper.Viper) (color func() *ite8291.Color) {

	r, g, b, col := addColorFlags(cmd, v, true)

	addValidationHook(cmd, func() (err error) {
		if *col == nil {
			*col = ite8291.NewColor(*r, *g, *b)
		}
		return nil
	})

	return func() *ite8291.Color { return *col }
}

// AddSingleModeColor adds color related flags to the provided
// 'single-color-mode' cmd.  It also adds hook to validate their
// values. It uses either configured or default value if color is not
// specified explicitly.  AddSingleModeColor returns function to
// retrieve current value of color for single-color mode.
func AddSingleModeColor(cmd *cobra.Command, v *viper.Viper) (color func() *ite8291.Color) {

	r, g, b, col := addColorFlags(cmd, v, false)

	addValidationHook(cmd, func() (err error) {
		if *col == nil {
			// no color via color name or rgb
			if cmd.Flag(ColorRedFlag).Changed || cmd.Flag(ColorGreenFlag).Changed || cmd.Flag(ColorBlueFlag).Changed {
				// color  via red, green, blue flags
				*col = ite8291.NewColor(*r, *g, *b)
				return nil
			}
			// try configured color
			c := v.GetString(SingleColorProp)
			if len(c) == 0 {
				c = SingleColorDefault // no configured -> use default
			}

			// try as color name
			if *col, err = colorNameToColor(c, v); err == nil {
				return nil
			}

			// it must be rgb
			if *col, err = ite8291.ParseColor(c); err != nil {
				return fmt.Errorf("%w %q for configured single mode color: %w",
					ErrInvalidOptVal, c, err)
			}
		}

		return nil
	})

	return func() *ite8291.Color { return *col }
}
