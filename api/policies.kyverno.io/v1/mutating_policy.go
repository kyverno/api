package v1

import (
	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	MutatingPolicySpec                    = v1beta1.MutatingPolicySpec
	MutatingPolicyStatus                  = v1beta1.MutatingPolicyStatus
	TargetMatchConstraints                = v1beta1.TargetMatchConstraints
	MutatingPolicyEvaluationConfiguration = v1beta1.MutatingPolicyEvaluationConfiguration
	MutatingPolicyAutogenConfiguration    = v1beta1.MutatingPolicyAutogenConfiguration
	MAPGenerationConfiguration            = v1beta1.MAPGenerationConfiguration
	MutateExistingConfiguration           = v1beta1.MutateExistingConfiguration
	MutationTarget                        = v1beta1.MutationTarget
)

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=mutatingpolicies,scope="Cluster",shortName=mpol,categories=kyverno
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="READY",type=string,JSONPath=`.status.conditionStatus.ready`
// +kubebuilder:selectablefield:JSONPath=`.spec.evaluation.mode`
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MutatingPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MutatingPolicySpec `json:"spec"`
	// Status contains policy runtime data.
	// +optional
	Status MutatingPolicyStatus `json:"status,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope="Namespaced",shortName=nmpol,categories=kyverno
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="READY",type=string,JSONPath=`.status.conditionStatus.ready`
// +kubebuilder:selectablefield:JSONPath=`.spec.evaluation.mode`
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NamespacedMutatingPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MutatingPolicySpec `json:"spec"`
	// Status contains policy runtime data.
	// +optional
	Status MutatingPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MutatingPolicyList is a list of MutatingPolicy instances
type MutatingPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []MutatingPolicy `json:"items"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespacedMutatingPolicyList is a list of NamespacedMutatingPolicy instances
type NamespacedMutatingPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []NamespacedMutatingPolicy `json:"items"`
}
