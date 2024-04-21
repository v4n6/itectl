package params

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// poll property default value
	pollDefault = false
	// poll interval property default value
	pollIntervalDefault = 200 * time.Millisecond
	// poll timeout property default value
	pollTimeoutDefault = time.Second
)

const (
	// name poll configuration property
	pollProp = "poll.always"
	// name poll flag
	pollFlag = "poll"

	// name poll interval configuration property
	pollIntervalProp = "poll.interval"
	// name poll interval flag
	pollIntervalFlag = "poll-interval"

	// name poll timeout configuration property
	pollTimeoutProp = "poll.timeout"
	// name poll timeout flag
	pollTimeoutFlag = "poll-timeout"
)

// AddPoll adds polling related flags to the provided cmd.
// It also adds hook to bind them to the corresponding viper config properties
// and to validate poll properties. It ensures that if poll is true poll interval is less than timeout.
// It returns functions to retrieve current poll, pollInterval and pollTimeout values.
func AddPoll(cmd *cobra.Command, v *viper.Viper) (poll func() bool, pollInterval, pollTimeout func() time.Duration) {

	cmd.PersistentFlags().Duration(pollIntervalFlag, pollIntervalDefault, "Time interval to wait between controller polls.")
	bindAndValidate(cmd, v, pollIntervalFlag, pollIntervalProp, nil)
	cmd.PersistentFlags().Duration(pollTimeoutFlag, pollTimeoutDefault, "Maximum time to wait for controller.")
	bindAndValidate(cmd, v, pollTimeoutFlag, pollTimeoutProp, nil)
	cmd.PersistentFlags().Bool(pollFlag, pollDefault, "Instructs application to poll presence of the controller before executing command")
	bindAndValidate(cmd, v, pollFlag, pollProp, func() error {
		if v.GetBool(pollProp) {
			i, t := v.GetDuration(pollIntervalProp), v.GetDuration(pollTimeoutProp)
			if i >= t {
				return fmt.Errorf("%w %q; expected poll-interval '%s' to be less than poll-timeout '%s' ",
					InvalidOptionValueError, "poll-interval", i, t)
			}
		}

		return nil
	})

	return func() bool { return v.GetBool(pollProp) },
		func() time.Duration { return v.GetDuration(pollIntervalProp) },
		func() time.Duration { return v.GetDuration(pollTimeoutProp) }
}
