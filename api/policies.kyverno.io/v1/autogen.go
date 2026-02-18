package v1

import "github.com/kyverno/api/api/policies.kyverno.io/v1beta1"

type (
	PodControllersGenerationConfiguration = v1beta1.PodControllersGenerationConfiguration
	Target                                = v1beta1.Target
	ValidatingPolicyAutogenStatus         = v1beta1.ValidatingPolicyAutogenStatus
	ValidatingPolicyAutogen               = v1beta1.ValidatingPolicyAutogen
	ImageValidatingPolicyAutogenStatus    = v1beta1.ImageValidatingPolicyAutogenStatus
	ImageValidatingPolicyAutogen          = v1beta1.ImageValidatingPolicyAutogen
	MutatingPolicyAutogenStatus           = v1beta1.MutatingPolicyAutogenStatus
	MutatingPolicyAutogen                 = v1beta1.MutatingPolicyAutogen
)
