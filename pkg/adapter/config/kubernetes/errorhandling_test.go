package kubernetes

import (
	"fmt"
	liberrors "github.com/cloudogu/ces-commons-lib/errors"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
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
		{"otherError", assert.AnError, domainservice.IsInternalError},
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
