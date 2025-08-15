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
	Blueprint          Blueprint
	BlueprintMask      BlueprintMask
	EffectiveBlueprint EffectiveBlueprint
	StateDiff          StateDiff
	Config             BlueprintConfiguration
	Status             StatusPhase
	Conditions         *[]Condition
	// PersistenceContext can hold generic values needed for persistence with repositories, e.g. version counters or transaction contexts.
	// This field has a generic map type as the values within it highly depend on the used type of repository.
	// This field should be ignored in the whole domain.
	PersistenceContext map[string]interface{}
	Events             []Event
}

type Condition = metav1.Condition

const (
	ConditionValid                = "Valid"
	ConditionExecutable           = "Executable"
	ConditionEcosystemHealthy     = "EcosystemHealthy"
	ConditionSelfUpgradeCompleted = "SelfUpgradeCompleted"
	ConditionConfigApplied        = "ConfigApplied"
	ConditionDogusApplied         = "DogusApplied"
	ConditionComponentsApplied    = "ComponentsApplied"
	ConditionBlueprintApplied     = "BlueprintApplied"
)

type StatusPhase string

const (
	// StatusPhaseNew marks a newly created blueprint-CR.
	StatusPhaseNew StatusPhase = ""
	// StatusPhaseInProgress marks that the blueprint is currently being processed.
	StatusPhaseInProgress StatusPhase = "inProgress"
	// StatusPhaseBlueprintApplicationFailed shows that the blueprint application failed.
	StatusPhaseBlueprintApplicationFailed StatusPhase = "blueprintApplicationFailed"
	// StatusPhaseBlueprintApplied indicates that the blueprint was applied but the ecosystem is not healthy yet.
	StatusPhaseBlueprintApplied StatusPhase = "blueprintApplied"
	// StatusPhaseFailed marks that an error occurred during processing of the blueprint.
	StatusPhaseFailed StatusPhase = "failed"
	// StatusPhaseCompleted marks the blueprint as successfully applied.
	StatusPhaseCompleted StatusPhase = "completed"
	// StatusPhaseRestartsTriggered indicates that a restart has been triggered for all Dogus that needed a restart.
	// Restarts are needed when the Dogu config changes.
	StatusPhaseRestartsTriggered StatusPhase = "restartsTriggered"
)

type BlueprintConfiguration struct {
	// IgnoreDoguHealth forces blueprint upgrades even if dogus are unhealthy
	IgnoreDoguHealth bool
	// IgnoreComponentHealth forces blueprint upgrades even if components are unhealthy
	IgnoreComponentHealth bool
	// AllowDoguNamespaceSwitch allows the blueprint upgrade to switch a dogus namespace
	AllowDoguNamespaceSwitch bool
	// DryRun lets the user test a blueprint run to check if all attributes of the blueprint are correct and avoid a result with a failure state.
	DryRun bool
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
		meta.SetStatusCondition(spec.Conditions, metav1.Condition{
			Type:    ConditionValid,
			Status:  metav1.ConditionFalse,
			Reason:  "blueprint invalid",
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
		if doguMask.TargetState == TargetStatePresent && dogu.TargetState == TargetStateAbsent {
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
		conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
			Type:    ConditionValid,
			Status:  metav1.ConditionFalse,
			Reason:  "inconsistent blueprint",
			Message: err.Error(),
		})
		if conditionChanged {
			spec.Events = append(spec.Events, BlueprintSpecInvalidEvent{ValidationError: err})
		}
	} else {
		meta.SetStatusCondition(spec.Conditions, metav1.Condition{
			Type:   ConditionValid,
			Status: metav1.ConditionTrue,
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
		Dogus:      effectiveDogus,
		Components: spec.Blueprint.Components,
		Config:     effectiveConfig,
	}
	validationError := spec.EffectiveBlueprint.validateOnlyConfigForDogusInBlueprint()
	if validationError != nil {
		conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
			Type:    ConditionValid,
			Status:  metav1.ConditionFalse,
			Reason:  "inconsistent blueprint",
			Message: validationError.Error(),
		})
		if conditionChanged {
			spec.Events = append(spec.Events, BlueprintSpecInvalidEvent{ValidationError: validationError})
		}
		return validationError
	}
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
		Name:               dogu.Name,
		Version:            dogu.Version,
		TargetState:        dogu.TargetState,
		MinVolumeSize:      dogu.MinVolumeSize,
		ReverseProxyConfig: dogu.ReverseProxyConfig,
		AdditionalMounts:   dogu.AdditionalMounts,
	}
	maskDogu, noMaskDoguErr := spec.BlueprintMask.FindDoguByName(dogu.Name.SimpleName)
	if noMaskDoguErr == nil {
		emptyVersion := core.Version{}
		if maskDogu.Version != emptyVersion {
			effectiveDogu.Version = maskDogu.Version
		}
		if maskDogu.Name.Namespace != dogu.Name.Namespace {
			if spec.Config.AllowDoguNamespaceSwitch {
				effectiveDogu.Name.Namespace = maskDogu.Name.Namespace
			} else {
				return Dogu{}, fmt.Errorf(
					"changing the dogu namespace is forbidden by default and can be allowed by a flag: %q -> %q", dogu.Name, maskDogu.Name)
			}
		}
		effectiveDogu.TargetState = maskDogu.TargetState
	}

	return effectiveDogu, nil
}

// It is not allowed to have config without the corresponding dogu, so this will clean up the unnecessary config.
func (spec *BlueprintSpec) removeConfigForMaskedDogus() Config {
	effectiveDoguConfig := maps.Clone(spec.Blueprint.Config.Dogus)

	for _, dogu := range spec.BlueprintMask.Dogus {
		if dogu.TargetState == TargetStateAbsent {
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
	spec.Events = append(spec.Events, NewMissingConfigReferencesEvent(error))
}

// DetermineStateDiff creates the StateDiff between the blueprint and the actual state of the ecosystem.
// if sth. is not in the lists of installed things, it is considered not installed.
// installedDogus are a map in the form of simpleDoguName->*DoguInstallation. There should be no nil values.
// The StateDiff is an 'as is' representation, therefore no error is thrown, e.g. if dogu namespaces are different and namespace changes are not allowed.
// If there are not allowed actions should be considered at the start of the execution of the blueprint.
// returns an error if the BlueprintSpec is not in the necessary state to determine the stateDiff.
func (spec *BlueprintSpec) DetermineStateDiff(
	ecosystemState ecosystem.EcosystemState,
	referencedSensitiveConfig map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue,
) error {
	doguDiffs := determineDoguDiffs(spec.EffectiveBlueprint.Dogus, ecosystemState.InstalledDogus)
	compDiffs, err := determineComponentDiffs(spec.EffectiveBlueprint.Components, ecosystemState.InstalledComponents)
	if err != nil {
		// FIXME: a proper state and event should be set, so that this error don't lead to an endless retry.
		// The error here occurs, if a targetState is not properly set in components. We can remove this case
		// when we introduce the absent flag in the domain or we just ignore this error like for dogu targetState
		return err
	}
	doguConfigDiffs, sensitiveDoguConfigDiffs, globalConfigDiffs := determineConfigDiffs(
		spec.EffectiveBlueprint.Config,
		ecosystemState.GlobalConfig,
		ecosystemState.ConfigByDogu,
		ecosystemState.SensitiveConfigByDogu,
		referencedSensitiveConfig,
	)

	spec.StateDiff = StateDiff{
		DoguDiffs:                doguDiffs,
		ComponentDiffs:           compDiffs,
		DoguConfigDiffs:          doguConfigDiffs,
		SensitiveDoguConfigDiffs: sensitiveDoguConfigDiffs,
		GlobalConfigDiffs:        globalConfigDiffs,
	}

	spec.Events = append(spec.Events, newStateDiffDoguEvent(spec.StateDiff.DoguDiffs))
	spec.Events = append(spec.Events, newStateDiffComponentEvent(spec.StateDiff.ComponentDiffs))
	spec.Events = append(spec.Events, GlobalConfigDiffDeterminedEvent{GlobalConfigDiffs: spec.StateDiff.GlobalConfigDiffs})
	spec.Events = append(spec.Events, NewDoguConfigDiffDeterminedEvent(spec.StateDiff.DoguConfigDiffs))
	spec.Events = append(spec.Events, NewSensitiveDoguConfigDiffDeterminedEvent(spec.StateDiff.SensitiveDoguConfigDiffs))

	invalidBlueprintError := spec.validateStateDiff()
	if invalidBlueprintError != nil {
		conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
			Type:    ConditionExecutable,
			Status:  metav1.ConditionFalse,
			Reason:  "forbidden operations needed",
			Message: invalidBlueprintError.Error(),
		})
		if conditionChanged {
			spec.Events = append(spec.Events, BlueprintSpecInvalidEvent{ValidationError: invalidBlueprintError})
		}
		return invalidBlueprintError
	}

	meta.SetStatusCondition(spec.Conditions, metav1.Condition{
		Type:   ConditionExecutable,
		Status: metav1.ConditionTrue,
	})
	//TODO: we cannot just deduplicate the events here by detecting a condition change,
	// because the blueprint could be executable even after a change of the blueprint.
	// Therefore, a check with "conditionChanged" is not enough to prevent, that we regenerate all events on every reconcile.

	//TODO: We could set all diff-related conditions here, so that don't need to call every step even if the state diff says "no change"
	return nil
}

// CheckEcosystemHealthUpfront checks if the ecosystem is healthy with the given health result and sets the next status phase depending on that.
func (spec *BlueprintSpec) CheckEcosystemHealthUpfront(healthResult ecosystem.HealthResult) error {
	// healthResult does not contain dogu info if IgnoreDoguHealth flag is set. (no need to load all doguInstallations then)
	// Therefore we don't need to exclude dogus while checking with AllHealthy()
	if healthResult.AllHealthy() {
		event := EcosystemHealthyUpfrontEvent{
			doguHealthIgnored:      spec.Config.IgnoreDoguHealth,
			componentHealthIgnored: spec.Config.IgnoreComponentHealth,
		}
		conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
			Type:    ConditionEcosystemHealthy,
			Status:  metav1.ConditionTrue,
			Message: event.Message(),
		})
		if conditionChanged {
			spec.Events = append(spec.Events, event)
		}
		return nil
	} else {
		event := EcosystemUnhealthyUpfrontEvent{HealthResult: healthResult}
		conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
			Type:    ConditionEcosystemHealthy,
			Status:  metav1.ConditionFalse,
			Message: event.Message(),
		})
		if conditionChanged {
			spec.Events = append(spec.Events, event)
		}
		return NewUnhealthyEcosystemError(nil, "ecosystem is unhealthy before applying the blueprint", healthResult)
	}
}

// ShouldBeApplied returns true if the blueprint should be applied or an early-exit should happen, e.g. while dry run.
func (spec *BlueprintSpec) ShouldBeApplied() bool {
	// wrote it in the long form to reduce complexity
	if spec.Config.DryRun {
		return false
	}
	return spec.StateDiff.HasChanges()
}

func (spec *BlueprintSpec) MarkWaitingForSelfUpgrade() {
	conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
		Type:   ConditionSelfUpgradeCompleted,
		Status: metav1.ConditionFalse,
		Reason: "await self upgrade",
	})
	if conditionChanged {
		spec.Events = append(spec.Events, AwaitSelfUpgradeEvent{})
	}
}

func (spec *BlueprintSpec) MarkSelfUpgradeCompleted() {
	conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
		Type:   ConditionSelfUpgradeCompleted,
		Status: metav1.ConditionTrue,
	})
	if conditionChanged {
		spec.Events = append(spec.Events, SelfUpgradeCompletedEvent{})
	}
}

// CheckEcosystemHealthAfterwards checks with the given health result if the ecosystem is healthy and the blueprint was therefore successful.
func (spec *BlueprintSpec) CheckEcosystemHealthAfterwards(healthResult ecosystem.HealthResult) error {
	if healthResult.AllHealthy() {
		event := EcosystemHealthyAfterwardsEvent{}
		conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
			Type:    ConditionEcosystemHealthy,
			Status:  metav1.ConditionTrue,
			Message: event.Message(),
		})
		if conditionChanged {
			spec.Events = append(spec.Events, event)
		}

		return nil
	} else {
		event := EcosystemUnhealthyAfterwardsEvent{HealthResult: healthResult}
		conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
			Type:    ConditionEcosystemHealthy,
			Status:  metav1.ConditionFalse,
			Message: event.Message(),
		})
		if conditionChanged {
			spec.Events = append(spec.Events, event)
		}
		return NewUnhealthyEcosystemError(nil, "ecosystem is unhealthy after applying the blueprint", healthResult)
	}
}

// StartApplying marks the blueprint as in progress, which indicates, that the system started applying the blueprint.
// This state is used to detect complete failures as this state will only stay persisted if the process failed before setting the state to blueprint applied.
func (spec *BlueprintSpec) StartApplying() {
	spec.Status = StatusPhaseInProgress
	spec.Events = append(spec.Events, InProgressEvent{})
}

// MarkBlueprintApplicationFailed sets the blueprint state to application failed, which indicates that the blueprint could not be applied completely.
// In reaction to this, further post-processing will happen.
func (spec *BlueprintSpec) MarkBlueprintApplicationFailed(err error) {
	spec.Status = StatusPhaseBlueprintApplicationFailed
	spec.Events = append(spec.Events, ExecutionFailedEvent{err: err})
}

// MarkBlueprintApplied sets the blueprint state to blueprint applied, which indicates that the blueprint was applied successful and further steps can happen then.
func (spec *BlueprintSpec) MarkBlueprintApplied() {
	spec.Status = StatusPhaseBlueprintApplied
	spec.Events = append(spec.Events, BlueprintAppliedEvent{})
}

// CompletePostProcessing is used to mark the blueprint as completed or failed , depending on the blueprint application result.
func (spec *BlueprintSpec) CompletePostProcessing() {
	// this function will not be called, if the ecosystem is not healthy
	switch spec.Status {
	case StatusPhaseInProgress:
		spec.Status = StatusPhaseFailed
		err := errors.New(handleInProgressMsg)
		spec.Events = append(spec.Events, ExecutionFailedEvent{err: err})

	case StatusPhaseBlueprintApplicationFailed:
		spec.Status = StatusPhaseFailed
		spec.Events = append(spec.Events, ExecutionFailedEvent{err: errors.New("could not apply blueprint")})
	default:
		//if healthy
		spec.Status = StatusPhaseCompleted
		spec.Events = append(spec.Events, CompletedEvent{})
	}

}

var notAllowedComponentActions = []Action{ActionDowngrade, ActionSwitchComponentNamespace}

// ActionSwitchDoguNamespace is an exception and should be handled with the blueprint config.
var notAllowedDoguActions = []Action{ActionDowngrade, ActionSwitchDoguNamespace}

func (spec *BlueprintSpec) validateStateDiff() error {
	var invalidBlueprintErrors []error

	for _, diff := range spec.StateDiff.DoguDiffs {
		invalidBlueprintErrors = append(invalidBlueprintErrors, spec.validateDoguDiffActions(diff)...)
	}

	for _, diff := range spec.StateDiff.ComponentDiffs {
		invalidBlueprintErrors = append(invalidBlueprintErrors, spec.validateComponentDiffActions(diff)...)
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

func (spec *BlueprintSpec) validateComponentDiffActions(diff ComponentDiff) []error {
	return util.Map(diff.NeededActions, func(action Action) error {
		if slices.Contains(notAllowedComponentActions, action) {
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

func (spec *BlueprintSpec) GetDogusThatNeedARestart() []cescommons.SimpleName {
	var dogusThatNeedRestart []cescommons.SimpleName
	dogusInEffectiveBlueprint := spec.EffectiveBlueprint.Dogus
	for _, dogu := range dogusInEffectiveBlueprint {
		//TODO: test this
		if spec.StateDiff.DoguConfigDiffs[dogu.Name.SimpleName].HasChanges() ||
			spec.StateDiff.SensitiveDoguConfigDiffs[dogu.Name.SimpleName].HasChanges() {
			dogusThatNeedRestart = append(dogusThatNeedRestart, dogu.Name.SimpleName)
		}
	}
	return dogusThatNeedRestart
}

func (spec *BlueprintSpec) StartApplyEcosystemConfig() {
	event := ApplyEcosystemConfigEvent{}
	conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
		Type:    ConditionConfigApplied,
		Status:  metav1.ConditionFalse,
		Message: event.Message(),
	})
	if conditionChanged {
		spec.Events = append(spec.Events, ApplyEcosystemConfigEvent{})
	}
}

func (spec *BlueprintSpec) MarkApplyEcosystemConfigFailed(err error) {
	conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
		Type:    ConditionConfigApplied,
		Status:  metav1.ConditionFalse,
		Message: err.Error(),
	})
	if conditionChanged {
		spec.Events = append(spec.Events, ApplyEcosystemConfigFailedEvent{err: err})
	}
}

func (spec *BlueprintSpec) MarkEcosystemConfigApplied() {
	conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
		Type:   ConditionConfigApplied,
		Status: metav1.ConditionTrue,
	})
	if conditionChanged {
		spec.Events = append(spec.Events, EcosystemConfigAppliedEvent{})
	}
}

const handleInProgressMsg = "cannot handle blueprint in state " + string(StatusPhaseInProgress) +
	" as this state shows that the appliance of the blueprint was interrupted before it could update the state " +
	"to either " + string(StatusPhaseFailed) + " or " + string(StatusPhaseCompleted)
