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
