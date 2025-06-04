package domainservice

import (
	_ "embed"
	"github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

//go:embed testdata/k8s-nginx-static-1-26-3-2.json
var doguSpecK8sNginxStaticString string

var doguSpecK8sNginxStatic, _, _ = core.ReadDoguFromString(doguSpecK8sNginxStaticString)

func TestValidateAdditionalMountsDomainUseCase_ValidateAdditionalMounts(t *testing.T) {
	t.Run("no dogus", func(t *testing.T) {
		//given
		registry := NewMockRemoteDoguRegistry(t)
		useCase := NewValidateAdditionalMountsDomainUseCase(registry)
		blueprint := domain.EffectiveBlueprint{}
		//when
		err := useCase.ValidateAdditionalMounts(ctx, blueprint)
		//then
		require.NoError(t, err)
	})

	t.Run("no dogus with additional mounts", func(t *testing.T) {
		//given
		registry := NewMockRemoteDoguRegistry(t)
		useCase := NewValidateAdditionalMountsDomainUseCase(registry)
		blueprint := domain.EffectiveBlueprint{
			Dogus: []domain.Dogu{
				{
					Name:        k8sNginxStatic,
					Version:     version1_26_3_2,
					TargetState: domain.TargetStatePresent,
				},
			},
		}
		//when
		err := useCase.ValidateAdditionalMounts(ctx, blueprint)
		//then
		require.NoError(t, err)
	})

	t.Run("error loading dogu spec", func(t *testing.T) {
		//given
		registry := NewMockRemoteDoguRegistry(t)
		useCase := NewValidateAdditionalMountsDomainUseCase(registry)
		blueprint := domain.EffectiveBlueprint{
			Dogus: []domain.Dogu{
				{
					Name:        k8sNginxStatic,
					Version:     version1_26_3_2,
					TargetState: domain.TargetStatePresent,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "myConfigMap",
							Volume:     "customhtml",
						},
					},
				},
			},
		}

		dogusToLoad := []dogu.QualifiedVersion{
			{Name: k8sNginxStatic, Version: version1_26_3_2},
		}
		registry.EXPECT().GetDogus(ctx, dogusToLoad).Return(nil, assert.AnError)

		//when
		err := useCase.ValidateAdditionalMounts(ctx, blueprint)

		//then
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("valid additional mount", func(t *testing.T) {
		//given
		registry := NewMockRemoteDoguRegistry(t)
		useCase := NewValidateAdditionalMountsDomainUseCase(registry)
		blueprint := domain.EffectiveBlueprint{
			Dogus: []domain.Dogu{
				{
					Name:        k8sNginxStatic,
					Version:     version1_26_3_2,
					TargetState: domain.TargetStatePresent,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "myConfigMap",
							Volume:     "customhtml",
						},
					},
				},
			},
		}

		dogusToLoad := []dogu.QualifiedVersion{
			{Name: k8sNginxStatic, Version: version1_26_3_2},
		}
		doguSpecsToReturn := map[dogu.QualifiedName]*core.Dogu{
			k8sNginxStatic: doguSpecK8sNginxStatic,
		}
		registry.EXPECT().GetDogus(ctx, dogusToLoad).Return(doguSpecsToReturn, nil)

		//when
		err := useCase.ValidateAdditionalMounts(ctx, blueprint)

		//then
		require.NoError(t, err)
	})

	t.Run("unknown volume", func(t *testing.T) {
		//given
		registry := NewMockRemoteDoguRegistry(t)
		useCase := NewValidateAdditionalMountsDomainUseCase(registry)
		blueprint := domain.EffectiveBlueprint{
			Dogus: []domain.Dogu{
				{
					Name:        k8sNginxStatic,
					Version:     version1_26_3_2,
					TargetState: domain.TargetStatePresent,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "myConfigMap",
							Volume:     "unknownVolume",
						},
					},
				},
			},
		}

		dogusToLoad := []dogu.QualifiedVersion{
			{Name: k8sNginxStatic, Version: version1_26_3_2},
		}
		doguSpecsToReturn := map[dogu.QualifiedName]*core.Dogu{
			k8sNginxStatic: doguSpecK8sNginxStatic,
		}
		registry.EXPECT().GetDogus(ctx, dogusToLoad).Return(doguSpecsToReturn, nil)

		//when
		err := useCase.ValidateAdditionalMounts(ctx, blueprint)
		//then
		assert.ErrorContains(t, err, `additionalMounts are invalid`)
		assert.ErrorContains(t, err, `volume "unknownVolume" in additional mount for dogu "k8s/nginx-static" is invalid`)
		assert.ErrorContains(t, err, `["app.conf.d" "customhtml" "menu-json" "localConfig"]`)
	})
}
