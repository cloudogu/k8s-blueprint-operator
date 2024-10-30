package application

import (
	"context"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
	"time"
)

const blueprintId = "blueprint1"

var version3211, _ = core.ParseVersion("3.2.1-1")
var version3212, _ = core.ParseVersion("3.2.1-2")

var postgresqlQualifiedName = common.QualifiedDoguName{
	Namespace:  "official",
	SimpleName: "postgresql",
}

// TODO Add test for proxy and volume actions
func TestDoguInstallationUseCase_applyDoguState(t *testing.T) {
	t.Run("action none", func(t *testing.T) {
		// given
		sut := NewDoguInstallationUseCase(nil, nil, nil)

		// when
		err := sut.applyDoguState(testCtx, domain.DoguDiff{
			DoguName: "postgresql",
			Actual: domain.DoguDiffState{
				Namespace:         "official",
				Version:           version3211,
				InstallationState: domain.TargetStatePresent,
			},
			Expected: domain.DoguDiffState{
				Namespace:         "official",
				Version:           version3211,
				InstallationState: domain.TargetStatePresent,
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
			RewriteTarget:    "/",
			AdditionalConfig: "additional",
		}
		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().
			Create(testCtx, ecosystem.InstallDogu(postgresqlQualifiedName, version3211, &volumeSize, config)).
			Return(nil)

		sut := NewDoguInstallationUseCase(nil, doguRepoMock, nil)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Actual: domain.DoguDiffState{
					Namespace:         "official",
					Version:           version3211,
					InstallationState: domain.TargetStateAbsent,
				},
				Expected: domain.DoguDiffState{
					Namespace:          "official",
					Version:            version3211,
					InstallationState:  domain.TargetStatePresent,
					MinVolumeSize:      &volumeSize,
					ReverseProxyConfig: config,
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
			Delete(testCtx, common.SimpleDoguName("postgresql")).
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
					Version: version3212,
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
					Version: version3211,
				},
				NeededActions: []domain.Action{domain.ActionDowngrade},
			},
			dogu,
			domain.BlueprintConfiguration{},
		)

		// then
		require.ErrorContains(t, err, getNoDowngradesExplanationTextForDogus())
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
			ReverseProxyConfig: ecosystem.ReverseProxyConfig{
				MaxBodySize: &expectedProxyBodySize,
			},
		}

		dogu := &ecosystem.DoguInstallation{
			Name: postgresqlQualifiedName,
			ReverseProxyConfig: ecosystem.ReverseProxyConfig{
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
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
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
		target := ecosystem.RewriteTarget("")
		expectedTarget := ecosystem.RewriteTarget("/")
		expectedDogu := &ecosystem.DoguInstallation{
			Name: postgresqlQualifiedName,
			ReverseProxyConfig: ecosystem.ReverseProxyConfig{
				RewriteTarget: expectedTarget,
			},
		}

		dogu := &ecosystem.DoguInstallation{
			Name: postgresqlQualifiedName,
			ReverseProxyConfig: ecosystem.ReverseProxyConfig{
				RewriteTarget: target,
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
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
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
		additionalConfig := ecosystem.AdditionalConfig("")
		expectedAdditionalConfig := ecosystem.AdditionalConfig("snippet")
		expectedDogu := &ecosystem.DoguInstallation{
			Name: postgresqlQualifiedName,
			ReverseProxyConfig: ecosystem.ReverseProxyConfig{
				AdditionalConfig: expectedAdditionalConfig,
			},
		}

		dogu := &ecosystem.DoguInstallation{
			Name: postgresqlQualifiedName,
			ReverseProxyConfig: ecosystem.ReverseProxyConfig{
				AdditionalConfig: additionalConfig,
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
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
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

	t.Run("should process multiple update actions", func(t *testing.T) {
		volumeSize := resource.MustParse("2Gi")
		expectedVolumeSize := resource.MustParse("3Gi")
		proxyBodySize := resource.MustParse("2G")
		expectedProxyBodySize := resource.MustParse("3G")
		target := ecosystem.RewriteTarget("")
		expectedTarget := ecosystem.RewriteTarget("/")
		additionalConfig := ecosystem.AdditionalConfig("")
		expectedAdditionalConfig := ecosystem.AdditionalConfig("snippet")
		expectedDogu := &ecosystem.DoguInstallation{
			Name:          postgresqlQualifiedName,
			Version:       version3212,
			MinVolumeSize: &expectedVolumeSize,
			ReverseProxyConfig: ecosystem.ReverseProxyConfig{
				MaxBodySize:      &expectedProxyBodySize,
				RewriteTarget:    expectedTarget,
				AdditionalConfig: expectedAdditionalConfig,
			},
		}

		dogu := &ecosystem.DoguInstallation{
			Name:          postgresqlQualifiedName,
			Version:       version3211,
			MinVolumeSize: &volumeSize,
			ReverseProxyConfig: ecosystem.ReverseProxyConfig{
				MaxBodySize:      &proxyBodySize,
				RewriteTarget:    target,
				AdditionalConfig: additionalConfig,
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
					Version:       version3212,
					MinVolumeSize: &expectedVolumeSize,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
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
		assert.Equal(t, common.DoguNamespace("premium"), dogu.Name.Namespace)
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
	t.Run("cannot load blueprintSpec", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(nil, assert.AnError)

		doguRepoMock := newMockDoguInstallationRepository(t)

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock, nil)

		// when
		err := sut.ApplyDoguStates(testCtx, blueprintId)

		// then
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("cannot load doguInstallations", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(nil, nil)

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().GetAll(testCtx).Return(nil, assert.AnError)

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock, nil)

		// when
		err := sut.ApplyDoguStates(testCtx, blueprintId)

		// then
		require.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot load dogu installations")
	})

	t.Run("success", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(&domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				DoguDiffs: []domain.DoguDiff{
					{
						DoguName:      "postgresql",
						NeededActions: []domain.Action{},
					},
				},
			},
			Config: domain.BlueprintConfiguration{},
		}, nil)

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().GetAll(testCtx).Return(map[common.SimpleDoguName]*ecosystem.DoguInstallation{}, nil)

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock, nil)

		// when
		err := sut.ApplyDoguStates(testCtx, blueprintId)

		// then
		require.NoError(t, err)
	})

	t.Run("action error", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(&domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				DoguDiffs: []domain.DoguDiff{
					{
						DoguName:      "postgresql",
						NeededActions: []domain.Action{domain.ActionDowngrade},
					},
				},
			},
			Config: domain.BlueprintConfiguration{},
		}, nil)

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().GetAll(testCtx).Return(map[common.SimpleDoguName]*ecosystem.DoguInstallation{
			"postgresql": {
				Name:          postgresqlQualifiedName,
				Version:       version3211,
				UpgradeConfig: ecosystem.UpgradeConfig{},
			},
		}, nil)

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock, nil)

		// when
		err := sut.ApplyDoguStates(testCtx, blueprintId)

		// then
		require.ErrorContains(t, err, getNoDowngradesExplanationTextForDogus())
		require.ErrorContains(t, err, "an error occurred while applying dogu state to the ecosystem")
	})
}

func TestDoguInstallationUseCase_WaitForHealthyDogus(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		// given
		doguRepoMock := newMockDoguInstallationRepository(t)
		timedCtx, cancel := context.WithTimeout(testCtx, 10*time.Millisecond)
		defer cancel()
		doguRepoMock.EXPECT().GetAll(timedCtx).Return(map[common.SimpleDoguName]*ecosystem.DoguInstallation{}, nil)

		waitConfigMock := newMockHealthWaitConfigProvider(t)
		waitConfigMock.EXPECT().GetWaitConfig(timedCtx).Return(ecosystem.WaitConfig{Interval: time.Millisecond}, nil)

		sut := DoguInstallationUseCase{
			blueprintSpecRepo:  nil,
			doguRepo:           doguRepoMock,
			waitConfigProvider: waitConfigMock,
		}

		// when
		result, err := sut.WaitForHealthyDogus(timedCtx)

		// then
		require.NoError(t, err)
		assert.True(t, result.AllHealthy())
	})

	t.Run("fail to get health check interval", func(t *testing.T) {
		t.Parallel()
		// given
		waitConfigMock := newMockHealthWaitConfigProvider(t)
		waitConfigMock.EXPECT().GetWaitConfig(testCtx).Return(ecosystem.WaitConfig{}, assert.AnError)

		sut := DoguInstallationUseCase{
			blueprintSpecRepo:  nil,
			doguRepo:           nil,
			waitConfigProvider: waitConfigMock,
		}

		// when
		_, err := sut.WaitForHealthyDogus(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get health check interval")
	})

	t.Run("timeout", func(t *testing.T) {
		t.Parallel()
		// given
		doguRepoMock := newMockDoguInstallationRepository(t)
		timedCtx, cancel := context.WithTimeout(testCtx, 0*time.Millisecond)
		defer cancel()
		// return unhealthy result
		doguRepoMock.EXPECT().GetAll(timedCtx).Return(map[common.SimpleDoguName]*ecosystem.DoguInstallation{
			"postgresql": {Health: ecosystem.DoguStatusInstalling},
		}, nil).Maybe()

		waitConfigMock := newMockHealthWaitConfigProvider(t)
		waitConfigMock.EXPECT().GetWaitConfig(timedCtx).Return(ecosystem.WaitConfig{Interval: 5 * time.Millisecond}, nil)

		sut := DoguInstallationUseCase{
			blueprintSpecRepo:  nil,
			doguRepo:           doguRepoMock,
			waitConfigProvider: waitConfigMock,
		}

		// when
		result, err := sut.WaitForHealthyDogus(timedCtx)

		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Equal(t, ecosystem.DoguHealthResult{}, result)
	})

	t.Run("cannot load dogus", func(t *testing.T) {
		t.Parallel()
		// given
		doguRepoMock := newMockDoguInstallationRepository(t)
		timedCtx, cancel := context.WithTimeout(testCtx, 10*time.Millisecond)
		defer cancel()
		doguRepoMock.EXPECT().GetAll(timedCtx).Return(nil, assert.AnError).Maybe()

		waitConfigMock := newMockHealthWaitConfigProvider(t)
		waitConfigMock.EXPECT().GetWaitConfig(timedCtx).Return(ecosystem.WaitConfig{Interval: time.Millisecond}, nil)

		sut := DoguInstallationUseCase{
			blueprintSpecRepo:  nil,
			doguRepo:           doguRepoMock,
			waitConfigProvider: waitConfigMock,
		}

		// when
		result, err := sut.WaitForHealthyDogus(timedCtx)

		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, ecosystem.DoguHealthResult{}, result)
	})

}
