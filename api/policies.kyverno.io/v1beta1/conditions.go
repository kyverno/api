package v1beta1

import (
	"github.com/kyverno/api/api/policies.kyverno.io/v1alpha1"
)

const (
	PolicyConditionTypeWebhookConfigured      PolicyConditionType = "WebhookConfigured"
	PolicyConditionTypePolicyCached           PolicyConditionType = "PolicyCached"
	PolicyConditionTypeRBACPermissionsGranted PolicyConditionType = "RBACPermissionsGranted"
)

type (
	PolicyConditionType = v1alpha1.PolicyConditionType
	ConditionStatus     = v1alpha1.ConditionStatus
)
