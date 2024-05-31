package params

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// PollIntervalDefault is poll interval property default value.
	PollIntervalDefault = 200 * time.Millisecond
	// PollTimeoutDefault is poll timeout property default value.
	PollTimeoutDefault = 0
)

const (
	// pollIntervalProp is name poll interval configuration property.
	pollIntervalProp = "poll.interval"
	// PollIntervalFlag is name poll interval flag.
	PollIntervalFlag = "poll-interval"

	// pollTimeoutProp is name poll timeout configuration property.
	pollTimeoutProp = "poll.timeout"
	// PollTimeoutFlag is name poll timeout flag.
	PollTimeoutFlag = "poll-timeout"
)

// AddPoll adds polling related flags to the provided cmd.
// It also adds hook to bind them to the corresponding viper config properties.
func AddPoll(cmd *cobra.Command, v *viper.Viper) {

	cmd.PersistentFlags().Duration(PollIntervalFlag, PollIntervalDefault,
		fmt.Sprintf("Time interval to wait between controller polls. The value is ignored if --%s is set to 0. %s",
			PollTimeoutFlag, configurationWarning))
	bindAndValidate(cmd, v, PollIntervalFlag, pollIntervalProp, nil)

	//nolint:lll
	cmd.PersistentFlags().Duration(PollTimeoutFlag, PollTimeoutDefault,
		"Maximum time to wait for controller to be available. Exit immediately, if it's set to 0 and no controller is available. "+
			configurationWarning)
	bindAndValidate(cmd, v, PollTimeoutFlag, pollTimeoutProp, nil)
}

// Polls returns polling interval and timeout property values.
// It also ensures that poll interval is a positive duration
// less than the timeout, if timeout is not 0.
func Polls(v *viper.Viper) (pollInterval, pollTimeout time.Duration, err error) {

	interval, timeout := v.GetDuration(pollIntervalProp), v.GetDuration(pollTimeoutProp)
	if interval <= 0 {
		return 0, 0, fmt.Errorf("%w %q for (--%s): poll interval must be positive",
			ErrInvalidOptVal, v.GetDuration(pollIntervalProp), PollIntervalFlag)
	}

	if timeout > 0 && interval >= timeout {
		return 0, 0, fmt.Errorf("%w %q (--%s) %q (--%s): expected poll interval to be less than poll timeout",
			ErrInvalidOptVal, interval, PollIntervalFlag, timeout, PollTimeoutFlag)
	}

	return interval, timeout, nil
}
