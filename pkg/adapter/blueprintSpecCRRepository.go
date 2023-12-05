package adapter

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/api/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
)

type blueprintSpecCRRepo struct {
	k8sNamespace string
	k8sClient    ecosystem.V1Alpha1Client
}

// GetById returns a Blueprint identified by its ID.
func (repo *blueprintSpecCRRepo) GetById(ctx context.Context, blueprintId string) (domain.BlueprintSpec, error) {
	blueprintCR, err := repo.k8sClient.Blueprints(repo.k8sNamespace).Get(ctx, blueprintId, metav1.GetOptions{})

	if err != nil {
		return domain.BlueprintSpec{}, err
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
func (repo *blueprintSpecCRRepo) Update(ctx context.Context, spec domain.BlueprintSpec) error {
	// TODO
	return nil
}
