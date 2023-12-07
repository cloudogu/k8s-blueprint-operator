package kubernetes

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/retry"
)

type blueprintSpecRepo struct {
	k8sNamespace    string
	ecosystemClient Interface
}

// NewBlueprintSpecRepository returns a new BlueprintSpecRepository to interact on BlueprintSpecs.
func NewBlueprintSpecRepository(k8sNamespace string, ecosystemClient Interface) domainservice.BlueprintSpecRepository {
	return &blueprintSpecRepo{k8sNamespace: k8sNamespace, ecosystemClient: ecosystemClient}
}

// GetById returns a Blueprint identified by its ID.
func (repo *blueprintSpecRepo) GetById(ctx context.Context, blueprintId string) (domain.BlueprintSpec, error) {
	blueprintCR, err := repo.ecosystemClient.EcosystemV1Alpha1().Blueprints(repo.k8sNamespace).Get(ctx, blueprintId, metav1.GetOptions{})
	if err != nil {
		return domain.BlueprintSpec{}, fmt.Errorf("error while accessing blueprint ID=%s: %w", blueprintId, err)
	}

	blueprint, err := serializer.DeserializeBlueprint([]byte(blueprintCR.Spec.Blueprint))
	if err != nil {
		return domain.BlueprintSpec{}, fmt.Errorf("could not deserialize blueprint for blueprint CR %s", blueprintId)
	}

	blueprintMask, err := serializer.DeserializeBlueprintMask([]byte(blueprintCR.Spec.BlueprintMask))
	if err != nil {
		return domain.BlueprintSpec{}, fmt.Errorf("could not deserialize blueprint mask for blueprint CR %s", blueprintId)
	}

	return domain.BlueprintSpec{
		Id:            blueprintId,
		Blueprint:     blueprint,
		BlueprintMask: blueprintMask,
		Status:        domain.StatusPhase(blueprintCR.Status.Phase),
	}, nil
}

// Update updates a given BlueprintSpec.
func (repo *blueprintSpecRepo) Update(ctx context.Context, spec domain.BlueprintSpec) error {
	return retry.OnConflict(func() error {
		blueprintCli := repo.ecosystemClient.EcosystemV1Alpha1().Blueprints(repo.k8sNamespace)

		updatedBlueprint, err := blueprintCli.Get(ctx, spec.Id, metav1.GetOptions{})
		if err != nil {
			// TODO add logging and event writing
			return err
		}

		blueprintString, err := serializer.SerializeBlueprint(spec.Blueprint)
		if err != nil {
			// TODO add logging and event writing
			_, err2 := blueprintCli.UpdateStatusFailed(ctx, updatedBlueprint)
			if err2 != nil {
				return fmt.Errorf("could not serialize blueprint for blueprint CR %s and also failed to update blueprint: %w", spec.Id, err)
			}
			return fmt.Errorf("could not serialize blueprint for blueprint CR %s: %w", spec.Id, err)
		}
		updatedBlueprint.Spec.Blueprint = blueprintString

		blueprintMaskString, err := serializer.SerializeBlueprintMask(spec.BlueprintMask)
		if err != nil {
			// TODO add logging and event writing
			_, err2 := blueprintCli.UpdateStatusFailed(ctx, updatedBlueprint)
			if err2 != nil {
				return fmt.Errorf("could not serialize blueprint mask for blueprint CR %s and also failed to update blueprint: %w", spec.Id, err)
			}
			return fmt.Errorf("could not serialize blueprint mask for blueprint CR %s: %w", spec.Id, err)
		}
		updatedBlueprint.Spec.BlueprintMask = blueprintMaskString

		// TODO add to Status: effective Blueprint, statediff, upgradePlan?

		return err
	})
}
