package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func req(key SubjectSelectorKey, op SubjectSelectorOperator, values ...string) SubjectSelectorRequirement {
	return SubjectSelectorRequirement{Key: key, Operator: op, Values: values}
}

func TestAuthorizingPolicy_GetKind(t *testing.T) {
	p := &AuthorizingPolicy{}
	assert.Equal(t, "AuthorizingPolicy", p.GetKind())
	assert.Same(t, &p.Spec, p.GetSpec())
	assert.Same(t, &p.Status.ConditionStatus, p.GetStatus().GetConditionStatus())
}

func TestSubjectSelectorRequirement_Matches(t *testing.T) {
	alice := AuthorizationSubject{User: "alice", Groups: []string{"dev", "system:authenticated"}}
	admin := AuthorizationSubject{User: "root", Groups: []string{"system:masters", "dev"}}
	nogroups := AuthorizationSubject{User: "bob"}
	sa := AuthorizationSubject{User: "system:serviceaccount:team-a:builder", Groups: []string{"system:serviceaccounts", "system:serviceaccounts:team-a"}}
	badSA := AuthorizationSubject{User: "system:serviceaccount:team-a"}

	tests := []struct {
		name    string
		req     SubjectSelectorRequirement
		subject AuthorizationSubject
		want    bool
	}{
		{"user In", req(SubjectSelectorKeyUser, SubjectSelectorOpIn, "alice", "carol"), alice, true},
		{"user In miss", req(SubjectSelectorKeyUser, SubjectSelectorOpIn, "carol"), alice, false},
		{"user In wildcard", req(SubjectSelectorKeyUser, SubjectSelectorOpIn, "*"), nogroups, true},
		{"user In wildcard anonymous", req(SubjectSelectorKeyUser, SubjectSelectorOpIn, "*"), AuthorizationSubject{User: "system:anonymous"}, true},
		{"user NotIn", req(SubjectSelectorKeyUser, SubjectSelectorOpNotIn, "carol"), alice, true},
		{"user NotIn miss", req(SubjectSelectorKeyUser, SubjectSelectorOpNotIn, "alice"), alice, false},
		{"user Exists invalid", req(SubjectSelectorKeyUser, SubjectSelectorOpExists), alice, false},
		{"user case-sensitive", req(SubjectSelectorKeyUser, SubjectSelectorOpIn, "Alice"), alice, false},
		{"groups In any", req(SubjectSelectorKeyGroups, SubjectSelectorOpIn, "ops", "dev"), alice, true},
		{"groups In none", req(SubjectSelectorKeyGroups, SubjectSelectorOpIn, "ops"), alice, false},
		{"groups In wildcard", req(SubjectSelectorKeyGroups, SubjectSelectorOpIn, "*"), alice, true},
		{"groups In wildcard no groups", req(SubjectSelectorKeyGroups, SubjectSelectorOpIn, "*"), nogroups, false},
		{"groups NotIn non member", req(SubjectSelectorKeyGroups, SubjectSelectorOpNotIn, "system:masters"), alice, true},
		// multi-group negation: a member of ANY listed group does not match, even with other groups
		{"groups NotIn member of one", req(SubjectSelectorKeyGroups, SubjectSelectorOpNotIn, "system:masters"), admin, false},
		{"groups NotIn multiple values", req(SubjectSelectorKeyGroups, SubjectSelectorOpNotIn, "ops", "dev"), admin, false},
		{"groups NotIn no groups", req(SubjectSelectorKeyGroups, SubjectSelectorOpNotIn, "system:masters"), nogroups, true},
		{"groups NotIn empty values invalid", req(SubjectSelectorKeyGroups, SubjectSelectorOpNotIn), alice, false},
		{"groups Exists", req(SubjectSelectorKeyGroups, SubjectSelectorOpExists), alice, true},
		{"groups Exists none", req(SubjectSelectorKeyGroups, SubjectSelectorOpExists), nogroups, false},
		{"groups DoesNotExist", req(SubjectSelectorKeyGroups, SubjectSelectorOpDoesNotExist), nogroups, true},
		{"groups DoesNotExist has groups", req(SubjectSelectorKeyGroups, SubjectSelectorOpDoesNotExist), alice, false},
		{"sa In", req(SubjectSelectorKeyServiceAccount, SubjectSelectorOpIn, "team-a/builder"), sa, true},
		{"sa In other", req(SubjectSelectorKeyServiceAccount, SubjectSelectorOpIn, "team-b/builder"), sa, false},
		{"sa In wildcard", req(SubjectSelectorKeyServiceAccount, SubjectSelectorOpIn, "*"), sa, true},
		{"sa In wildcard user", req(SubjectSelectorKeyServiceAccount, SubjectSelectorOpIn, "*"), alice, false},
		{"sa In malformed user", req(SubjectSelectorKeyServiceAccount, SubjectSelectorOpIn, "*"), badSA, false},
		{"sa NotIn user", req(SubjectSelectorKeyServiceAccount, SubjectSelectorOpNotIn, "team-a/builder"), alice, true},
		{"sa NotIn listed", req(SubjectSelectorKeyServiceAccount, SubjectSelectorOpNotIn, "team-a/builder"), sa, false},
		{"sa NotIn unlisted", req(SubjectSelectorKeyServiceAccount, SubjectSelectorOpNotIn, "team-b/x"), sa, true},
		{"sa Exists", req(SubjectSelectorKeyServiceAccount, SubjectSelectorOpExists), sa, true},
		{"sa Exists user", req(SubjectSelectorKeyServiceAccount, SubjectSelectorOpExists), alice, false},
		{"sa DoesNotExist", req(SubjectSelectorKeyServiceAccount, SubjectSelectorOpDoesNotExist), alice, true},
		{"sa DoesNotExist sa", req(SubjectSelectorKeyServiceAccount, SubjectSelectorOpDoesNotExist), sa, false},
		{"unknown key", req("uid", SubjectSelectorOpIn, "*"), alice, false},
		{"unknown operator", req(SubjectSelectorKeyUser, "Equals", "alice"), alice, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.req.Matches(tt.subject))
		})
	}
}

func TestAuthorizingPolicySpec_MatchesSubject(t *testing.T) {
	spec := AuthorizingPolicySpec{
		MatchSubjects: []SubjectSelectorTerm{
			// developers who are not cluster admins (AND)
			{MatchExpressions: []SubjectSelectorRequirement{
				req(SubjectSelectorKeyGroups, SubjectSelectorOpIn, "dev"),
				req(SubjectSelectorKeyGroups, SubjectSelectorOpNotIn, "system:masters"),
			}},
			// OR the builder service account
			{MatchExpressions: []SubjectSelectorRequirement{
				req(SubjectSelectorKeyServiceAccount, SubjectSelectorOpIn, "team-a/builder"),
			}},
		},
	}
	assert.True(t, spec.MatchesSubject(AuthorizationSubject{User: "alice", Groups: []string{"dev"}}))
	assert.False(t, spec.MatchesSubject(AuthorizationSubject{User: "root", Groups: []string{"dev", "system:masters"}}))
	assert.True(t, spec.MatchesSubject(AuthorizationSubject{User: "system:serviceaccount:team-a:builder"}))
	assert.False(t, spec.MatchesSubject(AuthorizationSubject{User: "bob"}))
	assert.False(t, (&AuthorizingPolicySpec{}).MatchesSubject(AuthorizationSubject{User: "bob"}), "no terms selects nobody")
	assert.False(t, (&SubjectSelectorTerm{}).Matches(AuthorizationSubject{User: "bob"}), "empty term selects nobody")
}

func TestAuthorizingResourceRule_Matches(t *testing.T) {
	pods := AuthorizationResourceAttributes{APIGroup: "", APIVersion: "v1", Resource: "pods", Namespace: "prod", Name: "web", Verb: "update"}
	withSub := func(a AuthorizationResourceAttributes, sub string) AuthorizationResourceAttributes {
		a.Subresource = sub
		return a
	}
	rule := func(resources ...string) AuthorizingResourceRule {
		return AuthorizingResourceRule{APIGroups: []string{""}, APIVersions: []string{"*"}, Resources: resources, Verbs: []string{"*"}}
	}
	tests := []struct {
		name  string
		rule  AuthorizingResourceRule
		attrs AuthorizationResourceAttributes
		want  bool
	}{
		{"exact resource", rule("pods"), pods, true},
		{"resource excludes subresource", rule("pods"), withSub(pods, "exec"), false},
		{"exact subresource", rule("pods/exec"), withSub(pods, "exec"), true},
		{"subresource excludes parent", rule("pods/exec"), pods, false},
		{"resource wildcard sub includes parent", rule("pods/*"), pods, true},
		{"resource wildcard sub includes sub", rule("pods/*"), withSub(pods, "status"), true},
		{"any resource scale", rule("*/scale"), withSub(AuthorizationResourceAttributes{APIVersion: "v1", Resource: "deployments", Verb: "update"}, "scale"), true},
		{"any resource scale not parent", rule("*/scale"), pods, false},
		{"star excludes subresources", rule("*"), withSub(pods, "exec"), false},
		{"star star", rule("*/*"), withSub(pods, "exec"), true},
		{"group mismatch", AuthorizingResourceRule{APIGroups: []string{"apps"}, APIVersions: []string{"*"}, Resources: []string{"*"}, Verbs: []string{"*"}}, pods, false},
		{"version mismatch", AuthorizingResourceRule{APIGroups: []string{""}, APIVersions: []string{"v2"}, Resources: []string{"*"}, Verbs: []string{"*"}}, pods, false},
		{"verb mismatch", AuthorizingResourceRule{APIGroups: []string{""}, APIVersions: []string{"*"}, Resources: []string{"*"}, Verbs: []string{"get"}}, pods, false},
		{"namespace listed", AuthorizingResourceRule{APIGroups: []string{"*"}, APIVersions: []string{"*"}, Resources: []string{"*"}, Verbs: []string{"*"}, Namespaces: []string{"prod"}}, pods, true},
		{"namespace not listed", AuthorizingResourceRule{APIGroups: []string{"*"}, APIVersions: []string{"*"}, Resources: []string{"*"}, Verbs: []string{"*"}, Namespaces: []string{"dev"}}, pods, false},
		{"namespace wildcard", AuthorizingResourceRule{APIGroups: []string{"*"}, APIVersions: []string{"*"}, Resources: []string{"*"}, Verbs: []string{"*"}, Namespaces: []string{"*"}}, pods, true},
		{"namespace wildcard excludes cluster scope", AuthorizingResourceRule{APIGroups: []string{"*"}, APIVersions: []string{"*"}, Resources: []string{"*"}, Verbs: []string{"*"}, Namespaces: []string{"*"}}, AuthorizationResourceAttributes{APIVersion: "v1", Resource: "nodes", Verb: "get"}, false},
		{"namespace omitted matches cluster scope", rule("*"), AuthorizationResourceAttributes{APIVersion: "v1", Resource: "namespaces", Verb: "list"}, true},
		{"resource name listed", AuthorizingResourceRule{APIGroups: []string{""}, APIVersions: []string{"*"}, Resources: []string{"pods"}, Verbs: []string{"*"}, ResourceNames: []string{"web"}}, pods, true},
		{"resource name not listed", AuthorizingResourceRule{APIGroups: []string{""}, APIVersions: []string{"*"}, Resources: []string{"pods"}, Verbs: []string{"*"}, ResourceNames: []string{"db"}}, pods, false},
		{"resource names exclude nameless", AuthorizingResourceRule{APIGroups: []string{""}, APIVersions: []string{"*"}, Resources: []string{"pods"}, Verbs: []string{"*"}, ResourceNames: []string{"web"}}, AuthorizationResourceAttributes{APIVersion: "v1", Resource: "pods", Namespace: "prod", Verb: "create"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.rule.Matches(tt.attrs))
			m := AuthorizingMatchResources{ResourceRules: []AuthorizingResourceRule{tt.rule}}
			assert.Equal(t, tt.want, m.MatchesResource(tt.attrs))
		})
	}
}

func TestAuthorizingNonResourceRule_Matches(t *testing.T) {
	r := AuthorizingNonResourceRule{Paths: []string{"/healthz", "/debug/*"}, Verbs: []string{"get"}}
	assert.True(t, r.Matches(AuthorizationNonResourceAttributes{Path: "/healthz", Verb: "get"}))
	assert.False(t, r.Matches(AuthorizationNonResourceAttributes{Path: "/healthz/ready", Verb: "get"}))
	assert.True(t, r.Matches(AuthorizationNonResourceAttributes{Path: "/debug/pprof", Verb: "get"}))
	assert.False(t, r.Matches(AuthorizationNonResourceAttributes{Path: "/debug/pprof", Verb: "post"}))
	m := AuthorizingMatchResources{NonResourceRules: []AuthorizingNonResourceRule{r}}
	assert.True(t, m.MatchesNonResource(AuthorizationNonResourceAttributes{Path: "/healthz", Verb: "get"}))
	assert.False(t, m.MatchesResource(AuthorizationResourceAttributes{Resource: "pods", Verb: "get"}))
}

func TestAuthorizingPolicySpec_BroadSelectorWarnings(t *testing.T) {
	narrow := AuthorizingPolicySpec{
		MatchSubjects: []SubjectSelectorTerm{{MatchExpressions: []SubjectSelectorRequirement{
			req(SubjectSelectorKeyGroups, SubjectSelectorOpIn, "dev"),
			req(SubjectSelectorKeyGroups, SubjectSelectorOpNotIn, "system:masters"),
		}}},
		MatchResources: AuthorizingMatchResources{
			ResourceRules:    []AuthorizingResourceRule{{APIGroups: []string{""}, APIVersions: []string{"v1"}, Resources: []string{"pods"}, Verbs: []string{"*"}}},
			NonResourceRules: []AuthorizingNonResourceRule{{Paths: []string{"/metrics"}, Verbs: []string{"get"}}},
		},
	}
	assert.Empty(t, narrow.BroadSelectorWarnings())

	broad := AuthorizingPolicySpec{
		MatchSubjects: []SubjectSelectorTerm{
			{MatchExpressions: []SubjectSelectorRequirement{req(SubjectSelectorKeyUser, SubjectSelectorOpIn, "*")}},
			{MatchExpressions: []SubjectSelectorRequirement{req(SubjectSelectorKeyGroups, SubjectSelectorOpNotIn, "system:masters")}},
			{MatchExpressions: []SubjectSelectorRequirement{req(SubjectSelectorKeyServiceAccount, SubjectSelectorOpDoesNotExist)}},
		},
		MatchResources: AuthorizingMatchResources{
			ResourceRules:    []AuthorizingResourceRule{{APIGroups: []string{"*"}, APIVersions: []string{"*"}, Resources: []string{"*/*"}, Verbs: []string{"*"}}},
			NonResourceRules: []AuthorizingNonResourceRule{{Paths: []string{"/*"}, Verbs: []string{"get"}}},
		},
	}
	assert.Equal(t, []string{
		`spec.matchSubjects[0]: term matches every subject (user In ["*"])`,
		"spec.matchSubjects[1]: term uses only NotIn/DoesNotExist requirements and matches almost every subject",
		"spec.matchSubjects[2]: term uses only NotIn/DoesNotExist requirements and matches almost every subject",
		"spec.matchResources.resourceRules[0]: rule matches every verb on every resource",
		"spec.matchResources.nonResourceRules[0]: rule matches every non-resource path",
	}, broad.BroadSelectorWarnings())
}
