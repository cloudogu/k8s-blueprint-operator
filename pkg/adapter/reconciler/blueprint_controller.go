package reconciler

import (
	"context"
	"errors"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"

	k8sv1 "github.com/cloudogu/k8s-blueprint-operator/pkg/api/v1"
)

// BlueprintReconciler reconciles a Blueprint object
type BlueprintReconciler struct {
	blueprintChangeHandler BlueprintChangeHandler
}

func NewBlueprintReconciler(
	blueprintChangeHandler BlueprintChangeHandler,
) *BlueprintReconciler {
	return &BlueprintReconciler{blueprintChangeHandler: blueprintChangeHandler}
}

//+kubebuilder:rbac:groups=k8s.cloudogu.com,resources=blueprints,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.cloudogu.com,resources=blueprints/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.cloudogu.com,resources=blueprints/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *BlueprintReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).
		WithName("BlueprintReconciler.Reconcile").
		WithValues("resourceName", req.Name)

	err := r.blueprintChangeHandler.HandleChange(ctx, req.Name)

	var internalError *domainservice.InternalError
	var conflictError *domainservice.ConflictError
	var notFoundError *domainservice.NotFoundError
	var invalidError *domain.InvalidBlueprintError
	if err != nil {
		logger := logger.WithValues("error", err)
		if errors.As(err, &internalError) {
			logger.Error(err, "An internal error occurred and can maybe be fixed by retrying it later")
			return ctrl.Result{}, err //automatic requeue
		} else if errors.As(err, &conflictError) {
			logger.Info("A concurrent update happened in conflict to the processing of the blueprint spec. A retry could fix this issue")
			return ctrl.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil // no error as this would lead to the ignorance of our own retry params
		} else if errors.As(err, &notFoundError) {
			logger.Info("blueprint was not found, so maybe it was deleted in the meantime. No further evaluation will happen")
			return ctrl.Result{Requeue: false}, nil
		} else if errors.As(err, &invalidError) {
			logger.Info("blueprint is invalid, therefore there will be no further evaluation.")
			return ctrl.Result{Requeue: false}, nil
		}
		logger.Error(err, "an unknown error type occurred. Retry with default backoff")
		return ctrl.Result{}, err //automatic requeue
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BlueprintReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8sv1.Blueprint{}).
		Complete(r)
}
