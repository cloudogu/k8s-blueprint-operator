package serializer

import (
	v2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_convertToDoguConfigDTO(t *testing.T) {
	tests := []struct {
		name   string
		config domain.DoguConfig
		want   *v2.DoguConfig
	}{
		{
			name:   "nil config",
			config: domain.DoguConfig{},
			want:   nil,
		},
		{
			name: "empty config",
			config: domain.DoguConfig{
				Present: map[common.DoguConfigKey]common.DoguConfigValue{},
				Absent:  []common.DoguConfigKey{},
			},
			want: nil,
		},
		{
			name: "convert present config",
			config: domain.DoguConfig{
				Present: map[common.DoguConfigKey]common.DoguConfigValue{
					testDoguKey1: "val1",
				},
			},
			want: &v2.DoguConfig{
				Present: map[string]string{
					testDoguKey1.Key.String(): "val1",
				},
			},
		},
		{
			name: "convert absent config",
			config: domain.DoguConfig{
				Absent: []common.DoguConfigKey{
					testDoguKey1,
				},
			},
			want: &v2.DoguConfig{
				Absent: []string{
					testDoguKey1.Key.String(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, convertToDoguConfigDTO(tt.config), "convertToDoguConfigDTO(%v)", tt.config)
		})
	}
}

func Test_convertToDoguConfigDomain(t *testing.T) {
	type args struct {
		doguName string
		config   *v2.DoguConfig
	}
	tests := []struct {
		name string
		args args
		want domain.DoguConfig
	}{
		{
			name: "nil",
			args: args{
				doguName: string(testDoguKey1.DoguName),
			},
			want: domain.DoguConfig{},
		},
		{
			name: "nil config",
			args: args{
				doguName: string(testDoguKey1.DoguName),
				config:   &v2.DoguConfig{},
			},
			want: domain.DoguConfig{},
		},
		{
			name: "convert present config",
			args: args{
				doguName: string(testDoguKey1.DoguName),
				config: &v2.DoguConfig{
					Present: map[string]string{
						testDoguKey1.Key.String(): "val1",
					},
				},
			},
			want: domain.DoguConfig{
				Present: map[common.DoguConfigKey]common.DoguConfigValue{
					testDoguKey1: "val1",
				},
			},
		},
		{
			name: "convert absent config",
			args: args{
				doguName: string(testDoguKey1.DoguName),
				config: &v2.DoguConfig{
					Absent: []string{
						testDoguKey1.Key.String(),
					},
				},
			},
			want: domain.DoguConfig{
				Absent: []common.DoguConfigKey{
					testDoguKey1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, convertToDoguConfigDomain(tt.args.doguName, tt.args.config), "convertToDoguConfigDomain(%v, %v)", tt.args.doguName, tt.args.config)
		})
	}
}

func Test_convertToGlobalConfigDTO(t *testing.T) {
	tests := []struct {
		name   string
		config domain.GlobalConfig
		want   v2.GlobalConfig
	}{
		{
			name:   "empty",
			config: domain.GlobalConfig{},
			want:   v2.GlobalConfig{},
		},
		{
			name: "convert present",
			config: domain.GlobalConfig{
				Present: map[common.GlobalConfigKey]common.GlobalConfigValue{
					"test": "val1",
				},
			},
			want: v2.GlobalConfig{
				Present: map[string]string{
					"test": "val1",
				},
			},
		},
		{
			name: "convert absent",
			config: domain.GlobalConfig{
				Absent: []common.GlobalConfigKey{
					"test",
				},
			},
			want: v2.GlobalConfig{
				Absent: []string{
					"test",
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
		config v2.GlobalConfig
		want   domain.GlobalConfig
	}{
		{
			name:   "empty",
			config: v2.GlobalConfig{},
			want:   domain.GlobalConfig{},
		},
		{
			name: "convert present",
			config: v2.GlobalConfig{
				Present: map[string]string{
					"test": "val1",
				},
			},
			want: domain.GlobalConfig{
				Present: map[common.GlobalConfigKey]common.GlobalConfigValue{
					"test": "val1",
				},
			},
		},
		{
			name: "convert present",
			config: v2.GlobalConfig{
				Absent: []string{
					"test",
				},
			},
			want: domain.GlobalConfig{
				Absent: []common.GlobalConfigKey{
					"test",
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
