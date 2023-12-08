package domain

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvalidBlueprintError_Error(t *testing.T) {
	t.Run("without wrapped error", func(t *testing.T) {
		actual := InvalidBlueprintError{WrappedError: nil, Message: "test"}
		assert.Equal(t, "test", actual.Error())
	})
	t.Run("with wrapped error", func(t *testing.T) {
		actual := InvalidBlueprintError{WrappedError: errors.New("test2"), Message: "test"}
		assert.Equal(t, "test: test2", actual.Error())
	})
}

func TestInvalidBlueprintError_Unwrap(t *testing.T) {
	t.Run("without wrapped error", func(t *testing.T) {
		actual := InvalidBlueprintError{WrappedError: nil, Message: "test"}
		assert.Nil(t, actual.Unwrap())
	})
	t.Run("with wrapped error", func(t *testing.T) {
		actual := InvalidBlueprintError{WrappedError: errors.New("test2"), Message: "test"}
		assert.Error(t, actual.Unwrap())
		assert.ErrorContains(t, actual.Unwrap(), "test2")
	})
}
