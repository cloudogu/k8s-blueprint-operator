package application

import (
	"context"
	"maps"
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	liberrors "github.com/cloudogu/ces-commons-lib/errors"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
)

const (
	redmine = cescommons.SimpleName("redmine")
	cas     = cescommons.SimpleName("cas")
)

var emptyDoguList []cescommons.SimpleName

func TestEcosystemConfigUseCase_ApplyConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		doguConfigMock := newMockDoguConfigRepository(t)
		sensitiveDoguConfigMock := newMockSensitiveDoguConfigRepository(t)
		globalConfigRepoMock := newMockGlobalConfigRepository(t)

		sensitiveRedmineDiff := getSensitiveDoguConfigEntryDiffForAction("key", "value", redmine, domain.ConfigActionSet)
		sensitiveCasDiff := getSensitiveDoguConfigEntryDiffForAction("key", "value", cas, domain.ConfigActionSet)
		blueprint := &domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				DoguDiffs: []domain.DoguDiff{
					{
						DoguName:      redmine,
						NeededActions: []domain.Action{},
					},
					{
						DoguName:      cas,
						NeededActions: []domain.Action{},
					},
				},
				DoguConfigDiffs: map[cescommons.SimpleName]domain.DoguConfigDiffs{
					redmine: {
						getSetDoguConfigEntryDiff("key", "value", redmine),
						getRemoveDoguConfigEntryDiff("key", redmine),
					},
					cas: {
						getSetDoguConfigEntryDiff("key", "value", cas),
						getRemoveDoguConfigEntryDiff("key", cas),
					},
				},
				SensitiveDoguConfigDiffs: map[cescommons.SimpleName]domain.SensitiveDoguConfigDiffs{
					redmine: {
						sensitiveRedmineDiff,
						getRemoveSensitiveDoguConfigEntryDiff("key", redmine),
					},
					cas: {
						sensitiveCasDiff,
						getRemoveSensitiveDoguConfigEntryDiff("key", cas),
					},
				},
				GlobalConfigDiffs: domain.GlobalConfigDiffs{
					getSetGlobalConfigEntryDiff("key", "value"),
					getRemoveGlobalConfigEntryDiff("key"),
				},
			},
		}

		// Just check if the routine hits the repos. Check values in concrete test of methods.
		doguConfigMock.EXPECT().
			GetAllExisting(testCtx, []cescommons.SimpleName{cas, redmine}).
			Return(map[cescommons.SimpleName]config.DoguConfig{
				redmine: config.CreateDoguConfig(redmine, map[config.Key]config.Value{}),
				cas:     config.CreateDoguConfig(cas, map[config.Key]config.Value{}),
			}, nil)
		doguConfigMock.EXPECT().UpdateOrCreate(testCtx, mock.Anything).Return(config.DoguConfig{}, nil).Times(2)

		sensitiveDoguConfigMock.EXPECT().
			GetAllExisting(testCtx, []cescommons.SimpleName{cas, redmine}).
			Return(map[cescommons.SimpleName]config.DoguConfig{
				redmine: config.CreateDoguConfig(redmine, map[config.Key]config.Value{}),
				cas:     config.CreateDoguConfig(cas, map[config.Key]config.Value{}),
			}, nil)
		sensitiveDoguConfigMock.EXPECT().UpdateOrCreate(testCtx, mock.Anything).Return(config.DoguConfig{}, nil).Times(2)

		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)
		globalConfigRepoMock.EXPECT().Update(testCtx, mock.Anything).Return(globalConfig, nil)

		blueprintRepoMock.EXPECT().Update(testCtx, mock.Anything).Return(nil)
		doguInstallaltionRepoMock := newMockDoguInstallationRepository(t)
		doguInstallaltionRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		sut := NewEcosystemConfigUseCase(blueprintRepoMock, doguConfigMock, sensitiveDoguConfigMock, globalConfigRepoMock, doguInstallaltionRepoMock)

		// when
		err := sut.ApplyConfig(testCtx, blueprint)

		// then
		require.NoError(t, err)
		require.Len(t, blueprint.Events, 1)
		assert.Equal(t, domain.EcosystemConfigAppliedEvent{}, blueprint.Events[0])
	})

	t.Run("pause reconciliation for dogus with config and version changes", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		doguConfigMock := newMockDoguConfigRepository(t)
		sensitiveDoguConfigMock := newMockSensitiveDoguConfigRepository(t)
		globalConfigRepoMock := newMockGlobalConfigRepository(t)

		sensitiveCasDiff := getSensitiveDoguConfigEntryDiffForAction("key", "value", cas, domain.ConfigActionSet)
		blueprint := &domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				DoguDiffs: []domain.DoguDiff{
					{
						DoguName:      redmine,
						NeededActions: []domain.Action{domain.ActionUpgrade},
					},
					{
						DoguName:      cas,
						NeededActions: []domain.Action{domain.ActionUpgrade},
					},
				},
				DoguConfigDiffs: map[cescommons.SimpleName]domain.DoguConfigDiffs{
					redmine: {
						getSetDoguConfigEntryDiff("key", "value", redmine),
					},
				},
				SensitiveDoguConfigDiffs: map[cescommons.SimpleName]domain.SensitiveDoguConfigDiffs{
					cas: {
						sensitiveCasDiff,
					},
				},
			},
		}

		// Just check if the routine hits the repos. Check values in concrete test of methods.
		doguConfigMock.EXPECT().
			GetAllExisting(testCtx, []cescommons.SimpleName{redmine}).
			Return(map[cescommons.SimpleName]config.DoguConfig{
				redmine: config.CreateDoguConfig(redmine, map[config.Key]config.Value{}),
			}, nil)
		doguConfigMock.EXPECT().UpdateOrCreate(testCtx, mock.Anything).Return(config.DoguConfig{}, nil)

		sensitiveDoguConfigMock.EXPECT().
			GetAllExisting(testCtx, []cescommons.SimpleName{cas}).
			Return(map[cescommons.SimpleName]config.DoguConfig{
				cas: config.CreateDoguConfig(cas, map[config.Key]config.Value{}),
			}, nil)
		sensitiveDoguConfigMock.EXPECT().UpdateOrCreate(testCtx, mock.Anything).Return(config.DoguConfig{}, nil)

		globalConfigRepoMock.EXPECT().Get(testCtx).Return(config.GlobalConfig{}, nil)

		blueprintRepoMock.EXPECT().Update(testCtx, mock.Anything).Return(nil)
		doguInstallaltionRepoMock := newMockDoguInstallationRepository(t)
		dogus := map[cescommons.SimpleName]*ecosystem.DoguInstallation{
			redmine: {Name: cescommons.QualifiedName{SimpleName: redmine, Namespace: "namespace"}},
			cas:     {Name: cescommons.QualifiedName{SimpleName: cas, Namespace: "namespace"}},
		}
		doguInstallaltionRepoMock.EXPECT().GetAll(testCtx).Return(dogus, nil)
		doguInstallaltionRepoMock.EXPECT().Update(testCtx, mock.Anything).Run(func(ctx context.Context, dogu *ecosystem.DoguInstallation) {
			assert.True(t, dogu.PauseReconciliation)
		}).Return(nil).Times(2)
		sut := NewEcosystemConfigUseCase(blueprintRepoMock, doguConfigMock, sensitiveDoguConfigMock, globalConfigRepoMock, doguInstallaltionRepoMock)

		// when
		err := sut.ApplyConfig(testCtx, blueprint)

		// then
		require.NoError(t, err)
		require.Len(t, blueprint.Events, 1)
		assert.Equal(t, domain.EcosystemConfigAppliedEvent{}, blueprint.Events[0])
	})

	t.Run("error applying dogu config", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
			StateDiff: domain.StateDiff{
				DoguConfigDiffs: map[cescommons.SimpleName]domain.DoguConfigDiffs{
					redmine: {
						getSetDoguConfigEntryDiff("key", "value", redmine),
					},
					cas: {
						getSetDoguConfigEntryDiff("key", "value", cas),
					},
				},
			},
		}

		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		sensitiveDoguConfigMock := newMockSensitiveDoguConfigRepository(t)
		globalConfigMock := newMockGlobalConfigRepository(t)
		doguConfigMock := newMockDoguConfigRepository(t)
		// Just check if the routine hits the repos. Check values in concrete test of methods.
		doguConfigMock.EXPECT().
			GetAllExisting(testCtx, []cescommons.SimpleName{cas, redmine}).
			Return(map[cescommons.SimpleName]config.DoguConfig{
				redmine: config.CreateDoguConfig(redmine, map[config.Key]config.Value{}),
				cas:     config.CreateDoguConfig(cas, map[config.Key]config.Value{}),
			}, nil)
		doguConfigMock.EXPECT().UpdateOrCreate(testCtx, mock.Anything).Return(config.DoguConfig{}, assert.AnError).Times(1)
		blueprintRepoMock.EXPECT().Update(testCtx, mock.Anything).Return(nil)
		doguInstallaltionRepoMock := newMockDoguInstallationRepository(t)
		doguInstallaltionRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		sut := NewEcosystemConfigUseCase(blueprintRepoMock, doguConfigMock, sensitiveDoguConfigMock, globalConfigMock, doguInstallaltionRepoMock)

		// when
		err := sut.ApplyConfig(testCtx, blueprint)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "could not apply normal dogu config")
		assert.True(t, meta.IsStatusConditionFalse(blueprint.Conditions, domain.ConditionLastApplySucceeded))

		require.Len(t, blueprint.Events, 1)
		assert.Equal(t, domain.NewExecutionFailedEvent(err), blueprint.Events[0])
	})
	t.Run("error applying sensitive config", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		doguConfigMock := newMockDoguConfigRepository(t)
		sensitiveDoguConfigMock := newMockSensitiveDoguConfigRepository(t)
		globalConfigMock := newMockGlobalConfigRepository(t)

		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
			StateDiff: domain.StateDiff{
				SensitiveDoguConfigDiffs: map[cescommons.SimpleName]domain.DoguConfigDiffs{
					redmine: {
						getSensitiveDoguConfigEntryDiffForAction("key", "value", redmine, domain.ConfigActionSet),
					},
					cas: {
						getSensitiveDoguConfigEntryDiffForAction("key", "value", cas, domain.ConfigActionSet),
					},
				},
			},
		}

		// Just check if the routine hits the repos. Check values in concrete test of methods.
		doguConfigMock.EXPECT().
			GetAllExisting(testCtx, emptyDoguList).
			Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)
		sensitiveDoguConfigMock.EXPECT().
			GetAllExisting(testCtx, []cescommons.SimpleName{cas, redmine}).
			Return(map[cescommons.SimpleName]config.DoguConfig{
				redmine: config.CreateDoguConfig(redmine, map[config.Key]config.Value{}),
				cas:     config.CreateDoguConfig(cas, map[config.Key]config.Value{}),
			}, nil)
		sensitiveDoguConfigMock.EXPECT().UpdateOrCreate(testCtx, mock.Anything).Return(config.DoguConfig{}, assert.AnError).Times(1)
		blueprintRepoMock.EXPECT().Update(testCtx, mock.Anything).Return(nil)
		doguInstallaltionRepoMock := newMockDoguInstallationRepository(t)
		doguInstallaltionRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)
		sut := NewEcosystemConfigUseCase(blueprintRepoMock, doguConfigMock, sensitiveDoguConfigMock, globalConfigMock, doguInstallaltionRepoMock)

		// when
		err := sut.ApplyConfig(testCtx, blueprint)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "could not apply sensitive dogu config")

		assert.True(t, meta.IsStatusConditionFalse(blueprint.Conditions, domain.ConditionLastApplySucceeded))

		require.Len(t, blueprint.Events, 1)
		assert.Equal(t, domain.NewExecutionFailedEvent(err), blueprint.Events[0])
		assert.Contains(t, blueprint.Events[0].Message(), "could not apply sensitive dogu config")
		// cannot check for dogu name here as the order of the events is not fixed. It could be either redmine or cas
		assert.Contains(t, blueprint.Events[0].Message(), "could not persist config for dogu")
		assert.Contains(t, blueprint.Events[0].Message(), "assert.AnError general error for testing")
	})
	t.Run("error applying global config", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		doguConfigMock := newMockDoguConfigRepository(t)
		sensitiveDoguConfigMock := newMockSensitiveDoguConfigRepository(t)
		globalConfigMock := newMockGlobalConfigRepository(t)

		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
			StateDiff: domain.StateDiff{
				DoguConfigDiffs:          map[cescommons.SimpleName]domain.DoguConfigDiffs{},
				SensitiveDoguConfigDiffs: map[cescommons.SimpleName]domain.SensitiveDoguConfigDiffs{},
				GlobalConfigDiffs: domain.GlobalConfigDiffs{
					getSetGlobalConfigEntryDiff("key", "value"),
				},
			},
		}

		// Just check if the routine hits the repos. Check values in concrete test of methods.

		doguConfigMock.EXPECT().
			GetAllExisting(testCtx, emptyDoguList).
			Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)
		sensitiveDoguConfigMock.EXPECT().
			GetAllExisting(testCtx, emptyDoguList).
			Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)

		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigMock.EXPECT().Get(testCtx).Return(globalConfig, nil)
		globalConfigMock.EXPECT().Update(testCtx, mock.Anything).Return(globalConfig, assert.AnError)

		blueprintRepoMock.EXPECT().Update(testCtx, blueprint).Return(nil)
		doguInstallaltionRepoMock := newMockDoguInstallationRepository(t)
		doguInstallaltionRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)
		sut := NewEcosystemConfigUseCase(blueprintRepoMock, doguConfigMock, sensitiveDoguConfigMock, globalConfigMock, doguInstallaltionRepoMock)

		// when
		err := sut.ApplyConfig(testCtx, blueprint)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "could not apply global config")

		assert.True(t, meta.IsStatusConditionFalse(blueprint.Conditions, domain.ConditionLastApplySucceeded))

		require.Len(t, blueprint.Events, 1)
		assert.Equal(t, domain.NewExecutionFailedEvent(err), blueprint.Events[0])
		assert.Contains(t, blueprint.Events[0].Message(), "could not apply global config")
		assert.Contains(t, blueprint.Events[0].Message(), "assert.AnError general error for testing")
	})
}

func TestEcosystemConfigUseCase_pauseReconciliationForDogus(t *testing.T) {
	t.Run("pause reconciliation for multiple dogus when dogu config changes", func(t *testing.T) {
		// given
		stateDiff := domain.StateDiff{
			DoguDiffs: []domain.DoguDiff{
				{
					DoguName:      redmine,
					NeededActions: []domain.Action{domain.ActionUpgrade},
				},
				{
					DoguName:      cas,
					NeededActions: []domain.Action{domain.ActionUpgrade},
				},
				{
					DoguName:      postfix,
					NeededActions: []domain.Action{}, // no action, so no pause required
				},
			},
			DoguConfigDiffs: map[cescommons.SimpleName]domain.DoguConfigDiffs{
				redmine: {
					getSetDoguConfigEntryDiff("key", "value", redmine),
				},
				cas: {
					getSetDoguConfigEntryDiff("key", "value", cas),
				},
				postfix: {
					getSetDoguConfigEntryDiff("key", "value", postfix),
				},
			},
		}

		doguInstallaltionRepoMock := newMockDoguInstallationRepository(t)
		dogus := map[cescommons.SimpleName]*ecosystem.DoguInstallation{
			redmine: {Name: cescommons.QualifiedName{SimpleName: redmine, Namespace: "namespace"}},
			cas:     {Name: cescommons.QualifiedName{SimpleName: cas, Namespace: "namespace"}},
			postfix: {Name: cescommons.QualifiedName{SimpleName: postfix, Namespace: "namespace"}},
		}
		doguInstallaltionRepoMock.EXPECT().GetAll(testCtx).Return(dogus, nil)
		doguInstallaltionRepoMock.EXPECT().Update(testCtx, mock.Anything).Run(func(ctx context.Context, dogu *ecosystem.DoguInstallation) {
			assert.True(t, dogu.PauseReconciliation)
		}).Return(nil).Times(2)

		// when
		err := pauseReconciliationForDogus(testCtx, doguInstallaltionRepoMock, stateDiff)

		// then
		require.NoError(t, err)
	})

	t.Run("pause reconciliation for multiple dogus when sensitive config changes", func(t *testing.T) {
		// given
		stateDiff := domain.StateDiff{
			DoguDiffs: []domain.DoguDiff{
				{
					DoguName:      redmine,
					NeededActions: []domain.Action{domain.ActionUpgrade},
				},
				{
					DoguName:      cas,
					NeededActions: []domain.Action{domain.ActionUpgrade},
				},
				{
					DoguName:      postfix,
					NeededActions: []domain.Action{}, // no action, so no pause required
				},
			},
			SensitiveDoguConfigDiffs: map[cescommons.SimpleName]domain.SensitiveDoguConfigDiffs{
				redmine: {
					getSetDoguConfigEntryDiff("key", "value", redmine),
				},
				cas: {
					getSetDoguConfigEntryDiff("key", "value", cas),
				},
				postfix: {
					getSetDoguConfigEntryDiff("key", "value", postfix),
				},
			},
		}

		doguInstallaltionRepoMock := newMockDoguInstallationRepository(t)
		dogus := map[cescommons.SimpleName]*ecosystem.DoguInstallation{
			redmine: {Name: cescommons.QualifiedName{SimpleName: redmine, Namespace: "namespace"}},
			cas:     {Name: cescommons.QualifiedName{SimpleName: cas, Namespace: "namespace"}},
			postfix: {Name: cescommons.QualifiedName{SimpleName: postfix, Namespace: "namespace"}},
		}
		doguInstallaltionRepoMock.EXPECT().GetAll(testCtx).Return(dogus, nil)
		doguInstallaltionRepoMock.EXPECT().Update(testCtx, mock.Anything).Run(func(ctx context.Context, dogu *ecosystem.DoguInstallation) {
			assert.True(t, dogu.PauseReconciliation)
		}).Return(nil).Times(2)

		// when
		err := pauseReconciliationForDogus(testCtx, doguInstallaltionRepoMock, stateDiff)

		// then
		require.NoError(t, err)
	})

	t.Run("pause reconciliation for multiple dogus when global config changes", func(t *testing.T) {
		// given
		stateDiff := domain.StateDiff{
			DoguDiffs: []domain.DoguDiff{
				{
					DoguName:      redmine,
					NeededActions: []domain.Action{domain.ActionUpgrade},
				},
				{
					DoguName:      cas,
					NeededActions: []domain.Action{domain.ActionUpgrade},
				},
				{
					DoguName:      postfix,
					NeededActions: []domain.Action{}, // no action, so no pause required
				},
			},
			GlobalConfigDiffs: domain.GlobalConfigDiffs{
				getSetGlobalConfigEntryDiff("key", "value"),
			},
		}

		doguInstallaltionRepoMock := newMockDoguInstallationRepository(t)
		dogus := map[cescommons.SimpleName]*ecosystem.DoguInstallation{
			redmine: {Name: cescommons.QualifiedName{SimpleName: redmine, Namespace: "namespace"}},
			cas:     {Name: cescommons.QualifiedName{SimpleName: cas, Namespace: "namespace"}},
			postfix: {Name: cescommons.QualifiedName{SimpleName: postfix, Namespace: "namespace"}},
		}
		doguInstallaltionRepoMock.EXPECT().GetAll(testCtx).Return(dogus, nil)
		doguInstallaltionRepoMock.EXPECT().Update(testCtx, mock.Anything).Run(func(ctx context.Context, dogu *ecosystem.DoguInstallation) {
			assert.True(t, dogu.PauseReconciliation)
		}).Return(nil).Times(2)

		// when
		err := pauseReconciliationForDogus(testCtx, doguInstallaltionRepoMock, stateDiff)

		// then
		require.NoError(t, err)
	})

	t.Run("do not pause reconciliation for dogus without config changes", func(t *testing.T) {
		// given
		stateDiff := domain.StateDiff{
			DoguDiffs: []domain.DoguDiff{
				{
					DoguName:      redmine,
					NeededActions: []domain.Action{domain.ActionUpgrade},
				},
				{
					DoguName:      cas,
					NeededActions: []domain.Action{domain.ActionUpgrade},
				},
				{
					DoguName:      postfix,
					NeededActions: []domain.Action{}, // no action, so no pause required
				},
			},
		}

		doguInstallaltionRepoMock := newMockDoguInstallationRepository(t)
		dogus := map[cescommons.SimpleName]*ecosystem.DoguInstallation{
			redmine: {Name: cescommons.QualifiedName{SimpleName: redmine, Namespace: "namespace"}},
			cas:     {Name: cescommons.QualifiedName{SimpleName: cas, Namespace: "namespace"}},
			postfix: {Name: cescommons.QualifiedName{SimpleName: postfix, Namespace: "namespace"}},
		}
		doguInstallaltionRepoMock.EXPECT().GetAll(testCtx).Return(dogus, nil)
		// No Update calls

		// when
		err := pauseReconciliationForDogus(testCtx, doguInstallaltionRepoMock, stateDiff)

		// then
		require.NoError(t, err)
	})

	t.Run("error on get all dogus error", func(t *testing.T) {
		// given
		doguInstallaltionRepoMock := newMockDoguInstallationRepository(t)
		doguInstallaltionRepoMock.EXPECT().GetAll(testCtx).Return(nil, assert.AnError)

		// when
		err := pauseReconciliationForDogus(testCtx, doguInstallaltionRepoMock, domain.StateDiff{})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "error while attempting to load dogus")
	})

	t.Run("error on update error", func(t *testing.T) {
		// given
		stateDiff := domain.StateDiff{
			DoguDiffs: []domain.DoguDiff{
				{
					DoguName:      redmine,
					NeededActions: []domain.Action{domain.ActionUpgrade},
				},
			},
			GlobalConfigDiffs: domain.GlobalConfigDiffs{
				getSetGlobalConfigEntryDiff("key", "value"),
			},
		}
		dogus := map[cescommons.SimpleName]*ecosystem.DoguInstallation{
			redmine: {Name: cescommons.QualifiedName{SimpleName: redmine, Namespace: "namespace"}},
		}
		doguInstallaltionRepoMock := newMockDoguInstallationRepository(t)
		doguInstallaltionRepoMock.EXPECT().GetAll(testCtx).Return(dogus, nil)
		doguInstallaltionRepoMock.EXPECT().Update(testCtx, mock.Anything).Return(assert.AnError)

		// when
		err := pauseReconciliationForDogus(testCtx, doguInstallaltionRepoMock, stateDiff)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "could not pause reconciliation for dogu")
	})
}

func TestEcosystemConfigUseCase_applyDoguConfigDiffs(t *testing.T) {
	t.Run("should save diffs with action set", func(t *testing.T) {
		// given
		doguConfigMock := newMockDoguConfigRepository(t)

		diff1 := getSetDoguConfigEntryDiff("key1", "update1", redmine)
		diff2 := getSetDoguConfigEntryDiff("key2", "update2", redmine)
		diffsByDogu := map[cescommons.SimpleName]domain.DoguConfigDiffs{
			redmine: {
				diff1,
				diff2,
			},
		}

		redmineConfig := config.CreateDoguConfig(redmine, map[config.Key]config.Value{
			"key1": "val1",
			"key2": "val2",
		})

		// do not use redmineConfig here, because there is a bug in the k8s-registry lib
		// TODO: remove workaround when bug #50007 is fixed
		updatedConfig := config.CreateDoguConfig(redmine, map[config.Key]config.Value{
			"key1": "val1",
			"key2": "val2",
		}).Config
		updatedConfig, err := updatedConfig.Set(diff1.Key.Key, config.Value(*diff1.Expected.Value))
		require.NoError(t, err)
		updatedConfig, err = updatedConfig.Set(diff2.Key.Key, config.Value(*diff2.Expected.Value))
		require.NoError(t, err)

		doguConfigMock.EXPECT().
			GetAllExisting(testCtx, []cescommons.SimpleName{redmine}).
			Return(map[cescommons.SimpleName]config.DoguConfig{redmine: redmineConfig}, nil)
		doguConfigMock.EXPECT().
			UpdateOrCreate(testCtx, config.DoguConfig{DoguName: redmine, Config: updatedConfig}).
			Return(config.DoguConfig{}, nil)

		// when
		err = applyDoguConfigDiffs(testCtx, doguConfigMock, diffsByDogu)

		// then
		require.NoError(t, err)
	})

	t.Run("should delete diffs with action remove", func(t *testing.T) {
		// given
		doguConfigMock := newMockDoguConfigRepository(t)
		diff1 := getRemoveDoguConfigEntryDiff("key1", redmine)
		diff2 := getRemoveDoguConfigEntryDiff("key2", redmine)
		diffsByDogu := map[cescommons.SimpleName]domain.DoguConfigDiffs{
			redmine: {diff1, diff2},
		}

		redmineConfig := config.CreateDoguConfig(redmine, map[config.Key]config.Value{
			"key1": "val1",
			"key2": "val2",
		})

		//TODO: this fixes a bug #50007 in the lib: Delete modifies redmineConfig and updatedConfig as well. The structs are not really immutable at the moment.
		updatedConfig := config.CreateConfig(maps.Clone(redmineConfig.GetAll()))

		updatedConfig = updatedConfig.
			Delete("key1").
			Delete("key2")

		doguConfigMock.EXPECT().
			GetAllExisting(testCtx, []cescommons.SimpleName{redmine}).
			Return(map[cescommons.SimpleName]config.DoguConfig{redmine: redmineConfig}, nil)
		doguConfigMock.EXPECT().
			UpdateOrCreate(testCtx, config.DoguConfig{DoguName: redmine, Config: updatedConfig}).
			Return(config.DoguConfig{}, nil)

		// when
		err := applyDoguConfigDiffs(testCtx, doguConfigMock, diffsByDogu)

		// then
		require.NoError(t, err)
	})

	t.Run("should apply nothing on action none", func(t *testing.T) {
		// given
		doguConfigMock := newMockDoguConfigRepository(t)
		diff1 := domain.DoguConfigEntryDiff{
			NeededAction: domain.ConfigActionNone,
		}
		diffsByDogu := map[cescommons.SimpleName]domain.DoguConfigDiffs{
			redmine: {diff1},
		}

		doguConfigMock.EXPECT().
			GetAllExisting(testCtx, emptyDoguList).
			Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)

		// when
		err := applyDoguConfigDiffs(testCtx, doguConfigMock, diffsByDogu)

		// then
		require.NoError(t, err)
	})

	t.Run("err when GetAllExisting fails", func(t *testing.T) {
		// given
		doguConfigMock := newMockDoguConfigRepository(t)
		diff1 := getSetDoguConfigEntryDiff("key1", "value", redmine)
		diffsByDogu := map[cescommons.SimpleName]domain.DoguConfigDiffs{
			redmine: {diff1},
		}

		expectedError := liberrors.NewConnectionError(assert.AnError)
		doguConfigMock.EXPECT().
			GetAllExisting(testCtx, []cescommons.SimpleName{redmine}).
			Return(map[cescommons.SimpleName]config.DoguConfig{}, expectedError)

		// when
		err := applyDoguConfigDiffs(testCtx, doguConfigMock, diffsByDogu)

		// then
		require.Error(t, err)
		require.ErrorContains(t, err, expectedError.Error())
	})

	t.Run("error while applying key", func(t *testing.T) {
		// given
		doguConfigMock := newMockDoguConfigRepository(t)
		diff1 := getSetDoguConfigEntryDiff("key1/key1_1", "value", redmine)
		diffsByDogu := map[cescommons.SimpleName]domain.DoguConfigDiffs{
			redmine: {diff1},
		}

		redmineConfig := config.CreateDoguConfig(redmine, map[config.Key]config.Value{
			"key1": "val1",
			"key2": "val2",
		})

		doguConfigMock.EXPECT().
			GetAllExisting(testCtx, []cescommons.SimpleName{redmine}).
			Return(map[cescommons.SimpleName]config.DoguConfig{redmine: redmineConfig}, nil)

		// when
		err := applyDoguConfigDiffs(testCtx, doguConfigMock, diffsByDogu)

		// then
		assert.Error(t, err, "should throw an error when trying to create a sub key for an existing key")
		require.ErrorContains(t, err, "key key1 already has Value set") //error msg from registry-lib
	})
}

func TestEcosystemConfigUseCase_applyGlobalConfigDiffs(t *testing.T) {
	t.Run("should save diffs with action set", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigRepository(t)
		sut := NewEcosystemConfigUseCase(nil, nil, nil, globalConfigMock, nil)
		diff1 := getSetGlobalConfigEntryDiff("key1", "value1")
		diff2 := getSetGlobalConfigEntryDiff("key2", "value2")
		byAction := map[domain.ConfigAction][]domain.GlobalConfigEntryDiff{domain.ConfigActionSet: {diff1, diff2}}

		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)

		updatedEntries, err := globalConfig.Set(diff1.Key, common.GlobalConfigValue(*diff1.Expected.Value))
		require.NoError(t, err)
		updatedEntries, err = updatedEntries.Set(diff2.Key, common.GlobalConfigValue(*diff2.Expected.Value))
		require.NoError(t, err)

		globalConfigMock.EXPECT().Get(testCtx).Return(globalConfig, nil)
		globalConfigMock.EXPECT().Update(testCtx, config.GlobalConfig{Config: updatedEntries}).Return(globalConfig, nil)

		// when
		err = sut.applyGlobalConfigDiffs(testCtx, byAction)

		// then
		require.NoError(t, err)
	})

	t.Run("should delete diffs with action remove", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigRepository(t)
		sut := NewEcosystemConfigUseCase(nil, nil, nil, globalConfigMock, nil)
		diff1 := getRemoveGlobalConfigEntryDiff("key")
		diff2 := getRemoveGlobalConfigEntryDiff("key1")
		byAction := map[domain.ConfigAction][]domain.GlobalConfigEntryDiff{domain.ConfigActionRemove: {diff1, diff2}}

		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)

		updatedEntries := globalConfig.Delete(diff1.Key)
		updatedEntries = updatedEntries.Delete(diff2.Key)

		globalConfigMock.EXPECT().Get(testCtx).Return(globalConfig, nil)
		globalConfigMock.EXPECT().Update(testCtx, config.GlobalConfig{Config: updatedEntries}).Return(globalConfig, nil)

		// when
		err := sut.applyGlobalConfigDiffs(testCtx, byAction)

		// then
		require.NoError(t, err)
	})

	t.Run("should return nil on action none", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigRepository(t)
		sut := NewEcosystemConfigUseCase(nil, nil, nil, globalConfigMock, nil)
		diff1 := domain.GlobalConfigEntryDiff{
			NeededAction: domain.ConfigActionNone,
		}
		byAction := map[domain.ConfigAction][]domain.GlobalConfigEntryDiff{domain.ConfigActionNone: {diff1}}

		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)

		globalConfigMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		// when
		err := sut.applyGlobalConfigDiffs(testCtx, byAction)

		// then
		require.NoError(t, err)
	})

	t.Run("err when get fails", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigRepository(t)
		sut := NewEcosystemConfigUseCase(nil, nil, nil, globalConfigMock, nil)
		diff1 := domain.GlobalConfigEntryDiff{
			NeededAction: domain.ConfigActionSet,
		}
		byAction := map[domain.ConfigAction][]domain.GlobalConfigEntryDiff{domain.ConfigActionNone: {diff1}}

		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)

		expectedError := liberrors.NewConnectionError(assert.AnError)
		globalConfigMock.EXPECT().Get(testCtx).Return(globalConfig, expectedError)

		// when
		err := sut.applyGlobalConfigDiffs(testCtx, byAction)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, expectedError.Error())
	})
}

func TestEcosystemConfigUseCase_handleFailedApplyEcosystemConfig(t *testing.T) {
	t.Run("should set applied condition and event", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
		}
		blueprintRepoMock := newMockBlueprintSpecRepository(t)

		blueprintRepoMock.EXPECT().Update(testCtx, mock.IsType(&domain.BlueprintSpec{})).Return(nil)

		sut := EcosystemConfigUseCase{blueprintRepository: blueprintRepoMock}

		// when
		err := sut.handleFailedApplyEcosystemConfig(testCtx, blueprint, assert.AnError)

		// then
		require.Error(t, err)
		condition := meta.FindStatusCondition(blueprint.Conditions, domain.ConditionLastApplySucceeded)
		require.NotNil(t, condition)
		assert.Equal(t, err.Error(), condition.Message)
		assert.True(t, meta.IsStatusConditionFalse(blueprint.Conditions, domain.ConditionLastApplySucceeded))
		assert.IsType(t, domain.ExecutionFailedEvent{}, blueprint.Events[0])
	})

	t.Run("should return error on update error", func(t *testing.T) {
		// given
		spec := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
		}
		blueprintRepoMock := newMockBlueprintSpecRepository(t)

		blueprintRepoMock.EXPECT().Update(testCtx, mock.IsType(&domain.BlueprintSpec{})).Return(assert.AnError)

		sut := EcosystemConfigUseCase{blueprintRepository: blueprintRepoMock}

		// when
		err := sut.handleFailedApplyEcosystemConfig(testCtx, spec, assert.AnError)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "cannot mark blueprint config apply as failed")
		assert.ErrorContains(t, err, assert.AnError.Error())
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestNewEcosystemConfigUseCase(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		doguConfigMock := newMockDoguConfigRepository(t)
		sensitiveDoguConfigMock := newMockSensitiveDoguConfigRepository(t)
		globalConfigMock := newMockGlobalConfigRepository(t)

		// when
		useCase := NewEcosystemConfigUseCase(blueprintRepoMock, doguConfigMock, sensitiveDoguConfigMock, globalConfigMock, nil)

		// then
		assert.Equal(t, blueprintRepoMock, useCase.blueprintRepository)
		assert.Equal(t, doguConfigMock, useCase.doguConfigRepository)
		assert.Equal(t, sensitiveDoguConfigMock, useCase.sensitiveDoguConfigRepository)
		assert.Equal(t, globalConfigMock, useCase.globalConfigRepository)
	})
}

func getSetDoguConfigEntryDiff(key, value string, doguName cescommons.SimpleName) domain.DoguConfigEntryDiff {
	return domain.DoguConfigEntryDiff{
		Key: common.DoguConfigKey{
			Key:      config.Key(key),
			DoguName: doguName,
		},
		Expected: domain.DoguConfigValueState{
			Value: &value,
		},
		NeededAction: domain.ConfigActionSet,
	}
}

func getRemoveDoguConfigEntryDiff(key string, doguName cescommons.SimpleName) domain.DoguConfigEntryDiff {
	return domain.DoguConfigEntryDiff{
		Key: common.DoguConfigKey{
			Key:      config.Key(key),
			DoguName: doguName,
		},
		NeededAction: domain.ConfigActionRemove,
	}
}

func getSensitiveDoguConfigEntryDiffForAction(key, value string, doguName cescommons.SimpleName, action domain.ConfigAction) domain.SensitiveDoguConfigEntryDiff {
	return domain.SensitiveDoguConfigEntryDiff{
		Key: common.DoguConfigKey{
			Key:      config.Key(key),
			DoguName: doguName,
		},
		Expected: domain.DoguConfigValueState{
			Value: &value,
		},
		NeededAction: action,
	}
}

func getRemoveSensitiveDoguConfigEntryDiff(key string, doguName cescommons.SimpleName) domain.SensitiveDoguConfigEntryDiff {
	return domain.SensitiveDoguConfigEntryDiff{
		Key: common.DoguConfigKey{
			Key:      config.Key(key),
			DoguName: doguName,
		},
		NeededAction: domain.ConfigActionRemove,
	}
}

func getSetGlobalConfigEntryDiff(key, value string) domain.GlobalConfigEntryDiff {
	return domain.GlobalConfigEntryDiff{
		Key: common.GlobalConfigKey(key),
		Expected: domain.GlobalConfigValueState{
			Value: &value,
		},
		NeededAction: domain.ConfigActionSet,
	}
}

func getRemoveGlobalConfigEntryDiff(key string) domain.GlobalConfigEntryDiff {
	return domain.GlobalConfigEntryDiff{
		Key:          common.GlobalConfigKey(key),
		NeededAction: domain.ConfigActionRemove,
	}
}
