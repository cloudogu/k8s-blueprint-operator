package serializer

import (
	"testing"

	bpv3 "github.com/cloudogu/k8s-blueprint-lib/v3/api/v3"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/stretchr/testify/assert"
)

var (
	val1 = "val1"
)

func Test_convertToDoguConfigDTO(t *testing.T) {
	tests := []struct {
		name   string
		config domain.DoguConfigEntries
		want   []bpv3.ConfigEntry
	}{
		{
			name:   "nil config",
			config: nil,
			want:   nil,
		},
		{
			name:   "empty config",
			config: domain.DoguConfigEntries{},
			want:   nil,
		},
		{
			name: "convert present config",
			config: domain.DoguConfigEntries{
				{
					Key:   testDoguKey1.Key,
					Value: (*config.Value)(&val1),
				},
			},
			want: []bpv3.ConfigEntry{
				{Key: testDoguKey1.Key.String(), Value: &val1},
			},
		},
		{
			name: "convert absent config",
			config: domain.DoguConfigEntries{
				{
					Key:    testDoguKey1.Key,
					Absent: true,
				},
			},
			want: []bpv3.ConfigEntry{
				{Key: testDoguKey1.Key.String(), Absent: &trueVar},
			},
		},
		{
			name: "censor sensitive config values",
			config: domain.DoguConfigEntries{
				{
					Key:       testDoguKey1.Key,
					Sensitive: true,
					Value:     (*config.Value)(&val1),
				},
			},
			want: []bpv3.ConfigEntry{
				{Key: testDoguKey1.Key.String(), Sensitive: &trueVar},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, convertToDoguConfigDTO(tt.config), "convertToDoguConfigDTO(%v)", tt.config)
		})
	}
}

func Test_convertToDoguConfigEntriesDomain(t *testing.T) {
	tests := []struct {
		name   string
		config []bpv3.ConfigEntry
		want   domain.DoguConfigEntries
	}{
		{
			name:   "nil",
			config: nil,
			want:   nil,
		},
		{
			name:   "empty config",
			config: []bpv3.ConfigEntry{},
			want:   nil,
		},
		{
			name: "convert present config",
			config: []bpv3.ConfigEntry{
				{
					Key:   testDoguKey1.Key.String(),
					Value: &val1,
				},
			},
			want: domain.DoguConfigEntries{
				{
					Key:   testDoguKey1.Key,
					Value: (*config.Value)(&val1),
				},
			},
		},
		{
			name: "convert absent config",
			config: []bpv3.ConfigEntry{
				{
					Key:    testDoguKey1.Key.String(),
					Absent: &trueVar,
				},
			},
			want: domain.DoguConfigEntries{
				{
					Key:    testDoguKey1.Key,
					Absent: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, convertToDoguConfigEntriesDomain(tt.config), "convertToDoguConfigDomain%v)", tt.config)
		})
	}
}

func Test_convertToGlobalConfigDTO(t *testing.T) {
	tests := []struct {
		name   string
		config domain.GlobalConfigEntries
		want   []bpv3.ConfigEntry
	}{
		{
			name:   "nil",
			config: nil,
			want:   nil,
		},
		{
			name:   "empty",
			config: domain.GlobalConfigEntries{},
			want:   nil,
		},
		{
			name: "convert present",
			config: domain.GlobalConfigEntries{
				{
					Key:   "test",
					Value: (*config.Value)(&val1),
				},
			},
			want: []bpv3.ConfigEntry{
				{
					Key:   "test",
					Value: &val1,
				},
			},
		},
		{
			name: "convert absent",
			config: domain.GlobalConfigEntries{
				{
					Key:    "test",
					Absent: true,
				},
			},
			want: []bpv3.ConfigEntry{
				{
					Key:    "test",
					Absent: &trueVar,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, convertToGlobalConfigDTO(tt.config), "convertToGlobalConfigDTO(%v)", tt.config)
		})
	}
}

func Test_convertToGlobalConfigDomain(t *testing.T) {
	tests := []struct {
		name   string
		config []bpv3.ConfigEntry
		want   domain.GlobalConfigEntries
	}{
		{
			name:   "nil",
			config: nil,
			want:   nil,
		},
		{
			name: "convert present",
			config: []bpv3.ConfigEntry{
				{
					Key:   "test",
					Value: &val1,
				},
			},
			want: domain.GlobalConfigEntries{
				{
					Key:   "test",
					Value: (*config.Value)(&val1),
				},
			},
		},
		{
			name: "convert absent",
			config: []bpv3.ConfigEntry{
				{
					Key:    "test",
					Absent: &trueVar,
				},
			},
			want: domain.GlobalConfigEntries{
				{
					Key:    "test",
					Absent: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, convertToGlobalConfigDomain(tt.config), "convertToGlobalConfigDomain(%v)", tt.config)
		})
	}
}
