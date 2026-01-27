package v1

import "github.com/kyverno/api/api/policies.kyverno.io/v1beta1"

const (
	EvaluationModeKubernetes EvaluationMode = "Kubernetes"
	EvaluationModeJSON       EvaluationMode = "JSON"
)

type (
	EvaluationMode          = v1beta1.EvaluationMode
	EvaluationConfiguration = v1beta1.EvaluationConfiguration
	AdmissionConfiguration  = v1beta1.AdmissionConfiguration
	BackgroundConfiguration = v1beta1.BackgroundConfiguration
)
