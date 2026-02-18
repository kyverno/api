package v1

import "github.com/kyverno/api/api/policies.kyverno.io/v1beta1"

const (
	PolicyConditionTypeWebhookConfigured      PolicyConditionType = "WebhookConfigured"
	PolicyConditionTypePolicyCached           PolicyConditionType = "PolicyCached"
	PolicyConditionTypeRBACPermissionsGranted PolicyConditionType = "RBACPermissionsGranted"
)

type (
	PolicyConditionType = v1beta1.PolicyConditionType
	ConditionStatus     = v1beta1.ConditionStatus
)
