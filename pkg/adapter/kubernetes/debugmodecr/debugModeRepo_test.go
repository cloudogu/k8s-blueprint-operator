package debugmodecr

import (
	"context"
	"errors"
	"testing"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	v1 "github.com/cloudogu/k8s-debug-mode-cr-lib/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var crResourceVersion = "abc"

var testCtx = context.Background()

func Test_doguInstallationRepo_GetByName(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// given
		debugModeClientMock := NewMockDebugModeInterface(t)
		repo := NewDebugModeRepo(debugModeClientMock)

		// when
		debugModeClientMock.EXPECT().Get(testCtx, debugModeSingletonCRName, metav1.GetOptions{}).Return(
			&v1.DebugMode{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:            debugModeSingletonCRName,
					ResourceVersion: crResourceVersion,
				},
				Spec: v1.DebugModeSpec{
					DeactivateTimestamp: metav1.Now(),
					TargetLogLevel:      "DEBUG",
				},
				Status: v1.DebugModeStatus{
					Phase: v1.DebugModeStatusSet,
				},
			}, nil)
		debugMode, err := repo.GetSingleton(testCtx)

		// then
		require.NoError(t, err)
		assert.Equal(t, &ecosystem.DebugMode{
			Phase: "SetDebugMode",
		}, debugMode)
	})

	t.Run("not found error", func(t *testing.T) {
		// given
		debugModeClientMock := NewMockDebugModeInterface(t)
		repo := NewDebugModeRepo(debugModeClientMock)

		// when
		debugModeClientMock.EXPECT().Get(testCtx, debugModeSingletonCRName, metav1.GetOptions{}).Return(
			nil,
			k8sErrors.NewNotFound(schema.GroupResource{}, debugModeSingletonCRName),
		)

		_, err := repo.GetSingleton(testCtx)

		// then
		require.Error(t, err)
		var expectedError *domainservice.NotFoundError
		assert.ErrorAs(t, err, &expectedError)
	})

	t.Run("internal error", func(t *testing.T) {
		// given
		debugModeClientMock := NewMockDebugModeInterface(t)
		repo := NewDebugModeRepo(debugModeClientMock)

		// when
		debugModeClientMock.EXPECT().Get(testCtx, debugModeSingletonCRName, metav1.GetOptions{}).Return(
			nil,
			k8sErrors.NewInternalError(errors.New("test-error")),
		)

		_, err := repo.GetSingleton(testCtx)

		// then
		require.Error(t, err)
		var expectedError *domainservice.InternalError
		assert.ErrorAs(t, err, &expectedError)
	})
}
