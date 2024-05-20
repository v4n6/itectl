package params

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// SaveDefault is "save" property default value
const SaveDefault = false

// SaveProp is name of "save" flag and "save" configuration property.
const SaveProp = "save"

// AddSave adds "save" flag to the provided cmd.
// It also adds hook to bind it to the corresponding viper config property.
func AddSave(cmd *cobra.Command, v *viper.Viper) {
	cmd.PersistentFlags().Bool(SaveProp, SaveDefault,
		fmt.Sprintf("Instruct the controller to save its state. %s", configurationWarning))
	bindAndValidate(cmd, v, SaveProp, SaveProp, nil)
}

// Save returns value of save property
func Save(v *viper.Viper) bool {

	return v.GetBool(SaveProp)
}
