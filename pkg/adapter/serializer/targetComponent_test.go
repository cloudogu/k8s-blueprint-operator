package serializer

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	compVersion3211    = semver.MustParse("3.2.1-1")
	k8sK8sDoguOperator = common.QualifiedComponentName{Namespace: "k8s", SimpleName: "k8s-dogu-operator"}
)

func TestConvertComponents(t *testing.T) {
	type args struct {
		components []TargetComponent
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
			args:    args{components: []TargetComponent{}},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name:    "normal component",
			args:    args{components: []TargetComponent{{Name: "k8s/k8s-dogu-operator", Version: version3211.Raw, TargetState: "present", DeployConfig: map[string]interface{}{"deployNamespace": "longhorn-system", "configOverwrite": map[string]string{"key": "value"}}}}},
			want:    []domain.Component{{Name: k8sK8sDoguOperator, Version: compVersion3211, TargetState: 0, DeployConfig: map[string]interface{}{"deployNamespace": "longhorn-system", "configOverwrite": map[string]string{"key": "value"}}}},
			wantErr: assert.NoError,
		},
		{
			name:    "unparsable version",
			args:    args{components: []TargetComponent{{Name: "k8s/k8s-dogu-operator", Version: "1.", TargetState: "present"}}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "invalid component name",
			args:    args{components: []TargetComponent{{Name: "k8s/k8s-dogu-operator/oh/no", Version: "1.0.0", TargetState: "present"}}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "unknown target state",
			args:    args{components: []TargetComponent{{Name: "k8s/k8s-dogu-operator", Version: version3211.Raw, TargetState: "unknown"}}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "does not contain distribution namespace",
			args:    args{components: []TargetComponent{{Name: "k8s-dogu-operator", Version: version3211.Raw, TargetState: "unknown"}}},
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
		want    []TargetComponent
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "nil",
			args:    args{},
			want:    []TargetComponent{},
			wantErr: assert.NoError,
		},
		{
			name:    "empty list",
			args:    args{components: []domain.Component{}},
			want:    []TargetComponent{},
			wantErr: assert.NoError,
		},
		{
			name:    "ok",
			args:    args{components: []domain.Component{{Name: k8sK8sDoguOperator, Version: compVersion3211, TargetState: domain.TargetStatePresent}}},
			want:    []TargetComponent{{Name: "k8s/k8s-dogu-operator", Version: version3211.Raw, TargetState: "present"}},
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
