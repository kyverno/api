package v1beta1

import "github.com/kyverno/api/api/policies.kyverno.io/v1alpha1"

type (
	GenericPolicy             = v1alpha1.GenericPolicy
	DeletingPolicyLike        = v1alpha1.DeletingPolicyLike
	GeneratingPolicyLike      = v1alpha1.GeneratingPolicyLike
	ImageValidatingPolicyLike = v1alpha1.ImageValidatingPolicyLike
	MutatingPolicyLike        = v1alpha1.MutatingPolicyLike
	ValidatingPolicyLike      = v1alpha1.ValidatingPolicyLike
)
