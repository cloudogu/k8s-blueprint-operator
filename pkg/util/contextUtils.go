package util

import (
	"context"
	"time"
)

// RetryUntilSuccessOrCancellation executes the given function in the given interval starting immediately
// until the given function returns that no retry should be done.
// The given function needs to return
//   - a generic value or nil,
//   - an error and
//   - a boolean which determines if a retry should happen.
//
// This function returns the result and the error of the given function or nil and the `ctx.Err` if the context is cancelled.
func RetryUntilSuccessOrCancellation[R interface{}](
	ctx context.Context,
	retryInterval time.Duration,
	function func(context.Context) (*R, error, bool),
) (*R, error) {
	ticker := time.NewTicker(retryInterval)
	defer ticker.Stop()

	for {
		result, userError, shouldRetry := function(ctx)
		if !shouldRetry {
			return result, userError
		}

		select {
		case <-ticker.C:
			continue // use ticker to start next try
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
