package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer/effectiveBlueprintV1"
	v1 "github.com/cloudogu/k8s-blueprint-operator/pkg/api/v1"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const resourceVersionKey = "resourceVersion"

type resourceVersionValue struct {
	string
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
func (repo *blueprintSpecRepo) GetById(ctx context.Context, blueprintId string) (domain.BlueprintSpec, error) {
	blueprintCR, err := repo.blueprintClient.Get(ctx, blueprintId, metav1.GetOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return domain.BlueprintSpec{}, &domainservice.NotFoundError{
				WrappedError: err,
				Message:      fmt.Sprintf("cannot load blueprint CR %q as it does not exist", blueprintId),
			}
		}
		return domain.BlueprintSpec{}, &domainservice.InternalError{
			WrappedError: err,
			Message:      fmt.Sprintf("error while loading blueprint CR %q", blueprintId),
		}
	}

	persistenceContext := make(map[string]interface{}, 1)
	persistenceContext[resourceVersionKey] = resourceVersionValue{blueprintCR.GetResourceVersion()}
	effectiveBlueprint, err := effectiveBlueprintV1.ConvertToEffectiveBlueprint(blueprintCR.Status.EffectiveBlueprint)
	if err != nil {
		return domain.BlueprintSpec{}, err
	}
	blueprintSpec := domain.BlueprintSpec{
		Id:                   blueprintId,
		EffectiveBlueprint:   effectiveBlueprint,
		StateDiff:            domain.StateDiff{},
		BlueprintUpgradePlan: domain.BlueprintUpgradePlan{},
		Config: domain.BlueprintConfiguration{
			IgnoreDoguHealth:         blueprintCR.Spec.IgnoreDoguHealth,
			AllowDoguNamespaceSwitch: blueprintCR.Spec.AllowDoguNamespaceSwitch,
		},
		Status:             blueprintCR.Status.Phase,
		PersistenceContext: persistenceContext,
	}

	blueprint, blueprintErr := repo.blueprintSerializer.Deserialize(blueprintCR.Spec.Blueprint)
	blueprintMask, maskErr := repo.blueprintMaskSerializer.Deserialize(blueprintCR.Spec.BlueprintMask)
	serializationErr := errors.Join(blueprintErr, maskErr)
	if serializationErr != nil {
		return blueprintSpec, fmt.Errorf("could not deserialize blueprint CR %q: %w", blueprintId, serializationErr)
	}

	blueprintSpec.Blueprint = blueprint
	blueprintSpec.BlueprintMask = blueprintMask
	return blueprintSpec, nil
}

// Update persists changes in the blueprint to the corresponding blueprint CR.
func (repo *blueprintSpecRepo) Update(ctx context.Context, spec domain.BlueprintSpec) error {
	logger := log.FromContext(ctx).WithName("blueprintSpecRepo.Update")
	resourceVersion, err := getResourceVersion(ctx, spec)
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
			ResourceVersion:   resourceVersion.string,
			CreationTimestamp: metav1.Time{},
		},
		Status: v1.BlueprintStatus{
			Phase:              spec.Status,
			EffectiveBlueprint: effectiveBlueprint,
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
	repo.publishEvents(CRAfterUpdate, spec.Events)

	return nil
}

// getResourceVersion reads the repo-specific resourceVersion from the domain.BlueprintSpec or returns an error.
func getResourceVersion(ctx context.Context, spec domain.BlueprintSpec) (resourceVersionValue, error) {
	logger := log.FromContext(ctx).WithName("blueprintSpecRepo.Update")
	rawResourceVersion, versionExists := spec.PersistenceContext[resourceVersionKey]
	if versionExists {
		resourceVersion, isString := rawResourceVersion.(resourceVersionValue)
		if isString {
			return resourceVersion, nil
		} else {
			err := fmt.Errorf("resourceVersion in blueprintSpec is not a 'resourceVersionValue' but '%T'", rawResourceVersion)
			logger.Error(err, "does this value come from a different repository?")
			return resourceVersionValue{}, err
		}
	} else {
		err := errors.New("no resourceVersion was provided over the persistenceContext in the given blueprintSpec")
		logger.Error(err, "This is normally written while loading the blueprintSpec over this repository. "+
			"Did you try to persist a new blueprintSpec with repo.Update()?")
		return resourceVersionValue{}, err
	}
}

func (repo *blueprintSpecRepo) publishEvents(blueprintCR *v1.Blueprint, events []interface{}) {
	for _, event := range events {
		switch ev := event.(type) {
		case domain.BlueprintSpecStaticallyValidatedEvent:
			repo.eventRecorder.Event(blueprintCR, corev1.EventTypeNormal, "BlueprintSpecStaticallyValidatedEvent", "")
		case domain.BlueprintSpecValidatedEvent:
			repo.eventRecorder.Event(blueprintCR, corev1.EventTypeNormal, "BlueprintSpecValidatedEvent", "")
		case domain.BlueprintSpecInvalidEvent:
			repo.eventRecorder.Event(blueprintCR, corev1.EventTypeNormal, "BlueprintSpecInvalidEvent", ev.ValidationError.Error())
		case domain.EffectiveBlueprintCalculatedEvent:
			repo.eventRecorder.Event(blueprintCR, corev1.EventTypeNormal, "EffectiveBlueprintCalculatedEvent", fmt.Sprintf("effective blueprint: %+v", ev.EffectiveBlueprint))
		default:
			repo.eventRecorder.Event(blueprintCR, corev1.EventTypeNormal, "Unknown", fmt.Sprintf("unknown event of type '%T': %+v", event, event))
		}
	}
}
