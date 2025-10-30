package domain

import (
	"testing"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/stretchr/testify/assert"
)

var version3_2_1_4, _ = core.ParseVersion("3.2.1-4")

func Test_validate_ok(t *testing.T) {
	dogus := []Dogu{
		{Name: officialDogu1, Version: &version3_2_1_0, Absent: true},
		{Name: officialDogu2, Absent: true},
		{Name: officialDogu3, Version: &version3_2_1_0, Absent: false},
		{Name: officialNexus, Version: &version3213},
	}

	blueprint := Blueprint{Dogus: dogus}

	err := blueprint.Validate()

	require.NoError(t, err)
}

func Test_validate_multipleErrors(t *testing.T) {

	dogus := []Dogu{
		{Version: nil},
	}
	blueprint := Blueprint{
		Dogus: dogus,
		Config: Config{
			Global: GlobalConfigEntries{
				{
					Key:    "",
					Absent: true,
				},
			},
		},
	}

	err := blueprint.Validate()

	require.Error(t, err)
	assert.ErrorContains(t, err, "blueprint is invalid")
	assert.ErrorContains(t, err, "dogu is invalid")
	assert.ErrorContains(t, err, "dogu version must not be empty")
	assert.ErrorContains(t, err, `key for global config should not be empty`)
}

func Test_validateDogus_ok(t *testing.T) {
	dogus := []Dogu{
		{Name: officialDogu1, Version: &version3_2_1_4, Absent: true},
		//versionIsOptionalForStateAbsent
		{Name: officialDogu2, Absent: true},
		{Name: officialDogu3, Version: &version3212, Absent: false},
		//StateDefaultsToPresent
		{Name: officialNexus, Version: &version3212},
	}
	blueprint := Blueprint{Dogus: dogus}

	err := blueprint.validateDogus()

	require.NoError(t, err)
}

func Test_validateDogus_multipleErrors(t *testing.T) {
	wrongBodySize := resource.MustParse("1Ki")
	dogus := []Dogu{
		{Name: officialDogu1},
		{Name: officialDogu2, ReverseProxyConfig: ecosystem.ReverseProxyConfig{MaxBodySize: &wrongBodySize}},
	}
	blueprint := Blueprint{Dogus: dogus}

	err := blueprint.validateDogus()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "dogu proxy body size is not in Decimal SI (\"M\" or \"G\")")
	assert.Contains(t, err.Error(), "dogu version must not be empty")
}

func Test_validateDoguUniqueness(t *testing.T) {
	dogus := []Dogu{
		{Name: officialDogu1, Version: &version3_2_1_0, Absent: false},
		{Name: officialDogu1, Version: &version3213},
		{Name: officialDogu2, Version: &version3213},
		{Name: officialDogu2, Version: &version3213},
	}

	blueprint := Blueprint{Dogus: dogus}

	err := blueprint.validateDoguUniqueness()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "there are duplicate dogus")
	assert.Contains(t, err.Error(), "dogu1")
	assert.Contains(t, err.Error(), "dogu2")
}
