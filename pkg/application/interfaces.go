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
	ApplyDoguStates(ctx context.Context, blueprintId string) error
}

type doguRestartUseCase interface {
	TriggerDoguRestarts(ctx context.Context, blueprintid string) error
}

type componentInstallationUseCase interface {
	ApplyComponentStates(ctx context.Context, blueprintId string) error
	CheckComponentHealth(ctx context.Context) (ecosystem.ComponentHealthResult, error)
	WaitForHealthyComponents(ctx context.Context) (ecosystem.ComponentHealthResult, error)
	applyComponentState(context.Context, domain.ComponentDiff, *ecosystem.ComponentInstallation) error
}

type applyBlueprintSpecUseCase interface {
	CheckEcosystemHealthUpfront(ctx context.Context, blueprintId string) error
	CheckEcosystemHealthAfterwards(ctx context.Context, blueprintId string) error
	PreProcessBlueprintApplication(ctx context.Context, blueprintId string) error
	PostProcessBlueprintApplication(ctx context.Context, blueprintId string) error
	ApplyBlueprintSpec(ctx context.Context, blueprintId string) error
}

type ecosystemHealthUseCase interface {
	CheckEcosystemHealth(ctx context.Context, ignoreDoguHealth bool, ignoreComponentHealth bool) (ecosystem.HealthResult, error)
	WaitForHealthyEcosystem(ctx context.Context, ignoreDoguHealth bool, ignoreComponentHealth bool) (ecosystem.HealthResult, error)
}

type selfUpgradeUseCase interface {
	HandleSelfUpgrade(ctx context.Context, blueprintId string) error
}

type ecosystemConfigUseCase interface {
	ApplyConfig(ctx context.Context, blueprintId string) error
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

type maintenanceMode interface {
	domainservice.MaintenanceMode
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
