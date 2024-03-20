package healthconfig

import (
	"context"
	"math/rand"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/stretchr/testify/assert"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

var testCtx = context.Background()

func TestHealthConfigRepository_GetRequiredComponents(t *testing.T) {
	notFoundErr := errors.NewNotFound(schema.GroupResource{
		Group:    "v1",
		Resource: "ConfigMap",
	}, "k8s-blueprint-operator-health-config")
	tests := []struct {
		name       string
		cmClientFn func(t *testing.T) configMapInterface
		want       []ecosystem.RequiredComponent
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "should return default on not found error",
			cmClientFn: func(t *testing.T) configMapInterface {
				cmMock := newMockConfigMapInterface(t)
				cmMock.EXPECT().Get(testCtx, "k8s-blueprint-operator-health-config", metav1.GetOptions{}).
					Return(nil, notFoundErr)
				return cmMock
			},
			want:    make([]ecosystem.RequiredComponent, 0),
			wantErr: assert.NoError,
		},
		{
			name: "should fail to get configmap due to other error",
			cmClientFn: func(t *testing.T) configMapInterface {
				cmMock := newMockConfigMapInterface(t)
				cmMock.EXPECT().Get(testCtx, "k8s-blueprint-operator-health-config", metav1.GetOptions{}).
					Return(nil, assert.AnError)
				return cmMock
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to get config map \"k8s-blueprint-operator-health-config\"", i)
			},
		},
		{
			name: "should return default if config key does not exist",
			cmClientFn: func(t *testing.T) configMapInterface {
				cmMock := newMockConfigMapInterface(t)
				configMap := &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-blueprint-operator-health-config"},
					Data:       map[string]string{},
				}
				cmMock.EXPECT().Get(testCtx, "k8s-blueprint-operator-health-config", metav1.GetOptions{}).
					Return(configMap, nil)
				return cmMock
			},
			want:    []ecosystem.RequiredComponent{},
			wantErr: assert.NoError,
		},
		{
			name: "should fail to parse component config",
			cmClientFn: func(t *testing.T) configMapInterface {
				cmMock := newMockConfigMapInterface(t)
				configMap := &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-blueprint-operator-health-config"},
					Data:       map[string]string{"components": `{"required": [{]}`},
				}
				cmMock.EXPECT().Get(testCtx, "k8s-blueprint-operator-health-config", metav1.GetOptions{}).
					Return(configMap, nil)
				return cmMock
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse component health config", i)
			},
		},
		{
			name: "should succeed",
			cmClientFn: func(t *testing.T) configMapInterface {
				cmMock := newMockConfigMapInterface(t)
				configMap := &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-blueprint-operator-health-config"},
					Data:       map[string]string{"components": `{"required": [{"name": "k8s-dogu-operator"}, {"name": "k8s-etcd"}]}`},
				}
				cmMock.EXPECT().Get(testCtx, "k8s-blueprint-operator-health-config", metav1.GetOptions{}).
					Return(configMap, nil)
				return cmMock
			},
			want:    []ecosystem.RequiredComponent{{Name: "k8s-dogu-operator"}, {Name: "k8s-etcd"}},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HealthConfigProvider{
				cmClient: tt.cmClientFn(t),
			}
			got, err := h.GetRequiredComponents(testCtx)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHealthConfigRepository_GetWaitConfig(t *testing.T) {
	notFoundErr := errors.NewNotFound(schema.GroupResource{
		Group:    "v1",
		Resource: "ConfigMap",
	}, "k8s-blueprint-operator-health-config")
	tests := []struct {
		name       string
		cmClientFn func(t *testing.T) configMapInterface
		want       ecosystem.WaitConfig
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "should return default on not found error",
			cmClientFn: func(t *testing.T) configMapInterface {
				cmMock := newMockConfigMapInterface(t)
				cmMock.EXPECT().Get(testCtx, "k8s-blueprint-operator-health-config", metav1.GetOptions{}).
					Return(nil, notFoundErr)
				return cmMock
			},
			want: ecosystem.WaitConfig{
				Timeout:  10 * time.Minute,
				Interval: 10 * time.Second,
			},
			wantErr: assert.NoError,
		},
		{
			name: "should fail to get config map due to other error",
			cmClientFn: func(t *testing.T) configMapInterface {
				cmMock := newMockConfigMapInterface(t)
				cmMock.EXPECT().Get(testCtx, "k8s-blueprint-operator-health-config", metav1.GetOptions{}).
					Return(nil, assert.AnError)
				return cmMock
			},
			want: ecosystem.WaitConfig{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to get config map \"k8s-blueprint-operator-health-config\"", i)
			},
		},
		{
			name: "should return default if config key does not exist",
			cmClientFn: func(t *testing.T) configMapInterface {
				cmMock := newMockConfigMapInterface(t)
				configMap := &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-blueprint-operator-health-config"},
					Data:       map[string]string{},
				}
				cmMock.EXPECT().Get(testCtx, "k8s-blueprint-operator-health-config", metav1.GetOptions{}).
					Return(configMap, nil)
				return cmMock
			},
			want: ecosystem.WaitConfig{
				Timeout:  10 * time.Minute,
				Interval: 10 * time.Second,
			},
			wantErr: assert.NoError,
		},
		{
			name: "should fail to parse wait config",
			cmClientFn: func(t *testing.T) configMapInterface {
				cmMock := newMockConfigMapInterface(t)
				configMap := &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-blueprint-operator-health-config"},
					Data:       map[string]string{"wait": `{"timeout": {[[{{}`},
				}
				cmMock.EXPECT().Get(testCtx, "k8s-blueprint-operator-health-config", metav1.GetOptions{}).
					Return(configMap, nil)
				return cmMock
			},
			want: ecosystem.WaitConfig{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse wait health config", i)
			},
		},
		{
			name: "should return default timeout if empty",
			cmClientFn: func(t *testing.T) configMapInterface {
				cmMock := newMockConfigMapInterface(t)
				configMap := &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-blueprint-operator-health-config"},
					Data:       map[string]string{"wait": `{"interval": 2}`},
				}
				cmMock.EXPECT().Get(testCtx, "k8s-blueprint-operator-health-config", metav1.GetOptions{}).
					Return(configMap, nil)
				return cmMock
			},
			want: ecosystem.WaitConfig{
				Timeout:  10 * time.Minute,
				Interval: 2,
			},
			wantErr: assert.NoError,
		},
		{
			name: "should return default interval if empty",
			cmClientFn: func(t *testing.T) configMapInterface {
				cmMock := newMockConfigMapInterface(t)
				configMap := &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-blueprint-operator-health-config"},
					Data:       map[string]string{"wait": `{"timeout": "50s"}`},
				}
				cmMock.EXPECT().Get(testCtx, "k8s-blueprint-operator-health-config", metav1.GetOptions{}).
					Return(configMap, nil)
				return cmMock
			},
			want: ecosystem.WaitConfig{
				Timeout:  50 * time.Second,
				Interval: 10 * time.Second,
			},
			wantErr: assert.NoError,
		},
		{
			name: "should succeed",
			cmClientFn: func(t *testing.T) configMapInterface {
				cmMock := newMockConfigMapInterface(t)
				configMap := &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-blueprint-operator-health-config"},
					Data:       map[string]string{"wait": `{"timeout": "120s", "interval": "13s"}`},
				}
				cmMock.EXPECT().Get(testCtx, "k8s-blueprint-operator-health-config", metav1.GetOptions{}).
					Return(configMap, nil)
				return cmMock
			},
			want: ecosystem.WaitConfig{
				Timeout:  120 * time.Second,
				Interval: 13 * time.Second,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HealthConfigProvider{
				cmClient: tt.cmClientFn(t),
			}
			got, err := h.GetWaitConfig(testCtx)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewHealthConfigRepository(t *testing.T) {
	mock := newMockConfigMapInterface(t)
	provider := NewHealthConfigProvider(mock)
	assert.Same(t, mock, provider.cmClient)
}

func Test_duration_Marshal_Unmarshal_JSON(t *testing.T) {
	// we provide a seed to produce always the exact same sequence of values
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		d1 := duration{time.Duration(rng.Int63())}
		json, err := d1.MarshalJSON()
		assert.NoErrorf(t, err, "failed to marshal duration %q to json", d1)

		var d2 duration
		err = d2.UnmarshalJSON(json)
		assert.NoErrorf(t, err, "failed to unmarshal json %q (original: %q)", string(json), d1)
	}
}
