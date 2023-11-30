package serializer

import (
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
	"strings"
)

type BlueprintApi string

const (
	V1        BlueprintApi = "v1"
	V2        BlueprintApi = "v2"
	TestEmpty BlueprintApi = "test/empty"
)

// GeneralBlueprint defines the minimum set to parse the blueprint API version string in order to select the right
// blueprint handling strategy. This is necessary in order to accommodate maximal changes in different blueprint API
// versions.
type GeneralBlueprint struct {
	// API is used to distinguish between different versions of the used API and impacts directly the interpretation of
	// this blueprint. Must not be empty.
	//
	// This field MUST NOT be MODIFIED or REMOVED because the API is paramount for distinguishing between different
	// blueprint version implementations.
	API BlueprintApi `json:"blueprintApi"`
}

// BlueprintV2 describes an abstraction of CES components that should be absent or present within one or more CES
// instances. When the same Blueprint is applied to two different CES instances it is required to leave two equal
// instances in terms of the components.
//
// In general additions without changing the version are fine, as long as they don't change semantics. Removal or
// renaming are breaking changes and require a new blueprint API version.
type BlueprintV2 struct {
	GeneralBlueprint
	// Dogus contains a set of exact dogu versions which should be present or absent in the CES instance after which this
	// blueprint was applied. Optional.
	Dogus []TargetDogu `json:"dogus,omitempty"`
	// Packages contains a set of exact package versions which should be present or absent in the CES instance after which
	// this blueprint was applied. The packages must correspond to the used operation system package manager. Optional.
	Components []TargetComponent `json:"components,omitempty"`
	// Used to configure registry globalRegistryEntries on blueprint upgrades
	RegistryConfig RegistryConfig `json:"registryConfig,omitempty"`
	// Used to remove registry globalRegistryEntries on blueprint upgrades
	RegistryConfigAbsent []string `json:"registryConfigAbsent,omitempty"`
	// Used to configure encrypted registry globalRegistryEntries on blueprint upgrades
	RegistryConfigEncrypted RegistryConfig `json:"registryConfigEncrypted,omitempty"`
}

type RegistryConfig map[string]map[string]interface{}

// TargetDogu defines a Dogu, its version, and the installation state in which it is supposed to be after a blueprint
// was applied.
type TargetDogu struct {
	// Name defines the name of the dogu including its namespace, f. i. "official/nginx". Must not be empty.
	Name string `json:"name"`
	// Version defines the version of the dogu that is to be installed. Must not be empty if the targetState is "present";
	// otherwise it is optional and is not going to be interpreted.
	Version string `json:"version"`
	// TargetState defines a state of installation of this dogu. Optional field, but defaults to "TargetStatePresent"
	TargetState TargetState `json:"targetState"`
}

type TargetComponent struct {
	// Name defines the name of the component including its namespace, f. i. "official/nginx". Must not be empty.
	Name string `json:"name"`
	// Version defines the version of the dogu that is to be installed. Must not be empty if the targetState is "present";
	// otherwise it is optional and is not going to be interpreted.
	Version string `json:"version"`
	// TargetState defines a state of installation of this component. Optional field, but defaults to "TargetStatePresent"
	TargetState TargetState `json:"targetState"`
}

func ConvertToBlueprintV2(spec domain.Blueprint) BlueprintV2 {
	return BlueprintV2{
		GeneralBlueprint: GeneralBlueprint{V2},
		Dogus: util.Map(spec.Dogus, func(dogu domain.TargetDogu) TargetDogu {
			return TargetDogu{
				Name:        dogu.GetQualifiedName(),
				Version:     dogu.Version,
				TargetState: TargetState(dogu.TargetState),
			}
		}),
		Components: util.Map(spec.Components, func(component domain.Component) TargetComponent {
			return TargetComponent{
				Name:        component.Name,
				Version:     component.Version,
				TargetState: TargetState(component.TargetState),
			}
		}),
		RegistryConfig:          RegistryConfig(spec.RegistryConfig),
		RegistryConfigAbsent:    spec.RegistryConfigAbsent,
		RegistryConfigEncrypted: RegistryConfig(spec.RegistryConfigEncrypted),
	}
}

func convertToBlueprint(blueprint BlueprintV2) (domain.Blueprint, error) {
	convertedDogus, err := convertDogus(blueprint.Dogus)
	if err != nil {
		return domain.Blueprint{}, fmt.Errorf("syntax of blueprintV2 is not correct: %w", err)
	}
	return domain.Blueprint{
		Dogus: convertedDogus,
		Components: util.Map(blueprint.Components, func(component TargetComponent) domain.Component {
			return domain.Component{
				Name:        component.Name,
				Version:     component.Version,
				TargetState: domain.TargetState(component.TargetState),
			}
		}),
		RegistryConfig:          domain.RegistryConfig(blueprint.RegistryConfig),
		RegistryConfigAbsent:    blueprint.RegistryConfigAbsent,
		RegistryConfigEncrypted: domain.RegistryConfig(blueprint.RegistryConfigEncrypted),
	}, nil
}

func convertDogus(dogus []TargetDogu) ([]domain.TargetDogu, error) {
	var convertedDogus []domain.TargetDogu
	var errorList []error

	for _, dogu := range dogus {
		doguNamespace, doguName, err := splitDoguName(dogu.Name)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}
		convertedDogus = append(convertedDogus, domain.TargetDogu{
			Namespace:   doguNamespace,
			Name:        doguName,
			Version:     dogu.Version,
			TargetState: domain.TargetState(dogu.TargetState),
		})
	}
	err := errors.Join(errorList...)

	return convertedDogus, err
}

func splitDoguName(doguName string) (string, string, error) {
	splitName := strings.Split(doguName, "/")
	if len(splitName) != 2 {
		return "", "", fmt.Errorf("dogu name needs to be in the form 'namespace/dogu' but is '%s'", doguName)
	}
	return splitName[0], splitName[1], nil
}
