package params

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ReactiveDefault is reactive property default value
const ReactiveDefault = false

// ReactiveProp is name of reactive flag and reactive configuration property.
const ReactiveProp = "reactive"

// AddReactive adds reactive flag to the provided cmd.
// It also adds hook to bind it to the corresponding viper config property.
// AddReactive returns function to retrieve current reactive property value.
func AddReactive(cmd *cobra.Command, v *viper.Viper) (reactive func() bool) {
	cmd.PersistentFlags().Bool(ReactiveProp, ReactiveDefault,
		fmt.Sprintf("Make the keyboard backlight effect react to user input. %s",
			configurationWarning))
	bindAndValidate(cmd, v, ReactiveProp, ReactiveProp, nil)

	return func() bool { return v.GetBool(ReactiveProp) }
}
