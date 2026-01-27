package v1beta1

import "github.com/kyverno/api/api/policies.kyverno.io/v1alpha1"

const (
	EvaluationModeKubernetes EvaluationMode = "Kubernetes"
	EvaluationModeJSON       EvaluationMode = "JSON"
)

type (
	EvaluationMode          = v1alpha1.EvaluationMode
	EvaluationConfiguration = v1alpha1.EvaluationConfiguration
	AdmissionConfiguration  = v1alpha1.AdmissionConfiguration
	BackgroundConfiguration = v1alpha1.BackgroundConfiguration
)
