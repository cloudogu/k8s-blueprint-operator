package application

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var postfixQualifiedDoguName = common.QualifiedDoguName{
	Namespace: "official",
	Name:      "postfix",
}

var ldapQualifiedDoguName = common.QualifiedDoguName{
	Namespace: "official",
	Name:      "ldap",
}

var nginxIngressQualifiedDoguName = common.QualifiedDoguName{
	Namespace: "k8s",
	Name:      "nginx-ingress",
}

var nginxStaticQualifiedDoguName = common.QualifiedDoguName{
	Namespace: "k8s",
	Name:      "nginx-static",
}

func TestStateDiffUseCase_DetermineStateDiff(t *testing.T) {
	t.Run("should fail to load blueprint spec", func(t *testing.T) {
		// given
		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(nil, assert.AnError)

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock)

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
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, assert.AnError)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot get installed dogus to determine state diff")
	})
	t.Run("should fail to get installed components", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{Id: "testBlueprint1"}

		blueprintRepoMock := newMockBlueprintSpecRepository(t)
		blueprintRepoMock.EXPECT().GetById(testCtx, "testBlueprint1").Return(blueprint, nil)

		doguInstallRepoMock := newMockDoguInstallationRepository(t)
		doguInstallRepoMock.EXPECT().GetAll(testCtx).Return(map[common.SimpleDoguName]*ecosystem.DoguInstallation{}, nil)
		componentInstallRepoMock := newMockComponentInstallationRepository(t)
		componentInstallRepoMock.EXPECT().GetAll(testCtx).Return(nil, assert.AnError)

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot get installed components to determine state diff")
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

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock)

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

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock)

		// when
		err := sut.DetermineStateDiff(testCtx, "testBlueprint1")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot save blueprint spec \"testBlueprint1\" after Determining the state diff to the ecosystem")
	})
	t.Run("should succeed", func(t *testing.T) {
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

		sut := NewStateDiffUseCase(blueprintRepoMock, doguInstallRepoMock, componentInstallRepoMock)

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
}

func mustParseVersion(t *testing.T, raw string) core.Version {
	version, err := core.ParseVersion(raw)
	assert.NoError(t, err)
	return version
}
