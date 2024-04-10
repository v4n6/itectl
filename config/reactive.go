package config

import (
	"github.com/spf13/cobra"
)

const reactiveDefault = false

type ReactiveProp struct {
	Reactive bool
}

func (c *ReactiveProp) ReactiveVal() bool {
	return c.Reactive
}

func AddReactiveFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&Config.Reactive, "reactive", reactiveDefault, "Make the keyboard backlight effect reactive.")
}
