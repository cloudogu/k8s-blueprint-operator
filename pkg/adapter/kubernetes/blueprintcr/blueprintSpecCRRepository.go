package blueprintcr

import (
	"context"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/kubernetes"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/kubernetes/blueprintcr/v1"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

const blueprintSpecRepoContextKey = "blueprintSpecRepoContext"

type blueprintSpecRepoContext struct {
	resourceVersion string
}

type blueprintSpecRepo struct {
	blueprintClient         blueprintInterface
	blueprintSerializer     serializer.BlueprintSerializer
	blueprintMaskSerializer serializer.BlueprintMaskSerializer
	eventRecorder           eventRecorder
}

// NewBlueprintSpecRepository returns a new BlueprintSpecRepository to interact on BlueprintSpecs.
func NewBlueprintSpecRepository(
	blueprintClient kubernetes.BlueprintInterface,
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

	effectiveBlueprint, err := v1.ConvertToEffectiveBlueprintDomain(blueprintCR.Status.EffectiveBlueprint)
	if err != nil {
		return nil, err
	}

	stateDiff, err := v1.ConvertToStateDiffDomain(blueprintCR.Status.StateDiff)
	if err != nil {
		return nil, err
	}

	println("DEEEEEEEEEEEEEEEEEEEBUUUUUUUUUUUUUUUG GET")
	println("DEEEEEEEEEEEEEEEEEEEBUUUUUUUUUUUUUUUG")
	println("Effective")
	println(fmt.Sprintf("%+v", effectiveBlueprint))
	println("Status")
	println(fmt.Sprintf("%+v", stateDiff))

	blueprintSpec := &domain.BlueprintSpec{
		Id:                 blueprintId,
		EffectiveBlueprint: effectiveBlueprint,
		StateDiff:          stateDiff,
		Config: domain.BlueprintConfiguration{
			IgnoreDoguHealth:         blueprintCR.Spec.IgnoreDoguHealth,
			IgnoreComponentHealth:    blueprintCR.Spec.IgnoreComponentHealth,
			AllowDoguNamespaceSwitch: blueprintCR.Spec.AllowDoguNamespaceSwitch,
			DryRun:                   blueprintCR.Spec.DryRun,
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

	effectiveBlueprint, err := v1.ConvertToEffectiveBlueprintDTO(spec.EffectiveBlueprint)
	if err != nil {
		return err
	}

	blueprintJson, err := repo.blueprintSerializer.Serialize(spec.Blueprint)
	if err != nil {
		return domainservice.NewInternalError(err, "failed to serialize blueprint json")
	}

	blueprintMaskJson, err := repo.blueprintMaskSerializer.Serialize(spec.BlueprintMask)
	if err != nil {
		return domainservice.NewInternalError(err, "failed to serialize blueprint mask json")
	}

	updatedBlueprint := v1.Blueprint{
		ObjectMeta: metav1.ObjectMeta{
			Name:              spec.Id,
			ResourceVersion:   persistenceContext.resourceVersion,
			CreationTimestamp: metav1.Time{},
		},
		Spec: v1.BlueprintSpec{
			Blueprint:                blueprintJson,
			BlueprintMask:            blueprintMaskJson,
			IgnoreDoguHealth:         spec.Config.IgnoreDoguHealth,
			IgnoreComponentHealth:    spec.Config.IgnoreComponentHealth,
			AllowDoguNamespaceSwitch: spec.Config.AllowDoguNamespaceSwitch,
			DryRun:                   spec.Config.DryRun,
		},
	}

	logger.Info("update blueprint", "blueprint to save", updatedBlueprint)
	CRAfterUpdate, err := repo.blueprintClient.Update(ctx, &updatedBlueprint, metav1.UpdateOptions{})
	if err != nil {
		if k8sErrors.IsConflict(err) {
			return &domainservice.ConflictError{
				WrappedError: err,
				Message:      fmt.Sprintf("cannot update blueprint CR %q as it was modified in the meantime", spec.Id),
			}
		}
		return domainservice.NewInternalError(err, "cannot update blueprint CR %q", spec.Id)
	}

	blueprintStatus := v1.BlueprintStatus{
		Phase:              spec.Status,
		EffectiveBlueprint: effectiveBlueprint,
		StateDiff:          v1.ConvertToStateDiffDTO(spec.StateDiff),
	}

	println("DEEEEEEEEEEEEEEEEEEEBUUUUUUUUUUUUUUUG Update")
	println("DEEEEEEEEEEEEEEEEEEEBUUUUUUUUUUUUUUUG")
	println("Effective")
	println(fmt.Sprintf("%+v", blueprintStatus.EffectiveBlueprint))
	println("Status")
	println(fmt.Sprintf("%+v", blueprintStatus.StateDiff))

	CRAfterUpdate.Status = blueprintStatus
	CRAfterUpdate, err = repo.blueprintClient.UpdateStatus(ctx, CRAfterUpdate, metav1.UpdateOptions{})
	if err != nil {
		if k8sErrors.IsConflict(err) {
			return &domainservice.ConflictError{
				WrappedError: err,
				Message:      fmt.Sprintf("cannot update blueprint CR status %q as it was modified in the meantime", spec.Id),
			}
		}
		return domainservice.NewInternalError(err, "cannot update blueprint CR status %q", spec.Id)
	}

	setPersistenceContext(CRAfterUpdate, spec)
	repo.publishEvents(CRAfterUpdate, spec.Events)
	spec.Events = []domain.Event{}

	return nil
}

func setPersistenceContext(blueprintCR *v1.Blueprint, spec *domain.BlueprintSpec) {
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

func (repo *blueprintSpecRepo) publishEvents(blueprintCR *v1.Blueprint, events []domain.Event) {
	for _, event := range events {
		repo.eventRecorder.Event(blueprintCR, corev1.EventTypeNormal, event.Name(), event.Message())
	}
}
