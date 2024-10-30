package ecosystem

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testDoguNamespace    cescommons.DoguNamespace  = "testing"
	postgresqlSimpleName cescommons.SimpleDoguName = "postgresql"
	postfixSimpleName    cescommons.SimpleDoguName = "postfix"
	ldapSimpleName       cescommons.SimpleDoguName = "ldap"
)

var (
	postgresqlQualifiedName = cescommons.QualifiedDoguName{
		Namespace:  testDoguNamespace,
		SimpleName: postgresqlSimpleName,
	}
	postfixQualifiedName = cescommons.QualifiedDoguName{
		Namespace:  testDoguNamespace,
		SimpleName: postfixSimpleName,
	}
	ldapQualifiedName = cescommons.QualifiedDoguName{
		Namespace:  testDoguNamespace,
		SimpleName: ldapSimpleName,
	}
)

func TestCalculateDoguHealthResult(t *testing.T) {
	tests := []struct {
		name  string
		dogus []*DoguInstallation
		want  DoguHealthResult
	}{
		{
			name: "",
			dogus: []*DoguInstallation{
				{
					Name:   postgresqlQualifiedName,
					Health: AvailableHealthStatus,
				},
				{
					Name:   postfixQualifiedName,
					Health: UnavailableHealthStatus,
				},
				{
					Name:   ldapQualifiedName,
					Health: PendingHealthStatus,
				},
			},
			want: DoguHealthResult{
				DogusByStatus: map[HealthStatus][]cescommons.SimpleDoguName{
					AvailableHealthStatus:   {postgresqlSimpleName},
					UnavailableHealthStatus: {postfixSimpleName},
					PendingHealthStatus:     {ldapSimpleName},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, CalculateDoguHealthResult(tt.dogus), "CalculateDoguHealthResult(%v)", tt.dogus)
		})
	}
}

func TestDoguHealthResult_String(t *testing.T) {
	tests := []struct {
		name         string
		healthStates map[HealthStatus][]cescommons.SimpleDoguName
		contains     []string
		notContains  []string
	}{
		{
			name:         "no dogus should result in 0 components unhealthy",
			healthStates: map[HealthStatus][]cescommons.SimpleDoguName{},
			contains:     []string{"0 dogu(s) are unhealthy: "},
		},
		{
			name: "only available dogus should result in 0 components unhealthy",
			healthStates: map[HealthStatus][]cescommons.SimpleDoguName{
				AvailableHealthStatus: {"nginx-ingress"},
			},
			contains:    []string{"0 dogu(s) are unhealthy: "},
			notContains: []string{"nginx-ingress"},
		},
		{
			name: "any dogus not available should be unhealthy",
			healthStates: map[HealthStatus][]cescommons.SimpleDoguName{
				AvailableHealthStatus:   {"nginx-static"},
				UnavailableHealthStatus: {postgresqlSimpleName, "redmine"},
				"other":                 {"scm"},
			},
			contains: []string{
				"3 dogu(s) are unhealthy: ",
				"postgresql",
				"redmine",
				"scm",
			},
			notContains: []string{"nginx-static"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DoguHealthResult{DogusByStatus: tt.healthStates}
			actual := result.String()
			for _, contains := range tt.contains {
				assert.Contains(t, actual, contains)
			}
			for _, notContains := range tt.notContains {
				assert.NotContains(t, actual, notContains)
			}
		})
	}
}

func TestDoguHealthResult_AllHealthy(t *testing.T) {
	tests := []struct {
		name         string
		healthStates map[HealthStatus][]cescommons.SimpleDoguName
		want         bool
	}{
		{
			name:         "should be healthy if empty",
			healthStates: map[HealthStatus][]cescommons.SimpleDoguName{},
			want:         true,
		},
		{
			name: "should be healthy if all are available",
			healthStates: map[HealthStatus][]cescommons.SimpleDoguName{
				AvailableHealthStatus: {"nginx-ingress", "nginx-static", "postfix"},
			},
			want: true,
		},
		{
			name: "should not be healthy if one is not available",
			healthStates: map[HealthStatus][]cescommons.SimpleDoguName{
				AvailableHealthStatus:   {"nginx-ingress", "nginx-static", "postfix"},
				UnavailableHealthStatus: {"postfix"},
			},
			want: false,
		},
		{
			name: "should not be healthy if multiple are not available",
			healthStates: map[HealthStatus][]cescommons.SimpleDoguName{
				AvailableHealthStatus:   {"nginx-ingress", "nginx-static", "postfix"},
				UnavailableHealthStatus: {"postfix", "redmine"},
				PendingHealthStatus:     {postgresqlSimpleName},
				"other":                 {"plantuml"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DoguHealthResult{DogusByStatus: tt.healthStates}
			assert.Equalf(t, tt.want, result.AllHealthy(), "AllHealthy()")
		})
	}
}
