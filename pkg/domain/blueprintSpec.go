package domain

import (
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

type BlueprintSpec struct {
	Id                   string
	Blueprint            Blueprint
	BlueprintMask        BlueprintMask
	EffectiveBlueprint   EffectiveBlueprint
	StateDiff            StateDiff
	BlueprintUpgradePlan BlueprintUpgradePlan
	Config               BlueprintConfiguration
	Status               StatusPhase
	// PersistenceContext can hold generic values needed for persistence with repositories, e.g. version counters or transaction contexts.
	// This field has a generic map type as the values within it highly depend on the used type of repository.
	// This field should be ignored in the whole domain.
	PersistenceContext map[string]interface{}
	Events             []Event
}

type StatusPhase string

const (
	// StatusPhaseNew marks a newly created blueprint-CR.
	StatusPhaseNew StatusPhase = ""
	// StatusPhaseStaticallyValidated marks the given blueprint spec as validated.
	StatusPhaseStaticallyValidated StatusPhase = "staticallyValidated"
	// StatusPhaseValidated marks the given blueprint spec as validated.
	StatusPhaseValidated StatusPhase = "validated"
	// StatusPhaseEffectiveBlueprintGenerated marks that the effective blueprint was generated out of the blueprint and the mask.
	StatusPhaseEffectiveBlueprintGenerated StatusPhase = "effectiveBlueprintGenerated"
	// StatusPhaseStateDiffDetermined marks that the diff to the ecosystem state was successfully determined.
	StatusPhaseStateDiffDetermined StatusPhase = "stateDiffDetermined"
	// StatusPhaseEcosystemHealthyUpfront marks that all currently installed dogus are healthy.
	StatusPhaseEcosystemHealthyUpfront StatusPhase = "ecosystemHealthyUpfront"
	// StatusPhaseIgnoreDoguHealth marks that dogu health checks have been skipped.
	StatusPhaseIgnoreDoguHealth StatusPhase = "ignoreDoguHealth"
	// StatusPhaseEcosystemUnhealthyUpfront marks that some currently installed dogus are unhealthy.
	StatusPhaseEcosystemUnhealthyUpfront StatusPhase = "dogusUnhealthy"
	// StatusPhaseInvalid marks the given blueprint spec is semantically incorrect.
	StatusPhaseInvalid StatusPhase = "invalid"
	// StatusPhaseInProgress marks that the blueprint is currently being processed.
	StatusPhaseInProgress StatusPhase = "inProgress"
	// StatusPhaseWaitForHealthyEcosystem indicates that the blueprint was applied but the ecosystem is not healthy yet.
	StatusPhaseWaitForHealthyEcosystem StatusPhase = "waitForHealthyEcosystem"
	// StatusPhaseFailed marks that an error occurred during processing of the blueprint.
	StatusPhaseFailed StatusPhase = "failed"
	// StatusPhaseCompleted marks the blueprint as successfully applied.
	StatusPhaseCompleted StatusPhase = "completed"
)

type BlueprintConfiguration struct {
	// Force blueprint upgrade even when a dogu is unhealthy
	IgnoreDoguHealth bool
	// allowNamespaceSwitch allows the blueprint upgrade to switch a dogus namespace
	AllowDoguNamespaceSwitch bool
}

type BlueprintUpgradePlan struct {
	DogusToInstall   []string
	DogusToUpgrade   []string
	DogusToUninstall []string

	ComponentsToInstall   []string
	ComponentsToUpgrade   []string
	ComponentsToUninstall []string

	RegistryConfigToAdd    []string
	RegistryConfigToUpdate []string
	RegistryConfigToRemove []string
}

// ValidateStatically checks the blueprintSpec for semantic errors and sets the status to the result.
// Here will be only checked, what can be checked without any external information, e.g. without dogu specification.
// returns a domain.InvalidBlueprintError if blueprint is invalid
// or nil otherwise.
func (spec *BlueprintSpec) ValidateStatically() error {
	switch spec.Status {
	case StatusPhaseNew: //continue
	case StatusPhaseInvalid: //do not validate again
		return &InvalidBlueprintError{Message: "blueprint spec was marked invalid before: do not revalidate"}
	default: //do not validate again. for all other status it must be either status validated or a status beyond that
		return nil
	}
	var errorList []error

	if spec.Id == "" {
		errorList = append(errorList, errors.New("blueprint spec doesn't have an ID"))
	}
	errorList = append(errorList, spec.Blueprint.Validate())
	errorList = append(errorList, spec.BlueprintMask.Validate())
	errorList = append(errorList, spec.validateMaskAgainstBlueprint())
	err := errors.Join(errorList...)
	if err != nil {
		err = &InvalidBlueprintError{
			WrappedError: err,
			Message:      "blueprint spec is invalid",
		}
		spec.Status = StatusPhaseInvalid
		spec.Events = append(spec.Events, BlueprintSpecInvalidEvent{ValidationError: err})
	} else {
		spec.Status = StatusPhaseStaticallyValidated
		spec.Events = append(spec.Events, BlueprintSpecStaticallyValidatedEvent{})
	}
	return err
}

func (spec *BlueprintSpec) validateMaskAgainstBlueprint() error {
	var errorList []error
	for _, doguMask := range spec.BlueprintMask.Dogus {
		dogu, noDoguFoundError := FindDoguByName(spec.Blueprint.Dogus, doguMask.Name)
		if noDoguFoundError != nil {
			errorList = append(errorList, fmt.Errorf("dogu %q is missing in the blueprint", doguMask.Name))
		}
		if !spec.Config.AllowDoguNamespaceSwitch && dogu.Namespace != doguMask.Namespace {
			errorList = append(errorList, fmt.Errorf(
				"namespace switch is not allowed by default for dogu %q: activate the feature flag for that", doguMask.Name),
			)
		}
	}

	err := errors.Join(errorList...)
	if err != nil {
		err = fmt.Errorf("blueprint mask does not match the blueprint: %w", err)
	}
	return err
}

// ValidateDynamically sets the Status either to StatusPhaseInvalid or StatusPhaseValidated
// depending on if the dependencies or versions of the elements in the blueprint are invalid.
// returns a domain.InvalidBlueprintError if blueprint is invalid
// or nil otherwise.
func (spec *BlueprintSpec) ValidateDynamically(possibleInvalidDependenciesError error) {
	if possibleInvalidDependenciesError != nil {
		err := &InvalidBlueprintError{
			WrappedError: possibleInvalidDependenciesError,
			Message:      "blueprint spec is invalid",
		}
		spec.Status = StatusPhaseInvalid
		spec.Events = append(spec.Events, BlueprintSpecInvalidEvent{ValidationError: err})
	} else {
		spec.Status = StatusPhaseValidated
		spec.Events = append(spec.Events, BlueprintSpecValidatedEvent{})
	}
}

func (spec *BlueprintSpec) CalculateEffectiveBlueprint() error {
	switch spec.Status {
	case StatusPhaseEffectiveBlueprintGenerated:
		return nil //do not regenerate effective blueprint
	case StatusPhaseNew: // stop
		return fmt.Errorf("cannot calculate effective blueprint before the blueprint spec is validated")
	case StatusPhaseInvalid: // stop
		return fmt.Errorf("cannot calculate effective blueprint on invalid blueprint spec")
	default: //continue: StatusPhaseValidated, StatusPhaseInProgress, StatusPhaseFailed, StatusPhaseCompleted
	}

	effectiveDogus, err := spec.calculateEffectiveDogus()
	if err != nil {
		return err
	}

	spec.EffectiveBlueprint = EffectiveBlueprint{
		Dogus:                   effectiveDogus,
		Components:              spec.Blueprint.Components,
		RegistryConfig:          spec.Blueprint.RegistryConfig,
		RegistryConfigAbsent:    spec.Blueprint.RegistryConfigAbsent,
		RegistryConfigEncrypted: spec.Blueprint.RegistryConfigEncrypted,
	}
	spec.Status = StatusPhaseEffectiveBlueprintGenerated
	spec.Events = append(spec.Events, EffectiveBlueprintCalculatedEvent{})
	return nil
}

func (spec *BlueprintSpec) calculateEffectiveDogus() ([]Dogu, error) {
	var effectiveDogus []Dogu
	for _, dogu := range spec.Blueprint.Dogus {
		effectiveDogu, err := spec.calculateEffectiveDogu(dogu)
		if err != nil {
			return nil, err
		}
		effectiveDogus = append(effectiveDogus, effectiveDogu)
	}
	return effectiveDogus, nil
}

func (spec *BlueprintSpec) calculateEffectiveDogu(dogu Dogu) (Dogu, error) {
	effectiveDogu := Dogu{
		Namespace:   dogu.Namespace,
		Name:        dogu.Name,
		Version:     dogu.Version,
		TargetState: dogu.TargetState,
	}
	maskDogu, noMaskDoguErr := spec.BlueprintMask.FindDoguByName(dogu.Name)
	if noMaskDoguErr == nil {
		emptyVersion := core.Version{}
		if maskDogu.Version != emptyVersion {
			effectiveDogu.Version = maskDogu.Version
		}
		if maskDogu.Namespace != dogu.Namespace {
			if spec.Config.AllowDoguNamespaceSwitch {
				effectiveDogu.Namespace = maskDogu.Namespace
			} else {
				return Dogu{}, fmt.Errorf(
					"changing the dogu namespace is forbidden by default and can be allowed by a flag: %q -> %q", dogu.GetQualifiedName(), maskDogu.GetQualifiedName())
			}
		}
		effectiveDogu.TargetState = maskDogu.TargetState
	}

	return effectiveDogu, nil
}

// MarkInvalid is used to mark the blueprint as invalid after dynamically validating it.
func (spec *BlueprintSpec) MarkInvalid(err error) {
	spec.Status = StatusPhaseInvalid
	spec.Events = append(spec.Events, BlueprintSpecInvalidEvent{ValidationError: err})
}

// DetermineStateDiff creates the StateDiff between the blueprint and the actual state of the ecosystem.
// if sth. is not in the lists of installed things, it is considered not installed.
// installedDogus are a map in the form of simpleDoguName->*DoguInstallation. There should be no nil values.
// The StateDiff is an 'as is' representation, therefore no error is thrown, e.g. if dogu namespaces are different and namespace changes are not allowed.
// If there are not allowed actions should be considered at the start of the execution of the blueprint.
// returns an error if the BlueprintSpec is not in the necessary state to determine the stateDiff.
func (spec *BlueprintSpec) DetermineStateDiff(installedDogus map[string]*ecosystem.DoguInstallation) error {
	switch spec.Status {
	case StatusPhaseNew:
		fallthrough
	case StatusPhaseStaticallyValidated:
		fallthrough
	case StatusPhaseEffectiveBlueprintGenerated:
		return fmt.Errorf("cannot determine state diff in status phase %q", spec.Status)
	case StatusPhaseValidated: // this is the state, the blueprint spec should be
	default:
		return nil // do not re-determine the state diff from status StatusPhaseStateDiffDetermined and above
	}
	spec.StateDiff = StateDiff{
		DoguDiffs: determineDoguDiffs(spec.EffectiveBlueprint.Dogus, installedDogus),
		// there will be more diffs, e.g. for components and registry keys
	}
	spec.Status = StatusPhaseStateDiffDetermined
	spec.Events = append(spec.Events, StateDiffDeterminedEvent{StateDiff: spec.StateDiff})
	return nil
}

func (spec *BlueprintSpec) CheckEcosystemHealthUpfront(healthResult ecosystem.HealthResult) {
	if spec.Config.IgnoreDoguHealth {
		spec.Status = StatusPhaseIgnoreDoguHealth
		spec.Events = append(spec.Events, IgnoreEcosystemHealthEvent{})
		return
	}

	if healthResult.AllHealthy() {
		spec.Status = StatusPhaseEcosystemHealthyUpfront
		spec.Events = append(spec.Events, EcosystemHealthyEvent{})
	} else {
		spec.Status = StatusPhaseEcosystemUnhealthyUpfront
		spec.Events = append(spec.Events, EcosystemUnhealthyUpfrontEvent{HealthResult: healthResult})
	}
}

func (spec *BlueprintSpec) MarkInProgress() {
	spec.Status = StatusPhaseInProgress
	spec.Events = append(spec.Events, InProgressEvent{})
}

func (spec *BlueprintSpec) MarkFailed(err error) {
	spec.Status = StatusPhaseFailed
	spec.Events = append(spec.Events, ExecutionFailedEvent{err: err})
}

func (spec *BlueprintSpec) MarkWaitingForHealthyEcosystem() {
	spec.Status = StatusPhaseWaitForHealthyEcosystem
	spec.Events = append(spec.Events, BlueprintAppliedEvent{})
}

func (spec *BlueprintSpec) MarkCompleted() {
	spec.Status = StatusPhaseCompleted
	spec.Events = append(spec.Events, CompletedEvent{})
}
