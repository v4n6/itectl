package params

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// reactiveDefault is reactive property default value
const reactiveDefault = false

// reactiveProp is name of reactive flag and reactive configuration property.
const reactiveProp = "reactive"

// AddReactive adds reactive flag to the provided cmd.
// It also adds hook to bind it to the corresponding viper config property.
// It returns functions to retrieve current reactive property value.
func AddReactive(cmd *cobra.Command, v *viper.Viper) (reactive func() bool) {
	cmd.PersistentFlags().Bool(reactiveProp, reactiveDefault, "Make the keyboard backlight mode reactive.")
	bindAndValidate(cmd, v, reactiveProp, reactiveProp, nil)

	return func() bool { return v.GetBool(reactiveProp) }
}
