package application

import (
	"context"
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InitiateBlueprintStatusUseCase contains all use cases which are needed to initiate the
// blueprint status after the determining the state diff.
//
// Use cases:
// - InitateConditions
type InitiateBlueprintStatusUseCase struct {
	repo blueprintSpecRepository
}

func NewInitiateBlueprintStatusUseCase(
	repo blueprintSpecRepository,
) *InitiateBlueprintStatusUseCase {
	return &InitiateBlueprintStatusUseCase{
		repo: repo,
	}
}

// InitateConditions handles the initial setting of the conditions to unknown if they are not set yet.
// returns a domainservice.InternalError on any error.
func (useCase *InitiateBlueprintStatusUseCase) InitateConditions(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	if len(blueprint.Conditions) != len(domain.BlueprintConditions) {
		for _, condition := range domain.BlueprintConditions {
			if meta.FindStatusCondition(blueprint.Conditions, condition) == nil {
				meta.SetStatusCondition(&blueprint.Conditions, metav1.Condition{
					Type:    condition,
					Status:  metav1.ConditionUnknown,
					Reason:  "InitialSyncPending",
					Message: "controller has not determined this condition yet",
				})
			}
		}
		err := useCase.repo.Update(ctx, blueprint)
		if err != nil {
			return fmt.Errorf("cannot save blueprint spec %q after initially setting the conditions to unknown: %w", blueprint.Id, err)
		}
	}
	return nil
}
