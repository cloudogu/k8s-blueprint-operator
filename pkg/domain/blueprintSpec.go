package domain

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BlueprintSpec struct {
	Id                 string
	DisplayName        string
	Blueprint          Blueprint
	BlueprintMask      BlueprintMask
	EffectiveBlueprint EffectiveBlueprint
	StateDiff          StateDiff
	Config             BlueprintConfiguration
	Conditions         []Condition
	// PersistenceContext can hold generic values needed for persistence with repositories, e.g. version counters or transaction contexts.
	// This field has a generic map type as the values within it highly depend on the used type of repository.
	// This field should be ignored in the whole domain.
	PersistenceContext map[string]interface{}
	Events             []Event
}

type Condition = metav1.Condition

const (
	ConditionValid              = "Valid"
	ConditionExecutable         = "Executable"
	ConditionEcosystemHealthy   = "EcosystemHealthy"
	ConditionCompleted          = "Completed"
	ConditionLastApplySucceeded = "LastApplySucceeded"

	ReasonLastApplyErrorAtDogus  = "DoguApplyFailure"
	ReasonLastApplyErrorAtConfig = "ConfigApplyFailure"

	loggingKey = "logging/root"
)

var (
	BlueprintConditions = []string{ConditionValid, ConditionExecutable, ConditionEcosystemHealthy, ConditionCompleted, ConditionLastApplySucceeded}

	// ActionSwitchDoguNamespace is an exception and should be handled with the blueprint config.
	notAllowedDoguActions = []Action{ActionDowngrade, ActionSwitchDoguNamespace}
)

type BlueprintConfiguration struct {
	// IgnoreDoguHealth forces blueprint upgrades even if dogus are unhealthy
	IgnoreDoguHealth bool
	// AllowDoguNamespaceSwitch allows the blueprint upgrade to switch a dogus namespace
	AllowDoguNamespaceSwitch bool
	// Stopped lets the user test a blueprint run to check if all attributes of the blueprint are correct and avoid a result with a failure state.
	Stopped bool
}

// ValidateStatically checks the blueprintSpec for semantic errors and sets the status to the result.
// Here will be only checked, what can be checked without any external information, e.g. without dogu specification.
// returns a domain.InvalidBlueprintError if blueprint is invalid
// or nil otherwise.
func (spec *BlueprintSpec) ValidateStatically() error {
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
		spec.Events = append(spec.Events, BlueprintSpecInvalidEvent{ValidationError: err})
		meta.SetStatusCondition(&spec.Conditions, metav1.Condition{
			Type:    ConditionValid,
			Status:  metav1.ConditionFalse,
			Reason:  "Invalid",
			Message: err.Error(),
		})
	}
	// Do not set condition to true here.
	// We reuse the condition for the dynamic validation.
	// If the blueprint is completely consistent and valid can only be decided there
	return err
}

func (spec *BlueprintSpec) validateMaskAgainstBlueprint() error {
	var errorList []error
	for _, doguMask := range spec.BlueprintMask.Dogus {
		dogu, found := FindDoguByName(spec.Blueprint.Dogus, cescommons.SimpleName(doguMask.Name.SimpleName))
		if !found {
			errorList = append(errorList, fmt.Errorf("dogu %q is missing in the blueprint", doguMask.Name))
		}
		if !doguMask.Absent && dogu.Absent {
			errorList = append(errorList, fmt.Errorf("absent dogu %q cannot be present in blueprint mask", dogu.Name.SimpleName))
		}
		if !spec.Config.AllowDoguNamespaceSwitch && dogu.Name.Namespace != doguMask.Name.Namespace {
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

// ValidateDynamically sets the ConditionValid
// depending on if the dependencies or versions of the elements in the blueprint are invalid.
// This function decides completely on the given error, therefore no error will be returned explicitly again.
func (spec *BlueprintSpec) ValidateDynamically(possibleInvalidDependenciesError error) {
	if possibleInvalidDependenciesError != nil {
		err := &InvalidBlueprintError{
			WrappedError: possibleInvalidDependenciesError,
			Message:      "blueprint spec is invalid",
		}
		conditionChanged := meta.SetStatusCondition(&spec.Conditions, metav1.Condition{
			Type:    ConditionValid,
			Status:  metav1.ConditionFalse,
			Reason:  "Inconsistent",
			Message: err.Error(),
		})
		if conditionChanged {
			spec.Events = append(spec.Events, BlueprintSpecInvalidEvent{ValidationError: err})
		}
	} else {
		meta.SetStatusCondition(&spec.Conditions, metav1.Condition{
			Type:   ConditionValid,
			Status: metav1.ConditionTrue,
			Reason: "Valid",
		})
	}
}

func (spec *BlueprintSpec) CalculateEffectiveBlueprint() error {
	effectiveDogus, err := spec.calculateEffectiveDogus()
	if err != nil {
		return err
	}

	effectiveConfig := spec.removeConfigForMaskedDogus()

	spec.EffectiveBlueprint = EffectiveBlueprint{
		Dogus:  effectiveDogus,
		Config: effectiveConfig,
	}
	validationError := spec.EffectiveBlueprint.validateOnlyConfigForDogusInBlueprint()
	if validationError != nil {
		conditionChanged := meta.SetStatusCondition(&spec.Conditions, metav1.Condition{
			Type:    ConditionValid,
			Status:  metav1.ConditionFalse,
			Reason:  "Inconsistent",
			Message: validationError.Error(),
		})
		if conditionChanged {
			spec.Events = append(spec.Events, BlueprintSpecInvalidEvent{ValidationError: validationError})
		}
		return validationError
	}
	return nil
}

// removeLogLevelChangesFromConfig creates a copy of the given config with all
// logging configuration entries removed from dogu configs. This is used in
// debug mode to prevent log level changes from being applied.
func removeLogLevelChangesFromConfig(config Config) Config {
	configCopy := Config{
		Dogus:  make(map[cescommons.SimpleName]DoguConfigEntries),
		Global: config.Global,
	}

	for doguName, doguConfig := range config.Dogus {
		newDoguConfig := make(DoguConfigEntries, 0, len(doguConfig))
		for _, configEntry := range doguConfig {
			if configEntry.Key != loggingKey {
				newDoguConfig = append(newDoguConfig, configEntry)
			}
		}
		configCopy.Dogus[doguName] = newDoguConfig
	}
	return configCopy
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
		Name:               dogu.Name,
		Version:            dogu.Version,
		Absent:             dogu.Absent,
		MinVolumeSize:      dogu.MinVolumeSize,
		ReverseProxyConfig: dogu.ReverseProxyConfig,
		AdditionalMounts:   dogu.AdditionalMounts,
	}
	maskDogu, noMaskDoguErr := spec.BlueprintMask.FindDoguByName(dogu.Name.SimpleName)
	if noMaskDoguErr == nil {
		emptyVersion := core.Version{}
		if maskDogu.Version != emptyVersion {
			effectiveDogu.Version = &maskDogu.Version
		}
		if maskDogu.Name.Namespace != dogu.Name.Namespace {
			if spec.Config.AllowDoguNamespaceSwitch {
				effectiveDogu.Name.Namespace = maskDogu.Name.Namespace
			} else {
				return Dogu{}, fmt.Errorf(
					"changing the dogu namespace is forbidden by default and can be allowed by a flag: %q -> %q", dogu.Name, maskDogu.Name)
			}
		}
		effectiveDogu.Absent = maskDogu.Absent
	}

	return effectiveDogu, nil
}

// It is not allowed to have config without the corresponding dogu, so this will clean up the unnecessary config.
func (spec *BlueprintSpec) removeConfigForMaskedDogus() Config {
	effectiveDoguConfig := maps.Clone(spec.Blueprint.Config.Dogus)

	for _, dogu := range spec.BlueprintMask.Dogus {
		if dogu.Absent {
			delete(effectiveDoguConfig, dogu.Name.SimpleName)
		}
	}

	return Config{
		Dogus:  effectiveDoguConfig,
		Global: spec.Blueprint.Config.Global,
	}
}

// MissingConfigReferences adds a given error, which was caused during preparations for determining the state diff
func (spec *BlueprintSpec) MissingConfigReferences(error error) {
	conditionChanged := meta.SetStatusCondition(&spec.Conditions, metav1.Condition{
		Type:    ConditionExecutable,
		Status:  metav1.ConditionFalse,
		Reason:  "MissingConfigReferences",
		Message: error.Error(),
	})
	if conditionChanged {
		spec.Events = append(spec.Events, NewMissingConfigReferencesEvent(error))
	}
}

// DetermineStateDiff creates the StateDiff between the blueprint and the actual state of the ecosystem.
// if sth. is not in the lists of installed things, it is considered not installed.
// installedDogus are a map in the form of simpleDoguName->*DoguInstallation. There should be no nil values.
// The StateDiff is an 'as is' representation, therefore no error is thrown, e.g. if dogu namespaces are different and namespace changes are not allowed.
// If there are not allowed actions should be considered at the start of the execution of the blueprint.
// returns an error if the BlueprintSpec is not in the necessary state to determine the stateDiff.
func (spec *BlueprintSpec) DetermineStateDiff(
	ecosystemState ecosystem.EcosystemState,
	referencedSensitiveConfig map[common.DoguConfigKey]common.SensitiveDoguConfigValue,
	isDebugModeActive bool,
) error {
	doguDiffs := determineDoguDiffs(spec.EffectiveBlueprint.Dogus, ecosystemState.InstalledDogus)
	config := spec.EffectiveBlueprint.Config
	if isDebugModeActive {
		config = removeLogLevelChangesFromConfig(config)
	}
	doguConfigDiffs, sensitiveDoguConfigDiffs, globalConfigDiffs := determineConfigDiffs(
		config,
		ecosystemState.GlobalConfig,
		ecosystemState.ConfigByDogu,
		ecosystemState.SensitiveConfigByDogu,
		referencedSensitiveConfig,
	)

	spec.StateDiff = StateDiff{
		DoguDiffs:                doguDiffs,
		DoguConfigDiffs:          doguConfigDiffs,
		SensitiveDoguConfigDiffs: sensitiveDoguConfigDiffs,
		GlobalConfigDiffs:        globalConfigDiffs,
	}

	spec.resetCompletedConditionAfterStateDiff()
	if spec.StateDiff.DoguDiffs.HasChanges() {
		spec.Events = append(spec.Events, newStateDiffEvent(spec.StateDiff))
	}

	invalidBlueprintError := spec.validateStateDiff()
	if invalidBlueprintError != nil {
		conditionChanged := meta.SetStatusCondition(&spec.Conditions, metav1.Condition{
			Type:    ConditionExecutable,
			Status:  metav1.ConditionFalse,
			Reason:  "ForbiddenOperations",
			Message: invalidBlueprintError.Error(),
		})
		if conditionChanged {
			spec.Events = append(spec.Events, BlueprintSpecInvalidEvent{ValidationError: invalidBlueprintError})
		}
		return invalidBlueprintError
	}

	meta.SetStatusCondition(&spec.Conditions, metav1.Condition{
		Type:   ConditionExecutable,
		Status: metav1.ConditionTrue,
		Reason: "Executable",
	})
	return nil
}

// HandleHealthResult sets the healthCondition accordingly to the healthResult and a possible error.
// if an error is given, the condition will be set to unknown.
// The function returns true if the condition changed, otherwise false.
func (spec *BlueprintSpec) HandleHealthResult(healthResult ecosystem.HealthResult, err error) bool {
	if err != nil {
		conditionChanged := meta.SetStatusCondition(&spec.Conditions, metav1.Condition{
			Type:    ConditionEcosystemHealthy,
			Status:  metav1.ConditionUnknown,
			Reason:  "CannotCheckHealth",
			Message: err.Error(),
		})
		return conditionChanged
	}

	if healthResult.AllHealthy() {
		event := EcosystemHealthyEvent{
			doguHealthIgnored: spec.Config.IgnoreDoguHealth,
		}
		conditionChanged := meta.SetStatusCondition(&spec.Conditions, metav1.Condition{
			Type:    ConditionEcosystemHealthy,
			Status:  metav1.ConditionTrue,
			Reason:  "Healthy",
			Message: event.Message(),
		})
		if conditionChanged {
			spec.Events = append(spec.Events, event)
		}
		return conditionChanged
	}

	event := EcosystemUnhealthyEvent{
		HealthResult: healthResult,
	}
	oldHealthyCondition := meta.FindStatusCondition(spec.Conditions, ConditionEcosystemHealthy)
	// determine here, because the condition is a pointer and will change with the SetStatusCondition call below
	isConditionStatusChanged := oldHealthyCondition == nil || oldHealthyCondition.Status == metav1.ConditionTrue
	conditionChanged := meta.SetStatusCondition(&spec.Conditions, metav1.Condition{
		Type:    ConditionEcosystemHealthy,
		Status:  metav1.ConditionFalse,
		Reason:  "Unhealthy",
		Message: "ecosystem health:\n  " + healthResult.String(),
	})
	// only throw an event the first unhealthy time to avoid having too many events
	if isConditionStatusChanged {
		spec.Events = append(spec.Events, event)
	}
	return conditionChanged
}

// ShouldBeApplied returns true if the blueprint should be applied or an early-exit should happen, e.g. while being stopped.
func (spec *BlueprintSpec) ShouldBeApplied() bool {
	if spec.Config.Stopped {
		return false
	}
	// not true does not equal IsStatusConditionFalse here, because not true includes status "unknown"
	return !meta.IsStatusConditionTrue(spec.Conditions, ConditionCompleted) || spec.StateDiff.HasChanges()
}

func (spec *BlueprintSpec) resetCompletedConditionAfterStateDiff() bool {
	if spec.StateDiff.HasChanges() {
		conditionChanged := meta.SetStatusCondition(&spec.Conditions, metav1.Condition{
			Type:    ConditionCompleted,
			Status:  metav1.ConditionFalse,
			Reason:  "StateDiffHasChanges",
			Message: "Blueprint is being applied.",
		})
		return conditionChanged
	}

	return false
}

func (spec *BlueprintSpec) validateStateDiff() error {
	var invalidBlueprintErrors []error

	for _, diff := range spec.StateDiff.DoguDiffs {
		invalidBlueprintErrors = append(invalidBlueprintErrors, spec.validateDoguDiffActions(diff)...)
	}

	return errors.Join(invalidBlueprintErrors...)
}

func (spec *BlueprintSpec) validateDoguDiffActions(diff DoguDiff) []error {
	return util.Map(diff.NeededActions, func(action Action) error {
		if slices.Contains(notAllowedDoguActions, action) {
			if action == ActionSwitchDoguNamespace && spec.Config.AllowDoguNamespaceSwitch {
				return nil
			}

			return getActionNotAllowedError(action)
		}

		return nil
	})
}

func getActionNotAllowedError(action Action) *InvalidBlueprintError {
	return &InvalidBlueprintError{
		Message: fmt.Sprintf("action %q is not allowed", action),
	}
}

// Complete is used to mark the blueprint as completed and to inform the user.
// Returns true if the condition changed, false otherwise.
func (spec *BlueprintSpec) Complete() bool {
	conditionChanged := meta.SetStatusCondition(&spec.Conditions, metav1.Condition{
		Type:   ConditionCompleted,
		Status: metav1.ConditionTrue,
		Reason: "Completed",
	})
	meta.SetStatusCondition(&spec.Conditions, metav1.Condition{
		Type:   ConditionLastApplySucceeded,
		Status: metav1.ConditionTrue,
		Reason: "ApplySucceeded",
	})

	if conditionChanged {
		spec.Events = append(spec.Events, CompletedEvent{})
	}
	return conditionChanged
}

func (spec *BlueprintSpec) SetLastApplySucceededConditionOnError(reason string, err error) bool {
	if err != nil {
		conditionChanged := meta.SetStatusCondition(&spec.Conditions, metav1.Condition{
			Type:    ConditionLastApplySucceeded,
			Status:  metav1.ConditionFalse,
			Reason:  reason,
			Message: err.Error(),
		})
		if conditionChanged {
			spec.Events = append(spec.Events, ExecutionFailedEvent{err: err})
		}

		return conditionChanged
	}

	return false
}

func (spec *BlueprintSpec) MarkEcosystemConfigApplied() {
	spec.Events = append(spec.Events, EcosystemConfigAppliedEvent{})
}

func (spec *BlueprintSpec) MarkBlueprintStopped() {
	spec.Events = append(spec.Events, BlueprintStoppedEvent{})
}

func (spec *BlueprintSpec) MarkDogusApplied(isDogusApplied bool, err error) bool {
	if isDogusApplied {
		spec.Events = append(spec.Events, DogusAppliedEvent{Diffs: spec.StateDiff.DoguDiffs})
	}
	conditionChanged := spec.SetLastApplySucceededConditionOnError(ReasonLastApplyErrorAtDogus, err)
	return conditionChanged
}
