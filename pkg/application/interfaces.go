package application

import (
	"context"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
)

type blueprintSpecValidationUseCase interface {
	ValidateBlueprintSpecStatically(ctx context.Context, blueprint *domain.BlueprintSpec) error
	ValidateBlueprintSpecDynamically(ctx context.Context, blueprint *domain.BlueprintSpec) error
}

type effectiveBlueprintUseCase interface {
	CalculateEffectiveBlueprint(ctx context.Context, blueprint *domain.BlueprintSpec) error
}

type stateDiffUseCase interface {
	DetermineStateDiff(ctx context.Context, blueprint *domain.BlueprintSpec) error
}

type doguInstallationUseCase interface {
	CheckDoguHealth(ctx context.Context) (ecosystem.DoguHealthResult, error)
	WaitForHealthyDogus(ctx context.Context) (ecosystem.DoguHealthResult, error)
	ApplyDoguStates(ctx context.Context, blueprint *domain.BlueprintSpec) error
}

type doguRestartUseCase interface {
	TriggerDoguRestarts(ctx context.Context, blueprint *domain.BlueprintSpec) error
}

type applyDogusUseCase interface {
	ApplyDogus(ctx context.Context, blueprint *domain.BlueprintSpec) error
}

type applyComponentsUseCase interface {
	ApplyComponents(ctx context.Context, blueprint *domain.BlueprintSpec) error
}

type componentInstallationUseCase interface {
	ApplyComponentStates(ctx context.Context, blueprint *domain.BlueprintSpec) error
	CheckComponentHealth(ctx context.Context) (ecosystem.ComponentHealthResult, error)
	WaitForHealthyComponents(ctx context.Context) (ecosystem.ComponentHealthResult, error)
	applyComponentState(context.Context, domain.ComponentDiff, *ecosystem.ComponentInstallation) error
}

type applyBlueprintSpecUseCase interface {
	PostProcessBlueprintApplication(ctx context.Context, blueprint *domain.BlueprintSpec) error
}

type ecosystemHealthUseCase interface {
	CheckEcosystemHealth(context.Context, *domain.BlueprintSpec) (ecosystem.HealthResult, error)
}

type selfUpgradeUseCase interface {
	HandleSelfUpgrade(ctx context.Context, blueprint *domain.BlueprintSpec) error
}

type ecosystemConfigUseCase interface {
	ApplyConfig(ctx context.Context, blueprint *domain.BlueprintSpec) error
}

type doguInstallationRepository interface {
	domainservice.DoguInstallationRepository
}

//nolint:unused
//goland:noinspection GoUnusedType
type componentInstallationRepository interface {
	domainservice.ComponentInstallationRepository
}

//nolint:unused
//goland:noinspection GoUnusedType
type blueprintSpecRepository interface {
	domainservice.BlueprintSpecRepository
}

type requiredComponentsProvider interface {
	domainservice.RequiredComponentsProvider
}

type healthWaitConfigProvider interface {
	domainservice.HealthWaitConfigProvider
}

type healthConfigProvider interface {
	requiredComponentsProvider
	healthWaitConfigProvider
}

// interface duplication for mocks

//nolint:unused
//goland:noinspection GoUnusedType
type remoteDoguRegistry interface {
	domainservice.RemoteDoguRegistry
}

type globalConfigRepository interface {
	domainservice.GlobalConfigRepository
}

type doguConfigRepository interface {
	domainservice.DoguConfigRepository
}

type sensitiveDoguConfigRepository interface {
	domainservice.SensitiveDoguConfigRepository
}

type sensitiveConfigRefReader interface {
	domainservice.SensitiveConfigRefReader
}

type doguRestartRepository interface {
	domainservice.DoguRestartRepository
}

// validateDependenciesDomainUseCase is an interface for the domain service for better testability
type validateDependenciesDomainUseCase interface {
	ValidateDependenciesForAllDogus(ctx context.Context, effectiveBlueprint domain.EffectiveBlueprint) error
}

// validateAdditionalMountsDomainUseCase is an interface for the domain service for better testability
type validateAdditionalMountsDomainUseCase interface {
	ValidateAdditionalMounts(ctx context.Context, effectiveBlueprint domain.EffectiveBlueprint) error
}
