package serializer

import (
	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_convertToDoguConfigEntryDiffsDTO(t *testing.T) {
	tests := []struct {
		name        string
		domainModel domain.DoguConfigDiffs
		want        []crd.DoguConfigEntryDiff
		isSensitive bool
	}{
		{
			name:        "should exit early if slices are empty",
			domainModel: domain.DoguConfigDiffs{},
			want:        nil,
		},
		{
			name: "should convert multiple dogu config diffs",
			domainModel: domain.DoguConfigDiffs{
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
			want: []crd.DoguConfigEntryDiff{
				{
					Key: "container_config/memory_limit",
					Actual: crd.DoguConfigValueState{
						Value:  "512m",
						Exists: true,
					},
					Expected: crd.DoguConfigValueState{
						Value:  "1024m",
						Exists: true,
					},
					NeededAction: "set",
				},
				{
					Key: "container_config/swap_limit",
					Actual: crd.DoguConfigValueState{
						Exists: false,
					},
					Expected: crd.DoguConfigValueState{
						Value:  "512m",
						Exists: true,
					},
					NeededAction: "set",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, convertToDoguConfigEntryDiffsDTO(tt.domainModel, tt.isSensitive), "convertToDoguConfigEntryDiffsDTO(%v)", tt.domainModel)
		})
	}
}

func Test_convertToDoguConfigDiffsDomain(t *testing.T) {
	tests := []struct {
		name string
		dto  crd.DoguConfigDiff
		want domain.DoguConfigDiffs
	}{
		{
			name: "should exit early if slices are empty",
			dto:  crd.DoguConfigDiff{},
			want: nil,
		},
		{
			name: "should convert multiple dogu config diffs",
			dto: crd.DoguConfigDiff{
				{
					Key: "container_config/memory_limit",
					Actual: crd.DoguConfigValueState{
						Value:  "512m",
						Exists: true,
					},
					Expected: crd.DoguConfigValueState{
						Value:  "1024m",
						Exists: true,
					},
					NeededAction: "set",
				},
				{
					Key: "container_config/swap_limit",
					Actual: crd.DoguConfigValueState{
						Exists: false,
					},
					Expected: crd.DoguConfigValueState{
						Value:  "512m",
						Exists: true,
					},
					NeededAction: "set",
				},
			},
			want: domain.DoguConfigDiffs{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, convertToDoguConfigDiffsDomain("ldap", tt.dto), "convertToDoguConfigDiffsDomain(%v, %v)", "ldap", tt.dto)
		})
	}
}
