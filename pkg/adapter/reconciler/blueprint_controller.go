package reconciler

import (
	"context"
	"errors"
	"time"

	v2 "github.com/cloudogu/k8s-dogu-lib/v2/api/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	bpv2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
)

// BlueprintReconciler reconciles a Blueprint object
type BlueprintReconciler struct {
	blueprintChangeHandler BlueprintChangeHandler
	blueprintRepo          BlueprintSpecRepository
	namespace              string
	debounce               SingletonDebounce
	window                 time.Duration
	errorHandler           *ErrorHandler
}

func NewBlueprintReconciler(
	blueprintChangeHandler BlueprintChangeHandler,
	repo domainservice.BlueprintSpecRepository,
	namespace string,
	window time.Duration,
) *BlueprintReconciler {
	return &BlueprintReconciler{
		blueprintChangeHandler: blueprintChangeHandler,
		blueprintRepo:          repo,
		namespace:              namespace,
		debounce:               SingletonDebounce{},
		window:                 window,
		errorHandler:           NewErrorHandler(),
	}
}

// +kubebuilder:rbac:groups=k8s.cloudogu.com,resources=blueprints,verbs=get;watch;update;patch
// +kubebuilder:rbac:groups=k8s.cloudogu.com,resources=blueprints/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=k8s.cloudogu.com,resources=blueprints/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *BlueprintReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).
		WithName("BlueprintReconciler.Reconcile").
		WithValues("resourceName", req.Name)

	err := r.blueprintChangeHandler.CheckForMultipleBlueprintResources(ctx)

	if err != nil {
		return r.errorHandler.handleError(logger, err)
	}

	err = r.blueprintChangeHandler.HandleUntilApplied(ctx, req.Name)

	if err != nil {
		return r.errorHandler.handleError(logger, err)
	}

	// Schedule a reconciliation after the cooldown period if there is one pending.
	if requeue, after := r.debounce.ShouldRequeue(); requeue {
		return ctrl.Result{RequeueAfter: after}, nil
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BlueprintReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if mgr == nil {
		return errors.New("must provide a non-nil Manager")
	}

	controllerOptions := mgr.GetControllerOptions()
	options := controller.TypedOptions[reconcile.Request]{
		SkipNameValidation: controllerOptions.SkipNameValidation,
		RecoverPanic:       controllerOptions.RecoverPanic,
		NeedLeaderElection: controllerOptions.NeedLeaderElection,
	}

	return ctrl.NewControllerManagedBy(mgr).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		WithOptions(options).
		For(&bpv2.Blueprint{}).
		WatchesRawSource(r.getConfigMapKind(mgr)).
		WatchesRawSource(r.getSecretKind(mgr)).
		WatchesRawSource(r.getDoguKind(mgr)).
		Complete(r)
}

func (r *BlueprintReconciler) getConfigMapKind(mgr ctrl.Manager) source.TypedSyncingSource[reconcile.Request] {
	return source.TypedKind(
		mgr.GetCache(),
		&corev1.ConfigMap{},
		handler.TypedEnqueueRequestsFromMapFunc(func(ctx context.Context, cm *corev1.ConfigMap) []reconcile.Request {
			return r.getBlueprintRequest(ctx)
		}),
		predicate.And(
			makeResourcePredicate[*corev1.ConfigMap](r.hasOperatorNamespace),
			makeResourcePredicate[*corev1.ConfigMap](hasCesLabel),
			makeResourcePredicate[*corev1.ConfigMap](hasNotDoguDescriptorLabel),
			makeContentPredicate(&r.debounce, r.window, configMapContentChanged),
		),
	)
}

func (r *BlueprintReconciler) getSecretKind(mgr ctrl.Manager) source.TypedSyncingSource[reconcile.Request] {
	return source.TypedKind(
		mgr.GetCache(),
		&corev1.Secret{},
		handler.TypedEnqueueRequestsFromMapFunc(func(ctx context.Context, s *corev1.Secret) []reconcile.Request {
			return r.getBlueprintRequest(ctx)
		}),
		predicate.And(
			makeResourcePredicate[*corev1.Secret](r.hasOperatorNamespace),
			makeResourcePredicate[*corev1.Secret](hasCesLabel),
			makeContentPredicate(&r.debounce, r.window, secretContentChanged),
		),
	)
}

func (r *BlueprintReconciler) getDoguKind(mgr ctrl.Manager) source.TypedSyncingSource[reconcile.Request] {
	return source.TypedKind(
		mgr.GetCache(),
		&v2.Dogu{},
		handler.TypedEnqueueRequestsFromMapFunc(func(ctx context.Context, d *v2.Dogu) []reconcile.Request {
			return r.getBlueprintRequest(ctx)
		}),
		predicate.And(
			makeResourcePredicate[*v2.Dogu](r.hasOperatorNamespace),
			makeContentPredicate(&r.debounce, r.window, doguSpecChanged),
		),
	)
}

func (r *BlueprintReconciler) getBlueprintRequest(ctx context.Context) []reconcile.Request {
	list, err := r.blueprintRepo.List(ctx)
	if err != nil || len(list.Items) != 1 {
		return nil
	}

	blueprintRequest := []reconcile.Request{{
		NamespacedName: types.NamespacedName{
			Name:      list.Items[0].Name,
			Namespace: list.Items[0].Namespace,
		},
	}}
	return blueprintRequest
}
