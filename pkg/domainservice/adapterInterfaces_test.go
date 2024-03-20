package domainservice

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInternalError(t *testing.T) {
	t.Run("should return error", func(t *testing.T) {
		err := NewInternalError(nil, "")
		require.Error(t, err)
	})
	t.Run("should return error without wrapped error or interpreted msgargs", func(t *testing.T) {
		err := NewInternalError(nil, "msg")
		assert.Equal(t, "msg", err.Error())
	})
	t.Run("should return error with wrapped error but no interpreted msgargs", func(t *testing.T) {
		err := NewInternalError(assert.AnError, "msg")
		assert.Equal(t, "msg: assert.AnError general error for testing", err.Error())
	})
	t.Run("should return error with wrapped error and interpret msgargs", func(t *testing.T) {
		err := NewInternalError(assert.AnError, "here is %s %s", "a", "string")
		assert.Equal(t, "here is a string: assert.AnError general error for testing", err.Error())
	})
	t.Run("should return error without wrapped error and interpret msgargs", func(t *testing.T) {
		err := NewInternalError(nil, "here is %s %s", "a", "string")
		assert.Equal(t, "here is a string", err.Error())
	})
}

func TestInternalError_Error(t *testing.T) {
	t.Run("without wrapped error", func(t *testing.T) {
		actual := InternalError{WrappedError: nil, Message: "test"}
		assert.Equal(t, "test", actual.Error())
	})
	t.Run("with wrapped error", func(t *testing.T) {
		actual := InternalError{WrappedError: errors.New("test2"), Message: "test"}
		assert.Equal(t, "test: test2", actual.Error())
	})
}

func TestInternalError_Unwrap(t *testing.T) {
	t.Run("without wrapped error", func(t *testing.T) {
		actual := InternalError{WrappedError: nil, Message: "test"}
		assert.Nil(t, actual.Unwrap())
	})
	t.Run("with wrapped error", func(t *testing.T) {
		actual := InternalError{WrappedError: errors.New("test2"), Message: "test"}
		assert.Error(t, actual.Unwrap())
		assert.ErrorContains(t, actual.Unwrap(), "test2")
	})
}

func TestNotFoundError_Error(t *testing.T) {
	t.Run("without wrapped error", func(t *testing.T) {
		actual := NotFoundError{WrappedError: nil, Message: "test"}
		assert.Equal(t, "test", actual.Error())
	})
	t.Run("with wrapped error", func(t *testing.T) {
		actual := NotFoundError{WrappedError: errors.New("test2"), Message: "test"}
		assert.Equal(t, "test: test2", actual.Error())
	})
}

func TestNotFoundError_Unwrap(t *testing.T) {
	t.Run("without wrapped error", func(t *testing.T) {
		actual := NotFoundError{WrappedError: nil, Message: "test"}
		assert.Nil(t, actual.Unwrap())
	})
	t.Run("with wrapped error", func(t *testing.T) {
		actual := NotFoundError{WrappedError: errors.New("test2"), Message: "test"}
		assert.Error(t, actual.Unwrap())
		assert.ErrorContains(t, actual.Unwrap(), "test2")
	})
}

func TestConflictError_Error(t *testing.T) {
	t.Run("without wrapped error", func(t *testing.T) {
		actual := ConflictError{WrappedError: nil, Message: "test"}
		assert.Equal(t, "test", actual.Error())
	})
	t.Run("with wrapped error", func(t *testing.T) {
		actual := ConflictError{WrappedError: errors.New("test2"), Message: "test"}
		assert.Equal(t, "test: test2", actual.Error())
	})
}

func TestConflictError_Unwrap(t *testing.T) {
	t.Run("without wrapped error", func(t *testing.T) {
		actual := ConflictError{WrappedError: nil, Message: "test"}
		assert.Nil(t, actual.Unwrap())
	})
	t.Run("with wrapped error", func(t *testing.T) {
		actual := ConflictError{WrappedError: errors.New("test2"), Message: "test"}
		assert.Error(t, actual.Unwrap())
		assert.ErrorContains(t, actual.Unwrap(), "test2")
	})
}

func TestIsNotFoundError(t *testing.T) {
	assert.True(t, IsNotFoundError(NewNotFoundError(assert.AnError, "test")))
	assert.True(t, IsNotFoundError(fmt.Errorf("test: %w", NewNotFoundError(assert.AnError, "test"))))
	assert.False(t, IsNotFoundError(assert.AnError))
}
