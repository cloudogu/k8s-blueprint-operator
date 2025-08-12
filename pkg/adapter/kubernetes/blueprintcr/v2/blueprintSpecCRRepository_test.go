package v2

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	bpv2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
)

var ctx = context.Background()

func Test_blueprintSpecRepo_GetById(t *testing.T) {
	blueprintId := "MyBlueprint"

	t.Run("all ok", func(t *testing.T) {
		// given
		restClientMock := newMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, eventRecorderMock)

		cr := &bpv2.Blueprint{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{ResourceVersion: "abc"},
			Spec: bpv2.BlueprintSpec{
				Blueprint:                bpv2.BlueprintManifest{},
				BlueprintMask:            bpv2.BlueprintMask{},
				AllowDoguNamespaceSwitch: true,
				IgnoreDoguHealth:         true,
				DryRun:                   true,
			},
			Status: bpv2.BlueprintStatus{},
		}
		restClientMock.EXPECT().Get(ctx, blueprintId, metav1.GetOptions{}).Return(cr, nil)

		// when
		spec, err := repo.GetById(ctx, blueprintId)

		// then
		require.NoError(t, err)
		persistenceContext := make(map[string]interface{})
		persistenceContext[blueprintSpecRepoContextKey] = blueprintSpecRepoContext{"abc"}
		assert.Equal(t, &domain.BlueprintSpec{
			Id: blueprintId,
			Config: domain.BlueprintConfiguration{
				IgnoreDoguHealth:         true,
				AllowDoguNamespaceSwitch: true,
				DryRun:                   true,
			},
			StateDiff:          domain.StateDiff{},
			PersistenceContext: persistenceContext,
		}, spec)
	})

	t.Run("invalid blueprint and mask", func(t *testing.T) {
		// given
		restClientMock := newMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)

		cr := &bpv2.Blueprint{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{ResourceVersion: "abc"},
			Spec: bpv2.BlueprintSpec{
				Blueprint: bpv2.BlueprintManifest{
					Dogus: []bpv2.Dogu{
						{Name: "invalid"},
					},
				},
				BlueprintMask: bpv2.BlueprintMask{
					Dogus: []bpv2.MaskDogu{
						{Name: "invalid"},
					},
				},
			},
			Status: bpv2.BlueprintStatus{},
		}
		eventRecorderMock.EXPECT().Event(cr, "Warning", "BlueprintSpecInvalid", "cannot deserialize blueprint: cannot convert blueprint dogus: dogu name needs to be in the form 'namespace/dogu' but is 'invalid'")
		eventRecorderMock.EXPECT().Event(cr, "Warning", "BlueprintSpecInvalid", "cannot deserialize blueprint mask: cannot convert blueprint dogus: dogu name needs to be in the form 'namespace/dogu' but is 'invalid'")
		repo := NewBlueprintSpecRepository(restClientMock, eventRecorderMock)
		restClientMock.EXPECT().Get(ctx, blueprintId, metav1.GetOptions{}).Return(cr, nil)

		// when
		_, err := repo.GetById(ctx, blueprintId)

		// then
		require.Error(t, err)
		var expectedErrorType *domain.InvalidBlueprintError
		assert.ErrorAs(t, err, &expectedErrorType)
		assert.ErrorContains(t, err, fmt.Sprintf("could not deserialize blueprint CR %q: ", blueprintId))
		assert.ErrorContains(t, err, "cannot deserialize blueprint")
		assert.ErrorContains(t, err, "cannot deserialize blueprint mask")
	})

	t.Run("internal error while loading", func(t *testing.T) {
		// given
		restClientMock := newMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, eventRecorderMock)

		restClientMock.EXPECT().Get(ctx, blueprintId, metav1.GetOptions{}).Return(nil, k8sErrors.NewInternalError(errors.New("test-error")))

		// when
		_, err := repo.GetById(ctx, blueprintId)

		// then
		require.Error(t, err)
		var expectedErrorType *domainservice.InternalError
		assert.ErrorAs(t, err, &expectedErrorType)
		assert.ErrorContains(t, err, fmt.Sprintf("error while loading blueprint CR %q:", blueprintId))
		assert.ErrorContains(t, err, "test-error")
	})

	t.Run("not found error while loading", func(t *testing.T) {
		// given
		restClientMock := newMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, eventRecorderMock)

		restClientMock.EXPECT().
			Get(ctx, blueprintId, metav1.GetOptions{}).
			Return(nil, k8sErrors.NewNotFound(schema.GroupResource{}, blueprintId))

		// when
		_, err := repo.GetById(ctx, blueprintId)

		// then
		require.Error(t, err)
		var expectedErrorType *domainservice.NotFoundError
		assert.ErrorAs(t, err, &expectedErrorType)
		assert.ErrorContains(t, err, fmt.Sprintf("cannot load blueprint CR %q as it does not exist:", blueprintId))
	})
}

func Test_blueprintSpecRepo_Update(t *testing.T) {
	blueprintId := "MyBlueprint"

	t.Run("all ok", func(t *testing.T) {
		// given
		restClientMock := newMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, eventRecorderMock)
		expectedStatus := bpv2.BlueprintStatus{
			Phase: bpv2.StatusPhaseValidated,
			EffectiveBlueprint: bpv2.BlueprintManifest{
				Dogus:      []bpv2.Dogu{},
				Components: []bpv2.Component{},
				Config:     bpv2.Config{},
			},
			StateDiff: bpv2.StateDiff{DoguDiffs: map[string]bpv2.DoguDiff{}, ComponentDiffs: map[string]bpv2.ComponentDiff{}},
		}
		restClientMock.EXPECT().
			UpdateStatus(ctx, mock.Anything, metav1.UpdateOptions{}).
			RunAndReturn(func(ctx2 context.Context, blueprint *bpv2.Blueprint, options metav1.UpdateOptions) (*bpv2.Blueprint, error) {
				assert.Equal(t, expectedStatus, blueprint.Status)
				return blueprint, nil
			})

		// when
		persistenceContext := make(map[string]interface{})
		persistenceContext[blueprintSpecRepoContextKey] = blueprintSpecRepoContext{"abc"}
		err := repo.Update(ctx, &domain.BlueprintSpec{
			Id:                 blueprintId,
			Status:             domain.StatusPhaseValidated,
			Events:             nil,
			PersistenceContext: persistenceContext,
		})

		// then
		require.NoError(t, err)
	})

	t.Run("no version counter", func(t *testing.T) {
		// given
		restClientMock := newMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, eventRecorderMock)

		// when
		err := repo.Update(ctx, &domain.BlueprintSpec{
			Id:     blueprintId,
			Status: domain.StatusPhaseValidated,
			Events: nil,
		})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "no blueprintSpecRepoContext was provided over the persistenceContext in the given blueprintSpec")
	})

	t.Run("version counter of different type", func(t *testing.T) {
		// given
		restClientMock := newMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, eventRecorderMock)

		// when
		persistenceContext := make(map[string]interface{})
		persistenceContext[blueprintSpecRepoContextKey] = 1
		err := repo.Update(ctx, &domain.BlueprintSpec{
			Id:                 blueprintId,
			Status:             domain.StatusPhaseValidated,
			Events:             nil,
			PersistenceContext: persistenceContext,
		})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "persistence context in blueprintSpec is not a 'blueprintSpecRepoContext' but 'int'")
	})

	t.Run("conflict error on status update", func(t *testing.T) {
		// given
		restClientMock := newMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, eventRecorderMock)
		expectedStatus := bpv2.BlueprintStatus{
			Phase: bpv2.StatusPhaseValidated,
			EffectiveBlueprint: bpv2.BlueprintManifest{
				Dogus:      []bpv2.Dogu{},
				Components: []bpv2.Component{},
				Config:     bpv2.Config{},
			},
			StateDiff: bpv2.StateDiff{DoguDiffs: map[string]bpv2.DoguDiff{}, ComponentDiffs: map[string]bpv2.ComponentDiff{}},
		}
		expectedError := k8sErrors.NewConflict(
			schema.GroupResource{Group: "blueprints", Resource: blueprintId},
			blueprintId,
			fmt.Errorf("test-error"),
		)
		restClientMock.EXPECT().
			UpdateStatus(ctx, mock.Anything, metav1.UpdateOptions{}).
			RunAndReturn(func(ctx2 context.Context, blueprint *bpv2.Blueprint, options metav1.UpdateOptions) (*bpv2.Blueprint, error) {
				assert.Equal(t, expectedStatus, blueprint.Status)
				return nil, expectedError
			})

		// when
		persistenceContext := make(map[string]interface{})
		persistenceContext[blueprintSpecRepoContextKey] = blueprintSpecRepoContext{"abc"}
		err := repo.Update(ctx, &domain.BlueprintSpec{
			Id:                 blueprintId,
			Status:             domain.StatusPhaseValidated,
			Events:             nil,
			PersistenceContext: persistenceContext,
		})

		// then
		require.Error(t, err)
		var expectedErrorType *domainservice.ConflictError
		assert.ErrorAs(t, err, &expectedErrorType)
		assert.ErrorIs(t, err, expectedError)
	})

	t.Run("internal error on status update", func(t *testing.T) {
		// given
		restClientMock := newMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, eventRecorderMock)
		expectedStatus := bpv2.BlueprintStatus{
			Phase: bpv2.StatusPhaseValidated,
			EffectiveBlueprint: bpv2.BlueprintManifest{
				Dogus:      []bpv2.Dogu{},
				Components: []bpv2.Component{},
				Config:     bpv2.Config{},
			},
			StateDiff: bpv2.StateDiff{DoguDiffs: map[string]bpv2.DoguDiff{}, ComponentDiffs: map[string]bpv2.ComponentDiff{}},
		}
		expectedError := fmt.Errorf("test-error")
		restClientMock.EXPECT().
			UpdateStatus(ctx, mock.Anything, metav1.UpdateOptions{}).
			RunAndReturn(func(ctx2 context.Context, blueprint *bpv2.Blueprint, options metav1.UpdateOptions) (*bpv2.Blueprint, error) {
				assert.Equal(t, expectedStatus, blueprint.Status)
				return nil, expectedError
			})

		// when
		persistenceContext := make(map[string]interface{})
		persistenceContext[blueprintSpecRepoContextKey] = blueprintSpecRepoContext{"abc"}
		err := repo.Update(ctx, &domain.BlueprintSpec{
			Id:                 blueprintId,
			Status:             domain.StatusPhaseValidated,
			Events:             nil,
			PersistenceContext: persistenceContext,
		})

		// then
		require.Error(t, err)
		var expectedErrorType *domainservice.InternalError
		assert.ErrorAs(t, err, &expectedErrorType)
		assert.ErrorIs(t, err, expectedError)
	})
}

func Test_blueprintSpecRepo_Update_publishEvents(t *testing.T) {
	blueprintId := "MyBlueprint"
	t.Run("publish events", func(t *testing.T) {
		// given
		restClientMock := newMockBlueprintInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(restClientMock, eventRecorderMock)
		restClientMock.EXPECT().
			UpdateStatus(ctx, mock.Anything, metav1.UpdateOptions{}).
			RunAndReturn(func(ctx2 context.Context, blueprint *bpv2.Blueprint, options metav1.UpdateOptions) (*bpv2.Blueprint, error) {
				// assert.Equal(t, &expected, blueprint)
				blueprint.ResourceVersion = "newVersion"
				return blueprint, nil
			})

		var events []domain.Event
		events = append(events,
			domain.BlueprintSpecValidatedEvent{},
			domain.EffectiveBlueprintCalculatedEvent{},
			domain.StateDiffDoguDeterminedEvent{},
			domain.StateDiffComponentDeterminedEvent{},
			domain.EcosystemHealthyUpfrontEvent{},
			domain.EcosystemUnhealthyUpfrontEvent{HealthResult: ecosystem.HealthResult{}},
			domain.BlueprintSpecInvalidEvent{ValidationError: errors.New("test-error")},
		)
		eventRecorderMock.EXPECT().Event(mock.Anything, corev1.EventTypeNormal, "BlueprintSpecValidated", "")
		eventRecorderMock.EXPECT().Event(mock.Anything, corev1.EventTypeNormal, "EffectiveBlueprintCalculated", "")
		eventRecorderMock.EXPECT().Event(mock.Anything, corev1.EventTypeNormal, "StateDiffDoguDetermined", "dogu state diff determined: 0 actions ()")
		eventRecorderMock.EXPECT().Event(mock.Anything, corev1.EventTypeNormal, "StateDiffComponentDetermined", "component state diff determined: 0 actions ()")
		eventRecorderMock.EXPECT().Event(mock.Anything, corev1.EventTypeNormal, "EcosystemHealthyUpfront", "dogu health ignored: false; component health ignored: false")
		eventRecorderMock.EXPECT().Event(mock.Anything, corev1.EventTypeNormal, "EcosystemUnhealthyUpfront", "ecosystem health:\n  0 dogu(s) are unhealthy: \n  0 component(s) are unhealthy: ")
		eventRecorderMock.EXPECT().Event(mock.Anything, corev1.EventTypeNormal, "BlueprintSpecInvalid", "test-error")

		// when
		persistenceContext := make(map[string]interface{})
		persistenceContext[blueprintSpecRepoContextKey] = blueprintSpecRepoContext{"abc"}
		spec := &domain.BlueprintSpec{Id: blueprintId, Events: events, PersistenceContext: persistenceContext}
		err := repo.Update(ctx, spec)

		// then
		require.NoError(t, err)
		newPersistenceContext, _ := getPersistenceContext(ctx, spec)
		assert.Equal(t, "newVersion", newPersistenceContext.resourceVersion)
		assert.Empty(t, spec.Events, "events in aggregate should be deleted after publishing them")
	})
}
