package v2

import (
	"context"
	"errors"
	"fmt"
	"testing"

	bpv2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	ctx           = context.Background()
	testCondition = metav1.Condition{
		Type:               domain.ConditionCompleted,
		Status:             metav1.ConditionUnknown,
		ObservedGeneration: 1,
		LastTransitionTime: metav1.Time{},
		Reason:             "Completed",
		Message:            "test",
	}
	trueVar = true
)

func Test_blueprintSpecRepo_GetById(t *testing.T) {
	blueprintId := "MyBlueprint"

	t.Run("all ok", func(t *testing.T) {
		// given
		blueprintClientMock := newMockBlueprintInterface(t)
		maskClientMock := newMockBlueprintMaskInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(blueprintClientMock, maskClientMock, eventRecorderMock)

		cr := &bpv2.Blueprint{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{ResourceVersion: "abc"},
			Spec: bpv2.BlueprintSpec{
				DisplayName:              "MyBlueprint",
				Blueprint:                bpv2.BlueprintManifest{},
				BlueprintMask:            &bpv2.BlueprintMaskManifest{},
				AllowDoguNamespaceSwitch: &trueVar,
				IgnoreDoguHealth:         &trueVar,
				Stopped:                  &trueVar,
			},
			Status: &bpv2.BlueprintStatus{
				Conditions: []metav1.Condition{testCondition},
			},
		}
		blueprintClientMock.EXPECT().Get(ctx, blueprintId, metav1.GetOptions{}).Return(cr, nil)

		// when
		spec, err := repo.GetById(ctx, blueprintId)

		// then
		require.NoError(t, err)
		persistenceContext := make(map[string]interface{})
		persistenceContext[blueprintSpecRepoContextKey] = blueprintSpecRepoContext{"abc"}
		assert.Equal(t, &domain.BlueprintSpec{
			Id:          blueprintId,
			DisplayName: "MyBlueprint",
			Config: domain.BlueprintConfiguration{
				IgnoreDoguHealth:         true,
				AllowDoguNamespaceSwitch: true,
				Stopped:                  true,
			},
			StateDiff:          domain.StateDiff{},
			PersistenceContext: persistenceContext,
			Conditions:         []domain.Condition{testCondition},
		}, spec)
	})

	t.Run("all ok without status", func(t *testing.T) {
		// given
		blueprintClientMock := newMockBlueprintInterface(t)
		maskClientMock := newMockBlueprintMaskInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(blueprintClientMock, maskClientMock, eventRecorderMock)

		cr := &bpv2.Blueprint{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{ResourceVersion: "abc"},
			Spec: bpv2.BlueprintSpec{
				Blueprint:                bpv2.BlueprintManifest{},
				BlueprintMask:            &bpv2.BlueprintMaskManifest{},
				AllowDoguNamespaceSwitch: &trueVar,
				IgnoreDoguHealth:         &trueVar,
				Stopped:                  &trueVar,
			},
		}
		blueprintClientMock.EXPECT().Get(ctx, blueprintId, metav1.GetOptions{}).Return(cr, nil)

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
				Stopped:                  true,
			},
			StateDiff:          domain.StateDiff{},
			PersistenceContext: persistenceContext,
			Conditions:         nil,
		}, spec)
	})

	t.Run("invalid if both mask and mask ref are set", func(t *testing.T) {
		// given
		blueprintClientMock := newMockBlueprintInterface(t)
		maskClientMock := newMockBlueprintMaskInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(blueprintClientMock, maskClientMock, eventRecorderMock)

		cr := &bpv2.Blueprint{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{ResourceVersion: "abc"},
			Spec: bpv2.BlueprintSpec{
				DisplayName:              "MyBlueprint",
				Blueprint:                bpv2.BlueprintManifest{},
				BlueprintMask:            &bpv2.BlueprintMaskManifest{},
				BlueprintMaskRef:         &bpv2.BlueprintMaskRef{Name: "my-blueprint-mask"},
				AllowDoguNamespaceSwitch: &trueVar,
				IgnoreDoguHealth:         &trueVar,
				Stopped:                  &trueVar,
			},
			Status: &bpv2.BlueprintStatus{
				Conditions: []metav1.Condition{testCondition},
			},
		}
		eventRecorderMock.EXPECT().Event(cr, "Warning", "BlueprintSpecInvalid", "blueprint mask and mask ref cannot be set at the same time")
		blueprintClientMock.EXPECT().Get(ctx, blueprintId, metav1.GetOptions{}).Return(cr, nil)

		// when
		_, err := repo.GetById(ctx, blueprintId)

		// then
		require.Error(t, err)
		var expectedErrorType *domain.InvalidBlueprintError
		assert.ErrorAs(t, err, &expectedErrorType)
		assert.ErrorContains(t, err, fmt.Sprintf("could not deserialize blueprint CR %q: ", blueprintId))
		assert.ErrorContains(t, err, "blueprint mask and mask ref cannot be set at the same time")
	})

	t.Run("internal error when not able to get blueprint mask from ref", func(t *testing.T) {
		// given
		blueprintClientMock := newMockBlueprintInterface(t)
		maskClientMock := newMockBlueprintMaskInterface(t)
		maskClientMock.EXPECT().Get(ctx, "my-blueprint-mask", metav1.GetOptions{}).Return(nil, assert.AnError)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(blueprintClientMock, maskClientMock, eventRecorderMock)

		cr := &bpv2.Blueprint{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{ResourceVersion: "abc"},
			Spec: bpv2.BlueprintSpec{
				DisplayName:              "MyBlueprint",
				Blueprint:                bpv2.BlueprintManifest{},
				BlueprintMaskRef:         &bpv2.BlueprintMaskRef{Name: "my-blueprint-mask"},
				AllowDoguNamespaceSwitch: &trueVar,
				IgnoreDoguHealth:         &trueVar,
				Stopped:                  &trueVar,
			},
			Status: &bpv2.BlueprintStatus{
				Conditions: []metav1.Condition{testCondition},
			},
		}
		blueprintClientMock.EXPECT().Get(ctx, blueprintId, metav1.GetOptions{}).Return(cr, nil)

		// when
		_, err := repo.GetById(ctx, blueprintId)

		// then
		require.Error(t, err)
		var expectedErrorType *domainservice.InternalError
		assert.ErrorAs(t, err, &expectedErrorType)
		assert.ErrorContains(t, err, fmt.Sprintf("could not get blueprint mask from ref %q in blueprint %q", "my-blueprint-mask", blueprintId))
	})

	t.Run("all ok when getting mask from CR", func(t *testing.T) {
		// given
		blueprintClientMock := newMockBlueprintInterface(t)
		maskClientMock := newMockBlueprintMaskInterface(t)
		mask := &bpv2.BlueprintMask{Spec: bpv2.BlueprintMaskSpec{BlueprintMaskManifest: &bpv2.BlueprintMaskManifest{}}}
		maskClientMock.EXPECT().Get(ctx, "my-blueprint-mask", metav1.GetOptions{}).Return(mask, nil)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(blueprintClientMock, maskClientMock, eventRecorderMock)

		cr := &bpv2.Blueprint{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{ResourceVersion: "abc"},
			Spec: bpv2.BlueprintSpec{
				DisplayName:              "MyBlueprint",
				Blueprint:                bpv2.BlueprintManifest{},
				BlueprintMaskRef:         &bpv2.BlueprintMaskRef{Name: "my-blueprint-mask"},
				AllowDoguNamespaceSwitch: &trueVar,
				IgnoreDoguHealth:         &trueVar,
				Stopped:                  &trueVar,
			},
			Status: &bpv2.BlueprintStatus{
				Conditions: []metav1.Condition{testCondition},
			},
		}
		eventRecorderMock.EXPECT().Event(cr, "Normal", "BlueprintMaskFromRef", "Using blueprint mask from ref \"my-blueprint-mask\"")
		blueprintClientMock.EXPECT().Get(ctx, blueprintId, metav1.GetOptions{}).Return(cr, nil)

		// when
		spec, err := repo.GetById(ctx, blueprintId)

		// then
		require.NoError(t, err)
		persistenceContext := make(map[string]interface{})
		persistenceContext[blueprintSpecRepoContextKey] = blueprintSpecRepoContext{"abc"}
		assert.Equal(t, &domain.BlueprintSpec{
			Id:          blueprintId,
			DisplayName: "MyBlueprint",
			Config: domain.BlueprintConfiguration{
				IgnoreDoguHealth:         true,
				AllowDoguNamespaceSwitch: true,
				Stopped:                  true,
			},
			StateDiff:          domain.StateDiff{},
			PersistenceContext: persistenceContext,
			Conditions:         []domain.Condition{testCondition},
		}, spec)
	})

	t.Run("invalid blueprint and mask", func(t *testing.T) {
		// given
		blueprintClientMock := newMockBlueprintInterface(t)
		maskClientMock := newMockBlueprintMaskInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(blueprintClientMock, maskClientMock, eventRecorderMock)

		cr := &bpv2.Blueprint{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{ResourceVersion: "abc"},
			Spec: bpv2.BlueprintSpec{
				Blueprint: bpv2.BlueprintManifest{
					Dogus: []bpv2.Dogu{
						{Name: "invalid"},
					},
				},
				BlueprintMask: &bpv2.BlueprintMaskManifest{
					Dogus: []bpv2.MaskDogu{
						{Name: "invalid"},
					},
				},
			},
			Status: &bpv2.BlueprintStatus{},
		}
		errMsg := "cannot deserialize blueprint: cannot convert blueprint dogus: dogu name needs to be in the form 'namespace/dogu' but is 'invalid'\n" +
			"cannot deserialize blueprint mask: cannot convert blueprint dogus: dogu name needs to be in the form 'namespace/dogu' but is 'invalid'"
		eventRecorderMock.EXPECT().Event(cr, "Warning", "BlueprintSpecInvalid", errMsg)
		blueprintClientMock.EXPECT().Get(ctx, blueprintId, metav1.GetOptions{}).Return(cr, nil)

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
		blueprintClientMock := newMockBlueprintInterface(t)
		maskClientMock := newMockBlueprintMaskInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(blueprintClientMock, maskClientMock, eventRecorderMock)

		blueprintClientMock.EXPECT().Get(ctx, blueprintId, metav1.GetOptions{}).Return(nil, k8sErrors.NewInternalError(errors.New("test-error")))

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
		blueprintClientMock := newMockBlueprintInterface(t)
		maskClientMock := newMockBlueprintMaskInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(blueprintClientMock, maskClientMock, eventRecorderMock)

		blueprintClientMock.EXPECT().
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
		blueprintClientMock := newMockBlueprintInterface(t)
		maskClientMock := newMockBlueprintMaskInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(blueprintClientMock, maskClientMock, eventRecorderMock)
		expectedStatus := &bpv2.BlueprintStatus{
			EffectiveBlueprint: &bpv2.BlueprintManifest{
				Dogus:  []bpv2.Dogu{},
				Config: nil,
			},
			StateDiff:  &bpv2.StateDiff{DoguDiffs: map[string]bpv2.DoguDiff{}},
			Conditions: []metav1.Condition{testCondition},
		}
		blueprintClientMock.EXPECT().
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
			Events:             nil,
			PersistenceContext: persistenceContext,
			Conditions:         []domain.Condition{testCondition},
		})

		// then
		require.NoError(t, err)
	})

	t.Run("no version counter", func(t *testing.T) {
		// given
		blueprintClientMock := newMockBlueprintInterface(t)
		maskClientMock := newMockBlueprintMaskInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(blueprintClientMock, maskClientMock, eventRecorderMock)

		// when
		err := repo.Update(ctx, &domain.BlueprintSpec{
			Id:     blueprintId,
			Events: nil,
		})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "no blueprintSpecRepoContext was provided over the persistenceContext in the given blueprintSpec")
	})

	t.Run("version counter of different type", func(t *testing.T) {
		// given
		blueprintClientMock := newMockBlueprintInterface(t)
		maskClientMock := newMockBlueprintMaskInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(blueprintClientMock, maskClientMock, eventRecorderMock)

		// when
		persistenceContext := make(map[string]interface{})
		persistenceContext[blueprintSpecRepoContextKey] = 1
		err := repo.Update(ctx, &domain.BlueprintSpec{
			Id:                 blueprintId,
			Events:             nil,
			PersistenceContext: persistenceContext,
		})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "persistence context in blueprintSpec is not a 'blueprintSpecRepoContext' but 'int'")
	})

	t.Run("conflict error on status update", func(t *testing.T) {
		// given
		blueprintClientMock := newMockBlueprintInterface(t)
		maskClientMock := newMockBlueprintMaskInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(blueprintClientMock, maskClientMock, eventRecorderMock)
		expectedStatus := &bpv2.BlueprintStatus{
			EffectiveBlueprint: &bpv2.BlueprintManifest{
				Dogus:  []bpv2.Dogu{},
				Config: nil,
			},
			StateDiff:  &bpv2.StateDiff{DoguDiffs: map[string]bpv2.DoguDiff{}},
			Conditions: []metav1.Condition{},
		}
		expectedError := k8sErrors.NewConflict(
			schema.GroupResource{Group: "blueprints", Resource: blueprintId},
			blueprintId,
			fmt.Errorf("test-error"),
		)
		blueprintClientMock.EXPECT().
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
			Events:             nil,
			PersistenceContext: persistenceContext,
			Conditions:         []domain.Condition{},
		})

		// then
		require.Error(t, err)
		var expectedErrorType *domainservice.ConflictError
		assert.ErrorAs(t, err, &expectedErrorType)
		assert.ErrorIs(t, err, expectedError)
	})

	t.Run("internal error on status update", func(t *testing.T) {
		// given
		blueprintClientMock := newMockBlueprintInterface(t)
		maskClientMock := newMockBlueprintMaskInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(blueprintClientMock, maskClientMock, eventRecorderMock)
		expectedStatus := &bpv2.BlueprintStatus{
			EffectiveBlueprint: &bpv2.BlueprintManifest{
				Dogus:  []bpv2.Dogu{},
				Config: nil,
			},
			StateDiff:  &bpv2.StateDiff{DoguDiffs: map[string]bpv2.DoguDiff{}},
			Conditions: []metav1.Condition{},
		}
		expectedError := fmt.Errorf("test-error")
		blueprintClientMock.EXPECT().
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
			Events:             nil,
			PersistenceContext: persistenceContext,
			Conditions:         []domain.Condition{},
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
		blueprintClientMock := newMockBlueprintInterface(t)
		maskClientMock := newMockBlueprintMaskInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(blueprintClientMock, maskClientMock, eventRecorderMock)
		blueprintClientMock.EXPECT().
			UpdateStatus(ctx, mock.Anything, metav1.UpdateOptions{}).
			RunAndReturn(func(ctx2 context.Context, blueprint *bpv2.Blueprint, options metav1.UpdateOptions) (*bpv2.Blueprint, error) {
				// assert.Equal(t, &expected, blueprint)
				blueprint.ResourceVersion = "newVersion"
				return blueprint, nil
			})

		var events []domain.Event
		events = append(events,
			domain.StateDiffDeterminedEvent{},
			domain.BlueprintSpecInvalidEvent{ValidationError: errors.New("test-error")},
		)
		eventRecorderMock.EXPECT().Event(mock.Anything, corev1.EventTypeNormal, "StateDiffDetermined", "state diff determined:\n  0 config changes ()\n  0 dogu actions ()")
		eventRecorderMock.EXPECT().Event(mock.Anything, corev1.EventTypeNormal, "BlueprintSpecInvalid", "test-error")

		// when
		persistenceContext := make(map[string]interface{})
		persistenceContext[blueprintSpecRepoContextKey] = blueprintSpecRepoContext{"abc"}
		spec := &domain.BlueprintSpec{
			Id:                 blueprintId,
			Events:             events,
			PersistenceContext: persistenceContext,
			Conditions:         []domain.Condition{},
		}
		err := repo.Update(ctx, spec)

		// then
		require.NoError(t, err)
		newPersistenceContext, _ := getPersistenceContext(ctx, spec)
		assert.Equal(t, "newVersion", newPersistenceContext.resourceVersion)
		assert.Empty(t, spec.Events, "events in aggregate should be deleted after publishing them")
	})
}

func Test_blueprintSpecRepo_Count(t *testing.T) {
	t.Run("0 when no blueprints found", func(t *testing.T) {
		blueprintClientMock := newMockBlueprintInterface(t)
		maskClientMock := newMockBlueprintMaskInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(blueprintClientMock, maskClientMock, eventRecorderMock)
		limit := 100
		blueprintClientMock.EXPECT().List(ctx, metav1.ListOptions{Limit: int64(limit)}).Return(nil, nil)

		// when
		count, err := repo.Count(ctx, limit)

		// then
		assert.Equal(t, 0, count)
		require.NoError(t, err)
	})

	t.Run("1 on single blueprint resource", func(t *testing.T) {
		blueprintClientMock := newMockBlueprintInterface(t)
		maskClientMock := newMockBlueprintMaskInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(blueprintClientMock, maskClientMock, eventRecorderMock)
		list := &bpv2.BlueprintList{
			Items: []bpv2.Blueprint{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "blueprint-1",
					},
				},
			},
		}
		limit := 2
		blueprintClientMock.EXPECT().List(ctx, metav1.ListOptions{Limit: int64(limit)}).Return(list, nil)

		// when
		count, err := repo.Count(ctx, limit)

		// then
		assert.Equal(t, 1, count)
		require.NoError(t, err)
	})

	t.Run("2 on two blueprint resources", func(t *testing.T) {
		blueprintClientMock := newMockBlueprintInterface(t)
		maskClientMock := newMockBlueprintMaskInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(blueprintClientMock, maskClientMock, eventRecorderMock)
		list := &bpv2.BlueprintList{
			Items: []bpv2.Blueprint{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "blueprint-1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "blueprint-2",
					},
				},
			},
		}
		limit := 2
		blueprintClientMock.EXPECT().List(ctx, metav1.ListOptions{Limit: int64(limit)}).Return(list, nil)

		// when
		count, err := repo.Count(ctx, limit)

		// then
		assert.Equal(t, 2, count)
		require.NoError(t, err)
	})

	t.Run("InternalError on List error", func(t *testing.T) {
		blueprintClientMock := newMockBlueprintInterface(t)
		maskClientMock := newMockBlueprintMaskInterface(t)
		eventRecorderMock := newMockEventRecorder(t)
		repo := NewBlueprintSpecRepository(blueprintClientMock, maskClientMock, eventRecorderMock)
		limit := 2
		blueprintClientMock.EXPECT().List(ctx, metav1.ListOptions{Limit: 2}).Return(nil, assert.AnError)

		// when
		_, err := repo.Count(ctx, limit)

		// then
		require.Error(t, err)
		var targetErr *domainservice.InternalError
		assert.ErrorAs(t, err, &targetErr)
		assert.ErrorContains(t, err, "error while listing blueprint resources")
	})
}
