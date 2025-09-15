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
	blueprintClient blueprintInterface
	eventRecorder   eventRecorder
}

// NewBlueprintSpecRepository returns a new BlueprintSpecRepository to interact on BlueprintSpecs.
func NewBlueprintSpecRepository(
	blueprintClient bpv2client.BlueprintInterface,
	eventRecorder eventRecorder,
) domainservice.BlueprintSpecRepository {
	return &blueprintSpecRepo{
		blueprintClient: blueprintClient,
		eventRecorder:   eventRecorder,
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
			}
		}
		return nil, &domainservice.InternalError{
			WrappedError: err,
			Message:      fmt.Sprintf("error while loading blueprint CR %q", blueprintId),
		}
	}

	var effectiveBlueprint domain.EffectiveBlueprint
	var stateDiff domain.StateDiff
	if blueprintCR.Status != nil {
		effectiveBlueprint, err = serializerv2.ConvertToEffectiveBlueprintDomain(blueprintCR.Status.EffectiveBlueprint)
		if err != nil {
			return nil, err
		}

		stateDiff, err = serializerv2.ConvertToStateDiffDomain(blueprintCR.Status.StateDiff)
		if err != nil {
			return nil, err
		}
	}

	var conditions []domain.Condition
	if blueprintCR.Status != nil && blueprintCR.Status.Conditions != nil {
		conditions = blueprintCR.Status.Conditions
	}

	blueprintSpec := &domain.BlueprintSpec{
		Id:                 blueprintId,
		EffectiveBlueprint: effectiveBlueprint,
		StateDiff:          stateDiff,
		Conditions:         conditions,
		Config: domain.BlueprintConfiguration{
			IgnoreDoguHealth:         boolPtrToValue(blueprintCR.Spec.IgnoreDoguHealth),
			IgnoreComponentHealth:    boolPtrToValue(blueprintCR.Spec.IgnoreComponentHealth),
			AllowDoguNamespaceSwitch: boolPtrToValue(blueprintCR.Spec.AllowDoguNamespaceSwitch),
			DryRun:                   boolPtrToValue(blueprintCR.Spec.Stopped),
		},
	}

	blueprint, blueprintErr := serializerv2.ConvertToBlueprintDomain(blueprintCR.Spec.Blueprint)
	if blueprintErr != nil {
		blueprintErrorEvent := domain.BlueprintSpecInvalidEvent{ValidationError: blueprintErr}
		repo.eventRecorder.Event(blueprintCR, corev1.EventTypeWarning, blueprintErrorEvent.Name(), blueprintErrorEvent.Message())
	}

	blueprintMask, maskErr := serializerv2.ConvertToBlueprintMaskDomain(blueprintCR.Spec.BlueprintMask)
	if maskErr != nil {
		blueprintMaskErrorEvent := domain.BlueprintSpecInvalidEvent{ValidationError: maskErr}
		repo.eventRecorder.Event(blueprintCR, corev1.EventTypeWarning, blueprintMaskErrorEvent.Name(), blueprintMaskErrorEvent.Message())
	}

	serializationErr := errors.Join(blueprintErr, maskErr)
	if serializationErr != nil {
		return nil, fmt.Errorf("could not deserialize blueprint CR %q: %w", blueprintId, serializationErr)
	}

	setPersistenceContext(blueprintCR, blueprintSpec)
	blueprintSpec.Blueprint = blueprint
	blueprintSpec.BlueprintMask = blueprintMask
	return blueprintSpec, nil
}

func boolPtrToValue(boolPtr *bool) bool {
	if boolPtr != nil {
		return *boolPtr
	}
	return false
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
