package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer/blueprintMaskV1"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer/blueprintV2"
	v1 "github.com/cloudogu/k8s-blueprint-operator/pkg/api/v1"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"testing"
)

var ctx = context.Background()

func Test_blueprintSpecRepo_GetById(t *testing.T) {
	blueprintId := "MyBlueprint"

	t.Run("all ok", func(t *testing.T) {
		//given
		restClientMock := NewMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, blueprintV2.Serializer{}, blueprintMaskV1.Serializer{}, eventRecorderMock)

		cr := &v1.Blueprint{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{ResourceVersion: "abc"},
			Spec: v1.BlueprintSpec{
				Blueprint:                `{"blueprintApi": "v2"}`,
				BlueprintMask:            `{"blueprintMaskAPI": "v1"}`,
				AllowDoguNamespaceSwitch: true,
				IgnoreDoguHealth:         true,
			},
			Status: v1.BlueprintStatus{},
		}
		restClientMock.EXPECT().Get(ctx, blueprintId, metav1.GetOptions{}).Return(cr, nil)

		//when
		spec, err := repo.GetById(ctx, blueprintId)

		//then
		require.NoError(t, err)
		persistenceContext := make(map[string]interface{})
		persistenceContext[resourceVersionKey] = resourceVersionValue{"abc"}
		assert.Equal(t, domain.BlueprintSpec{
			Id: blueprintId,
			Config: domain.BlueprintConfiguration{
				IgnoreDoguHealth:         true,
				AllowDoguNamespaceSwitch: true,
			},
			PersistenceContext: persistenceContext,
		}, spec)
	})

	t.Run("invalid blueprint and mask", func(t *testing.T) {
		//given
		restClientMock := NewMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, blueprintV2.Serializer{}, blueprintMaskV1.Serializer{}, eventRecorderMock)

		cr := &v1.Blueprint{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{ResourceVersion: "abc"},
			Spec: v1.BlueprintSpec{
				Blueprint:     `{}`,
				BlueprintMask: `{}`,
			},
			Status: v1.BlueprintStatus{},
		}
		restClientMock.EXPECT().Get(ctx, blueprintId, metav1.GetOptions{}).Return(cr, nil)

		//when
		_, err := repo.GetById(ctx, blueprintId)

		//then
		require.Error(t, err)
		var expectedErrorType *domain.InvalidBlueprintError
		assert.ErrorAs(t, err, &expectedErrorType)
		assert.ErrorContains(t, err, fmt.Sprintf("could not deserialize blueprint CR %q: ", blueprintId))
		assert.ErrorContains(t, err, "cannot deserialize blueprint")
		assert.ErrorContains(t, err, "cannot deserialize blueprint mask")
	})

	t.Run("internal error while loading", func(t *testing.T) {
		//given
		restClientMock := NewMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, blueprintV2.Serializer{}, blueprintMaskV1.Serializer{}, eventRecorderMock)

		restClientMock.EXPECT().Get(ctx, blueprintId, metav1.GetOptions{}).Return(nil, k8sErrors.NewInternalError(errors.New("test-error")))

		//when
		_, err := repo.GetById(ctx, blueprintId)

		//then
		require.Error(t, err)
		var expectedErrorType *domainservice.InternalError
		assert.ErrorAs(t, err, &expectedErrorType)
		assert.ErrorContains(t, err, fmt.Sprintf("error while loading blueprint CR %q:", blueprintId))
		assert.ErrorContains(t, err, "test-error")
	})

	t.Run("not found error while loading", func(t *testing.T) {
		//given
		restClientMock := NewMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, blueprintV2.Serializer{}, blueprintMaskV1.Serializer{}, eventRecorderMock)

		restClientMock.EXPECT().
			Get(ctx, blueprintId, metav1.GetOptions{}).
			Return(nil, k8sErrors.NewNotFound(schema.GroupResource{}, blueprintId))

		//when
		_, err := repo.GetById(ctx, blueprintId)

		//then
		require.Error(t, err)
		var expectedErrorType *domainservice.NotFoundError
		assert.ErrorAs(t, err, &expectedErrorType)
		assert.ErrorContains(t, err, fmt.Sprintf("cannot load blueprint CR %q as it does not exist:", blueprintId))
	})
}

func Test_blueprintSpecRepo_Update(t *testing.T) {
	blueprintId := "MyBlueprint"

	t.Run("all ok", func(t *testing.T) {
		//given
		restClientMock := NewMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, blueprintV2.Serializer{}, blueprintMaskV1.Serializer{}, eventRecorderMock)
		expected := v1.Blueprint{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:            blueprintId,
				ResourceVersion: "abc",
			},
			Spec: v1.BlueprintSpec{},
			Status: v1.BlueprintStatus{
				Phase: v1.StatusPhase(domain.StatusPhaseValidated),
			},
		}
		restClientMock.EXPECT().
			UpdateStatus(ctx, mock.Anything, metav1.UpdateOptions{}).
			RunAndReturn(func(ctx2 context.Context, blueprint *v1.Blueprint, options metav1.UpdateOptions) (*v1.Blueprint, error) {
				assert.Equal(t, &expected, blueprint)
				return blueprint, nil
			})

		//when
		persistenceContext := make(map[string]interface{})
		persistenceContext[resourceVersionKey] = resourceVersionValue{"abc"}
		err := repo.Update(ctx, domain.BlueprintSpec{
			Id:                 blueprintId,
			Status:             domain.StatusPhaseValidated,
			Events:             nil,
			PersistenceContext: persistenceContext,
		})

		//then
		require.NoError(t, err)
	})

	t.Run("no version counter", func(t *testing.T) {
		//given
		restClientMock := NewMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, blueprintV2.Serializer{}, blueprintMaskV1.Serializer{}, eventRecorderMock)

		//when
		err := repo.Update(ctx, domain.BlueprintSpec{
			Id:     blueprintId,
			Status: domain.StatusPhaseValidated,
			Events: nil,
		})

		//then
		require.Error(t, err)
		assert.ErrorContains(t, err, "no resourceVersion was provided over the persistenceContext in the given blueprintSpec")
	})

	t.Run("version counter of different type", func(t *testing.T) {
		//given
		restClientMock := NewMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, blueprintV2.Serializer{}, blueprintMaskV1.Serializer{}, eventRecorderMock)

		//when
		persistenceContext := make(map[string]interface{})
		persistenceContext[resourceVersionKey] = 1
		err := repo.Update(ctx, domain.BlueprintSpec{
			Id:                 blueprintId,
			Status:             domain.StatusPhaseValidated,
			Events:             nil,
			PersistenceContext: persistenceContext,
		})

		//then
		require.Error(t, err)
		assert.ErrorContains(t, err, "resourceVersion in blueprintSpec is not a 'resourceVersionValue' but 'int'")
	})

	t.Run("conflict error", func(t *testing.T) {
		//given
		restClientMock := NewMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, blueprintV2.Serializer{}, blueprintMaskV1.Serializer{}, eventRecorderMock)
		expected := v1.Blueprint{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:            blueprintId,
				ResourceVersion: "abc",
			},
			Spec: v1.BlueprintSpec{},
			Status: v1.BlueprintStatus{
				Phase: v1.StatusPhase(domain.StatusPhaseValidated),
			},
		}
		expectedError := k8sErrors.NewConflict(
			schema.GroupResource{Group: "blueprints", Resource: blueprintId},
			blueprintId,
			fmt.Errorf("test-error"),
		)
		restClientMock.EXPECT().
			UpdateStatus(ctx, mock.Anything, metav1.UpdateOptions{}).
			RunAndReturn(func(ctx2 context.Context, blueprint *v1.Blueprint, options metav1.UpdateOptions) (*v1.Blueprint, error) {
				assert.Equal(t, &expected, blueprint)
				return nil, expectedError
			})

		//when
		persistenceContext := make(map[string]interface{})
		persistenceContext[resourceVersionKey] = resourceVersionValue{"abc"}
		err := repo.Update(ctx, domain.BlueprintSpec{
			Id:                 blueprintId,
			Status:             domain.StatusPhaseValidated,
			Events:             nil,
			PersistenceContext: persistenceContext,
		})

		//then
		require.Error(t, err)
		var expectedErrorType *domainservice.ConflictError
		assert.ErrorAs(t, err, &expectedErrorType)
		assert.ErrorIs(t, err, expectedError)
	})

	t.Run("internal error", func(t *testing.T) {
		//given
		restClientMock := NewMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, blueprintV2.Serializer{}, blueprintMaskV1.Serializer{}, eventRecorderMock)
		expected := v1.Blueprint{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:            blueprintId,
				ResourceVersion: "abc",
			},
			Spec: v1.BlueprintSpec{},
			Status: v1.BlueprintStatus{
				Phase: v1.StatusPhase(domain.StatusPhaseValidated),
			},
		}
		expectedError := fmt.Errorf("test-error")
		restClientMock.EXPECT().
			UpdateStatus(ctx, mock.Anything, metav1.UpdateOptions{}).
			RunAndReturn(func(ctx2 context.Context, blueprint *v1.Blueprint, options metav1.UpdateOptions) (*v1.Blueprint, error) {
				assert.Equal(t, &expected, blueprint)
				return nil, expectedError
			})

		//when
		persistenceContext := make(map[string]interface{})
		persistenceContext[resourceVersionKey] = resourceVersionValue{"abc"}
		err := repo.Update(ctx, domain.BlueprintSpec{
			Id:                 blueprintId,
			Status:             domain.StatusPhaseValidated,
			Events:             nil,
			PersistenceContext: persistenceContext,
		})

		//then
		require.Error(t, err)
		var expectedErrorType *domainservice.InternalError
		assert.ErrorAs(t, err, &expectedErrorType)
		assert.ErrorIs(t, err, expectedError)
	})
}

func Test_blueprintSpecRepo_Update_publishEvents(t *testing.T) {
	blueprintId := "MyBlueprint"
	t.Run("publish events", func(t *testing.T) {
		//given
		restClientMock := NewMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, blueprintV2.Serializer{}, blueprintMaskV1.Serializer{}, eventRecorderMock)
		restClientMock.EXPECT().UpdateStatus(ctx, mock.Anything, metav1.UpdateOptions{}).Return(nil, nil)

		var events []interface{}
		events = append(events,
			domain.BlueprintSpecValidatedEvent{},
			domain.BlueprintSpecInvalidEvent{ValidationError: errors.New("test-error")},
		)
		eventRecorderMock.EXPECT().Event(mock.Anything, corev1.EventTypeNormal, "BlueprintSpecValidatedEvent", "")
		eventRecorderMock.EXPECT().Event(mock.Anything, corev1.EventTypeNormal, "BlueprintSpecInvalidEvent", "test-error")

		//when
		persistenceContext := make(map[string]interface{})
		persistenceContext[resourceVersionKey] = resourceVersionValue{"abc"}
		err := repo.Update(ctx, domain.BlueprintSpec{Id: blueprintId, Events: events, PersistenceContext: persistenceContext})

		//then
		require.NoError(t, err)
	})
	t.Run("publish unknown event", func(t *testing.T) {
		//given
		restClientMock := NewMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, blueprintV2.Serializer{}, blueprintMaskV1.Serializer{}, eventRecorderMock)
		restClientMock.EXPECT().UpdateStatus(ctx, mock.Anything, metav1.UpdateOptions{}).Return(nil, nil)

		var events []interface{}
		events = append(
			events,
			"myString",
			struct{ key1 string }{key1: "val1"},
		)
		eventRecorderMock.EXPECT().Event(mock.Anything, corev1.EventTypeNormal, "Unknown", "unknown event of type 'string': myString")
		eventRecorderMock.EXPECT().Event(mock.Anything, corev1.EventTypeNormal, "Unknown", "unknown event of type 'struct { key1 string }': {key1:val1}")

		//when
		persistenceContext := make(map[string]interface{})
		persistenceContext[resourceVersionKey] = resourceVersionValue{"abc"}
		err := repo.Update(ctx, domain.BlueprintSpec{Id: blueprintId, Events: events, PersistenceContext: persistenceContext})

		//then
		require.NoError(t, err)
	})
}
