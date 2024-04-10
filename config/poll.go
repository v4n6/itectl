package config

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

const (
	pollIntervalDefault = 200 * time.Millisecond
	pollTimeoutDefault  = time.Second
)

type PollProp struct {
	Poll         bool
	PollInterval time.Duration
	PollTimeout  time.Duration
}

func (c *PollProp) PollVal() (bool, error) {
	return c.Poll, nil
}

func (c *PollProp) validate() error {

	if c.PollInterval >= c.PollTimeout {
		return fmt.Errorf("%w %q; expected poll-interval '%s' to be less than poll-timeout '%s' ",
			InvalidOptionValueError, "poll-interval", c.PollInterval, c.PollTimeout)
	}

	return nil
}

func (c *PollProp) PollIntervalVal() (time.Duration, error) {

	err := c.validate()
	if err != nil {
		return 0, err
	}

	return c.PollInterval, nil
}

func (c *PollProp) PollTimeoutVal() (time.Duration, error) {

	err := c.validate()
	if err != nil {
		return 0, err
	}

	return c.PollTimeout, nil
}

func AddPollFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&Config.Poll, "poll", "p", false, "Instructs application to poll presence of the controller before executing command")
}

func AddPollIntervalFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().DurationVar(&Config.PollInterval, "poll-interval", pollIntervalDefault, "Time interval to wait between controller polls.")
}

func AddPollTimeoutFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().DurationVar(&Config.PollTimeout, "poll-timeout", pollTimeoutDefault, "Maximum time to wait for controller.")
}
