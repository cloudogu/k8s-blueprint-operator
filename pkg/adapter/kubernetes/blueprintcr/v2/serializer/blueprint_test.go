package serializer

import (
	"fmt"
	"testing"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"

	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

var (
	version3211, _ = core.ParseVersion("3.2.1-1")
	version3212, _ = core.ParseVersion("3.2.1-2")
	version1233, _ = core.ParseVersion("1.2.3-3")
)

var (
	trueVar  = true
	falseVar = false
)

func TestConvertToBlueprintDTO(t *testing.T) {
	t.Run("convert dogus", func(t *testing.T) {
		dogus := []domain.Dogu{
			{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "dogu1"}, Version: &version3211, Absent: true},
			{Name: cescommons.QualifiedName{Namespace: "premium", SimpleName: "dogu3"}, Version: &version3212, Absent: false},
		}
		blueprint := domain.EffectiveBlueprint{
			Dogus: dogus,
		}

		//when
		blueprintV2 := ConvertToBlueprintDTO(blueprint)

		//then
		convertedDogus := []crd.Dogu{
			{Name: "official/dogu1", Version: &version3211.Raw, Absent: &trueVar},
			{Name: "premium/dogu3", Version: &version3212.Raw, Absent: &falseVar},
		}
		expectedManifest := crd.BlueprintManifest{
			Dogus: convertedDogus,
		}
		assert.Empty(t, cmp.Diff(expectedManifest, blueprintV2))
	})

	t.Run("convert config", func(t *testing.T) {
		value42 := "42"
		blueprint := domain.EffectiveBlueprint{
			Config: domain.Config{
				Dogus: map[cescommons.SimpleName]domain.DoguConfigEntries{
					"my-dogu": {
						{
							Key:   "config",
							Value: (*config.Value)(&value42),
						},
						{
							Key:       "sensitive-config",
							Sensitive: true,
							SecretRef: &domain.SensitiveValueRef{
								SecretName: "mySecret",
								SecretKey:  "myKey",
							},
						},
					},
				},
				Global: domain.GlobalConfigEntries{
					{
						Key:    "test/key",
						Absent: true,
					},
				},
			},
		}
		blueprintV2 := ConvertToBlueprintDTO(blueprint)

		expectedManifest := crd.BlueprintManifest{
			Dogus: []crd.Dogu{},
			Config: &crd.Config{
				Dogus: map[string][]crd.ConfigEntry{
					"my-dogu": {
						crd.ConfigEntry{
							Key:   "config",
							Value: &value42,
						},
						crd.ConfigEntry{
							Key:       "sensitive-config",
							Sensitive: &trueVar,
							SecretRef: &crd.SecretReference{
								Name: "mySecret",
								Key:  "myKey",
							},
						},
					},
				},
				Global: []crd.ConfigEntry{
					{
						Key:    "test/key",
						Absent: &trueVar,
					},
				},
			},
		}
		assert.Empty(t, cmp.Diff(expectedManifest, blueprintV2))
	})
}

func TestConvertToEffectiveBlueprintDomain(t *testing.T) {
	t.Run("all ok", func(t *testing.T) {
		//given
		subfolder := "subfolder"
		convertedDogus := []crd.Dogu{
			{Name: "official/dogu1", Version: &version3211.Raw, Absent: &trueVar},
			{Name: "official/dogu2", Absent: &trueVar},
			{Name: "premium/dogu3", Version: &version3212.Raw, Absent: &falseVar},
			{Name: "premium/dogu4", Version: &version1233.Raw, Absent: &falseVar},
			{
				Name:    "premium/dogu5",
				Version: &version1233.Raw,
				Absent:  &falseVar,
				PlatformConfig: &crd.PlatformConfig{
					ResourceConfig:     nil,
					ReverseProxyConfig: nil,
					AdditionalMountsConfig: []crd.AdditionalMount{
						{
							SourceType: crd.DataSourceConfigMap,
							Name:       "config",
							Volume:     "volume",
							Subfolder:  &subfolder,
						},
					},
				},
			},
		}

		value42 := "42"
		dto := crd.BlueprintManifest{
			Dogus: convertedDogus,
			Config: &crd.Config{
				Dogus: map[string][]crd.ConfigEntry{
					"my-dogu": {
						crd.ConfigEntry{
							Key:   "config",
							Value: &value42,
						},
						crd.ConfigEntry{
							Key:       "sensitive-config",
							Sensitive: &trueVar,
							SecretRef: &crd.SecretReference{
								Name: "mySecret",
								Key:  "myKey",
							},
						},
					},
				},
				Global: []crd.ConfigEntry{
					{
						Key:    "test/key",
						Absent: &trueVar,
					},
				},
			},
		}
		//when
		blueprint, err := ConvertToEffectiveBlueprintDomain(&dto)

		//then
		dogus := []domain.Dogu{
			{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "dogu1"}, Version: &version3211, Absent: true},
			{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "dogu2"}, Absent: true},
			{Name: cescommons.QualifiedName{Namespace: "premium", SimpleName: "dogu3"}, Version: &version3212, Absent: false},
			{Name: cescommons.QualifiedName{Namespace: "premium", SimpleName: "dogu4"}, Version: &version1233},
			{
				Name:    cescommons.QualifiedName{Namespace: "premium", SimpleName: "dogu5"},
				Version: &version1233,
				AdditionalMounts: []ecosystem.AdditionalMount{
					{
						SourceType: ecosystem.DataSourceConfigMap,
						Name:       "config",
						Volume:     "volume",
						Subfolder:  subfolder,
					},
				},
			},
		}

		expected := domain.EffectiveBlueprint{
			Dogus: dogus,
			Config: domain.Config{
				Dogus: map[cescommons.SimpleName]domain.DoguConfigEntries{
					"my-dogu": {
						{
							Key:   "config",
							Value: (*config.Value)(&value42),
						},
						{
							Key:       "sensitive-config",
							Sensitive: true,
							SecretRef: &domain.SensitiveValueRef{
								SecretName: "mySecret",
								SecretKey:  "myKey",
							},
						},
					},
				},
				Global: domain.GlobalConfigEntries{
					{
						Key:    "test/key",
						Absent: true,
					},
				},
			},
		}

		require.NoError(t, err)
		assert.Empty(t, cmp.Diff(expected, blueprint))
	})

	t.Run("when manifest is nil return empty effective blueprint", func(t *testing.T) {
		//when
		blueprint, err := ConvertToEffectiveBlueprintDomain(nil)

		//then
		require.NoError(t, err)
		assert.Equal(t, domain.EffectiveBlueprint{}, blueprint)
	})

	t.Run("when convert dogu error return empty effective blueprint and error", func(t *testing.T) {
		//given
		dto := crd.BlueprintManifest{
			Dogus: []crd.Dogu{
				{Name: "dogu1", Version: &version3211.Raw}, // name contains no "/"
			},
		}

		//when
		blueprint, err := ConvertToEffectiveBlueprintDomain(&dto)

		//then
		require.Error(t, err)
		assert.ErrorContains(t, err, "cannot deserialize effective blueprint")
		assert.Equal(t, domain.EffectiveBlueprint{}, blueprint)
	})
}

func TestConvertToBlueprintDomain(t *testing.T) {
	t.Run("convert empty", func(t *testing.T) {
		given := crd.BlueprintManifest{}
		converted, err := ConvertToBlueprintDomain(given)

		require.NoError(t, err)
		assert.Equal(t, domain.Blueprint{}, converted)
	})

	t.Run("error converting", func(t *testing.T) {
		versionUnparsable := "unparsable"
		given := crd.BlueprintManifest{
			Dogus: []crd.Dogu{
				{
					Name:    "official/redmine",
					Version: &versionUnparsable,
				},
			},
		}
		_, err := ConvertToBlueprintDomain(given)

		require.Error(t, err)
		assert.ErrorContains(t, err, "cannot deserialize blueprint")
	})

	t.Run("convert dogu", func(t *testing.T) {
		given := crd.BlueprintManifest{
			Dogus: []crd.Dogu{
				{
					Name:    "official/redmine",
					Version: &version1233.Raw,
				},
			},
		}
		converted, err := ConvertToBlueprintDomain(given)

		require.NoError(t, err)
		assert.Equal(t, domain.Blueprint{
			Dogus: []domain.Dogu{
				{
					Name:    cescommons.QualifiedName{Namespace: "official", SimpleName: "redmine"},
					Version: &version1233,
				},
			},
		}, converted)
	})
}

func TestConvertToBlueprintMaskDomain(t *testing.T) {
	type args struct {
		mask *crd.BlueprintMaskManifest
	}
	tests := []struct {
		name    string
		args    args
		want    domain.BlueprintMask
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "nil",
			args:    args{mask: nil},
			want:    domain.BlueprintMask{},
			wantErr: assert.NoError,
		},
		{
			name:    "empty",
			args:    args{mask: &crd.BlueprintMaskManifest{}},
			want:    domain.BlueprintMask{},
			wantErr: assert.NoError,
		},
		{
			name: "will convert a MaskDogu",
			args: args{mask: &crd.BlueprintMaskManifest{
				Dogus: []crd.MaskDogu{
					{
						Name:    "official/ldap",
						Version: &version1233.Raw,
					},
				},
			}},
			want: domain.BlueprintMask{
				Dogus: []domain.MaskDogu{
					{
						Name: cescommons.QualifiedName{
							Namespace:  "official",
							SimpleName: "ldap",
						},
						Version: version1233,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error if invalid mask",
			args: args{mask: &crd.BlueprintMaskManifest{
				Dogus: []crd.MaskDogu{
					{
						Name:    "invalid name",
						Version: &version1233.Raw,
					},
				},
			}},
			want: domain.BlueprintMask{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "cannot deserialize blueprint mask", i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToBlueprintMaskDomain(tt.args.mask)
			if !tt.wantErr(t, err, fmt.Sprintf("ConvertToBlueprintMaskDomain(%v)", tt.args.mask)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ConvertToBlueprintMaskDomain(%v)", tt.args.mask)
		})
	}
}
