package application

import (
	"context"
	"errors"
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
)

var testCtx = context.Background()
var testBlueprintId = "testBlueprint1"

func TestBlueprintSpecChangeUseCase_HandleChange(t *testing.T) {
	//FIXME: this test runs endless at the moment. Refactoring the changeUseCase is the last step
	require.True(t, false, "tests deactivated until the refactoring is done.")

	t.Run("do all steps with blueprint", func(t *testing.T) {
		// given
		logger := log.FromContext(testCtx).
			WithName("BlueprintSpecChangeUseCase.HandleUntilApplied").
			WithValues("blueprintId", blueprintId)
		log.SetLogger(logger)
		ctxWithLogger := log.IntoContext(testCtx, logger)

		repoMock := newMockBlueprintSpecRepository(t)
		validationMock := newMockBlueprintSpecValidationUseCase(t)
		effectiveBlueprintMock := newMockEffectiveBlueprintUseCase(t)
		stateDiffMock := newMockStateDiffUseCase(t)
		applyMock := newMockApplyBlueprintSpecUseCase(t)
		ecosystemConfigUseCaseMock := newMockEcosystemConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		componentUseCase := newMockApplyComponentUseCase(t)
		doguUseCase := newMockApplyDogusUseCase(t)
		healthUseCase := newMockEcosystemHealthUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock,
			applyMock, ecosystemConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase, componentUseCase, doguUseCase, healthUseCase)

		blueprintSpec := &domain.BlueprintSpec{
			Id:     testBlueprintId,
			Status: domain.StatusPhaseNew,
		}
		repoMock.EXPECT().GetById(mock.Anything, testBlueprintId).Return(blueprintSpec, nil)
		validationMock.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, blueprintSpec).Return(nil).
			Run(func(ctx context.Context, blueprint *domain.BlueprintSpec) {
				meta.SetStatusCondition(blueprint.Conditions, metav1.Condition{
					Type:   domain.ConditionValid,
					Status: metav1.ConditionTrue,
				})
			})
		effectiveBlueprintMock.EXPECT().CalculateEffectiveBlueprint(mock.Anything, blueprintSpec).Return(nil)
		validationMock.EXPECT().ValidateBlueprintSpecDynamically(mock.Anything, blueprintSpec).Return(nil).
			Run(func(ctx context.Context, blueprint *domain.BlueprintSpec) {
				meta.SetStatusCondition(blueprint.Conditions, metav1.Condition{
					Type:   domain.ConditionValid,
					Status: metav1.ConditionTrue,
				})
			})
		stateDiffMock.EXPECT().DetermineStateDiff(mock.Anything, blueprintSpec).Return(nil).
			Run(func(ctx context.Context, blueprint *domain.BlueprintSpec) {
				meta.SetStatusCondition(blueprint.Conditions, metav1.Condition{
					Type:   domain.ConditionExecutable,
					Status: metav1.ConditionTrue,
				})
			})
		healthUseCase.EXPECT().CheckEcosystemHealth(mock.Anything, blueprintSpec).Return(ecosystem.HealthResult{}, nil)
		ecosystemConfigUseCaseMock.EXPECT().ApplyConfig(mock.Anything, blueprintSpec).Return(nil)
		selfUpgradeUseCase.EXPECT().HandleSelfUpgrade(mock.Anything, blueprintSpec).Return(nil)
		healthUseCase.EXPECT().CheckEcosystemHealth(mock.Anything, blueprintSpec).Return(ecosystem.HealthResult{}, nil)
		applyMock.EXPECT().PostProcessBlueprintApplication(mock.Anything, blueprintSpec).Return(nil).
			Run(func(ctx context.Context, blueprint *domain.BlueprintSpec) {
				blueprint.Status = domain.StatusPhaseCompleted
			})
		doguRestartUseCaseMock.EXPECT().TriggerDoguRestarts(mock.Anything, blueprintSpec).Return(nil).
			Run(func(ctx context.Context, blueprint *domain.BlueprintSpec) {
				blueprint.Status = domain.StatusPhaseRestartsTriggered
			})

		// when
		err := useCase.HandleUntilApplied(ctxWithLogger, testBlueprintId)
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
		ecosystemConfigUseCaseMock := newMockEcosystemConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		componentUseCase := newMockApplyComponentUseCase(t)
		doguUseCase := newMockApplyDogusUseCase(t)
		healthUseCase := newMockEcosystemHealthUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock,
			applyMock, ecosystemConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase, componentUseCase, doguUseCase, healthUseCase)

		expectedError := &domainservice.InternalError{
			WrappedError: nil,
			Message:      "test-error",
		}
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(nil, expectedError)

		// when
		err := useCase.HandleUntilApplied(testCtx, blueprintId)

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
		ecosystemConfigUseCaseMock := newMockEcosystemConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		componentUseCase := newMockApplyComponentUseCase(t)
		doguUseCase := newMockApplyDogusUseCase(t)
		healthUseCase := newMockEcosystemHealthUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock,
			applyMock, ecosystemConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase, componentUseCase, doguUseCase, healthUseCase)

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseNew,
			Blueprint: domain.Blueprint{Dogus: []domain.Dogu{
				{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "DoguWithNoVersion"}, TargetState: domain.TargetStatePresent},
			}},
		}, nil)
		validationMock.EXPECT().ValidateBlueprintSpecStatically(testCtx, "testBlueprint1").Return(assert.AnError)

		// when
		err := useCase.HandleUntilApplied(testCtx, "testBlueprint1")

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
		ecosystemConfigUseCaseMock := newMockEcosystemConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		componentUseCase := newMockApplyComponentUseCase(t)
		doguUseCase := newMockApplyDogusUseCase(t)
		healthUseCase := newMockEcosystemHealthUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock,
			applyMock, ecosystemConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase, componentUseCase, doguUseCase, healthUseCase)

		updatedSpec := &domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseNew,
		}

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(updatedSpec, nil)
		validationMock.EXPECT().ValidateBlueprintSpecStatically(testCtx, "testBlueprint1").Return(nil)
		effectiveBlueprintMock.EXPECT().CalculateEffectiveBlueprint(testCtx, "testBlueprint1").Return(assert.AnError)

		// when
		err := useCase.HandleUntilApplied(testCtx, "testBlueprint1")

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
		ecosystemConfigUseCaseMock := newMockEcosystemConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		componentUseCase := newMockApplyComponentUseCase(t)
		doguUseCase := newMockApplyDogusUseCase(t)
		healthUseCase := newMockEcosystemHealthUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock,
			applyMock, ecosystemConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase, componentUseCase, doguUseCase, healthUseCase)

		updatedSpec := &domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseNew,
		}

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(updatedSpec, nil)
		validationMock.EXPECT().ValidateBlueprintSpecStatically(testCtx, "testBlueprint1").Return(nil)
		effectiveBlueprintMock.EXPECT().CalculateEffectiveBlueprint(testCtx, "testBlueprint1").Return(nil)
		validationMock.EXPECT().ValidateBlueprintSpecDynamically(testCtx, "testBlueprint1").Return(assert.AnError)

		// when
		err := useCase.HandleUntilApplied(testCtx, "testBlueprint1")

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
		ecosystemConfigUseCaseMock := newMockEcosystemConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		componentUseCase := newMockApplyComponentUseCase(t)
		doguUseCase := newMockApplyDogusUseCase(t)
		healthUseCase := newMockEcosystemHealthUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock,
			applyMock, ecosystemConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase, componentUseCase, doguUseCase, healthUseCase)

		updatedSpec := &domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseNew,
		}

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(updatedSpec, nil)
		validationMock.EXPECT().ValidateBlueprintSpecStatically(testCtx, "testBlueprint1").Return(nil)
		effectiveBlueprintMock.EXPECT().CalculateEffectiveBlueprint(testCtx, "testBlueprint1").Return(nil)
		validationMock.EXPECT().ValidateBlueprintSpecDynamically(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprint *domain.BlueprintSpec) {
				meta.SetStatusCondition(blueprint.Conditions, metav1.Condition{
					Type:   domain.ConditionValid,
					Status: metav1.ConditionTrue,
				})
			})
		stateDiffMock.EXPECT().DetermineStateDiff(testCtx, "testBlueprint1").Return(assert.AnError)
		// when
		err := useCase.HandleUntilApplied(testCtx, "testBlueprint1")

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
		ecosystemConfigUseCaseMock := newMockEcosystemConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		componentUseCase := newMockApplyComponentUseCase(t)
		doguUseCase := newMockApplyDogusUseCase(t)
		healthUseCase := newMockEcosystemHealthUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock,
			applyMock, ecosystemConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase, componentUseCase, doguUseCase, healthUseCase)

		updatedSpec := &domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseNew,
		}

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(updatedSpec, nil)
		validationMock.EXPECT().ValidateBlueprintSpecStatically(testCtx, "testBlueprint1").Return(nil)
		effectiveBlueprintMock.EXPECT().CalculateEffectiveBlueprint(testCtx, "testBlueprint1").Return(nil)
		validationMock.EXPECT().ValidateBlueprintSpecDynamically(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprint *domain.BlueprintSpec) {
				meta.SetStatusCondition(blueprint.Conditions, metav1.Condition{
					Type:   domain.ConditionValid,
					Status: metav1.ConditionTrue,
				})
			})
		stateDiffMock.EXPECT().DetermineStateDiff(testCtx, "testBlueprint1").Return(nil).
			Run(func(ctx context.Context, blueprint *domain.BlueprintSpec) {
				meta.SetStatusCondition(blueprint.Conditions, metav1.Condition{
					Type:   domain.ConditionExecutable,
					Status: metav1.ConditionTrue,
				})
			})
		healthUseCase.EXPECT().CheckEcosystemHealth(mock.Anything, updatedSpec).Return(ecosystem.HealthResult{}, assert.AnError)

		// when
		err := useCase.HandleUntilApplied(testCtx, "testBlueprint1")

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
		ecosystemConfigUseCaseMock := newMockEcosystemConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		componentUseCase := newMockApplyComponentUseCase(t)
		doguUseCase := newMockApplyDogusUseCase(t)
		healthUseCase := newMockEcosystemHealthUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock,
			applyMock, ecosystemConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase, componentUseCase, doguUseCase, healthUseCase)

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id: "testBlueprint1",
		}, nil)
		// when
		err := useCase.HandleUntilApplied(testCtx, "testBlueprint1")
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
		ecosystemConfigUseCaseMock := newMockEcosystemConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		componentUseCase := newMockApplyComponentUseCase(t)
		doguUseCase := newMockApplyDogusUseCase(t)
		healthUseCase := newMockEcosystemHealthUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock,
			applyMock, ecosystemConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase, componentUseCase, doguUseCase, healthUseCase)

		blueprintSpec := &domain.BlueprintSpec{
			Id:         "testBlueprint1",
			Conditions: &[]domain.Condition{},
		}
		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprintSpec, nil)
		ecosystemConfigUseCaseMock.EXPECT().ApplyConfig(testCtx, blueprintSpec.Id).Return(assert.AnError)
		// when
		actualErr := useCase.HandleUntilApplied(testCtx, "testBlueprint1")
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
		ecosystemConfigUseCaseMock := newMockEcosystemConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		componentUseCase := newMockApplyComponentUseCase(t)
		doguUseCase := newMockApplyDogusUseCase(t)
		healthUseCase := newMockEcosystemHealthUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock,
			applyMock, ecosystemConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase, componentUseCase, doguUseCase, healthUseCase)

		blueprintSpec := &domain.BlueprintSpec{
			Id:     testBlueprintId,
			Status: domain.StatusPhaseBlueprintApplied,
		}
		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprintSpec, nil)
		healthUseCase.EXPECT().CheckEcosystemHealth(mock.Anything, blueprintSpec).Return(ecosystem.HealthResult{}, assert.AnError)
		doguRestartUseCaseMock.EXPECT().TriggerDoguRestarts(testCtx, testBlueprintId).Return(nil).
			Run(func(ctx context.Context, blueprint *domain.BlueprintSpec) {
				blueprint.Status = domain.StatusPhaseRestartsTriggered
			})

		// when
		err := useCase.HandleUntilApplied(testCtx, testBlueprintId)
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
		ecosystemConfigUseCaseMock := newMockEcosystemConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		componentUseCase := newMockApplyComponentUseCase(t)
		doguUseCase := newMockApplyDogusUseCase(t)
		healthUseCase := newMockEcosystemHealthUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock,
			applyMock, ecosystemConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase, componentUseCase, doguUseCase, healthUseCase)

		blueprintSpec := &domain.BlueprintSpec{
			Id: "testBlueprint1",
		}
		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprintSpec, nil)
		applyMock.EXPECT().PostProcessBlueprintApplication(testCtx, "testBlueprint1").Return(nil).Run(func(ctx context.Context, blueprint *domain.BlueprintSpec) {
			blueprint.Status = domain.StatusPhaseCompleted
		})
		// when
		err := useCase.HandleUntilApplied(testCtx, "testBlueprint1")
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
		ecosystemConfigUseCaseMock := newMockEcosystemConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		componentUseCase := newMockApplyComponentUseCase(t)
		doguUseCase := newMockApplyDogusUseCase(t)
		healthUseCase := newMockEcosystemHealthUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock,
			applyMock, ecosystemConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase, componentUseCase, doguUseCase, healthUseCase)

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseCompleted,
		}, nil)
		// when
		err := useCase.HandleUntilApplied(testCtx, "testBlueprint1")
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
		ecosystemConfigUseCaseMock := newMockEcosystemConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		componentUseCase := newMockApplyComponentUseCase(t)
		doguUseCase := newMockApplyDogusUseCase(t)
		healthUseCase := newMockEcosystemHealthUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock,
			applyMock, ecosystemConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase, componentUseCase, doguUseCase, healthUseCase)

		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(&domain.BlueprintSpec{
			Id:     blueprintId,
			Status: domain.StatusPhaseBlueprintApplicationFailed,
		}, nil)
		applyMock.EXPECT().PostProcessBlueprintApplication(testCtx, blueprintId).Return(nil)
		// when
		err := useCase.HandleUntilApplied(testCtx, blueprintId)
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
		ecosystemConfigUseCaseMock := newMockEcosystemConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		componentUseCase := newMockApplyComponentUseCase(t)
		doguUseCase := newMockApplyDogusUseCase(t)
		healthUseCase := newMockEcosystemHealthUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock,
			applyMock, ecosystemConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase, componentUseCase, doguUseCase, healthUseCase)

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseFailed,
		}, nil)
		// when
		err := useCase.HandleUntilApplied(testCtx, "testBlueprint1")
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
		ecosystemConfigUseCaseMock := newMockEcosystemConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		componentUseCase := newMockApplyComponentUseCase(t)
		doguUseCase := newMockApplyDogusUseCase(t)
		healthUseCase := newMockEcosystemHealthUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock,
			applyMock, ecosystemConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase, componentUseCase, doguUseCase, healthUseCase)

		repoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: "unknown",
		}, nil)
		// when
		err := useCase.HandleUntilApplied(testCtx, "testBlueprint1")
		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "could not handle unknown status of blueprint")
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
		ecosystemConfigUseCaseMock := newMockEcosystemConfigUseCase(t)
		doguRestartUseCaseMock := newMockDoguRestartUseCase(t)
		selfUpgradeUseCase := newMockSelfUpgradeUseCase(t)
		componentUseCase := newMockApplyComponentUseCase(t)
		doguUseCase := newMockApplyDogusUseCase(t)
		healthUseCase := newMockEcosystemHealthUseCase(t)
		useCase := NewBlueprintSpecChangeUseCase(repoMock, validationMock, effectiveBlueprintMock, stateDiffMock,
			applyMock, ecosystemConfigUseCaseMock, doguRestartUseCaseMock, selfUpgradeUseCase, componentUseCase, doguUseCase, healthUseCase)

		blueprintSpec := &domain.BlueprintSpec{
			Id:     testBlueprintId,
			Status: domain.StatusPhaseBlueprintApplied,
		}
		repoMock.EXPECT().GetById(mock.Anything, "testBlueprint1").Return(blueprintSpec, nil)
		doguRestartUseCaseMock.EXPECT().TriggerDoguRestarts(mock.Anything, testBlueprintId).Return(errors.New("testerror"))

		// when
		err := useCase.HandleUntilApplied(testCtx, testBlueprintId)
		// then
		require.Error(t, err)
		assert.Equal(t, "testerror", err.Error())
	})
}
