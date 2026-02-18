package v1alpha1

import (
	"time"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type GenericPolicy interface {
	metav1.Object
	GetMatchConstraints() admissionregistrationv1.MatchResources
	GetMatchConditions() []admissionregistrationv1.MatchCondition
	GetFailurePolicy(bool) admissionregistrationv1.FailurePolicyType
	GetTimeoutSeconds() *int32
	GetVariables() []admissionregistrationv1.Variable
}

// DeletingPolicyLike captures the common behavior shared by deleting policies regardless of scope.
type DeletingPolicyLike interface {
	metav1.Object
	runtime.Object
	GetDeletingPolicySpec() *DeletingPolicySpec
	GetKind() string
	GetExecutionTime() (*time.Time, error)
	GetNextExecutionTime(time.Time) (*time.Time, error)
}

// GeneratingPolicyLike captures the common behaviour shared by generating policies regardless of scope.
type GeneratingPolicyLike interface {
	metav1.Object
	runtime.Object
	GetSpec() *GeneratingPolicySpec
	GetStatus() *GeneratingPolicyStatus
	GetMatchConstraints() admissionregistrationv1.MatchResources
	GetMatchConditions() []admissionregistrationv1.MatchCondition
	GetVariables() []admissionregistrationv1.Variable
	GetKind() string
}

// ImageValidatingPolicyLike captures the common behaviour shared by image validating policies regardless of scope.
type ImageValidatingPolicyLike interface {
	metav1.Object
	runtime.Object
	GetSpec() *ImageValidatingPolicySpec
	GetStatus() *ImageValidatingPolicyStatus
	GetFailurePolicy(bool) admissionregistrationv1.FailurePolicyType
	GetMatchConstraints() admissionregistrationv1.MatchResources
	GetMatchConditions() []admissionregistrationv1.MatchCondition
	GetVariables() []admissionregistrationv1.Variable
	GetWebhookConfiguration() *WebhookConfiguration
	BackgroundEnabled() bool
	GetKind() string
}

// MutatingPolicyLike captures the common behaviour shared by mutating policies regardless of scope.
type MutatingPolicyLike interface {
	metav1.Object
	runtime.Object
	GetSpec() *MutatingPolicySpec
	GetStatus() *MutatingPolicyStatus
	GetFailurePolicy(bool) admissionregistrationv1.FailurePolicyType
	GetMatchConstraints() admissionregistrationv1.MatchResources
	GetTargetMatchConstraints() TargetMatchConstraints
	GetMatchConditions() []admissionregistrationv1.MatchCondition
	GetVariables() []admissionregistrationv1.Variable
	GetWebhookConfiguration() *WebhookConfiguration
	BackgroundEnabled() bool
	GetKind() string
}

// ValidatingPolicyLike captures the common behaviour shared by validating policies regardless of scope.
type ValidatingPolicyLike interface {
	metav1.Object
	runtime.Object
	GetSpec() *ValidatingPolicySpec
	GetStatus() *ValidatingPolicyStatus
	GetFailurePolicy(bool) admissionregistrationv1.FailurePolicyType
	GetMatchConstraints() admissionregistrationv1.MatchResources
	GetMatchConditions() []admissionregistrationv1.MatchCondition
	GetVariables() []admissionregistrationv1.Variable
	GetValidatingPolicySpec() *ValidatingPolicySpec
	BackgroundEnabled() bool
	GetKind() string
}
