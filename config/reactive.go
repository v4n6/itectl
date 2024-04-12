package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const reactiveDefault = false

const reactiveProp = "reactive"

// Reactive returns either specified, configured or default value of the reactive flag.
func Reactive() bool {
	return viper.GetBool(reactiveProp)
}

// AddReactiveFlag adds save flag to the provided cmd and binds it to the corresponding viper config property.
func AddReactiveFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool(reactiveProp, reactiveDefault, "Make the keyboard backlight effect reactive.")

	bindAndValidate(cmd, reactiveProp, reactiveProp, nil)
}
