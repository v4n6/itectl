package config

import (
	"github.com/spf13/cobra"
)

const saveDefault = false

type SaveProp struct {
	Save bool
}

func (c *SaveProp) SaveVal() bool {
	return c.Save
}

func AddSaveFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&Config.Save, "save", saveDefault, "Instruct the controller to save the settings.")
}
