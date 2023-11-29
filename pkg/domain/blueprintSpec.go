package domain

import (
	"errors"
	"fmt"
)

type BlueprintSpec struct {
	Id                     string
	Blueprint              Blueprint
	BlueprintMask          BlueprintMask
	EffectiveBlueprint     EffectiveBlueprint
	StateDiff              StateDiff
	BlueprintUpgradePlan   BlueprintUpgradePlan
	BlueprintConfiguration BlueprintConfiguration
	Status                 StatusPhase
	Events                 []interface{}
}

type StatusPhase string

const (
	// StatusPhaseNew marks a newly created blueprint-CR.
	StatusPhaseNew StatusPhase = ""

	// StatusPhaseInvalid marks the given blueprint or the blueprint mask as not correct.
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

type InvalidBlueprintEvent struct {
	ValidationError error
}

func (spec *BlueprintSpec) Validate() error {
	var errorList []error

	if spec.Id == "" {
		errorList = append(errorList, errors.New("blueprint spec don't have an ID"))
	}
	errorList = append(errorList, spec.Blueprint.Validate())
	errorList = append(errorList, spec.BlueprintMask.Validate())
	err := errors.Join(errorList...)
	if err != nil {
		err = fmt.Errorf("blueprint spec is invalid: %w", err)
		spec.Status = StatusPhaseInvalid
		spec.Events = append(spec.Events, InvalidBlueprintEvent{ValidationError: err})
	}
	return err
}

type EffectiveBlueprintCalculatedEvent struct {
	effectiveBlueprint EffectiveBlueprint
}

func (spec *BlueprintSpec) CalculateEffectiveBlueprint() error {
	//TODO: do deep copy
	spec.EffectiveBlueprint = EffectiveBlueprint{
		Dogus:                   spec.calculateEffectiveDogus(),
		Components:              spec.Blueprint.Components,
		RegistryConfig:          spec.Blueprint.RegistryConfig,
		RegistryConfigAbsent:    spec.Blueprint.RegistryConfigAbsent,
		RegistryConfigEncrypted: spec.Blueprint.RegistryConfigEncrypted,
	}

	spec.Events = append(spec.Events, EffectiveBlueprintCalculatedEvent{effectiveBlueprint: spec.EffectiveBlueprint})
	return nil
}

func (spec *BlueprintSpec) calculateEffectiveDogus() []TargetDogu {
	var effectiveDogus []TargetDogu
	for _, dogu := range spec.Blueprint.Dogus {
		effectiveDogus = append(effectiveDogus, spec.calculateEffectiveDogu(dogu))
	}
	return effectiveDogus
}

func (spec *BlueprintSpec) calculateEffectiveDogu(dogu TargetDogu) TargetDogu {
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
		effectiveDogu.Namespace = maskDogu.Namespace
		effectiveDogu.TargetState = maskDogu.TargetState
	}
	return effectiveDogu
}
