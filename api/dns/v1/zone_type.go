package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ZoneSpec defines the desired state of Zone
type ZoneSpec struct {
	Domain      string `json:"domain"`
	TTL         int    `json:"ttl"`
	Email       string `json:"email"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
}

// ZoneStatus defines the observed state of Zone
type ZoneStatus struct {
	// +kubebuilder:validation:Default=Pending
	ZoneID string `json:"zoneId"`
}

// +kubebuilder:object:root=true

// Zone is the Schema for the Zones API
type Zone struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ZoneSpec   `json:"spec,omitempty"`
	Status ZoneStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ZoneList contains a list of Zone
type ZoneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Zone `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Zone{}, &ZoneList{})
}
