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
	//FIXME: check if we can use a non-pointer type. The only reason for a pointer was the meta.SetStatusCondition function
	// The pointer is really ugly, because we need to explicitly set conditions in every test
	Conditions *[]Condition
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
	ConditionComponentsApplied    = "ComponentsApplied"
	ConditionDogusApplied         = "DogusApplied"
	ConditionCompleted            = "Completed"

	ReasonCannotApply = "CannotApply"
)

var (
	notAllowedComponentActions = []Action{ActionDowngrade, ActionSwitchComponentNamespace}
	// ActionSwitchDoguNamespace is an exception and should be handled with the blueprint config.
	notAllowedDoguActions = []Action{ActionDowngrade, ActionSwitchDoguNamespace}
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
// changed a domain.InvalidBlueprintError if blueprint is invalid
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
			Reason:  "Inconsistent",
			Message: err.Error(),
		})
		if conditionChanged {
			spec.Events = append(spec.Events, BlueprintSpecInvalidEvent{ValidationError: err})
		}
	} else {
		meta.SetStatusCondition(spec.Conditions, metav1.Condition{
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
		Dogus:      effectiveDogus,
		Components: spec.Blueprint.Components,
		Config:     effectiveConfig,
	}
	validationError := spec.EffectiveBlueprint.validateOnlyConfigForDogusInBlueprint()
	if validationError != nil {
		conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
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
// changed an error if the BlueprintSpec is not in the necessary state to determine the stateDiff.
func (spec *BlueprintSpec) DetermineStateDiff(
	ecosystemState ecosystem.EcosystemState,
	referencedSensitiveConfig map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue,
) error {
	doguDiffs := determineDoguDiffs(spec.EffectiveBlueprint.Dogus, ecosystemState.InstalledDogus)
	compDiffs, err := determineComponentDiffs(spec.EffectiveBlueprint.Components, ecosystemState.InstalledComponents)
	if err != nil {
		// FIXME: a proper state and event should be set, so that this error doesn't lead to an endless retry.
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

	//TODO: we need the possible error from the use case to set the condition to Unknown
	spec.setDogusAppliedConditionAfterStateDiff(nil)
	spec.Events = append(spec.Events, newStateDiffComponentEvent(spec.StateDiff.ComponentDiffs))
	spec.Events = append(spec.Events, GlobalConfigDiffDeterminedEvent{GlobalConfigDiffs: spec.StateDiff.GlobalConfigDiffs})
	spec.Events = append(spec.Events, NewDoguConfigDiffDeterminedEvent(spec.StateDiff.DoguConfigDiffs))
	spec.Events = append(spec.Events, NewSensitiveDoguConfigDiffDeterminedEvent(spec.StateDiff.SensitiveDoguConfigDiffs))

	invalidBlueprintError := spec.validateStateDiff()
	if invalidBlueprintError != nil {
		conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
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

	meta.SetStatusCondition(spec.Conditions, metav1.Condition{
		Type:   ConditionExecutable,
		Status: metav1.ConditionTrue,
		Reason: "Executable",
	})
	//TODO: we cannot just deduplicate the events here by detecting a condition change,
	// because the blueprint could be executable even after a change of the blueprint.
	// Therefore, a check with "conditionChanged" is not enough to prevent, that we regenerate all events on every reconcile.

	//TODO: We could set all diff-related conditions here
	// Con: this could override Reason and Message.
	//		They will be set again when we try to apply.
	//		We also could check, if there is a reason and a message and just change nothing then if there is still a diff.
	// Pro: we don't need to update the blueprint so often -> big update here and only setting conditionTrue while applying things
	// Con: StateDiff as it is could problematic because it hides conflict errors
	//		(is this a real problem if we set it anyways in the next run?).
	// 		The idea is, that we just check the needed cluster-state when we try to apply.
	//		Then there is no central stateDiff anymore, just some independent ApplyUseCases.
	//      I am not sure with this yet.
	// Con: a central stateDiff could be problematic because we then check every state and
	//		we cannot optimize it if we know, which watch triggered the reconciliation.

	//TODO: The state diff will be generated at every run and will be written on the blueprint-CR.Status.
	//		After every apply, the state diff will be empty and therefore will be deleted on the blueprint-CR.
	//		It is only useful for dry-run or error states, where no further work happens.
	//		Maybe we just not write it in the status anymore? How to debug then?
	//		The old classic blueprint-process was hard to understand because there was no overview.
	//		A separate BlueprintExecution-CR could be a solution but i am not sure, if we have enough time for that and if it is the right choice.
	//			CR could have the diff as it's spec.
	//			It is a single execution like k8s-jobs.
	//			The operator will always spawn more blueprintExecutions if they fail or there is a diff left after applying.
	return nil
}

// HandleHealthResult sets the healthCondition accordingly to the healthResult and a possible error.
// if an error is given, the condition will be set to unknown.
// The function changed true if the condition changed, otherwise false.
func (spec *BlueprintSpec) HandleHealthResult(healthResult ecosystem.HealthResult, err error) bool {
	if err != nil {
		conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
			Type:    ConditionEcosystemHealthy,
			Status:  metav1.ConditionUnknown,
			Reason:  "CannotCheckHealth",
			Message: err.Error(),
		})
		return conditionChanged
	}

	if healthResult.AllHealthy() {
		event := EcosystemHealthyEvent{
			doguHealthIgnored:      spec.Config.IgnoreDoguHealth,
			componentHealthIgnored: spec.Config.IgnoreComponentHealth,
		}
		conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
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
	conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
		Type:    ConditionEcosystemHealthy,
		Status:  metav1.ConditionFalse,
		Reason:  "Unhealthy",
		Message: event.Message(),
	})
	if conditionChanged {
		spec.Events = append(spec.Events, event)
	}
	return conditionChanged
}

// ShouldBeApplied changed true if the blueprint should be applied or an early-exit should happen, e.g. while dry run.
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
		Reason: "AwaitSelfUpgrade",
	})
	if conditionChanged {
		spec.Events = append(spec.Events, AwaitSelfUpgradeEvent{})
	}
}

func (spec *BlueprintSpec) MarkSelfUpgradeCompleted() {
	conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
		Type:   ConditionSelfUpgradeCompleted,
		Status: metav1.ConditionTrue,
		Reason: "Completed",
	})
	if conditionChanged {
		spec.Events = append(spec.Events, SelfUpgradeCompletedEvent{})
	}
}

// setDogusAppliedConditionAfterStateDiff sets the ConditionDogusApplied based on the diff, and it's current state.
// This function gets called after creating the stateDiff. The ConditionDogusApplied could be applied in a previous run
// either by the state diff or later while applying the dogus.
// If there is a condition from the DoguApply-step, then do not override it, unless it is obviously outdated.
//
//	decision table:
//	current condition   withDiff     withoutDiff  DiffError
//	=======================================================
//	none                NeedToApply	 Applied      Unknown
//	Unknown             NeedToApply	 Applied      Unknown
//	NeedToApply         NeedToApply	 Applied      Unknown
//	Applied             NeedToApply	 no change    Unknown
//	CannotApply         no change    no change    Unknown
func (spec *BlueprintSpec) setDogusAppliedConditionAfterStateDiff(diffErr error) bool {
	condition := meta.FindStatusCondition(*spec.Conditions, ConditionDogusApplied)

	//	current condition DiffError
	//	===========================
	//	none              Unknown
	//	Unknown           Unknown
	//	NeedToApply       Unknown
	//	Applied           Unknown
	//	CannotApply       Unknown
	if diffErr != nil {
		conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
			Type:    ConditionDogusApplied,
			Status:  metav1.ConditionUnknown,
			Reason:  "CannotDetermineStateDiff",
			Message: diffErr.Error(),
		})
		return conditionChanged
	}

	//	current condition   withDiff     withoutDiff
	//	============================================
	//	none                NeedToApply	 Applied
	//	Unknown             NeedToApply	 Applied
	//	NeedToApply         NeedToApply	 Applied
	if condition == nil || condition.Status == metav1.ConditionUnknown || condition.Reason == "NeedToApply" {
		if spec.StateDiff.DoguDiffs.HasChanges() {
			return spec.setDogusNeedToApply()
		} else {
			event := newStateDiffDoguEvent(spec.StateDiff.DoguDiffs)
			conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
				Type:    ConditionDogusApplied,
				Status:  metav1.ConditionFalse,
				Reason:  "Applied",
				Message: event.Message(),
			})
			if conditionChanged {
				spec.Events = append(spec.Events, event)
			}
			return conditionChanged
		}
	}
	//	current condition   withDiff     withoutDiff
	//	============================================
	//	CannotApply         no change    no change
	if condition.Reason == ReasonCannotApply {
		return false
	}

	//	current condition   withDiff     withoutDiff
	//	============================================
	//	Applied             NeedToApply	 no change
	if condition.Reason == "Applied" {
		if spec.StateDiff.DoguDiffs.HasChanges() {
			return spec.setDogusNeedToApply()
		} else {
			return false
		}
	}
	// should never happen
	return false
}

func (spec *BlueprintSpec) setDogusNeedToApply() bool {
	event := newStateDiffDoguEvent(spec.StateDiff.DoguDiffs)
	conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
		Type:    ConditionDogusApplied,
		Status:  metav1.ConditionFalse,
		Reason:  "NeedToApply",
		Message: event.Message(),
	})
	if conditionChanged {
		spec.Events = append(spec.Events, event)
	}
	return conditionChanged
}

// SetComponentsAppliedCondition informs the user about the state of the component apply.
// If an error is given, it will set the condition to failed accordingly, otherwise it marks it as a success.
// Returns true if the condition changed, otherwise false.
func (spec *BlueprintSpec) SetComponentsAppliedCondition(err error) bool {
	if err != nil {
		conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
			Type:    ConditionComponentsApplied,
			Status:  metav1.ConditionFalse,
			Reason:  ReasonCannotApply,
			Message: err.Error(),
		})
		if conditionChanged {
			spec.Events = append(spec.Events, ExecutionFailedEvent{err: err})
		}
		return conditionChanged
	}
	event := ComponentsAppliedEvent{Diffs: spec.StateDiff.ComponentDiffs}
	conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
		Type:    ConditionComponentsApplied,
		Status:  metav1.ConditionTrue,
		Reason:  "Applied",
		Message: event.Message(),
	})
	if conditionChanged && spec.StateDiff.ComponentDiffs.HasChanges() {
		spec.Events = append(spec.Events, event)
	}
	return conditionChanged
}

// SetDogusAppliedCondition informs the user about the state of the dogu apply.
// If an error is given, it will set the condition to failed accordingly, otherwise it marks it as a success.
// Returns true if the condition changed, otherwise false.
func (spec *BlueprintSpec) SetDogusAppliedCondition(err error) bool {
	if err != nil {
		conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
			Type:    ConditionDogusApplied,
			Status:  metav1.ConditionFalse,
			Reason:  ReasonCannotApply,
			Message: err.Error(),
		})
		if conditionChanged {
			spec.Events = append(spec.Events, ExecutionFailedEvent{err: err})
		}
		return conditionChanged
	}
	event := DogusAppliedEvent{Diffs: spec.StateDiff.DoguDiffs}
	conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
		Type:    ConditionDogusApplied,
		Status:  metav1.ConditionTrue,
		Reason:  "Applied",
		Message: event.Message(),
	})
	if conditionChanged && spec.StateDiff.DoguDiffs.HasChanges() {
		spec.Events = append(spec.Events, event)
	}
	return conditionChanged
}

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

func (spec *BlueprintSpec) StartApplyEcosystemConfig() {
	event := ApplyEcosystemConfigEvent{}
	conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
		Type:    ConditionConfigApplied,
		Status:  metav1.ConditionFalse,
		Reason:  "Applying",
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
		Reason:  "ApplyingFailed",
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
		Reason: "Applied",
	})
	if conditionChanged {
		spec.Events = append(spec.Events, EcosystemConfigAppliedEvent{})
	}
}

// Complete is used to mark the blueprint as completed and to inform the user.
// Returns true if anything changed, false otherwise.
func (spec *BlueprintSpec) Complete() bool {
	conditionChanged := meta.SetStatusCondition(spec.Conditions, metav1.Condition{
		Type:   ConditionCompleted,
		Status: metav1.ConditionTrue,
		Reason: "Completed",
	})
	if conditionChanged {
		spec.Events = append(spec.Events, CompletedEvent{})
	}
	return conditionChanged
}
