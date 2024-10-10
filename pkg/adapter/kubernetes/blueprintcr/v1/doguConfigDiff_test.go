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
		domainModel domain.CombinedDoguConfigDiffs
		want        CombinedDoguConfigDiff
	}{
		{
			name:        "should exit early if slices are empty",
			domainModel: domain.CombinedDoguConfigDiffs{},
			want:        CombinedDoguConfigDiff{},
		},
		{
			name: "should convert multiple dogu config diffs",
			domainModel: domain.CombinedDoguConfigDiffs{
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
						NeededAction: domain.ConfigActionSet,
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
						NeededAction: domain.ConfigActionSet,
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
						NeededAction: "set",
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
						NeededAction: "set",
					},
				},
			},
		},
		{
			name: "should convert multiple sensitive dogu config diffs",
			domainModel: domain.CombinedDoguConfigDiffs{
				SensitiveDoguConfigDiff: []domain.SensitiveDoguConfigEntryDiff{
					{
						Key: common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
							DoguName: "ldap",
							Key:      "container_config/memory_limit",
						}},
						Actual: domain.DoguConfigValueState{
							Value:  "512m",
							Exists: true,
						},
						Expected: domain.DoguConfigValueState{
							Value:  "1024m",
							Exists: true,
						},
						NeededAction: domain.ConfigActionSet,
					},
					{
						Key: common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
							DoguName: "ldap",
							Key:      "container_config/swap_limit",
						}},
						Actual: domain.DoguConfigValueState{
							Exists: false,
						},
						Expected: domain.DoguConfigValueState{
							Value:  "512m",
							Exists: true,
						},
						NeededAction: domain.ConfigActionSet,
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
						NeededAction: "set",
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
						NeededAction: "set",
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

func Test_convertToCombinedDoguConfigDiffDomain(t *testing.T) {
	tests := []struct {
		name string
		dto  CombinedDoguConfigDiff
		want domain.CombinedDoguConfigDiffs
	}{
		{
			name: "should exit early if slices are empty",
			dto:  CombinedDoguConfigDiff{},
			want: domain.CombinedDoguConfigDiffs{},
		},
		{
			name: "should convert multiple dogu config diffs",
			dto: CombinedDoguConfigDiff{
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
						NeededAction: "set",
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
						NeededAction: "set",
					},
				},
			},
			want: domain.CombinedDoguConfigDiffs{
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
						NeededAction: domain.ConfigActionSet,
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
						NeededAction: domain.ConfigActionSet,
					},
				},
			},
		},
		{
			name: "should convert multiple sensitive dogu config diffs",
			dto: CombinedDoguConfigDiff{
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
						NeededAction: "set",
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
						NeededAction: "set",
					},
				},
			},
			want: domain.CombinedDoguConfigDiffs{
				SensitiveDoguConfigDiff: []domain.SensitiveDoguConfigEntryDiff{
					{
						Key: common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
							DoguName: "ldap",
							Key:      "container_config/memory_limit",
						}},
						Actual: domain.DoguConfigValueState{
							Value:  "512m",
							Exists: true,
						},
						Expected: domain.DoguConfigValueState{
							Value:  "1024m",
							Exists: true,
						},
						NeededAction: domain.ConfigActionSet,
					},
					{
						Key: common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
							DoguName: "ldap",
							Key:      "container_config/swap_limit",
						}},
						Actual: domain.DoguConfigValueState{
							Exists: false,
						},
						Expected: domain.DoguConfigValueState{
							Value:  "512m",
							Exists: true,
						},
						NeededAction: domain.ConfigActionSet,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, convertToCombinedDoguConfigDiffDomain("ldap", tt.dto), "convertToCombinedDoguConfigDiffDomain(%v, %v)", "ldap", tt.dto)
		})
	}
}
