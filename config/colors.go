package config

import (
	"github.com/spf13/cobra"
)

var (
	colorRed   byte
	colorGreen byte
	colorBlue  byte
)

const (
	redDefault   = 0
	greenDefault = 0
	blueDefault  = 0
)

const (
	colorRedFlag   = "red"
	colorGreenFlag = "green"
	colorBlueFlag  = "blue"
)

// ColorRed returns specified value of the red flag.
func ColorRed() uint8 {
	return colorRed
}

// ColorGreen returns specified value of the green flag.
func ColorGreen() uint8 {
	return colorGreen
}

// ColorBlue returns specified value of the blue flag.
func ColorBlue() uint8 {
	return colorBlue
}

// AddRedColorFlag adds red flag to the provided cmd.
func AddRedColorFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Uint8Var(&colorRed, colorRedFlag, redDefault, "Red atom of RGB color.")
}

// AddGreenColorFlag adds green flag to the provided cmd.
func AddGreenColorFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Uint8Var(&colorGreen, colorGreenFlag, greenDefault, "Green atom of RGB color.")
}

// AddBlueColorFlag adds blue flag to the provided cmd.
func AddBlueColorFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Uint8Var(&colorBlue, colorBlueFlag, blueDefault, "Blue atom of RGB color.")
}
