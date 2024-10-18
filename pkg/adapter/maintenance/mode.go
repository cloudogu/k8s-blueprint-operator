package maintenance

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	liberrors "github.com/cloudogu/k8s-registry-lib/errors"
	"github.com/cloudogu/k8s-registry-lib/repository"
)

// Mode contains methods to Activate and Deactivate the switcher.
// When it is active, a page is displayed to the user, telling them that there is maintenance going on.
type Mode struct {
	libAdapter libMaintenanceModeAdapter
}

func NewMaintenanceModeAdapter(libAdapter libMaintenanceModeAdapter) *Mode {
	return &Mode{
		libAdapter: libAdapter,
	}
}

// Activate enables the maintenance mode, setting the given MaintenancePageModel
func (m *Mode) Activate(ctx context.Context, title, text string) error {
	err := m.libAdapter.Activate(ctx, repository.MaintenanceModeDescription{
		Title: title,
		Text:  text,
	})
	err = mapToBlueprintError(err)
	if err != nil {
		return fmt.Errorf("could not activate maintenance mode: %w", err)
	}
	return nil
}

// Deactivate disables the maintenance mode.
func (m *Mode) Deactivate(ctx context.Context) error {
	err := m.libAdapter.Deactivate(ctx)
	err = mapToBlueprintError(err)
	if err != nil {
		return fmt.Errorf("could not activate maintenance mode: %w", err)
	}
	return nil
}

func mapToBlueprintError(err error) error {
	if err != nil {
		if liberrors.IsConflictError(err) {
			return domainservice.NewConflictError(err, "there were conflicting changes to the maintenance mode")
		} else if liberrors.IsConnectionError(err) {
			return domainservice.NewInternalError(err, "could not update maintenance mode due to connection problems")
		} else {
			// GenericError and NotFoundError and fallback if even that would not match the error
			return domainservice.NewInternalError(err, "could not update maintenance mode due to an unknown problem")
		}
	}
	return nil
}
