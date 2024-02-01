package application

import (
	"context"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

var testCtx = context.Background()

func TestBlueprintSpecChangeUseCase_HandleChange(t *testing.T) {

	t.Run("do all steps with blueprint", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock)

		blueprintSpec := &domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseNew,
		}
		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprintSpec, nil)
		validationMock.EXPECT().ValidateBlueprintSpecStatically(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseStaticallyValidated
			})
		effectiveBlueprintMock.EXPECT().CalculateEffectiveBlueprint(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseEffectiveBlueprintGenerated
			})
		validationMock.EXPECT().ValidateBlueprintSpecDynamically(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseValidated
			})
		stateDiffMock.EXPECT().DetermineStateDiff(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseStateDiffDetermined
			})
		applyMock.EXPECT().CheckEcosystemHealthUpfront(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseEcosystemHealthyUpfront
			})
		applyMock.EXPECT().ApplyBlueprintSpec(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseCompleted
			})

		// when
		err := useCase.HandleChange(testCtx, "testBlueprint1")
		// then
		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseCompleted, blueprintSpec.Status)
	})

	t.Run("should return nil and not handle blueprint spec on dry run", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock)

		blueprintSpec := &domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseNew,
			Config: domain.BlueprintConfiguration{DryRun: true},
		}

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprintSpec, nil)
		validationMock.EXPECT().ValidateBlueprintSpecStatically(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseStaticallyValidated
			})
		effectiveBlueprintMock.EXPECT().CalculateEffectiveBlueprint(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseEffectiveBlueprintGenerated
			})
		validationMock.EXPECT().ValidateBlueprintSpecDynamically(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseValidated
			})
		stateDiffMock.EXPECT().DetermineStateDiff(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseStateDiffDetermined
			})
		applyMock.EXPECT().CheckEcosystemHealthUpfront(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseEcosystemHealthyUpfront
			})
		applyMock.EXPECT().ApplyBlueprintSpec(testCtx, "testBlueprint1").Return(nil)

		// when
		err := useCase.HandleChange(testCtx, "testBlueprint1")

		// then
		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseEcosystemHealthyUpfront, blueprintSpec.Status)
	})

	t.Run("cannot load blueprint spec initially", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock)

		expectedError := &domainservice.InternalError{
			WrappedError: nil,
			Message:      "test-error",
		}
		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(nil, expectedError)

		// when
		err := useCase.HandleChange(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, expectedError)
	})

	t.Run("new with static validation error", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock)

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseNew,
			Blueprint: domain.Blueprint{Dogus: []domain.Dogu{
				{Namespace: "official", Name: "DoguWithNoVersion", TargetState: domain.TargetStatePresent},
			}},
		}, nil)
		validationMock.EXPECT().ValidateBlueprintSpecStatically(testCtx, "testBlueprint1").Return(assert.AnError)

		// when
		err := useCase.HandleChange(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("new with error calculating effective blueprint", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock)

		updatedSpec := &domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseNew,
		}

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(updatedSpec, nil)
		validationMock.EXPECT().ValidateBlueprintSpecStatically(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				updatedSpec.Status = domain.StatusPhaseStaticallyValidated
			})
		effectiveBlueprintMock.EXPECT().CalculateEffectiveBlueprint(testCtx, "testBlueprint1").Return(assert.AnError)

		// when
		err := useCase.HandleChange(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("new with dynamic validation error", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock)

		updatedSpec := &domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseNew,
		}

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(updatedSpec, nil)
		validationMock.EXPECT().ValidateBlueprintSpecStatically(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				updatedSpec.Status = domain.StatusPhaseStaticallyValidated
			})
		effectiveBlueprintMock.EXPECT().CalculateEffectiveBlueprint(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				updatedSpec.Status = domain.StatusPhaseEffectiveBlueprintGenerated
			})
		validationMock.EXPECT().ValidateBlueprintSpecDynamically(testCtx, "testBlueprint1").Return(assert.AnError)

		// when
		err := useCase.HandleChange(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("new with error determining state diff", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock)

		updatedSpec := &domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseNew,
		}

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(updatedSpec, nil)
		validationMock.EXPECT().ValidateBlueprintSpecStatically(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				updatedSpec.Status = domain.StatusPhaseStaticallyValidated
			})
		effectiveBlueprintMock.EXPECT().CalculateEffectiveBlueprint(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				updatedSpec.Status = domain.StatusPhaseEffectiveBlueprintGenerated
			})
		validationMock.EXPECT().ValidateBlueprintSpecDynamically(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				updatedSpec.Status = domain.StatusPhaseValidated
			})
		stateDiffMock.EXPECT().DetermineStateDiff(testCtx, "testBlueprint1").Return(assert.AnError)
		// when
		err := useCase.HandleChange(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("new with error checking dogu health", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock)

		updatedSpec := &domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseNew,
		}

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(updatedSpec, nil)
		validationMock.EXPECT().ValidateBlueprintSpecStatically(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				updatedSpec.Status = domain.StatusPhaseStaticallyValidated
			})
		effectiveBlueprintMock.EXPECT().CalculateEffectiveBlueprint(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				updatedSpec.Status = domain.StatusPhaseEffectiveBlueprintGenerated
			})
		validationMock.EXPECT().ValidateBlueprintSpecDynamically(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				updatedSpec.Status = domain.StatusPhaseValidated
			})
		stateDiffMock.EXPECT().DetermineStateDiff(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				updatedSpec.Status = domain.StatusPhaseStateDiffDetermined
			})
		applyMock.EXPECT().CheckEcosystemHealthUpfront(testCtx, "testBlueprint1").Return(assert.AnError)

		// when
		err := useCase.HandleChange(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("handle invalid blueprint", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock)

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseInvalid,
		}, nil)
		// when
		err := useCase.HandleChange(testCtx, "testBlueprint1")
		// then
		require.NoError(t, err)
	})

	t.Run("handle healthy ecosystem upfront", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock)

		blueprintSpec := &domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}
		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprintSpec, nil)
		applyMock.EXPECT().ApplyBlueprintSpec(testCtx, "testBlueprint1").Return(nil).Run(func(ctx context.Context, blueprintId string) {
			blueprintSpec.MarkCompleted()
		})
		// when
		err := useCase.HandleChange(testCtx, "testBlueprint1")
		// then
		require.NoError(t, err)
	})

	t.Run("handle unhealthy ecosystem upfront", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock)

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseEcosystemUnhealthyUpfront,
		}, nil)
		// when
		err := useCase.HandleChange(testCtx, "testBlueprint1")
		// then
		require.NoError(t, err)
	})

	t.Run("handle in progress blueprint", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock)

		blueprintSpec := &domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseInProgress,
		}
		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprintSpec, nil)
		applyMock.EXPECT().MarkFailed(testCtx, blueprintSpec, mock.Anything).Return(nil).Run(func(ctx context.Context, blueprintSpec *domain.BlueprintSpec, err error) {
			assert.ErrorContains(t, err, handleInProgressMsg)
			blueprintSpec.MarkFailed(err)
		})
		// when
		actualErr := useCase.HandleChange(testCtx, "testBlueprint1")
		// then
		require.NoError(t, actualErr)
	})

	t.Run("handle blueprint applied", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock)

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseBlueprintApplied,
		}, nil)
		applyMock.EXPECT().CheckEcosystemHealthAfterwards(testCtx, "testBlueprint1").Return(nil)
		// when
		err := useCase.HandleChange(testCtx, "testBlueprint1")
		// then
		require.NoError(t, err)
	})

	t.Run("handle ecosystem healthy afterwards", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock)

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseEcosystemHealthyAfterwards,
		}, nil)
		// when
		err := useCase.HandleChange(testCtx, "testBlueprint1")
		// then
		require.NoError(t, err)
	})

	t.Run("handle ecosystem unhealthy afterwards", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock)

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseEcosystemUnhealthyAfterwards,
		}, nil)
		// when
		err := useCase.HandleChange(testCtx, "testBlueprint1")
		// then
		require.NoError(t, err)
	})

	t.Run("handle completed blueprint", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock)

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseCompleted,
		}, nil)
		// when
		err := useCase.HandleChange(testCtx, "testBlueprint1")
		// then
		require.NoError(t, err)
	})

	t.Run("handle failed blueprint", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock)

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseFailed,
		}, nil)
		// when
		err := useCase.HandleChange(testCtx, "testBlueprint1")
		// then
		require.NoError(t, err)
	})

	t.Run("handle unknown status in blueprint", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock)

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: "unknown",
		}, nil)
		// when
		err := useCase.HandleChange(testCtx, "testBlueprint1")
		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "could not handle unknown status of blueprint")
	})
}
