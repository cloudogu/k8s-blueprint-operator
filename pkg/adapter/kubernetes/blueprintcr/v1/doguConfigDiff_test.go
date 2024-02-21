package v1

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_convertToCombinedDoguConfigDiffDTO(t *testing.T) {
	tests := []struct {
		name        string
		domainModel domain.CombinedDoguConfigDiff
		want        CombinedDoguConfigDiff
	}{
		{
			name:        "should exit early if slices are empty",
			domainModel: domain.CombinedDoguConfigDiff{},
			want:        CombinedDoguConfigDiff{},
		},
		{
			name: "should convert multiple dogu config diffs",
			domainModel: domain.CombinedDoguConfigDiff{
				DoguConfigDiff: []domain.DoguConfigEntryDiff{
					{
						Key: common.DoguConfigKey{
							DoguName: "ldap",
							Key:      "container_config/memory_limit",
						},
						Actual: domain.DoguConfigValueState{
							Value:  "512m",
							Exists: true,
						},
						Expected: domain.DoguConfigValueState{
							Value:  "1024m",
							Exists: true,
						},
						Action: domain.ConfigActionSet,
					},
					{
						Key: common.DoguConfigKey{
							DoguName: "ldap",
							Key:      "container_config/swap_limit",
						},
						Actual: domain.DoguConfigValueState{
							Exists: false,
						},
						Expected: domain.DoguConfigValueState{
							Value:  "512m",
							Exists: true,
						},
						Action: domain.ConfigActionSet,
					},
				},
			},
			want: CombinedDoguConfigDiff{
				DoguConfigDiff: DoguConfigDiff{
					{
						Key: "container_config/memory_limit",
						Actual: DoguConfigValueState{
							Value:  "512m",
							Exists: true,
						},
						Expected: DoguConfigValueState{
							Value:  "1024m",
							Exists: true,
						},
						NeededAction: ConfigActionSet,
					},
					{
						Key: "container_config/swap_limit",
						Actual: DoguConfigValueState{
							Exists: false,
						},
						Expected: DoguConfigValueState{
							Value:  "512m",
							Exists: true,
						},
						NeededAction: ConfigActionSet,
					},
				},
			},
		},
		{
			name: "should convert multiple sensitive dogu config diffs",
			domainModel: domain.CombinedDoguConfigDiff{
				SensitiveDoguConfigDiff: []domain.SensitiveDoguConfigEntryDiff{
					{
						Key: common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
							DoguName: "ldap",
							Key:      "container_config/memory_limit",
						}},
						Actual: domain.EncryptedDoguConfigValueState{
							Value:  "512m",
							Exists: true,
						},
						Expected: domain.EncryptedDoguConfigValueState{
							Value:  "1024m",
							Exists: true,
						},
						Action: domain.ConfigActionSet,
					},
					{
						Key: common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
							DoguName: "ldap",
							Key:      "container_config/swap_limit",
						}},
						Actual: domain.EncryptedDoguConfigValueState{
							Exists: false,
						},
						Expected: domain.EncryptedDoguConfigValueState{
							Value:  "512m",
							Exists: true,
						},
						Action: domain.ConfigActionSet,
					},
				},
			},
			want: CombinedDoguConfigDiff{
				SensitiveDoguConfigDiff: SensitiveDoguConfigDiff{
					{
						Key: "container_config/memory_limit",
						Actual: DoguConfigValueState{
							Value:  "512m",
							Exists: true,
						},
						Expected: DoguConfigValueState{
							Value:  "1024m",
							Exists: true,
						},
						NeededAction: ConfigActionSet,
					},
					{
						Key: "container_config/swap_limit",
						Actual: DoguConfigValueState{
							Exists: false,
						},
						Expected: DoguConfigValueState{
							Value:  "512m",
							Exists: true,
						},
						NeededAction: ConfigActionSet,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, convertToCombinedDoguConfigDiffDTO(tt.domainModel), "convertToCombinedDoguConfigDiffDTO(%v)", tt.domainModel)
		})
	}
}
