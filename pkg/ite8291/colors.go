package ite8291

import (
	"errors"
	"fmt"
)

// ErrInvalidColorFormat error signaling invalid format of color text representation.
var ErrInvalidColorFormat = errors.New("invalid color format")

// SupportedColorStringFormats provides supported color text formats.
var SupportedColorStringFormats = []string{
	"0xHHHHHH",
	"#xHHHHHH",
	"#HHHHHH",
	"HHHHHH",
	"#HHH",
	"HHH",
}

// Color type provides Color RGB representation.
type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

// String - stringer implementation.
func (c *Color) String() string {
	return fmt.Sprintf("#%02X%02X%02X", c.Red, c.Green, c.Blue)
}

// NewColor creates Color based on r,g,b params.
func NewColor(r, g, b uint8) *Color {
	return &Color{Red: r, Green: g, Blue: b}
}

// FromRGB creates Color based on given rgb uint.
func FromRGB(rgb uint32) *Color {
	return &Color{Red: byte(rgb >> 16 & 0xFF), Green: byte(rgb >> 8 & 0xFF), Blue: byte(rgb & 0xFF)}
}

// ParseColor parses s to Color.
// It accepts following formats 0xHHHHHH,#xHHHHHH,#HHHHHH,HHHHHH,#HHH,HHH.
//
//nolint:funlen,cyclop
func ParseColor(s string) (*Color, error) {

	var i, l = 0, len(s)

	newErr := func() error {
		return fmt.Errorf("%w: expected one of 0xHHHHHH,#xHHHHHH,#HHHHHH,HHHHHH,#HHH,HHH was %q", ErrInvalidColorFormat, s)
	}

	switch l {
	case 8:
		if (s[0] != '0' && s[0] != '#') || (s[1] != 'x' && s[1] != 'X') {
			return nil, newErr()
		} else {
			i += 2
		}

	case 7:
		if s[0] != '#' {
			return nil, newErr()
		} else {
			i++
		}

	case 6:

	case 4:
		if s[0] != '#' {
			return nil, newErr()
		} else {
			i++
		}

	case 3:

	default:
		return nil, newErr()
	}

	if l-i == 6 {
		r, ok := hexToUint8(s, i)
		if !ok {
			return nil, newErr()
		}

		g, ok := hexToUint8(s, i+2)
		if !ok {
			return nil, newErr()
		}

		b, ok := hexToUint8(s, i+4)
		if !ok {
			return nil, newErr()
		}

		return NewColor(r, g, b), nil
	}

	r, ok := hexToUint4(s, i)
	if !ok {
		return nil, newErr()
	}

	g, ok := hexToUint4(s, i+1)
	if !ok {
		return nil, newErr()
	}

	b, ok := hexToUint4(s, i+2)
	if !ok {
		return nil, newErr()
	}

	return NewColor(r*17, g*17, b*17), nil

}

func hexToUint4(s string, idx int) (uint8, bool) {

	c := s[idx]

	switch {
	case c >= '0' && c <= '9':
		return c - '0', true
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10, true
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10, true
	}

	return 0, false
}

func hexToUint8(s string, idx int) (b uint8, ok bool) {

	for k := range 2 {
		c, ok := hexToUint4(s, idx+k)
		if !ok {
			return 0, false
		}

		b = (b << k * 8) | (c & 0xF)
	}

	return b, true
}
