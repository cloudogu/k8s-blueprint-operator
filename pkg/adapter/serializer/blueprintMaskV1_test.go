package serializer

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ConvertToBlueprintMaskV1(t *testing.T) {
	dogus := []domain.MaskTargetDogu{
		{Namespace: "absent", Name: "dogu1", Version: "3.2.1-1", TargetState: TargetStateAbsent},
		{Namespace: "absent", Name: "dogu2", TargetState: TargetStateAbsent},
		{Namespace: "present", Name: "dogu3", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Namespace: "present", Name: "dogu4", Version: "1.2.3-3"},
	}
	blueprint := domain.BlueprintMask{Dogus: dogus}

	maskV1 := ConvertToBlueprintMaskV1(blueprint)

	convertedDogus := []MaskTargetDogu{
		{Name: "absent/dogu1", Version: "3.2.1-1", TargetState: TargetStateAbsent},
		{Name: "absent/dogu2", TargetState: TargetStateAbsent},
		{Name: "present/dogu3", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Name: "present/dogu4", Version: "1.2.3-3"},
	}

	assert.Equal(t, BlueprintMaskV1{
		GeneralBlueprintMask: GeneralBlueprintMask{BlueprintMaskAPIV1},
		Dogus:                convertedDogus,
	}, maskV1)
}

func Test_ConvertToBlueprintMask(t *testing.T) {
	dogus := []MaskTargetDogu{
		{Name: "absent/dogu1", Version: "3.2.1-1", TargetState: TargetStateAbsent},
		{Name: "absent/dogu2", TargetState: TargetStateAbsent},
		{Name: "present/dogu3", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Name: "present/dogu4", Version: "1.2.3-3"},
	}

	blueprintV2 := BlueprintMaskV1{
		GeneralBlueprintMask: GeneralBlueprintMask{BlueprintMaskAPIV1},
		Dogus:                dogus,
	}
	blueprint, err := convertToBlueprintMask(blueprintV2)

	require.Nil(t, err)

	convertedDogus := []domain.MaskTargetDogu{
		{Namespace: "absent", Name: "dogu1", Version: "3.2.1-1", TargetState: TargetStateAbsent},
		{Namespace: "absent", Name: "dogu2", TargetState: TargetStateAbsent},
		{Namespace: "present", Name: "dogu3", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Namespace: "present", Name: "dogu4", Version: "1.2.3-3"},
	}

	assert.Equal(t, domain.BlueprintMask{
		Dogus: convertedDogus,
	}, blueprint)
}

func Test_ConvertToBlueprintMask_invalidDoguName(t *testing.T) {
	blueprintV2 := BlueprintV2{
		GeneralBlueprint: GeneralBlueprint{V2},
		Dogus: []TargetDogu{
			{Name: "dogu1", Version: "3.2.1-1"},
			{Name: "name/space/dogu2", Version: "3.2.1-2"},
		},
	}

	_, err := convertToBlueprint(blueprintV2)

	require.ErrorContains(t, err, "syntax of blueprintV2 is not correct")
	require.ErrorContains(t, err, "dogu name needs to be in the form 'namespace/dogu' but is 'dogu1'")
	require.ErrorContains(t, err, "dogu name needs to be in the form 'namespace/dogu' but is 'name/space/dogu2'")
}
