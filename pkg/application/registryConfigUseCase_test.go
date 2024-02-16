package application

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testSimpleDoguNameRedmine = common.SimpleDoguName("redmine")
	testSimpleDoguNameCas     = common.SimpleDoguName("cas")
	testBlueprintID           = "blueprint1"
)

func TestEcosystemRegistryUseCase_ApplyConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		doguConfigMock := newMockDoguConfigRepository(t)
		sensitiveDoguConfigMock := newMockDoguSensitiveConfigRepository(t)
		globalConfigMock := newMockGlobalConfigRepository(t)

		spec := &domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				DoguConfigDiff: map[common.SimpleDoguName]domain.CombinedDoguConfigDiff{
					testSimpleDoguNameRedmine: {
						DoguConfigDiff: []domain.DoguConfigEntryDiff{
							getSetDoguConfigEntryDiff("key", "value", testSimpleDoguNameRedmine),
							getRemoveDoguConfigEntryDiff("key", testSimpleDoguNameRedmine),
						},
						SensitiveDoguConfigDiff: []domain.SensitiveDoguConfigEntryDiff{
							getSetSensitiveDoguConfigEntryDiff("key", "value", testSimpleDoguNameRedmine),
							getRemoveSensitiveDoguConfigEntryDiff("key", testSimpleDoguNameRedmine),
						},
					},
					testSimpleDoguNameCas: {
						DoguConfigDiff: []domain.DoguConfigEntryDiff{
							getSetDoguConfigEntryDiff("key", "value", testSimpleDoguNameCas),
							getRemoveDoguConfigEntryDiff("key", testSimpleDoguNameCas),
						},
						SensitiveDoguConfigDiff: []domain.SensitiveDoguConfigEntryDiff{
							getSetSensitiveDoguConfigEntryDiff("key", "value", testSimpleDoguNameCas),
							getRemoveSensitiveDoguConfigEntryDiff("key", testSimpleDoguNameCas),
						},
					},
				},
				GlobalConfigDiff: domain.GlobalConfigDiff{
					getSetGlobalConfigEntryDiff("key", "value"),
					getRemoveGlobalConfigEntryDiff("key"),
				},
			},
		}

		// Just check if the routine hits the repos. Check values in concrete test of methods.
		doguConfigMock.EXPECT().Save(testCtx, mock.Anything).Return(nil).Times(2)
		doguConfigMock.EXPECT().Delete(testCtx, mock.Anything).Return(nil).Times(2)
		sensitiveDoguConfigMock.EXPECT().Save(testCtx, mock.Anything).Return(nil).Times(2)
		sensitiveDoguConfigMock.EXPECT().Delete(testCtx, mock.Anything).Return(nil).Times(2)
		globalConfigMock.EXPECT().Save(testCtx, mock.Anything).Return(nil).Times(1)
		globalConfigMock.EXPECT().Delete(testCtx, mock.Anything).Return(nil).Times(1)

		blueprintRepoMock.EXPECT().GetById(testCtx, testBlueprintID).Return(spec, nil)
		blueprintRepoMock.EXPECT().Update(testCtx, mock.Anything).Return(nil).Times(2)

		sut := EcosystemRegistryUseCase{blueprintRepository: blueprintRepoMock, doguConfigRepository: doguConfigMock, doguSensitiveConfigRepository: sensitiveDoguConfigMock, globalConfigRepository: globalConfigMock}

		// when
		err := sut.ApplyConfig(testCtx, testBlueprintID)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on fetch blueprint error", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)

		blueprintRepoMock.EXPECT().GetById(testCtx, testBlueprintID).Return(nil, assert.AnError)

		sut := EcosystemRegistryUseCase{blueprintRepository: blueprintRepoMock}

		// when
		err := sut.ApplyConfig(testCtx, testBlueprintID)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "cannot load blueprint to apply config")
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("mark applied if diffs are empty", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)

		spec := &domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				DoguConfigDiff:   map[common.SimpleDoguName]domain.CombinedDoguConfigDiff{},
				GlobalConfigDiff: domain.GlobalConfigDiff{},
			},
		}

		blueprintRepoMock.EXPECT().GetById(testCtx, testBlueprintID).Return(spec, nil)
		blueprintRepoMock.EXPECT().Update(testCtx, mock.Anything).Return(nil).Times(1)

		sut := EcosystemRegistryUseCase{blueprintRepository: blueprintRepoMock}

		// when
		err := sut.ApplyConfig(testCtx, testBlueprintID)

		// then
		require.NoError(t, err)
	})

	t.Run("should return on mark start error", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)

		spec := &domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				DoguConfigDiff: map[common.SimpleDoguName]domain.CombinedDoguConfigDiff{},
				GlobalConfigDiff: domain.GlobalConfigDiff{
					getSetGlobalConfigEntryDiff("key", "value"),
				},
			},
		}

		blueprintRepoMock.EXPECT().GetById(testCtx, testBlueprintID).Return(spec, nil)
		blueprintRepoMock.EXPECT().Update(testCtx, mock.Anything).Return(assert.AnError).Times(1)
		blueprintRepoMock.EXPECT().Update(testCtx, mock.Anything).Return(nil).Times(1)

		sut := EcosystemRegistryUseCase{blueprintRepository: blueprintRepoMock}

		// when
		err := sut.ApplyConfig(testCtx, testBlueprintID)

		// then
		require.NoError(t, err)
		assert.Equal(t, spec.Status, domain.StatusPhaseApplyRegistryConfigFailed)
	})

	t.Run("should join repo errors", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		doguConfigMock := newMockDoguConfigRepository(t)
		sensitiveDoguConfigMock := newMockDoguSensitiveConfigRepository(t)
		globalConfigMock := newMockGlobalConfigRepository(t)

		spec := &domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				DoguConfigDiff: map[common.SimpleDoguName]domain.CombinedDoguConfigDiff{
					testSimpleDoguNameRedmine: {
						DoguConfigDiff: []domain.DoguConfigEntryDiff{
							getSetDoguConfigEntryDiff("key", "value", testSimpleDoguNameRedmine),
						},
						SensitiveDoguConfigDiff: []domain.SensitiveDoguConfigEntryDiff{
							getSetSensitiveDoguConfigEntryDiff("key", "value", testSimpleDoguNameRedmine),
						},
					},
					testSimpleDoguNameCas: {
						DoguConfigDiff: []domain.DoguConfigEntryDiff{
							getSetDoguConfigEntryDiff("key", "value", testSimpleDoguNameCas),
						},
						SensitiveDoguConfigDiff: []domain.SensitiveDoguConfigEntryDiff{
							getSetSensitiveDoguConfigEntryDiff("key", "value", testSimpleDoguNameCas),
						},
					},
				},
				GlobalConfigDiff: domain.GlobalConfigDiff{
					getSetGlobalConfigEntryDiff("key", "value"),
				},
			},
		}

		// Just check if the routine hits the repos. Check values in concrete test of methods.
		doguConfigMock.EXPECT().Save(testCtx, mock.Anything).Return(assert.AnError).Times(2)
		sensitiveDoguConfigMock.EXPECT().Save(testCtx, mock.Anything).Return(assert.AnError).Times(2)
		globalConfigMock.EXPECT().Save(testCtx, mock.Anything).Return(assert.AnError).Times(1)

		blueprintRepoMock.EXPECT().GetById(testCtx, testBlueprintID).Return(spec, nil)
		blueprintRepoMock.EXPECT().Update(testCtx, mock.Anything).Return(nil).Times(2)

		sut := EcosystemRegistryUseCase{blueprintRepository: blueprintRepoMock, doguConfigRepository: doguConfigMock, doguSensitiveConfigRepository: sensitiveDoguConfigMock, globalConfigRepository: globalConfigMock}

		// when
		err := sut.ApplyConfig(testCtx, testBlueprintID)

		// then
		require.NoError(t, err)
		assert.Equal(t, spec.Status, domain.StatusPhaseApplyRegistryConfigFailed)
		assert.Len(t, spec.Events, 2)
		assert.Equal(t, spec.Events[1].Message(), "assert.AnError general error for testing\nassert.AnError general error for testing\nassert.AnError general error for testing\nassert.AnError general error for testing\nassert.AnError general error for testing")
	})
}

func TestEcosystemRegistryUseCase_applyDoguConfigDiffs(t *testing.T) {
	t.Run("should save diffs to with action set", func(t *testing.T) {
		// given
		doguConfigMock := newMockDoguConfigRepository(t)
		sut := NewEcosystemRegistryUseCase(nil, doguConfigMock, nil, nil)
		diff1 := getSetDoguConfigEntryDiff("/key", "value", testSimpleDoguNameRedmine)
		diff2 := getSetDoguConfigEntryDiff("/key1", "value1", testSimpleDoguNameRedmine)
		diffs := domain.DoguConfigDiff{diff1, diff2}

		expectedEntry1 := &ecosystem.DoguConfigEntry{
			Key:   common.DoguConfigKey{DoguName: testSimpleDoguNameRedmine, Key: diff1.Key.Key},
			Value: common.DoguConfigValue(diff1.Expected.Value),
		}
		expectedEntry2 := &ecosystem.DoguConfigEntry{
			Key:   common.DoguConfigKey{DoguName: testSimpleDoguNameRedmine, Key: diff2.Key.Key},
			Value: common.DoguConfigValue(diff2.Expected.Value),
		}

		doguConfigMock.EXPECT().Save(testCtx, expectedEntry1).Return(nil).Times(1)
		doguConfigMock.EXPECT().Save(testCtx, expectedEntry2).Return(nil).Times(1)

		// when
		err := sut.applyDoguConfigDiffs(testCtx, testSimpleDoguNameRedmine, diffs)

		// then
		require.NoError(t, err)
	})

	t.Run("should delete diffs with action remove", func(t *testing.T) {
		// given
		doguConfigMock := newMockDoguConfigRepository(t)
		sut := NewEcosystemRegistryUseCase(nil, doguConfigMock, nil, nil)
		diff1 := getRemoveDoguConfigEntryDiff("/key", testSimpleDoguNameRedmine)
		diff2 := getRemoveDoguConfigEntryDiff("/key1", testSimpleDoguNameRedmine)
		diffs := domain.DoguConfigDiff{diff1, diff2}

		expectedKey1 := common.DoguConfigKey{DoguName: testSimpleDoguNameRedmine, Key: diff1.Key.Key}
		expectedKey2 := common.DoguConfigKey{DoguName: testSimpleDoguNameRedmine, Key: diff2.Key.Key}

		doguConfigMock.EXPECT().Delete(testCtx, expectedKey1).Return(nil).Times(1)
		doguConfigMock.EXPECT().Delete(testCtx, expectedKey2).Return(nil).Times(1)

		// when
		err := sut.applyDoguConfigDiffs(testCtx, testSimpleDoguNameRedmine, diffs)

		// then
		require.NoError(t, err)
	})

	t.Run("should return nil on action none", func(t *testing.T) {
		// given
		doguConfigMock := newMockDoguConfigRepository(t)
		sut := NewEcosystemRegistryUseCase(nil, doguConfigMock, nil, nil)
		diff1 := domain.DoguConfigEntryDiff{
			Action: domain.ConfigActionNone,
		}

		diffs := domain.DoguConfigDiff{diff1}

		// when
		err := sut.applyDoguConfigDiffs(testCtx, testSimpleDoguNameRedmine, diffs)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on unknown action", func(t *testing.T) {
		// given
		sut := NewEcosystemRegistryUseCase(nil, nil, nil, nil)
		diff1 := domain.DoguConfigEntryDiff{
			Key:    common.DoguConfigKey{Key: "key"},
			Action: "unknown",
		}

		diffs := domain.DoguConfigDiff{diff1}

		// when
		err := sut.applyDoguConfigDiffs(testCtx, testSimpleDoguNameRedmine, diffs)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "cannot perform unknown action \"unknown\" for dogu \"redmine\" with key \"key\"")
	})
}

func TestEcosystemRegistryUseCase_applyGlobalConfigDiffs(t *testing.T) {
	t.Run("should save diffs to with action set", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigRepository(t)
		sut := NewEcosystemRegistryUseCase(nil, nil, nil, globalConfigMock)
		diff1 := getSetGlobalConfigEntryDiff("/key", "value")
		diff2 := getSetGlobalConfigEntryDiff("/key1", "value1")
		diffs := domain.GlobalConfigDiff{diff1, diff2}

		expectedEntry1 := &ecosystem.GlobalConfigEntry{
			Key:   diff1.Key,
			Value: common.GlobalConfigValue(diff1.Expected.Value),
		}
		expectedEntry2 := &ecosystem.GlobalConfigEntry{
			Key:   diff2.Key,
			Value: common.GlobalConfigValue(diff2.Expected.Value),
		}

		globalConfigMock.EXPECT().Save(testCtx, expectedEntry1).Return(nil).Times(1)
		globalConfigMock.EXPECT().Save(testCtx, expectedEntry2).Return(nil).Times(1)

		// when
		err := sut.applyGlobalConfigDiffs(testCtx, diffs)

		// then
		require.NoError(t, err)
	})

	t.Run("should delete diffs with action remove", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigRepository(t)
		sut := NewEcosystemRegistryUseCase(nil, nil, nil, globalConfigMock)
		diff1 := getRemoveGlobalConfigEntryDiff("/key")
		diff2 := getRemoveGlobalConfigEntryDiff("/key1")
		diffs := domain.GlobalConfigDiff{diff1, diff2}

		globalConfigMock.EXPECT().Delete(testCtx, diff1.Key).Return(nil).Times(1)
		globalConfigMock.EXPECT().Delete(testCtx, diff2.Key).Return(nil).Times(1)

		// when
		err := sut.applyGlobalConfigDiffs(testCtx, diffs)

		// then
		require.NoError(t, err)
	})

	t.Run("should return nil on action none", func(t *testing.T) {
		// given
		sut := NewEcosystemRegistryUseCase(nil, nil, nil, nil)
		diff1 := domain.GlobalConfigEntryDiff{
			Action: domain.ConfigActionNone,
		}

		diffs := domain.GlobalConfigDiff{diff1}

		// when
		err := sut.applyGlobalConfigDiffs(testCtx, diffs)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on unknown action", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigRepository(t)
		sut := NewEcosystemRegistryUseCase(nil, nil, nil, globalConfigMock)
		diff1 := domain.GlobalConfigEntryDiff{
			Key:    "key",
			Action: "unknown",
		}

		diffs := domain.GlobalConfigDiff{diff1}

		// when
		err := sut.applyGlobalConfigDiffs(testCtx, diffs)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "cannot perform unknown action \"unknown\" for global config with key \"key\"")
	})
}

func TestEcosystemRegistryUseCase_applySensitiveDoguConfigDiffs(t *testing.T) {
	t.Run("should save diffs to with action set", func(t *testing.T) {
		// given
		sensitiveDoguConfigMock := newMockDoguSensitiveConfigRepository(t)
		sut := NewEcosystemRegistryUseCase(nil, nil, sensitiveDoguConfigMock, nil)
		diff1 := getSetSensitiveDoguConfigEntryDiff("key", "value", testSimpleDoguNameRedmine)
		diff2 := getSetSensitiveDoguConfigEntryDiff("key1", "value1", testSimpleDoguNameRedmine)
		diffs := domain.SensitiveDoguConfigDiff{diff1, diff2}

		expectedEntry1 := &ecosystem.SensitiveDoguConfigEntry{
			Key:   common.SensitiveDoguConfigKey{DoguName: testSimpleDoguNameRedmine, Key: diff1.Key.Key},
			Value: common.EncryptedDoguConfigValue(diff1.Expected.Value),
		}
		expectedEntry2 := &ecosystem.SensitiveDoguConfigEntry{
			Key:   common.SensitiveDoguConfigKey{DoguName: testSimpleDoguNameRedmine, Key: diff2.Key.Key},
			Value: common.EncryptedDoguConfigValue(diff2.Expected.Value),
		}

		sensitiveDoguConfigMock.EXPECT().Save(testCtx, expectedEntry1).Return(nil).Times(1)
		sensitiveDoguConfigMock.EXPECT().Save(testCtx, expectedEntry2).Return(nil).Times(1)

		// when
		err := sut.applySensitiveDoguConfigDiffs(testCtx, testSimpleDoguNameRedmine, diffs)

		// then
		require.NoError(t, err)
	})

	t.Run("should delete diffs with action remove", func(t *testing.T) {
		// given
		sensitiveDoguConfigMock := newMockDoguSensitiveConfigRepository(t)
		sut := NewEcosystemRegistryUseCase(nil, nil, sensitiveDoguConfigMock, nil)
		diff1 := getRemoveSensitiveDoguConfigEntryDiff("key", testSimpleDoguNameRedmine)
		diff2 := getRemoveSensitiveDoguConfigEntryDiff("key", testSimpleDoguNameRedmine)
		diffs := domain.SensitiveDoguConfigDiff{diff1, diff2}

		expectedKey1 := common.SensitiveDoguConfigKey{DoguName: testSimpleDoguNameRedmine, Key: diff1.Key.Key}
		expectedKey2 := common.SensitiveDoguConfigKey{DoguName: testSimpleDoguNameRedmine, Key: diff2.Key.Key}

		sensitiveDoguConfigMock.EXPECT().Delete(testCtx, expectedKey1).Return(nil).Times(1)
		sensitiveDoguConfigMock.EXPECT().Delete(testCtx, expectedKey2).Return(nil).Times(1)

		// when
		err := sut.applySensitiveDoguConfigDiffs(testCtx, testSimpleDoguNameRedmine, diffs)

		// then
		require.NoError(t, err)
	})

	t.Run("should return nil on action none", func(t *testing.T) {
		// given
		sut := NewEcosystemRegistryUseCase(nil, nil, nil, nil)
		diff1 := domain.SensitiveDoguConfigEntryDiff{
			Action: domain.ConfigActionNone,
		}

		diffs := domain.SensitiveDoguConfigDiff{diff1}

		// when
		err := sut.applySensitiveDoguConfigDiffs(testCtx, testSimpleDoguNameRedmine, diffs)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on unknown action", func(t *testing.T) {
		// given
		sensitiveDoguConfigMock := newMockDoguSensitiveConfigRepository(t)
		sut := NewEcosystemRegistryUseCase(nil, nil, sensitiveDoguConfigMock, nil)
		diff1 := domain.SensitiveDoguConfigEntryDiff{
			Key:    common.SensitiveDoguConfigKey{Key: "key"},
			Action: "unknown",
		}

		diffs := domain.SensitiveDoguConfigDiff{diff1}

		// when
		err := sut.applySensitiveDoguConfigDiffs(testCtx, testSimpleDoguNameRedmine, diffs)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "cannot perform unknown action \"unknown\" for dogu \"redmine\" with key \"key\"")
	})
}

func TestEcosystemRegistryUseCase_markConfigApplied(t *testing.T) {
	t.Run("should set applied status and event", func(t *testing.T) {
		// given
		spec := &domain.BlueprintSpec{}
		expectedSpec := &domain.BlueprintSpec{}
		expectedSpec.Status = domain.StatusPhaseRegistryConfigApplied
		expectedSpec.Events = append(spec.Events, domain.RegistryConfigAppliedEvent{})
		blueprintRepoMock := newMockBlueprintSpecRepository(t)

		blueprintRepoMock.EXPECT().Update(testCtx, expectedSpec).Return(nil)

		sut := EcosystemRegistryUseCase{blueprintRepository: blueprintRepoMock}

		// when
		err := sut.markConfigApplied(testCtx, spec)

		// then
		require.NoError(t, err)
	})

	t.Run("should return an error on update error", func(t *testing.T) {
		// given
		spec := &domain.BlueprintSpec{}
		expectedSpec := &domain.BlueprintSpec{}
		expectedSpec.Status = domain.StatusPhaseRegistryConfigApplied
		expectedSpec.Events = append(spec.Events, domain.RegistryConfigAppliedEvent{})
		blueprintRepoMock := newMockBlueprintSpecRepository(t)

		blueprintRepoMock.EXPECT().Update(testCtx, expectedSpec).Return(assert.AnError)

		sut := EcosystemRegistryUseCase{blueprintRepository: blueprintRepoMock}

		// when
		err := sut.markConfigApplied(testCtx, spec)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to mark registry config applied")
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestEcosystemRegistryUseCase_markApplyConfigStart(t *testing.T) {
	t.Run("should set status and event apply config", func(t *testing.T) {
		// given
		spec := &domain.BlueprintSpec{}
		expectedSpec := &domain.BlueprintSpec{}
		expectedSpec.Status = domain.StatusPhaseApplyRegistryConfig
		expectedSpec.Events = append(spec.Events, domain.ApplyRegistryConfigEvent{})
		blueprintRepoMock := newMockBlueprintSpecRepository(t)

		blueprintRepoMock.EXPECT().Update(testCtx, expectedSpec).Return(nil)

		sut := EcosystemRegistryUseCase{blueprintRepository: blueprintRepoMock}

		// when
		err := sut.markApplyConfigStart(testCtx, spec)

		// then
		require.NoError(t, err)
	})

	t.Run("should return an error on update error", func(t *testing.T) {
		// given
		spec := &domain.BlueprintSpec{}
		expectedSpec := &domain.BlueprintSpec{}
		expectedSpec.Status = domain.StatusPhaseApplyRegistryConfig
		expectedSpec.Events = append(spec.Events, domain.ApplyRegistryConfigEvent{})
		blueprintRepoMock := newMockBlueprintSpecRepository(t)

		blueprintRepoMock.EXPECT().Update(testCtx, expectedSpec).Return(assert.AnError)

		sut := EcosystemRegistryUseCase{blueprintRepository: blueprintRepoMock}

		// when
		err := sut.markApplyConfigStart(testCtx, spec)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "cannot mark blueprint as applying config")
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestEcosystemRegistryUseCase_handleFailedApplyRegistryConfig(t *testing.T) {
	t.Run("should set applied status and event", func(t *testing.T) {
		// given
		spec := &domain.BlueprintSpec{}
		blueprintRepoMock := newMockBlueprintSpecRepository(t)

		blueprintRepoMock.EXPECT().Update(testCtx, mock.IsType(&domain.BlueprintSpec{})).Return(nil)

		sut := EcosystemRegistryUseCase{blueprintRepository: blueprintRepoMock}

		// when
		err := sut.handleFailedApplyRegistryConfig(testCtx, spec, assert.AnError)

		// then
		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseApplyRegistryConfigFailed, spec.Status)
		assert.IsType(t, domain.ApplyRegistryConfigFailedEvent{}, spec.Events[0])
	})

	t.Run("should return error on update error", func(t *testing.T) {
		// given
		spec := &domain.BlueprintSpec{}
		blueprintRepoMock := newMockBlueprintSpecRepository(t)

		blueprintRepoMock.EXPECT().Update(testCtx, mock.IsType(&domain.BlueprintSpec{})).Return(assert.AnError)

		sut := EcosystemRegistryUseCase{blueprintRepository: blueprintRepoMock}

		// when
		err := sut.handleFailedApplyRegistryConfig(testCtx, spec, assert.AnError)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "cannot mark blueprint config apply as failed while handling \"applyRegistryConfigFailed\" status")
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestNewEcosystemRegistryUseCase(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		doguConfigMock := newMockDoguConfigRepository(t)
		sensitiveDoguConfigMock := newMockDoguSensitiveConfigRepository(t)
		globalConfigMock := newMockGlobalConfigRepository(t)

		// when
		useCase := NewEcosystemRegistryUseCase(blueprintRepoMock, doguConfigMock, sensitiveDoguConfigMock, globalConfigMock)

		// then
		assert.Equal(t, blueprintRepoMock, useCase.blueprintRepository)
		assert.Equal(t, doguConfigMock, useCase.doguConfigRepository)
		assert.Equal(t, sensitiveDoguConfigMock, useCase.doguSensitiveConfigRepository)
		assert.Equal(t, globalConfigMock, useCase.globalConfigRepository)
	})
}

func getSetDoguConfigEntryDiff(key, value string, doguName common.SimpleDoguName) domain.DoguConfigEntryDiff {
	return domain.DoguConfigEntryDiff{
		Key: common.DoguConfigKey{
			Key:      key,
			DoguName: doguName,
		},
		Expected: domain.DoguConfigValueState{
			Value: value,
		},
		Action: domain.ConfigActionSet,
	}
}

func getRemoveDoguConfigEntryDiff(key string, doguName common.SimpleDoguName) domain.DoguConfigEntryDiff {
	return domain.DoguConfigEntryDiff{
		Key: common.DoguConfigKey{
			Key:      key,
			DoguName: doguName,
		},
		Action: domain.ConfigActionRemove,
	}
}

func getSetSensitiveDoguConfigEntryDiff(key, value string, doguName common.SimpleDoguName) domain.SensitiveDoguConfigEntryDiff {
	return domain.SensitiveDoguConfigEntryDiff{
		Key: common.SensitiveDoguConfigKey{
			Key:      key,
			DoguName: doguName,
		},
		Expected: domain.EncryptedDoguConfigValueState{
			Value: value,
		},
		Action: domain.ConfigActionSet,
	}
}

func getRemoveSensitiveDoguConfigEntryDiff(key string, doguName common.SimpleDoguName) domain.SensitiveDoguConfigEntryDiff {
	return domain.SensitiveDoguConfigEntryDiff{
		Key: common.SensitiveDoguConfigKey{
			Key:      key,
			DoguName: doguName,
		},
		Action: domain.ConfigActionRemove,
	}
}

func getSetGlobalConfigEntryDiff(key, value string) domain.GlobalConfigEntryDiff {
	return domain.GlobalConfigEntryDiff{
		Key: common.GlobalConfigKey(key),
		Expected: domain.GlobalConfigValueState{
			Value: value,
		},
		Action: domain.ConfigActionSet,
	}
}

func getRemoveGlobalConfigEntryDiff(key string) domain.GlobalConfigEntryDiff {
	return domain.GlobalConfigEntryDiff{
		Key:    common.GlobalConfigKey(key),
		Action: domain.ConfigActionRemove,
	}
}
