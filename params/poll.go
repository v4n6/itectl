package params

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// PollIntervalDefault is poll interval property default value
	PollIntervalDefault = 200 * time.Millisecond
	// PollTimeoutDefault is poll timeout property default value
	PollTimeoutDefault = 0
)

const (
	// pollIntervalProp is name poll interval configuration property
	pollIntervalProp = "poll.interval"
	// PollIntervalFlag is name poll interval flag
	PollIntervalFlag = "poll-interval"

	// pollTimeoutProp is name poll timeout configuration property
	pollTimeoutProp = "poll.timeout"
	// PollTimeoutFlag is name poll timeout flag
	PollTimeoutFlag = "poll-timeout"
)

// AddPoll adds polling related flags to the provided cmd.
// It also adds hook to bind them to the corresponding viper config properties
// and to validate poll properties. It ensures that poll interval is a positive duration,
// poll timeout is a nonnegative duration and that interval is less than the timeout, if timeout is not 0.
// AddPoll returns functions to retrieve current poll interval and poll timeout values.
func AddPoll(cmd *cobra.Command, v *viper.Viper) (pollInterval, pollTimeout func() time.Duration) {

	cmd.PersistentFlags().Duration(PollIntervalFlag, PollIntervalDefault,
		fmt.Sprintf("Time interval to wait between controller polls. The value is ignored if --%s is set to 0. %s",
			PollTimeoutFlag, configurationWarning))
	bindAndValidate(cmd, v, PollIntervalFlag, pollIntervalProp, func() error {
		if v.GetDuration(pollIntervalProp) > 0 {
			return nil
		}
		return fmt.Errorf("%w %q for (--%s): poll interval must be positive",
			InvalidOptionValueError, v.GetDuration(pollIntervalProp), PollIntervalFlag)
	})

	cmd.PersistentFlags().Duration(PollTimeoutFlag, PollTimeoutDefault,
		fmt.Sprintf("Maximum time to wait for controller to be available. Exit immediately, if it's set to 0 and no controller is available. %s",
			configurationWarning))
	bindAndValidate(cmd, v, PollTimeoutFlag, pollTimeoutProp, func() error {
		if v.GetDuration(pollTimeoutProp) > 0 {
			if i, t := v.GetDuration(pollIntervalProp), v.GetDuration(pollTimeoutProp); i >= t {
				return fmt.Errorf("%w %q (--%s) %q (--%s): expected poll interval to be less than poll timeout",
					InvalidOptionValueError, i, PollIntervalFlag, t, PollTimeoutFlag)
			}
		}

		return nil
	})

	return func() time.Duration { return v.GetDuration(pollIntervalProp) },
		func() time.Duration { return v.GetDuration(pollTimeoutProp) }
}
