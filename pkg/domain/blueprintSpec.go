package domain

import (
	"errors"
	"fmt"
)

type BlueprintSpec struct {
	Id                   string
	Blueprint            Blueprint
	BlueprintMask        BlueprintMask
	EffectiveBlueprint   EffectiveBlueprint
	StateDiff            StateDiff
	BlueprintUpgradePlan BlueprintUpgradePlan
	config               BlueprintConfiguration
	Status               StatusPhase
	Events               []interface{}
}

type StatusPhase string

const (
	// StatusPhaseNew marks a newly created blueprint-CR.
	StatusPhaseNew StatusPhase = ""
	// StatusPhaseValidated marks the given blueprint spec as validated.
	StatusPhaseValidated StatusPhase = "validated"
	// StatusPhaseInvalid marks the given blueprint spec is semantically incorrect.
	StatusPhaseInvalid StatusPhase = "invalid"
	// StatusPhaseInProgress marks that the blueprint is currently being processed.
	StatusPhaseInProgress StatusPhase = "inProgress"
	// StatusPhaseFailed marks that an error occurred during processing of the blueprint.
	StatusPhaseFailed StatusPhase = "failed"
	// StatusPhaseCompleted marks the blueprint as successfully applied.
	StatusPhaseCompleted StatusPhase = "completed"
)

type BlueprintConfiguration struct {
	// Force blueprint upgrade even when a dogu is unhealthy
	ignoreDoguHealth bool
	// allowNamespaceSwitch allows the blueprint upgrade to switch a dogus namespace
	allowDoguNamespaceSwitch bool
}

type StateDiff struct{}

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

type BlueprintSpecInvalidEvent struct {
	ValidationError error
}

type BlueprintSpecValidatedEvent struct{}

// Validate checks the blueprintSpec for semantic errors and sets the status to the result.
// returns a domain.InvalidBlueprintError if blueprint is invalid
// or nil otherwise.
func (spec *BlueprintSpec) Validate() error {
	switch spec.Status {
	case StatusPhaseNew: //continue
	case StatusPhaseInvalid: //do not validate again
		return errors.New("blueprint spec was marked invalid before. Do not revalidate")
	default: //do not validate again. for all other status it must be either status validated or a status beyond that
		return nil
	}
	var errorList []error

	if spec.Id == "" {
		errorList = append(errorList, errors.New("blueprint spec don't have an ID"))
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
		spec.Status = StatusPhaseValidated
		spec.Events = append(spec.Events, BlueprintSpecValidatedEvent{})
	}
	return err
}

type EffectiveBlueprintCalculatedEvent struct {
	effectiveBlueprint EffectiveBlueprint
}

func (spec *BlueprintSpec) validateMaskAgainstBlueprint() error {
	var errorList []error
	for _, doguMask := range spec.BlueprintMask.Dogus {
		dogu, noDoguFoundError := FindDoguByName(spec.Blueprint.Dogus, doguMask.Name)
		if noDoguFoundError != nil {
			errorList = append(errorList, fmt.Errorf("dogu %s is missing in the blueprint", doguMask.Name))
		}
		if !spec.config.allowDoguNamespaceSwitch && dogu.Namespace != doguMask.Namespace {
			errorList = append(errorList, fmt.Errorf(
				"namespace switch is not allowed by default for dogu %s. Activate feature flag for that", doguMask.Name),
			)
		}
	}

	err := errors.Join(errorList...)
	if err != nil {
		err = fmt.Errorf("blueprint mask does not match the blueprint: %w", err)
	}
	return err
}

func (spec *BlueprintSpec) CalculateEffectiveBlueprint() error {
	switch spec.Status {
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

	spec.Events = append(spec.Events, EffectiveBlueprintCalculatedEvent{effectiveBlueprint: spec.EffectiveBlueprint})
	return nil
}

func (spec *BlueprintSpec) calculateEffectiveDogus() ([]TargetDogu, error) {
	var effectiveDogus []TargetDogu
	for _, dogu := range spec.Blueprint.Dogus {
		effectiveDogu, err := spec.calculateEffectiveDogu(dogu)
		if err != nil {
			return nil, err
		}
		effectiveDogus = append(effectiveDogus, effectiveDogu)
	}
	return effectiveDogus, nil
}

func (spec *BlueprintSpec) calculateEffectiveDogu(dogu TargetDogu) (TargetDogu, error) {
	effectiveDogu := TargetDogu{
		Namespace:   dogu.Namespace,
		Name:        dogu.Name,
		Version:     dogu.Version,
		TargetState: dogu.TargetState,
	}
	maskDogu, noMaskDoguErr := spec.BlueprintMask.FindDoguByName(dogu.Name)
	if noMaskDoguErr == nil {
		if maskDogu.Version != "" {
			effectiveDogu.Version = maskDogu.Version
		}
		if maskDogu.Namespace != dogu.Namespace {
			if spec.config.allowDoguNamespaceSwitch {
				effectiveDogu.Namespace = maskDogu.Namespace
			} else {
				return TargetDogu{}, errors.New("changing the dogu namespace is only allowed with the changeDoguNamespace flag")
			}
		}
		effectiveDogu.TargetState = maskDogu.TargetState
	}

	return effectiveDogu, nil
}
