package domainservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-registry-lib/config"

	"github.com/cloudogu/cesapp-lib/core"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
)

type DoguInstallationRepository interface {
	// GetByName returns the ecosystem.DoguInstallation or
	//  - a NotFoundError if the dogu is not installed or
	//  - an InternalError if there is any other error.
	GetByName(ctx context.Context, doguName common.SimpleDoguName) (*ecosystem.DoguInstallation, error)
	// GetAll returns the installation info of all installed dogus or
	//  - an InternalError if there is any other error.
	GetAll(ctx context.Context) (map[common.SimpleDoguName]*ecosystem.DoguInstallation, error)
	// Create saves a new ecosystem.DoguInstallation. This initiates a dogu installation. It returns
	//  - a ConflictError if there is already a DoguInstallation with this name or
	//  - an InternalError if there is any error while saving the DoguInstallation
	Create(ctx context.Context, dogu *ecosystem.DoguInstallation) error
	// Update updates an ecosystem.DoguInstallation in the ecosystem.
	//  - returns a ConflictError if there were changes on the DoguInstallation in the meantime or
	//  - returns a NotFoundError if the DoguInstallation was not found or
	//  - returns an InternalError if there is any other error
	Update(ctx context.Context, dogu *ecosystem.DoguInstallation) error
	// Delete removes the given ecosystem.DoguInstallation completely from the ecosystem.
	// We delete DoguInstallations with the object not just the name as this way we can detect concurrent updates.
	//  - returns a ConflictError if there were changes on the DoguInstallation in the meantime or
	//  - returns an InternalError if there is any other error
	Delete(ctx context.Context, doguName common.SimpleDoguName) error
}

type ComponentInstallationRepository interface {
	// GetByName loads an installed component from the ecosystem and returns
	//  - the ecosystem.ComponentInstallation or
	//  - a NotFoundError if the component is not installed or
	//  - an InternalError if there is any other error.
	GetByName(ctx context.Context, componentName common.SimpleComponentName) (*ecosystem.ComponentInstallation, error)
	// GetAll returns
	//  - the installation info of all installed components or
	//  - an InternalError if there is any other error.
	GetAll(ctx context.Context) (map[common.SimpleComponentName]*ecosystem.ComponentInstallation, error)
	// Delete deletes the component by name from the ecosystem.
	// returns an InternalError if there is an error.
	Delete(ctx context.Context, componentName common.SimpleComponentName) error
	// Create creates the ecosystem.ComponentInstallation in the ecosystem.
	// returns an InternalError if there is an error.
	Create(ctx context.Context, component *ecosystem.ComponentInstallation) error
	// Update updates the ecosystem.ComponentInstallation in the ecosystem.
	// returns an InternalError if anything went wrong.
	Update(ctx context.Context, component *ecosystem.ComponentInstallation) error
}

type RequiredComponentsProvider interface {
	GetRequiredComponents(ctx context.Context) ([]ecosystem.RequiredComponent, error)
}

type HealthWaitConfigProvider interface {
	GetWaitConfig(ctx context.Context) (ecosystem.WaitConfig, error)
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
	// GetDogu returns the dogu specification for the given dogu and version or
	// an NotFoundError indicating that there was no dogu spec found or
	// an InternalError indicating that the caller has no fault.
	GetDogu(doguName common.QualifiedDoguName, version string) (*core.Dogu, error)

	// GetDogus returns the all requested dogu specifications or
	// an NotFoundError indicating that any dogu spec was not found or
	// an InternalError indicating that the caller has no fault.
	GetDogus(dogusToLoad []DoguToLoad) (map[common.QualifiedDoguName]*core.Dogu, error)
}

type DoguToLoad struct {
	DoguName common.QualifiedDoguName
	Version  string
}

type MaintenanceMode interface {
	// Activate enables the maintenance mode with the given title and text.
	// May throw
	//  - a ConflictError if another party activated the maintenance mode or
	//  - a generic InternalError due to a connection error or an unknown error.
	Activate(ctx context.Context, title, text string) error
	// Deactivate disables the maintenance mode if it is active.
	// May throw
	//  - a ConflictError if another party initially activated the maintenance mode or
	//  - a generic InternalError due to a connection error or an unknown error.
	Deactivate(ctx context.Context) error
}

type DoguRestartRepository interface {
	// RestartAll restarts all provided Dogus
	RestartAll(context.Context, []common.SimpleDoguName) error
}

// GlobalConfigRepository is used to get the whole global config of the ecosystem to make changes and persist it as a whole.
type GlobalConfigRepository interface {
	// Get retrieves the whole global config.
	// It can throw the following errors:
	// 	- NotFoundError if the global config was not found.
	// 	- InternalError if any other error happens.
	Get(ctx context.Context) (config.GlobalConfig, error)
	// Update persists the whole global config.
	// It can throw the following errors:
	//  - NotFoundError if the global config was not found to update it.
	//  - ConflictError if there were concurrent write accesses.
	//  - InternalError if any other error happens.
	Update(ctx context.Context, config config.GlobalConfig) (config.GlobalConfig, error)
}

// DoguConfigRepository to get and update normal dogu config. The config is always handled as a whole.
type DoguConfigRepository interface {
	// Get retrieves the normal config for the given dogu.
	// It can throw the following errors:
	// 	- NotFoundError if the dogu config was not found.
	// 	- InternalError if any other error happens.
	Get(ctx context.Context, doguName common.SimpleDoguName) (config.DoguConfig, error)
	// GetAll retrieves the normal config for all given dogus as a map from doguName to config.
	// It can throw the following errors:
	// 	- NotFoundError if the dogu config was not found.
	// 	- InternalError if any other error happens.
	GetAll(ctx context.Context, doguNames []common.SimpleDoguName) (map[common.SimpleDoguName]config.DoguConfig, error)
	// GetAllExisting retrieves the normal config for all given dogus as a map from doguName to config and
	// includes not found configs as empty configs.
	// It can throw the following errors:
	// 	- InternalError if any other error happens.
	GetAllExisting(ctx context.Context, doguNames []common.SimpleDoguName) (map[common.SimpleDoguName]config.DoguConfig, error)
	// Update persists the whole given config.
	// It can throw the following errors:
	//  - NotFoundError if the dogu config was not found to update it.
	//  - ConflictError if there were concurrent write accesses.
	//  - InternalError if any other error happens.
	Update(ctx context.Context, config config.DoguConfig) (config.DoguConfig, error)
	// Create creates the data structure for the config and persists the given entries.
	// It can throw the following errors:
	//  - ConflictError if there already is a config.
	//  - InternalError if any other error happens.
	Create(ctx context.Context, config config.DoguConfig) (config.DoguConfig, error)
	// UpdateOrCreate updates the config if it already exists, otherwise creates it with the given content.
	// It can throw the following errors:
	//  - ConflictError if there already is a config.
	//  - InternalError if any other error happens.
	UpdateOrCreate(ctx context.Context, config config.DoguConfig) (config.DoguConfig, error)
}

// SensitiveDoguConfigRepository to get and update sensitive dogu config. The config is always handled as a whole.
type SensitiveDoguConfigRepository interface {
	DoguConfigRepository //interfaces are the same yet. Don't hesitate to split them if they diverge
}

// NewNotFoundError creates a NotFoundError with a given message. The wrapped error may be nil. The error message must
// omit the fmt.Errorf verb %w because this is done by NotFoundError.Error().
func NewNotFoundError(wrappedError error, message string, msgArgs ...any) *NotFoundError {
	return &NotFoundError{WrappedError: wrappedError, Message: fmt.Sprintf(message, msgArgs...)}
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

func IsNotFoundError(err error) bool {
	var notFoundError *NotFoundError
	return errors.As(err, &notFoundError)
}

// NewInternalError creates an InternalError with a given message. The wrapped error may be nil. The error message must
// omit the fmt.Errorf verb %w because this is done by InternalError.Error().
func NewInternalError(wrappedError error, message string, msgArgs ...any) *InternalError {
	return &InternalError{WrappedError: wrappedError, Message: fmt.Sprintf(message, msgArgs...)}
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

func IsInternalError(err error) bool {
	var internalError *InternalError
	return errors.As(err, &internalError)
}

// ConflictError is a common error indicating that the aggregate was modified in the meantime.
type ConflictError struct {
	WrappedError error
	Message      string
}

// NewConflictError creates an ConflictError with a given message. The wrapped error may be nil. The error message must
// omit the fmt.Errorf verb %w because this is done by ConflictError.Error().
func NewConflictError(wrappedError error, message string, msgArgs ...any) *ConflictError {
	return &ConflictError{WrappedError: wrappedError, Message: fmt.Sprintf(message, msgArgs...)}
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

func IsConflictError(err error) bool {
	var conflictError *ConflictError
	return errors.As(err, &conflictError)
}
