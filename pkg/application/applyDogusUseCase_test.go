package application

import (
	"testing"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
)

func TestApplyDogusUseCase_ApplyDogus(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: &[]domain.Condition{},
			StateDiff: domain.StateDiff{
				DoguDiffs: domain.DoguDiffs{
					{
						DoguName: "cas",
						NeededActions: []domain.Action{
							domain.ActionUpgrade,
						},
					},
				},
			},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil)
		doguInstallUseCaseMock := newMockDoguInstallationUseCase(t)
		doguInstallUseCaseMock.EXPECT().ApplyDoguStates(testCtx, blueprint).Return(nil)
		useCase := NewApplyDogusUseCase(repoMock, doguInstallUseCaseMock)

		err := useCase.ApplyDogus(testCtx, blueprint)

		require.NoError(t, err)
		assert.True(t, meta.IsStatusConditionTrue(*blueprint.Conditions, domain.ConditionDogusApplied))
		require.Equal(t, 1, len(blueprint.Events))
		assert.Equal(t, domain.DogusAppliedEvent{Diffs: blueprint.StateDiff.DoguDiffs}, blueprint.Events[0])
	})

	t.Run("no update without condition change", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: &[]domain.Condition{},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil).Once()
		doguInstallUseCaseMock := newMockDoguInstallationUseCase(t)
		doguInstallUseCaseMock.EXPECT().ApplyDoguStates(testCtx, blueprint).Return(nil)
		useCase := NewApplyDogusUseCase(repoMock, doguInstallUseCaseMock)

		err := useCase.ApplyDogus(testCtx, blueprint)
		require.NoError(t, err)
		assert.True(t, meta.IsStatusConditionTrue(*blueprint.Conditions, domain.ConditionDogusApplied))
		err = useCase.ApplyDogus(testCtx, blueprint)
		require.NoError(t, err)
		assert.True(t, meta.IsStatusConditionTrue(*blueprint.Conditions, domain.ConditionDogusApplied))
	})

	t.Run("fail to apply dogus", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: &[]domain.Condition{},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil)
		doguInstallUseCaseMock := newMockDoguInstallationUseCase(t)
		doguInstallUseCaseMock.EXPECT().ApplyDoguStates(testCtx, blueprint).Return(assert.AnError)
		useCase := NewApplyDogusUseCase(repoMock, doguInstallUseCaseMock)

		err := useCase.ApplyDogus(testCtx, blueprint)

		require.ErrorIs(t, err, assert.AnError)
		assert.True(t, meta.IsStatusConditionFalse(*blueprint.Conditions, domain.ConditionDogusApplied))
	})

	t.Run("fail to update blueprint", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: &[]domain.Condition{},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(assert.AnError)
		doguInstallUseCaseMock := newMockDoguInstallationUseCase(t)
		doguInstallUseCaseMock.EXPECT().ApplyDoguStates(testCtx, blueprint).Return(nil)
		useCase := NewApplyDogusUseCase(repoMock, doguInstallUseCaseMock)

		err := useCase.ApplyDogus(testCtx, blueprint)

		require.ErrorIs(t, err, assert.AnError)
		assert.True(t, meta.IsStatusConditionTrue(*blueprint.Conditions, domain.ConditionDogusApplied))
	})
}
