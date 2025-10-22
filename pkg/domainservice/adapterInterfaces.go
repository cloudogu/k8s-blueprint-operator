package domainservice

import (
	"context"
	"errors"
	"fmt"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-registry-lib/config"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

type DoguInstallationRepository interface {
	// GetByName returns the ecosystem.DoguInstallation or
	//  - a NotFoundError if the dogu is not installed or
	//  - an InternalError if there is any other error.
	GetByName(ctx context.Context, doguName cescommons.SimpleName) (*ecosystem.DoguInstallation, error)
	// GetAll returns the installation info of all installed dogus or
	//  - an InternalError if there is any other error.
	GetAll(ctx context.Context) (map[cescommons.SimpleName]*ecosystem.DoguInstallation, error)
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
	Delete(ctx context.Context, doguName cescommons.SimpleName) error
}

type BlueprintSpecRepository interface {
	// GetById returns a BlueprintSpec identified by its ID or
	// a NotFoundError if the BlueprintSpec was not found or
	// a domain.InvalidBlueprintError together with a BlueprintSpec without blueprint and mask if the BlueprintSpec could not be parsed or
	// an InternalError if there is any other error.
	GetById(ctx context.Context, blueprintId string) (*domain.BlueprintSpec, error)

	// Count counts the Blueprint-resources in the namespace of the repository up to the given limit and
	//  - returns the amount of blueprints or
	//  - returns an InternalError if there is any error, e.g. a connection error.
	Count(ctx context.Context, limit int) (int, error)

	// ListIds retrieves all Blueprint-Ids.
	//  - It returns a list of Ids containing all blueprint Ids, or
	//  - an InternalError if the operation fails.
	ListIds(ctx context.Context) ([]string, error)

	// Update updates a given BlueprintSpec.
	// returns a ConflictError if there were changes on the BlueprintSpec in the meantime or
	// returns an InternalError if there is any other error
	Update(ctx context.Context, blueprintSpec *domain.BlueprintSpec) error
}

type RemoteDoguRegistry interface {
	// GetDogu returns the dogu specification for the given dogu and version or
	// an NotFoundError indicating that there was no dogu spec found or
	// an InternalError indicating that the caller has no fault.
	GetDogu(ctx context.Context, qualifiedDoguVersion cescommons.QualifiedVersion) (*core.Dogu, error)

	// GetDogus returns the all requested dogu specifications or
	// an NotFoundError indicating that any dogu spec was not found or
	// an InternalError indicating that the caller has no fault.
	GetDogus(ctx context.Context, dogusToLoad []cescommons.QualifiedVersion) (map[cescommons.QualifiedName]*core.Dogu, error)
}

type DoguToLoad struct {
	DoguName cescommons.QualifiedName
	Version  string
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
	Get(ctx context.Context, doguName cescommons.SimpleName) (config.DoguConfig, error)
	// GetAll retrieves the normal config for all given dogus as a map from doguName to config.
	// It can throw the following errors:
	// 	- NotFoundError if the dogu config was not found.
	// 	- InternalError if any other error happens.
	GetAll(ctx context.Context, doguNames []cescommons.SimpleName) (map[cescommons.SimpleName]config.DoguConfig, error)
	// GetAllExisting retrieves the normal config for all given dogus as a map from doguName to config and
	// includes not found configs as empty configs.
	// It can throw the following errors:
	// 	- InternalError if any other error happens.
	GetAllExisting(ctx context.Context, doguNames []cescommons.SimpleName) (map[cescommons.SimpleName]config.DoguConfig, error)
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

// SensitiveConfigRefReader resolves given domain.SensitiveValueRef's and loads the referenced values.
// As sensitive config should not be written directly into the blueprint, only references are used.
// This way, the blueprint e.g. can securely get committed to git.
type SensitiveConfigRefReader interface {
	// GetValues reads all common.SensitiveDoguConfigValue's from the given domain.SensitiveValueRef's by common.SensitiveDoguConfigKey.
	// It can throw the following errors:
	//  - NotFoundError if any reference cannot be resolved.
	//  - InternalError if any other error happens.
	GetValues(
		ctx context.Context,
		refs map[common.DoguConfigKey]domain.SensitiveValueRef,
	) (
		map[common.DoguConfigKey]common.SensitiveDoguConfigValue,
		error,
	)
}

type DebugModeRepository interface {
	// GetSingleton returns the ecosystem.DebugMode or
	//  - a NotFoundError if the debugMode information is not found in the ecosystem or
	//  - an InternalError if there is any other error.
	GetSingleton(ctx context.Context) (*ecosystem.DebugMode, error)
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
	DoNotRetry   bool
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
