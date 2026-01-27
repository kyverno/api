package v1

import (
	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	GeneratingPolicySpec                        = v1beta1.GeneratingPolicySpec
	GeneratingPolicyStatus                      = v1beta1.GeneratingPolicyStatus
	GeneratingPolicyEvaluationConfiguration     = v1beta1.GeneratingPolicyEvaluationConfiguration
	GenerateExistingConfiguration               = v1beta1.GenerateExistingConfiguration
	SynchronizationConfiguration                = v1beta1.SynchronizationConfiguration
	OrphanDownstreamOnPolicyDeleteConfiguration = v1beta1.OrphanDownstreamOnPolicyDeleteConfiguration
	Generation                                  = v1beta1.Generation
)

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=generatingpolicies,scope="Cluster",shortName=gpol,categories=kyverno
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type GeneratingPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              GeneratingPolicySpec `json:"spec"`
	// Status contains policy runtime data.
	// +optional
	Status GeneratingPolicyStatus `json:"status,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope="Namespaced",shortName=ngpol,categories=kyverno
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//
// NamespacedGeneratingPolicy is the namespaced CEL-based generating policy.
type NamespacedGeneratingPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              GeneratingPolicySpec `json:"spec"`
	// Status contains policy runtime data.
	// +optional
	Status GeneratingPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GeneratingPolicyList is a list of GeneratingPolicy instances
type GeneratingPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []GeneratingPolicy `json:"items"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//
// NamespacedGeneratingPolicyList is a list of NamespacedGeneratingPolicy instances
type NamespacedGeneratingPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []NamespacedGeneratingPolicy `json:"items"`
}
