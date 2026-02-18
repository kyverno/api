package v1

import (
	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	ValidatingPolicySpec                 = v1beta1.ValidatingPolicySpec
	ValidatingPolicyStatus               = v1beta1.ValidatingPolicyStatus
	WebhookConfiguration                 = v1beta1.WebhookConfiguration
	ValidatingPolicyAutogenConfiguration = v1beta1.ValidatingPolicyAutogenConfiguration
	VapGenerationConfiguration           = v1beta1.VapGenerationConfiguration
)

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope="Namespaced",shortName=nvpol,categories=kyverno
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="READY",type=string,JSONPath=`.status.conditionStatus.ready`
// +kubebuilder:selectablefield:JSONPath=`.spec.evaluation.mode`
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NamespacedValidatingPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ValidatingPolicySpec `json:"spec"`
	// Status contains policy runtime data.
	// +optional
	Status ValidatingPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespacedValidatingPolicyList is a list of NamespacedValidatingPolicy instances
type NamespacedValidatingPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []NamespacedValidatingPolicy `json:"items"`
}

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=validatingpolicies,scope="Cluster",shortName=vpol,categories=kyverno
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="READY",type=string,JSONPath=`.status.conditionStatus.ready`
// +kubebuilder:selectablefield:JSONPath=`.spec.evaluation.mode`
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ValidatingPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ValidatingPolicySpec `json:"spec"`
	// Status contains policy runtime data.
	// +optional
	Status ValidatingPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ValidatingPolicyList is a list of ValidatingPolicy instances
type ValidatingPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ValidatingPolicy `json:"items"`
}
