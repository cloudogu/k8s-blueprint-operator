package application

import (
	"fmt"
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
)

const blueprintId = "blueprint1"

var version3211, _ = core.ParseVersion("3.2.1-1")
var version3212, _ = core.ParseVersion("3.2.1-2")

var postgresqlQualifiedName = cescommons.QualifiedName{
	Namespace:  "official",
	SimpleName: "postgresql",
}

var (
	rewriteTarget    = "/"
	additionalConfig = "additional"
	subfolder        = "subfolder"
	subfolder2       = "secsubfolder"
)

func TestDoguInstallationUseCase_applyDoguState(t *testing.T) {
	t.Run("action none", func(t *testing.T) {
		// given
		sut := NewDoguInstallationUseCase(nil, nil, nil)

		// when
		err := sut.applyDoguState(testCtx, domain.DoguDiff{
			DoguName: "postgresql",
			Actual: domain.DoguDiffState{
				Namespace: "official",
				Version:   &version3211,
				Absent:    false,
			},
			Expected: domain.DoguDiffState{
				Namespace: "official",
				Version:   &version3211,
				Absent:    false,
			},
			NeededActions: []domain.Action{},
		}, &ecosystem.DoguInstallation{
			Name:    postgresqlQualifiedName,
			Version: version3211,
		}, domain.BlueprintConfiguration{})

		// then
		require.NoError(t, err)
	})

	t.Run("action install", func(t *testing.T) {
		volumeSize := resource.MustParse("2Gi")
		bodySize := resource.MustParse("2G")
		config := ecosystem.ReverseProxyConfig{
			MaxBodySize:      &bodySize,
			RewriteTarget:    &rewriteTarget,
			AdditionalConfig: &additionalConfig,
		}
		additionalMounts := []ecosystem.AdditionalMount{
			{
				SourceType: ecosystem.DataSourceConfigMap,
				Name:       "configmap",
				Volume:     "volume",
				Subfolder:  &subfolder,
			},
		}

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().
			Create(testCtx,
				ecosystem.InstallDogu(postgresqlQualifiedName, &version3211, &volumeSize, &config, additionalMounts)).
			Return(nil)

		sut := NewDoguInstallationUseCase(nil, doguRepoMock, nil)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Actual: domain.DoguDiffState{
					Namespace: "official",
					Version:   &version3211,
					Absent:    true,
				},
				Expected: domain.DoguDiffState{
					Namespace:          "official",
					Version:            &version3211,
					Absent:             false,
					MinVolumeSize:      &volumeSize,
					ReverseProxyConfig: &config,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "configmap",
							Volume:     "volume",
							Subfolder:  &subfolder,
						},
					},
				},
				NeededActions: []domain.Action{domain.ActionInstall},
			},
			nil,
			domain.BlueprintConfiguration{},
		)

		// then
		require.NoError(t, err)
	})

	t.Run("action uninstall", func(t *testing.T) {
		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().
			Delete(testCtx, cescommons.SimpleName("postgresql")).
			Return(nil)

		sut := NewDoguInstallationUseCase(nil, doguRepoMock, nil)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName:      "postgresql",
				NeededActions: []domain.Action{domain.ActionUninstall},
			},
			&ecosystem.DoguInstallation{
				Name:    postgresqlQualifiedName,
				Version: version3211,
			},
			domain.BlueprintConfiguration{},
		)

		// then
		require.NoError(t, err)
	})

	t.Run("action upgrade", func(t *testing.T) {
		dogu := &ecosystem.DoguInstallation{
			Name:    postgresqlQualifiedName,
			Version: version3211,
		}
		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().
			Update(testCtx, dogu).
			Return(nil)

		sut := NewDoguInstallationUseCase(nil, doguRepoMock, nil)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Expected: domain.DoguDiffState{
					Version: &version3212,
				},
				NeededActions: []domain.Action{domain.ActionUpgrade},
			},
			dogu,
			domain.BlueprintConfiguration{},
		)

		// then
		require.NoError(t, err)
		assert.Equal(t, version3212, dogu.Version)
	})

	t.Run("action downgrade", func(t *testing.T) {

		dogu := &ecosystem.DoguInstallation{
			Name:    postgresqlQualifiedName,
			Version: version3212,
		}

		sut := NewDoguInstallationUseCase(nil, nil, nil)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Expected: domain.DoguDiffState{
					Version: &version3211,
				},
				NeededActions: []domain.Action{domain.ActionDowngrade},
			},
			dogu,
			domain.BlueprintConfiguration{},
		)

		// then
		require.ErrorContains(t, err, fmt.Sprintf(noDowngradesExplanationTextFmt, "dogu", "dogus"))
		assert.Equal(t, version3212, dogu.Version)
	})

	t.Run("action update volume size", func(t *testing.T) {
		volumeSize := resource.MustParse("2Gi")
		expectedVolumeSize := resource.MustParse("3Gi")
		expectedDogu := &ecosystem.DoguInstallation{
			Name:          postgresqlQualifiedName,
			MinVolumeSize: &expectedVolumeSize,
		}

		dogu := &ecosystem.DoguInstallation{
			Name:          postgresqlQualifiedName,
			MinVolumeSize: &volumeSize,
		}

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().Update(testCtx, expectedDogu).Return(nil)

		sut := NewDoguInstallationUseCase(nil, doguRepoMock, nil)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Expected: domain.DoguDiffState{
					MinVolumeSize: &expectedVolumeSize,
				},
				NeededActions: []domain.Action{domain.ActionUpdateDoguResourceMinVolumeSize},
			},
			dogu,
			domain.BlueprintConfiguration{},
		)

		// then
		require.NoError(t, err)
	})

	t.Run("action update proxy body size", func(t *testing.T) {
		proxyBodySize := resource.MustParse("2G")
		expectedProxyBodySize := resource.MustParse("3G")
		expectedDogu := &ecosystem.DoguInstallation{
			Name: postgresqlQualifiedName,
			ReverseProxyConfig: &ecosystem.ReverseProxyConfig{
				MaxBodySize: &expectedProxyBodySize,
			},
		}

		dogu := &ecosystem.DoguInstallation{
			Name: postgresqlQualifiedName,
			ReverseProxyConfig: &ecosystem.ReverseProxyConfig{
				MaxBodySize: &proxyBodySize,
			},
		}

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().Update(testCtx, expectedDogu).Return(nil)

		sut := NewDoguInstallationUseCase(nil, doguRepoMock, nil)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Expected: domain.DoguDiffState{
					ReverseProxyConfig: &ecosystem.ReverseProxyConfig{
						MaxBodySize: &expectedProxyBodySize,
					},
				},
				NeededActions: []domain.Action{domain.ActionUpdateDoguProxyBodySize},
			},
			dogu,
			domain.BlueprintConfiguration{},
		)

		// then
		require.NoError(t, err)
	})

	t.Run("action update proxy rewrite target", func(t *testing.T) {
		expectedTarget := &rewriteTarget
		expectedDogu := &ecosystem.DoguInstallation{
			Name: postgresqlQualifiedName,
			ReverseProxyConfig: &ecosystem.ReverseProxyConfig{
				RewriteTarget: expectedTarget,
			},
		}

		dogu := &ecosystem.DoguInstallation{
			Name: postgresqlQualifiedName,
			ReverseProxyConfig: &ecosystem.ReverseProxyConfig{
				RewriteTarget: nil,
			},
		}

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().Update(testCtx, expectedDogu).Return(nil)

		sut := NewDoguInstallationUseCase(nil, doguRepoMock, nil)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Expected: domain.DoguDiffState{
					ReverseProxyConfig: &ecosystem.ReverseProxyConfig{
						RewriteTarget: expectedTarget,
					},
				},
				NeededActions: []domain.Action{domain.ActionUpdateDoguProxyRewriteTarget},
			},
			dogu,
			domain.BlueprintConfiguration{},
		)

		// then
		require.NoError(t, err)
	})

	t.Run("action update proxy additional config", func(t *testing.T) {
		expectedAdditionalConfig := &additionalConfig
		expectedDogu := &ecosystem.DoguInstallation{
			Name: postgresqlQualifiedName,
			ReverseProxyConfig: &ecosystem.ReverseProxyConfig{
				AdditionalConfig: expectedAdditionalConfig,
			},
		}

		dogu := &ecosystem.DoguInstallation{
			Name: postgresqlQualifiedName,
			ReverseProxyConfig: &ecosystem.ReverseProxyConfig{
				AdditionalConfig: nil,
			},
		}

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().Update(testCtx, expectedDogu).Return(nil)

		sut := NewDoguInstallationUseCase(nil, doguRepoMock, nil)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Expected: domain.DoguDiffState{
					ReverseProxyConfig: &ecosystem.ReverseProxyConfig{
						AdditionalConfig: expectedAdditionalConfig,
					},
				},
				NeededActions: []domain.Action{domain.ActionUpdateDoguProxyAdditionalConfig},
			},
			dogu,
			domain.BlueprintConfiguration{},
		)

		// then
		require.NoError(t, err)
	})

	t.Run("should update additional mounts", func(t *testing.T) {
		expectedDogu := &ecosystem.DoguInstallation{
			Name: postgresqlQualifiedName,
			AdditionalMounts: []ecosystem.AdditionalMount{
				{
					SourceType: ecosystem.DataSourceConfigMap,
					Name:       "configmap",
					Volume:     "volume",
					Subfolder:  &subfolder,
				},
				{
					SourceType: ecosystem.DataSourceSecret,
					Name:       "secret",
					Volume:     "secvolume",
					Subfolder:  &subfolder2,
				},
			},
		}

		dogu := &ecosystem.DoguInstallation{
			Name: postgresqlQualifiedName,
			AdditionalMounts: []ecosystem.AdditionalMount{
				{
					SourceType: ecosystem.DataSourceConfigMap,
					Name:       "configmap",
					Volume:     "volume",
					Subfolder:  &subfolder,
				},
				{
					SourceType: ecosystem.DataSourceSecret,
					Name:       "secret",
					Volume:     "secvolume",
					Subfolder:  &subfolder2,
				},
			},
		}

		diff := domain.DoguDiff{
			DoguName: "postgresql",
			Expected: domain.DoguDiffState{
				AdditionalMounts: []ecosystem.AdditionalMount{
					{
						SourceType: ecosystem.DataSourceConfigMap,
						Name:       "configmap",
						Volume:     "volume",
						Subfolder:  &subfolder,
					},
					{
						SourceType: ecosystem.DataSourceSecret,
						Name:       "secret",
						Volume:     "secvolume",
						Subfolder:  &subfolder2,
					},
				},
			},
			NeededActions: []domain.Action{domain.ActionUpdateAdditionalMounts},
		}

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().Update(testCtx, expectedDogu).Return(nil)

		sut := NewDoguInstallationUseCase(nil, doguRepoMock, nil)

		// when
		err := sut.applyDoguState(testCtx, diff, dogu, domain.BlueprintConfiguration{})

		// then
		require.NoError(t, err)
	})

	t.Run("should process multiple update actions", func(t *testing.T) {
		volumeSize := resource.MustParse("2Gi")
		expectedVolumeSize := resource.MustParse("3Gi")
		proxyBodySize := resource.MustParse("2G")
		expectedProxyBodySize := resource.MustParse("3G")
		expectedTarget := &rewriteTarget
		expectedAdditionalConfig := &additionalConfig
		expectedDogu := &ecosystem.DoguInstallation{
			Name:          postgresqlQualifiedName,
			Version:       version3212,
			MinVolumeSize: &expectedVolumeSize,
			ReverseProxyConfig: &ecosystem.ReverseProxyConfig{
				MaxBodySize:      &expectedProxyBodySize,
				RewriteTarget:    expectedTarget,
				AdditionalConfig: expectedAdditionalConfig,
			},
		}

		dogu := &ecosystem.DoguInstallation{
			Name:          postgresqlQualifiedName,
			Version:       version3211,
			MinVolumeSize: &volumeSize,
			ReverseProxyConfig: &ecosystem.ReverseProxyConfig{
				MaxBodySize:      &proxyBodySize,
				RewriteTarget:    nil,
				AdditionalConfig: nil,
			},
		}

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().Update(testCtx, expectedDogu).Return(nil)

		sut := NewDoguInstallationUseCase(nil, doguRepoMock, nil)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Expected: domain.DoguDiffState{
					Version:       &version3212,
					MinVolumeSize: &expectedVolumeSize,
					ReverseProxyConfig: &ecosystem.ReverseProxyConfig{
						MaxBodySize:      &expectedProxyBodySize,
						RewriteTarget:    expectedTarget,
						AdditionalConfig: expectedAdditionalConfig,
					},
				},
				NeededActions: []domain.Action{
					domain.ActionUpgrade,
					domain.ActionUpdateDoguProxyAdditionalConfig,
					domain.ActionUpdateDoguProxyBodySize,
					domain.ActionUpdateDoguProxyRewriteTarget,
					domain.ActionUpdateDoguResourceMinVolumeSize,
				},
			},
			dogu,
			domain.BlueprintConfiguration{},
		)

		// then
		require.NoError(t, err)
	})

	t.Run("action SwitchNamespace not allowed", func(t *testing.T) {
		dogu := &ecosystem.DoguInstallation{
			Name:    postgresqlQualifiedName,
			Version: version3212,
		}

		sut := NewDoguInstallationUseCase(nil, nil, nil)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Expected: domain.DoguDiffState{
					Namespace: "premium",
				},
				NeededActions: []domain.Action{domain.ActionSwitchDoguNamespace},
			},
			dogu,
			domain.BlueprintConfiguration{
				AllowDoguNamespaceSwitch: false,
			},
		)

		// then
		require.ErrorContains(t, err, "not allowed to switch dogu namespace")
	})

	t.Run("action SwitchNamespace allowed", func(t *testing.T) {
		dogu := &ecosystem.DoguInstallation{
			Name:    postgresqlQualifiedName,
			Version: version3212,
		}
		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().Update(testCtx, dogu).Return(nil)

		sut := NewDoguInstallationUseCase(nil, doguRepoMock, nil)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Expected: domain.DoguDiffState{
					Namespace: "premium",
				},
				NeededActions: []domain.Action{domain.ActionSwitchDoguNamespace},
			},
			dogu,
			domain.BlueprintConfiguration{
				AllowDoguNamespaceSwitch: true,
			},
		)

		// then
		require.NoError(t, err)
		assert.Equal(t, cescommons.Namespace("premium"), dogu.Name.Namespace)
	})

	t.Run("unknown action", func(t *testing.T) {
		// given
		sut := NewDoguInstallationUseCase(nil, nil, nil)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Expected: domain.DoguDiffState{
					Namespace: "premium",
				},
				NeededActions: []domain.Action{"unknown"},
			},
			nil,
			domain.BlueprintConfiguration{},
		)

		// then
		require.ErrorContains(t, err, "cannot perform unknown action \"unknown\"")
	})

	t.Run("should no fail with no actions", func(t *testing.T) {
		// given
		sut := NewDoguInstallationUseCase(nil, nil, nil)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Expected: domain.DoguDiffState{
					Namespace: "premium",
				},
				NeededActions: []domain.Action{},
			},
			nil,
			domain.BlueprintConfiguration{},
		)

		// then
		require.NoError(t, err)
	})
}

func TestDoguInstallationUseCase_ApplyDoguStates(t *testing.T) {
	t.Run("cannot load doguInstallations", func(t *testing.T) {
		// given
		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().GetAll(testCtx).Return(nil, assert.AnError)

		sut := NewDoguInstallationUseCase(nil, doguRepoMock, nil)

		// when
		err := sut.ApplyDoguStates(testCtx, &domain.BlueprintSpec{})

		// then
		require.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot load dogu installations")
	})

	t.Run("success", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				DoguDiffs: []domain.DoguDiff{
					{
						DoguName:      "postgresql",
						NeededActions: []domain.Action{},
					},
				},
			},
			Config: domain.BlueprintConfiguration{},
		}
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().GetAll(testCtx).Return(map[cescommons.SimpleName]*ecosystem.DoguInstallation{}, nil)

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock, nil)

		// when
		err := sut.ApplyDoguStates(testCtx, blueprint)

		// then
		require.NoError(t, err)
	})

	t.Run("action error", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				DoguDiffs: []domain.DoguDiff{
					{
						DoguName:      "postgresql",
						NeededActions: []domain.Action{domain.ActionDowngrade},
					},
				},
			},
			Config: domain.BlueprintConfiguration{},
		}

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().GetAll(testCtx).Return(map[cescommons.SimpleName]*ecosystem.DoguInstallation{
			"postgresql": {
				Name:          postgresqlQualifiedName,
				Version:       version3211,
				UpgradeConfig: ecosystem.UpgradeConfig{},
			},
		}, nil)

		sut := NewDoguInstallationUseCase(nil, doguRepoMock, nil)

		// when
		err := sut.ApplyDoguStates(testCtx, blueprint)

		// then
		require.ErrorContains(t, err, fmt.Sprintf(noDowngradesExplanationTextFmt, "dogu", "dogus"))
		require.ErrorContains(t, err, "an error occurred while applying dogu state to the ecosystem")
	})
}
