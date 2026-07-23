package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EvaluationMode = string

type EvaluationConfiguration struct {
	// Mode is the mode of policy evaluation.
	// Allowed values are "Kubernetes", "HTTP" or "Envoy".
	// Optional. Default value is "Kubernetes".
	// +optional
	Mode EvaluationMode `json:"mode,omitempty"`

	// Admission controls policy evaluation during admission.
	// +optional
	Admission *AdmissionConfiguration `json:"admission,omitempty"`

	// Background  controls policy evaluation during background scan.
	// +optional
	Background *BackgroundConfiguration `json:"background,omitempty"`
}

type AdmissionConfiguration struct {
	// Enabled controls if rules are applied during admission.
	// Optional. Default value is "true".
	// +optional
	// +kubebuilder:default=true
	Enabled *bool `json:"enabled,omitempty"`
}

type BackgroundConfiguration struct {
	// Enabled controls if rules are applied to existing resources during a background scan.
	// Optional. Default value is "true". The value must be set to "false" if the policy rule
	// uses variables that are only available in the admission review request (e.g. user name).
	// +optional
	// +kubebuilder:default=true
	Enabled *bool `json:"enabled,omitempty"`

	// Interval defines how often background scanning is performed for this specific policy.
    // Overrides the global background scan interval set at the controller level.
    // The value must be a valid duration string (e.g., "15m", "1h", "30s").
    // Optional. If not set, the global background scan interval is used.
    // +optional
    Interval *metav1.Duration `json:"interval,omitempty"`
}
