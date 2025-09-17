package application

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

var testCtx = context.Background()
var testBlueprintId = "testBlueprint1"

func TestNewBlueprintSpecChangeUseCase(t *testing.T) {
	// given
	blueprintSpecRepositoryMock := newMockBlueprintSpecRepository(t)
	validationUseCaseMock := newMockBlueprintSpecValidationUseCase(t)
	effectiveUseCaseMock := newMockEffectiveBlueprintUseCase(t)
	stateDiffUseCaseMock := newMockStateDiffUseCase(t)
	completeUseCaseMock := newMockCompleteBlueprintUseCase(t)
	ecosystemConfigUseCaseMock := newMockEcosystemConfigUseCase(t)
	selfUpgradeUseCaseMock := newMockSelfUpgradeUseCase(t)
	applyComponentUseCaseMock := newMockApplyComponentsUseCase(t)
	applyDoguUseCaseMock := newMockApplyDogusUseCase(t)
	ecosystemHealthUseCaseMock := newMockEcosystemHealthUseCase(t)

	// when
	result := NewBlueprintSpecChangeUseCase(
		blueprintSpecRepositoryMock,
		validationUseCaseMock,
		effectiveUseCaseMock,
		stateDiffUseCaseMock,
		completeUseCaseMock,
		ecosystemConfigUseCaseMock,
		selfUpgradeUseCaseMock,
		applyComponentUseCaseMock,
		applyDoguUseCaseMock,
		ecosystemHealthUseCaseMock,
	)

	// then
	require.NotNil(t, result)
	assert.Equal(t, blueprintSpecRepositoryMock, result.repo)
	assert.Equal(t, validationUseCaseMock, result.validation)
	assert.Equal(t, effectiveUseCaseMock, result.effectiveBlueprint)
	assert.Equal(t, stateDiffUseCaseMock, result.stateDiff)
	assert.Equal(t, ecosystemConfigUseCaseMock, result.ecosystemConfigUseCase)
	assert.Equal(t, selfUpgradeUseCaseMock, result.selfUpgradeUseCase)
	assert.Equal(t, applyComponentUseCaseMock, result.applyComponentUseCase)
	assert.Equal(t, applyDoguUseCaseMock, result.applyDogusUseCase)
	assert.Equal(t, ecosystemHealthUseCaseMock, result.healthUseCase)
	assert.Equal(t, completeUseCaseMock, result.applyUseCase)
}

func TestBlueprintSpecChangeUseCase_HandleUntilApplied(t *testing.T) {
	testBlueprintSpec := &domain.BlueprintSpec{
		Id:        testBlueprintId,
		StateDiff: domain.StateDiff{DoguDiffs: domain.DoguDiffs{{NeededActions: []domain.Action{domain.ActionInstall}}}},
	}
	testDryRunBlueprintSpec := &domain.BlueprintSpec{
		Id:     testBlueprintId,
		Config: domain.BlueprintConfiguration{Stopped: true},
	}

	type fields struct {
		repo                   func(t *testing.T) blueprintSpecRepository
		validation             func(t *testing.T) blueprintSpecValidationUseCase
		effectiveBlueprint     func(t *testing.T) effectiveBlueprintUseCase
		stateDiff              func(t *testing.T) stateDiffUseCase
		applyUseCase           func(t *testing.T) completeBlueprintUseCase
		ecosystemConfigUseCase func(t *testing.T) ecosystemConfigUseCase
		selfUpgradeUseCase     func(t *testing.T) selfUpgradeUseCase
		applyComponentUseCase  func(t *testing.T) applyComponentsUseCase
		applyDogusUseCase      func(t *testing.T) applyDogusUseCase
		healthUseCase          func(t *testing.T) ecosystemHealthUseCase
	}
	type args struct {
		givenCtx    context.Context
		blueprintId string
	}

	testArgs := args{
		givenCtx:    testCtx,
		blueprintId: testBlueprintId,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should return error on error getting blueprint by id",
			fields: fields{
				repo: func(t *testing.T) blueprintSpecRepository {
					m := newMockBlueprintSpecRepository(t)
					m.EXPECT().GetById(mock.Anything, testBlueprintId).Return(nil, assert.AnError).Run(func(ctx context.Context, blueprintId string) {
						logger, err := logr.FromContext(ctx)
						require.NoError(t, err)
						assert.NotNil(t, logger)
					})
					return m
				},
			},
			args: testArgs,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "cannot load blueprint spec")
			},
		},
		{
			name: "should return error on error validating blueprint statically",
			fields: fields{
				repo: func(t *testing.T) blueprintSpecRepository {
					m := newMockBlueprintSpecRepository(t)
					m.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
					return m
				},
				validation: func(t *testing.T) blueprintSpecValidationUseCase {
					m := newMockBlueprintSpecValidationUseCase(t)
					m.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, testBlueprintSpec).Return(assert.AnError)

					return m
				},
			},
			args: testArgs,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err)
			},
		},
		{
			name: "should return error on error calculate effective blueprint",
			fields: fields{
				repo: func(t *testing.T) blueprintSpecRepository {
					m := newMockBlueprintSpecRepository(t)
					m.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
					return m
				},
				validation: func(t *testing.T) blueprintSpecValidationUseCase {
					m := newMockBlueprintSpecValidationUseCase(t)
					m.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, testBlueprintSpec).Return(nil)

					return m
				},
				effectiveBlueprint: func(t *testing.T) effectiveBlueprintUseCase {
					m := newMockEffectiveBlueprintUseCase(t)
					m.EXPECT().CalculateEffectiveBlueprint(mock.Anything, testBlueprintSpec).Return(assert.AnError)
					return m
				},
			},
			args: testArgs,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err)
			},
		},
		{
			name: "should return error on error validating blueprint dynamically",
			fields: fields{
				repo: func(t *testing.T) blueprintSpecRepository {
					m := newMockBlueprintSpecRepository(t)
					m.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
					return m
				},
				validation: func(t *testing.T) blueprintSpecValidationUseCase {
					m := newMockBlueprintSpecValidationUseCase(t)
					m.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, testBlueprintSpec).Return(nil)
					m.EXPECT().ValidateBlueprintSpecDynamically(mock.Anything, testBlueprintSpec).Return(assert.AnError)
					return m
				},
				effectiveBlueprint: func(t *testing.T) effectiveBlueprintUseCase {
					m := newMockEffectiveBlueprintUseCase(t)
					m.EXPECT().CalculateEffectiveBlueprint(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
			},
			args: testArgs,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err)
			},
		},
		{
			name: "should return error on error determining state diff",
			fields: fields{
				repo: func(t *testing.T) blueprintSpecRepository {
					m := newMockBlueprintSpecRepository(t)
					m.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
					return m
				},
				validation: func(t *testing.T) blueprintSpecValidationUseCase {
					m := newMockBlueprintSpecValidationUseCase(t)
					m.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, testBlueprintSpec).Return(nil)
					m.EXPECT().ValidateBlueprintSpecDynamically(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				effectiveBlueprint: func(t *testing.T) effectiveBlueprintUseCase {
					m := newMockEffectiveBlueprintUseCase(t)
					m.EXPECT().CalculateEffectiveBlueprint(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				stateDiff: func(t *testing.T) stateDiffUseCase {
					m := newMockStateDiffUseCase(t)
					m.EXPECT().DetermineStateDiff(mock.Anything, testBlueprintSpec).Return(assert.AnError)
					return m
				},
			},
			args: testArgs,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err)
			},
		},
		{
			name: "should return error on error checking ecosystem health",
			fields: fields{
				repo: func(t *testing.T) blueprintSpecRepository {
					m := newMockBlueprintSpecRepository(t)
					m.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
					return m
				},
				validation: func(t *testing.T) blueprintSpecValidationUseCase {
					m := newMockBlueprintSpecValidationUseCase(t)
					m.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, testBlueprintSpec).Return(nil)
					m.EXPECT().ValidateBlueprintSpecDynamically(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				effectiveBlueprint: func(t *testing.T) effectiveBlueprintUseCase {
					m := newMockEffectiveBlueprintUseCase(t)
					m.EXPECT().CalculateEffectiveBlueprint(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				stateDiff: func(t *testing.T) stateDiffUseCase {
					m := newMockStateDiffUseCase(t)
					m.EXPECT().DetermineStateDiff(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				healthUseCase: func(t *testing.T) ecosystemHealthUseCase {
					m := newMockEcosystemHealthUseCase(t)
					m.EXPECT().CheckEcosystemHealth(mock.Anything, testBlueprintSpec).Return(ecosystem.HealthResult{}, assert.AnError)
					return m
				},
			},
			args: testArgs,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err)
			},
		},
		{
			name: "should return nil and do nothing on dry run",
			fields: fields{
				repo: func(t *testing.T) blueprintSpecRepository {
					m := newMockBlueprintSpecRepository(t)
					m.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testDryRunBlueprintSpec, nil)
					return m
				},
				validation: func(t *testing.T) blueprintSpecValidationUseCase {
					m := newMockBlueprintSpecValidationUseCase(t)
					m.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, testDryRunBlueprintSpec).Return(nil)
					m.EXPECT().ValidateBlueprintSpecDynamically(mock.Anything, testDryRunBlueprintSpec).Return(nil)
					return m
				},
				effectiveBlueprint: func(t *testing.T) effectiveBlueprintUseCase {
					m := newMockEffectiveBlueprintUseCase(t)
					m.EXPECT().CalculateEffectiveBlueprint(mock.Anything, testDryRunBlueprintSpec).Return(nil)
					return m
				},
				stateDiff: func(t *testing.T) stateDiffUseCase {
					m := newMockStateDiffUseCase(t)
					m.EXPECT().DetermineStateDiff(mock.Anything, testDryRunBlueprintSpec).Return(nil)
					return m
				},
				healthUseCase: func(t *testing.T) ecosystemHealthUseCase {
					m := newMockEcosystemHealthUseCase(t)
					m.EXPECT().CheckEcosystemHealth(mock.Anything, testDryRunBlueprintSpec).Return(ecosystem.HealthResult{}, nil)
					return m
				},
			},
			args:    testArgs,
			wantErr: assert.NoError,
		},
		{
			name: "should return error on error handle self upgrade",
			fields: fields{
				repo: func(t *testing.T) blueprintSpecRepository {
					m := newMockBlueprintSpecRepository(t)
					m.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
					return m
				},
				validation: func(t *testing.T) blueprintSpecValidationUseCase {
					m := newMockBlueprintSpecValidationUseCase(t)
					m.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, testBlueprintSpec).Return(nil)
					m.EXPECT().ValidateBlueprintSpecDynamically(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				effectiveBlueprint: func(t *testing.T) effectiveBlueprintUseCase {
					m := newMockEffectiveBlueprintUseCase(t)
					m.EXPECT().CalculateEffectiveBlueprint(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				stateDiff: func(t *testing.T) stateDiffUseCase {
					m := newMockStateDiffUseCase(t)
					m.EXPECT().DetermineStateDiff(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				healthUseCase: func(t *testing.T) ecosystemHealthUseCase {
					m := newMockEcosystemHealthUseCase(t)
					m.EXPECT().CheckEcosystemHealth(mock.Anything, testBlueprintSpec).Return(ecosystem.HealthResult{}, nil)
					return m
				},
				selfUpgradeUseCase: func(t *testing.T) selfUpgradeUseCase {
					m := newMockSelfUpgradeUseCase(t)
					m.EXPECT().HandleSelfUpgrade(mock.Anything, testBlueprintSpec).Return(assert.AnError)
					return m
				},
			},
			args: testArgs,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err)
			},
		},
		{
			name: "should return error on error apply config",
			fields: fields{
				repo: func(t *testing.T) blueprintSpecRepository {
					m := newMockBlueprintSpecRepository(t)
					m.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
					return m
				},
				validation: func(t *testing.T) blueprintSpecValidationUseCase {
					m := newMockBlueprintSpecValidationUseCase(t)
					m.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, testBlueprintSpec).Return(nil)
					m.EXPECT().ValidateBlueprintSpecDynamically(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				effectiveBlueprint: func(t *testing.T) effectiveBlueprintUseCase {
					m := newMockEffectiveBlueprintUseCase(t)
					m.EXPECT().CalculateEffectiveBlueprint(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				stateDiff: func(t *testing.T) stateDiffUseCase {
					m := newMockStateDiffUseCase(t)
					m.EXPECT().DetermineStateDiff(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				healthUseCase: func(t *testing.T) ecosystemHealthUseCase {
					m := newMockEcosystemHealthUseCase(t)
					m.EXPECT().CheckEcosystemHealth(mock.Anything, testBlueprintSpec).Return(ecosystem.HealthResult{}, nil)
					return m
				},
				selfUpgradeUseCase: func(t *testing.T) selfUpgradeUseCase {
					m := newMockSelfUpgradeUseCase(t)
					m.EXPECT().HandleSelfUpgrade(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				ecosystemConfigUseCase: func(t *testing.T) ecosystemConfigUseCase {
					m := newMockEcosystemConfigUseCase(t)
					m.EXPECT().ApplyConfig(mock.Anything, testBlueprintSpec).Return(assert.AnError)
					return m
				},
			},
			args: testArgs,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err)
			},
		},
		{
			name: "should return error on error apply components",
			fields: fields{
				repo: func(t *testing.T) blueprintSpecRepository {
					m := newMockBlueprintSpecRepository(t)
					m.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
					return m
				},
				validation: func(t *testing.T) blueprintSpecValidationUseCase {
					m := newMockBlueprintSpecValidationUseCase(t)
					m.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, testBlueprintSpec).Return(nil)
					m.EXPECT().ValidateBlueprintSpecDynamically(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				effectiveBlueprint: func(t *testing.T) effectiveBlueprintUseCase {
					m := newMockEffectiveBlueprintUseCase(t)
					m.EXPECT().CalculateEffectiveBlueprint(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				stateDiff: func(t *testing.T) stateDiffUseCase {
					m := newMockStateDiffUseCase(t)
					m.EXPECT().DetermineStateDiff(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				healthUseCase: func(t *testing.T) ecosystemHealthUseCase {
					m := newMockEcosystemHealthUseCase(t)
					m.EXPECT().CheckEcosystemHealth(mock.Anything, testBlueprintSpec).Return(ecosystem.HealthResult{}, nil)
					return m
				},
				selfUpgradeUseCase: func(t *testing.T) selfUpgradeUseCase {
					m := newMockSelfUpgradeUseCase(t)
					m.EXPECT().HandleSelfUpgrade(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				ecosystemConfigUseCase: func(t *testing.T) ecosystemConfigUseCase {
					m := newMockEcosystemConfigUseCase(t)
					m.EXPECT().ApplyConfig(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				applyComponentUseCase: func(t *testing.T) applyComponentsUseCase {
					m := newMockApplyComponentsUseCase(t)
					m.EXPECT().ApplyComponents(mock.Anything, testBlueprintSpec).Return(false, assert.AnError)
					return m
				},
			},
			args: testArgs,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err)
			},
		},
		{
			name: "should return error on error check health after component apply",
			fields: fields{
				repo: func(t *testing.T) blueprintSpecRepository {
					m := newMockBlueprintSpecRepository(t)
					m.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
					return m
				},
				validation: func(t *testing.T) blueprintSpecValidationUseCase {
					m := newMockBlueprintSpecValidationUseCase(t)
					m.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, testBlueprintSpec).Return(nil)
					m.EXPECT().ValidateBlueprintSpecDynamically(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				effectiveBlueprint: func(t *testing.T) effectiveBlueprintUseCase {
					m := newMockEffectiveBlueprintUseCase(t)
					m.EXPECT().CalculateEffectiveBlueprint(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				stateDiff: func(t *testing.T) stateDiffUseCase {
					m := newMockStateDiffUseCase(t)
					m.EXPECT().DetermineStateDiff(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				healthUseCase: func(t *testing.T) ecosystemHealthUseCase {
					m := newMockEcosystemHealthUseCase(t)
					m.EXPECT().CheckEcosystemHealth(mock.Anything, testBlueprintSpec).Return(ecosystem.HealthResult{}, nil).Times(1)
					m.EXPECT().CheckEcosystemHealth(mock.Anything, testBlueprintSpec).Return(ecosystem.HealthResult{}, assert.AnError).Times(1)
					return m
				},
				selfUpgradeUseCase: func(t *testing.T) selfUpgradeUseCase {
					m := newMockSelfUpgradeUseCase(t)
					m.EXPECT().HandleSelfUpgrade(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				ecosystemConfigUseCase: func(t *testing.T) ecosystemConfigUseCase {
					m := newMockEcosystemConfigUseCase(t)
					m.EXPECT().ApplyConfig(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				applyComponentUseCase: func(t *testing.T) applyComponentsUseCase {
					m := newMockApplyComponentsUseCase(t)
					m.EXPECT().ApplyComponents(mock.Anything, testBlueprintSpec).Return(true, nil)
					return m
				},
			},
			args: testArgs,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err)
			},
		},
		{
			name: "should return error on error apply dogus",
			fields: fields{
				repo: func(t *testing.T) blueprintSpecRepository {
					m := newMockBlueprintSpecRepository(t)
					m.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
					return m
				},
				validation: func(t *testing.T) blueprintSpecValidationUseCase {
					m := newMockBlueprintSpecValidationUseCase(t)
					m.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, testBlueprintSpec).Return(nil)
					m.EXPECT().ValidateBlueprintSpecDynamically(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				effectiveBlueprint: func(t *testing.T) effectiveBlueprintUseCase {
					m := newMockEffectiveBlueprintUseCase(t)
					m.EXPECT().CalculateEffectiveBlueprint(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				stateDiff: func(t *testing.T) stateDiffUseCase {
					m := newMockStateDiffUseCase(t)
					m.EXPECT().DetermineStateDiff(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				healthUseCase: func(t *testing.T) ecosystemHealthUseCase {
					m := newMockEcosystemHealthUseCase(t)
					m.EXPECT().CheckEcosystemHealth(mock.Anything, testBlueprintSpec).Return(ecosystem.HealthResult{}, nil)
					return m
				},
				selfUpgradeUseCase: func(t *testing.T) selfUpgradeUseCase {
					m := newMockSelfUpgradeUseCase(t)
					m.EXPECT().HandleSelfUpgrade(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				ecosystemConfigUseCase: func(t *testing.T) ecosystemConfigUseCase {
					m := newMockEcosystemConfigUseCase(t)
					m.EXPECT().ApplyConfig(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				applyComponentUseCase: func(t *testing.T) applyComponentsUseCase {
					m := newMockApplyComponentsUseCase(t)
					m.EXPECT().ApplyComponents(mock.Anything, testBlueprintSpec).Return(false, nil)
					return m
				},
				applyDogusUseCase: func(t *testing.T) applyDogusUseCase {
					m := newMockApplyDogusUseCase(t)
					m.EXPECT().ApplyDogus(mock.Anything, testBlueprintSpec).Return(false, assert.AnError)
					return m
				},
			},
			args: testArgs,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err)
			},
		},
		{
			name: "should return error on error check health after dogu apply",
			fields: fields{
				repo: func(t *testing.T) blueprintSpecRepository {
					m := newMockBlueprintSpecRepository(t)
					m.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
					return m
				},
				validation: func(t *testing.T) blueprintSpecValidationUseCase {
					m := newMockBlueprintSpecValidationUseCase(t)
					m.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, testBlueprintSpec).Return(nil)
					m.EXPECT().ValidateBlueprintSpecDynamically(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				effectiveBlueprint: func(t *testing.T) effectiveBlueprintUseCase {
					m := newMockEffectiveBlueprintUseCase(t)
					m.EXPECT().CalculateEffectiveBlueprint(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				stateDiff: func(t *testing.T) stateDiffUseCase {
					m := newMockStateDiffUseCase(t)
					m.EXPECT().DetermineStateDiff(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				healthUseCase: func(t *testing.T) ecosystemHealthUseCase {
					m := newMockEcosystemHealthUseCase(t)
					m.EXPECT().CheckEcosystemHealth(mock.Anything, testBlueprintSpec).Return(ecosystem.HealthResult{}, nil).Times(1)
					m.EXPECT().CheckEcosystemHealth(mock.Anything, testBlueprintSpec).Return(ecosystem.HealthResult{}, assert.AnError)
					return m
				},
				selfUpgradeUseCase: func(t *testing.T) selfUpgradeUseCase {
					m := newMockSelfUpgradeUseCase(t)
					m.EXPECT().HandleSelfUpgrade(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				ecosystemConfigUseCase: func(t *testing.T) ecosystemConfigUseCase {
					m := newMockEcosystemConfigUseCase(t)
					m.EXPECT().ApplyConfig(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				applyComponentUseCase: func(t *testing.T) applyComponentsUseCase {
					m := newMockApplyComponentsUseCase(t)
					m.EXPECT().ApplyComponents(mock.Anything, testBlueprintSpec).Return(false, nil)
					return m
				},
				applyDogusUseCase: func(t *testing.T) applyDogusUseCase {
					m := newMockApplyDogusUseCase(t)
					m.EXPECT().ApplyDogus(mock.Anything, testBlueprintSpec).Return(true, nil)
					return m
				},
			},
			args: testArgs,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err)
			},
		},
		{
			name: "should return error on error complete blueprint",
			fields: fields{
				repo: func(t *testing.T) blueprintSpecRepository {
					m := newMockBlueprintSpecRepository(t)
					m.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
					return m
				},
				validation: func(t *testing.T) blueprintSpecValidationUseCase {
					m := newMockBlueprintSpecValidationUseCase(t)
					m.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, testBlueprintSpec).Return(nil)
					m.EXPECT().ValidateBlueprintSpecDynamically(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				effectiveBlueprint: func(t *testing.T) effectiveBlueprintUseCase {
					m := newMockEffectiveBlueprintUseCase(t)
					m.EXPECT().CalculateEffectiveBlueprint(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				stateDiff: func(t *testing.T) stateDiffUseCase {
					m := newMockStateDiffUseCase(t)
					m.EXPECT().DetermineStateDiff(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				healthUseCase: func(t *testing.T) ecosystemHealthUseCase {
					m := newMockEcosystemHealthUseCase(t)
					m.EXPECT().CheckEcosystemHealth(mock.Anything, testBlueprintSpec).Return(ecosystem.HealthResult{}, nil)
					return m
				},
				selfUpgradeUseCase: func(t *testing.T) selfUpgradeUseCase {
					m := newMockSelfUpgradeUseCase(t)
					m.EXPECT().HandleSelfUpgrade(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				ecosystemConfigUseCase: func(t *testing.T) ecosystemConfigUseCase {
					m := newMockEcosystemConfigUseCase(t)
					m.EXPECT().ApplyConfig(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				applyComponentUseCase: func(t *testing.T) applyComponentsUseCase {
					m := newMockApplyComponentsUseCase(t)
					m.EXPECT().ApplyComponents(mock.Anything, testBlueprintSpec).Return(false, nil)
					return m
				},
				applyDogusUseCase: func(t *testing.T) applyDogusUseCase {
					m := newMockApplyDogusUseCase(t)
					m.EXPECT().ApplyDogus(mock.Anything, testBlueprintSpec).Return(false, nil)
					return m
				},
				applyUseCase: func(t *testing.T) completeBlueprintUseCase {
					m := newMockCompleteBlueprintUseCase(t)
					m.EXPECT().CompleteBlueprint(mock.Anything, testBlueprintSpec).Return(assert.AnError)
					return m
				},
			},
			args: testArgs,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err)
			},
		},
		{
			name: "should return nil on success",
			fields: fields{
				repo: func(t *testing.T) blueprintSpecRepository {
					m := newMockBlueprintSpecRepository(t)
					m.EXPECT().GetById(mock.Anything, testBlueprintId).Return(testBlueprintSpec, nil)
					return m
				},
				validation: func(t *testing.T) blueprintSpecValidationUseCase {
					m := newMockBlueprintSpecValidationUseCase(t)
					m.EXPECT().ValidateBlueprintSpecStatically(mock.Anything, testBlueprintSpec).Return(nil)
					m.EXPECT().ValidateBlueprintSpecDynamically(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				effectiveBlueprint: func(t *testing.T) effectiveBlueprintUseCase {
					m := newMockEffectiveBlueprintUseCase(t)
					m.EXPECT().CalculateEffectiveBlueprint(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				stateDiff: func(t *testing.T) stateDiffUseCase {
					m := newMockStateDiffUseCase(t)
					m.EXPECT().DetermineStateDiff(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				healthUseCase: func(t *testing.T) ecosystemHealthUseCase {
					m := newMockEcosystemHealthUseCase(t)
					m.EXPECT().CheckEcosystemHealth(mock.Anything, testBlueprintSpec).Return(ecosystem.HealthResult{}, nil)
					return m
				},
				selfUpgradeUseCase: func(t *testing.T) selfUpgradeUseCase {
					m := newMockSelfUpgradeUseCase(t)
					m.EXPECT().HandleSelfUpgrade(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				ecosystemConfigUseCase: func(t *testing.T) ecosystemConfigUseCase {
					m := newMockEcosystemConfigUseCase(t)
					m.EXPECT().ApplyConfig(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
				applyComponentUseCase: func(t *testing.T) applyComponentsUseCase {
					m := newMockApplyComponentsUseCase(t)
					m.EXPECT().ApplyComponents(mock.Anything, testBlueprintSpec).Return(false, nil)
					return m
				},
				applyDogusUseCase: func(t *testing.T) applyDogusUseCase {
					m := newMockApplyDogusUseCase(t)
					m.EXPECT().ApplyDogus(mock.Anything, testBlueprintSpec).Return(false, nil)
					return m
				},
				applyUseCase: func(t *testing.T) completeBlueprintUseCase {
					m := newMockCompleteBlueprintUseCase(t)
					m.EXPECT().CompleteBlueprint(mock.Anything, testBlueprintSpec).Return(nil)
					return m
				},
			},
			args:    testArgs,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var repo blueprintSpecRepository
			if tt.fields.repo != nil {
				repo = tt.fields.repo(t)
			}

			var validation blueprintSpecValidationUseCase
			if tt.fields.validation != nil {
				validation = tt.fields.validation(t)
			}

			var effectiveBlueprint effectiveBlueprintUseCase
			if tt.fields.effectiveBlueprint != nil {
				effectiveBlueprint = tt.fields.effectiveBlueprint(t)
			}

			var stateDiff stateDiffUseCase
			if tt.fields.stateDiff != nil {
				stateDiff = tt.fields.stateDiff(t)
			}

			var applyUseCase completeBlueprintUseCase
			if tt.fields.applyUseCase != nil {
				applyUseCase = tt.fields.applyUseCase(t)
			}

			var ecoConfigUseCase ecosystemConfigUseCase
			if tt.fields.ecosystemConfigUseCase != nil {
				ecoConfigUseCase = tt.fields.ecosystemConfigUseCase(t)
			}

			var selfUpgrade selfUpgradeUseCase
			if tt.fields.selfUpgradeUseCase != nil {
				selfUpgrade = tt.fields.selfUpgradeUseCase(t)
			}

			var applyComponentUseCase applyComponentsUseCase
			if tt.fields.applyComponentUseCase != nil {
				applyComponentUseCase = tt.fields.applyComponentUseCase(t)
			}

			var applyDoguUseCase applyDogusUseCase
			if tt.fields.applyDogusUseCase != nil {
				applyDoguUseCase = tt.fields.applyDogusUseCase(t)
			}

			var ecoHealthUseCase ecosystemHealthUseCase
			if tt.fields.healthUseCase != nil {
				ecoHealthUseCase = tt.fields.healthUseCase(t)
			}

			useCase := &BlueprintSpecChangeUseCase{
				repo:                   repo,
				validation:             validation,
				effectiveBlueprint:     effectiveBlueprint,
				stateDiff:              stateDiff,
				applyUseCase:           applyUseCase,
				ecosystemConfigUseCase: ecoConfigUseCase,
				selfUpgradeUseCase:     selfUpgrade,
				applyComponentUseCase:  applyComponentUseCase,
				applyDogusUseCase:      applyDoguUseCase,
				healthUseCase:          ecoHealthUseCase,
			}
			tt.wantErr(t, useCase.HandleUntilApplied(tt.args.givenCtx, tt.args.blueprintId), fmt.Sprintf("HandleUntilApplied(%v, %v)", tt.args.givenCtx, tt.args.blueprintId))
		})
	}
}
