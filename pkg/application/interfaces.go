package application

import (
	"context"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
)

type initialBlueprintStatusUseCase interface {
	InitateConditions(ctx context.Context, blueprint *domain.BlueprintSpec) error
}

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
	CheckDogusUpToDate(ctx context.Context) ([]cescommons.SimpleName, error)
	ApplyDoguStates(ctx context.Context, blueprint *domain.BlueprintSpec) error
}

type applyDogusUseCase interface {
	ApplyDogus(ctx context.Context, blueprint *domain.BlueprintSpec) (bool, error)
}

type completeBlueprintUseCase interface {
	CompleteBlueprint(ctx context.Context, blueprint *domain.BlueprintSpec) error
}

type ecosystemHealthUseCase interface {
	CheckEcosystemHealth(context.Context, *domain.BlueprintSpec) (ecosystem.HealthResult, error)
}

type dogusUpToDateUseCase interface {
	CheckDogus(ctx context.Context, blueprint *domain.BlueprintSpec) error
}

type ecosystemConfigUseCase interface {
	ApplyConfig(ctx context.Context, blueprint *domain.BlueprintSpec) error
}

type restoreInProgressUseCase interface {
	CheckRestoreInProgress(context.Context) error
}

type doguInstallationRepository interface {
	domainservice.DoguInstallationRepository
}

//nolint:unused
//goland:noinspection GoUnusedType
type blueprintSpecRepository interface {
	domainservice.BlueprintSpecRepository
}

//nolint:unused
//goland:noinspection GoUnusedType
type debugModeRepository interface {
	domainservice.DebugModeRepository
}

//nolint:unused
//goland:noinspection GoUnusedType
type restoreRepository interface {
	domainservice.RestoreRepository
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

type configRefReader interface {
	domainservice.ConfigRefReader
}

// validateDependenciesDomainUseCase is an interface for the domain service for better testability
type validateDependenciesDomainUseCase interface {
	ValidateDependenciesForAllDogus(ctx context.Context, effectiveBlueprint domain.EffectiveBlueprint) error
}

// validateAdditionalMountsDomainUseCase is an interface for the domain service for better testability
type validateAdditionalMountsDomainUseCase interface {
	ValidateAdditionalMounts(ctx context.Context, effectiveBlueprint domain.EffectiveBlueprint) error
}
