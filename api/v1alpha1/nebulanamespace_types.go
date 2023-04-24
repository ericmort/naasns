package v1alpha1

import (
	"encoding/json"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type NaasNamespaceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

func (nns NaasNamespaceSpec) ToJson() string {
	res, _ := json.Marshal(nns)
	return string(res)
}

type NaasNamespaceStatus struct {
}

type NaasNamespace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NaasNamespaceSpec   `json:"spec,omitempty"`
	Status NaasNamespaceStatus `json:"status,omitempty"`
}

type NaasNamespaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NaasNamespace `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NaasNamespace{}, &NaasNamespaceList{})
}
