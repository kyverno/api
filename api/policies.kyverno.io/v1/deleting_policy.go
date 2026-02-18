package v1

import (
	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	DeletingPolicySpec   = v1beta1.DeletingPolicySpec
	DeletingPolicyStatus = v1beta1.DeletingPolicyStatus
)

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=deletingpolicies,scope="Cluster",shortName=dpol,categories=kyverno
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="READY",type=string,JSONPath=`.status.conditionStatus.ready`
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DeletingPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DeletingPolicySpec `json:"spec"`
	// Status contains policy runtime data.
	// +optional
	Status DeletingPolicyStatus `json:"status,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope="Namespaced",shortName=ndpol,categories=kyverno
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="READY",type=string,JSONPath=`.status.conditionStatus.ready`
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NamespacedDeletingPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DeletingPolicySpec `json:"spec"`
	// Status contains policy runtime data.
	// +optional
	Status DeletingPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DeletingPolicyList is a list of DeletingPolicy instances
type DeletingPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []DeletingPolicy `json:"items"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespacedDeletingPolicyList is a list of NamespacedDeletingPolicy instances
type NamespacedDeletingPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []NamespacedDeletingPolicy `json:"items"`
}
