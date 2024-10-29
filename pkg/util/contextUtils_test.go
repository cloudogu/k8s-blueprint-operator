package util

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type testType struct {
	name string
}

func TestRetryUntilSuccessOrCancellation(t *testing.T) {
	expected := &testType{"test"}

	t.Run("ok", func(t *testing.T) {
		testCtx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		result, err := RetryUntilSuccessOrCancellation(
			testCtx,
			1*time.Second,
			func(ctx context.Context) (*testType, error, bool) {
				return expected, nil, false
			},
		)

		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("ok after retry", func(t *testing.T) {
		testCtx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		callCounter := 0

		result, err := RetryUntilSuccessOrCancellation(
			testCtx,
			5*time.Millisecond,
			func(ctx context.Context) (*testType, error, bool) {
				callCounter++
				if callCounter >= 2 { // success at second try
					return expected, nil, false
				} else {
					return nil, nil, true // retry at first try
				}
			},
		)

		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("timeout in given function", func(t *testing.T) {
		testCtx, cancel := context.WithTimeout(context.Background(), 0*time.Millisecond)
		defer cancel()

		result, err := RetryUntilSuccessOrCancellation(
			testCtx,
			1*time.Second,
			func(ctx context.Context) (*testType, error, bool) {
				select {
				case <-ctx.Done():
					return expected, ctx.Err(), false
				default:
					return nil, nil, true
				}
			},
		)

		require.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Equal(t, expected, result)
	})

	t.Run("timeout outside of given function", func(t *testing.T) {
		testCtx, cancel := context.WithTimeout(context.Background(), 0*time.Millisecond)
		defer cancel()

		result, err := RetryUntilSuccessOrCancellation(
			testCtx,
			1*time.Second,
			func(ctx context.Context) (*testType, error, bool) {
				//ignore timeout here, so that the helper function need to stop
				return nil, nil, true
			},
		)

		require.ErrorIs(t, err, context.DeadlineExceeded)
		var expectedNil *testType = nil
		assert.Equal(t, expectedNil, result)
	})

}
