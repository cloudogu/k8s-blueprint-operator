package serializer

import (
	"testing"

	crd "github.com/cloudogu/k8s-blueprint-lib/v3/api/v3"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/stretchr/testify/assert"
)

var (
	limit512  = "512m"
	limit1024 = "1024m"
)

func Test_convertToDoguConfigEntryDiffsDTO(t *testing.T) {

	tests := []struct {
		name        string
		domainModel domain.DoguConfigDiffs
		want        crd.DoguConfigDiff
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
						Value:  &limit512,
						Exists: true,
					},
					Expected: domain.DoguConfigValueState{
						Value:  &limit1024,
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
						Value:  &limit512,
						Exists: true,
					},
					NeededAction: domain.ConfigActionSet,
				},
			},
			want: crd.DoguConfigDiff{
				{
					Key: "container_config/memory_limit",
					Actual: crd.ConfigValueState{
						Value:  &limit512,
						Exists: true,
					},
					Expected: crd.ConfigValueState{
						Value:  &limit1024,
						Exists: true,
					},
					NeededAction: "set",
				},
				{
					Key: "container_config/swap_limit",
					Actual: crd.ConfigValueState{
						Exists: false,
					},
					Expected: crd.ConfigValueState{
						Value:  &limit512,
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
