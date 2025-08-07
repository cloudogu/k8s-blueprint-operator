package application

import (
	"context"
	"errors"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
)

type EcosystemHealthUseCase struct {
	doguUseCase      doguInstallationUseCase
	componentUseCase componentInstallationUseCase
}

func NewEcosystemHealthUseCase(
	doguUseCase doguInstallationUseCase,
	componentUseCase componentInstallationUseCase,
) *EcosystemHealthUseCase {
	return &EcosystemHealthUseCase{
		doguUseCase:      doguUseCase,
		componentUseCase: componentUseCase,
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
