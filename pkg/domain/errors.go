package domain

import (
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
)

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

type UnhealthyEcosystemError struct {
	WrappedError error
	Message      string
	healthResult ecosystem.HealthResult
}

func (e *UnhealthyEcosystemError) Error() string {
	unhealthyComponentsText := e.healthResult.ComponentHealth.String()
	unhealthyDogusText := e.healthResult.DoguHealth.String()
	combinedMessage := fmt.Sprintf("%s - %s - %s", e.Message, unhealthyDogusText, unhealthyComponentsText)
	if e.WrappedError != nil {
		return fmt.Errorf("%s: %w", combinedMessage, e.WrappedError).Error()
	}
	return combinedMessage
}

// Unwrap is used to make it work with errors.Is, errors.As.
func (e *UnhealthyEcosystemError) Unwrap() error {
	return e.WrappedError
}

func NewUnhealthyEcosystemError(
	wrappedError error,
	message string,
	healthResult ecosystem.HealthResult,
) *UnhealthyEcosystemError {
	return &UnhealthyEcosystemError{WrappedError: wrappedError, Message: message, healthResult: healthResult}
}

// DogusNotUpToDateError indicates that there are dogus that are not yet up to date.
type DogusNotUpToDateError struct {
	Message string
}

func (e *DogusNotUpToDateError) Error() string {
	return e.Message
}

// MultipleBlueprintsError indicates that there are multiple blueprint-resources in this namespace, which the controller cannot handle.
type MultipleBlueprintsError struct {
	Message string
}

// Error marks the struct as an error.
func (e *MultipleBlueprintsError) Error() string {
	return e.Message
}

type AwaitSelfUpgradeError struct {
	Message string
}

func (e *AwaitSelfUpgradeError) Error() string {
	return e.Message
}

type StateDiffNotEmptyError struct {
	Message string
}

func (e *StateDiffNotEmptyError) Error() string {
	return e.Message
}
