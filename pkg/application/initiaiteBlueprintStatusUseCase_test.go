package application

import (
	"context"
	"testing"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestInitiateBlueprintStatusUseCase_InitateConditions(t *testing.T) {
	//type expectFn func(m *mockBlueprintSpecRepository) *mock.Call
	type args struct {
		blueprint *domain.BlueprintSpec
	}
	tests := []struct {
		name                  string
		args                  args
		wantUnknownConditions []string
		wantErr               error
	}{
		{
			name: "nil conditions",
			args: args{
				blueprint: &domain.BlueprintSpec{
					Conditions: nil,
				},
			},
			wantUnknownConditions: domain.BlueprintConditions,
			wantErr:               nil,
		},
		{
			name: "empty conditions",
			args: args{
				blueprint: &domain.BlueprintSpec{
					Conditions: []domain.Condition{},
				},
			},
			wantUnknownConditions: domain.BlueprintConditions,
			wantErr:               nil,
		},
		{
			name: "some conditions",
			args: args{
				blueprint: &domain.BlueprintSpec{
					Conditions: []domain.Condition{
						{
							Type: domain.ConditionValid,
						},
						{
							Type: domain.ConditionLastApplySucceeded,
						},
					},
				},
			},
			wantUnknownConditions: []string{domain.ConditionExecutable, domain.ConditionEcosystemHealthy, domain.ConditionCompleted, domain.ConditionSelfUpgradeCompleted},
			wantErr:               nil,
		},
		{
			name: "all conditions",
			args: args{
				blueprint: &domain.BlueprintSpec{
					Conditions: []domain.Condition{
						{
							Type: domain.ConditionValid,
						},
						{
							Type: domain.ConditionExecutable,
						},
						{
							Type: domain.ConditionEcosystemHealthy,
						},
						{
							Type: domain.ConditionSelfUpgradeCompleted,
						},
						{
							Type: domain.ConditionCompleted,
						},
						{
							Type: domain.ConditionLastApplySucceeded,
						},
					},
				},
			},
			wantUnknownConditions: nil,
			wantErr:               nil,
		},
		{
			name: "update error",
			args: args{
				blueprint: &domain.BlueprintSpec{
					Conditions: nil,
				},
			},
			wantUnknownConditions: domain.BlueprintConditions,
			wantErr:               assert.AnError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			repoMock := newMockBlueprintSpecRepository(t)
			if len(tt.wantUnknownConditions) > 0 {
				repoMock.EXPECT().Update(ctx, mock.AnythingOfType("*domain.BlueprintSpec")).RunAndReturn(func(ctx context.Context, bp *domain.BlueprintSpec) error {
					assert.Len(t, bp.Conditions, len(domain.BlueprintConditions))
					for _, condition := range tt.wantUnknownConditions {
						assert.True(t, meta.IsStatusConditionPresentAndEqual(bp.Conditions, condition, metav1.ConditionUnknown))
						bpCondition := meta.FindStatusCondition(bp.Conditions, condition)
						assert.Equal(t, "InitialSyncPending", bpCondition.Reason)
						assert.Equal(t, "controller has not determined this condition yet", bpCondition.Message)
					}
					return tt.wantErr
				})
			}

			useCase := &InitiateBlueprintStatusUseCase{
				repo: repoMock,
			}
			err := useCase.InitateConditions(ctx, tt.args.blueprint)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.ErrorContains(t, err, "cannot save blueprint spec \"\" after initially setting the conditions to unknown")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewInitiateBlueprintStatusUseCase(t *testing.T) {
	type args struct {
		repo blueprintSpecRepository
	}
	tests := []struct {
		name string
		args args
		want *InitiateBlueprintStatusUseCase
	}{
		{
			name: "nil repo",
			args: args{
				repo: nil,
			},
			want: &InitiateBlueprintStatusUseCase{},
		},
		{
			name: "repo",
			args: args{
				repo: &mockBlueprintSpecRepository{},
			},
			want: &InitiateBlueprintStatusUseCase{
				repo: &mockBlueprintSpecRepository{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewInitiateBlueprintStatusUseCase(tt.args.repo), "NewInitiateBlueprintStatusUseCase(%v)", tt.args.repo)
		})
	}
}
