package ecosystem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebugMode_IsActive(t *testing.T) {
	t.Run("is not active on completed", func(t *testing.T) {
		debugMode := DebugMode{Phase: DebugModeStatusComplete}
		assert.False(t, debugMode.IsActive())
	})

	t.Run("is not active on failed", func(t *testing.T) {
		debugMode := DebugMode{Phase: DebugModeStatusFailed}
		assert.False(t, debugMode.IsActive())
	})

	t.Run("is active on everything else", func(t *testing.T) {
		debugMode := DebugMode{Phase: "WaitForRollback"}
		assert.True(t, debugMode.IsActive())
	})
}
