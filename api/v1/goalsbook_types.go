/*
Copyright 2024.

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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GoalsbookSpec defines the desired state of Goalsbook.
type GoalsbookSpec struct {
	RevisionHistoryLimit *int32 `json:"revisionHistoryLimit,omitempty"`

	Replicas *int32 `json:"replicas,omitempty"`

	ContainerName string `json:"containerName,omitempty"`
	ImageName     string `json:"imageName,omitempty"`
	ContainerPort int32  `json:"containerPort"`
}

// GoalsbookStatus defines the observed state of Goalsbook.
type GoalsbookStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Goalsbook is the Schema for the goalsbooks API.
type Goalsbook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GoalsbookSpec   `json:"spec,omitempty"`
	Status GoalsbookStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GoalsbookList contains a list of Goalsbook.
type GoalsbookList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Goalsbook `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Goalsbook{}, &GoalsbookList{})
}
