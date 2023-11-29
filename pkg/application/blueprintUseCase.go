package application

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
)
import "github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"

type BlueprintUseCase struct {
	doguUseCase DoguInstallationUseCase
}

func (useCase BlueprintUseCase) validateBlueprintStatically(blueprint domain.BlueprintV2) error {
	err := blueprint.Validate()
	if err != nil {
		err = fmt.Errorf("blueprint is invalid: %w", err)
	}
	return err
}

func (useCase BlueprintUseCase) validateEcosystemState(blueprint domain.BlueprintV2) error {
	err := blueprint.Validate()
	if err != nil {
		err = fmt.Errorf("blueprint is invalid: %w", err)
	}
	return err
}

func (useCase BlueprintUseCase) MakeBlueprintUpgrade(blueprint domain.BlueprintV2, mask domain.BlueprintMaskV1) {
	//TODO
	_, _ = domainservice.CalculateEffectiveBlueprint(blueprint, mask)
}
