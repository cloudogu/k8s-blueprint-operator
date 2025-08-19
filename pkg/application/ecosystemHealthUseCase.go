package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
)

type EcosystemHealthUseCase struct {
	doguUseCase      doguInstallationUseCase
	componentUseCase componentInstallationUseCase
	blueprintRepo    blueprintSpecRepository
}

func NewEcosystemHealthUseCase(
	doguUseCase doguInstallationUseCase,
	componentUseCase componentInstallationUseCase,
	blueprintRepo blueprintSpecRepository,
) *EcosystemHealthUseCase {
	return &EcosystemHealthUseCase{
		doguUseCase:      doguUseCase,
		componentUseCase: componentUseCase,
		blueprintRepo:    blueprintRepo,
	}
}

// CheckEcosystemHealth checks the ecosystem health once and sets the health condition accordingly.
// Returns the health result.
// Returns a domain.UnhealthyEcosystemError and the ecosystem.HealthResult if the ecosystem is unhealthy or
// returns a domainservice.ConflictError if there was a conflicting update to the blueprint or
// returns a domainservice.InternalError if the health status could not be determined or the there was any another problem.
func (useCase *EcosystemHealthUseCase) CheckEcosystemHealth(
	ctx context.Context,
	blueprint *domain.BlueprintSpec,
) (ecosystem.HealthResult, error) {
	health, determineHealthError := useCase.getEcosystemHealth(
		ctx,
		blueprint.Config.IgnoreDoguHealth,
		blueprint.Config.IgnoreComponentHealth,
	)
	infoChanged := blueprint.HandleHealthResult(health, determineHealthError)
	if infoChanged {
		updateErr := useCase.blueprintRepo.Update(ctx, blueprint)
		if updateErr != nil {
			return ecosystem.HealthResult{}, fmt.Errorf(
				"could not update health condition after health check: %w",
				errors.Join(updateErr, determineHealthError),
			)
		}
	}
	if determineHealthError == nil && !health.AllHealthy() {
		return health, domain.NewUnhealthyEcosystemError(nil, "ecosystem is unhealthy", health)
	}

	return health, determineHealthError
}

// getEcosystemHealth checks the ecosystem health once.
// Returns a HealthResult even if parts are unhealthy or
// returns an error if the health state could not be fetched.
func (useCase *EcosystemHealthUseCase) getEcosystemHealth(
	ctx context.Context,
	ignoreDoguHealth bool,
	ignoreComponentHealth bool,
) (ecosystem.HealthResult, error) {
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
