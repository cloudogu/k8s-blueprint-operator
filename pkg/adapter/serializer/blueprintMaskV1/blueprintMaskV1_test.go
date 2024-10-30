package blueprintMaskV1

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var version3211, _ = core.ParseVersion("3.2.1-1")
var version3212, _ = core.ParseVersion("3.2.1-2")
var version1_2_3_3, _ = core.ParseVersion("1.2.3-3")

func Test_ConvertToBlueprintMaskV1_ok(t *testing.T) {
	dogus := []domain.MaskDogu{
		{Name: common.QualifiedDoguName{Namespace: "absent", SimpleName: "dogu1"}, Version: version3211, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedDoguName{Namespace: "absent", SimpleName: "dogu2"}, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedDoguName{Namespace: "present", SimpleName: "dogu3"}, Version: version3212, TargetState: domain.TargetStatePresent},
		{Name: common.QualifiedDoguName{Namespace: "present", SimpleName: "dogu4"}, Version: version1_2_3_3},
	}
	blueprint := domain.BlueprintMask{Dogus: dogus}

	maskV1, err := ConvertToBlueprintMaskV1(blueprint)

	convertedDogus := []MaskTargetDogu{
		{Name: "absent/dogu1", Version: version3211.Raw, TargetState: "absent"},
		{Name: "absent/dogu2", TargetState: "absent"},
		{Name: "present/dogu3", Version: version3212.Raw, TargetState: "present"},
		{Name: "present/dogu4", Version: version1_2_3_3.Raw, TargetState: "present"},
	}

	require.Nil(t, err)
	assert.Equal(t, BlueprintMaskV1{
		GeneralBlueprintMask: serializer.GeneralBlueprintMask{API: serializer.BlueprintMaskAPIV1},
		Dogus:                convertedDogus,
	}, maskV1)
}

func Test_ConvertToBlueprintMaskV1_error(t *testing.T) {
	dogus := []domain.MaskDogu{
		{Name: common.QualifiedDoguName{Namespace: "absent", SimpleName: "dogu1"}, Version: version3211, TargetState: -1},
	}
	blueprint := domain.BlueprintMask{Dogus: dogus}

	_, err := ConvertToBlueprintMaskV1(blueprint)

	require.NotNil(t, err)
	assert.ErrorContains(t, err, "cannot convert blueprintMask to BlueprintMaskV1 DTO: ")
	assert.ErrorContains(t, err, "unknown target state ID: '-1'")
}

func Test_ConvertToBlueprintMask2(t *testing.T) {
	var version32112, err = core.ParseVersion("3.2.1-12")
	require.NoError(t, err)
	require.Equal(t, "3.2.1-12", version32112.Raw)
}

func Test_ConvertToBlueprintMask(t *testing.T) {
	dogus := []MaskTargetDogu{
		{Name: "absent/dogu1", Version: version3211.Raw, TargetState: "absent"},
		{Name: "absent/dogu2", TargetState: "absent"},
		{Name: "present/dogu3", Version: version3212.Raw, TargetState: "present"},
		{Name: "present/dogu4", Version: version1_2_3_3.Raw},
	}

	blueprintV2 := BlueprintMaskV1{
		GeneralBlueprintMask: serializer.GeneralBlueprintMask{API: serializer.BlueprintMaskAPIV1},
		Dogus:                dogus,
	}
	blueprint, err := convertToBlueprintMask(blueprintV2)

	require.Nil(t, err)

	convertedDogus := []domain.MaskDogu{
		{Name: common.QualifiedDoguName{Namespace: "absent", SimpleName: "dogu1"}, Version: version3211, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedDoguName{Namespace: "absent", SimpleName: "dogu2"}, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedDoguName{Namespace: "present", SimpleName: "dogu3"}, Version: version3212, TargetState: domain.TargetStatePresent},
		{Name: common.QualifiedDoguName{Namespace: "present", SimpleName: "dogu4"}, Version: version1_2_3_3},
	}

	assert.Equal(t, domain.BlueprintMask{
		Dogus: convertedDogus,
	}, blueprint)
}

func Test_ConvertToBlueprintMask_errors(t *testing.T) {
	maskV1 := BlueprintMaskV1{
		GeneralBlueprintMask: serializer.GeneralBlueprintMask{API: serializer.BlueprintMaskAPIV1},
		Dogus: []MaskTargetDogu{
			{Name: "dogu1", Version: version3211.Raw, TargetState: "unknown"},
			{Name: "official/dogu1", Version: version3211.Raw, TargetState: "unknown"},
			{Name: "name/space/dogu2", Version: version3212.Raw},
			{Name: "official/dogu3", Version: "abc"},
		},
	}

	_, err := convertToBlueprintMask(maskV1)

	require.ErrorContains(t, err, "syntax of blueprintMaskV1 is not correct: ")
	require.ErrorContains(t, err, "dogu name needs to be in the form 'namespace/dogu' but is 'dogu1'")
	require.ErrorContains(t, err, "dogu name needs to be in the form 'namespace/dogu' but is 'name/space/dogu2'")
	require.ErrorContains(t, err, "unknown target state \"unknown\"")
	require.ErrorContains(t, err, "could not parse version of MaskTargetDogu: failed to parse major version abc")
}
