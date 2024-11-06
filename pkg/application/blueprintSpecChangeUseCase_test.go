package application

import (
	"context"
	"errors"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
)

var testCtx = context.Background()
var testBlueprintId = "testBlueprint1"

func TestBlueprintSpecChangeUseCase_HandleChange(t *testing.T) {

	t.Run("do all steps with blueprint", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

		blueprintSpec := &domain.BlueprintSpec{
			Id:     testBlueprintId,
			Status: domain.StatusPhaseNew,
		}
		repoMock.EXPECT().GetById(testCtx, testBlueprintId).Return(blueprintSpec, nil)
		validationMock.EXPECT().ValidateBlueprintSpecStatically(testCtx, testBlueprintId).Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseStaticallyValidated
			})
		effectiveBlueprintMock.EXPECT().CalculateEffectiveBlueprint(testCtx, testBlueprintId).Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseEffectiveBlueprintGenerated
			})
		validationMock.EXPECT().ValidateBlueprintSpecDynamically(testCtx, testBlueprintId).Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseValidated
			})
		stateDiffMock.EXPECT().DetermineStateDiff(testCtx, testBlueprintId).Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseStateDiffDetermined
			})
		applyMock.EXPECT().CheckEcosystemHealthUpfront(testCtx, testBlueprintId).Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseEcosystemHealthyUpfront
			})
		applyMock.EXPECT().PreProcessBlueprintApplication(testCtx, testBlueprintId).Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseBlueprintApplicationPreProcessed
			})
		registryConfigUseCaseMock.EXPECT().ApplyConfig(testCtx, testBlueprintId).Return(nil).Run(func(ctx context.Context, blueprintId string) {
			blueprintSpec.Status = domain.StatusPhaseRegistryConfigApplied
		})
		selfUpgradeUseCase.EXPECT().HandleSelfUpgrade(testCtx, "testBlueprint1").Return(nil).Run(func(ctx context.Context, blueprintId string) {
			blueprintSpec.Status = domain.StatusPhaseSelfUpgradeCompleted
		})
		applyMock.EXPECT().ApplyBlueprintSpec(testCtx, testBlueprintId).Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseBlueprintApplied
			})
		applyMock.EXPECT().CheckEcosystemHealthAfterwards(testCtx, testBlueprintId).Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseEcosystemHealthyAfterwards
			})
		applyMock.EXPECT().PostProcessBlueprintApplication(testCtx, testBlueprintId).Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseCompleted
			})
		doguRestartUseCaseMock.EXPECT().TriggerDoguRestarts(testCtx, testBlueprintId).Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseRestartsTriggered
			})

		// when
		err := useCase.HandleChange(testCtx, testBlueprintId)
		// then
		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseCompleted, blueprintSpec.Status)
	})

	t.Run("cannot load blueprint spec initially", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

		expectedError := &domainservice.InternalError{
			WrappedError: nil,
			Message:      "test-error",
		}
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(nil, expectedError)

		// when
		err := useCase.HandleChange(testCtx, blueprintId)

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
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseNew,
			Blueprint: domain.Blueprint{Dogus: []domain.Dogu{
				{Name: cescommons.QualifiedDoguName{Namespace: "official", SimpleName: "DoguWithNoVersion"}, TargetState: domain.TargetStatePresent},
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
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

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
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

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
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

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
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

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
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseInvalid,
		}, nil)
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
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseEcosystemUnhealthyUpfront,
		}, nil)
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
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseEcosystemUnhealthyUpfront,
		}, nil)
		// when
		err := useCase.HandleChange(testCtx, "testBlueprint1")
		// then
		require.NoError(t, err)
	})

	t.Run("handle error apply registry config", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

		blueprintSpec := &domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseSelfUpgradeCompleted,
		}
		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprintSpec, nil)
		registryConfigUseCaseMock.EXPECT().ApplyConfig(testCtx, blueprintSpec.Id).Return(assert.AnError)
		// when
		actualErr := useCase.HandleChange(testCtx, "testBlueprint1")
		// then
		require.ErrorIs(t, actualErr, assert.AnError)
	})

	t.Run("handle in progress blueprint", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

		blueprintSpec := &domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseInProgress,
		}
		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprintSpec, nil)
		applyMock.EXPECT().PostProcessBlueprintApplication(testCtx, blueprintSpec.Id).Return(nil)
		// when
		actualErr := useCase.HandleChange(testCtx, "testBlueprint1")
		// then
		require.NoError(t, actualErr)
	})

	t.Run("handle error when blueprint is in progress", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

		blueprintSpec := &domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseInProgress,
		}
		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprintSpec, nil)
		applyMock.EXPECT().PostProcessBlueprintApplication(testCtx, blueprintSpec.Id).Return(assert.AnError)
		// when
		actualErr := useCase.HandleChange(testCtx, "testBlueprint1")
		// then
		require.ErrorIs(t, actualErr, assert.AnError)
	})

	t.Run("handle error after blueprint was applied", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

		blueprintSpec := &domain.BlueprintSpec{
			Id:     testBlueprintId,
			Status: domain.StatusPhaseBlueprintApplied,
		}
		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprintSpec, nil)
		applyMock.EXPECT().CheckEcosystemHealthAfterwards(testCtx, testBlueprintId).Return(assert.AnError)
		doguRestartUseCaseMock.EXPECT().TriggerDoguRestarts(testCtx, testBlueprintId).Return(nil).
			Run(func(ctx context.Context, blueprintId string) {
				blueprintSpec.Status = domain.StatusPhaseRestartsTriggered
			})

		// when
		err := useCase.HandleChange(testCtx, testBlueprintId)
		// then
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("handle ecosystem healthy afterwards", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

		blueprintSpec := &domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseEcosystemHealthyAfterwards,
		}
		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprintSpec, nil)
		applyMock.EXPECT().PostProcessBlueprintApplication(testCtx, "testBlueprint1").Return(nil).Run(func(ctx context.Context, blueprintId string) {
			blueprintSpec.Status = domain.StatusPhaseCompleted
		})
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
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

		blueprintSpec := &domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseEcosystemUnhealthyAfterwards,
		}
		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprintSpec, nil)
		applyMock.EXPECT().PostProcessBlueprintApplication(testCtx, "testBlueprint1").Return(nil).Run(func(ctx context.Context, blueprintId string) {
			blueprintSpec.Status = domain.StatusPhaseCompleted
		})
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
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseCompleted,
		}, nil)
		// when
		err := useCase.HandleChange(testCtx, "testBlueprint1")
		// then
		require.NoError(t, err)
	})

	t.Run("handle blueprint application failed", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(&domain.BlueprintSpec{
			Id:     blueprintId,
			Status: domain.StatusPhaseBlueprintApplicationFailed,
		}, nil)
		applyMock.EXPECT().PostProcessBlueprintApplication(testCtx, blueprintId).Return(nil)
		// when
		err := useCase.HandleChange(testCtx, blueprintId)
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
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

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
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

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

func TestBlueprintSpecChangeUseCase_preProcessBlueprintApplication(t *testing.T) {
	t.Run("stop on dry run", func(t *testing.T) {
		// given
		spec := &domain.BlueprintSpec{
			Id:     blueprintId,
			Config: domain.BlueprintConfiguration{DryRun: true},
		}
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		applyMock.EXPECT().PreProcessBlueprintApplication(testCtx, blueprintId).Return(nil)
		useCase := NewBlueprintSpecChangeUseCase(nil, nil, nil, nil, applyMock, nil, nil, nil)
		// when
		err := useCase.preProcessBlueprintApplication(testCtx, spec)
		// then
		require.NoError(t, err)
	})
	t.Run("error", func(t *testing.T) {
		// given
		spec := &domain.BlueprintSpec{
			Id: blueprintId,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		applyMock.EXPECT().PreProcessBlueprintApplication(testCtx, blueprintId).Return(assert.AnError)
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)
		// when
		err := useCase.preProcessBlueprintApplication(testCtx, spec)
		// then
		require.ErrorIs(t, err, assert.AnError)
	})
}

func TestBlueprintSpecChangeUseCase_applyBlueprintSpec(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		// given
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		applyMock.EXPECT().ApplyBlueprintSpec(testCtx, blueprintId).Return(assert.AnError)
		useCase := NewBlueprintSpecChangeUseCase(nil, nil, nil, nil, applyMock, nil, nil, nil)
		// when
		err := useCase.applyBlueprintSpec(testCtx, blueprintId)
		// then
		require.ErrorIs(t, err, assert.AnError)
	})
}

func TestBlueprintSpecChangeUseCase_checkEcosystemHealthAfterwards(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		// given
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		applyMock.EXPECT().CheckEcosystemHealthAfterwards(testCtx, blueprintId).Return(assert.AnError)
		useCase := NewBlueprintSpecChangeUseCase(nil, nil, nil, nil, applyMock, nil, nil, nil)
		// when
		err := useCase.checkEcosystemHealthAfterwards(testCtx, blueprintId)
		// then
		require.ErrorIs(t, err, assert.AnError)
	})
}

func TestBlueprintSpecChangeUseCase_triggerDoguRestarts(t *testing.T) {
	t.Run("handle error in TriggerDoguRestarts", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		registryConfigUseCaseMock := newMockRegistryConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock, applyMock, registryConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase)

		blueprintSpec := &domain.BlueprintSpec{
			Id:     testBlueprintId,
			Status: domain.StatusPhaseBlueprintApplied,
		}
		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprintSpec, nil)
		doguRestartUseCaseMock.EXPECT().TriggerDoguRestarts(testCtx, testBlueprintId).Return(errors.New("testerror"))

		// when
		err := useCase.HandleChange(testCtx, testBlueprintId)
		// then
		require.Error(t, err)
		assert.Equal(t, "testerror", err.Error())
	})
}
