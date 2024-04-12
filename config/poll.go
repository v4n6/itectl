package config

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	pollDefault = false

	pollIntervalDefault = 200 * time.Millisecond

	pollTimeoutDefault = time.Second
)

const (
	pollProp = "poll.always"

	pollIntervalProp = "poll.interval"

	pollTimeoutProp = "poll.timeout"
)

const (
	pollFlag = "poll"

	pollIntervalFlag = "poll-interval"

	pollTimeoutFlag = "poll-timeout"
)

func Poll() bool {
	return viper.GetBool(pollProp)
}

func PollInterval() time.Duration {
	return viper.GetDuration(pollIntervalProp)
}

func PollTimeout() time.Duration {
	return viper.GetDuration(pollTimeoutProp)
}

func AddPollFlags(cmd *cobra.Command) {

	cmd.PersistentFlags().Bool(pollFlag, pollDefault, "Instructs application to poll presence of the controller before executing command")
	cmd.PersistentFlags().Duration(pollIntervalFlag, pollIntervalDefault, "Time interval to wait between controller polls.")
	cmd.PersistentFlags().Duration(pollTimeoutFlag, pollTimeoutDefault, "Maximum time to wait for controller.")

	bindAndValidate(cmd, pollFlag, pollProp, func() error {

		err := viper.BindPFlag(pollProp, cmd.Flag(pollFlag))
		if err != nil {
			return err
		}

		err = viper.BindPFlag(pollIntervalProp, cmd.Flag(pollIntervalFlag))
		if err != nil {
			return err
		}

		err = viper.BindPFlag(pollTimeoutProp, cmd.Flag(pollTimeoutFlag))
		if err != nil {
			return err
		}

		if Poll() {
			i, t := PollInterval(), PollTimeout()
			if i >= t {
				return fmt.Errorf("%w %q; expected poll-interval '%s' to be less than poll-timeout '%s' ",
					InvalidOptionValueError, "poll-interval", i, t)
			}
		}

		return nil
	})
}
