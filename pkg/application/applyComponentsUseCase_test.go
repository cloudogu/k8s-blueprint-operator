package application

import (
	"testing"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
)

func TestApplyComponentsUseCase_ApplyComponents(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				ComponentDiffs: []domain.ComponentDiff{
					{
						Name: "k8s-dogu-operator",
						NeededActions: []domain.Action{
							domain.ActionUpgrade,
						},
					},
				},
			},
			Conditions: []domain.Condition{},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil)
		componentInstallUseCaseMock := newMockComponentInstallationUseCase(t)
		componentInstallUseCaseMock.EXPECT().ApplyComponentStates(testCtx, blueprint).Return(nil)
		useCase := NewApplyComponentsUseCase(repoMock, componentInstallUseCaseMock)

		changed, err := useCase.ApplyComponents(testCtx, blueprint)

		require.NoError(t, err)
		assert.True(t, meta.IsStatusConditionTrue(blueprint.Conditions, domain.ConditionComponentsApplied))
		assert.True(t, changed)
		require.Equal(t, 1, len(blueprint.Events))
		assert.Equal(t, domain.ComponentsAppliedEvent{Diffs: blueprint.StateDiff.ComponentDiffs}, blueprint.Events[0])
	})

	t.Run("no update without condition change", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			StateDiff:  domain.StateDiff{},
			Conditions: []domain.Condition{},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil).Once()
		componentInstallUseCaseMock := newMockComponentInstallationUseCase(t)
		componentInstallUseCaseMock.EXPECT().ApplyComponentStates(testCtx, blueprint).Return(nil).Twice()
		useCase := NewApplyComponentsUseCase(repoMock, componentInstallUseCaseMock)

		changed, err := useCase.ApplyComponents(testCtx, blueprint)
		require.NoError(t, err)
		assert.True(t, meta.IsStatusConditionTrue(blueprint.Conditions, domain.ConditionComponentsApplied))
		assert.True(t, changed)
		changed, err = useCase.ApplyComponents(testCtx, blueprint)
		require.NoError(t, err)
		assert.True(t, meta.IsStatusConditionTrue(blueprint.Conditions, domain.ConditionComponentsApplied))
		assert.False(t, changed)
	})

	t.Run("fail to apply components", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil)
		componentInstallUseCaseMock := newMockComponentInstallationUseCase(t)
		componentInstallUseCaseMock.EXPECT().ApplyComponentStates(testCtx, blueprint).Return(assert.AnError)
		useCase := NewApplyComponentsUseCase(repoMock, componentInstallUseCaseMock)

		changed, err := useCase.ApplyComponents(testCtx, blueprint)

		require.ErrorIs(t, err, assert.AnError)
		assert.True(t, meta.IsStatusConditionFalse(blueprint.Conditions, domain.ConditionComponentsApplied))
		assert.True(t, changed)
	})

	t.Run("fail to update blueprint", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(assert.AnError)
		componentInstallUseCaseMock := newMockComponentInstallationUseCase(t)
		componentInstallUseCaseMock.EXPECT().ApplyComponentStates(testCtx, blueprint).Return(nil)
		useCase := NewApplyComponentsUseCase(repoMock, componentInstallUseCaseMock)

		changed, err := useCase.ApplyComponents(testCtx, blueprint)

		require.ErrorIs(t, err, assert.AnError)
		assert.True(t, meta.IsStatusConditionTrue(blueprint.Conditions, domain.ConditionComponentsApplied))
		assert.True(t, changed)
	})
}
