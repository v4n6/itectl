package params

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/v4n6/ite8291r3tool/pkg/ite8291"
)

const (
	// redDefault is default value of red flag.
	redDefault = 0
	// greenDefault is default value of green flag.
	greenDefault = 0
	// blueDefault is default value of blue flag.
	blueDefault = 0

	// singleColorDefault is default value of single color configuration property.
	singleColorDefault = "#FFFFFF"
)

const (
	// colorRedFlag is name of red color flag.
	colorRedFlag = "red"
	// colorGreenFlag is name of green color flag.
	colorGreenFlag = "green"
	// colorBlueFlag is name of blue color flag.
	colorBlueFlag = "blue"

	// colorNameFlag is name of color name flag.
	colorNameFlag = "color-name"
	// colorRGBFlag is name of RGB flag.
	colorRGBFlag = "rgb"

	// singleColorProp is name of single color configuration property.
	singleColorProp = "singleModeColor"

	// namedColorsProp is name of named colors configuration property.
	namedColorsProp = "namedColors"
)

// colorNameToColor converts given color name to the corresponding instance of ite8291.Color.
// It returns InvalidOptionValueError if color name was not configured or
// color's value is not a valid color.
func colorNameToColor(name string, v *viper.Viper) (color *ite8291.Color, err error) {

	v_ := v.Sub(namedColorsProp)
	if v_ != nil {
		if val := v_.GetString(name); val != "" {
			color, err = ite8291.ParseColor(val)
			if err != nil {
				return nil, fmt.Errorf("%w: color name %q: %w", InvalidOptionValueError, name, err)
			}

			return color, nil
		}
	}

	return nil, fmt.Errorf("%w: unknown color name %q", InvalidOptionValueError, name)
}

// addColorFlags adds color related flags to the provided cmd.
// It also adds hook to validate their values.
// required parameter specifies whether color must be specified explicitly.
// addColorFlags returns functions to retrieve current red, grenn, blue flags and target color values.
func addColorFlags(cmd *cobra.Command, v *viper.Viper, required bool) (red, green, blue *byte, color **ite8291.Color) {

	var r, g, b byte
	var name, rgb string
	var col *ite8291.Color

	cmd.PersistentFlags().Uint8Var(&r, colorRedFlag, redDefault, "Red atom of RGB color.")
	cmd.PersistentFlags().Uint8Var(&g, colorGreenFlag, greenDefault, "Green atom of RGB color.")
	cmd.PersistentFlags().Uint8Var(&b, colorBlueFlag, blueDefault, "Blue atom of RGB color.")

	cmd.PersistentFlags().StringVar(&name, colorNameFlag, "", "Name of the color")
	cmd.PersistentFlags().StringVar(&rgb, colorRGBFlag, "", "Color as RGB value")

	cmd.MarkFlagsMutuallyExclusive(colorRGBFlag, colorNameFlag, colorRedFlag)
	cmd.MarkFlagsMutuallyExclusive(colorRGBFlag, colorNameFlag, colorGreenFlag)
	cmd.MarkFlagsMutuallyExclusive(colorRGBFlag, colorNameFlag, colorBlueFlag)
	if required {
		cmd.MarkFlagsOneRequired(colorNameFlag, colorRGBFlag,
			colorRedFlag, colorGreenFlag, colorBlueFlag)
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
				return fmt.Errorf("%w: color RGB %q: %w",
					InvalidOptionValueError, rgb, err)
			}
		}

		return nil
	})

	return &r, &g, &b, &col
}

// AddColor adds color related flags to the provided cmd.
// It also adds hook to validate their values. It ensures that the color is specified explicitly.
// AddColor returns function to retrieve current color value.
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

// AddSingleColor adds color related flags to the provided cmd.
// It also adds hook to validate their values. If color is not specified explicitly,
// it will be read from configuration. If color is neither specified nor configured
// the default value will be used.
// AddSingleColor returns function to retrieve current color value.
func AddSingleColor(cmd *cobra.Command, v *viper.Viper) (color func() *ite8291.Color) {

	r, g, b, col := addColorFlags(cmd, v, false)

	addValidationHook(cmd, func() (err error) {
		if *col == nil {
			if cmd.Flag(colorRedFlag).Changed || cmd.Flag(colorGreenFlag).Changed || cmd.Flag(colorBlueFlag).Changed {
				*col = ite8291.NewColor(*r, *g, *b)
				return nil
			}

			c := v.GetString(singleColorProp)
			if len(c) > 0 {
				if *col, err = colorNameToColor(c, v); err == nil {
					return nil
				}

				if *col, err = ite8291.ParseColor(c); err != nil {
					return fmt.Errorf("%w: single color %q: %w",
						InvalidOptionValueError, c, err)
				}
				return nil
			}

			*col, _ = ite8291.ParseColor(singleColorDefault)
		}

		return nil
	})

	return func() *ite8291.Color { return *col }
}
