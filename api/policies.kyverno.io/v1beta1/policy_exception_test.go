package v1beta1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestCELPolicyException_GetKind(t *testing.T) {
	tests := []struct {
		name   string
		policy *PolicyException
		want   string
	}{{
		name:   "not set",
		policy: &PolicyException{},
		want:   "PolicyException",
	}, {
		name: "kind overridden",
		policy: &PolicyException{
			TypeMeta: v1.TypeMeta{
				Kind: "Foo",
			},
		},
		want: "PolicyException",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.policy.GetKind()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCELPolicyExceptionSpec_Validate(t *testing.T) {
	tests := []struct {
		name     string
		policy   *PolicyException
		wantErrs field.ErrorList
	}{{
		name:   "no refs",
		policy: &PolicyException{},
		wantErrs: field.ErrorList{{
			Type:     field.ErrorTypeInvalid,
			Field:    "spec.policyRefs",
			BadValue: []PolicyRef(nil),
			Detail:   "must specify at least one policy ref",
		}},
	}, {
		name: "one ref",
		policy: &PolicyException{
			Spec: PolicyExceptionSpec{
				PolicyRefs: []PolicyRef{{
					Name: "foo",
					Kind: "Foo",
				}},
			},
		},
		wantErrs: nil,
	}, {
		name: "ref no kind",
		policy: &PolicyException{
			Spec: PolicyExceptionSpec{
				PolicyRefs: []PolicyRef{{
					Name: "foo",
				}},
			},
		},
		wantErrs: field.ErrorList{{
			Type:     field.ErrorTypeInvalid,
			Field:    "spec.policyRefs[0].kind",
			BadValue: "",
			Detail:   "must specify policy kind",
		}},
	}, {
		name: "ref no name",
		policy: &PolicyException{
			Spec: PolicyExceptionSpec{
				PolicyRefs: []PolicyRef{{
					Kind: "Foo",
				}},
			},
		},
		wantErrs: field.ErrorList{{
			Type:     field.ErrorTypeInvalid,
			Field:    "spec.policyRefs[0].name",
			BadValue: "",
			Detail:   "must specify policy name",
		}},
	}, {
		name: "ref no kind and name",
		policy: &PolicyException{
			Spec: PolicyExceptionSpec{
				PolicyRefs: []PolicyRef{{}},
			},
		},
		wantErrs: field.ErrorList{{
			Type:     field.ErrorTypeInvalid,
			Field:    "spec.policyRefs[0].name",
			BadValue: "",
			Detail:   "must specify policy name",
		}, {
			Type:     field.ErrorTypeInvalid,
			Field:    "spec.policyRefs[0].kind",
			BadValue: "",
			Detail:   "must specify policy kind",
		}},
	}, {
		name: "multiple refs with indexed errors",
		policy: &PolicyException{
			Spec: PolicyExceptionSpec{
				PolicyRefs: []PolicyRef{{
					Name: "foo",
					Kind: "Foo",
				}, {
					Name: "bar",
				}, {
					Kind: "Baz",
				}},
			},
		},
		wantErrs: field.ErrorList{{
			Type:     field.ErrorTypeInvalid,
			Field:    "spec.policyRefs[1].kind",
			BadValue: "",
			Detail:   "must specify policy kind",
		}, {
			Type:     field.ErrorTypeInvalid,
			Field:    "spec.policyRefs[2].name",
			BadValue: "",
			Detail:   "must specify policy name",
		}},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErrs := tt.policy.Validate()
			assert.Equal(t, tt.wantErrs, gotErrs)
		})
	}
}
func TestCELPolicyException_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		policy   *PolicyException
		expected bool
	}{
		{
			name: "expiresAt not set",
			policy: &PolicyException{
				Spec: PolicyExceptionSpec{},
			},
			expected: false,
		},
		{
			name: "expiresAt in future",
			policy: &PolicyException{
				Spec: PolicyExceptionSpec{
					ExpiresAt: &v1.Time{Time: time.Now().Add(1 * time.Hour)},
				},
			},
			expected: false,
		},
		{
			name: "expiresAt in past",
			policy: &PolicyException{
				Spec: PolicyExceptionSpec{
					ExpiresAt: &v1.Time{Time: time.Now().Add(-1 * time.Hour)},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.policy.IsExpired())
		})
	}
}

func TestCELPolicyExceptionSpec_Properties(t *testing.T) {
	policy := &PolicyException{
		Spec: PolicyExceptionSpec{
			PolicyRefs: []PolicyRef{{
				Name: "foo",
				Kind: "Foo",
			}},
			Properties: map[string]string{
				"reason":      "temporary exception",
				"approved-by": "alice,bob,carol",
				"ticket":      "SNOW-12345",
			},
		},
	}

	assert.Equal(t, field.ErrorList(nil), policy.Validate())
	assert.Equal(t, "alice,bob,carol", policy.Spec.Properties["approved-by"])
}
