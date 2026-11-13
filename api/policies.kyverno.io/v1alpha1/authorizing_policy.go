package v1alpha1

import (
	policieskyvernoio "github.com/kyverno/api/api/policies.kyverno.io"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=authorizingpolicies,scope="Cluster",shortName=apol,categories=kyverno
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="READY",type=string,JSONPath=`.status.conditionStatus.ready`
// +kubebuilder:validation:XValidation:rule="!(self.metadata.name == 'k8s.io' || self.metadata.name.endsWith('.k8s.io') || self.metadata.name == 'kubernetes.io' || self.metadata.name.endsWith('.kubernetes.io'))",message="metadata.name must not be in the reserved k8s.io or kubernetes.io domains because it prefixes authorization condition IDs"
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AuthorizingPolicy is an alpha, cluster-scoped policy that the Kyverno authorization webhook
// evaluates for SubjectAccessReview requests. Rules from all applicable AuthorizingPolicy
// objects are merged into one unordered set and combined with Deny > NoOpinion > Allow
// precedence. Rule expressions that depend on the request object are automatically turned
// into conditional decisions that Kubernetes enforces during admission.
//
// Creating or updating an AuthorizingPolicy is equivalent to holding cluster-admin
// permissions, because an Allow rule can grant any permission to any subject. Restrict
// write access to this resource to cluster administrators.
type AuthorizingPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AuthorizingPolicySpec `json:"spec"`
	// Status contains policy runtime data.
	// +optional
	Status AuthorizingPolicyStatus `json:"status,omitempty"`
}

// GetKind returns the kind of the policy.
func (p *AuthorizingPolicy) GetKind() string {
	return policieskyvernoio.AuthorizingPolicyKind
}

// GetSpec returns the policy spec.
func (p *AuthorizingPolicy) GetSpec() *AuthorizingPolicySpec {
	return &p.Spec
}

// GetStatus returns the policy status.
func (p *AuthorizingPolicy) GetStatus() *AuthorizingPolicyStatus {
	return &p.Status
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AuthorizingPolicyList is a list of AuthorizingPolicy instances.
type AuthorizingPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []AuthorizingPolicy `json:"items"`
}

// AuthorizingPolicySpec is the specification of the desired behavior of the AuthorizingPolicy.
type AuthorizingPolicySpec struct {
	// MatchSubjects selects the subjects (users, groups, and service accounts) this policy applies to.
	// The policy applies to a request when ANY term matches. Enrollment is explicit: an empty or
	// omitted list is invalid. To select every subject, use a term with the requirement
	// `{key: user, operator: In, values: ["*"]}`.
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=8
	// +required
	MatchSubjects []SubjectSelectorTerm `json:"matchSubjects"`

	// MatchResources selects the resource and non-resource requests this policy applies to.
	// +required
	MatchResources AuthorizingMatchResources `json:"matchResources"`

	// Variables contain definitions of variables that can be used in composition of other expressions.
	// Each variable is defined as a named CEL expression and is available under `variables` in rule
	// expressions. The expression of a variable can refer to other variables defined earlier in the
	// list but not those after.
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=64
	// +optional
	Variables []admissionregistrationv1.Variable `json:"variables,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// Rules is an unordered set of authorization rules. Rule order has no effect: rules from all
	// applicable policies are merged and combined with Deny > NoOpinion > Allow precedence.
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=64
	// +required
	Rules []AuthorizingRule `json:"rules"`
}

// SubjectSelectorTerm matches a subject when ALL of its requirements match.
type SubjectSelectorTerm struct {
	// MatchExpressions is a list of subject requirements that are ANDed.
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=8
	// +required
	MatchExpressions []SubjectSelectorRequirement `json:"matchExpressions"`
}

// SubjectSelectorKey is the subject attribute a requirement applies to.
// +kubebuilder:validation:Enum=user;groups;serviceAccount
type SubjectSelectorKey string

const (
	// SubjectSelectorKeyUser selects on the authenticated user name (SubjectAccessReviewSpec.user).
	SubjectSelectorKeyUser SubjectSelectorKey = "user"
	// SubjectSelectorKeyGroups selects on the set of groups the user belongs to (SubjectAccessReviewSpec.groups).
	SubjectSelectorKeyGroups SubjectSelectorKey = "groups"
	// SubjectSelectorKeyServiceAccount selects on the service account identity derived from a user name
	// of the form `system:serviceaccount:<namespace>:<name>`, written as `<namespace>/<name>`.
	SubjectSelectorKeyServiceAccount SubjectSelectorKey = "serviceAccount"
)

// SubjectSelectorOperator is a set operator for subject requirements.
// +kubebuilder:validation:Enum=In;NotIn;Exists;DoesNotExist
type SubjectSelectorOperator string

const (
	SubjectSelectorOpIn           SubjectSelectorOperator = "In"
	SubjectSelectorOpNotIn        SubjectSelectorOperator = "NotIn"
	SubjectSelectorOpExists       SubjectSelectorOperator = "Exists"
	SubjectSelectorOpDoesNotExist SubjectSelectorOperator = "DoesNotExist"
)

// SubjectSelectorWildcard is the value that matches any value. It is valid only as the single value of an In requirement.
const SubjectSelectorWildcard = "*"

// SubjectSelectorRequirement is a selector requirement over a subject attribute.
//
// Semantics by key:
//   - user: In matches when the user name is in values; NotIn matches when it is not.
//     Exists and DoesNotExist are invalid.
//   - groups: In matches when the user is a member of ANY listed group; NotIn matches when the
//     user is a member of NONE of the listed groups; Exists matches when the user has at least one
//     group; DoesNotExist matches when the user has no groups.
//   - serviceAccount: In matches a service account whose `<namespace>/<name>` is in values;
//     NotIn matches a non-service-account user or a service account not in values; Exists matches
//     any service account; DoesNotExist matches any user that is not a service account.
//
// The wildcard value "*" matches any value and is only valid as the single value of an In requirement.
// Matching is exact and case-sensitive.
//
// +kubebuilder:validation:XValidation:rule="self.operator in ['In', 'NotIn'] ? (has(self.values) && size(self.values) > 0) : (!has(self.values) || size(self.values) == 0)",message="values must be non-empty for In and NotIn, and empty for Exists and DoesNotExist"
// +kubebuilder:validation:XValidation:rule="self.key != 'user' || self.operator in ['In', 'NotIn']",message="the user key supports only the In and NotIn operators"
// +kubebuilder:validation:XValidation:rule="!has(self.values) || !self.values.exists(v, v == '*') || (self.operator == 'In' && size(self.values) == 1)",message="the wildcard value \"*\" is valid only as the single value of an In requirement"
// +kubebuilder:validation:XValidation:rule="self.key != 'serviceAccount' || !has(self.values) || self.values.all(v, v == '*' || v.matches('^[a-z0-9]([-a-z0-9]*[a-z0-9])?/[a-z0-9]([-a-z0-9.]*[a-z0-9])?$'))",message="serviceAccount values must be \"*\" or of the form <namespace>/<name>"
type SubjectSelectorRequirement struct {
	// Key is the subject attribute: user, groups, or serviceAccount.
	// +required
	Key SubjectSelectorKey `json:"key"`

	// Operator is one of In, NotIn, Exists, or DoesNotExist.
	// +required
	Operator SubjectSelectorOperator `json:"operator"`

	// Values is the set of values for In and NotIn. It must be empty for Exists and DoesNotExist.
	// +listType=set
	// +kubebuilder:validation:MaxItems=64
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=253
	// +optional
	Values []string `json:"values,omitempty"`
}

// AuthorizingMatchResources selects the requests a policy applies to. A request matches when it
// matches ANY resource rule (resource requests) or ANY non-resource rule (non-resource requests).
// +kubebuilder:validation:XValidation:rule="(has(self.resourceRules) && size(self.resourceRules) > 0) || (has(self.nonResourceRules) && size(self.nonResourceRules) > 0)",message="at least one of resourceRules or nonResourceRules is required"
type AuthorizingMatchResources struct {
	// ResourceRules selects resource requests.
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=32
	// +optional
	ResourceRules []AuthorizingResourceRule `json:"resourceRules,omitempty"`

	// NonResourceRules selects non-resource requests, such as /healthz or /metrics.
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=32
	// +optional
	NonResourceRules []AuthorizingNonResourceRule `json:"nonResourceRules,omitempty"`
}

// AuthorizingResourceRule selects resource requests. All set fields must match.
type AuthorizingResourceRule struct {
	// APIGroups is the list of API groups. "" is the core API group and "*" matches all groups.
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=32
	// +required
	APIGroups []string `json:"apiGroups"`

	// APIVersions is the list of API versions. "*" matches all versions.
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=32
	// +kubebuilder:validation:items:MinLength=1
	// +required
	APIVersions []string `json:"apiVersions"`

	// Resources is the list of resources, using the same syntax as admission rules:
	// "pods" matches pods but not its subresources, "pods/exec" matches one subresource,
	// "pods/*" matches pods and all of its subresources, "*/scale" matches the scale subresource
	// of every resource, "*" matches every resource but no subresources, and "*/*" matches everything.
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=32
	// +kubebuilder:validation:items:MinLength=1
	// +required
	Resources []string `json:"resources"`

	// Verbs is the list of authorization verbs, such as get, list, watch, create, update, patch,
	// delete, deletecollection, impersonate, bind, or escalate. "*" matches all verbs.
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=32
	// +kubebuilder:validation:items:MinLength=1
	// +required
	Verbs []string `json:"verbs"`

	// Namespaces restricts the rule to requests in the listed namespaces. "*" matches any
	// non-empty namespace. When omitted, the rule matches namespaced, cluster-scoped, and
	// all-namespace requests. When set, cluster-scoped and all-namespace requests never match.
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=64
	// +kubebuilder:validation:items:MinLength=1
	// +optional
	Namespaces []string `json:"namespaces,omitempty"`

	// ResourceNames restricts the rule to requests for the listed object names. When omitted,
	// any name matches. When set, requests without a name (such as create, list, watch, or
	// deletecollection) never match.
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=64
	// +kubebuilder:validation:items:MinLength=1
	// +optional
	ResourceNames []string `json:"resourceNames,omitempty"`
}

// AuthorizingNonResourceRule selects non-resource requests. All fields must match.
type AuthorizingNonResourceRule struct {
	// Paths is the list of non-resource URL paths. A path matches exactly, or a path ending in
	// "*" matches every path with that prefix.
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=32
	// +kubebuilder:validation:items:MinLength=1
	// +required
	Paths []string `json:"paths"`

	// Verbs is the list of lowercase HTTP verbs, such as get, post, put, or delete. "*" matches all verbs.
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	// +kubebuilder:validation:items:MinLength=1
	// +required
	Verbs []string `json:"verbs"`
}

// AuthorizingRuleEffect is the effect of a rule whose expression evaluates to true.
// +kubebuilder:validation:Enum=Allow;Deny;NoOpinion
type AuthorizingRuleEffect string

const (
	AuthorizingRuleEffectAllow     AuthorizingRuleEffect = "Allow"
	AuthorizingRuleEffectDeny      AuthorizingRuleEffect = "Deny"
	AuthorizingRuleEffectNoOpinion AuthorizingRuleEffect = "NoOpinion"
)

// AuthorizingRule is one authorization rule.
//
// The expression has access to `request` (the SubjectAccessReview spec), `variables`, and
// `object` and `oldObject` (the admission objects). When the result depends on `object` or
// `oldObject`, Kyverno returns a conditional decision that Kubernetes evaluates during
// admission; users do not author separate authorization and admission conditions.
//
// When the expression evaluates to an error, Deny and NoOpinion rules take effect and Allow
// rules do not.
type AuthorizingRule struct {
	// Name uniquely identifies the rule within the policy. It is a stable identifier used in
	// decision reasons and condition IDs (`<policy name>/<rule name>`).
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	// +required
	Name string `json:"name"`

	// Effect is the decision the rule contributes when its expression evaluates to true.
	// +required
	Effect AuthorizingRuleEffect `json:"effect"`

	// Expression is a CEL expression that must evaluate to a bool.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=10240
	// +required
	Expression string `json:"expression"`

	// Description is a human-readable explanation, returned with conditional decisions and used in reasons.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Description string `json:"description,omitempty"`
}

// AuthorizingPolicyStatus is the status of an AuthorizingPolicy.
type AuthorizingPolicyStatus struct {
	// +optional
	ConditionStatus ConditionStatus `json:"conditionStatus,omitempty"`
}

// GetConditionStatus returns the shared condition status.
func (status *AuthorizingPolicyStatus) GetConditionStatus() *ConditionStatus {
	return &status.ConditionStatus
}
