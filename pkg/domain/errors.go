package domain

import "fmt"

// InvalidBlueprintError indicates that a given blueprintSpec is invalid for any reason.
type InvalidBlueprintError struct {
	WrappedError error
	Message      string
}

// Error marks the struct as an error.
func (e *InvalidBlueprintError) Error() string {
	if e.WrappedError != nil {
		return fmt.Errorf("%s: %w", e.Message, e.WrappedError).Error()
	}
	return e.Message
}

// Unwrap is used to make it work with errors.Is, errors.As.
func (e *InvalidBlueprintError) Unwrap() error {
	return e.WrappedError
}
