package config

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/stretchr/testify/assert"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

var testCtx = context.Background()

func TestHealthConfigRepository_Get(t *testing.T) {
	notFoundErr := errors.NewNotFound(schema.GroupResource{
		Group:    "v1",
		Resource: "ConfigMap",
	}, "k8s-blueprint-operator-health-config")
	tests := []struct {
		name       string
		cmClientFn func(t *testing.T) configMapInterface
		want       domain.HealthConfig
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "should fail to get configmap due to not found error",
			cmClientFn: func(t *testing.T) configMapInterface {
				cmMock := newMockConfigMapInterface(t)
				cmMock.EXPECT().Get(testCtx, "k8s-blueprint-operator-health-config", metav1.GetOptions{}).
					Return(nil, notFoundErr)
				return cmMock
			},
			want: domain.HealthConfig{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				expectedNotFound := &domainservice.NotFoundError{}
				return assert.ErrorIs(t, err, notFoundErr, i) &&
					assert.ErrorAs(t, err, &expectedNotFound, i) &&
					assert.ErrorContains(t, err, "could not find health config map \"k8s-blueprint-operator-health-config\"", i)
			},
		},
		{
			name: "should fail to get configmap due to other error",
			cmClientFn: func(t *testing.T) configMapInterface {
				cmMock := newMockConfigMapInterface(t)
				cmMock.EXPECT().Get(testCtx, "k8s-blueprint-operator-health-config", metav1.GetOptions{}).
					Return(nil, assert.AnError)
				return cmMock
			},
			want: domain.HealthConfig{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to get config map \"k8s-blueprint-operator-health-config\"", i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HealthConfigRepository{
				cmClient: tt.cmClientFn(t),
			}
			got, err := h.Get(testCtx)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
