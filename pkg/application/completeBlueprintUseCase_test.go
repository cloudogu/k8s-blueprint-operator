package application

import (
	"testing"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
)

func TestNewApplyBlueprintSpecUseCase(t *testing.T) {
	repoMock := newMockBlueprintSpecRepository(t)

	sut := NewCompleteBlueprintUseCase(repoMock)

	assert.Equal(t, repoMock, sut.repo)
}

func TestApplyBlueprintSpecUseCase_CompleteBlueprint(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil)
		useCase := NewCompleteBlueprintUseCase(repoMock)

		err := useCase.CompleteBlueprint(testCtx, blueprint)

		require.NoError(t, err)

		assert.True(t, meta.IsStatusConditionTrue(blueprint.Conditions, domain.ConditionCompleted))
	})
	t.Run("stateDiffNotEnmptyError, if StateDiff not empty", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
			StateDiff: domain.StateDiff{
				DoguDiffs: []domain.DoguDiff{
					{
						DoguName:      "ldap",
						NeededActions: []domain.Action{domain.ActionUpgrade},
					},
				},
			},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		useCase := NewCompleteBlueprintUseCase(repoMock)

		err := useCase.CompleteBlueprint(testCtx, blueprint)

		require.Error(t, err)
		var targetErr *domain.StateDiffNotEmptyError
		assert.ErrorAs(t, err, &targetErr)
		assert.ErrorContains(t, err, "cannot complete blueprint because the StateDiff has still changes")
		assert.Empty(t, blueprint.Conditions)
	})
	t.Run("no change if already completed", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil).Once()
		useCase := NewCompleteBlueprintUseCase(repoMock)

		err := useCase.CompleteBlueprint(testCtx, blueprint)
		require.NoError(t, err)
		err = useCase.CompleteBlueprint(testCtx, blueprint)
		require.NoError(t, err)
		assert.True(t, meta.IsStatusConditionTrue(blueprint.Conditions, domain.ConditionCompleted))
	})
	t.Run("repo error while saving", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(assert.AnError)
		useCase := NewCompleteBlueprintUseCase(repoMock)

		err := useCase.CompleteBlueprint(testCtx, blueprint)

		require.ErrorIs(t, err, assert.AnError)
	})
}
