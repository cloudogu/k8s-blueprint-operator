package ecosystem

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEcosystemState_GetInstalledDoguNames(t *testing.T) {
	nginx := common.SimpleDoguName("nginx")
	officialNginx := common.QualifiedDoguName{
		Namespace:  "official",
		SimpleName: nginx,
	}
	state := EcosystemState{
		InstalledDogus: map[common.SimpleDoguName]*DoguInstallation{
			nginx: {Name: officialNginx},
		},
	}

	names := state.GetInstalledDoguNames()

	assert.Equal(t, []common.SimpleDoguName{nginx}, names)
}
