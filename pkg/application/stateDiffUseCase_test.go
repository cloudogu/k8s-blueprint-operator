package application

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var postfixQualifiedDoguName = common.QualifiedDoguName{
	Namespace:  "official",
	SimpleName: "postfix",
}

var ldapQualifiedDoguName = common.QualifiedDoguName{
	Namespace:  "official",
	SimpleName: "ldap",
}

var nginxIngressQualifiedDoguName = common.QualifiedDoguName{
	Namespace:  "k8s",
	SimpleName: "nginx-ingress",
}

var nginxStaticQualifiedDoguName = common.QualifiedDoguName{
	Namespace:  "k8s",
	SimpleName: "nginx-static",
}

var (
	internalTestError                      = domainservice.NewInternalError(assert.AnError, "internal error")
	nginxStaticConfigKeyNginxKey1          = common.DoguConfigKey{DoguName: "nginx-static", Key: "nginxKey1"}
	nginxStaticConfigKeyNginxKey2          = common.DoguConfigKey{DoguName: "nginx-static", Key: "nginxKey2"}
	nginxStaticSensitiveConfigKeyNginxKey1 = common.SensitiveDoguConfigKey{DoguConfigKey: nginxStaticConfigKeyNginxKey1}
	nginxStaticSensitiveConfigKeyNginxKey2 = common.SensitiveDoguConfigKey{DoguConfigKey: nginxStaticConfigKeyNginxKey2}
)

func TestStateDiffUseCase_DetermineStateDiff(t *testing.T) {
	t.Run("should fail to load blueprint spec", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(nil, assert.AnError)

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, nil, nil, nil, nil)

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
		globalConfigRepoMock := newMockGlobalConfigEntryRepository(t)
		globalConfigRepoMock.EXPECT().GetAllByKey(testCtx, blueprint.EffectiveBlueprint.Config.Global.GetGlobalConfigKeys()).Return(map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry{}, nil)
		doguConfigRepoMock := newMockDoguConfigEntryRepository(t)
		doguConfigRepoMock.EXPECT().GetAllByKey(testCtx, []common.DoguConfigKey(nil)).Return(map[common.DoguConfigKey]*ecosystem.DoguConfigEntry{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigEntryRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllByKey(testCtx, []common.SensitiveDoguConfigKey(nil)).Return(map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry{}, nil)
		encryptionAdapterMock := newMockConfigEncryptionAdapter(t)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, encryptionAdapterMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, internalTestError)
		assert.ErrorContains(t, err, "could not determine state diff")
		assert.ErrorContains(t, err, "could not collect cluster state")
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
		globalConfigRepoMock := newMockGlobalConfigEntryRepository(t)
		globalConfigRepoMock.EXPECT().GetAllByKey(testCtx, blueprint.EffectiveBlueprint.Config.Global.GetGlobalConfigKeys()).Return(map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry{}, nil)
		doguConfigRepoMock := newMockDoguConfigEntryRepository(t)
		doguConfigRepoMock.EXPECT().GetAllByKey(testCtx, []common.DoguConfigKey(nil)).Return(map[common.DoguConfigKey]*ecosystem.DoguConfigEntry{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigEntryRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllByKey(testCtx, []common.SensitiveDoguConfigKey(nil)).Return(map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry{}, nil)
		encryptionAdapterMock := newMockConfigEncryptionAdapter(t)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, encryptionAdapterMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, internalTestError)
		assert.ErrorContains(t, err, "could not determine state diff")
		assert.ErrorContains(t, err, "could not collect cluster state")
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
		globalConfigRepoMock := newMockGlobalConfigEntryRepository(t)
		globalConfigRepoMock.EXPECT().
			GetAllByKey(testCtx, blueprint.EffectiveBlueprint.Config.Global.GetGlobalConfigKeys()).
			Return(map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry{}, domainservice.NewInternalError(assert.AnError, "internal error"))
		doguConfigRepoMock := newMockDoguConfigEntryRepository(t)
		doguConfigRepoMock.EXPECT().GetAllByKey(testCtx, []common.DoguConfigKey(nil)).Return(map[common.DoguConfigKey]*ecosystem.DoguConfigEntry{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigEntryRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllByKey(testCtx, []common.SensitiveDoguConfigKey(nil)).Return(map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry{}, nil)
		encryptionAdapterMock := newMockConfigEncryptionAdapter(t)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, encryptionAdapterMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		var internalError *domainservice.InternalError
		assert.ErrorAs(t, err, &internalError)
		assert.ErrorContains(t, err, "could not determine state diff")
		assert.ErrorContains(t, err, "could not collect cluster state")
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
		globalConfigRepoMock := newMockGlobalConfigEntryRepository(t)
		globalConfigRepoMock.EXPECT().GetAllByKey(testCtx, blueprint.EffectiveBlueprint.Config.Global.GetGlobalConfigKeys()).Return(map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry{}, nil)
		doguConfigRepoMock := newMockDoguConfigEntryRepository(t)
		doguConfigRepoMock.EXPECT().
			GetAllByKey(testCtx, []common.DoguConfigKey(nil)).
			Return(map[common.DoguConfigKey]*ecosystem.DoguConfigEntry{}, domainservice.NewInternalError(assert.AnError, "internal error"))
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigEntryRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllByKey(testCtx, []common.SensitiveDoguConfigKey(nil)).Return(map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry{}, nil)
		encryptionAdapterMock := newMockConfigEncryptionAdapter(t)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, encryptionAdapterMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		var internalError *domainservice.InternalError
		assert.ErrorAs(t, err, &internalError)
		assert.ErrorContains(t, err, "could not determine state diff")
		assert.ErrorContains(t, err, "could not collect cluster state")
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
		globalConfigRepoMock := newMockGlobalConfigEntryRepository(t)
		globalConfigRepoMock.EXPECT().GetAllByKey(testCtx, blueprint.EffectiveBlueprint.Config.Global.GetGlobalConfigKeys()).Return(map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry{}, nil)
		doguConfigRepoMock := newMockDoguConfigEntryRepository(t)
		doguConfigRepoMock.EXPECT().GetAllByKey(testCtx, []common.DoguConfigKey(nil)).Return(map[common.DoguConfigKey]*ecosystem.DoguConfigEntry{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigEntryRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().
			GetAllByKey(testCtx, []common.SensitiveDoguConfigKey(nil)).
			Return(map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry{}, domainservice.NewInternalError(assert.AnError, "internal error"))
		encryptionAdapterMock := newMockConfigEncryptionAdapter(t)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, encryptionAdapterMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		var internalError *domainservice.InternalError
		assert.ErrorAs(t, err, &internalError)
		assert.ErrorContains(t, err, "could not determine state diff")
		assert.ErrorContains(t, err, "could not collect cluster state")
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
		globalConfigRepoMock := newMockGlobalConfigEntryRepository(t)
		globalConfigRepoMock.EXPECT().GetAllByKey(testCtx, blueprint.EffectiveBlueprint.Config.Global.GetGlobalConfigKeys()).Return(map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry{}, nil)
		doguConfigRepoMock := newMockDoguConfigEntryRepository(t)
		doguConfigRepoMock.EXPECT().GetAllByKey(testCtx, []common.DoguConfigKey(nil)).Return(map[common.DoguConfigKey]*ecosystem.DoguConfigEntry{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigEntryRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllByKey(testCtx, []common.SensitiveDoguConfigKey(nil)).Return(map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry{}, nil)
		encryptionAdapterMock := newMockConfigEncryptionAdapter(t)
		encryptionAdapterMock.EXPECT().DecryptAll(testCtx, map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue{}).Return(nil, nil)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, encryptionAdapterMock)

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
		globalConfigRepoMock := newMockGlobalConfigEntryRepository(t)
		globalConfigRepoMock.EXPECT().GetAllByKey(testCtx, blueprint.EffectiveBlueprint.Config.Global.GetGlobalConfigKeys()).Return(map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry{}, nil)
		doguConfigRepoMock := newMockDoguConfigEntryRepository(t)
		doguConfigRepoMock.EXPECT().GetAllByKey(testCtx, []common.DoguConfigKey(nil)).Return(map[common.DoguConfigKey]*ecosystem.DoguConfigEntry{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigEntryRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllByKey(testCtx, []common.SensitiveDoguConfigKey(nil)).Return(map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry{}, nil)
		encryptionAdapterMock := newMockConfigEncryptionAdapter(t)
		encryptionAdapterMock.EXPECT().DecryptAll(testCtx, map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue{}).Return(nil, nil)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, encryptionAdapterMock)

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
		globalConfigRepoMock := newMockGlobalConfigEntryRepository(t)
		globalConfigRepoMock.EXPECT().GetAllByKey(testCtx, []common.GlobalConfigKey(nil)).Return(map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry{}, nil)
		doguConfigRepoMock := newMockDoguConfigEntryRepository(t)
		doguConfigRepoMock.EXPECT().GetAllByKey(testCtx, []common.DoguConfigKey(nil)).Return(map[common.DoguConfigKey]*ecosystem.DoguConfigEntry{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigEntryRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllByKey(testCtx, []common.SensitiveDoguConfigKey(nil)).Return(map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry{}, nil)
		encryptionAdapterMock := newMockConfigEncryptionAdapter(t)
		encryptionAdapterMock.EXPECT().DecryptAll(testCtx, map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue{}).Return(nil, nil)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, encryptionAdapterMock)

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
				NeededAction: domain.ActionInstall,
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
				NeededAction: domain.ActionUpgrade,
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
				NeededAction: domain.ActionNone,
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
				NeededAction: domain.ActionUninstall,
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
		globalConfigRepoMock := newMockGlobalConfigEntryRepository(t)
		globalConfigRepoMock.EXPECT().
			GetAllByKey(testCtx, blueprint.EffectiveBlueprint.Config.Global.GetGlobalConfigKeys()).
			Return(map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry{}, nil)
		doguConfigRepoMock := newMockDoguConfigEntryRepository(t)
		doguConfigRepoMock.EXPECT().GetAllByKey(testCtx, []common.DoguConfigKey(nil)).Return(map[common.DoguConfigKey]*ecosystem.DoguConfigEntry{}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigEntryRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().GetAllByKey(testCtx, []common.SensitiveDoguConfigKey(nil)).Return(map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry{}, nil)
		encryptionAdapterMock := newMockConfigEncryptionAdapter(t)
		encryptionAdapterMock.EXPECT().DecryptAll(testCtx, map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue{}).Return(nil, nil)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, encryptionAdapterMock)

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
		globalConfigRepoMock := newMockGlobalConfigEntryRepository(t)
		globalConfigRepoMock.EXPECT().GetAllByKey(testCtx, blueprint.EffectiveBlueprint.Config.Global.GetGlobalConfigKeys()).Return(map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry{}, nil)
		doguConfigRepoMock := newMockDoguConfigEntryRepository(t)
		doguConfigRepoMock.EXPECT().
			GetAllByKey(testCtx, []common.DoguConfigKey{
				nginxStaticConfigKeyNginxKey1,
				nginxStaticConfigKeyNginxKey2,
			}).
			Return(map[common.DoguConfigKey]*ecosystem.DoguConfigEntry{
				nginxStaticConfigKeyNginxKey1: {
					Key:   nginxStaticConfigKeyNginxKey1,
					Value: "val1",
				},
				nginxStaticConfigKeyNginxKey2: {
					Key:   nginxStaticConfigKeyNginxKey2,
					Value: "val2",
				},
			}, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigEntryRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().
			GetAllByKey(testCtx, []common.SensitiveDoguConfigKey(nil)).
			Return(map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry{}, nil)
		encryptionAdapterMock := newMockConfigEncryptionAdapter(t)
		encryptionAdapterMock.EXPECT().DecryptAll(testCtx, map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue{}).Return(nil, nil)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, encryptionAdapterMock)

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
		globalConfigRepoMock := newMockGlobalConfigEntryRepository(t)
		globalConfigRepoMock.EXPECT().GetAllByKey(testCtx, blueprint.EffectiveBlueprint.Config.Global.GetGlobalConfigKeys()).Return(map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry{}, nil)
		doguConfigRepoMock := newMockDoguConfigEntryRepository(t)
		doguConfigRepoMock.EXPECT().GetAllByKey(testCtx, []common.DoguConfigKey(nil)).Return(nil, nil)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigEntryRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().
			GetAllByKey(testCtx, []common.SensitiveDoguConfigKey{
				nginxStaticSensitiveConfigKeyNginxKey1,
				nginxStaticSensitiveConfigKeyNginxKey2,
			}).
			Return(map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry{
				nginxStaticSensitiveConfigKeyNginxKey1: {
					Key:   nginxStaticSensitiveConfigKeyNginxKey1,
					Value: "encrypted",
				},
				nginxStaticSensitiveConfigKeyNginxKey2: {
					Key:   nginxStaticSensitiveConfigKeyNginxKey2,
					Value: "encrypted",
				},
			}, nil)
		encryptionAdapterMock := newMockConfigEncryptionAdapter(t)
		encryptionAdapterMock.EXPECT().
			DecryptAll(testCtx, map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue{
				nginxStaticSensitiveConfigKeyNginxKey1: "encrypted",
				nginxStaticSensitiveConfigKeyNginxKey2: "encrypted",
			}).
			Return(map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{
				nginxStaticSensitiveConfigKeyNginxKey1: "val1",
				nginxStaticSensitiveConfigKeyNginxKey2: "val2",
			}, nil)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, encryptionAdapterMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.NoError(t, err)

		expectedConfigDiff := map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{
			"nginx-static": {
				SensitiveDoguConfigDiff: []domain.SensitiveDoguConfigEntryDiff{
					{
						Key:                  nginxStaticSensitiveConfigKeyNginxKey1,
						Actual:               domain.DoguConfigValueState{Value: "val1", Exists: true},
						Expected:             domain.DoguConfigValueState{Value: "nginxVal1", Exists: true},
						NeededAction:         domain.ConfigActionSet,
						DoguAlreadyInstalled: true,
					},
					{
						Key:                  nginxStaticSensitiveConfigKeyNginxKey2,
						Actual:               domain.DoguConfigValueState{Value: "val2", Exists: true},
						Expected:             domain.DoguConfigValueState{Value: "", Exists: false},
						NeededAction:         domain.ConfigActionRemove,
						DoguAlreadyInstalled: true,
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

func TestStateDiffUseCase_collectClusterState(t *testing.T) {
	t.Run("ignore not found errors", func(t *testing.T) {
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
		doguConfigNotFoundError := domainservice.NewNotFoundError(assert.AnError, "dogu config not found")
		sensitiveConfigNotFoundError := domainservice.NewNotFoundError(assert.AnError, "sensitive config not found")

		encryptedEntry := &ecosystem.SensitiveDoguConfigEntry{
			Key:   nginxStaticSensitiveConfigKeyNginxKey1,
			Value: "encrypted",
		}

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)
		globalConfigRepoMock := newMockGlobalConfigEntryRepository(t)
		globalConfigRepoMock.EXPECT().
			GetAllByKey(testCtx, effectiveBlueprint.Config.Global.GetGlobalConfigKeys()).
			Return(
				map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry{},
				globalConfigNotFoundError,
			)
		doguConfigRepoMock := newMockDoguConfigEntryRepository(t)
		doguConfigRepoMock.EXPECT().
			GetAllByKey(testCtx, effectiveBlueprint.Config.GetDoguConfigKeys()).
			Return(
				map[common.DoguConfigKey]*ecosystem.DoguConfigEntry{},
				doguConfigNotFoundError,
			)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigEntryRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().
			GetAllByKey(testCtx, effectiveBlueprint.Config.GetSensitiveDoguConfigKeys()).
			Return(
				map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry{
					nginxStaticSensitiveConfigKeyNginxKey1: encryptedEntry,
				},
				sensitiveConfigNotFoundError,
			)
		encryptionAdapterMock := newMockConfigEncryptionAdapter(t)
		encryptionAdapterMock.EXPECT().
			DecryptAll(testCtx, map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue{
				nginxStaticSensitiveConfigKeyNginxKey1: "encrypted",
			}).
			Return(map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{
				nginxStaticSensitiveConfigKeyNginxKey1: "val1",
			}, nil)

		sut := NewStateDiffUseCase(nil, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, encryptionAdapterMock)

		// when
		clusterState, err := sut.collectClusterState(testCtx, effectiveBlueprint)

		// then
		assert.NoError(t, err)
		assert.Equal(t, ecosystem.ClusterState{
			GlobalConfig: map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry{},
			DoguConfig:   map[common.DoguConfigKey]*ecosystem.DoguConfigEntry{},
			EncryptedDoguConfig: map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry{
				nginxStaticSensitiveConfigKeyNginxKey1: encryptedEntry,
			},
			DecryptedSensitiveDoguConfig: map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{
				nginxStaticSensitiveConfigKeyNginxKey1: "val1",
			},
		}, clusterState)
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
		doguConfigNotFoundError := domainservice.NewNotFoundError(assert.AnError, "dogu config not found")

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, nil)
		globalConfigRepoMock := newMockGlobalConfigEntryRepository(t)
		globalConfigRepoMock.EXPECT().
			GetAllByKey(testCtx, effectiveBlueprint.Config.Global.GetGlobalConfigKeys()).
			Return(
				map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry{},
				globalConfigNotFoundError,
			)
		doguConfigRepoMock := newMockDoguConfigEntryRepository(t)
		doguConfigRepoMock.EXPECT().
			GetAllByKey(testCtx, effectiveBlueprint.Config.GetDoguConfigKeys()).
			Return(
				map[common.DoguConfigKey]*ecosystem.DoguConfigEntry{},
				doguConfigNotFoundError,
			)
		sensitiveDoguConfigRepoMock := newMockSensitiveDoguConfigEntryRepository(t)
		sensitiveDoguConfigRepoMock.EXPECT().
			GetAllByKey(testCtx, effectiveBlueprint.Config.GetSensitiveDoguConfigKeys()).
			Return(
				map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry{},
				internalTestError,
			)
		encryptionAdapterMock := newMockConfigEncryptionAdapter(t)

		sut := NewStateDiffUseCase(nil, doguInstallRepoMock, componentInstallRepoMock, globalConfigRepoMock, doguConfigRepoMock, sensitiveDoguConfigRepoMock, encryptionAdapterMock)

		// when
		clusterState, err := sut.collectClusterState(testCtx, effectiveBlueprint)

		// then
		assert.ErrorIs(t, err, internalTestError)
		assert.ErrorIs(t, err, globalConfigNotFoundError)
		assert.ErrorIs(t, err, doguConfigNotFoundError)
		assert.Equal(t, ecosystem.ClusterState{}, clusterState)
	})
}
