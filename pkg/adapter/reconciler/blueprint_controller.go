package reconciler

import (
	"context"
	"errors"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/application"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
	"time"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	bpv2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
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

// +kubebuilder:rbac:groups=k8s.cloudogu.com,resources=blueprints,verbs=get;list;watch;update;patch
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

	err := r.blueprintChangeHandler.HandleUntilApplied(ctx, req.Name)

	if err != nil {
		return decideRequeueForError(logger, err)
	}

	return ctrl.Result{}, nil
}

func decideRequeueForError(logger logr.Logger, err error) (ctrl.Result, error) {
	errLogger := logger.WithValues("error", err)

	var internalError *domainservice.InternalError
	var conflictError *domainservice.ConflictError
	var notFoundError *domainservice.NotFoundError
	var invalidBlueprintError *domain.InvalidBlueprintError
	var healthError *domain.UnhealthyEcosystemError
	switch {
	case errors.As(err, &internalError):
		errLogger.Error(err, "An internal error occurred and can maybe be fixed by retrying it later")
		return ctrl.Result{}, err // automatic requeue because of non-nil err
	case errors.As(err, &conflictError):
		errLogger.Info("A concurrent update happened in conflict to the processing of the blueprint spec. A retry could fix this issue")
		return ctrl.Result{RequeueAfter: 1 * time.Second}, nil // no error as this would lead to the ignorance of our own retry params
	case errors.As(err, &notFoundError):
		if strings.Contains(err.Error(), application.REFERENCED_CONFIG_NOT_FOUND) {
			// retry in this case because maybe the user will create the secret in a few seconds.
			// we want to be declarative, so our API should not care about the order
			errLogger.Error(err, "Referenced config not found. Retry later")
			return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
		}
		errLogger.Error(err, "Blueprint was not found, so maybe it was deleted in the meantime. No further evaluation will happen")
		return ctrl.Result{}, nil
	case errors.As(err, &invalidBlueprintError):
		errLogger.Info("Blueprint is invalid, therefore there will be no further evaluation.")
		return ctrl.Result{}, nil
	case errors.As(err, &healthError):
		errLogger.Info("Ecosystem is unhealthy. Retry later")
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	default:
		errLogger.Error(err, "An unknown error type occurred. Retry with default backoff")
		return ctrl.Result{}, err // automatic requeue because of non-nil err
	}
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
		Complete(r)
}
