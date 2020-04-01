package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RabbitmqPolicySpec defines the Rabbitmq Policy Spec for the Vhost
type RabbitmqPolicyDefinitionSpec struct {
	Vhost      string                   `json:"vhost,omitempty"`
	Name       string                   `json:"name"`
	Pattern    string                   `json:"pattern"`
	Definition RabbitmqPolicyDefinition `json:"definition"`
	Priority   int64                    `json:"priority"`
	ApplyTo    string                   `json:"apply-to"`
}

// RabbitmqPolicyDefinition defines the Rabbitmq Policy content
type RabbitmqPolicyDefinition struct {
	FederationUpstreamSet string `json:"federation-upstream-set,omitempty"`
	HaMode                string `json:"ha-mode,omitempty"`
	HaParams              int    `json:"ha-params,omitempty"`
	HaSyncMode            string `json:"ha-sync-mode,omitempty"`
	Expires               int    `json:"expires,omitempty"`
	MessageTTL            int    `json:"message-ttl,omitempty"`
	MaxLen                int    `json:"max-length,omitempty"`
	MaxLenBytes           int    `json:"max-length-bytes,omitempty"`
}

// RabbitmqSpec defines the desired state of Rabbitmq
type RabbitmqSpec struct {
	AuthSecret   string                         `json:"authSecret"`
	Policies     []RabbitmqPolicyDefinitionSpec `json:"policies,omitempty"`
	NodeSelector map[string]string              `json:"nodeSelector,omitempty"`
	Tolerations  []v1.Toleration                `json:"tolerations,omitempty"`
}

// RabbitmqStatus defines the observed state of Rabbitmq
type RabbitmqStatus struct {
	// +kubebuilder:validation:Default=Pending
	Phase string `json:"phase"`
}

// Rabbitmq is the Schema for the Rabbitmqs API
// +kubebuilder:object:root=true
type Rabbitmq struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RabbitmqSpec   `json:"spec,omitempty"`
	Status RabbitmqStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RabbitmqList contains a list of Rabbitmq
type RabbitmqList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Rabbitmq `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Rabbitmq{}, &RabbitmqList{})
}
