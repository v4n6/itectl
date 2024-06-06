package params

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ReactiveDefault - default value of the reactive property.
const ReactiveDefault = false

// ReactiveProp - name of reactive flag and reactive configuration property.
const ReactiveProp = "reactive"

// AddReactive adds reactive flag to the provided cmd. It also adds
// hook to bind it to the corresponding viper configuration property.
func AddReactive(cmd *cobra.Command, v *viper.Viper) {
	cmd.PersistentFlags().Bool(ReactiveProp, ReactiveDefault,
		"Make the keyboard backlight effect react to user input. "+
			configurationWarning)
	bindAndValidate(cmd, v, ReactiveProp, ReactiveProp, nil)
}

// Reactive returns reactive property value.
func Reactive(v *viper.Viper) bool {

	return v.GetBool(ReactiveProp)
}
