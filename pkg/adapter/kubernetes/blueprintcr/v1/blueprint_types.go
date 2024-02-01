package v1

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BlueprintSpec defines the desired state of Blueprint
type BlueprintSpec struct {
	// Blueprint json with the desired state of the ecosystem.
	Blueprint string `json:"blueprint"`
	// BlueprintMask json can further restrict the desired state from the blueprint.
	BlueprintMask string `json:"blueprintMask,omitempty"`
	// IgnoreDoguHealth lets the user execute the blueprint even if dogus are unhealthy at the moment.
	IgnoreDoguHealth bool `json:"ignoreDoguHealth,omitempty"`
	// AllowDoguNamespaceSwitch lets the user switch the namespace of dogus in the blueprint mask
	// in comparison to the blueprint.
	AllowDoguNamespaceSwitch bool `json:"allowDoguNamespaceSwitch,omitempty"`
	// DryRun lets the user test a blueprint run to check if all attributes of the blueprint are correct and avoid a result with a failure state.
	DryRun bool `json:"dryRun,omitempty"`
}

// BlueprintStatus defines the observed state of Blueprint
type BlueprintStatus struct {
	// Phase represents the processing state of the blueprint
	Phase domain.StatusPhase `json:"phase,omitempty"`
	// EffectiveBlueprint is the blueprint after applying the blueprint mask.
	EffectiveBlueprint EffectiveBlueprint `json:"effectiveBlueprint,omitempty"`
	// StateDiff is the result of comparing the EffectiveBlueprint to the current cluster state.
	// It describes what operations need to be done to achieve the desired state of the blueprint.
	StateDiff StateDiff `json:"stateDiff,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Blueprint is the Schema for the blueprints API
type Blueprint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired state of the Blueprint.
	Spec BlueprintSpec `json:"spec,omitempty"`
	// Status defines the observed state of the Blueprint.
	Status BlueprintStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BlueprintList contains a list of Blueprint
type BlueprintList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Blueprint `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Blueprint{}, &BlueprintList{})
}
