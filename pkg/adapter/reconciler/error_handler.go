package reconciler

import (
	"errors"
	"fmt"
	"time"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
)

// ErrorHandler handles different types of errors and determines the appropriate requeue strategy.
type ErrorHandler struct{}

// NewErrorHandler creates a new ErrorHandler instance.
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// handleError processes an error and returns the appropriate reconcile result.
func (h *ErrorHandler) handleError(logger logr.Logger, err error) (ctrl.Result, error) {
	errLogger := logger.WithValues("error", err)

	var internalError *domainservice.InternalError
	var conflictError *domainservice.ConflictError
	var notFoundError *domainservice.NotFoundError
	var invalidBlueprintError *domain.InvalidBlueprintError
	var healthError *domain.UnhealthyEcosystemError
	var stateDiffNotEmptyError *domain.StateDiffNotEmptyError
	var multipleBlueprintsError *domain.MultipleBlueprintsError
	var dogusNotUpToDateError *domain.DogusNotUpToDateError
	var restoreInProgressError *domain.RestoreInProgressError
	switch {
	case errors.As(err, &internalError):
		return h.handleInternalError(errLogger, err)
	case errors.As(err, &conflictError):
		return h.handleConflictError(errLogger)
	case errors.As(err, &notFoundError):
		return h.handleNotFoundError(errLogger, notFoundError)
	case errors.As(err, &invalidBlueprintError):
		return h.handleInvalidBlueprintError(errLogger)
	case errors.As(err, &healthError):
		return h.handleHealthError(errLogger)
	case errors.As(err, &stateDiffNotEmptyError):
		return h.handleStateDiffNotEmptyError(errLogger)
	case errors.As(err, &multipleBlueprintsError):
		return h.handleMultipleBlueprintsError(errLogger, err)
	case errors.As(err, &dogusNotUpToDateError):
		return h.handleDogusNotUpToDateError(errLogger, err)
	case errors.As(err, &restoreInProgressError):
		return h.handleRestoreInProgressError(errLogger, err)
	default:
		return h.handleUnknownError(errLogger, err)
	}
}

func (h *ErrorHandler) handleInternalError(logger logr.Logger, err error) (ctrl.Result, error) {
	logger.Error(err, "An internal error occurred and can maybe be fixed by retrying it later")
	return ctrl.Result{}, err // automatic requeue because of non-nil err
}

func (h *ErrorHandler) handleConflictError(logger logr.Logger) (ctrl.Result, error) {
	logger.Info("A concurrent update happened in conflict to the processing of the blueprint spec. A retry could fix this issue")
	return ctrl.Result{RequeueAfter: 1 * time.Second}, nil // no error as this would lead to the ignorance of our own retry params
}

func (h *ErrorHandler) handleNotFoundError(logger logr.Logger, err *domainservice.NotFoundError) (ctrl.Result, error) {
	if err.DoNotRetry {
		// do not retry in this case, because if f.e. the blueprint is not found, nothing will bring it back, except the
		// user, and this would trigger the reconciler by itself.
		logger.Error(err, "Did not find resource and a retry is not expected to fix this issue. There will be no further automatic evaluation.")
		return ctrl.Result{}, nil
	}
	logger.Error(err, "Resource was not found, so maybe it was deleted in the meantime. Retry later")
	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
}

func (h *ErrorHandler) handleInvalidBlueprintError(logger logr.Logger) (ctrl.Result, error) {
	logger.Info("Blueprint is invalid, therefore there will be no further evaluation.")
	return ctrl.Result{}, nil
}

func (h *ErrorHandler) handleHealthError(logger logr.Logger) (ctrl.Result, error) {
	// really normal case
	logger.Info("Ecosystem is unhealthy. Retry later")
	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
}

func (h *ErrorHandler) handleStateDiffNotEmptyError(logger logr.Logger) (ctrl.Result, error) {
	logger.Info("requeue until state diff is empty")
	// fast requeue here since state diff has to be determined again
	return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
}

func (h *ErrorHandler) handleMultipleBlueprintsError(logger logr.Logger, err error) (ctrl.Result, error) {
	logger.Error(err, "Ecosystem contains multiple blueprints - delete all but one. Retry later")
	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
}

func (h *ErrorHandler) handleDogusNotUpToDateError(logger logr.Logger, err error) (ctrl.Result, error) {
	// really normal case
	logger.Info(fmt.Sprintf("Dogus are not up to date yet. Retry later: %s", err.Error()))
	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
}

func (h *ErrorHandler) handleRestoreInProgressError(logger logr.Logger, err error) (ctrl.Result, error) {
	// really normal case
	logger.Info(fmt.Sprintf("A restore is currently in progress. Retry later: %s", err.Error()))
	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
}

func (h *ErrorHandler) handleUnknownError(logger logr.Logger, err error) (ctrl.Result, error) {
	logger.Error(err, "An unknown error type occurred. Retry with default backoff")
	return ctrl.Result{}, err // automatic requeue because of non-nil err
}
