package serializer

import (
	"testing"

	v2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
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
		want   []v2.ConfigEntry
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
			want: []v2.ConfigEntry{
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
			want: []v2.ConfigEntry{
				{Key: testDoguKey1.Key.String(), Absent: &trueVar},
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
		config []v2.ConfigEntry
		want   domain.DoguConfigEntries
	}{
		{
			name:   "nil",
			config: nil,
			want:   nil,
		},
		{
			name:   "empty config",
			config: []v2.ConfigEntry{},
			want:   nil,
		},
		{
			name: "convert present config",
			config: []v2.ConfigEntry{
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
			config: []v2.ConfigEntry{
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
		want   []v2.ConfigEntry
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
			want: []v2.ConfigEntry{
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
			want: []v2.ConfigEntry{
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
		config []v2.ConfigEntry
		want   domain.GlobalConfigEntries
	}{
		{
			name:   "nil",
			config: nil,
			want:   nil,
		},
		{
			name: "convert present",
			config: []v2.ConfigEntry{
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
			config: []v2.ConfigEntry{
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
