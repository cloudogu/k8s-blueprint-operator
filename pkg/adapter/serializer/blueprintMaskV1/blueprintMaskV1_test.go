package blueprintMaskV1

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ConvertToBlueprintMaskV1_ok(t *testing.T) {
	dogus := []domain.MaskTargetDogu{
		{Namespace: "absent", Name: "dogu1", Version: "3.2.1-1", TargetState: domain.TargetStateAbsent},
		{Namespace: "absent", Name: "dogu2", TargetState: domain.TargetStateAbsent},
		{Namespace: "present", Name: "dogu3", Version: "3.2.1-2", TargetState: domain.TargetStatePresent},
		{Namespace: "present", Name: "dogu4", Version: "1.2.3-3"},
	}
	blueprint := domain.BlueprintMask{Dogus: dogus}

	maskV1, err := ConvertToBlueprintMaskV1(blueprint)

	convertedDogus := []MaskTargetDogu{
		{Name: "absent/dogu1", Version: "3.2.1-1", TargetState: "absent"},
		{Name: "absent/dogu2", TargetState: "absent"},
		{Name: "present/dogu3", Version: "3.2.1-2", TargetState: "present"},
		{Name: "present/dogu4", Version: "1.2.3-3", TargetState: "present"},
	}

	require.Nil(t, err)
	assert.Equal(t, BlueprintMaskV1{
		GeneralBlueprintMask: serializer.GeneralBlueprintMask{API: serializer.BlueprintMaskAPIV1},
		Dogus:                convertedDogus,
	}, maskV1)
}

func Test_ConvertToBlueprintMaskV1_error(t *testing.T) {
	dogus := []domain.MaskTargetDogu{
		{Namespace: "absent", Name: "dogu1", Version: "3.2.1-1", TargetState: -1},
	}
	blueprint := domain.BlueprintMask{Dogus: dogus}

	_, err := ConvertToBlueprintMaskV1(blueprint)

	require.NotNil(t, err)
	assert.ErrorContains(t, err, "cannot convert blueprintMask to BlueprintMaskV1 DTO: ")
	assert.ErrorContains(t, err, "unknown target state ID: '-1'")
}

func Test_ConvertToBlueprintMask(t *testing.T) {
	dogus := []MaskTargetDogu{
		{Name: "absent/dogu1", Version: "3.2.1-1", TargetState: "absent"},
		{Name: "absent/dogu2", TargetState: "absent"},
		{Name: "present/dogu3", Version: "3.2.1-2", TargetState: "present"},
		{Name: "present/dogu4", Version: "1.2.3-3"},
	}

	blueprintV2 := BlueprintMaskV1{
		GeneralBlueprintMask: serializer.GeneralBlueprintMask{API: serializer.BlueprintMaskAPIV1},
		Dogus:                dogus,
	}
	blueprint, err := convertToBlueprintMask(blueprintV2)

	require.Nil(t, err)

	convertedDogus := []domain.MaskTargetDogu{
		{Namespace: "absent", Name: "dogu1", Version: "3.2.1-1", TargetState: domain.TargetStateAbsent},
		{Namespace: "absent", Name: "dogu2", TargetState: domain.TargetStateAbsent},
		{Namespace: "present", Name: "dogu3", Version: "3.2.1-2", TargetState: domain.TargetStatePresent},
		{Namespace: "present", Name: "dogu4", Version: "1.2.3-3"},
	}

	assert.Equal(t, domain.BlueprintMask{
		Dogus: convertedDogus,
	}, blueprint)
}

func Test_ConvertToBlueprintMask_errors(t *testing.T) {
	maskV1 := BlueprintMaskV1{
		GeneralBlueprintMask: serializer.GeneralBlueprintMask{API: serializer.BlueprintMaskAPIV1},
		Dogus: []MaskTargetDogu{
			{Name: "dogu1", Version: "3.2.1-1", TargetState: "unknown"},
			{Name: "official/dogu1", Version: "3.2.1-1", TargetState: "unknown"},
			{Name: "name/space/dogu2", Version: "3.2.1-2"},
		},
	}

	_, err := convertToBlueprintMask(maskV1)

	require.ErrorContains(t, err, "syntax of blueprintMaskV1 is not correct: ")
	require.ErrorContains(t, err, "dogu name needs to be in the form 'namespace/dogu' but is 'dogu1'")
	require.ErrorContains(t, err, "dogu name needs to be in the form 'namespace/dogu' but is 'name/space/dogu2'")
	require.ErrorContains(t, err, "unknown targetState 'unknown'")
}
