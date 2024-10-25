package kubernetes

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/cloudogu/k8s-registry-lib/errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

var testCtx = context.TODO()

// testGlobalConfig is used as immutable structure for tests
var testGlobalConfig = config.CreateGlobalConfig(map[config.Key]config.Value{
	"fqdn": "domain.local",
})

func TestNewGlobalConfigRepository(t *testing.T) {
	repoMock := newMockK8sGlobalConfigRepo(t)
	//given
	repo := NewGlobalConfigRepository(repoMock)
	//when
	assert.Equal(t, repoMock, repo.repo)
}

func TestGlobalConfigRepository_Get(t *testing.T) {

	t.Run("get", func(t *testing.T) {
		repoMock := newMockK8sGlobalConfigRepo(t)
		//given
		repoMock.EXPECT().Get(testCtx).Return(testGlobalConfig, nil)
		repo := NewGlobalConfigRepository(repoMock)
		//when
		actualConfig, err := repo.Get(testCtx)
		//then
		assert.NoError(t, err)
		assert.Equal(t, testGlobalConfig, actualConfig)
	})
	t.Run("global config not found", func(t *testing.T) {
		repoMock := newMockK8sGlobalConfigRepo(t)
		//given
		// wrap the error because the original implementation does it too
		givenError := fmt.Errorf("wrapping error: %w", errors.NewNotFoundError(assert.AnError))
		repoMock.EXPECT().Get(testCtx).Return(testGlobalConfig, givenError)
		repo := NewGlobalConfigRepository(repoMock)
		//when
		_, err := repo.Get(testCtx)
		//then
		assert.ErrorContains(t, err, givenError.Error())
		assert.ErrorContains(t, err, "could not load global config")
		assert.True(t, domainservice.IsNotFoundError(err), "error is no NotFoundError")
	})
	t.Run("internal error if connection error happens", func(t *testing.T) {
		repoMock := newMockK8sGlobalConfigRepo(t)
		//given
		// wrap the error because the original implementation does it too
		givenError := fmt.Errorf("wrapping error: %w", errors.NewConnectionError(assert.AnError))
		repoMock.EXPECT().Get(testCtx).Return(testGlobalConfig, givenError)
		repo := NewGlobalConfigRepository(repoMock)
		//when
		_, err := repo.Get(testCtx)
		//then
		assert.ErrorContains(t, err, givenError.Error())
		assert.ErrorContains(t, err, "could not load global config")
		assert.True(t, domainservice.IsInternalError(err), "error is no InternalError")
	})
	t.Run("internal error on all other errors", func(t *testing.T) {
		repoMock := newMockK8sGlobalConfigRepo(t)
		//given
		// wrap the error because the original implementation does it too
		givenError := fmt.Errorf("wrapping error: %w", errors.NewGenericError(assert.AnError))
		repoMock.EXPECT().Get(testCtx).Return(testGlobalConfig, givenError)
		repo := NewGlobalConfigRepository(repoMock)
		//when
		_, err := repo.Get(testCtx)
		//then
		assert.ErrorContains(t, err, givenError.Error())
		assert.ErrorContains(t, err, "could not load global config")
		assert.True(t, domainservice.IsInternalError(err), "error is no InternalError")
	})
}

func TestGlobalConfigRepository_Update(t *testing.T) {
	t.Run("update global config", func(t *testing.T) {
		repoMock := newMockK8sGlobalConfigRepo(t)
		//given
		repoMock.EXPECT().Update(testCtx, testGlobalConfig).Return(testGlobalConfig, nil)
		repo := NewGlobalConfigRepository(repoMock)
		//when
		actualConfig, err := repo.Update(testCtx, testGlobalConfig)
		//then
		assert.NoError(t, err)
		assert.Equal(t, testGlobalConfig, actualConfig)
	})
	t.Run("global config not found", func(t *testing.T) {
		repoMock := newMockK8sGlobalConfigRepo(t)
		//given
		// wrap the error because the original implementation does it too
		givenError := fmt.Errorf("wrapping error: %w", errors.NewNotFoundError(assert.AnError))
		repoMock.EXPECT().Update(testCtx, testGlobalConfig).Return(testGlobalConfig, givenError)
		repo := NewGlobalConfigRepository(repoMock)
		//when
		_, err := repo.Update(testCtx, testGlobalConfig)
		//then
		assert.ErrorContains(t, err, givenError.Error())
		assert.ErrorContains(t, err, "could not update global config")
		assert.True(t, domainservice.IsNotFoundError(err), "error is no NotFoundError")
	})
	t.Run("conflicts while updating global config", func(t *testing.T) {
		repoMock := newMockK8sGlobalConfigRepo(t)
		//given
		// wrap the error because the original implementation does it too
		givenError := fmt.Errorf("wrapping error: %w", errors.NewConflictError(assert.AnError))
		repoMock.EXPECT().Update(testCtx, testGlobalConfig).Return(testGlobalConfig, givenError)
		repo := NewGlobalConfigRepository(repoMock)
		//when
		_, err := repo.Update(testCtx, testGlobalConfig)
		//then
		assert.ErrorContains(t, err, givenError.Error())
		assert.ErrorContains(t, err, "could not update global config")
		assert.True(t, domainservice.IsConflictError(err), "error is no ConflictError")
	})
	t.Run("internal error if connection error happens", func(t *testing.T) {
		repoMock := newMockK8sGlobalConfigRepo(t)
		//given
		// wrap the error because the original implementation does it too
		givenError := fmt.Errorf("wrapping error: %w", errors.NewConnectionError(assert.AnError))
		repoMock.EXPECT().Update(testCtx, testGlobalConfig).Return(testGlobalConfig, givenError)
		repo := NewGlobalConfigRepository(repoMock)
		//when
		_, err := repo.Update(testCtx, testGlobalConfig)
		//then
		assert.ErrorContains(t, err, givenError.Error())
		assert.ErrorContains(t, err, "could not update global config")
		assert.True(t, domainservice.IsInternalError(err), "error is no InternalError")
	})
	t.Run("internal error on all other errors", func(t *testing.T) {
		repoMock := newMockK8sGlobalConfigRepo(t)
		//given
		// wrap the error because the original implementation does it too
		givenError := fmt.Errorf("wrapping error: %w", errors.NewGenericError(assert.AnError))
		repoMock.EXPECT().Update(testCtx, testGlobalConfig).Return(testGlobalConfig, givenError)
		repo := NewGlobalConfigRepository(repoMock)
		//when
		_, err := repo.Update(testCtx, testGlobalConfig)
		//then
		assert.ErrorContains(t, err, givenError.Error())
		assert.ErrorContains(t, err, "could not update global config")
		assert.True(t, domainservice.IsInternalError(err), "error is no InternalError")
	})
}
