/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BlueprintSpec defines the desired state of Blueprint
type BlueprintSpec struct {
	// Blueprint json with the desired state of the ecosystem.
	Blueprint string `json:"blueprint"`
	// BlueprintMask json can further restrict the desired state from the blueprint.
	BlueprintMask string `json:"blueprintMask,omitempty"`
}

// BlueprintStatus defines the observed state of Blueprint
type BlueprintStatus struct {
	// Phase represents the processing state of the blueprint
	Phase StatusPhase `json:"phase,omitempty"`
	// RequeueTimeNanos contains the time in nanoseconds to wait until the next requeue.
	RequeueTimeNanos time.Duration `json:"requeueTimeNanos,omitempty"`
}

type StatusPhase string

const (
	// StatusPhaseNew marks a newly created blueprint-CR.
	StatusPhaseNew StatusPhase = ""
	// StatusPhaseCompleted marks the blueprint as successfully applied.
	StatusPhaseCompleted StatusPhase = "completed"
	// StatusPhaseInvalid marks the given blueprint or the blueprint mask as not correct.
	StatusPhaseInvalid StatusPhase = "invalid"
	// StatusPhaseRetrying marks the blueprint as not applicable for now (e.g. dogu health state) but a retry is queued.
	StatusPhaseRetrying StatusPhase = "retrying"
	// StatusPhaseFailed marks that an error occurred during processing of the blueprint.
	StatusPhaseFailed StatusPhase = "failed"
	// StatusPhaseInProgress marks that the blueprint is currently being processed.
	StatusPhaseInProgress StatusPhase = "inProgress"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Blueprint is the Schema for the blueprints API
type Blueprint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired state of the Blueprint.
	Spec BlueprintSpec `json:"spec,omitempty"`
	// Status defines the observed state of the Blueprint.
	Status BlueprintStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BlueprintList contains a list of Blueprint
type BlueprintList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Blueprint `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Blueprint{}, &BlueprintList{})
}
