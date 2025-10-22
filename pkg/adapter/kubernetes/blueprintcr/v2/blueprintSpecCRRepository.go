package v2

import (
	"context"
	"errors"
	"fmt"

	v2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	serializerv2 "github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/kubernetes/blueprintcr/v2/serializer"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/log"

	bpv2client "github.com/cloudogu/k8s-blueprint-lib/v2/client"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
)

const blueprintSpecRepoContextKey = "blueprintSpecRepoContext"

type blueprintSpecRepoContext struct {
	resourceVersion string
}

type blueprintSpecRepo struct {
	blueprintClient     blueprintInterface
	blueprintMaskClient blueprintMaskInterface
	eventRecorder       eventRecorder
}

// NewBlueprintSpecRepository returns a new BlueprintSpecRepository to interact on BlueprintSpecs.
func NewBlueprintSpecRepository(
	blueprintClient bpv2client.BlueprintInterface,
	blueprintMaskClient bpv2client.BlueprintMaskInterface,
	eventRecorder eventRecorder,
) domainservice.BlueprintSpecRepository {
	return &blueprintSpecRepo{
		blueprintClient:     blueprintClient,
		blueprintMaskClient: blueprintMaskClient,
		eventRecorder:       eventRecorder,
	}
}

// GetById returns a Blueprint identified by its ID.
func (repo *blueprintSpecRepo) GetById(ctx context.Context, blueprintId string) (*domain.BlueprintSpec, error) {
	blueprintCR, err := repo.blueprintClient.Get(ctx, blueprintId, metav1.GetOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return nil, &domainservice.NotFoundError{
				WrappedError: err,
				Message:      fmt.Sprintf("cannot load blueprint CR %q as it does not exist", blueprintId),
				DoNotRetry:   true,
			}
		}
		return nil, &domainservice.InternalError{
			WrappedError: err,
			Message:      fmt.Sprintf("error while loading blueprint CR %q", blueprintId),
		}
	}

	effectiveBlueprint, err := serializerv2.ConvertBlueprintStatus(blueprintCR)
	if err != nil {
		return nil, err
	}

	var conditions []domain.Condition
	if blueprintCR.Status != nil && blueprintCR.Status.Conditions != nil {
		conditions = blueprintCR.Status.Conditions
	}

	blueprintSpec := &domain.BlueprintSpec{
		Id:                 blueprintId,
		DisplayName:        blueprintCR.Spec.DisplayName,
		EffectiveBlueprint: effectiveBlueprint,
		Conditions:         conditions,
		Config: domain.BlueprintConfiguration{
			IgnoreDoguHealth:         ptr.Deref(blueprintCR.Spec.IgnoreDoguHealth, false),
			AllowDoguNamespaceSwitch: ptr.Deref(blueprintCR.Spec.AllowDoguNamespaceSwitch, false),
			Stopped:                  ptr.Deref(blueprintCR.Spec.Stopped, false),
		},
	}

	maskManifest, err := repo.getMaskManifest(ctx, blueprintId, blueprintCR)
	if err != nil {
		return nil, err
	}

	err = serializerv2.SerializeBlueprintAndMask(blueprintSpec, blueprintCR.Spec.Blueprint, maskManifest)
	if err != nil {
		invalidErrorEvent := domain.BlueprintSpecInvalidEvent{ValidationError: err}
		repo.eventRecorder.Event(blueprintCR, corev1.EventTypeWarning, invalidErrorEvent.Name(), invalidErrorEvent.Message())
		return nil, fmt.Errorf("could not deserialize blueprint CR %q: %w", blueprintId, err)
	}

	setPersistenceContext(blueprintCR, blueprintSpec)
	return blueprintSpec, nil
}

func (repo *blueprintSpecRepo) getMaskManifest(ctx context.Context, blueprintId string, blueprintCR *v2.Blueprint) (*v2.BlueprintMaskManifest, error) {
	if blueprintCR.Spec.BlueprintMask != nil && blueprintCR.Spec.BlueprintMaskRef != nil {
		err := &domain.InvalidBlueprintError{Message: "blueprint mask and mask ref cannot be set at the same time"}
		invalidErrorEvent := domain.BlueprintSpecInvalidEvent{ValidationError: err}
		repo.eventRecorder.Event(blueprintCR, corev1.EventTypeWarning, invalidErrorEvent.Name(), invalidErrorEvent.Message())
		return nil, fmt.Errorf("could not deserialize blueprint CR %q: %w", blueprintId, err)
	}

	var maskManifest = blueprintCR.Spec.BlueprintMask
	if blueprintCR.Spec.BlueprintMaskRef != nil {
		blueprintMask, maskErr := repo.blueprintMaskClient.Get(ctx, blueprintCR.Spec.BlueprintMaskRef.Name, metav1.GetOptions{})
		if maskErr != nil {
			return nil, &domainservice.InternalError{WrappedError: maskErr, Message: fmt.Sprintf("could not get blueprint mask from ref %q in blueprint %q", blueprintCR.Spec.BlueprintMaskRef.Name, blueprintId)}
		}

		maskManifest = blueprintMask.Spec.BlueprintMaskManifest
		maskRefEvent := domain.BlueprintMaskFromRefEvent{MaskRef: blueprintCR.Spec.BlueprintMaskRef.Name}
		repo.eventRecorder.Event(blueprintCR, corev1.EventTypeNormal, maskRefEvent.Name(), maskRefEvent.Message())
	}
	return maskManifest, nil
}

func (repo *blueprintSpecRepo) Count(ctx context.Context, limit int) (int, error) {
	limit64 := int64(limit)

	list, err := repo.blueprintClient.List(ctx, metav1.ListOptions{Limit: limit64})
	if err != nil {
		return 0, &domainservice.InternalError{
			WrappedError: err,
			Message:      "error while listing blueprint resources",
		}
	}

	if list == nil {
		return 0, nil
	}

	return len(list.Items), nil
}

// Update persists changes in the blueprint to the corresponding blueprint CR.
func (repo *blueprintSpecRepo) Update(ctx context.Context, spec *domain.BlueprintSpec) error {
	logger := log.FromContext(ctx).WithName("blueprintSpecRepo.Update")

	persistenceContext, err := getPersistenceContext(ctx, spec)
	if err != nil {
		return err
	}

	effectiveBlueprint := serializerv2.ConvertToBlueprintDTO(spec.EffectiveBlueprint)

	updatedBlueprint := &v2.Blueprint{
		ObjectMeta: metav1.ObjectMeta{
			Name:              spec.Id,
			ResourceVersion:   persistenceContext.resourceVersion,
			CreationTimestamp: metav1.Time{},
		},
		Status: &v2.BlueprintStatus{
			EffectiveBlueprint: &effectiveBlueprint,
			StateDiff:          serializerv2.ConvertToStateDiffDTO(spec.StateDiff),
			Conditions:         spec.Conditions,
		},
	}

	logger.V(2).Info("update blueprint CR status")

	CRAfterUpdate, err := repo.blueprintClient.UpdateStatus(ctx, updatedBlueprint, metav1.UpdateOptions{})
	if err != nil {
		if k8sErrors.IsConflict(err) {
			return domainservice.NewConflictError(err, "cannot update blueprint CR status %q as it was modified in the meantime", spec.Id)
		}
		return domainservice.NewInternalError(err, "cannot update blueprint CR status %q", spec.Id)
	}

	setPersistenceContext(CRAfterUpdate, spec)
	repo.publishEvents(CRAfterUpdate, spec.Events)
	spec.Events = []domain.Event{}

	return nil
}

func setPersistenceContext(blueprintCR *v2.Blueprint, spec *domain.BlueprintSpec) {
	if spec.PersistenceContext == nil {
		spec.PersistenceContext = make(map[string]interface{}, 1)
	}
	spec.PersistenceContext[blueprintSpecRepoContextKey] = blueprintSpecRepoContext{
		resourceVersion: blueprintCR.GetResourceVersion(),
	}
}

// getPersistenceContext reads the repo-specific resourceVersion from the domain.BlueprintSpec or returns an error.
func getPersistenceContext(ctx context.Context, spec *domain.BlueprintSpec) (blueprintSpecRepoContext, error) {
	logger := log.FromContext(ctx).WithName("blueprintSpecRepo.Update")
	rawField, versionExists := spec.PersistenceContext[blueprintSpecRepoContextKey]
	if versionExists {
		repoContext, isContext := rawField.(blueprintSpecRepoContext)
		if isContext {
			return repoContext, nil
		} else {
			err := fmt.Errorf("persistence context in blueprintSpec is not a 'blueprintSpecRepoContext' but '%T'", rawField)
			logger.Error(err, "does this value come from a different repository?")
			return blueprintSpecRepoContext{}, err
		}
	} else {
		err := errors.New("no blueprintSpecRepoContext was provided over the persistenceContext in the given blueprintSpec")
		logger.Error(err, "This is normally written while loading the blueprintSpec over this repository. "+
			"Did you try to persist a new blueprintSpec with repo.Update()?")
		return blueprintSpecRepoContext{}, err
	}
}

func (repo *blueprintSpecRepo) publishEvents(blueprintCR *v2.Blueprint, events []domain.Event) {
	for _, event := range events {
		repo.eventRecorder.Event(blueprintCR, corev1.EventTypeNormal, event.Name(), event.Message())
	}
}
