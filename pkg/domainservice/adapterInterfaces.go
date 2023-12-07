package domainservice

import (
	"context"
	"fmt"

	"github.com/cloudogu/cesapp-lib/core"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

type DoguInstallationRepository interface {
	GetByName(doguName string) (ecosystem.DoguInstallation, error)
	GetAll() ([]ecosystem.DoguInstallation, error)
	Create(ecosystem.DoguInstallation) error
	Update(ecosystem.DoguInstallation) error
	Delete(ecosystem.DoguInstallation) error
}

type BlueprintSpecRepository interface {
	// GetById returns a BlueprintSpec identified by its ID or
	// a NotFoundError if the BlueprintSpec was not found or
	// an InternalError if there is any other error.
	GetById(ctx context.Context, blueprintId string) (domain.BlueprintSpec, error)
	// Update updates a given BlueprintSpec.
	// returns an InternalError if there is any error.
	//TODO: Maybe we also need an ConcurrentModificationError
	Update(ctx context.Context, blueprintSpec domain.BlueprintSpec) error
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
