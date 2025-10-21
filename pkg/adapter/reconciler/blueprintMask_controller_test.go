package reconciler

import (
	"testing"

	bpv2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestNewBlueprintMaskReconciler(t *testing.T) {
	blueprintClient := newMockBlueprintInterface(t)
	blueprintMaskClient := newMockBlueprintMaskInterface(t)
	blueprintEvents := make(chan<- event.TypedGenericEvent[*bpv2.Blueprint])
	reconciler := NewBlueprintMaskReconciler(blueprintClient, blueprintMaskClient, blueprintEvents)

	require.NotEmpty(t, reconciler)
	assert.Same(t, blueprintClient, reconciler.blueprintInterface)
	assert.Same(t, blueprintMaskClient, reconciler.blueprintMaskInterface)
	assert.Equal(t, blueprintEvents, reconciler.blueprintEvents)
}

func TestBlueprintMaskReconciler_SetupWithManager(t *testing.T) {
	tests := []struct {
		name    string
		mgrFn   func(t *testing.T) controllerruntime.Manager
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "fail on nil manager",
			mgrFn: func(t *testing.T) controllerruntime.Manager {
				return nil
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "must provide a non-nil manager")
			},
		},
		{
			name: "success",
			mgrFn: func(t *testing.T) controllerruntime.Manager {
				ctrlManMock := newMockControllerManager(t)
				ctrlManMock.EXPECT().GetControllerOptions().Return(config.Controller{})
				ctrlManMock.EXPECT().GetScheme().Return(createScheme(t))
				logger := log.FromContext(testCtx)
				ctrlManMock.EXPECT().GetLogger().Return(logger)
				ctrlManMock.EXPECT().Add(mock.Anything).Return(nil)
				ctrlManMock.EXPECT().GetCache().Return(nil)
				return ctrlManMock
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &BlueprintMaskReconciler{}
			tt.wantErr(t, r.SetupWithManager(tt.mgrFn(t)))
		})
	}
}

func TestBlueprintMaskReconciler_Reconcile(t *testing.T) {
	type fields struct {
		blueprintInterface     func(t *testing.T) blueprintInterface
		blueprintMaskInterface func(t *testing.T) blueprintMaskInterface
	}
	tests := []struct {
		name           string
		fields         fields
		expectedEvents []event.TypedGenericEvent[*bpv2.Blueprint]
		want           controllerruntime.Result
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "fail to list blueprint masks",
			fields: fields{
				blueprintInterface: func(t *testing.T) blueprintInterface {
					mck := newMockBlueprintInterface(t)
					return mck
				},
				blueprintMaskInterface: func(t *testing.T) blueprintMaskInterface {
					mck := newMockBlueprintMaskInterface(t)
					mck.EXPECT().List(testCtx, metav1.ListOptions{}).Return(nil, assert.AnError)
					return mck
				},
			},
			want: controllerruntime.Result{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError) &&
					assert.ErrorContains(t, err, "failed to list blueprint masks")
			},
		},
		{
			name: "abort if no blueprint masks exist",
			fields: fields{
				blueprintInterface: func(t *testing.T) blueprintInterface {
					mck := newMockBlueprintInterface(t)
					return mck
				},
				blueprintMaskInterface: func(t *testing.T) blueprintMaskInterface {
					mck := newMockBlueprintMaskInterface(t)
					mck.EXPECT().List(testCtx, metav1.ListOptions{}).Return(&bpv2.BlueprintMaskList{}, nil)
					return mck
				},
			},
			want:    controllerruntime.Result{},
			wantErr: assert.NoError,
		},
		{
			name: "fail to list blueprint masks",
			fields: fields{
				blueprintInterface: func(t *testing.T) blueprintInterface {
					mck := newMockBlueprintInterface(t)
					mck.EXPECT().List(testCtx, metav1.ListOptions{}).Return(nil, assert.AnError)
					return mck
				},
				blueprintMaskInterface: func(t *testing.T) blueprintMaskInterface {
					mck := newMockBlueprintMaskInterface(t)
					mck.EXPECT().List(testCtx, metav1.ListOptions{}).Return(&bpv2.BlueprintMaskList{Items: make([]bpv2.BlueprintMask, 2)}, nil)
					return mck
				},
			},
			want: controllerruntime.Result{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError) &&
					assert.ErrorContains(t, err, "failed to list blueprints")
			},
		},
		{
			name: "succeed with expected events",
			fields: fields{
				blueprintInterface: func(t *testing.T) blueprintInterface {
					mck := newMockBlueprintInterface(t)
					mck.EXPECT().List(testCtx, metav1.ListOptions{}).Return(&bpv2.BlueprintList{Items: []bpv2.Blueprint{
						{Spec: bpv2.BlueprintSpec{BlueprintMaskRef: &bpv2.BlueprintMaskRef{Name: "baz"}}},
						{Spec: bpv2.BlueprintSpec{BlueprintMaskRef: &bpv2.BlueprintMaskRef{Name: "foo"}}},
					}}, nil)
					return mck
				},
				blueprintMaskInterface: func(t *testing.T) blueprintMaskInterface {
					mck := newMockBlueprintMaskInterface(t)
					mck.EXPECT().List(testCtx, metav1.ListOptions{}).Return(&bpv2.BlueprintMaskList{Items: []bpv2.BlueprintMask{
						{ObjectMeta: metav1.ObjectMeta{Name: "foo"}},
						{ObjectMeta: metav1.ObjectMeta{Name: "bar"}},
						{ObjectMeta: metav1.ObjectMeta{Name: "baz"}},
					}}, nil)
					return mck
				},
			},
			expectedEvents: []event.TypedGenericEvent[*bpv2.Blueprint]{
				{Object: &bpv2.Blueprint{Spec: bpv2.BlueprintSpec{BlueprintMaskRef: &bpv2.BlueprintMaskRef{Name: "foo"}}}},
				{Object: &bpv2.Blueprint{Spec: bpv2.BlueprintSpec{BlueprintMaskRef: &bpv2.BlueprintMaskRef{Name: "baz"}}}},
			},
			want:    controllerruntime.Result{},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blueprintEvents := make(chan event.TypedGenericEvent[*bpv2.Blueprint])
			r := &BlueprintMaskReconciler{
				blueprintInterface:     tt.fields.blueprintInterface(t),
				blueprintMaskInterface: tt.fields.blueprintMaskInterface(t),
				blueprintEvents:        blueprintEvents,
			}

			go func() {
				var actualEvents []event.TypedGenericEvent[*bpv2.Blueprint]
				for e := range blueprintEvents {
					actualEvents = append(actualEvents, e)
				}

				assert.ElementsMatch(t, tt.expectedEvents, actualEvents)
			}()

			got, err := r.Reconcile(testCtx, reconcile.Request{NamespacedName: types.NamespacedName{Name: "test-mask"}})

			close(blueprintEvents)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
