package params

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// saveDefault is "save" property default value
const saveDefault = false

// saveProp is name of "save" flag and "save" configuration property.
const saveProp = "save"

// AddSave adds "save" flag to the provided cmd.
// It also adds hook to bind it to the corresponding viper config property.
// It returns functions to retrieve current "save" property value.
func AddSave(cmd *cobra.Command, v *viper.Viper) (save func() bool) {
	cmd.PersistentFlags().Bool(saveProp, saveDefault, "Instruct the controller to save the settings.")
	bindAndValidate(cmd, v, saveProp, saveProp, nil)

	return func() bool { return v.GetBool(saveProp) }
}
