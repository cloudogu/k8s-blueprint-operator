package domain

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
)

type BlueprintSpec struct {
	Id                 string
	Blueprint          Blueprint
	BlueprintMask      BlueprintMask
	EffectiveBlueprint EffectiveBlueprint
	StateDiff          StateDiff
	Config             BlueprintConfiguration
	Status             StatusPhase
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
	// StatusPhaseInvalid marks the given blueprint spec is semantically incorrect.
	StatusPhaseInvalid StatusPhase = "invalid"
	// StatusPhaseEcosystemHealthyUpfront marks that all currently installed dogus are healthy.
	StatusPhaseEcosystemHealthyUpfront StatusPhase = "ecosystemHealthyUpfront"
	// StatusPhaseEcosystemUnhealthyUpfront marks that some currently installed dogus are unhealthy.
	StatusPhaseEcosystemUnhealthyUpfront StatusPhase = "ecosystemUnhealthyUpfront"
	// StatusPhaseBlueprintApplicationPreProcessed shows that all pre-processing steps for the blueprint application
	// were successful.
	StatusPhaseBlueprintApplicationPreProcessed StatusPhase = "blueprintApplicationPreProcessed"
	// StatusPhaseAwaitSelfUpgrade marks that the blueprint operator waits for termination for a self upgrade.
	StatusPhaseAwaitSelfUpgrade StatusPhase = "awaitSelfUpgrade"
	// StatusPhaseSelfUpgradeCompleted marks that the blueprint operator itself got successfully upgraded.
	StatusPhaseSelfUpgradeCompleted StatusPhase = "selfUpgradeCompleted"
	// StatusPhaseInProgress marks that the blueprint is currently being processed.
	StatusPhaseInProgress StatusPhase = "inProgress"
	// StatusPhaseBlueprintApplicationFailed shows that the blueprint application failed.
	StatusPhaseBlueprintApplicationFailed StatusPhase = "blueprintApplicationFailed"
	// StatusPhaseBlueprintApplied indicates that the blueprint was applied but the ecosystem is not healthy yet.
	StatusPhaseBlueprintApplied StatusPhase = "blueprintApplied"
	// StatusPhaseEcosystemHealthyAfterwards shows that the ecosystem got healthy again after applying the blueprint.
	StatusPhaseEcosystemHealthyAfterwards StatusPhase = "ecosystemHealthyAfterwards"
	// StatusPhaseEcosystemUnhealthyAfterwards shows that the ecosystem got not healthy again after applying the blueprint.
	StatusPhaseEcosystemUnhealthyAfterwards StatusPhase = "ecosystemUnhealthyAfterwards"
	// StatusPhaseFailed marks that an error occurred during processing of the blueprint.
	StatusPhaseFailed StatusPhase = "failed"
	// StatusPhaseCompleted marks the blueprint as successfully applied.
	StatusPhaseCompleted StatusPhase = "completed"
	// StatusPhaseApplyRegistryConfig indicates that the apply registry config phase is active.
	StatusPhaseApplyRegistryConfig StatusPhase = "applyRegistryConfig"
	// StatusPhaseApplyRegistryConfigFailed indicates that the phase to apply registry config phase failed.
	StatusPhaseApplyRegistryConfigFailed StatusPhase = "applyRegistryConfigFailed"
	// StatusPhaseRegistryConfigApplied indicates that the phase to apply registry config phase succeeded.
	StatusPhaseRegistryConfigApplied StatusPhase = "registryConfigApplied"
)

// censorValue is the value for censoring sensitive blueprint configuration data.
const censorValue = "*****"

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
	switch spec.Status {
	case StatusPhaseNew: // continue
	case StatusPhaseInvalid: // do not validate again
		return &InvalidBlueprintError{Message: "blueprint spec was marked invalid before: do not revalidate"}
	default: // do not validate again. for all other status it must be either status validated or a status beyond that
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
		dogu, noDoguFoundError := FindDoguByName(spec.Blueprint.Dogus, doguMask.Name.SimpleName)
		if noDoguFoundError != nil {
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
		return nil // do not regenerate effective blueprint
	case StatusPhaseNew: // stop
		return fmt.Errorf("cannot calculate effective blueprint before the blueprint spec is validated")
	case StatusPhaseInvalid: // stop
		return fmt.Errorf("cannot calculate effective blueprint on invalid blueprint spec")
	default: // continue: StatusPhaseValidated, StatusPhaseInProgress, StatusPhaseFailed, StatusPhaseCompleted
	}

	effectiveDogus, err := spec.calculateEffectiveDogus()
	if err != nil {
		return err
	}

	spec.EffectiveBlueprint = EffectiveBlueprint{
		Dogus:      effectiveDogus,
		Components: spec.Blueprint.Components,
		Config:     spec.Blueprint.Config,
	}
	validationError := spec.EffectiveBlueprint.validateOnlyConfigForDogusInBlueprint()
	if validationError != nil {
		spec.Status = StatusPhaseInvalid
		spec.Events = append(spec.Events, BlueprintSpecInvalidEvent{ValidationError: validationError})
		return validationError
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
		Name:        dogu.Name,
		Version:     dogu.Version,
		TargetState: dogu.TargetState,
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
func (spec *BlueprintSpec) DetermineStateDiff(
	ecosystemState ecosystem.EcosystemState,
) error {
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

	doguDiffs := determineDoguDiffs(spec.EffectiveBlueprint.Dogus, ecosystemState.InstalledDogus)
	compDiffs, err := determineComponentDiffs(spec.EffectiveBlueprint.Components, ecosystemState.InstalledComponents)
	if err != nil {
		//FIXME: a proper state and event should be set, so that this error don't lead to an endless retry.
		// we need to analyze first, what kind of error this is. Why do we need one?
		return err
	}
	doguConfigDiffs, globalConfigDiffs := determineConfigDiffs(
		spec.EffectiveBlueprint.Config,
		ecosystemState,
	)

	spec.StateDiff = StateDiff{
		DoguDiffs:         doguDiffs,
		ComponentDiffs:    compDiffs,
		DoguConfigDiffs:   doguConfigDiffs,
		GlobalConfigDiffs: globalConfigDiffs,
	}

	spec.Events = append(spec.Events, newStateDiffDoguEvent(spec.StateDiff.DoguDiffs))
	spec.Events = append(spec.Events, newStateDiffComponentEvent(spec.StateDiff.ComponentDiffs))
	spec.Events = append(spec.Events, GlobalConfigDiffDeterminedEvent{GlobalConfigDiffs: spec.StateDiff.GlobalConfigDiffs})
	spec.Events = append(spec.Events, DoguConfigDiffDeterminedEvent{spec.StateDiff.DoguConfigDiffs})

	invalidBlueprintError := spec.validateStateDiff()
	if invalidBlueprintError != nil {
		spec.Status = StatusPhaseInvalid
		spec.Events = append(spec.Events, BlueprintSpecInvalidEvent{ValidationError: invalidBlueprintError})
		return invalidBlueprintError
	}

	spec.Status = StatusPhaseStateDiffDetermined

	return nil
}

// CheckEcosystemHealthUpfront checks if the ecosystem is healthy with the given health result and sets the next status phase depending on that.
func (spec *BlueprintSpec) CheckEcosystemHealthUpfront(healthResult ecosystem.HealthResult) {
	// healthResult does not contain dogu info if IgnoreDoguHealth flag is set. (no need to load all doguInstallations then)
	// Therefore we don't need to exclude dogus while checking with AllHealthy()
	if healthResult.AllHealthy() {
		spec.Status = StatusPhaseEcosystemHealthyUpfront
		spec.Events = append(spec.Events, EcosystemHealthyUpfrontEvent{doguHealthIgnored: spec.Config.IgnoreDoguHealth,
			componentHealthIgnored: spec.Config.IgnoreComponentHealth})
	} else {
		spec.Status = StatusPhaseEcosystemUnhealthyUpfront
		spec.Events = append(spec.Events, EcosystemUnhealthyUpfrontEvent{HealthResult: healthResult})
	}
}

// ShouldBeApplied returns true if the blueprint should be applied or an early-exit should happen, e.g. while dry run.
func (spec *BlueprintSpec) ShouldBeApplied() bool {
	// TODO: also check if an early-exit is possible if no changes need to be applied, see PR #29
	return !spec.Config.DryRun
}

// CompletePreProcessing decides if the blueprint is ready to be applied or not by setting the fitting next status phase.
func (spec *BlueprintSpec) CompletePreProcessing() {
	if spec.Config.DryRun {
		spec.Events = append(spec.Events, BlueprintDryRunEvent{})
	} else {
		spec.Status = StatusPhaseBlueprintApplicationPreProcessed
		spec.Events = append(spec.Events, BlueprintApplicationPreProcessedEvent{})
	}
}

// HandleSelfUpgrade checks if a self upgrade is needed and sets the appropriate status.
// if the operator is not installed in the usual way, the actualInstalledVersion can be nil.
// Returns the ComponentDiff for the given component name so that it can be used to initiate further steps.
func (spec *BlueprintSpec) HandleSelfUpgrade(ownComponentName common.SimpleComponentName, actualInstalledVersion *semver.Version) ComponentDiff {
	ownDiff := spec.StateDiff.ComponentDiffs.GetComponentDiffByName(ownComponentName)
	// if already everything is as it should
	if isExpectedVersionInstalled(ownDiff.Expected.Version, actualInstalledVersion) {
		// no self upgrade planned
		spec.Status = StatusPhaseSelfUpgradeCompleted
		spec.Events = append(spec.Events, SelfUpgradeCompletedEvent{})
		return ownDiff
	}
	// if there is sth. left to be done
	spec.Status = StatusPhaseAwaitSelfUpgrade
	spec.Events = append(spec.Events, AwaitSelfUpgradeEvent{})
	return ownDiff
}

func isExpectedVersionInstalled(expected, actual *semver.Version) bool {
	if expected == nil {
		return true
	}
	if actual == nil {
		return false
	}
	return expected.Equal(actual)
}

// CheckEcosystemHealthAfterwards checks with the given health result if the ecosystem is healthy and the blueprint was therefore successful.
func (spec *BlueprintSpec) CheckEcosystemHealthAfterwards(healthResult ecosystem.HealthResult) {
	if healthResult.AllHealthy() {
		spec.Status = StatusPhaseEcosystemHealthyAfterwards
		spec.Events = append(spec.Events, EcosystemHealthyAfterwardsEvent{})
	} else {
		spec.Status = StatusPhaseEcosystemUnhealthyAfterwards
		spec.Events = append(spec.Events, EcosystemUnhealthyAfterwardsEvent{HealthResult: healthResult})
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

// CensorSensitiveData censors all sensitive configuration data of the blueprint, effective blueprint and the statediff,
// to make the values unrecognisable.
func (spec *BlueprintSpec) CensorSensitiveData() {
	spec.Blueprint.Config = spec.Blueprint.Config.censorValues()
	spec.EffectiveBlueprint.Config = spec.EffectiveBlueprint.Config.censorValues()
	for k, v := range spec.StateDiff.DoguConfigDiffs {
		spec.StateDiff.DoguConfigDiffs[k] = v.censorValues()
	}

	spec.Events = append(spec.Events, SensitiveConfigDataCensoredEvent{})
}

// CompletePostProcessing is used to mark the blueprint as completed or failed , depending on the blueprint application result.
func (spec *BlueprintSpec) CompletePostProcessing() {
	switch spec.Status {
	case StatusPhaseEcosystemHealthyAfterwards:
		spec.Status = StatusPhaseCompleted
		spec.Events = append(spec.Events, CompletedEvent{})
	case StatusPhaseInProgress:
		spec.Status = StatusPhaseFailed
		err := errors.New(handleInProgressMsg)
		spec.Events = append(spec.Events, ExecutionFailedEvent{err: err})
	case StatusPhaseEcosystemUnhealthyAfterwards:
		spec.Status = StatusPhaseFailed
		spec.Events = append(spec.Events, ExecutionFailedEvent{err: errors.New("ecosystem is unhealthy")})
	case StatusPhaseBlueprintApplicationFailed:
		spec.Status = StatusPhaseFailed
		spec.Events = append(spec.Events, ExecutionFailedEvent{err: errors.New("could not apply blueprint")})
	}
}

var notAllowedComponentActions = []Action{ActionSwitchComponentNamespace}

// ActionSwitchDoguNamespace is an exception and should be handled with the blueprint config.
var notAllowedDoguActions = []Action{ActionDowngrade, ActionSwitchDoguNamespace}

func (spec *BlueprintSpec) validateStateDiff() error {
	dogusByAction := util.GroupBy(spec.StateDiff.DoguDiffs, func(doguDiff DoguDiff) Action {
		return doguDiff.NeededAction
	})
	var invalidBlueprintErrors []error

	for _, action := range notAllowedDoguActions {
		if action == ActionSwitchDoguNamespace && spec.Config.AllowDoguNamespaceSwitch {
			continue
		}
		invalidBlueprintErrors = evaluateInvalidAction(action, dogusByAction, invalidBlueprintErrors)
	}

	componentsByAction := util.GroupBy(spec.StateDiff.ComponentDiffs, func(componentDiff ComponentDiff) Action {
		return componentDiff.NeededAction
	})

	for _, action := range notAllowedComponentActions {
		invalidBlueprintErrors = evaluateInvalidAction(action, componentsByAction, invalidBlueprintErrors)
	}

	return errors.Join(invalidBlueprintErrors...)
}

func evaluateInvalidAction[T any](action Action, mapByAction map[Action][]T, invalidBlueprintErrors []error) []error {
	invalidElement := mapByAction[action]
	if len(invalidElement) != 0 {
		err := &InvalidBlueprintError{
			Message: fmt.Sprintf("action %q is not allowed: %v", action, invalidElement),
		}
		invalidBlueprintErrors = append(invalidBlueprintErrors, err)
	}

	return invalidBlueprintErrors
}

func (spec *BlueprintSpec) StartApplyRegistryConfig() {
	spec.Status = StatusPhaseApplyRegistryConfig
	spec.Events = append(spec.Events, ApplyRegistryConfigEvent{})
}

func (spec *BlueprintSpec) MarkApplyRegistryConfigFailed(err error) {
	spec.Status = StatusPhaseApplyRegistryConfigFailed
	spec.Events = append(spec.Events, ApplyRegistryConfigFailedEvent{err: err})
}

func (spec *BlueprintSpec) MarkRegistryConfigApplied() {
	spec.Status = StatusPhaseRegistryConfigApplied
	spec.Events = append(spec.Events, RegistryConfigAppliedEvent{})
}

const handleInProgressMsg = "cannot handle blueprint in state " + string(StatusPhaseInProgress) +
	" as this state shows that the appliance of the blueprint was interrupted before it could update the state " +
	"to either " + string(StatusPhaseFailed) + " or " + string(StatusPhaseCompleted)
