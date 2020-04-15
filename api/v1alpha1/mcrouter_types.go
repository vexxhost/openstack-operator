// Copyright 2020 VEXXHOST, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// McrouterPoolSpec defines the desired state of an Mcrouter pool
type McrouterPoolSpec struct {
	Servers []string `json:"servers"`
}

// McrouterSpec defines the desired state of Mcrouter
type McrouterSpec struct {
	Pools        map[string]McrouterPoolSpec `json:"pools"`
	Route        string                      `json:"route"`
	NodeSelector map[string]string           `json:"nodeSelector,omitempty"`
	Tolerations  []v1.Toleration             `json:"tolerations,omitempty"`
}

// McrouterStatus defines the observed state of Mcrouter
type McrouterStatus struct {
	// +kubebuilder:validation:Default=Pending
	Phase string `json:"phase"`
}

// +kubebuilder:object:root=true

// Mcrouter is the Schema for the mcrouters API
type Mcrouter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   McrouterSpec   `json:"spec,omitempty"`
	Status McrouterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// McrouterList contains a list of Mcrouter
type McrouterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Mcrouter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Mcrouter{}, &McrouterList{})
}
