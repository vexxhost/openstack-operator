package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DesignateSpec defines the desired state of Designate
type DesignateSpec struct {
	Credentials string `json:"credentials"`
	CloudName   string `json:"cloudname"`
}

// DesignateStatus defines the observed state of Designate
type DesignateStatus struct {
}

// +kubebuilder:object:root=true

// Designate is the Schema for the Designates API
type Designate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DesignateSpec   `json:"spec,omitempty"`
	Status DesignateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DesignateList contains a list of Designate
type DesignateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Designate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Designate{}, &DesignateList{})
}
