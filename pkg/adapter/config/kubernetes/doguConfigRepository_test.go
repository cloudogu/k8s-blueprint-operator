package kubernetes

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/cloudogu/k8s-registry-lib/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

var doguCas = common.SimpleDoguName("cas")
var doguScm = common.SimpleDoguName("scm")
var testCasConfig = config.CreateDoguConfig(doguCas, map[config.Key]config.Value{
	"key1": "val1",
})
var testScmConfig = config.CreateDoguConfig(doguScm, map[config.Key]config.Value{
	"key1": "val1",
})

func TestNewDoguConfigRepository(t *testing.T) {
	repoMock := newMockK8sDoguConfigRepo(t)
	//given
	repo := NewDoguConfigRepository(repoMock)
	//when
	assert.Equal(t, repoMock, repo.repo)
}

func TestNewSensitiveDoguConfigRepository(t *testing.T) {
	repoMock := newMockK8sDoguConfigRepo(t)
	//given
	repo := NewSensitiveDoguConfigRepository(repoMock)
	//when
	assert.Equal(t, repoMock, repo.repo)
}

func TestDoguConfigRepository_Get(t *testing.T) {

	t.Run("get", func(t *testing.T) {
		repoMock := newMockK8sDoguConfigRepo(t)
		//given
		repoMock.EXPECT().Get(testCtx, doguCas).Return(testCasConfig, nil)
		repo := NewDoguConfigRepository(repoMock)
		//when
		actualConfig, err := repo.Get(testCtx, doguCas)
		//then
		assert.NoError(t, err)
		assert.Equal(t, testCasConfig, actualConfig)
	})
	t.Run("dogu config not found", func(t *testing.T) {
		repoMock := newMockK8sDoguConfigRepo(t)
		//given
		// wrap the error because the original implementation does it too
		givenError := fmt.Errorf("wrapping error: %w", errors.NewNotFoundError(assert.AnError))
		repoMock.EXPECT().Get(testCtx, doguCas).Return(testCasConfig, givenError)
		repo := NewDoguConfigRepository(repoMock)
		//when
		_, err := repo.Get(testCtx, doguCas)
		//then
		assert.ErrorContains(t, err, givenError.Error())
		assert.True(t, domainservice.IsNotFoundError(err), "error is no NotFoundError")
	})
	t.Run("internal error if connection error happens", func(t *testing.T) {
		repoMock := newMockK8sDoguConfigRepo(t)
		//given
		// wrap the error because the original implementation does it too
		givenError := fmt.Errorf("wrapping error: %w", errors.NewConnectionError(assert.AnError))
		repoMock.EXPECT().Get(testCtx, doguCas).Return(testCasConfig, givenError)
		repo := NewDoguConfigRepository(repoMock)
		//when
		_, err := repo.Get(testCtx, doguCas)
		//then
		assert.ErrorContains(t, err, givenError.Error())
		assert.ErrorContains(t, err, fmt.Sprintf("could not load normal dogu config for %s", doguCas.String()))
		assert.True(t, domainservice.IsInternalError(err), "error is no InternalError")
	})
	t.Run("internal error on all other errors", func(t *testing.T) {
		repoMock := newMockK8sDoguConfigRepo(t)
		//given
		// wrap the error because the original implementation does it too
		givenError := fmt.Errorf("wrapping error: %w", errors.NewGenericError(assert.AnError))
		repoMock.EXPECT().Get(testCtx, doguCas).Return(testCasConfig, givenError)
		repo := NewDoguConfigRepository(repoMock)
		//when
		_, err := repo.Get(testCtx, doguCas)
		//then
		assert.ErrorContains(t, err, givenError.Error())
		assert.ErrorContains(t, err, fmt.Sprintf("could not load normal dogu config for %s", doguCas.String()))
		assert.True(t, domainservice.IsInternalError(err), "error is no InternalError")
	})
}

func TestDoguConfigRepository_Update(t *testing.T) {
	t.Run("update dogu config", func(t *testing.T) {
		repoMock := newMockK8sDoguConfigRepo(t)
		//given
		repoMock.EXPECT().Update(testCtx, testCasConfig).Return(testCasConfig, nil)
		repo := NewDoguConfigRepository(repoMock)
		//when
		actualConfig, err := repo.Update(testCtx, testCasConfig)
		//then
		assert.NoError(t, err)
		assert.Equal(t, testCasConfig, actualConfig)
	})
	t.Run("dogu config not found", func(t *testing.T) {
		repoMock := newMockK8sDoguConfigRepo(t)
		//given
		// wrap the error because the original implementation does it too
		givenError := fmt.Errorf("wrapping error: %w", errors.NewNotFoundError(assert.AnError))
		repoMock.EXPECT().Update(testCtx, testCasConfig).Return(testCasConfig, givenError)
		repo := NewDoguConfigRepository(repoMock)
		//when
		_, err := repo.Update(testCtx, testCasConfig)
		//then
		assert.ErrorContains(t, err, givenError.Error())
		assert.True(t, domainservice.IsNotFoundError(err), "error is no NotFoundError")
	})
	t.Run("conflicts while updating dogu config", func(t *testing.T) {
		repoMock := newMockK8sDoguConfigRepo(t)
		//given
		// wrap the error because the original implementation does it too
		givenError := fmt.Errorf("wrapping error: %w", errors.NewConflictError(assert.AnError))
		repoMock.EXPECT().Update(testCtx, testCasConfig).Return(testCasConfig, givenError)
		repo := NewDoguConfigRepository(repoMock)
		//when
		_, err := repo.Update(testCtx, testCasConfig)
		//then
		assert.ErrorContains(t, err, givenError.Error())
		assert.ErrorContains(t, err, fmt.Sprintf("could not update normal dogu config for %s", doguCas.String()))
		assert.True(t, domainservice.IsConflictError(err), "error is no ConflictError")
	})
	t.Run("internal error if connection error happens", func(t *testing.T) {
		repoMock := newMockK8sDoguConfigRepo(t)
		//given
		// wrap the error because the original implementation does it too
		givenError := fmt.Errorf("wrapping error: %w", errors.NewConnectionError(assert.AnError))
		repoMock.EXPECT().Update(testCtx, testCasConfig).Return(testCasConfig, givenError)
		repo := NewDoguConfigRepository(repoMock)
		//when
		_, err := repo.Update(testCtx, testCasConfig)
		//then
		assert.ErrorContains(t, err, givenError.Error())
		assert.ErrorContains(t, err, fmt.Sprintf("could not update normal dogu config for %s", doguCas.String()))
		assert.True(t, domainservice.IsInternalError(err), "error is no InternalError")
	})
	t.Run("internal error on all other errors", func(t *testing.T) {
		repoMock := newMockK8sDoguConfigRepo(t)
		//given
		// wrap the error because the original implementation does it too
		givenError := fmt.Errorf("wrapping error: %w", errors.NewGenericError(assert.AnError))
		repoMock.EXPECT().Update(testCtx, testCasConfig).Return(testCasConfig, givenError)
		repo := NewDoguConfigRepository(repoMock)
		//when
		_, err := repo.Update(testCtx, testCasConfig)
		//then
		assert.ErrorContains(t, err, givenError.Error())
		assert.ErrorContains(t, err, fmt.Sprintf("could not update normal dogu config for %s", doguCas.String()))
		assert.True(t, domainservice.IsInternalError(err), "error is no InternalError")
	})
}

func TestDoguConfigRepository_GetAll(t *testing.T) {
	t.Run("getAll dogu config", func(t *testing.T) {
		repoMock := newMockK8sDoguConfigRepo(t)
		dogus := []common.SimpleDoguName{doguCas, doguScm}
		//given
		repoMock.EXPECT().Get(testCtx, doguCas).Return(testCasConfig, nil)
		repoMock.EXPECT().Get(testCtx, doguScm).Return(testScmConfig, nil)
		repo := NewDoguConfigRepository(repoMock)
		//when
		configByDogu, err := repo.GetAll(testCtx, dogus)
		//then
		assert.NoError(t, err)
		assert.Equal(t, map[common.SimpleDoguName]config.DoguConfig{
			doguCas: testCasConfig,
			doguScm: testScmConfig,
		}, configByDogu)
	})
	t.Run("getAll dogu config with error", func(t *testing.T) {
		repoMock := newMockK8sDoguConfigRepo(t)
		dogus := []common.SimpleDoguName{doguCas, doguScm}
		//given
		repoMock.EXPECT().Get(testCtx, doguCas).Return(testCasConfig, nil).Maybe()
		givenError := errors.NewNotFoundError(assert.AnError)
		repoMock.EXPECT().Get(testCtx, doguScm).Return(testScmConfig, givenError)
		repo := NewDoguConfigRepository(repoMock)
		//when
		_, err := repo.GetAll(testCtx, dogus)
		//then
		assert.ErrorContains(t, err, givenError.Error())
	})
}

func TestDoguConfigRepository_GetAllExisting(t *testing.T) {
	t.Run("all ok", func(t *testing.T) {
		repoMock := newMockK8sDoguConfigRepo(t)
		dogus := []common.SimpleDoguName{doguCas, doguScm}
		//given
		repoMock.EXPECT().Get(testCtx, doguCas).Return(testCasConfig, nil)
		repoMock.EXPECT().Get(testCtx, doguScm).Return(testScmConfig, nil)
		repo := NewDoguConfigRepository(repoMock)
		//when
		configByDogu, err := repo.GetAllExisting(testCtx, dogus)
		//then
		assert.NoError(t, err)
		assert.Equal(t, map[common.SimpleDoguName]config.DoguConfig{
			doguCas: testCasConfig,
			doguScm: testScmConfig,
		}, configByDogu)
	})
	t.Run("with NotFoundError", func(t *testing.T) {
		repoMock := newMockK8sDoguConfigRepo(t)
		dogus := []common.SimpleDoguName{doguCas, doguScm}
		//given
		repoMock.EXPECT().Get(testCtx, doguCas).Return(testCasConfig, nil)
		givenError := errors.NewNotFoundError(assert.AnError)
		repoMock.EXPECT().Get(testCtx, doguScm).Return(
			config.CreateDoguConfig(doguScm, map[config.Key]config.Value{}),
			givenError,
		)
		repo := NewDoguConfigRepository(repoMock)
		//when
		result, err := repo.GetAllExisting(testCtx, dogus)
		//then
		assert.Equal(t, map[config.SimpleDoguName]config.DoguConfig{
			doguCas: testCasConfig,
			doguScm: config.CreateDoguConfig(doguScm, map[config.Key]config.Value{}),
		}, result)
		assert.NoError(t, err, "should not throw an error if it is only a NotFoundError")
	})
	t.Run("with ConnectionError", func(t *testing.T) {
		repoMock := newMockK8sDoguConfigRepo(t)
		dogus := []common.SimpleDoguName{doguCas, doguScm}
		//given
		repoMock.EXPECT().Get(testCtx, doguCas).Return(testCasConfig, nil).Maybe()
		givenError := errors.NewConnectionError(assert.AnError)
		repoMock.EXPECT().Get(testCtx, doguScm).Return(
			config.CreateDoguConfig(doguScm, map[config.Key]config.Value{}),
			givenError,
		)

		repo := NewDoguConfigRepository(repoMock)
		//when
		_, err := repo.GetAllExisting(testCtx, dogus)
		//then
		assert.ErrorContains(t, err, err.Error())
	})

}
