package application

import (
	"context"
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
	CheckDoguHealth(ctx context.Context, blueprintId string) error
}

// interface duplication for mocks

//nolint:unused
//goland:noinspection GoUnusedType
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
type remoteDoguRegistry interface {
	domainservice.RemoteDoguRegistry
}
