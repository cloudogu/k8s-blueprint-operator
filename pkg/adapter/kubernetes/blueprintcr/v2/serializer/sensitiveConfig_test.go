package serializer

import (
	"testing"

	v2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/stretchr/testify/assert"
)

func Test_convertToSensitiveDoguConfigDTO(t *testing.T) {

	tests := []struct {
		name   string
		config domain.SensitiveDoguConfig
		want   []v2.ConfigEntry
	}{
		{
			name:   "empty struct to nil",
			config: domain.SensitiveDoguConfig{},
			want:   nil,
		},
		{
			name: "empty config to nil",
			config: domain.SensitiveDoguConfig{
				Present: nil,
				Absent:  nil,
			},
			want: nil,
		},
		{
			name: "convert present config",
			config: domain.SensitiveDoguConfig{
				Present: map[common.DoguConfigKey]domain.SensitiveValueRef{
					testDoguKey1: {
						SecretName: "mySecret",
						SecretKey:  "myKey",
					},
				},
			},
			want: []v2.ConfigEntry{
				{
					Key: string(testDoguKey1.Key),
					SecretRef: &v2.SecretReference{
						Name: "mySecret",
						Key:  "myKey",
					},
					Sensitive: &trueVar,
				},
			},
		},
		{
			name: "convert absent config",
			config: domain.SensitiveDoguConfig{
				Absent: []common.SensitiveDoguConfigKey{
					testDoguKey1,
				},
			},
			want: []v2.ConfigEntry{
				{
					Key:    string(testDoguKey1.Key),
					Absent: &trueVar,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, convertToSensitiveDoguConfigDTO(tt.config), "convertToSensitiveDoguConfigDTO(%v)", tt.config)
		})
	}
}

func Test_convertToSensitiveDoguConfigDomain(t *testing.T) {
	type args struct {
		doguName   string
		doguConfig []v2.ConfigEntry
	}
	tests := []struct {
		name string
		args args
		want domain.SensitiveDoguConfig
	}{
		{
			name: "nil -> empty struct",
			args: args{
				doguName:   string(testDoguKey1.DoguName),
				doguConfig: nil,
			},
			want: domain.SensitiveDoguConfig{},
		},
		{
			name: "empty",
			args: args{
				doguName:   string(testDoguKey1.DoguName),
				doguConfig: []v2.ConfigEntry{},
			},
			want: domain.SensitiveDoguConfig{},
		},
		{
			name: "convert present config",
			args: args{
				doguName: string(testDoguKey1.DoguName),
				doguConfig: []v2.ConfigEntry{
					{
						Key: string(testDoguKey1.Key),
						SecretRef: &v2.SecretReference{
							Name: "mySecret",
							Key:  "myKey",
						},
						Sensitive: &trueVar,
					},
				},
			},
			want: domain.SensitiveDoguConfig{
				Present: map[common.DoguConfigKey]domain.SensitiveValueRef{
					testDoguKey1: {
						SecretName: "mySecret",
						SecretKey:  "myKey",
					},
				},
			},
		},
		{
			name: "convert absent config",
			args: args{
				doguName: string(testDoguKey1.DoguName),
				doguConfig: []v2.ConfigEntry{
					{
						Key:    string(testDoguKey1.Key),
						Absent: &trueVar,
					},
				},
			},
			want: domain.SensitiveDoguConfig{
				Absent: []common.SensitiveDoguConfigKey{
					testDoguKey1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, convertToSensitiveDoguConfigDomain(tt.args.doguName, tt.args.doguConfig), "convertToSensitiveDoguConfigDomain(%v, %v)", tt.args.doguName, tt.args.doguConfig)
		})
	}
}
