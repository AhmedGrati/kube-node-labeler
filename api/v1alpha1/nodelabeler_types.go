/*
Copyright 2022.

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

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NodeLabelerSpec defines the desired state of NodeLabeler
type NodeLabelerSpec struct {
	v1.NodeSelector `json:",inline"`

	// // +optional
	DryRun bool `json:"dryRun,omitempty"`

	// // +optional
	NodeNamePatterns []string `json:"nodeNamePatterns,omitempty" protobuf:"bytes,3,rep,name=nodeNamePatterns"`

	Merge MergeSpec `json:"merge,omitempty"`

	Overwrite OverwriteSpec `json:"overwrite,omitempty"`
}

// NodeLabelerStatus defines the observed state of NodeLabeler
type NodeLabelerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

type OverwriteSpec struct {
	metav1.ObjectMeta `json:",omitempty"`

	v1.NodeSpec `json:",inline" protobuf:"bytes,2,opt,name=spec"`
}

type MergeSpec struct {
	metav1.ObjectMeta `json:",omitempty"`

	v1.NodeSpec `json:",inline" protobuf:"bytes,2,opt,name=spec"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// NodeLabeler is the Schema for the nodelabelers API
type NodeLabeler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeLabelerSpec   `json:"spec,omitempty"`
	Status NodeLabelerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NodeLabelerList contains a list of NodeLabeler
type NodeLabelerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeLabeler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodeLabeler{}, &NodeLabelerList{})
}
