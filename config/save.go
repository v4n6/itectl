package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const saveDefault = false

const saveProp = "save"

// Save returns either specified, configured or default value of the save flag.
func Save() bool {
	return viper.GetBool(saveProp)
}

// AddSaveFlag adds save flag to the provided cmd and binds it to the corresponding viper config property.
func AddSaveFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool(saveProp, saveDefault, "Instruct the controller to save the settings.")

	bindAndValidate(cmd, saveProp, saveProp, nil)
}
