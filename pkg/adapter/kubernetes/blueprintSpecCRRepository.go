package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	v1 "github.com/cloudogu/k8s-blueprint-operator/pkg/api/v1"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/retry"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// Update persists changes in the blueprint to the corresponding blueprint CR.
func (repo *blueprintSpecRepo) Update(ctx context.Context, spec domain.BlueprintSpec) error {
	return retry.OnConflict(func() error {
		updatedBlueprint := v1.Blueprint{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:              spec.Id,
				ResourceVersion:   "",
				CreationTimestamp: metav1.Time{},
			},
			Status: v1.BlueprintStatus{
				Phase: v1.StatusPhase(spec.Status),
			},
		}

		_, err := repo.blueprintClient.UpdateStatus(ctx, &updatedBlueprint, metav1.UpdateOptions{})
		if err != nil {
			if k8sErrors.IsConflict(err) {
				return &domainservice.ConflictError{
					WrappedError: err,
					Message:      fmt.Sprintf("cannot update blueprint CR '%s' as it was modified in the meantime", spec.Id),
				}
			}
			return &domainservice.InternalError{WrappedError: err, Message: fmt.Sprintf("Cannot update blueprint CR '%s'", spec.Id)}
		}
		// TODO add to Status: effective Blueprint, stateDiff, upgradePlan?

		return err
	})
}
