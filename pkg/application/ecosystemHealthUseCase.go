package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
)

type EcosystemHealthUseCase struct {
	doguUseCase        doguInstallationUseCase
	componentUseCase   componentInstallationUseCase
	waitConfigProvider healthWaitConfigProvider
}

func NewEcosystemHealthUseCase(
	doguUseCase doguInstallationUseCase,
	componentUseCase componentInstallationUseCase,
	waitConfigProvider domainservice.HealthWaitConfigProvider,
) *EcosystemHealthUseCase {
	return &EcosystemHealthUseCase{
		doguUseCase:        doguUseCase,
		componentUseCase:   componentUseCase,
		waitConfigProvider: waitConfigProvider,
	}
}

// CheckEcosystemHealth checks the ecosystem health once.
// Returns a HealthResult even if parts are unhealthy or
// returns an error if the health state could not be fetched.
func (useCase *EcosystemHealthUseCase) CheckEcosystemHealth(ctx context.Context, ignoreDoguHealth bool, ignoreComponentHealth bool) (ecosystem.HealthResult, error) {
	var doguHealth ecosystem.DoguHealthResult
	var doguHealthErr error
	if !ignoreDoguHealth {
		doguHealth, doguHealthErr = useCase.doguUseCase.CheckDoguHealth(ctx)
	}

	var componentHealth ecosystem.ComponentHealthResult
	var componentHealthErr error
	if !ignoreComponentHealth {
		componentHealth, componentHealthErr = useCase.componentUseCase.CheckComponentHealth(ctx)
	}

	return ecosystem.HealthResult{
		DoguHealth:      doguHealth,
		ComponentHealth: componentHealth,
	}, errors.Join(doguHealthErr, componentHealthErr)
}

// WaitForHealthyEcosystem waits for a healthy ecosystem and returns an HealthResult.
func (useCase *EcosystemHealthUseCase) WaitForHealthyEcosystem(ctx context.Context, ignoreDoguHealth bool, ignoreComponentHealth bool) (ecosystem.HealthResult, error) {
	waitConfig, err := useCase.waitConfigProvider.GetWaitConfig(ctx)
	if err != nil {
		return ecosystem.HealthResult{}, fmt.Errorf("failed to get health check timeout: %w", err)
	}

	timedCtx, cancel := context.WithTimeout(ctx, waitConfig.Timeout)
	defer cancel()

	doguHealthChan := make(chan ecosystem.DoguHealthResult)
	doguErrChan := make(chan error)
	if !ignoreDoguHealth {
		go useCase.asyncWaitForHealthyDogus(timedCtx, doguErrChan, doguHealthChan)
	} else {
		//send empty result, so that wait routine terminates
		doguHealthChan <- ecosystem.DoguHealthResult{}
	}

	componentHealthChan := make(chan ecosystem.ComponentHealthResult)
	componentErrChan := make(chan error)
	if !ignoreComponentHealth {
		go useCase.asyncWaitForHealthyComponents(timedCtx, componentErrChan, componentHealthChan)
	} else {
		//send empty result, so that wait routine terminates
		componentHealthChan <- ecosystem.ComponentHealthResult{}
	}

	return waitForHealthResult(doguHealthChan, doguErrChan, componentHealthChan, componentErrChan)
}

func waitForHealthResult(
	doguHealthChan chan ecosystem.DoguHealthResult,
	doguErrChan chan error,
	componentHealthChan chan ecosystem.ComponentHealthResult,
	componentErrChan chan error,
) (ecosystem.HealthResult, error) {
	var doguHealth ecosystem.DoguHealthResult
	var doguErr error
	var componentHealth ecosystem.ComponentHealthResult
	var componentErr error
	for i := 0; i < 2; i++ {
		select {
		case doguHealth = <-doguHealthChan:
		case doguErr = <-doguErrChan:
		case componentHealth = <-componentHealthChan:
		case componentErr = <-componentErrChan:
		}
	}

	return ecosystem.HealthResult{
		DoguHealth:      doguHealth,
		ComponentHealth: componentHealth,
	}, errors.Join(doguErr, componentErr)
}

func (useCase *EcosystemHealthUseCase) asyncWaitForHealthyComponents(ctx context.Context, componentErrChan chan error, componentHealthChan chan ecosystem.ComponentHealthResult) {
	componentHealth, err := useCase.componentUseCase.WaitForHealthyComponents(ctx)
	if err != nil {
		componentErrChan <- fmt.Errorf("failed to wait for healthy components: %w", err)
		return
	}
	componentHealthChan <- componentHealth
}

func (useCase *EcosystemHealthUseCase) asyncWaitForHealthyDogus(ctx context.Context, doguErrChan chan error, doguHealthChan chan ecosystem.DoguHealthResult) {
	doguHealth, err := useCase.doguUseCase.WaitForHealthyDogus(ctx)
	if err != nil {
		doguErrChan <- fmt.Errorf("failed to wait for healthy dogus: %w", err)
		return
	}
	doguHealthChan <- doguHealth
}
