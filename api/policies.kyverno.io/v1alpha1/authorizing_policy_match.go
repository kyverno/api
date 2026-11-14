package v1alpha1

import (
	"fmt"
	"slices"
	"strings"
)

const serviceAccountUserPrefix = "system:serviceaccount:"

// AuthorizationSubject holds the subject attributes of an authorization request.
type AuthorizationSubject struct {
	User   string
	Groups []string
}

// serviceAccount returns `<namespace>/<name>` when the user is a service account.
func (s AuthorizationSubject) serviceAccount() (string, bool) {
	rest, ok := strings.CutPrefix(s.User, serviceAccountUserPrefix)
	if !ok {
		return "", false
	}
	namespace, name, ok := strings.Cut(rest, ":")
	if !ok || namespace == "" || name == "" || strings.Contains(name, ":") {
		return "", false
	}
	return namespace + "/" + name, true
}

// AuthorizationResourceAttributes holds the resource attributes of an authorization request.
type AuthorizationResourceAttributes struct {
	APIGroup    string
	APIVersion  string
	Resource    string
	Subresource string
	Namespace   string
	Name        string
	Verb        string
}

// AuthorizationNonResourceAttributes holds the non-resource attributes of an authorization request.
type AuthorizationNonResourceAttributes struct {
	Path string
	Verb string
}

// MatchesSubject reports whether the subject is selected by ANY term.
func (s *AuthorizingPolicySpec) MatchesSubject(subject AuthorizationSubject) bool {
	for i := range s.MatchSubjects {
		if s.MatchSubjects[i].Matches(subject) {
			return true
		}
	}
	return false
}

// Matches reports whether ALL requirements of the term match the subject. An empty term matches nothing.
func (t *SubjectSelectorTerm) Matches(subject AuthorizationSubject) bool {
	if len(t.MatchExpressions) == 0 {
		return false
	}
	for i := range t.MatchExpressions {
		if !t.MatchExpressions[i].Matches(subject) {
			return false
		}
	}
	return true
}

func (r *SubjectSelectorRequirement) isWildcard() bool {
	return len(r.Values) == 1 && r.Values[0] == SubjectSelectorWildcard
}

// Matches reports whether the requirement matches the subject. Invalid requirements never match.
func (r *SubjectSelectorRequirement) Matches(subject AuthorizationSubject) bool {
	switch r.Key {
	case SubjectSelectorKeyUser:
		switch r.Operator {
		case SubjectSelectorOpIn:
			return r.isWildcard() || slices.Contains(r.Values, subject.User)
		case SubjectSelectorOpNotIn:
			return len(r.Values) > 0 && !slices.Contains(r.Values, subject.User)
		}
	case SubjectSelectorKeyGroups:
		switch r.Operator {
		case SubjectSelectorOpIn:
			if r.isWildcard() {
				return len(subject.Groups) > 0
			}
			return slices.ContainsFunc(subject.Groups, func(g string) bool { return slices.Contains(r.Values, g) })
		case SubjectSelectorOpNotIn:
			return len(r.Values) > 0 && !slices.ContainsFunc(subject.Groups, func(g string) bool { return slices.Contains(r.Values, g) })
		case SubjectSelectorOpExists:
			return len(subject.Groups) > 0
		case SubjectSelectorOpDoesNotExist:
			return len(subject.Groups) == 0
		}
	case SubjectSelectorKeyServiceAccount:
		sa, isSA := subject.serviceAccount()
		switch r.Operator {
		case SubjectSelectorOpIn:
			return isSA && (r.isWildcard() || slices.Contains(r.Values, sa))
		case SubjectSelectorOpNotIn:
			return len(r.Values) > 0 && (!isSA || !slices.Contains(r.Values, sa))
		case SubjectSelectorOpExists:
			return isSA
		case SubjectSelectorOpDoesNotExist:
			return !isSA
		}
	}
	return false
}

func matchesValue(values []string, value string) bool {
	return slices.Contains(values, "*") || slices.Contains(values, value)
}

// MatchesResource reports whether a resource request is selected by ANY resource rule.
func (m *AuthorizingMatchResources) MatchesResource(attrs AuthorizationResourceAttributes) bool {
	for i := range m.ResourceRules {
		if m.ResourceRules[i].Matches(attrs) {
			return true
		}
	}
	return false
}

// MatchesNonResource reports whether a non-resource request is selected by ANY non-resource rule.
func (m *AuthorizingMatchResources) MatchesNonResource(attrs AuthorizationNonResourceAttributes) bool {
	for i := range m.NonResourceRules {
		if m.NonResourceRules[i].Matches(attrs) {
			return true
		}
	}
	return false
}

// Matches reports whether the resource request matches the rule.
func (r *AuthorizingResourceRule) Matches(attrs AuthorizationResourceAttributes) bool {
	if !matchesValue(r.APIGroups, attrs.APIGroup) || !matchesValue(r.APIVersions, attrs.APIVersion) || !matchesValue(r.Verbs, attrs.Verb) {
		return false
	}
	if !slices.ContainsFunc(r.Resources, func(entry string) bool {
		resource, subresource, _ := strings.Cut(entry, "/")
		return (resource == "*" || resource == attrs.Resource) && (subresource == "*" || subresource == attrs.Subresource)
	}) {
		return false
	}
	if len(r.Namespaces) > 0 {
		if attrs.Namespace == "" || !matchesValue(r.Namespaces, attrs.Namespace) {
			return false
		}
	}
	if len(r.ResourceNames) > 0 {
		if attrs.Name == "" || !slices.Contains(r.ResourceNames, attrs.Name) {
			return false
		}
	}
	return true
}

// Matches reports whether the non-resource request matches the rule.
func (r *AuthorizingNonResourceRule) Matches(attrs AuthorizationNonResourceAttributes) bool {
	if !matchesValue(r.Verbs, attrs.Verb) {
		return false
	}
	return slices.ContainsFunc(r.Paths, func(path string) bool {
		if prefix, ok := strings.CutSuffix(path, "*"); ok {
			return strings.HasPrefix(attrs.Path, prefix)
		}
		return path == attrs.Path
	})
}

// BroadSelectorWarnings returns human-readable warnings for selectors that match (almost) every
// subject or request. Kyverno returns them as admission warnings; they are not validation errors.
func (s *AuthorizingPolicySpec) BroadSelectorWarnings() []string {
	var warnings []string
	for i, term := range s.MatchSubjects {
		path := fmt.Sprintf("spec.matchSubjects[%d]", i)
		negationOnly := len(term.MatchExpressions) > 0
		for _, req := range term.MatchExpressions {
			if req.Key == SubjectSelectorKeyUser && req.Operator == SubjectSelectorOpIn && req.isWildcard() {
				warnings = append(warnings, path+": term matches every subject (user In [\"*\"])")
				negationOnly = false
				break
			}
			if req.Operator != SubjectSelectorOpNotIn && req.Operator != SubjectSelectorOpDoesNotExist {
				negationOnly = false
			}
		}
		if negationOnly {
			warnings = append(warnings, path+": term uses only NotIn/DoesNotExist requirements and matches almost every subject")
		}
	}
	for i, rule := range s.MatchResources.ResourceRules {
		allResources := slices.Contains(rule.Resources, "*") || slices.Contains(rule.Resources, "*/*")
		if allResources && slices.Contains(rule.APIGroups, "*") && slices.Contains(rule.Verbs, "*") {
			warnings = append(warnings, fmt.Sprintf("spec.matchResources.resourceRules[%d]: rule matches every verb on every resource", i))
		}
	}
	for i, rule := range s.MatchResources.NonResourceRules {
		if slices.Contains(rule.Paths, "*") || slices.Contains(rule.Paths, "/*") {
			warnings = append(warnings, fmt.Sprintf("spec.matchResources.nonResourceRules[%d]: rule matches every non-resource path", i))
		}
	}
	return warnings
}
