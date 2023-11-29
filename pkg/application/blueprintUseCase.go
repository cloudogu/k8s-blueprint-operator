package application

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
)
import "github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"

type BlueprintUseCase struct {
	repo        BlueprintRepository
	doguUseCase DoguInstallationUseCase
}

func (useCase BlueprintUseCase) ValidateBlueprintStatically(blueprintId string) error {
	blueprint, err := useCase.repo.getById(blueprintId)
	if err != nil {
		return fmt.Errorf("cannot validate blueprint: %w", err)
	}
	return blueprint.Validate()
}

func (useCase BlueprintUseCase) MakeBlueprintUpgrade(blueprint domain.Blueprint, mask domain.BlueprintMask) {
	//TODO
	_, _ = domainservice.CalculateEffectiveBlueprint(blueprint, mask)
}
