package reconciler

import (
	"context"
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	bpv2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
)

// BlueprintMaskReconciler recognizes changes from BlueprintMask objects and trigger reconciliation of the blueprint.
type BlueprintMaskReconciler struct {
	blueprintInterface     blueprintInterface
	blueprintMaskInterface blueprintMaskInterface
	blueprintEvents        chan<- event.TypedGenericEvent[*bpv2.Blueprint]
}

func NewBlueprintMaskReconciler(blueprintInterface blueprintInterface, blueprintMaskInterface blueprintMaskInterface, blueprintEvents chan<- event.TypedGenericEvent[*bpv2.Blueprint]) *BlueprintMaskReconciler {
	return &BlueprintMaskReconciler{
		blueprintInterface:     blueprintInterface,
		blueprintMaskInterface: blueprintMaskInterface,
		blueprintEvents:        blueprintEvents,
	}
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *BlueprintMaskReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).
		WithName("BlueprintMaskReconciler.Reconcile").
		WithValues("resourceName", req.Name)

	maskList, err := r.blueprintMaskInterface.List(ctx, metav1.ListOptions{})
	if err != nil {
		return ctrl.Result{}, err
	}

	if len(maskList.Items) == 0 {
		return ctrl.Result{}, nil
	}

	if len(maskList.Items) > 1 {
		logger.Info("found multiple blueprint mask resources")
	}

	blueprintList, err := r.blueprintInterface.List(ctx, metav1.ListOptions{})
	if err != nil {
		return ctrl.Result{}, err
	}

	for _, mask := range maskList.Items {
		for _, blueprint := range blueprintList.Items {
			ref := blueprint.Spec.BlueprintMaskRef
			if ref != nil && *ref == mask.Name {
				r.blueprintEvents <- event.TypedGenericEvent[*bpv2.Blueprint]{Object: &blueprint}
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BlueprintMaskReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if mgr == nil {
		return errors.New("must provide a non-nil Manager")
	}

	// TODO Is this necessary?
	controllerOptions := mgr.GetControllerOptions()
	options := controller.TypedOptions[reconcile.Request]{
		SkipNameValidation: controllerOptions.SkipNameValidation,
		RecoverPanic:       controllerOptions.RecoverPanic,
		NeedLeaderElection: controllerOptions.NeedLeaderElection,
	}
	return ctrl.NewControllerManagedBy(mgr).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		WithOptions(options).
		For(&bpv2.BlueprintMask{}).
		Complete(r)
}
