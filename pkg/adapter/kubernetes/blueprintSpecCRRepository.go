package kubernetes

import (
	"context"
	"errors"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/retry"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
)

type blueprintSpecRepo struct {
	blueprintClient         BlueprintInterface
	blueprintSerializer     serializer.BlueprintSerializer
	blueprintMaskSerializer serializer.BlueprintMaskSerializer
}

// NewBlueprintSpecRepository returns a new BlueprintSpecRepository to interact on BlueprintSpecs.
func NewBlueprintSpecRepository(
	blueprintClient BlueprintInterface,
	blueprintSerializer serializer.BlueprintSerializer,
	blueprintMaskSerializer serializer.BlueprintMaskSerializer,
) domainservice.BlueprintSpecRepository {
	return &blueprintSpecRepo{
		blueprintClient:         blueprintClient,
		blueprintSerializer:     blueprintSerializer,
		blueprintMaskSerializer: blueprintMaskSerializer,
	}
}

// GetById returns a Blueprint identified by its ID.
func (repo *blueprintSpecRepo) GetById(ctx context.Context, blueprintId string) (domain.BlueprintSpec, error) {
	blueprintCR, err := repo.blueprintClient.Get(ctx, blueprintId, metav1.GetOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return domain.BlueprintSpec{}, &domainservice.NotFoundError{
				WrappedError: err,
				Message:      fmt.Sprintf("cannot load Blueprint CR '%s' as it does not exist", blueprintId),
			}
		}
		return domain.BlueprintSpec{}, &domainservice.InternalError{
			WrappedError: err,
			Message:      fmt.Sprintf("error while loading blueprint CR '%s'", blueprintId),
		}
	}

	blueprint, blueprintErr := repo.blueprintSerializer.Deserialize(blueprintCR.Spec.Blueprint)
	blueprintMask, maskErr := repo.blueprintMaskSerializer.Deserialize(blueprintCR.Spec.BlueprintMask)
	err = errors.Join(blueprintErr, maskErr)
	if err != nil {
		return domain.BlueprintSpec{}, fmt.Errorf("could not deserialize Blueprint CR %s: %w", blueprintId, err)
	}

	return domain.BlueprintSpec{
		Id:                   blueprintId,
		Blueprint:            blueprint,
		BlueprintMask:        blueprintMask,
		EffectiveBlueprint:   domain.EffectiveBlueprint{},
		StateDiff:            domain.StateDiff{},
		BlueprintUpgradePlan: domain.BlueprintUpgradePlan{},
		Status:               domain.StatusPhase(blueprintCR.Status.Phase),
		Config: domain.BlueprintConfiguration{
			IgnoreDoguHealth:         blueprintCR.Spec.IgnoreDoguHealth,
			AllowDoguNamespaceSwitch: blueprintCR.Spec.AllowDoguNamespaceSwitch,
		},
	}, nil
}

// Update updates a given BlueprintSpec.
func (repo *blueprintSpecRepo) Update(ctx context.Context, spec domain.BlueprintSpec) error {
	return retry.OnConflict(func() error {

		updatedBlueprint, err := repo.blueprintClient.Get(ctx, spec.Id, metav1.GetOptions{})
		if err != nil {
			// TODO add logging and event writing
			return err
		}

		blueprintString, err := repo.blueprintSerializer.Serialize(spec.Blueprint)
		if err != nil {
			// TODO add logging and event writing
			_, err2 := repo.blueprintClient.UpdateStatusFailed(ctx, updatedBlueprint)
			if err2 != nil {
				return fmt.Errorf("could not serialize blueprint for blueprint CR %s and also failed to update blueprint: %w", spec.Id, err)
			}
			return fmt.Errorf("could not serialize blueprint for blueprint CR %s: %w", spec.Id, err)
		}
		updatedBlueprint.Spec.Blueprint = blueprintString

		blueprintMaskString, err := repo.blueprintMaskSerializer.Serialize(spec.BlueprintMask)
		if err != nil {
			// TODO add logging and event writing
			_, err2 := repo.blueprintClient.UpdateStatusFailed(ctx, updatedBlueprint)
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
