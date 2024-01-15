package application

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"time"
)

type EcosystemHealthUseCase struct {
	doguUseCase        doguInstallationUseCase
	healthCheckTimeOut time.Duration
}

func NewEcosystemHealthUseCase(
	doguUseCase doguInstallationUseCase,
	healthCheckTimeOut time.Duration,
) *EcosystemHealthUseCase {
	return &EcosystemHealthUseCase{
		doguUseCase:        doguUseCase,
		healthCheckTimeOut: healthCheckTimeOut,
	}
}

func (useCase *EcosystemHealthUseCase) CheckEcosystemHealth(ctx context.Context, ignoreDoguHealth bool) (ecosystem.HealthResult, error) {
	doguHealth := ecosystem.DoguHealthResult{}
	if !ignoreDoguHealth {
		var err error
		doguHealth, err = useCase.doguUseCase.CheckDoguHealthStates(ctx)
		if err != nil {
			return ecosystem.HealthResult{}, err
		}
	}

	return ecosystem.HealthResult{
		DoguHealth: doguHealth,
	}, nil
}

func (useCase *EcosystemHealthUseCase) WaitForHealthyEcosystem(ctx context.Context) (ecosystem.HealthResult, error) {
	timedCtx, cancel := context.WithTimeout(ctx, useCase.healthCheckTimeOut)
	defer cancel()

	doguHealth, err := useCase.doguUseCase.WaitForHealthyDogus(timedCtx)
	if err != nil {
		return ecosystem.HealthResult{}, err
	}

	return ecosystem.HealthResult{
		DoguHealth: doguHealth,
	}, nil
}
