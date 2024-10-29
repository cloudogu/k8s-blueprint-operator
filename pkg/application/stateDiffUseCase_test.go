package application

import (
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
	nsOfficial = common.DoguNamespace("official")
	nsK8s      = common.DoguNamespace("k8s")

	postfix      = common.SimpleDoguName("postfix")
	ldap         = common.SimpleDoguName("ldap")
	nginxIngress = common.SimpleDoguName("nginx-ingress")
	nginxStatic  = common.SimpleDoguName("nginx-static")

	postfixQualifiedDoguName = common.QualifiedDoguName{
		Namespace:  nsOfficial,
		SimpleName: postfix,
	}
	ldapQualifiedDoguName = common.QualifiedDoguName{
		Namespace:  nsOfficial,
		SimpleName: ldap,
	}
	nginxIngressQualifiedDoguName = common.QualifiedDoguName{
		Namespace:  nsK8s,
		SimpleName: nginxIngress,
	}
	nginxStaticQualifiedDoguName = common.QualifiedDoguName{
		Namespace:  nsK8s,
		SimpleName: nginxStatic,
	}

	nilDoguNameList []common.SimpleDoguName
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
		doguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[config.SimpleDoguName]config.DoguConfig{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[config.SimpleDoguName]config.DoguConfig{}, nil)

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
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(map[common.SimpleDoguName]*ecosystem.DoguInstallation{}, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, internalTestError)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		doguConfigRepoMock := newMockDoguConfigRepository(t)
		doguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[config.SimpleDoguName]config.DoguConfig{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[config.SimpleDoguName]config.DoguConfig{}, nil)

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
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(map[common.SimpleDoguName]*ecosystem.DoguInstallation{}, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, domainservice.NewInternalError(assert.AnError, "internal error"))

		doguConfigRepoMock := newMockDoguConfigRepository(t)
		doguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[config.SimpleDoguName]config.DoguConfig{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[config.SimpleDoguName]config.DoguConfig{}, nil)

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
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(map[common.SimpleDoguName]*ecosystem.DoguInstallation{}, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		doguConfigRepoMock := newMockDoguConfigRepository(t)
		doguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[config.SimpleDoguName]config.DoguConfig{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().
			GetAllExisting(testCtx, nilDoguNameList).
			Return(map[config.SimpleDoguName]config.DoguConfig{}, domainservice.NewInternalError(assert.AnError, "internal error"))

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
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(map[common.SimpleDoguName]*ecosystem.DoguInstallation{}, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		doguConfigRepoMock := newMockDoguConfigRepository(t)
		doguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[config.SimpleDoguName]config.DoguConfig{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().
			GetAllExisting(testCtx, nilDoguNameList).
			Return(map[config.SimpleDoguName]config.DoguConfig{}, domainservice.NewInternalError(assert.AnError, "internal error"))

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
		doguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[config.SimpleDoguName]config.DoguConfig{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[config.SimpleDoguName]config.DoguConfig{}, nil)

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
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(map[common.SimpleDoguName]*ecosystem.DoguInstallation{}, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		doguConfigRepoMock := newMockDoguConfigRepository(t)
		doguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[config.SimpleDoguName]config.DoguConfig{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[config.SimpleDoguName]config.DoguConfig{}, nil)

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
					{
						Name:        nginxIngressQualifiedDoguName,
						Version:     mustParseVersion(t, "1.8.5"),
						TargetState: domain.TargetStatePresent,
					},
					{
						Name:        nginxStaticQualifiedDoguName,
						TargetState: domain.TargetStateAbsent,
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
		installedDogus := map[common.SimpleDoguName]*ecosystem.DoguInstallation{
			"ldap":          {Name: ldapQualifiedDoguName, Version: mustParseVersion(t, "1.1.1")},
			"nginx-ingress": {Name: nginxIngressQualifiedDoguName, Version: mustParseVersion(t, "1.8.5")},
			"nginx-static":  {Name: nginxStaticQualifiedDoguName, Version: mustParseVersion(t, "1.8.6")},
		}
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(installedDogus, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		doguConfigRepoMock := newMockDoguConfigRepository(t)
		doguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[config.SimpleDoguName]config.DoguConfig{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[config.SimpleDoguName]config.DoguConfig{}, nil)

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
			{
				DoguName: "nginx-ingress",
				Actual: domain.DoguDiffState{
					Namespace:         "k8s",
					Version:           mustParseVersion(t, "1.8.5"),
					InstallationState: domain.TargetStatePresent,
				},
				Expected: domain.DoguDiffState{
					Namespace:         "k8s",
					Version:           mustParseVersion(t, "1.8.5"),
					InstallationState: domain.TargetStatePresent,
				},
				NeededActions: nil,
			},
			{
				DoguName: "nginx-static",
				Actual: domain.DoguDiffState{
					Namespace:         "k8s",
					Version:           mustParseVersion(t, "1.8.6"),
					InstallationState: domain.TargetStatePresent,
				},
				Expected: domain.DoguDiffState{
					Namespace:         "k8s",
					InstallationState: domain.TargetStateAbsent,
				},
				NeededActions: []domain.Action{domain.ActionUninstall},
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
		installedDogus := map[common.SimpleDoguName]*ecosystem.DoguInstallation{
			"nginx-static": {Name: nginxStaticQualifiedDoguName, Version: mustParseVersion(t, "1.8.6")},
		}
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(installedDogus, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		doguConfigRepoMock := newMockDoguConfigRepository(t)
		doguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[config.SimpleDoguName]config.DoguConfig{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[config.SimpleDoguName]config.DoguConfig{}, nil)

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
					Dogus: map[common.SimpleDoguName]domain.CombinedDoguConfig{
						nginxStaticQualifiedDoguName.SimpleName: {
							DoguName: nginxStaticQualifiedDoguName.SimpleName,
							Config: domain.DoguConfig{
								Present: map[common.DoguConfigKey]common.DoguConfigValue{
									nginxStaticConfigKeyNginxKey1: "nginxVal1",
								},
								Absent: []common.DoguConfigKey{
									nginxStaticConfigKeyNginxKey2,
								},
							},
							SensitiveConfig: domain.SensitiveDoguConfig{},
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
		installedDogus := map[common.SimpleDoguName]*ecosystem.DoguInstallation{
			"nginx-static": {Name: nginxStaticQualifiedDoguName, Version: mustParseVersion(t, "1.8.6")},
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
			"nginxKey1": "val1",
			"nginxKey2": "val2",
		})
		require.NoError(t, entryErr)
		doguConfigRepoMock.EXPECT().
			GetAllExisting(testCtx, []common.SimpleDoguName{
				nginxStatic,
			}).
			Return(map[config.SimpleDoguName]config.DoguConfig{
				nginxStatic: config.CreateDoguConfig(nginxStatic, doguConfigEntries),
			}, nil)

		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[config.SimpleDoguName]config.DoguConfig{}, nil)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.NoError(t, err)

		expectedConfigDiff := map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{
			"nginx-static": {
				DoguConfigDiff: []domain.DoguConfigEntryDiff{
					{
						Key:          nginxStaticConfigKeyNginxKey1,
						Actual:       domain.DoguConfigValueState{Value: "val1", Exists: true},
						Expected:     domain.DoguConfigValueState{Value: "nginxVal1", Exists: true},
						NeededAction: domain.ConfigActionSet,
					},
					{
						Key:          nginxStaticConfigKeyNginxKey2,
						Actual:       domain.DoguConfigValueState{Value: "val2", Exists: true},
						Expected:     domain.DoguConfigValueState{Value: "", Exists: false},
						NeededAction: domain.ConfigActionRemove,
					},
				},
				SensitiveDoguConfigDiff: nil,
			},
		}
		assert.Equal(t, expectedConfigDiff, blueprint.StateDiff.DoguConfigDiffs)
	})
	t.Run("should succeed for sensitive dogu config diff", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			Id: "testBlueprint1",
			EffectiveBlueprint: domain.EffectiveBlueprint{
				Config: domain.Config{
					Dogus: map[common.SimpleDoguName]domain.CombinedDoguConfig{
						nginxStaticQualifiedDoguName.SimpleName: {
							DoguName: nginxStaticQualifiedDoguName.SimpleName,
							SensitiveConfig: domain.SensitiveDoguConfig{
								Present: map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{
									nginxStaticSensitiveConfigKeyNginxKey1: "nginxVal1",
								},
								Absent: []common.SensitiveDoguConfigKey{
									nginxStaticSensitiveConfigKeyNginxKey2,
								},
							},
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
		installedDogus := map[common.SimpleDoguName]*ecosystem.DoguInstallation{
			"nginx-static": {Name: nginxStaticQualifiedDoguName, Version: mustParseVersion(t, "1.8.6")},
		}
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(installedDogus, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)

		globalConfigRepoMock := newMockGlobalConfigRepository(t)
		entries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(entries)
		globalConfigRepoMock.EXPECT().Get(testCtx).Return(globalConfig, nil)

		doguConfigRepoMock := newMockDoguConfigRepository(t)
		doguConfigRepoMock.EXPECT().GetAllExisting(testCtx, nilDoguNameList).Return(map[config.SimpleDoguName]config.DoguConfig{}, nil)

		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		doguConfigEntries, entryErr := config.MapToEntries(map[string]any{
			"nginxKey1": "val1",
			"nginxKey2": "val2",
		})
		require.NoError(t, entryErr)
		sensitiveDoguConfigRepoMock.EXPECT().
			GetAllExisting(testCtx, []common.SimpleDoguName{
				nginxStatic,
			}).
			Return(map[config.SimpleDoguName]config.DoguConfig{
				nginxStatic: config.CreateDoguConfig(nginxStatic, doguConfigEntries),
			}, nil)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.NoError(t, err)

		expectedConfigDiff := map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{
			"nginx-static": {
				SensitiveDoguConfigDiff: []domain.SensitiveDoguConfigEntryDiff{
					{
						Key:          nginxStaticSensitiveConfigKeyNginxKey1,
						Actual:       domain.DoguConfigValueState{Value: "val1", Exists: true},
						Expected:     domain.DoguConfigValueState{Value: "nginxVal1", Exists: true},
						NeededAction: domain.ConfigActionSet,
					},
					{
						Key:          nginxStaticSensitiveConfigKeyNginxKey2,
						Actual:       domain.DoguConfigValueState{Value: "val2", Exists: true},
						Expected:     domain.DoguConfigValueState{Value: "", Exists: false},
						NeededAction: domain.ConfigActionRemove,
					},
				},
			},
		}
		assert.Equal(t, expectedConfigDiff, blueprint.StateDiff.DoguConfigDiffs)
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
				Dogus: map[common.SimpleDoguName]domain.CombinedDoguConfig{
					nginxStaticQualifiedDoguName.SimpleName: {
						DoguName: nginxStaticQualifiedDoguName.SimpleName,
						Config: domain.DoguConfig{
							Present: map[common.DoguConfigKey]common.DoguConfigValue{
								nginxStaticConfigKeyNginxKey1: "nginxVal1",
							},
							Absent: []common.DoguConfigKey{
								nginxStaticConfigKeyNginxKey2,
							},
						},
						SensitiveConfig: domain.SensitiveDoguConfig{
							Present: map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{
								nginxStaticSensitiveConfigKeyNginxKey1: "nginxVal1",
							},
							Absent: []common.SensitiveDoguConfigKey{
								nginxStaticSensitiveConfigKeyNginxKey2,
							},
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
			Return(map[config.SimpleDoguName]config.DoguConfig{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)

		nginxStaticConfig := config.CreateDoguConfig(nginxStatic, map[config.Key]config.Value{
			"nginxKey1": "val1",
		})
		sensitiveDoguConfigRepoMock.EXPECT().
			GetAllExisting(testCtx, effectiveBlueprint.Config.GetDogusWithChangedSensitiveConfig()).
			Return(map[config.SimpleDoguName]config.DoguConfig{
				nginxStatic: nginxStaticConfig,
			}, nil)

		sut := NewStateDiffUseCase(nil, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock)

		// when
		ecosystemState, err := sut.collectEcosystemState(testCtx, effectiveBlueprint)

		// then
		assert.NoError(t, err)

		assert.Equal(t, ecosystem.EcosystemState{
			GlobalConfig: globalConfig,
			ConfigByDogu: map[common.SimpleDoguName]config.DoguConfig{},
			SensitiveConfigByDogu: map[config.SimpleDoguName]config.DoguConfig{
				nginxStatic: nginxStaticConfig,
			},
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
				Dogus: map[common.SimpleDoguName]domain.CombinedDoguConfig{
					nginxStaticQualifiedDoguName.SimpleName: {
						DoguName: nginxStaticQualifiedDoguName.SimpleName,
						Config: domain.DoguConfig{
							Present: map[common.DoguConfigKey]common.DoguConfigValue{
								nginxStaticConfigKeyNginxKey1: "nginxVal1",
							},
							Absent: []common.DoguConfigKey{
								nginxStaticConfigKeyNginxKey2,
							},
						},
						SensitiveConfig: domain.SensitiveDoguConfig{
							Present: map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{
								nginxStaticSensitiveConfigKeyNginxKey1: "nginxVal1",
							},
							Absent: []common.SensitiveDoguConfigKey{
								nginxStaticSensitiveConfigKeyNginxKey2,
							},
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
			Return(map[config.SimpleDoguName]config.DoguConfig{}, internalTestError)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().
			GetAllExisting(testCtx, effectiveBlueprint.Config.GetDogusWithChangedSensitiveConfig()).
			Return(map[config.SimpleDoguName]config.DoguConfig{}, internalTestError)

		sut := NewStateDiffUseCase(nil, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock)

		// when
		ecosystemState, err := sut.collectEcosystemState(testCtx, effectiveBlueprint)

		// then
		assert.ErrorIs(t, err, internalTestError)
		assert.ErrorIs(t, err, globalConfigNotFoundError)
		assert.Equal(t, ecosystem.EcosystemState{}, ecosystemState)
	})
}
