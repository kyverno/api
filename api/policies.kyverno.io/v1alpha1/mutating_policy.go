package v1alpha1

import (
	policieskyvernoio "github.com/kyverno/api/api/policies.kyverno.io"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	admissionregistrationv1alpha1 "k8s.io/api/admissionregistration/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=mutatingpolicies,scope="Cluster",shortName=mpol,categories=kyverno
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="READY",type=string,JSONPath=`.status.conditionStatus.ready`
// +kubebuilder:selectablefield:JSONPath=`.spec.evaluation.mode`
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:deprecatedversion

type MutatingPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MutatingPolicySpec `json:"spec"`
	// Status contains policy runtime data.
	// +optional
	Status MutatingPolicyStatus `json:"status,omitempty"`
}

type MutatingPolicyStatus struct {
	// +optional
	ConditionStatus ConditionStatus `json:"conditionStatus,omitempty"`

	// +optional
	Autogen MutatingPolicyAutogenStatus `json:"autogen,omitempty"`

	// Generated indicates whether a MutatingAdmissionPolicy is generated from the policy or not
	// +optional
	Generated bool `json:"generated"`
}

func (status *MutatingPolicyStatus) GetConditionStatus() *ConditionStatus {
	return &status.ConditionStatus
}

// MutatingPolicySpec is the specification of the desired behavior of the MutatingPolicy.
type MutatingPolicySpec struct {
	// MatchConstraints specifies the trigger resources this policy is designed to evaluate.
	// The AdmissionPolicy cares about a request if it matches _all_ Constraints.
	// Trigger constraints and MatchConditions are evaluated before target resolution.
	// Required.
	MatchConstraints *admissionregistrationv1.MatchResources `json:"matchConstraints,omitempty"`

	// failurePolicy defines how to handle failures for the admission policy. Failures can
	// occur from CEL expression parse errors, type check errors, runtime errors and invalid
	// or mis-configured policy definitions or bindings.
	//
	// failurePolicy does not define how validations that evaluate to false are handled.
	//
	// When failurePolicy is set to Fail, the validationActions field define how failures are enforced.
	//
	// Allowed values are Ignore or Fail. Defaults to Fail.
	// +optional
	// +kubebuilder:validation:Enum=Ignore;Fail
	FailurePolicy *admissionregistrationv1.FailurePolicyType `json:"failurePolicy,omitempty"`

	// MatchConditions is a list of conditions that must be met for a request to be validated.
	// Match conditions filter requests that have already been matched by the rules,
	// namespaceSelector, and objectSelector. An empty list of matchConditions matches all requests.
	// There are a maximum of 64 match conditions allowed.
	//
	// The exact matching logic is (in order):
	//   1. If ANY matchCondition evaluates to FALSE, the policy is skipped.
	//   2. If ALL matchConditions evaluate to TRUE, the policy is evaluated.
	//   3. If any matchCondition evaluates to an error (but none are FALSE):
	//      - If failurePolicy=Fail, reject the request
	//      - If failurePolicy=Ignore, the policy is skipped
	//
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	// +optional
	MatchConditions []admissionregistrationv1.MatchCondition `json:"matchConditions,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// Variables contain definitions of variables that can be used in composition of other expressions.
	// Each variable is defined as a named CEL expression.
	// The variables defined here will be available under `variables` in other expressions of the policy,
	// including MatchConditions where they are evaluated lazily on first reference.
	// Note that a native Kubernetes MutatingAdmissionPolicy does not support variables in match
	// conditions; policies that generate a MutatingAdmissionPolicy through autogen should not
	// reference variables in MatchConditions.
	//
	// The expression of a variable can refer to other variables defined earlier in the list but not those after.
	// Thus, Variables must be sorted by the order of first appearance and acyclic.
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	// +optional
	Variables []admissionregistrationv1.Variable `json:"variables,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// AutogenConfiguration defines the configuration for the generation controller.
	// +optional
	AutogenConfiguration *MutatingPolicyAutogenConfiguration `json:"autogen,omitempty"`

	// TargetMatchConstraints resolves the resources to mutate after Variables are evaluated.
	// +optional
	TargetMatchConstraints *TargetMatchConstraints `json:"targetMatchConstraints,omitempty"`

	// TargetMatchConditions is a list of conditions that must be met for a resolved target resource.
	// Target match conditions are evaluated after variables and target resolution. They can reference
	// variables and use Object to refer to the target resource.
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	// +optional
	TargetMatchConditions []admissionregistrationv1.MatchCondition `json:"targetMatchConditions,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// mutations contain operations to perform on matching objects.
	// mutations may not be empty; a minimum of one mutation is required.
	// mutations are evaluated in order, and are reinvoked according to
	// the reinvocationPolicy.
	// The mutations of a policy are invoked for each binding of this policy
	// and reinvocation of mutations occurs on a per binding basis.
	//
	// +listType=atomic
	// +optional
	Mutations []admissionregistrationv1alpha1.Mutation `json:"mutations,omitempty" protobuf:"bytes,4,rep,name=mutations"`

	// AuditAnnotations contains CEL expressions which are used to produce audit annotations that are surfaced
	// as properties in policy report results. auditAnnotations are evaluated after the mutations have been
	// applied successfully. The results of evaluating the expressions are attached to the report result as
	// properties with the annotation key.
	// If the expression evaluates to an empty string or null the annotation will not be included.
	// +optional
	AuditAnnotations []admissionregistrationv1.AuditAnnotation `json:"auditAnnotations,omitempty"`

	// WebhookConfiguration defines the configuration for the webhook.
	// +optional
	WebhookConfiguration *WebhookConfiguration `json:"webhookConfiguration,omitempty"`

	// EvaluationConfiguration defines the configuration for mutating policy evaluation.
	// +optional
	EvaluationConfiguration *MutatingPolicyEvaluationConfiguration `json:"evaluation,omitempty"`

	// reinvocationPolicy indicates whether mutations may be called multiple times per MutatingAdmissionPolicyBinding
	// as part of a single admission evaluation.
	// Allowed values are "Never" and "IfNeeded".
	//
	// Never: These mutations will not be called more than once per binding in a single admission evaluation.
	//
	// IfNeeded: These mutations may be invoked more than once per binding for a single admission request and there is no guarantee of
	// order with respect to other admission plugins, admission webhooks, bindings of this policy and admission policies.  Mutations are only
	// reinvoked when mutations change the object after this mutation is invoked.
	// Required.
	ReinvocationPolicy admissionregistrationv1.ReinvocationPolicyType `json:"reinvocationPolicy,omitempty" protobuf:"bytes,7,opt,name=reinvocationPolicy,casttype=ReinvocationPolicyType"`
}

// GenerateMutatingAdmissionPolicyEnabled checks if mutating admission policy generation is enabled
func (s MutatingPolicySpec) GenerateMutatingAdmissionPolicyEnabled() bool {
	const defaultValue = false
	if s.AutogenConfiguration == nil {
		return defaultValue
	}
	if s.AutogenConfiguration.MutatingAdmissionPolicy == nil {
		return defaultValue
	}
	if s.AutogenConfiguration.MutatingAdmissionPolicy.Enabled == nil {
		return defaultValue
	}
	return *s.AutogenConfiguration.MutatingAdmissionPolicy.Enabled
}

// AdmissionEnabled checks if admission is set to true
func (s MutatingPolicySpec) AdmissionEnabled() bool {
	const defaultValue = true
	if s.EvaluationConfiguration == nil || s.EvaluationConfiguration.Admission == nil || s.EvaluationConfiguration.Admission.Enabled == nil {
		return defaultValue
	}
	return *s.EvaluationConfiguration.Admission.Enabled
}

// BackgroundEnabled checks if background is set to true
func (s MutatingPolicySpec) BackgroundEnabled() bool {
	const defaultValue = true
	if s.EvaluationConfiguration == nil || s.EvaluationConfiguration.Background == nil || s.EvaluationConfiguration.Background.Enabled == nil {
		return defaultValue
	}
	return *s.EvaluationConfiguration.Background.Enabled
}

// SkipBackgroundRequestsEnabled checks if background-controller requests should be skipped.
func (s MutatingPolicySpec) SkipBackgroundRequestsEnabled() bool {
	const defaultValue = true
	if s.EvaluationConfiguration == nil || s.EvaluationConfiguration.SkipBackgroundRequests == nil {
		return defaultValue
	}
	return *s.EvaluationConfiguration.SkipBackgroundRequests
}

// EvaluationMode returns the evaluation mode of the policy.
func (s MutatingPolicySpec) EvaluationMode() EvaluationMode {
	const defaultValue = policieskyvernoio.EvaluationModeKubernetes
	if s.EvaluationConfiguration == nil || s.EvaluationConfiguration.Mode == "" {
		return defaultValue
	}
	return s.EvaluationConfiguration.Mode
}

// GetReinvocationPolicy returns the reinvocation policy of the MutatingPolicy
func (s *MutatingPolicySpec) GetReinvocationPolicy() admissionregistrationv1.ReinvocationPolicyType {
	const defaultValue = admissionregistrationv1.NeverReinvocationPolicy
	if s.ReinvocationPolicy == "" {
		return defaultValue
	}
	return s.ReinvocationPolicy
}

// MutateExistingEnabled checks if mutate existing is set to true
func (s MutatingPolicySpec) MutateExistingEnabled() bool {
	if s.EvaluationConfiguration == nil ||
		s.EvaluationConfiguration.MutateExistingConfiguration == nil ||
		s.EvaluationConfiguration.MutateExistingConfiguration.Enabled == nil {
		return false
	}
	return *s.EvaluationConfiguration.MutateExistingConfiguration.Enabled
}

func (s *MutatingPolicySpec) SetMatchConstraints(in admissionregistrationv1.MatchResources) {
	out := &admissionregistrationv1.MatchResources{}
	out.NamespaceSelector = in.NamespaceSelector
	out.ObjectSelector = in.ObjectSelector
	for _, ex := range in.ExcludeResourceRules {
		out.ExcludeResourceRules = append(out.ExcludeResourceRules, admissionregistrationv1.NamedRuleWithOperations{
			ResourceNames:      ex.ResourceNames,
			RuleWithOperations: ex.RuleWithOperations,
		})
	}
	for _, ex := range in.ResourceRules {
		out.ResourceRules = append(out.ResourceRules, admissionregistrationv1.NamedRuleWithOperations{
			ResourceNames:      ex.ResourceNames,
			RuleWithOperations: ex.RuleWithOperations,
		})
	}
	if in.MatchPolicy != nil {
		out.MatchPolicy = in.MatchPolicy
	}
	s.MatchConstraints = out
}

type MutatingPolicyEvaluationConfiguration struct {
	// Mode is the mode of policy evaluation.
	// Allowed values are "Kubernetes" or "JSON".
	// Optional. Default value is "Kubernetes".
	// +optional
	Mode EvaluationMode `json:"mode,omitempty"`

	// Admission controls policy evaluation during admission.
	// +optional
	Admission *AdmissionConfiguration `json:"admission,omitempty"`

	// Background controls policy evaluation during background scan.
	// +optional
	Background *BackgroundConfiguration `json:"background,omitempty"`

	// MutateExisting controls whether existing resources are mutated.
	// +optional
	MutateExistingConfiguration *MutateExistingConfiguration `json:"mutateExisting,omitempty"`

	// SkipBackgroundRequests bypasses admission requests that are sent by the background controller.
	// The default value is set to "true", it must be set to "false" to apply mutateExisting rules to those requests.
	// +optional
	// +kubebuilder:default=true
	SkipBackgroundRequests *bool `json:"skipBackgroundRequests,omitempty"`

	// UseServerSideApply applies ApplyConfiguration patches with Server-Side Apply semantics,
	// which allows setting atomic fields (for example a container's args or a projected volume)
	// that the default MutatingAdmissionPolicy behaviour rejects. When true, an atomic value is
	// replaced as a whole, so any field the object owner set but the patch does not is dropped.
	// The default is false, which keeps parity with a native MutatingAdmissionPolicy.
	// +optional
	UseServerSideApply bool `json:"useServerSideApply,omitempty"`
}

type MutatingPolicyAutogenConfiguration struct {
	// PodControllers specifies whether to generate a pod controllers rules.
	PodControllers *PodControllersGenerationConfiguration `json:"podControllers,omitempty"`
	// MutatingAdmissionPolicy specifies whether to generate a Kubernetes MutatingAdmissionPolicy.
	MutatingAdmissionPolicy *MAPGenerationConfiguration `json:"mutatingAdmissionPolicy,omitempty"`
}

type MAPGenerationConfiguration struct {
	// Enabled specifies whether to generate a Kubernetes MutatingAdmissionPolicy.
	// Optional. Defaults to "false" if not specified.
	Enabled *bool `json:"enabled,omitempty"`
}

type TargetMatchConstraints struct {
	// Expression resolves one target object or a list of target objects. The expression can use
	// variables computed from the trigger request.
	// +optional
	Expression string `json:"expression,omitempty"`

	// MatchResources constrains target resources and supplies their API route. When Expression is
	// set, these rules still determine whether a target uses a subresource endpoint.
	// +optional
	admissionregistrationv1.MatchResources `json:",inline"`
}

type MutateExistingConfiguration struct {
	// Enabled enables mutation of existing resources. Default is false.
	// When spec.targetMatchConstraints is not defined, Kyverno mutates existing resources matched in spec.matchConstraints.
	// +optional
	// +kubebuilder:default=false
	Enabled *bool `json:"enabled,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MutatingPolicyList is a list of MutatingPolicy instances
type MutatingPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []MutatingPolicy `json:"items"`
}

// MutationTarget specifies the target of the mutation.
type MutationTarget struct {
	// Group specifies the API group of the target resource.
	// +optional
	Group string `json:"group,omitempty"`

	// Version specifies the API version of the target resource.
	// +optional
	Version string `json:"version,omitempty"`

	// Resource specifies the resource name of the target resource.
	// +optional
	Resource string `json:"resource,omitempty"`

	// Kind specifies the kind of the target resource.
	// +optional
	Kind string `json:"kind,omitempty"`
}
