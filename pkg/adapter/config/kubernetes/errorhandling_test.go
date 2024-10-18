package kubernetes

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	liberrors "github.com/cloudogu/k8s-registry-lib/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mapToBlueprintError(t *testing.T) {
	tests := []struct {
		name          string
		givenError    error
		errTypeAssert func(err error) bool
	}{
		{"no error", nil, nil},
		{"connectionError", liberrors.NewConnectionError(assert.AnError), domainservice.IsInternalError},
		{"notFoundError", liberrors.NewNotFoundError(assert.AnError), domainservice.IsNotFoundError},
		{"conflictError", liberrors.NewConflictError(assert.AnError), domainservice.IsConflictError},
		{"alreadyExistsError", liberrors.NewAlreadyExistsError(assert.AnError), domainservice.IsConflictError},
		{"genericError", liberrors.NewGenericError(assert.AnError), domainservice.IsInternalError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testString := fmt.Sprintf("mapToBlueprintError(%v)", tt.givenError)
			resultError := mapToBlueprintError(tt.givenError)
			if tt.givenError != nil {
				assert.ErrorContains(t, resultError, tt.givenError.Error(), testString)
			}
			if tt.errTypeAssert != nil {
				assert.True(t, tt.errTypeAssert(resultError), testString)
			}
		})
	}
}
