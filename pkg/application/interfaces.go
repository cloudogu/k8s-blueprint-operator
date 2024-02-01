package application

import (
	"context"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type blueprintSpecValidationUseCase interface {
	ValidateBlueprintSpecStatically(ctx context.Context, blueprintId string) error
	ValidateBlueprintSpecDynamically(ctx context.Context, blueprintId string) error
}

type effectiveBlueprintUseCase interface {
	CalculateEffectiveBlueprint(ctx context.Context, blueprintId string) error
}

type stateDiffUseCase interface {
	DetermineStateDiff(ctx context.Context, blueprintId string) error
}

type doguInstallationUseCase interface {
	CheckDoguHealth(ctx context.Context) (ecosystem.DoguHealthResult, error)
	WaitForHealthyDogus(ctx context.Context) (ecosystem.DoguHealthResult, error)
	ApplyDoguStates(ctx context.Context, blueprintId string) error
}

type componentInstallationUseCase interface {
	CheckComponentHealth(ctx context.Context) (ecosystem.ComponentHealthResult, error)
	WaitForHealthyComponents(ctx context.Context) (ecosystem.ComponentHealthResult, error)
}

type applyBlueprintSpecUseCase interface {
	CheckEcosystemHealthUpfront(ctx context.Context, blueprintId string) error
	CheckEcosystemHealthAfterwards(ctx context.Context, blueprintId string) error
	ApplyBlueprintSpec(ctx context.Context, blueprintId string) error
	MarkFailed(ctx context.Context, blueprintSpec *domain.BlueprintSpec, err error) error
}

type ecosystemHealthUseCase interface {
	CheckEcosystemHealth(ctx context.Context, ignoreDoguHealth bool, ignoreComponentHealth bool) (ecosystem.HealthResult, error)
	WaitForHealthyEcosystem(ctx context.Context) (ecosystem.HealthResult, error)
}

type doguInstallationRepository interface {
	domainservice.DoguInstallationRepository
}

type componentInstallationRepository interface {
	domainservice.ComponentInstallationRepository
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

type blueprintSpecRepository interface {
	domainservice.BlueprintSpecRepository
}

// interface duplication for mocks

//nolint:unused
//goland:noinspection GoUnusedType
type remoteDoguRegistry interface {
	domainservice.RemoteDoguRegistry
}
