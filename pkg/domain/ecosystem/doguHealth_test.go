package ecosystem

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnhealthyDogu_String(t *testing.T) {
	type fields struct {
		Namespace string
		Name      string
		Version   core.Version
		Health    HealthStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ok",
			fields: fields{
				Namespace: "official",
				Name:      "postgresql",
				Version:   version1_2_3_1,
				Health:    UnavailableHealthStatus,
			},
			want: "official/postgresql:1.2.3-1 is unavailable",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ud := UnhealthyDogu{
				Namespace: tt.fields.Namespace,
				Name:      tt.fields.Name,
				Version:   tt.fields.Version,
				Health:    tt.fields.Health,
			}
			assert.Equalf(t, tt.want, ud.String(), "String()")
		})
	}
}
