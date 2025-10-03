package ecosystem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebugMode_IsActive(t *testing.T) {
	t.Run("is active on status set", func(t *testing.T) {
		debugMode := DebugMode{Phase: DebugModeStatusSet}
		assert.True(t, debugMode.IsActive())
	})

	t.Run("is active on wait for rollback", func(t *testing.T) {
		debugMode := DebugMode{Phase: DebugModeStatusWaitForRollback}
		assert.True(t, debugMode.IsActive())
	})

	t.Run("is not active on anything else", func(t *testing.T) {
		debugMode := DebugMode{Phase: "not active"}
		assert.False(t, debugMode.IsActive())
	})
}
