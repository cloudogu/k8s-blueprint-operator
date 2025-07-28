package sensitiveconfigref

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"testing"
)

var testCtx = context.TODO()

func TestSecretRefReader_ExistAll(t *testing.T) {

	type fields struct {
		secretClient v1.SecretInterface
	}
	type args struct {
		ctx  context.Context
		refs []domain.SensitiveValueRef
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "nothing to check",
			fields: fields{
				secretClient: newMockSecretClient(t),
			},
			args: args{
				ctx:  testCtx,
				refs: nil,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &SecretRefReader{
				secretClient: tt.fields.secretClient,
			}
			got, err := reader.ExistAll(tt.args.ctx, tt.args.refs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExistAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExistAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}
