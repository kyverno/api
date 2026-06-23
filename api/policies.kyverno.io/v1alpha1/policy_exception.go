package v1alpha1

import (
	"time"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// +genclient
// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:deprecatedversion

// PolicyException declares resources to be excluded from specified policies.
type PolicyException struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec declares policy exception behaviors.
	Spec PolicyExceptionSpec `json:"spec"`
}

// IsExpired checks if the policy exception has expired based on the expiresAt field.
func (p *PolicyException) IsExpired() bool {
	if p.Spec.ExpiresAt == nil {
		return false
	}
	return time.Now().After(p.Spec.ExpiresAt.Time)
}

// PolicyExceptionSpec stores policy exception spec
type PolicyExceptionSpec struct {
	// PolicyRefs identifies the policies to which the exception is applied.
	PolicyRefs []PolicyRef `json:"policyRefs"`

	// MatchConditions is a list of CEL expressions that must be met for a resource to be excluded.
	// +optional
	MatchConditions []admissionregistrationv1.MatchCondition `json:"matchConditions,omitempty"`

	// Images specifies container images to be excluded from policy evaluation.
	// These excluded images can be referenced in CEL expressions via `exceptions.allowedImages`.
	// +optional
	Images []string `json:"images,omitempty"`

	// AllowedValues specifies values that can be used in CEL expressions to bypass policy checks.
	// These values can be referenced in CEL expressions via `exceptions.allowedValues`.
	// +optional
	AllowedValues []string `json:"allowedValues,omitempty"`

	// ReportResult indicates whether the policy exception should be reported in the policy report
	// as a skip result or pass result. Defaults to "skip".
	// +optional
	// +kubebuilder:validation:Enum=skip;pass
	// +kubebuilder:default=skip
	ReportResult string `json:"reportResult,omitempty"`

	// ExpiresAt specifies the time when the policy exception expires.
	// Once expired, the exception will no longer be applied to incoming requests.
	// The expected format is RFC3339 date-time (for example "2026-05-01T00:00:00Z").
	// +optional
	ExpiresAt *metav1.Time `json:"expiresAt,omitempty"`

	// Properties is an optional map for additional metadata attached to this exception.
	// For example:
	// - reason: why this exception is needed
	// - ticket: external approval/request identifier
	// - approved-by: comma-separated approver list
	// +optional
	Properties map[string]string `json:"properties,omitempty"`
}

// Validate implements programmatic validation
func (p *PolicyExceptionSpec) Validate(path *field.Path) (errs field.ErrorList) {
	if len(p.PolicyRefs) == 0 {
		errs = append(errs, field.Invalid(path.Child("policyRefs"), p.PolicyRefs, "must specify at least one policy ref"))
	} else {
		for i, policyRef := range p.PolicyRefs {
			errs = append(errs, policyRef.Validate(path.Child("policyRefs").Index(i))...)
		}
	}
	return errs
}

type PolicyRef struct {
	// Name is the name of the policy
	Name string `json:"name"`

	// Kind is the kind of the policy
	Kind string `json:"kind"`
}

func (p *PolicyRef) Validate(path *field.Path) (errs field.ErrorList) {
	if p.Name == "" {
		errs = append(errs, field.Invalid(path.Child("name"), p.Name, "must specify policy name"))
	}
	if p.Kind == "" {
		errs = append(errs, field.Invalid(path.Child("kind"), p.Kind, "must specify policy kind"))
	}
	return errs
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PolicyExceptionList is a list of Policy Exceptions
type PolicyExceptionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []PolicyException `json:"items"`
}
