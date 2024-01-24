package application

import (
	"context"
	"fmt"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
)

type ComponentInstallationUseCase struct {
	componentRepo       domainservice.ComponentInstallationRepository
	healthConfigRepo    domainservice.HealthConfigRepository
	healthCheckInterval time.Duration
}

func NewComponentInstallationUseCase(
	componentRepo domainservice.ComponentInstallationRepository,
	healthConfigRepo domainservice.HealthConfigRepository,
	healthCheckInterval time.Duration,
) *ComponentInstallationUseCase {
	return &ComponentInstallationUseCase{
		componentRepo:       componentRepo,
		healthConfigRepo:    healthConfigRepo,
		healthCheckInterval: healthCheckInterval,
	}
}

func (useCase *ComponentInstallationUseCase) CheckComponentHealth(ctx context.Context) (ecosystem.ComponentHealthResult, error) {
	logger := log.FromContext(ctx).WithName("ComponentInstallationUseCase.CheckComponentHealth")
	logger.Info("check component health...")
	installedComponents, err := useCase.componentRepo.GetAll(ctx)
	if err != nil {
		return ecosystem.ComponentHealthResult{}, fmt.Errorf("cannot retrieve installed components: %w", err)
	}

	healthConfig, err := useCase.healthConfigRepo.Get(ctx)
	if err != nil {
		return ecosystem.ComponentHealthResult{}, fmt.Errorf("cannot retrieve health config: %w", err)
	}

	return ecosystem.CalculateComponentHealthResult(installedComponents, healthConfig.RequiredComponents), nil
}

func (useCase *ComponentInstallationUseCase) WaitForHealthyComponents(ctx context.Context) (ecosystem.ComponentHealthResult, error) {
	logger := log.FromContext(ctx).WithName("ComponentInstallationUseCase.WaitForHealthyComponents")
	logger.Info("start waiting for component health")
	healthResult, err := util.RetryUntilSuccessOrCancellation(
		ctx,
		useCase.healthCheckInterval,
		useCase.checkComponentHealthStatesRetryable,
	)
	var result ecosystem.ComponentHealthResult
	if healthResult == nil {
		result = ecosystem.ComponentHealthResult{}
	} else {
		result = *healthResult
	}

	if err != nil {
		err = fmt.Errorf("stop waiting for component health: %w", err)
		logger.Error(err, "stop waiting for component health because of an error or time out")
	}

	return result, err
}

func (useCase *ComponentInstallationUseCase) checkComponentHealthStatesRetryable(ctx context.Context) (result *ecosystem.ComponentHealthResult, err error, shouldRetry bool) {
	// use named return values to make their meaning clear
	health, err := useCase.CheckComponentHealth(ctx)
	if err != nil {
		// no retry on error while loading components
		return &ecosystem.ComponentHealthResult{}, err, false
	}
	result = &health
	shouldRetry = !health.AllHealthy()
	return
}
