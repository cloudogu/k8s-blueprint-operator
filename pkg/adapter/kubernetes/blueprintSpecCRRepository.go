package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer/effectiveBlueprintV1"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer/stateDiffV1"
	v1 "github.com/cloudogu/k8s-blueprint-operator/pkg/api/v1"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const blueprintSpecRepoContextKey = "blueprintSpecRepoContext"

type blueprintSpecRepoContext struct {
	resourceVersion string
}

type blueprintSpecRepo struct {
	blueprintClient         BlueprintInterface
	blueprintSerializer     serializer.BlueprintSerializer
	blueprintMaskSerializer serializer.BlueprintMaskSerializer
	eventRecorder           eventRecorder
}

// NewBlueprintSpecRepository returns a new BlueprintSpecRepository to interact on BlueprintSpecs.
func NewBlueprintSpecRepository(
	blueprintClient BlueprintInterface,
	blueprintSerializer serializer.BlueprintSerializer,
	blueprintMaskSerializer serializer.BlueprintMaskSerializer,
	eventRecorder eventRecorder,
) domainservice.BlueprintSpecRepository {
	return &blueprintSpecRepo{
		blueprintClient:         blueprintClient,
		blueprintSerializer:     blueprintSerializer,
		blueprintMaskSerializer: blueprintMaskSerializer,
		eventRecorder:           eventRecorder,
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

	effectiveBlueprint, err := effectiveBlueprintV1.ConvertToEffectiveBlueprint(blueprintCR.Status.EffectiveBlueprint)
	if err != nil {
		return nil, err
	}

	stateDiff, err := stateDiffV1.ConvertToDomainModel(blueprintCR.Status.StateDiff)
	if err != nil {
		return nil, err
	}

	blueprintSpec := &domain.BlueprintSpec{
		Id:                   blueprintId,
		EffectiveBlueprint:   effectiveBlueprint,
		StateDiff:            stateDiff,
		BlueprintUpgradePlan: domain.BlueprintUpgradePlan{},
		Config: domain.BlueprintConfiguration{
			IgnoreDoguHealth:         blueprintCR.Spec.IgnoreDoguHealth,
			AllowDoguNamespaceSwitch: blueprintCR.Spec.AllowDoguNamespaceSwitch,
		},
		Status: blueprintCR.Status.Phase,
	}

	blueprint, blueprintErr := repo.blueprintSerializer.Deserialize(blueprintCR.Spec.Blueprint)
	blueprintMask, maskErr := repo.blueprintMaskSerializer.Deserialize(blueprintCR.Spec.BlueprintMask)
	serializationErr := errors.Join(blueprintErr, maskErr)
	if serializationErr != nil {
		return nil, fmt.Errorf("could not deserialize blueprint CR %q: %w", blueprintId, serializationErr)
	}

	setPersistenceContext(blueprintCR, blueprintSpec)
	blueprintSpec.Blueprint = blueprint
	blueprintSpec.BlueprintMask = blueprintMask
	return blueprintSpec, nil
}

// Update persists changes in the blueprint to the corresponding blueprint CR.
func (repo *blueprintSpecRepo) Update(ctx context.Context, spec *domain.BlueprintSpec) error {
	logger := log.FromContext(ctx).WithName("blueprintSpecRepo.Update")
	persistenceContext, err := getPersistenceContext(ctx, spec)
	if err != nil {
		return err
	}

	effectiveBlueprint, err := effectiveBlueprintV1.ConvertToEffectiveBlueprintV1(spec.EffectiveBlueprint)
	if err != nil {
		return err
	}

	updatedBlueprint := v1.Blueprint{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:              spec.Id,
			ResourceVersion:   persistenceContext.resourceVersion,
			CreationTimestamp: metav1.Time{},
		},
		Status: v1.BlueprintStatus{
			Phase:              spec.Status,
			EffectiveBlueprint: effectiveBlueprint,
			StateDiff:          stateDiffV1.ConvertToDTO(spec.StateDiff),
		},
	}

	logger.Info("update blueprint", "blueprint to save", updatedBlueprint)
	CRAfterUpdate, err := repo.blueprintClient.UpdateStatus(ctx, &updatedBlueprint, metav1.UpdateOptions{})
	if err != nil {
		if k8sErrors.IsConflict(err) {
			return &domainservice.ConflictError{
				WrappedError: err,
				Message:      fmt.Sprintf("cannot update blueprint CR %q as it was modified in the meantime", spec.Id),
			}
		}
		return &domainservice.InternalError{WrappedError: err, Message: fmt.Sprintf("Cannot update blueprint CR %q", spec.Id)}
	}

	setPersistenceContext(CRAfterUpdate, spec)
	repo.publishEvents(CRAfterUpdate, spec.Events)
	spec.Events = []domain.Event{}

	return nil
}

func setPersistenceContext(blueprintCR *v1.Blueprint, spec *domain.BlueprintSpec) {
	persistenceContext := spec.PersistenceContext
	if persistenceContext == nil {
		persistenceContext = make(map[string]interface{}, 1)
	}
	persistenceContext[blueprintSpecRepoContextKey] = blueprintSpecRepoContext{
		resourceVersion: blueprintCR.GetResourceVersion(),
	}
	spec.PersistenceContext = persistenceContext
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

func (repo *blueprintSpecRepo) publishEvents(blueprintCR *v1.Blueprint, events []domain.Event) {
	for _, event := range events {
		repo.eventRecorder.Event(blueprintCR, corev1.EventTypeNormal, event.Name(), event.Message())
	}
}
