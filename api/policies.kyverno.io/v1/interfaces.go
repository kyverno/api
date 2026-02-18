package v1

import "github.com/kyverno/api/api/policies.kyverno.io/v1beta1"

type (
	GenericPolicy             = v1beta1.GenericPolicy
	DeletingPolicyLike        = v1beta1.DeletingPolicyLike
	GeneratingPolicyLike      = v1beta1.GeneratingPolicyLike
	ImageValidatingPolicyLike = v1beta1.ImageValidatingPolicyLike
	MutatingPolicyLike        = v1beta1.MutatingPolicyLike
	ValidatingPolicyLike      = v1beta1.ValidatingPolicyLike
)
