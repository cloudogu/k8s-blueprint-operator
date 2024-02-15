package domainservice

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"

	"github.com/cloudogu/cesapp-lib/core"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
)

type DoguInstallationRepository interface {
	// GetByName returns the ecosystem.DoguInstallation or
	// a NotFoundError if the dogu is not installed or
	// an InternalError if there is any other error.
	GetByName(ctx context.Context, doguName common.SimpleDoguName) (*ecosystem.DoguInstallation, error)
	// GetAll returns the installation info of all installed dogus or
	// a NotFoundError if any dogu is not installed or
	// an InternalError if there is any other error.
	GetAll(ctx context.Context) (map[common.SimpleDoguName]*ecosystem.DoguInstallation, error)
	// Create saves a new ecosystem.DoguInstallation. This initiates a dogu installation. It returns
	// a ConflictError if there is already a DoguInstallation with this name or
	// an InternalError if there is any error while saving the DoguInstallation
	Create(ctx context.Context, dogu *ecosystem.DoguInstallation) error
	// Update updates an ecosystem.DoguInstallation in the ecosystem.
	// returns a ConflictError if there were changes on the DoguInstallation in the meantime or
	// TODO: also return NotFoundErrors? Does k8s supply this error to us?
	// returns an InternalError if there is any other error
	Update(ctx context.Context, dogu *ecosystem.DoguInstallation) error
	// Delete removes the given ecosystem.DoguInstallation completely from the ecosystem.
	// We delete DoguInstallations with the object not just the name as this way we can detect concurrent updates.
	// returns a ConflictError if there were changes on the DoguInstallation in the meantime or
	// returns an InternalError if there is any other error
	Delete(ctx context.Context, doguName common.SimpleDoguName) error
}

type ComponentInstallationRepository interface {
	// GetByName returns the ecosystem.ComponentInstallation or
	// a NotFoundError if the component is not installed or
	// an InternalError if there is any other error.
	GetByName(ctx context.Context, componentName common.SimpleComponentName) (*ecosystem.ComponentInstallation, error)
	// GetAll returns the installation info of all installed components or
	// a NotFoundError if any component is not installed or
	// an InternalError if there is any other error.
	GetAll(ctx context.Context) (map[common.SimpleComponentName]*ecosystem.ComponentInstallation, error)
	// Delete deletes the component by name from the ecosystem.
	// returns an InternalError if there is an error.
	Delete(ctx context.Context, componentName common.SimpleComponentName) error
	// Create creates the ecosystem.ComponentInstallation in the cluster.
	// returns an InternalError if there is an error.
	Create(ctx context.Context, component *ecosystem.ComponentInstallation) error
	// Update updates the ecosystem.ComponentInstallation with a patch operation in the cluster.
	// returns an InternalError on patch error.
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
	// Activate enables the maintenance mode with the given page model.
	// May throw a generic InternalError or a ConflictError if another party activated the maintenance mode.
	Activate(content MaintenancePageModel) error
	// Deactivate disables the maintenance mode if it is active.
	// May throw a generic InternalError or a ConflictError if another party holds the maintenance mode lock.
	Deactivate() error
}

// MaintenancePageModel contains data that gets displayed when the maintenance mode is active.
type MaintenancePageModel struct {
	Title string
	Text  string
}

type ConfigEncryptionAdapter interface {
	// Encrypt encrypts the given value for a dogu.
	// It can throw an InternalError if the encryption did not succeed, public key is missing or config store is not reachable.
	Encrypt(context.Context, common.SimpleDoguName, common.SensitiveDoguConfigValue) (common.EncryptedDoguConfigValue, error)
	// EncryptAll encrypts the given values for a dogu.
	// It can throw an InternalError if the encryption did not succeed, public key is missing or config store is not reachable.
	EncryptAll(context.Context, map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue) (map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue, error)
}

type GlobalConfigEntryRepository interface {
	// GetAll retrieves all keys from the global config.
	// It can throw the following errors:
	// 	- NotFoundError if there is no global config.
	// 	- InternalError if any other error happens.
	GetAll(ctx context.Context) ([]*ecosystem.GlobalConfigEntry, error)
	// Save persists the global config.
	// It can throw the following errors:
	//  - ConflictError if there were concurrent write accesses.
	//  - InternalError if any other error happens.
	Save(context.Context, *ecosystem.GlobalConfigEntry) error
	// SaveAll persists all given config keys.
	// It can throw the following errors:
	//	- ConflictError if there were concurrent write accesses.
	//	- InternalError if any other error happens.
	SaveAll(ctx context.Context, keys []*ecosystem.GlobalConfigEntry) error
	// Delete deletes a global config key.
	// It can throw an InternalError if any error happens.
	// If the key is not existent no error will be returned.
	Delete(ctx context.Context, key common.SimpleDoguName) error
}

type DoguConfigEntryRepository interface {
	// GetAllByKey retrieves all ecosystem.DoguConfigEntry's for the given ecosystem.DoguConfigKey's config.
	// It can trow the following errors:
	// 	- NotFoundError if there is no config for the dogu.
	// 	- InternalError if any other error happens.
	GetAllByKey(ctx context.Context, keys []common.DoguConfigKey) (map[common.SimpleDoguName][]*ecosystem.DoguConfigEntry, error)
	// Save persists the config for the given dogu. Config can be set even if the dogu is not yet installed.
	// It can throw the following errors:
	//	- ConflictError if there were concurrent write accesses.
	//	- InternalError if any other error happens.
	Save(context.Context, *ecosystem.DoguConfigEntry) error
	// SaveAll persists all given config keys. Configs can be set even if the dogus are not installed.
	// It can throw the following errors:
	//	- ConflictError if there were concurrent write accesses.
	//	- InternalError if any other error happens.
	SaveAll(ctx context.Context, keys []*ecosystem.DoguConfigEntry) error
	// Delete deletes a dogu config key.
	// It can throw an InternalError if any error happens.
	// If the key is not existent no error will be returned.
	Delete(ctx context.Context, key common.DoguConfigKey) error
}

type SensitiveDoguConfigEntryRepository interface {
	// GetAllByKey retrieves a dogu's sensitive config for the given keys.
	// It can trow the following errors:
	// 	- NotFoundError if there is no config for the dogu.
	// 	- InternalError if any other error happens.
	GetAllByKey(ctx context.Context, keys []common.SensitiveDoguConfigKey) (map[common.SimpleDoguName][]*ecosystem.SensitiveDoguConfigEntry, error)
	// Save persists the sensitive config for the given dogu. Config can be set even if the dogu is not yet installed.
	// It can throw the following errors:
	//	- ConflictError if there were concurrent write accesses.
	//	- InternalError if any other error happens.
	Save(context.Context, *ecosystem.SensitiveDoguConfigEntry) error
	// SaveAll persists all given sensitive config keys. sensitive configs can be set even if the dogus are not installed.
	// It can throw the following errors:
	//	- ConflictError if there were concurrent write accesses.
	//	- InternalError if any other error happens.
	SaveAll(ctx context.Context, keys []*ecosystem.SensitiveDoguConfigEntry) error
	// Delete deletes a sensitive dogu config key.
	// It can throw an InternalError if any error happens.
	// If the key is not existent no error will be returned.
	Delete(ctx context.Context, key common.SensitiveDoguConfigKey) error
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

// NewInternalError creates an InternalError with a given message. The wrapped error may be nil. The error message must
// omit the fmt.Errorf verb %w because this is done by InternalError.Error().
func NewInternalError(wrappedError error, message string, msgArgs ...interface{}) *InternalError {
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
