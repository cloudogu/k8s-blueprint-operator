package ecosystem

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
					Name:   "postgresql",
					Health: AvailableHealthStatus,
				},
				{
					Name:   "postfix",
					Health: UnavailableHealthStatus,
				},
				{
					Name:   "ldap",
					Health: PendingHealthStatus,
				},
			},
			want: DoguHealthResult{
				DogusByStatus: map[HealthStatus][]DoguName{
					AvailableHealthStatus:   {"postgresql"},
					UnavailableHealthStatus: {"postfix"},
					PendingHealthStatus:     {"ldap"},
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
		healthStates map[HealthStatus][]DoguName
		contains     []string
		notContains  []string
	}{
		{
			name:         "no dogus should result in 0 components unhealthy",
			healthStates: map[HealthStatus][]DoguName{},
			contains:     []string{"0 dogus are unhealthy: "},
		},
		{
			name: "only available dogus should result in 0 components unhealthy",
			healthStates: map[HealthStatus][]DoguName{
				AvailableHealthStatus: {DoguName("nginx-ingress")},
			},
			contains:    []string{"0 dogus are unhealthy: "},
			notContains: []string{"nginx-ingress"},
		},
		{
			name: "any dogus not available should be unhealthy",
			healthStates: map[HealthStatus][]DoguName{
				AvailableHealthStatus:   {DoguName("nginx-static")},
				UnavailableHealthStatus: {DoguName("postgresql"), DoguName("redmine")},
				"other":                 {DoguName("scm")},
			},
			contains: []string{
				"3 dogus are unhealthy: ",
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
