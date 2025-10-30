package application

import (
	"context"
	"testing"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	testCtx           = context.Background()
	testBlueprintId   = "testBlueprint1"
	testBlueprintSpec = &domain.BlueprintSpec{
		Id:        testBlueprintId,
		StateDiff: domain.StateDiff{DoguDiffs: domain.DoguDiffs{{NeededActions: []domain.Action{domain.ActionInstall}}}},
	}
	testBlueprintSpecEmptyDiff = &domain.BlueprintSpec{
		Id: testBlueprintId,
	}
	testStoppedBlueprintSpec = &domain.BlueprintSpec{
		Id:     testBlueprintId,
		Config: domain.BlueprintConfiguration{Stopped: true},
	}
	testCompletedBlueprintSpec = &domain.BlueprintSpec{
		Id: testBlueprintId,
		Conditions: []domain.Condition{{
			Type:   domain.ConditionCompleted,
			Status: metav1.ConditionTrue,
		}},
	}
)

func TestNewBlueprintSpecChangeUseCase(t *testing.T) {
	// given
	mocks := createAllMocks(t)
	preparationUseCases := NewBlueprintPreparationUseCase(
		mocks.initialStatus,
		mocks.validation,
		mocks.effectiveBlueprint,
		mocks.stateDiff,
		mocks.ecosystemHealth,
		mocks.restoreInProgress,
	)
	applyUseCases := NewBlueprintApplyUseCase(
		mocks.completeBlueprint,
		mocks.ecosystemConfig,
		mocks.applyDogus,
		mocks.ecosystemHealth,
		mocks.dogusUpToDate,
	)

	// when
	result := NewBlueprintSpecChangeUseCase(mocks.repo, preparationUseCases, applyUseCases)

	// then
	require.NotNil(t, result)
	assert.Equal(t, mocks.repo, result.repo)
	assertPreparationUseCases(t, result.preparationUseCase, mocks)
	assertApplyUseCases(t, result.applyUseCase, mocks)
}

type allMocks struct {
	repo               *mockBlueprintSpecRepository
	initialStatus      *mockInitialBlueprintStatusUseCase
	validation         *mockBlueprintSpecValidationUseCase
	effectiveBlueprint *mockEffectiveBlueprintUseCase
	stateDiff          *mockStateDiffUseCase
	completeBlueprint  *mockCompleteBlueprintUseCase
	ecosystemConfig    *mockEcosystemConfigUseCase
	applyDogus         *mockApplyDogusUseCase
	ecosystemHealth    *mockEcosystemHealthUseCase
	dogusUpToDate      *mockDogusUpToDateUseCase
	restoreInProgress  *mockRestoreInProgressUseCase
}

func createAllMocks(t *testing.T) *allMocks {
	return &allMocks{
		repo:               newMockBlueprintSpecRepository(t),
		initialStatus:      newMockInitialBlueprintStatusUseCase(t),
		validation:         newMockBlueprintSpecValidationUseCase(t),
		effectiveBlueprint: newMockEffectiveBlueprintUseCase(t),
		stateDiff:          newMockStateDiffUseCase(t),
		completeBlueprint:  newMockCompleteBlueprintUseCase(t),
		ecosystemConfig:    newMockEcosystemConfigUseCase(t),
		applyDogus:         newMockApplyDogusUseCase(t),
		ecosystemHealth:    newMockEcosystemHealthUseCase(t),
		dogusUpToDate:      newMockDogusUpToDateUseCase(t),
		restoreInProgress:  newMockRestoreInProgressUseCase(t),
	}
}

func assertPreparationUseCases(t *testing.T, useCases BlueprintPreparationUseCase, mocks *allMocks) {
	assert.Equal(t, mocks.validation, useCases.validation)
	assert.Equal(t, mocks.effectiveBlueprint, useCases.effectiveBlueprint)
	assert.Equal(t, mocks.stateDiff, useCases.stateDiff)
	assert.Equal(t, mocks.ecosystemHealth, useCases.healthUseCase)
	assert.Equal(t, mocks.restoreInProgress, useCases.restoreInProgressUseCase)
}

func assertApplyUseCases(t *testing.T, useCases BlueprintApplyUseCase, mocks *allMocks) {
	assert.Equal(t, mocks.ecosystemConfig, useCases.ecosystemConfigUseCase)
	assert.Equal(t, mocks.applyDogus, useCases.applyDogusUseCase)
	assert.Equal(t, mocks.completeBlueprint, useCases.completeUseCase)
	assert.Equal(t, mocks.ecosystemHealth, useCases.healthUseCase)
	assert.Equal(t, mocks.dogusUpToDate, useCases.dogusUpToDateUseCase)
}

func TestBlueprintSpecChangeUseCase_HandleUntilApplied_RepositoryErrors(t *testing.T) {
	t.Run("should return error on error getting blueprint by id", func(t *testing.T) {
		// given
		mocks := createAllMocks(t)
		mocks.repo.EXPECT().GetById(mock.Anything, testBlueprintId).Return(nil, assert.AnError).Run(func(ctx context.Context, blueprintId string) {
			logger, err := logr.FromContext(ctx)
			require.NoError(t, err)
			assert.NotNil(t, logger)
		})

		useCase := createUseCase(mocks)

		// when
		err := useCase.HandleUntilApplied(testCtx, testBlueprintId)

		// then
		assert.ErrorContains(t, err, "cannot load blueprint spec")
	})
}

func TestBlueprintSpecChangeUseCase_HandleUntilApplied_PreparationPhaseErrors(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func(*allMocks)
		wantErrTest func(*testing.T, error)
	}{
		{
			name: "should return error on error initially setting the blueprint status",
			setupMocks: func(mocks *allMocks) {
				mocks.repo.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
				mocks.initialStatus.EXPECT().InitateConditions(mock.Anything, testBlueprintSpec).Return(assert.AnError)
			},
			wantErrTest: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should return error on error validating blueprint statically",
			setupMocks: func(mocks *allMocks) {
				mocks.repo.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
				mocks.initialStatus.EXPECT().InitateConditions(mock.Anything, mock.Anything).Return(nil)
				mocks.validation.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, testBlueprintSpec).Return(assert.AnError)
			},
			wantErrTest: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should return error on error calculate effective blueprint",
			setupMocks: func(mocks *allMocks) {
				mocks.repo.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
				mocks.initialStatus.EXPECT().InitateConditions(mock.Anything, mock.Anything).Return(nil)
				mocks.validation.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, mock.Anything).Return(nil)
				mocks.effectiveBlueprint.EXPECT().CalculateEffectiveBlueprint(mock.Anything, testBlueprintSpec).Return(assert.AnError)
			},
			wantErrTest: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should return error on error validating blueprint dynamically",
			setupMocks: func(mocks *allMocks) {
				mocks.repo.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
				mocks.initialStatus.EXPECT().InitateConditions(mock.Anything, mock.Anything).Return(nil)
				mocks.validation.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, mock.Anything).Return(nil)
				mocks.effectiveBlueprint.EXPECT().CalculateEffectiveBlueprint(mock.Anything, testBlueprintSpec).Return(nil)
				mocks.validation.EXPECT().ValidateBlueprintSpecDynamically(mock.Anything, testBlueprintSpec).Return(assert.AnError)
			},
			wantErrTest: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should return error on error checking ecosystem health",
			setupMocks: func(mocks *allMocks) {
				mocks.repo.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
				mocks.initialStatus.EXPECT().InitateConditions(mock.Anything, mock.Anything).Return(nil)
				mocks.validation.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, mock.Anything).Return(nil)
				mocks.effectiveBlueprint.EXPECT().CalculateEffectiveBlueprint(mock.Anything, testBlueprintSpec).Return(nil)
				mocks.validation.EXPECT().ValidateBlueprintSpecDynamically(mock.Anything, testBlueprintSpec).Return(nil)
				mocks.ecosystemHealth.EXPECT().CheckEcosystemHealth(mock.Anything, testBlueprintSpec).Return(ecosystem.HealthResult{}, assert.AnError)
			},
			wantErrTest: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should return error on error determining state diff",
			setupMocks: func(mocks *allMocks) {
				mocks.repo.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
				mocks.initialStatus.EXPECT().InitateConditions(mock.Anything, mock.Anything).Return(nil)
				mocks.validation.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, mock.Anything).Return(nil)
				mocks.effectiveBlueprint.EXPECT().CalculateEffectiveBlueprint(mock.Anything, testBlueprintSpec).Return(nil)
				mocks.validation.EXPECT().ValidateBlueprintSpecDynamically(mock.Anything, testBlueprintSpec).Return(nil)
				mocks.ecosystemHealth.EXPECT().CheckEcosystemHealth(mock.Anything, testBlueprintSpec).Return(ecosystem.HealthResult{}, nil)
				mocks.stateDiff.EXPECT().DetermineStateDiff(mock.Anything, testBlueprintSpec).Return(assert.AnError)
			},
			wantErrTest: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should return error on error checking if a restore is in progress",
			setupMocks: func(mocks *allMocks) {
				mocks.repo.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
				mocks.initialStatus.EXPECT().InitateConditions(mock.Anything, mock.Anything).Return(nil)
				mocks.validation.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, mock.Anything).Return(nil)
				mocks.effectiveBlueprint.EXPECT().CalculateEffectiveBlueprint(mock.Anything, testBlueprintSpec).Return(nil)
				mocks.validation.EXPECT().ValidateBlueprintSpecDynamically(mock.Anything, testBlueprintSpec).Return(nil)
				mocks.ecosystemHealth.EXPECT().CheckEcosystemHealth(mock.Anything, testBlueprintSpec).Return(ecosystem.HealthResult{}, nil)
				mocks.stateDiff.EXPECT().DetermineStateDiff(mock.Anything, testBlueprintSpec).Return(nil)
				mocks.restoreInProgress.EXPECT().CheckRestoreInProgress(mock.Anything).Return(assert.AnError)
			},
			wantErrTest: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks := createAllMocks(t)
			tt.setupMocks(mocks)
			useCase := createUseCase(mocks)

			err := useCase.HandleUntilApplied(testCtx, testBlueprintId)
			tt.wantErrTest(t, err)
		})
	}
}

func TestBlueprintSpecChangeUseCase_HandleUntilApplied_SpecialBlueprintStates(t *testing.T) {
	t.Run("should handle stopped blueprint", func(t *testing.T) {
		// given
		mocks := createAllMocks(t)
		mocks.repo.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testStoppedBlueprintSpec, nil)
		setupSuccessfulPreparationPhase(mocks, testStoppedBlueprintSpec)
		mocks.repo.EXPECT().Update(mock.Anything, mock.Anything).Run(func(ctx context.Context, blueprint *domain.BlueprintSpec) {
			assert.Len(t, blueprint.Events, 1)
			assert.Equal(t, domain.BlueprintStoppedEvent{}.Name(), blueprint.Events[0].Name())
		}).Return(nil)

		useCase := createUseCase(mocks)

		// when
		err := useCase.HandleUntilApplied(testCtx, testBlueprintId)

		// then
		assert.NoError(t, err)
	})

	t.Run("should not apply completed blueprint with no diff", func(t *testing.T) {
		// given
		mocks := createAllMocks(t)
		mocks.repo.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testCompletedBlueprintSpec, nil)
		setupSuccessfulPreparationPhase(mocks, testCompletedBlueprintSpec)

		useCase := createUseCase(mocks)

		// when
		err := useCase.HandleUntilApplied(testCtx, testBlueprintId)

		// then
		assert.NoError(t, err)
	})
}

func TestBlueprintSpecChangeUseCase_HandleUntilApplied_ApplyPhaseErrors(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func(*allMocks)
		wantErrTest func(*testing.T, error)
	}{
		{
			name: "should return error on error apply config",
			setupMocks: func(mocks *allMocks) {
				mocks.repo.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
				setupSuccessfulPreparationPhase(mocks, testBlueprintSpec)
				mocks.ecosystemConfig.EXPECT().ApplyConfig(mock.Anything, testBlueprintSpec).Return(assert.AnError)
			},
			wantErrTest: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should return error on error apply dogus",
			setupMocks: func(mocks *allMocks) {
				mocks.repo.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
				setupSuccessfulPreparationPhase(mocks, testBlueprintSpec)
				mocks.ecosystemConfig.EXPECT().ApplyConfig(mock.Anything, testBlueprintSpec).Return(nil)
				mocks.applyDogus.EXPECT().ApplyDogus(mock.Anything, testBlueprintSpec).Return(false, assert.AnError)
			},
			wantErrTest: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should return error on error check health after dogu apply",
			setupMocks: func(mocks *allMocks) {
				mocks.repo.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
				setupSuccessfulPreparationPhase(mocks, testBlueprintSpec)
				mocks.ecosystemConfig.EXPECT().ApplyConfig(mock.Anything, testBlueprintSpec).Return(nil)
				mocks.applyDogus.EXPECT().ApplyDogus(mock.Anything, testBlueprintSpec).Return(true, nil)
				mocks.ecosystemHealth.EXPECT().CheckEcosystemHealth(mock.Anything, testBlueprintSpec).Return(ecosystem.HealthResult{}, assert.AnError)
			},
			wantErrTest: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should return error on error checking if dogus are up to date",
			setupMocks: func(mocks *allMocks) {
				mocks.repo.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
				setupSuccessfulPreparationPhase(mocks, testBlueprintSpec)
				mocks.ecosystemConfig.EXPECT().ApplyConfig(mock.Anything, testBlueprintSpec).Return(nil)
				mocks.applyDogus.EXPECT().ApplyDogus(mock.Anything, testBlueprintSpec).Return(false, nil)
				mocks.dogusUpToDate.EXPECT().CheckDogus(mock.Anything, testBlueprintSpec).Return(assert.AnError)
			},
			wantErrTest: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks := createAllMocks(t)
			tt.setupMocks(mocks)
			useCase := createUseCase(mocks)

			err := useCase.HandleUntilApplied(testCtx, testBlueprintId)
			tt.wantErrTest(t, err)
		})
	}
}

func TestBlueprintSpecChangeUseCase_HandleUntilApplied_CompletionScenarios(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func(*allMocks)
		wantErrTest func(*testing.T, error)
	}{
		{
			name: "should return nil on complete success",
			setupMocks: func(mocks *allMocks) {
				mocks.repo.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpecEmptyDiff, nil)
				setupSuccessfulPreparationPhase(mocks, testBlueprintSpecEmptyDiff)
				setupSuccessfulApplyPhaseExceptComplete(mocks, testBlueprintSpecEmptyDiff)
				mocks.completeBlueprint.EXPECT().CompleteBlueprint(mock.Anything, testBlueprintSpecEmptyDiff).Return(nil)
			},
			wantErrTest: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "should return error on error complete blueprint",
			setupMocks: func(mocks *allMocks) {
				mocks.repo.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpecEmptyDiff, nil)
				setupSuccessfulPreparationPhase(mocks, testBlueprintSpecEmptyDiff)
				setupSuccessfulApplyPhaseExceptComplete(mocks, testBlueprintSpecEmptyDiff)
				mocks.completeBlueprint.EXPECT().CompleteBlueprint(mock.Anything, testBlueprintSpecEmptyDiff).Return(assert.AnError)
			},
			wantErrTest: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should return StateDiffNotEmptyError when diff not empty",
			setupMocks: func(mocks *allMocks) {
				mocks.repo.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
				setupSuccessfulPreparationPhase(mocks, testBlueprintSpec)
				setupSuccessfulApplyPhaseExceptComplete(mocks, testBlueprintSpec)
			},
			wantErrTest: func(t *testing.T, err error) {
				var targetError *domain.StateDiffNotEmptyError
				assert.Error(t, err)
				assert.ErrorAs(t, err, &targetError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks := createAllMocks(t)
			tt.setupMocks(mocks)
			useCase := createUseCase(mocks)

			err := useCase.HandleUntilApplied(testCtx, testBlueprintId)
			tt.wantErrTest(t, err)
		})
	}
}

func createUseCase(mocks *allMocks) *BlueprintSpecChangeUseCase {
	preparationUseCases := BlueprintPreparationUseCase{
		initialStatus:            mocks.initialStatus,
		validation:               mocks.validation,
		effectiveBlueprint:       mocks.effectiveBlueprint,
		stateDiff:                mocks.stateDiff,
		healthUseCase:            mocks.ecosystemHealth,
		restoreInProgressUseCase: mocks.restoreInProgress,
	}
	applyUseCases := BlueprintApplyUseCase{
		completeUseCase:        mocks.completeBlueprint,
		ecosystemConfigUseCase: mocks.ecosystemConfig,
		applyDogusUseCase:      mocks.applyDogus,
		healthUseCase:          mocks.ecosystemHealth,
		dogusUpToDateUseCase:   mocks.dogusUpToDate,
	}

	return &BlueprintSpecChangeUseCase{
		repo:               mocks.repo,
		preparationUseCase: preparationUseCases,
		applyUseCase:       applyUseCases,
	}
}

func setupSuccessfulPreparationPhase(mocks *allMocks, spec *domain.BlueprintSpec) {
	mocks.initialStatus.EXPECT().InitateConditions(mock.Anything, mock.Anything).Return(nil)
	mocks.validation.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, mock.Anything).Return(nil)
	mocks.effectiveBlueprint.EXPECT().CalculateEffectiveBlueprint(mock.Anything, spec).Return(nil)
	mocks.validation.EXPECT().ValidateBlueprintSpecDynamically(mock.Anything, spec).Return(nil)
	mocks.ecosystemHealth.EXPECT().CheckEcosystemHealth(mock.Anything, spec).Return(ecosystem.HealthResult{}, nil).Times(1)
	mocks.stateDiff.EXPECT().DetermineStateDiff(mock.Anything, spec).Return(nil)
	mocks.restoreInProgress.EXPECT().CheckRestoreInProgress(mock.Anything).Return(nil)
}

func setupSuccessfulApplyPhaseExceptComplete(mocks *allMocks, spec *domain.BlueprintSpec) {
	mocks.ecosystemConfig.EXPECT().ApplyConfig(mock.Anything, spec).Return(nil)
	mocks.applyDogus.EXPECT().ApplyDogus(mock.Anything, spec).Return(false, nil)
	mocks.dogusUpToDate.EXPECT().CheckDogus(mock.Anything, spec).Return(nil)
	// Note: no completeBlueprint expectation - this allows the completion steps to be tested
}

func TestBlueprintSpecChangeUseCase_CheckForMultipleBlueprintResources(t *testing.T) {
	t.Run("should succeed without error", func(t *testing.T) {
		// given
		mockRepo := newMockBlueprintSpecRepository(t)
		mockRepo.EXPECT().Count(t.Context(), 2).Return(1, nil)
		useCase := &BlueprintSpecChangeUseCase{
			repo: mockRepo,
		}

		//when
		err := useCase.CheckForMultipleBlueprintResources(t.Context())

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on repo error", func(t *testing.T) {
		// given
		mockRepo := newMockBlueprintSpecRepository(t)
		mockRepo.EXPECT().Count(t.Context(), 2).Return(0, assert.AnError)
		useCase := &BlueprintSpecChangeUseCase{
			repo: mockRepo,
		}

		//when
		err := useCase.CheckForMultipleBlueprintResources(t.Context())

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "check for multiple blueprints not successful")
	})

	t.Run("should return error on multiple blueprints", func(t *testing.T) {
		// given
		mockRepo := newMockBlueprintSpecRepository(t)
		mockRepo.EXPECT().Count(t.Context(), 2).Return(2, nil)
		useCase := &BlueprintSpecChangeUseCase{
			repo: mockRepo,
		}

		//when
		err := useCase.CheckForMultipleBlueprintResources(t.Context())

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "check for multiple blueprints not successful: more than one blueprint CR found")
	})
}
