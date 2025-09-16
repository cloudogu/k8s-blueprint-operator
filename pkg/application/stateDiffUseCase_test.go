package application

import (
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	nsOfficial = cescommons.Namespace("official")
	nsK8s      = cescommons.Namespace("k8s")

	postfix = cescommons.SimpleName("postfix")
	ldap    = cescommons.SimpleName("ldap")

	postfixQualifiedDoguName = cescommons.QualifiedName{
		Namespace:  nsOfficial,
		SimpleName: postfix,
	}
	ldapQualifiedDoguName = cescommons.QualifiedName{
		Namespace:  nsOfficial,
		SimpleName: ldap,
	}

	nilDoguNameList []cescommons.SimpleName
)

var (
	internalTestError                      = domainservice.NewInternalError(assert.AnError, "internal error")
	ldapConfigKeyNginxKey1          = common.DoguConfigKey{DoguName: "ldap", Key: "ldapKey1"}
	ldapConfigKeyNginxKey2          = common.DoguConfigKey{DoguName: "ldap", Key: "ldapKey2"}
	ldapSensitiveConfigKeyNginxKey1 = ldapConfigKeyNginxKey1
	ldapSensitiveConfigKeyNginxKey2 = ldapConfigKeyNginxKey2
	val1                                   = "val1"
	val2                                   = "val2"
	val3                                   = "val3"
)

func TestStateDiffUseCase_DetermineStateDiff(t *testing.T) {
	t.Run("should fail to get installed dogus", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{Id: "testBlueprint1"}

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, internalTestError)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		sut := NewStateDiffUseCase(nil, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, nil, nil, nil)

		// when
		err := sut.DetermineStateDiff(testCtx, blueprint)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, internalTestError)
		assert.ErrorContains(t, err, "could not determine state diff")
		assert.ErrorContains(t, err, "could not collect ecosystem state")
	})
	t.Run("should fail to get installed components", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{Id: "testBlueprint1"}

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(map[cescommons.SimpleName]*ecosystem.DoguInstallation{}, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, internalTestError)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		sut := NewStateDiffUseCase(nil, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, nil, nil, nil)

		// when
		err := sut.DetermineStateDiff(testCtx, blueprint)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, internalTestError)
		assert.ErrorContains(t, err, "could not determine state diff")
		assert.ErrorContains(t, err, "could not collect ecosystem state")
	})
	t.Run("should fail to get global config", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{Id: "testBlueprint1"}

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(map[cescommons.SimpleName]*ecosystem.DoguInstallation{}, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, domainservice.NewInternalError(assert.AnError, "internal error"))

		sut := NewStateDiffUseCase(nil, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, nil, nil, nil)

		// when
		err := sut.DetermineStateDiff(testCtx, blueprint)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		var internalError *domainservice.InternalError
		assert.ErrorAs(t, err, &internalError)
		assert.ErrorContains(t, err, "could not determine state diff")
		assert.ErrorContains(t, err, "could not collect ecosystem state")
	})
	t.Run("should fail to get dogu config", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			Id: "testBlueprint1",
			EffectiveBlueprint: domain.EffectiveBlueprint{
				Config: &domain.Config{},
			},
		}

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(map[cescommons.SimpleName]*ecosystem.DoguInstallation{}, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		doguConfigRepoMock := newMockDoguConfigRepository(t)
		doguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[cescommons.SimpleName]config.DoguConfig{}, domainservice.NewInternalError(assert.AnError, "internal error"))
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().
			GetAllExisting(testCtx, nilDoguNameList).
			Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)
		configRefReaderMock := newMockSensitiveConfigRefReader(t)
		configRefReaderMock.EXPECT().
			GetValues(testCtx, map[common.DoguConfigKey]domain.SensitiveValueRef{}).
			Return(map[common.DoguConfigKey]config.Value{}, nil)

		sut := NewStateDiffUseCase(nil, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, configRefReaderMock)

		// when
		err := sut.DetermineStateDiff(testCtx, blueprint)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		var internalError *domainservice.InternalError
		assert.ErrorAs(t, err, &internalError)
		assert.ErrorContains(t, err, "could not determine state diff")
		assert.ErrorContains(t, err, "could not collect ecosystem state")
	})
	//t.Run("should fail to get sensitive dogu config", func(t *testing.T) {
	//	// given
	//	blueprint := &domain.BlueprintSpec{Id: "testBlueprint1"}
	//
	//	doguInstallRepoMock := newMockDoguInstallationRepository(t)
	//	doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(map[cescommons.SimpleName]*ecosystem.DoguInstallation{}, nil)
	//	componentInstallRepoMock := newMockComponentInstallationRepository(t)
	//	componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)
	//
	//	globalConfigRepoMock := newMockGlobalConfigRepository(t)
	//	entries, _ := config.MapToEntries(map[string]any{})
	//	globalConfig := config.CreateGlobalConfig(entries)
	//	globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)
	//
	//	doguConfigRepoMock := newMockDoguConfigRepository(t)
	//	doguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)
	//	sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
	//	sensitiveDoguConfigRepoMock.EXPECT().
	//		GetAllExisting(testCtx, nilDoguNameList).
	//		Return(map[cescommons.SimpleName]config.DoguConfig{}, domainservice.NewInternalError(assert.AnError, "internal error"))
	//	configRefReaderMock := newMockSensitiveConfigRefReader(t)
	//	configRefReaderMock.EXPECT().
	//		GetValues(testCtx, map[common.DoguConfigKey]domain.SensitiveValueRef{}).
	//		Return(map[common.DoguConfigKey]config.Value{}, nil)
	//
	//	sut := NewStateDiffUseCase(nil, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, configRefReaderMock)
	//
	//	// when
	//	err := sut.DetermineStateDiff(testCtx, blueprint)
	//
	//	// then
	//	require.Error(t, err)
	//	assert.ErrorIs(t, err, assert.AnError)
	//	var internalError *domainservice.InternalError
	//	assert.ErrorAs(t, err, &internalError)
	//	assert.ErrorContains(t, err, "could not determine state diff")
	//	assert.ErrorContains(t, err, "could not collect ecosystem state")
	//})
	// TODO: Instead we should have a test with a forbidden diff action
	t.Run("should fail to update blueprint", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			Id:         "testBlueprint1",
			Conditions: []domain.Condition{},
		}

		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().Update(testCtx, blueprint).Return(assert.AnError)

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(map[cescommons.SimpleName]*ecosystem.DoguInstallation{}, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, nil, nil, nil)

		// when
		err := sut.DetermineStateDiff(testCtx, blueprint)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot save blueprint spec \"testBlueprint1\" after determining the state diff to the ecosystem")
	})
	t.Run("should succeed for dogu diff", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			Id:         "testBlueprint1",
			Conditions: []domain.Condition{},
			EffectiveBlueprint: domain.EffectiveBlueprint{
				Dogus: []domain.Dogu{
					{
						Name:    postfixQualifiedDoguName,
						Version: mustParseVersionToPtr(t, "2.9.0"),
						Absent:  false,
					},
					{
						Name:    ldapQualifiedDoguName,
						Version: mustParseVersionToPtr(t, "1.2.3"),
						Absent:  false,
					},
				},
			},
			// TODO: add config to test
		}

		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().Update(testCtx, blueprint).Return(nil)

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		installedDogus := map[cescommons.SimpleName]*ecosystem.DoguInstallation{
			"ldap": {Name: ldapQualifiedDoguName, Version: mustParseVersion(t, "1.1.1")},
		}
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(installedDogus, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, nil, nil, nil)

		// when
		err := sut.DetermineStateDiff(testCtx, blueprint)

		// then
		require.NoError(t, err)
		expectedDoguDiffs := []domain.DoguDiff{
			{
				DoguName: "postfix",
				Actual:   domain.DoguDiffState{Absent: true},
				Expected: domain.DoguDiffState{
					Namespace: "official",
					Version:   mustParseVersionToPtr(t, "2.9.0"),
					Absent:    false,
				},
				NeededActions: []domain.Action{domain.ActionInstall},
			},
			{
				DoguName: "ldap",
				Actual: domain.DoguDiffState{
					Namespace: "official",
					Version:   mustParseVersionToPtr(t, "1.1.1"),
					Absent:    false,
				},
				Expected: domain.DoguDiffState{
					Namespace: "official",
					Version:   mustParseVersionToPtr(t, "1.2.3"),
					Absent:    false,
				},
				NeededActions: []domain.Action{domain.ActionUpgrade},
			},
		}
		assert.ElementsMatch(t, expectedDoguDiffs, blueprint.StateDiff.DoguDiffs)
	})
	t.Run("should succeed for global config diff", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			Id:         "testBlueprint1",
			Conditions: []domain.Condition{},
			EffectiveBlueprint: domain.EffectiveBlueprint{
				Config: &domain.Config{
					Global: domain.GlobalConfigEntries{
						{
							Key:   "globalKey1",
							Value: (*config.Value)(&val1),
						},
						{
							Key:    "globalKey2",
							Absent: true,
						},
					},
				},
			},
		}

		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().Update(testCtx, blueprint).Return(nil)

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		installedDogus := map[cescommons.SimpleName]*ecosystem.DoguInstallation{}
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(installedDogus, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		doguConfigRepoMock := newMockDoguConfigRepository(t)
		doguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)
		configRefReaderMock := newMockSensitiveConfigRefReader(t)
		configRefReaderMock.EXPECT().
			GetValues(testCtx, map[common.DoguConfigKey]domain.SensitiveValueRef{}).
			Return(map[common.DoguConfigKey]config.Value{}, nil)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, configRefReaderMock)

		// when
		err := sut.DetermineStateDiff(testCtx, blueprint)

		// then
		require.NoError(t, err)

		// only changes are expected
		expectedConfigDiff := []domain.GlobalConfigEntryDiff{
			{
				Key:          "globalKey1",
				Actual:       domain.GlobalConfigValueState{Value: nil, Exists: false},
				Expected:     domain.GlobalConfigValueState{Value: &val1, Exists: true},
				NeededAction: domain.ConfigActionSet,
			},
		}
		assert.ElementsMatch(t, expectedConfigDiff, blueprint.StateDiff.GlobalConfigDiffs)
	})
	t.Run("should succeed for dogu config diff", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			Id:         "testBlueprint1",
			Conditions: []domain.Condition{},
			EffectiveBlueprint: domain.EffectiveBlueprint{
				Config: &domain.Config{
					Dogus: map[cescommons.SimpleName]domain.DoguConfigEntries{
						ldapQualifiedDoguName.SimpleName: {
							{
								Key:   ldapConfigKeyNginxKey1.Key,
								Value: (*config.Value)(&val3),
							},
							{
								Key:    ldapConfigKeyNginxKey2.Key,
								Absent: true,
							},
						},
					},
				},
			},
		}

		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().Update(testCtx, blueprint).Return(nil)

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		installedDogus := map[cescommons.SimpleName]*ecosystem.DoguInstallation{
			"ldap": {Name: ldapQualifiedDoguName, Version: mustParseVersion(t, "1.8.6")},
		}
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(installedDogus, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		doguConfigRepoMock := newMockDoguConfigRepository(t)
		doguConfigEntries, entryErr := config.MapToEntries(map[string]any{
			"ldapKey1": val1,
			"ldapKey2": val2,
		})
		require.NoError(t, entryErr)
		doguConfigRepoMock.EXPECT().
			GetAllExisting(testCtx, []cescommons.SimpleName{
				ldap,
			}).
			Return(map[cescommons.SimpleName]config.DoguConfig{
				ldap: config.CreateDoguConfig(ldap, doguConfigEntries),
			}, nil)

		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)
		configRefReaderMock := newMockSensitiveConfigRefReader(t)
		configRefReaderMock.EXPECT().
			GetValues(testCtx, map[common.DoguConfigKey]domain.SensitiveValueRef{}).
			Return(map[common.DoguConfigKey]config.Value{}, nil)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, configRefReaderMock)

		// when
		err := sut.DetermineStateDiff(testCtx, blueprint)

		// then
		require.NoError(t, err)

		expectedConfigDiff := map[cescommons.SimpleName]domain.DoguConfigDiffs{
			ldap: {
				domain.DoguConfigEntryDiff{
					Key:          ldapConfigKeyNginxKey1,
					Actual:       domain.DoguConfigValueState{Value: &val1, Exists: true},
					Expected:     domain.DoguConfigValueState{Value: &val3, Exists: true},
					NeededAction: domain.ConfigActionSet,
				},
				domain.DoguConfigEntryDiff{
					Key:          ldapConfigKeyNginxKey2,
					Actual:       domain.DoguConfigValueState{Value: &val2, Exists: true},
					Expected:     domain.DoguConfigValueState{Value: nil, Exists: false},
					NeededAction: domain.ConfigActionRemove,
				},
			},
		}
		assert.Equal(t, expectedConfigDiff, blueprint.StateDiff.DoguConfigDiffs)
	})
	t.Run("should succeed for sensitive dogu config diff", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			Id:         "testBlueprint1",
			Conditions: []domain.Condition{},
			EffectiveBlueprint: domain.EffectiveBlueprint{
				Config: &domain.Config{
					Dogus: map[cescommons.SimpleName]domain.DoguConfigEntries{
						ldap: {
							{
								Key:       ldapSensitiveConfigKeyNginxKey1.Key,
								Sensitive: true,
								SecretRef: &domain.SensitiveValueRef{
									SecretName: "ldap-conf",
									SecretKey:  "ldapKey1",
								}, // val3
							},
							{
								Key:       ldapSensitiveConfigKeyNginxKey2.Key,
								Sensitive: true,
								Absent:    true,
							},
						},
					},
				},
			},
		}

		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().Update(testCtx, blueprint).Return(nil)

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		installedDogus := map[cescommons.SimpleName]*ecosystem.DoguInstallation{
			"ldap": {Name: ldapQualifiedDoguName, Version: mustParseVersion(t, "1.8.6")},
		}
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(installedDogus, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		doguConfigRepoMock := newMockDoguConfigRepository(t)
		doguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)

		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		doguConfigEntries, entryErr := config.MapToEntries(map[string]any{
			"ldapKey1": val1,
			"ldapKey2": val2,
		})
		require.NoError(t, entryErr)
		sensitiveDoguConfigRepoMock.EXPECT().
			GetAllExisting(testCtx, []cescommons.SimpleName{
				ldap,
			}).
			Return(map[cescommons.SimpleName]config.DoguConfig{
				ldap: config.CreateDoguConfig(ldap, doguConfigEntries),
			}, nil)
		configRefReaderMock := newMockSensitiveConfigRefReader(t)
		configRefReaderMock.EXPECT().
			GetValues(
				testCtx,
				blueprint.EffectiveBlueprint.Config.GetSensitiveConfigReferences(),
			).
			Return(map[common.DoguConfigKey]config.Value{
				ldapConfigKeyNginxKey1: config.Value(val3),
			}, nil)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, configRefReaderMock)

		// when
		err := sut.DetermineStateDiff(testCtx, blueprint)

		// then
		require.NoError(t, err)

		expectedConfigDiff := map[cescommons.SimpleName]domain.DoguConfigDiffs{
			ldap: {
				{
					Key:          ldapSensitiveConfigKeyNginxKey1,
					Actual:       domain.DoguConfigValueState{Value: &val1, Exists: true},
					Expected:     domain.DoguConfigValueState{Value: &val3, Exists: true},
					NeededAction: domain.ConfigActionSet,
				},
				{
					Key:          ldapSensitiveConfigKeyNginxKey2,
					Actual:       domain.DoguConfigValueState{Value: &val2, Exists: true},
					Expected:     domain.DoguConfigValueState{Value: nil, Exists: false},
					NeededAction: domain.ConfigActionRemove,
				},
			},
		}
		assert.Equal(t, expectedConfigDiff, blueprint.StateDiff.SensitiveDoguConfigDiffs)
	})
}

func mustParseVersion(t *testing.T, raw string) core.Version {
	version, err := core.ParseVersion(raw)
	assert.NoError(t, err)
	return version
}

func mustParseVersionToPtr(t *testing.T, raw string) *core.Version {
	version := mustParseVersion(t, raw)
	return &version
}

func TestStateDiffUseCase_collectEcosystemState(t *testing.T) {
	t.Run("all ok", func(t *testing.T) {
		// given
		effectiveBlueprint := domain.EffectiveBlueprint{
			Config: &domain.Config{
				Global: domain.GlobalConfigEntries{
					{
						Key:   "globalKey1",
						Value: (*config.Value)(&val1),
					},
					{
						Key:    "globalKey2",
						Absent: true,
					},
				},
				Dogus: map[cescommons.SimpleName]domain.DoguConfigEntries{
					ldapQualifiedDoguName.SimpleName: {
						{
							Key:   ldapConfigKeyNginxKey1.Key,
							Value: (*config.Value)(&val1),
						},
						{
							Key:    ldapConfigKeyNginxKey2.Key,
							Absent: true,
						},
						{
							Key:       ldapSensitiveConfigKeyNginxKey1.Key,
							Sensitive: true,
							SecretRef: &domain.SensitiveValueRef{
								SecretName: "nginx-conf",
								SecretKey:  "nginxKey1",
							},
						},
						{
							Key:       ldapSensitiveConfigKeyNginxKey2.Key,
							Sensitive: true,
							Absent:    true,
						},
					},
				},
			},
		}

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		doguConfigRepoMock := newMockDoguConfigRepository(t)
		doguConfigRepoMock.EXPECT().
			GetAllExisting(testCtx, effectiveBlueprint.Config.GetDogusWithChangedConfig()).
			Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)

		ldapConfig := config.CreateDoguConfig(ldap, map[config.Key]config.Value{
			"ldapKey1": "val1",
		})
		sensitiveDoguConfigRepoMock.EXPECT().
			GetAllExisting(testCtx, effectiveBlueprint.Config.GetDogusWithChangedSensitiveConfig()).
			Return(map[cescommons.SimpleName]config.DoguConfig{
				ldap: ldapConfig,
			}, nil)
		configRefReaderMock := newMockSensitiveConfigRefReader(t)

		sut := NewStateDiffUseCase(nil, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, configRefReaderMock)

		// when
		ecosystemState, err := sut.collectEcosystemState(testCtx, effectiveBlueprint)

		// then
		assert.NoError(t, err)

		assert.Equal(t, ecosystem.EcosystemState{
			GlobalConfig: globalConfig,
			ConfigByDogu: map[cescommons.SimpleName]config.DoguConfig{},
			SensitiveConfigByDogu: map[cescommons.SimpleName]config.DoguConfig{
				ldap: ldapConfig,
			},
		}, ecosystemState)
	})
	t.Run("fail with internalError and notFoundError", func(t *testing.T) {
		// given
		effectiveBlueprint := domain.EffectiveBlueprint{
			Config: &domain.Config{
				Global: domain.GlobalConfigEntries{
					{
						Key:   "globalKey1",
						Value: (*config.Value)(&val1),
					},
					{
						Key:    "globalKey2",
						Absent: true,
					},
				},
				Dogus: map[cescommons.SimpleName]domain.DoguConfigEntries{
					ldapQualifiedDoguName.SimpleName: {
						{
							Key:   ldapConfigKeyNginxKey1.Key,
							Value: (*config.Value)(&val1),
						},
						{
							Key:    ldapConfigKeyNginxKey2.Key,
							Absent: true,
						},
						{
							Key:       ldapSensitiveConfigKeyNginxKey1.Key,
							Sensitive: true,
							SecretRef: &domain.SensitiveValueRef{
								SecretName: "nginx-conf",
								SecretKey:  "nginxKey1",
							},
						},
						{
							Key:       ldapSensitiveConfigKeyNginxKey2.Key,
							Sensitive: true,
							Absent:    true,
						},
					},
				},
			},
		}

		globalConfigNotFoundError := domainservice.NewNotFoundError(assert.AnError, "global config not found")

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, globalConfigNotFoundError)

		doguConfigRepoMock := newMockDoguConfigRepository(t)
		doguConfigRepoMock.EXPECT().
			GetAllExisting(testCtx, effectiveBlueprint.Config.GetDogusWithChangedConfig()).
			Return(map[cescommons.SimpleName]config.DoguConfig{}, internalTestError)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().
			GetAllExisting(testCtx, effectiveBlueprint.Config.GetDogusWithChangedSensitiveConfig()).
			Return(map[cescommons.SimpleName]config.DoguConfig{}, internalTestError)
		configRefReaderMock := newMockSensitiveConfigRefReader(t)

		sut := NewStateDiffUseCase(nil, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, configRefReaderMock)

		// when
		ecosystemState, err := sut.collectEcosystemState(testCtx, effectiveBlueprint)

		// then
		assert.ErrorIs(t, err, internalTestError)
		assert.ErrorIs(t, err, globalConfigNotFoundError)
		assert.Equal(t, ecosystem.EcosystemState{}, ecosystemState)
	})
}
