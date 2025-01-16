package domain

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bpv2 "github.com/cloudogu/blueprint-lib/v2"
	"github.com/cloudogu/cesapp-lib/core"
)

var version3_2_1_4, _ = core.ParseVersion("3.2.1-4")

var (
	compVersion3210 = semver.MustParse("3.2.1-0")
	compVersion3212 = semver.MustParse("3.2.1-2")
	compVersion3213 = semver.MustParse("3.2.1-3")

	testComponentName1 = bpv2.QualifiedComponentName{Namespace: "k8s", SimpleName: "my-component1"}
	testComponentName2 = bpv2.QualifiedComponentName{Namespace: "official", SimpleName: "my-component2"}
	testComponentName3 = bpv2.QualifiedComponentName{Namespace: "testing", SimpleName: "my-component3"}
	testComponentName4 = bpv2.QualifiedComponentName{Namespace: "k8s", SimpleName: "my-component4"}
)

func Test_validate_ok(t *testing.T) {
	dogus := []bpv2.Dogu{
		{Name: officialDogu1, Version: version3_2_1_0, TargetState: bpv2.TargetStateAbsent},
		{Name: officialDogu2, TargetState: bpv2.TargetStateAbsent},
		{Name: officialDogu3, Version: version3_2_1_0, TargetState: bpv2.TargetStatePresent},
		{Name: officialNexus, Version: version3213},
	}

	components := []bpv2.Component{
		{Name: testComponentName1, Version: compVersion3210, TargetState: bpv2.TargetStateAbsent},
		{Name: testComponentName2, TargetState: bpv2.TargetStateAbsent},
		{Name: testComponentName3, Version: compVersion3212, TargetState: bpv2.TargetStatePresent},
		{Name: testComponentName4, Version: compVersion3213},
	}
	blueprint := bpv2.Blueprint{Dogus: dogus, Components: components}

	err := newBlueprintValidator(blueprint).validate()

	require.NoError(t, err)
}

func Test_validate_multipleErrors(t *testing.T) {
	dogus := []bpv2.Dogu{
		{Version: version3212, TargetState: 666},
	}
	components := []bpv2.Component{
		{Version: compVersion3212},
		{Name: testComponentName, Version: compVersion3212},
	}
	blueprint := bpv2.Blueprint{
		Dogus:      dogus,
		Components: components,
		Config: bpv2.Config{
			Global: bpv2.GlobalConfig{
				Present: nil,
				Absent: []bpv2.GlobalConfigKey{
					"",
				},
			},
		},
	}

	err := newBlueprintValidator(blueprint).validate()

	require.Error(t, err)
	assert.ErrorContains(t, err, "blueprint is invalid")
	assert.ErrorContains(t, err, "dogu is invalid")
	assert.ErrorContains(t, err, "dogu target state is invalid")
	assert.ErrorContains(t, err, "component name must not be empty")
	assert.ErrorContains(t, err, `namespace of component "" must not be empty`)
	assert.ErrorContains(t, err, `key for absent global config should not be empty`)
}

func Test_validateDogus_ok(t *testing.T) {
	dogus := []bpv2.Dogu{
		{Name: officialDogu1, Version: version3_2_1_4, TargetState: bpv2.TargetStateAbsent},
		//versionIsOptionalForStateAbsent
		{Name: officialDogu2, TargetState: bpv2.TargetStateAbsent},
		{Name: officialDogu3, Version: version3212, TargetState: bpv2.TargetStatePresent},
		//StateDefaultsToPresent
		{Name: officialNexus, Version: version3212},
	}
	blueprint := bpv2.Blueprint{Dogus: dogus}

	err := newBlueprintValidator(blueprint).validate()

	require.NoError(t, err)
}

func Test_validateDogus_multipleErrors(t *testing.T) {
	dogus := []bpv2.Dogu{
		{Name: officialDogu1},
		{Name: officialDogu2, TargetState: 666},
	}
	blueprint := bpv2.Blueprint{Dogus: dogus}

	err := newBlueprintValidator(blueprint).validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "dogu target state is invalid")
	assert.Contains(t, err.Error(), "dogu version must not be empty")
}

func Test_validateComponents_ok(t *testing.T) {
	components := []bpv2.Component{
		{
			Name: bpv2.QualifiedComponentName{
				Namespace:  "k8s",
				SimpleName: "absent-component",
			},
			TargetState: bpv2.TargetStateAbsent,
		},
		{
			Name: bpv2.QualifiedComponentName{
				SimpleName: "present-component",
				Namespace:  "k8s",
			},
			Version:     compVersion3212,
			TargetState: bpv2.TargetStatePresent,
		},
	}
	blueprint := bpv2.Blueprint{Components: components}

	err := NewComponentValidator(blueprint).validate()

	require.NoError(t, err)
}

func Test_validateComponents_multipleErrors(t *testing.T) {
	components := []bpv2.Component{
		{Name: testComponentName},
		{Version: compVersion3212},
	}
	blueprint := bpv2.Blueprint{Components: components}
	err := NewComponentValidator(blueprint).validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "component name must not be empty")
	assert.Contains(t, err.Error(), `version of component "k8s/my-component" must not be empty`)
}

func Test_validateDoguUniqueness(t *testing.T) {
	dogus := []bpv2.Dogu{
		{Name: officialDogu1, Version: version3_2_1_0, TargetState: bpv2.TargetStatePresent},
		{Name: officialDogu1, Version: version3213},
		{Name: officialDogu2, Version: version3213},
		{Name: officialDogu2, Version: version3213},
	}

	blueprint := bpv2.Blueprint{Dogus: dogus}

	err := newBlueprintValidator(blueprint).validateDoguUniqueness()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "there are duplicate dogus")
	assert.Contains(t, err.Error(), "dogu1")
	assert.Contains(t, err.Error(), "dogu2")
}

func Test_validateComponentUniqueness(t *testing.T) {
	components := []bpv2.Component{
		{
			Name: bpv2.QualifiedComponentName{
				Namespace:  "present",
				SimpleName: "component1",
			},
			Version:     compVersion3210,
			TargetState: bpv2.TargetStatePresent,
		},
		{
			Name: bpv2.QualifiedComponentName{
				Namespace:  "present",
				SimpleName: "component1",
			},
			Version: compVersion3213},
		{
			Name: bpv2.QualifiedComponentName{
				Namespace:  "present",
				SimpleName: "component2",
			},
			Version: compVersion3213},
		{
			Name: bpv2.QualifiedComponentName{
				Namespace:  "present",
				SimpleName: "component2",
			},
			Version: compVersion3213},
	}

	blueprint := bpv2.Blueprint{Components: components}

	err := newBlueprintValidator(blueprint).validateComponentUniqueness()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "there are duplicate components")
	assert.Contains(t, err.Error(), "component1")
	assert.Contains(t, err.Error(), "component2")
}
