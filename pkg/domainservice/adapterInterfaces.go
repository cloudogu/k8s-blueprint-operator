package domainservice

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"

	"github.com/cloudogu/cesapp-lib/core"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
)

type DoguInstallationRepository interface {
	// GetByName returns the ecosystem.DoguInstallation or
	// a NotFoundError if the dogu is not installed or
	// an InternalError if there is any other error.
	GetByName(ctx context.Context, doguName string) (*ecosystem.DoguInstallation, error)
	// GetAll returns the installation info of all installed dogus or
	// a NotFoundError if any dogu is not installed or
	// an InternalError if there is any other error.
	GetAll(ctx context.Context) (map[string]*ecosystem.DoguInstallation, error)
	//Create(ctx context.Context, dogu ecosystem.DoguInstallation) error
	//Update(ctx context.Context, dogu ecosystem.DoguInstallation) error
	//Delete(ctx context.Context, dogu ecosystem.DoguInstallation) error
}

type BlueprintSpecRepository interface {
	// GetById returns a BlueprintSpec identified by its ID or
	// a NotFoundError if the BlueprintSpec was not found or
	// a domain.InvalidBlueprintError together with a BlueprintSpec without blueprint and mask if the BlueprintSpec could not be parsed or
	// an InternalError if there is any other error.
	GetById(ctx context.Context, blueprintId string) (*domain.BlueprintSpec, error)
	// Update updates a given BlueprintSpec.
	// returns a ConflictError if there were changes on the BlueprintSpec in the meantime or
	// returns an InternalError if there is any other error
	Update(ctx context.Context, blueprintSpec *domain.BlueprintSpec) error
}

type RemoteDoguRegistry interface {
	//GetDogu returns the dogu specification for the given dogu and version or
	//an NotFoundError indicating that there was no dogu spec found or
	//an InternalError indicating that the caller has no fault.
	GetDogu(qualifiedDoguName string, version string) (*core.Dogu, error)

	//GetDogus returns the all requested dogu specifications or
	//an NotFoundError indicating that any dogu spec was not found or
	//an InternalError indicating that the caller has no fault.
	GetDogus(dogusToLoad []DoguToLoad) (map[string]*core.Dogu, error)
}

type DoguToLoad struct {
	QualifiedDoguName string
	Version           string
}

type MaintenanceMode interface {
	GetLock() (MaintenanceLock, error)
	Activate(content MaintenancePageModel) error
	Deactivate() error
}

type MaintenanceLock interface {
	IsActive() bool
	IsOurs() bool
}

// NotFoundError is a common error indicating that sth. was requested but not found on the other side.
type NotFoundError struct {
	WrappedError error
	Message      string
}

// Error marks the struct as an error.
func (e *NotFoundError) Error() string {
	if e.WrappedError != nil {
		return fmt.Errorf("%s: %w", e.Message, e.WrappedError).Error()
	}
	return e.Message
}

// Unwrap is used to make it work with errors.Is, errors.As.
func (e *NotFoundError) Unwrap() error {
	return e.WrappedError
}

// InternalError is a common error indicating that there was an error at the called side independent of the specific call.
type InternalError struct {
	WrappedError error
	Message      string
}

// Error marks the struct as an error.
func (e *InternalError) Error() string {
	if e.WrappedError != nil {
		return fmt.Errorf("%s: %w", e.Message, e.WrappedError).Error()
	}
	return e.Message
}

// Unwrap is used to make it work with errors.Is, errors.As.
func (e *InternalError) Unwrap() error {
	return e.WrappedError
}

// ConflictError is a common error indicating that the aggregate was modified in the meantime.
type ConflictError struct {
	WrappedError error
	Message      string
}

// Error marks the struct as an error.
func (e *ConflictError) Error() string {
	if e.WrappedError != nil {
		return fmt.Errorf("%s: %w", e.Message, e.WrappedError).Error()
	}
	return e.Message
}

// Unwrap is used to make it work with errors.Is, errors.As.
func (e *ConflictError) Unwrap() error {
	return e.WrappedError
}
