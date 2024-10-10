package application

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"maps"
	"testing"
)

const (
	redmine         = common.SimpleDoguName("redmine")
	cas             = common.SimpleDoguName("cas")
	testBlueprintID = "blueprint1"
)

func TestEcosystemConfigUseCase_ApplyConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		doguConfigMock := newMockDoguConfigRepository(t)
		sensitiveDoguConfigMock := newMockSensitiveDoguConfigEntryRepository(t)
		globalConfigRepoMock := newMockGlobalConfigRepository(t)

		redmineDiffToEncrypt := getSensitiveDoguConfigEntryDiffForAction("key", "value", redmine, domain.ConfigActionSet)
		casDiffToEncrypt := getSensitiveDoguConfigEntryDiffForAction("key", "value", cas, domain.ConfigActionSet)
		spec := &domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				DoguConfigDiffs: map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{
					redmine: {
						DoguConfigDiff: []domain.DoguConfigEntryDiff{
							getSetDoguConfigEntryDiff("key", "value", redmine),
							getRemoveDoguConfigEntryDiff("key", redmine),
						},
						SensitiveDoguConfigDiff: []domain.SensitiveDoguConfigEntryDiff{
							redmineDiffToEncrypt,
							getRemoveSensitiveDoguConfigEntryDiff("key", redmine),
						},
					},
					cas: {
						DoguConfigDiff: []domain.DoguConfigEntryDiff{
							getSetDoguConfigEntryDiff("key", "value", cas),
							getRemoveDoguConfigEntryDiff("key", cas),
						},
						SensitiveDoguConfigDiff: []domain.SensitiveDoguConfigEntryDiff{
							casDiffToEncrypt,
							getRemoveSensitiveDoguConfigEntryDiff("key", cas),
						},
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
			GetAll(testCtx, []common.SimpleDoguName{cas, redmine}).
			Return(map[config.SimpleDoguName]config.DoguConfig{
				redmine: config.CreateDoguConfig(redmine, map[config.Key]config.Value{}),
				cas:     config.CreateDoguConfig(cas, map[config.Key]config.Value{}),
			}, nil)
		doguConfigMock.EXPECT().Update(testCtx, mock.Anything).Return(config.DoguConfig{}, nil).Times(2)

		sensitiveDoguConfigMock.EXPECT().SaveAll(testCtx, mock.Anything).Return(nil).Times(1)
		sensitiveDoguConfigMock.EXPECT().DeleteAllByKeys(testCtx, mock.Anything).Return(nil).Times(1)

		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)
		globalConfigRepoMock.EXPECT().Update(testCtx, mock.Anything).Return(globalConfig, nil)

		blueprintRepoMock.EXPECT().GetById(testCtx, testBlueprintID).Return(spec, nil)
		blueprintRepoMock.EXPECT().Update(testCtx, mock.Anything).Return(nil).Times(2)

		sut := NewEcosystemConfigUseCase(blueprintRepoMock, doguConfigMock, sensitiveDoguConfigMock, globalConfigRepoMock)

		// when
		err := sut.ApplyConfig(testCtx, testBlueprintID)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on fetch blueprint error", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)

		blueprintRepoMock.EXPECT().GetById(testCtx, testBlueprintID).Return(nil, assert.AnError)

		sut := EcosystemConfigUseCase{blueprintRepository: blueprintRepoMock}

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
				DoguConfigDiffs:   map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{},
				GlobalConfigDiffs: domain.GlobalConfigDiffs{},
			},
		}

		blueprintRepoMock.EXPECT().GetById(testCtx, testBlueprintID).Return(spec, nil)
		blueprintRepoMock.EXPECT().Update(testCtx, mock.Anything).Return(nil).Times(1)

		sut := EcosystemConfigUseCase{blueprintRepository: blueprintRepoMock}

		// when
		err := sut.ApplyConfig(testCtx, testBlueprintID)

		// then
		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseRegistryConfigApplied, spec.Status)
	})

	t.Run("should return on mark apply config start error", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)

		spec := &domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				DoguConfigDiffs: map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{},
				GlobalConfigDiffs: domain.GlobalConfigDiffs{
					getSetGlobalConfigEntryDiff("key", "value"),
				},
			},
		}

		blueprintRepoMock.EXPECT().GetById(testCtx, testBlueprintID).Return(spec, nil)
		blueprintRepoMock.EXPECT().Update(testCtx, mock.Anything).Return(assert.AnError).Times(1)
		blueprintRepoMock.EXPECT().Update(testCtx, mock.Anything).Return(nil).Times(1)

		sut := EcosystemConfigUseCase{blueprintRepository: blueprintRepoMock}

		// when
		err := sut.ApplyConfig(testCtx, testBlueprintID)

		// then
		require.NoError(t, err)
		assert.Equal(t, spec.Status, domain.StatusPhaseApplyRegistryConfigFailed)
	})

	t.Run("error applying dogu config", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		doguConfigMock := newMockDoguConfigRepository(t)
		sensitiveDoguConfigMock := newMockSensitiveDoguConfigEntryRepository(t)
		globalConfigMock := newMockGlobalConfigRepository(t)

		spec := &domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				DoguConfigDiffs: map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{
					redmine: {
						DoguConfigDiff: []domain.DoguConfigEntryDiff{
							getSetDoguConfigEntryDiff("key", "value", redmine),
						},
					},
					cas: {
						DoguConfigDiff: []domain.DoguConfigEntryDiff{
							getSetDoguConfigEntryDiff("key", "value", cas),
						},
					},
				},
			},
		}

		// Just check if the routine hits the repos. Check values in concrete test of methods.
		doguConfigMock.EXPECT().
			GetAll(testCtx, []common.SimpleDoguName{cas, redmine}).
			Return(map[config.SimpleDoguName]config.DoguConfig{
				redmine: config.CreateDoguConfig(redmine, map[config.Key]config.Value{}),
				cas:     config.CreateDoguConfig(cas, map[config.Key]config.Value{}),
			}, nil)
		doguConfigMock.EXPECT().Update(testCtx, mock.Anything).Return(config.DoguConfig{}, assert.AnError).Times(1)

		blueprintRepoMock.EXPECT().GetById(testCtx, testBlueprintID).Return(spec, nil)
		blueprintRepoMock.EXPECT().Update(testCtx, mock.Anything).Return(nil).Times(2)

		sut := NewEcosystemConfigUseCase(blueprintRepoMock, doguConfigMock, sensitiveDoguConfigMock, globalConfigMock)

		// when
		err := sut.ApplyConfig(testCtx, testBlueprintID)

		// then
		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseApplyRegistryConfigFailed, spec.Status)
		assert.Len(t, spec.Events, 2)
		assert.Equal(t, "assert.AnError general error for testing", spec.Events[1].Message())
	})
	t.Run("error applying sensitive config", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		doguConfigMock := newMockDoguConfigRepository(t)
		sensitiveDoguConfigMock := newMockSensitiveDoguConfigEntryRepository(t)
		globalConfigMock := newMockGlobalConfigRepository(t)

		casDiff := getSensitiveDoguConfigEntryDiffForAction("key", "value", cas, domain.ConfigActionSet)
		spec := &domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				DoguConfigDiffs: map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{
					redmine: {
						SensitiveDoguConfigDiff: []domain.SensitiveDoguConfigEntryDiff{
							getSensitiveDoguConfigEntryDiffForAction("key", "value", redmine, domain.ConfigActionSet),
						},
					},
					cas: {
						SensitiveDoguConfigDiff: []domain.SensitiveDoguConfigEntryDiff{
							casDiff,
						},
					},
				},
			},
		}

		// Just check if the routine hits the repos. Check values in concrete test of methods.
		sensitiveDoguConfigMock.EXPECT().SaveAll(testCtx, mock.Anything).Return(assert.AnError)

		blueprintRepoMock.EXPECT().GetById(testCtx, testBlueprintID).Return(spec, nil)
		blueprintRepoMock.EXPECT().Update(testCtx, mock.Anything).Return(nil).Times(2)

		sut := NewEcosystemConfigUseCase(blueprintRepoMock, doguConfigMock, sensitiveDoguConfigMock, globalConfigMock)

		// when
		err := sut.ApplyConfig(testCtx, testBlueprintID)

		// then
		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseApplyRegistryConfigFailed, spec.Status)
		assert.Len(t, spec.Events, 2)
		assert.Equal(t, "assert.AnError general error for testing", spec.Events[1].Message())
	})
	t.Run("error applying global config", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		doguConfigMock := newMockDoguConfigRepository(t)
		sensitiveDoguConfigMock := newMockSensitiveDoguConfigEntryRepository(t)
		globalConfigMock := newMockGlobalConfigRepository(t)

		spec := &domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				DoguConfigDiffs: map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{},
				GlobalConfigDiffs: domain.GlobalConfigDiffs{
					getSetGlobalConfigEntryDiff("key", "value"),
				},
			},
		}

		// Just check if the routine hits the repos. Check values in concrete test of methods.

		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigMock.EXPECT().Get(testCtx).Return(globalConfig, nil)
		globalConfigMock.EXPECT().Update(testCtx, mock.Anything).Return(globalConfig, assert.AnError)

		blueprintRepoMock.EXPECT().GetById(testCtx, testBlueprintID).Return(spec, nil)
		blueprintRepoMock.EXPECT().Update(testCtx, mock.Anything).Return(nil).Times(2)

		sut := NewEcosystemConfigUseCase(blueprintRepoMock, doguConfigMock, sensitiveDoguConfigMock, globalConfigMock)

		// when
		err := sut.ApplyConfig(testCtx, testBlueprintID)

		// then
		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseApplyRegistryConfigFailed, spec.Status)
		assert.Len(t, spec.Events, 2)
		assert.Equal(t, "assert.AnError general error for testing", spec.Events[1].Message())
	})
}

func TestEcosystemConfigUseCase_applyDoguConfigDiffs(t *testing.T) {
	t.Run("should save diffs with action set", func(t *testing.T) {
		// given
		doguConfigMock := newMockDoguConfigRepository(t)
		sut := NewEcosystemConfigUseCase(nil, doguConfigMock, nil, nil)
		diff1 := getSetDoguConfigEntryDiff("key1", "update1", redmine)
		diff2 := getSetDoguConfigEntryDiff("key2", "update2", redmine)
		diffsByDogu := map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{
			"redmine": {DoguConfigDiff: []domain.DoguConfigEntryDiff{diff1, diff2}},
		}

		redmineConfig := config.CreateDoguConfig("redmine", map[config.Key]config.Value{
			"key1": "val1",
			"key2": "val2",
		})

		updatedConfig := redmineConfig.Config
		updatedConfig, err := updatedConfig.Set(diff1.Key.Key, config.Value(diff1.Expected.Value))
		require.NoError(t, err)
		updatedConfig, err = updatedConfig.Set(diff2.Key.Key, config.Value(diff2.Expected.Value))
		require.NoError(t, err)

		doguConfigMock.EXPECT().
			GetAll(testCtx, []common.SimpleDoguName{redmine}).
			Return(map[config.SimpleDoguName]config.DoguConfig{redmine: redmineConfig}, nil)
		doguConfigMock.EXPECT().
			Update(testCtx, config.DoguConfig{DoguName: redmine, Config: updatedConfig}).
			Return(config.DoguConfig{}, nil)

		// when
		err = sut.applyDoguConfigDiffs(testCtx, diffsByDogu)

		// then
		require.NoError(t, err)
	})

	t.Run("should delete diffs with action remove", func(t *testing.T) {
		// given
		doguConfigMock := newMockDoguConfigRepository(t)
		sut := NewEcosystemConfigUseCase(nil, doguConfigMock, nil, nil)
		diff1 := getRemoveDoguConfigEntryDiff("key1", redmine)
		diff2 := getRemoveDoguConfigEntryDiff("key2", redmine)
		diffsByDogu := map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{
			"redmine": {DoguConfigDiff: []domain.DoguConfigEntryDiff{diff1, diff2}},
		}

		redmineConfig := config.CreateDoguConfig("redmine", map[config.Key]config.Value{
			"key1": "val1",
			"key2": "val2",
		})

		//TODO: this fixes a bug in the lib: Delete modifies redmineConfig and updatedConfig as well. The structs are not really immutable at the moment.
		updatedConfig := config.CreateConfig(maps.Clone(redmineConfig.GetAll()))

		updatedConfig = updatedConfig.
			Delete("key1").
			Delete("key2")

		doguConfigMock.EXPECT().
			GetAll(testCtx, []common.SimpleDoguName{"redmine"}).
			Return(map[config.SimpleDoguName]config.DoguConfig{"redmine": redmineConfig}, nil)
		doguConfigMock.EXPECT().
			Update(testCtx, config.DoguConfig{DoguName: redmine, Config: updatedConfig}).
			Return(config.DoguConfig{}, nil)

		// when
		err := sut.applyDoguConfigDiffs(testCtx, diffsByDogu)

		// then
		require.NoError(t, err)
	})

	t.Run("should return nil on action none", func(t *testing.T) {
		// given
		doguConfigMock := newMockDoguConfigRepository(t)
		sut := NewEcosystemConfigUseCase(nil, doguConfigMock, nil, nil)
		diff1 := domain.DoguConfigEntryDiff{
			NeededAction: domain.ConfigActionNone,
		}
		diffsByDogu := map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{
			"redmine": {DoguConfigDiff: []domain.DoguConfigEntryDiff{diff1}},
		}

		// when
		err := sut.applyDoguConfigDiffs(testCtx, diffsByDogu)

		// then
		require.NoError(t, err)
	})
}

func TestEcosystemConfigUseCase_applyGlobalConfigDiffs(t *testing.T) {
	t.Run("should save diffs with action set", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigRepository(t)
		sut := NewEcosystemConfigUseCase(nil, nil, nil, globalConfigMock)
		diff1 := getSetGlobalConfigEntryDiff("key1", "value1")
		diff2 := getSetGlobalConfigEntryDiff("key2", "value2")
		byAction := map[domain.ConfigAction][]domain.GlobalConfigEntryDiff{domain.ConfigActionSet: {diff1, diff2}}

		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)

		updatedEntries, err := globalConfig.Set(diff1.Key, common.GlobalConfigValue(diff1.Expected.Value))
		require.NoError(t, err)
		updatedEntries, err = updatedEntries.Set(diff2.Key, common.GlobalConfigValue(diff2.Expected.Value))
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
		sut := NewEcosystemConfigUseCase(nil, nil, nil, globalConfigMock)
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
		sut := NewEcosystemConfigUseCase(nil, nil, nil, globalConfigMock)
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
}

func TestEcosystemConfigUseCase_applySensitiveDoguConfigDiffs(t *testing.T) {
	t.Run("should save diffs with action set", func(t *testing.T) {
		// given
		sensitiveDoguConfigMock := newMockSensitiveDoguConfigEntryRepository(t)
		sut := NewEcosystemConfigUseCase(nil, nil, sensitiveDoguConfigMock, nil)
		diff1 := getSensitiveDoguConfigEntryDiffForAction("key1", "value1", redmine, domain.ConfigActionSet)
		diff2 := getSensitiveDoguConfigEntryDiffForAction("key2", "value2", redmine, domain.ConfigActionSet)
		byAction := map[domain.ConfigAction]domain.SensitiveDoguConfigDiffs{domain.ConfigActionSet: {diff1, diff2}}

		expectedEntry1 := &ecosystem.SensitiveDoguConfigEntry{
			Key:   common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{DoguName: redmine, Key: diff1.Key.Key}},
			Value: common.SensitiveDoguConfigValue("value1"),
		}
		expectedEntry2 := &ecosystem.SensitiveDoguConfigEntry{
			Key:   common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{DoguName: redmine, Key: diff2.Key.Key}},
			Value: common.SensitiveDoguConfigValue("value2"),
		}

		sensitiveDoguConfigMock.EXPECT().SaveAll(testCtx, []*ecosystem.SensitiveDoguConfigEntry{expectedEntry1, expectedEntry2}).Return(nil).Times(1)

		// when
		err := sut.applySensitiveDoguConfigDiffs(testCtx, byAction)

		// then
		require.NoError(t, err)
	})

	t.Run("should delete diffs with action remove", func(t *testing.T) {
		// given
		sensitiveDoguConfigMock := newMockSensitiveDoguConfigEntryRepository(t)
		sut := NewEcosystemConfigUseCase(nil, nil, sensitiveDoguConfigMock, nil)
		diff1 := getRemoveSensitiveDoguConfigEntryDiff("key", redmine)
		diff2 := getRemoveSensitiveDoguConfigEntryDiff("key", redmine)
		byAction := map[domain.ConfigAction]domain.SensitiveDoguConfigDiffs{domain.ConfigActionRemove: {diff1, diff2}}

		expectedKey1 := common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{DoguName: redmine, Key: diff1.Key.Key}}
		expectedKey2 := common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{DoguName: redmine, Key: diff2.Key.Key}}

		sensitiveDoguConfigMock.EXPECT().DeleteAllByKeys(testCtx, []common.SensitiveDoguConfigKey{expectedKey1, expectedKey2}).Return(nil).Times(1)

		// when
		err := sut.applySensitiveDoguConfigDiffs(testCtx, byAction)

		// then
		require.NoError(t, err)
	})

	t.Run("should return nil on action none", func(t *testing.T) {
		// given
		sut := NewEcosystemConfigUseCase(nil, nil, newMockSensitiveDoguConfigEntryRepository(t), nil)
		diff1 := domain.SensitiveDoguConfigEntryDiff{
			NeededAction: domain.ConfigActionNone,
		}
		byAction := map[domain.ConfigAction]domain.SensitiveDoguConfigDiffs{domain.ConfigActionNone: {diff1}}

		// when
		err := sut.applySensitiveDoguConfigDiffs(testCtx, byAction)

		// then
		require.NoError(t, err)
	})
}

func TestEcosystemConfigUseCase_markConfigApplied(t *testing.T) {
	t.Run("should set applied status and event", func(t *testing.T) {
		// given
		spec := &domain.BlueprintSpec{}
		expectedSpec := &domain.BlueprintSpec{}
		expectedSpec.Status = domain.StatusPhaseRegistryConfigApplied
		expectedSpec.Events = append(spec.Events, domain.RegistryConfigAppliedEvent{})
		blueprintRepoMock := newMockBlueprintSpecRepository(t)

		blueprintRepoMock.EXPECT().Update(testCtx, expectedSpec).Return(nil)

		sut := EcosystemConfigUseCase{blueprintRepository: blueprintRepoMock}

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

		sut := EcosystemConfigUseCase{blueprintRepository: blueprintRepoMock}

		// when
		err := sut.markConfigApplied(testCtx, spec)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to mark registry config applied")
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestEcosystemConfigUseCase_markApplyConfigStart(t *testing.T) {
	t.Run("should set status and event apply config", func(t *testing.T) {
		// given
		spec := &domain.BlueprintSpec{}
		expectedSpec := &domain.BlueprintSpec{}
		expectedSpec.Status = domain.StatusPhaseApplyRegistryConfig
		expectedSpec.Events = append(spec.Events, domain.ApplyRegistryConfigEvent{})
		blueprintRepoMock := newMockBlueprintSpecRepository(t)

		blueprintRepoMock.EXPECT().Update(testCtx, expectedSpec).Return(nil)

		sut := EcosystemConfigUseCase{blueprintRepository: blueprintRepoMock}

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

		sut := EcosystemConfigUseCase{blueprintRepository: blueprintRepoMock}

		// when
		err := sut.markApplyConfigStart(testCtx, spec)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "cannot mark blueprint as applying config")
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestEcosystemConfigUseCase_handleFailedApplyRegistryConfig(t *testing.T) {
	t.Run("should set applied status and event", func(t *testing.T) {
		// given
		spec := &domain.BlueprintSpec{}
		blueprintRepoMock := newMockBlueprintSpecRepository(t)

		blueprintRepoMock.EXPECT().Update(testCtx, mock.IsType(&domain.BlueprintSpec{})).Return(nil)

		sut := EcosystemConfigUseCase{blueprintRepository: blueprintRepoMock}

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

		sut := EcosystemConfigUseCase{blueprintRepository: blueprintRepoMock}

		// when
		err := sut.handleFailedApplyRegistryConfig(testCtx, spec, assert.AnError)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "cannot mark blueprint config apply as failed while handling \"applyRegistryConfigFailed\" status")
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestNewEcosystemConfigUseCase(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		doguConfigMock := newMockDoguConfigRepository(t)
		sensitiveDoguConfigMock := newMockSensitiveDoguConfigEntryRepository(t)
		globalConfigMock := newMockGlobalConfigRepository(t)

		// when
		useCase := NewEcosystemConfigUseCase(blueprintRepoMock, doguConfigMock, sensitiveDoguConfigMock, globalConfigMock)

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
			Key:      config.Key(key),
			DoguName: doguName,
		},
		Expected: domain.DoguConfigValueState{
			Value: value,
		},
		NeededAction: domain.ConfigActionSet,
	}
}

func getRemoveDoguConfigEntryDiff(key string, doguName common.SimpleDoguName) domain.DoguConfigEntryDiff {
	return domain.DoguConfigEntryDiff{
		Key: common.DoguConfigKey{
			Key:      config.Key(key),
			DoguName: doguName,
		},
		NeededAction: domain.ConfigActionRemove,
	}
}

func getSensitiveDoguConfigEntryDiffForAction(key, value string, doguName common.SimpleDoguName, action domain.ConfigAction) domain.SensitiveDoguConfigEntryDiff {
	return domain.SensitiveDoguConfigEntryDiff{
		Key: common.SensitiveDoguConfigKey{
			DoguConfigKey: common.DoguConfigKey{
				Key:      config.Key(key),
				DoguName: doguName,
			},
		},
		Expected: domain.DoguConfigValueState{
			Value: value,
		},
		NeededAction: action,
	}
}

func getRemoveSensitiveDoguConfigEntryDiff(key string, doguName common.SimpleDoguName) domain.SensitiveDoguConfigEntryDiff {
	return domain.SensitiveDoguConfigEntryDiff{
		Key: common.SensitiveDoguConfigKey{
			DoguConfigKey: common.DoguConfigKey{
				Key:      config.Key(key),
				DoguName: doguName,
			},
		},
		NeededAction: domain.ConfigActionRemove,
	}
}

func getSetGlobalConfigEntryDiff(key, value string) domain.GlobalConfigEntryDiff {
	return domain.GlobalConfigEntryDiff{
		Key: common.GlobalConfigKey(key),
		Expected: domain.GlobalConfigValueState{
			Value: value,
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
