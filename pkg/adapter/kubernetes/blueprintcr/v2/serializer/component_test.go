package serializer

import (
	"fmt"
	bpv2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
)

var (
	k8sK8sDoguOperator = common.QualifiedComponentName{Namespace: "k8s", SimpleName: "k8s-dogu-operator"}
)

func TestConvertComponents(t *testing.T) {
	type args struct {
		components []bpv2.Component
	}
	tests := []struct {
		name    string
		args    args
		want    []domain.Component
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "nil",
			args:    args{components: nil},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name:    "empty list",
			args:    args{components: []bpv2.Component{}},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name:    "normal component",
			args:    args{components: []bpv2.Component{{Name: "k8s/k8s-dogu-operator", Version: version3211.Raw, Absent: false, DeployConfig: map[string]interface{}{"deployNamespace": "longhorn-system", "configOverwrite": map[string]string{"key": "value"}}}}},
			want:    []domain.Component{{Name: k8sK8sDoguOperator, Version: compVersion3211, TargetState: domain.TargetStatePresent, DeployConfig: map[string]interface{}{"deployNamespace": "longhorn-system", "configOverwrite": map[string]string{"key": "value"}}}},
			wantErr: assert.NoError,
		},
		{
			name:    "absent component",
			args:    args{components: []bpv2.Component{{Name: "k8s/k8s-dogu-operator", Absent: true}}},
			want:    []domain.Component{{Name: k8sK8sDoguOperator, TargetState: domain.TargetStateAbsent}},
			wantErr: assert.NoError,
		},
		{
			name:    "unparsable version",
			args:    args{components: []bpv2.Component{{Name: "k8s/k8s-dogu-operator", Version: "1."}}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "invalid component name",
			args:    args{components: []bpv2.Component{{Name: "k8s/k8s-dogu-operator/oh/no", Version: "1.0.0"}}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "does not contain distribution namespace",
			args:    args{components: []bpv2.Component{{Name: "k8s-dogu-operator", Version: version3211.Raw}}},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertComponents(tt.args.components)
			if !tt.wantErr(t, err, fmt.Sprintf("ConvertComponents(%v)", tt.args.components)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ConvertComponents(%v)", tt.args.components)
		})
	}
}

func TestConvertToComponentDTOs(t *testing.T) {
	type args struct {
		components []domain.Component
	}
	tests := []struct {
		name    string
		args    args
		want    []bpv2.Component
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "nil",
			args:    args{},
			want:    []bpv2.Component{},
			wantErr: assert.NoError,
		},
		{
			name:    "empty list",
			args:    args{components: []domain.Component{}},
			want:    []bpv2.Component{},
			wantErr: assert.NoError,
		},
		{
			name:    "ok",
			args:    args{components: []domain.Component{{Name: k8sK8sDoguOperator, Version: compVersion3211, TargetState: domain.TargetStatePresent}}},
			want:    []bpv2.Component{{Name: "k8s/k8s-dogu-operator", Version: version3211.Raw, Absent: false}},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToComponentDTOs(tt.args.components)
			if !tt.wantErr(t, err, fmt.Sprintf("ConvertToComponentDTOs(%v)", tt.args.components)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ConvertToComponentDTOs(%v)", tt.args.components)
		})
	}
}
