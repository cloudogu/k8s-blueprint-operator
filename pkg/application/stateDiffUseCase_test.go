package application

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
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
	nginxStaticConfigKeyNginxKey1          = common.DoguConfigKey{DoguName: "nginx-static", Key: "nginxKey1"}
	nginxStaticConfigKeyNginxKey2          = common.DoguConfigKey{DoguName: "nginx-static", Key: "nginxKey2"}
	nginxStaticSensitiveConfigKeyNginxKey1 = nginxStaticConfigKeyNginxKey1
	nginxStaticSensitiveConfigKeyNginxKey2 = nginxStaticConfigKeyNginxKey2
)

func TestStateDiffUseCase_DetermineStateDiff(t *testing.T) {
	t.Run("should fail to load blueprint spec", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(nil, assert.AnError)

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, nil, nil, nil)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot load blueprint spec \"testBlueprint1\" to determine state diff")
	})
	t.Run("should fail to get installed dogus", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{Id: "testBlueprint1"}

		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprint, nil)

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, internalTestError)
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

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, internalTestError)
		assert.ErrorContains(t, err, "could not determine state diff")
		assert.ErrorContains(t, err, "could not collect ecosystem state")
	})
	t.Run("should fail to get installed components", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{Id: "testBlueprint1"}

		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprint, nil)

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(map[cescommons.SimpleName]*ecosystem.DoguInstallation{}, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, internalTestError)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		doguConfigRepoMock := newMockDoguConfigRepository(t)
		doguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, internalTestError)
		assert.ErrorContains(t, err, "could not determine state diff")
		assert.ErrorContains(t, err, "could not collect ecosystem state")
	})
	t.Run("should fail to get global config", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{Id: "testBlueprint1"}

		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprint, nil)

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(map[cescommons.SimpleName]*ecosystem.DoguInstallation{}, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, domainservice.NewInternalError(assert.AnError, "internal error"))

		doguConfigRepoMock := newMockDoguConfigRepository(t)
		doguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

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
		blueprint := &domain.BlueprintSpec{Id: "testBlueprint1"}

		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprint, nil)

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(map[cescommons.SimpleName]*ecosystem.DoguInstallation{}, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		doguConfigRepoMock := newMockDoguConfigRepository(t)
		doguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().
			GetAllExisting(testCtx, nilDoguNameList).
			Return(map[cescommons.SimpleName]config.DoguConfig{}, domainservice.NewInternalError(assert.AnError, "internal error"))

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		var internalError *domainservice.InternalError
		assert.ErrorAs(t, err, &internalError)
		assert.ErrorContains(t, err, "could not determine state diff")
		assert.ErrorContains(t, err, "could not collect ecosystem state")
	})
	t.Run("should fail to get sensitive dogu config", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{Id: "testBlueprint1"}

		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprint, nil)

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(map[cescommons.SimpleName]*ecosystem.DoguInstallation{}, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		doguConfigRepoMock := newMockDoguConfigRepository(t)
		doguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().
			GetAllExisting(testCtx, nilDoguNameList).
			Return(map[cescommons.SimpleName]config.DoguConfig{}, domainservice.NewInternalError(assert.AnError, "internal error"))

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		var internalError *domainservice.InternalError
		assert.ErrorAs(t, err, &internalError)
		assert.ErrorContains(t, err, "could not determine state diff")
		assert.ErrorContains(t, err, "could not collect ecosystem state")
	})
	t.Run("should fail to determine state diff for blueprint", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{Id: "testBlueprint1"}

		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprint, nil)

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)
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

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to determine state diff for blueprint \"testBlueprint1\"")
	})
	t.Run("should fail to update blueprint", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{Id: "testBlueprint1", Status: domain.StatusPhaseValidated}

		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprint, nil)
		blueprintRepoMock.EXPECT().Update(testCtx, blueprint).Return(assert.AnError)

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(map[cescommons.SimpleName]*ecosystem.DoguInstallation{}, nil)
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

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot save blueprint spec \"testBlueprint1\" after determining the state diff to the ecosystem")
	})
	t.Run("should succeed for dogu diff", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			Id: "testBlueprint1",
			EffectiveBlueprint: domain.EffectiveBlueprint{
				Dogus: []domain.Dogu{
					{
						Name:        postfixQualifiedDoguName,
						Version:     mustParseVersion(t, "2.9.0"),
						TargetState: domain.TargetStatePresent,
					},
					{
						Name:        ldapQualifiedDoguName,
						Version:     mustParseVersion(t, "1.2.3"),
						TargetState: domain.TargetStatePresent,
					},
				},
			},
			Status: domain.StatusPhaseValidated,
			// TODO: add config to test
		}

		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprint, nil)
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

		doguConfigRepoMock := newMockDoguConfigRepository(t)
		doguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.NoError(t, err)
		expectedDoguDiffs := []domain.DoguDiff{
			{
				DoguName: "postfix",
				Actual:   domain.DoguDiffState{InstallationState: domain.TargetStateAbsent},
				Expected: domain.DoguDiffState{
					Namespace:         "official",
					Version:           mustParseVersion(t, "2.9.0"),
					InstallationState: domain.TargetStatePresent,
				},
				NeededActions: []domain.Action{domain.ActionInstall},
			},
			{
				DoguName: "ldap",
				Actual: domain.DoguDiffState{
					Namespace:         "official",
					Version:           mustParseVersion(t, "1.1.1"),
					InstallationState: domain.TargetStatePresent,
				},
				Expected: domain.DoguDiffState{
					Namespace:         "official",
					Version:           mustParseVersion(t, "1.2.3"),
					InstallationState: domain.TargetStatePresent,
				},
				NeededActions: []domain.Action{domain.ActionUpgrade},
			},
		}
		assert.ElementsMatch(t, expectedDoguDiffs, blueprint.StateDiff.DoguDiffs)
	})
	t.Run("should succeed for global config diff", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			Id: "testBlueprint1",
			EffectiveBlueprint: domain.EffectiveBlueprint{
				Config: domain.Config{
					Global: domain.GlobalConfig{
						Present: map[common.GlobalConfigKey]common.GlobalConfigValue{
							"globalKey1": "globalValue",
						},
						Absent: []common.GlobalConfigKey{
							"globalKey2",
						},
					},
				},
			},
			Status: domain.StatusPhaseValidated,
		}

		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprint, nil)
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

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.NoError(t, err)

		expectedConfigDiff := []domain.GlobalConfigEntryDiff{
			{
				Key:          "globalKey1",
				Actual:       domain.GlobalConfigValueState{Value: "", Exists: false},
				Expected:     domain.GlobalConfigValueState{Value: "globalValue", Exists: true},
				NeededAction: domain.ConfigActionSet,
			},
			{
				Key:          "globalKey2",
				Actual:       domain.GlobalConfigValueState{Value: "", Exists: false},
				Expected:     domain.GlobalConfigValueState{Value: "", Exists: false},
				NeededAction: domain.ConfigActionNone,
			},
		}
		assert.ElementsMatch(t, expectedConfigDiff, blueprint.StateDiff.GlobalConfigDiffs)
	})
	t.Run("should succeed for dogu config diff", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			Id: "testBlueprint1",
			EffectiveBlueprint: domain.EffectiveBlueprint{
				Config: domain.Config{
					Dogus: map[cescommons.SimpleName]domain.CombinedDoguConfig{},
				},
			},
			Status: domain.StatusPhaseValidated,
		}

		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprint, nil)
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
		var nilDogus []cescommons.SimpleName
		doguConfigRepoMock.EXPECT().
			GetAllExisting(testCtx, nilDogus).
			Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)

		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.NoError(t, err)

		expectedConfigDiff := map[cescommons.SimpleName]domain.DoguConfigDiffs{}
		assert.Equal(t, expectedConfigDiff, blueprint.StateDiff.DoguConfigDiffs)
	})
	t.Run("should succeed for sensitive dogu config diff", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			Id: "testBlueprint1",
			EffectiveBlueprint: domain.EffectiveBlueprint{
				Config: domain.Config{
					Dogus: map[cescommons.SimpleName]domain.CombinedDoguConfig{},
				},
			},
			Status: domain.StatusPhaseValidated,
		}

		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprint, nil)
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

		sensitiveDoguConfigRepoMock.EXPECT().
			GetAllExisting(testCtx, nilDoguNameList).
			Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.NoError(t, err)

		expectedConfigDiff := map[cescommons.SimpleName]domain.DoguConfigDiffs{}
		assert.Equal(t, expectedConfigDiff, blueprint.StateDiff.SensitiveDoguConfigDiffs)
	})
}

func mustParseVersion(t *testing.T, raw string) core.Version {
	version, err := core.ParseVersion(raw)
	assert.NoError(t, err)
	return version
}

func TestStateDiffUseCase_collectEcosystemState(t *testing.T) {
	t.Run("all ok", func(t *testing.T) {
		// given
		effectiveBlueprint := domain.EffectiveBlueprint{
			Config: domain.Config{
				Global: domain.GlobalConfig{
					Present: map[common.GlobalConfigKey]common.GlobalConfigValue{
						"globalKey1": "globalValue",
					},
					Absent: []common.GlobalConfigKey{
						"globalKey2",
					},
				},
				Dogus: map[cescommons.SimpleName]domain.CombinedDoguConfig{},
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

		sensitiveDoguConfigRepoMock.EXPECT().
			GetAllExisting(testCtx, effectiveBlueprint.Config.GetDogusWithChangedSensitiveConfig()).
			Return(map[cescommons.SimpleName]config.DoguConfig{}, nil)

		sut := NewStateDiffUseCase(nil, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock)

		// when
		ecosystemState, err := sut.collectEcosystemState(testCtx, effectiveBlueprint)

		// then
		assert.NoError(t, err)

		assert.Equal(t, ecosystem.EcosystemState{
			GlobalConfig:          globalConfig,
			ConfigByDogu:          map[cescommons.SimpleName]config.DoguConfig{},
			SensitiveConfigByDogu: map[cescommons.SimpleName]config.DoguConfig{},
		}, ecosystemState)
	})
	t.Run("fail with internalError and notFoundError", func(t *testing.T) {
		// given
		effectiveBlueprint := domain.EffectiveBlueprint{
			Config: domain.Config{
				Global: domain.GlobalConfig{
					Present: map[common.GlobalConfigKey]common.GlobalConfigValue{
						"globalKey1": "globalValue",
					},
					Absent: []common.GlobalConfigKey{
						"globalKey2",
					},
				},
				Dogus: map[cescommons.SimpleName]domain.CombinedDoguConfig{},
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

		sut := NewStateDiffUseCase(nil, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock)

		// when
		ecosystemState, err := sut.collectEcosystemState(testCtx, effectiveBlueprint)

		// then
		assert.ErrorIs(t, err, internalTestError)
		assert.ErrorIs(t, err, globalConfigNotFoundError)
		assert.Equal(t, ecosystem.EcosystemState{}, ecosystemState)
	})
}
