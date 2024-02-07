package componentcr

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	compV1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

var testCtx = context.Background()

const (
	testComponentName         = "my-component"
	testDistributionNamespace = "k8s"
)

func Test_componentInstallationRepo_GetAll(t *testing.T) {
	t.Run("should return error when k8s client fails generically", func(t *testing.T) {
		// given
		mockRepo := newMockComponentRepo(t)
		mockRepo.EXPECT().List(testCtx, mock.Anything).Return(nil, assert.AnError)
		sut := componentInstallationRepo{componentClient: mockRepo}

		// when
		_, err := sut.GetAll(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should return InternalError when resource parsing fails", func(t *testing.T) {
		// given
		mockRepo := newMockComponentRepo(t)
		sut := componentInstallationRepo{componentClient: mockRepo}
		listWithErroneousElement := &compV1.ComponentList{Items: []compV1.Component{{
			Spec: compV1.ComponentSpec{
				Name:      testComponentName,
				Namespace: testDistributionNamespace,
				Version:   "a-b.c:d@1.2@parse-fail-here",
			},
		}}}
		mockRepo.EXPECT().List(testCtx, mock.Anything).Return(listWithErroneousElement, nil)

		// when
		_, err := sut.GetAll(testCtx)

		// then
		require.Error(t, err)
		assert.IsType(t, err, &domainservice.InternalError{})
		assert.ErrorContains(t, err, "failed to parse component CR")

	})
	t.Run("should return all existing blueprint resources", func(t *testing.T) {
		// given
		mockRepo := newMockComponentRepo(t)
		sut := componentInstallationRepo{componentClient: mockRepo}
		list := &compV1.ComponentList{Items: []compV1.Component{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:            testComponentName,
					ResourceVersion: "42",
				},
				Spec: compV1.ComponentSpec{
					Name:      testComponentName,
					Namespace: testDistributionNamespace,
					Version:   "1.2.3-4",
				},
				Status: compV1.ComponentStatus{
					Status: compV1.ComponentStatusInstalled,
					Health: compV1.PendingHealthStatus,
				},
			},
		}}
		mockRepo.EXPECT().List(testCtx, mock.Anything).Return(list, nil)

		// when
		actual, err := sut.GetAll(testCtx)

		// then
		require.NoError(t, err)

		expected := map[string]*ecosystem.ComponentInstallation{}
		version, _ := semver.NewVersion("1.2.3-4")
		expected[testComponentName] = &ecosystem.ComponentInstallation{
			Name:                  testComponentName,
			DistributionNamespace: testDistributionNamespace,
			Version:               version,
			Status:                "installed",
			Health:                "",
			PersistenceContext:    nil,
		}
		assert.Equal(t, expected[testComponentName].Name, actual[testComponentName].Name)
		assert.Equal(t, expected[testComponentName].DistributionNamespace, actual[testComponentName].DistributionNamespace)
		assert.Equal(t, expected[testComponentName].Status, actual[testComponentName].Status)
		assert.Equal(t, expected[testComponentName].Version, actual[testComponentName].Version)
		assert.Equal(t, expected[testComponentName].Health, actual[testComponentName].Health)
		// map pointers are hard to compare, test each field individually
		assert.Equal(t,
			map[string]any{componentInstallationRepoContextKey: componentInstallationRepoContext{resourceVersion: "42"}},
			actual[testComponentName].PersistenceContext)
	})
}

func Test_componentInstallationRepo_GetByName(t *testing.T) {
	t.Run("should return error when k8s client fails generically", func(t *testing.T) {
		// given
		mockRepo := newMockComponentRepo(t)
		mockRepo.EXPECT().Get(testCtx, testComponentName, mock.Anything).Return(nil, assert.AnError)
		sut := componentInstallationRepo{componentClient: mockRepo}

		// when
		_, err := sut.GetByName(testCtx, testComponentName)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should return InternalError when resource parsing fails", func(t *testing.T) {
		// given
		mockRepo := newMockComponentRepo(t)
		sut := componentInstallationRepo{componentClient: mockRepo}
		erroneousComponent := &compV1.Component{
			Spec: compV1.ComponentSpec{
				Name:      testComponentName,
				Namespace: testDistributionNamespace,
				Version:   "a-b.c:d@1.2@parse-fail-here",
			},
		}
		mockRepo.EXPECT().Get(testCtx, testComponentName, mock.Anything).Return(erroneousComponent, nil)

		// when
		_, err := sut.GetByName(testCtx, testComponentName)

		// then
		require.Error(t, err)
		assert.IsType(t, err, &domainservice.InternalError{})
		assert.ErrorContains(t, err, "cannot load component CR as it cannot be parsed correctly")
	})
	t.Run("should return InternalError when resource is nil", func(t *testing.T) {
		// given
		mockRepo := newMockComponentRepo(t)
		sut := componentInstallationRepo{componentClient: mockRepo}
		mockRepo.EXPECT().Get(testCtx, testComponentName, mock.Anything).Return(nil, nil)

		// when
		_, err := sut.GetByName(testCtx, testComponentName)

		// then
		require.Error(t, err)
		assert.IsType(t, err, &domainservice.InternalError{})
		assert.ErrorContains(t, err, "cannot parse component CR as it is nil")
	})
	t.Run("should return NotFoundError when resource does not exist", func(t *testing.T) {
		// given
		mockRepo := newMockComponentRepo(t)
		sut := componentInstallationRepo{componentClient: mockRepo}
		errNotFound := errors.NewNotFound(
			schema.GroupResource{Group: compV1.GroupVersion.Group, Resource: "component"},
			testComponentName)
		mockRepo.EXPECT().Get(testCtx, testComponentName, mock.Anything).Return(nil, errNotFound)

		// when
		_, err := sut.GetByName(testCtx, testComponentName)

		// then
		require.Error(t, err)
		assert.IsType(t, err, &domainservice.NotFoundError{})
		assert.ErrorContains(t, err, `cannot read component CR "my-component" as it does not exist`)
	})
	t.Run("should return all existing blueprint resources", func(t *testing.T) {
		// given
		mockRepo := newMockComponentRepo(t)
		sut := componentInstallationRepo{componentClient: mockRepo}
		result := &compV1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name:            testComponentName,
				ResourceVersion: "42",
			},
			Spec: compV1.ComponentSpec{
				Name:      testComponentName,
				Namespace: testDistributionNamespace,
				Version:   "1.2.3-4",
			},
			Status: compV1.ComponentStatus{
				Status: compV1.ComponentStatusInstalled,
				Health: compV1.PendingHealthStatus,
			},
		}
		mockRepo.EXPECT().Get(testCtx, testComponentName, mock.Anything).Return(result, nil)

		// when
		actual, err := sut.GetByName(testCtx, testComponentName)

		// then
		require.NoError(t, err)

		version, _ := semver.NewVersion("1.2.3-4")
		expected := ecosystem.ComponentInstallation{
			Name:                  testComponentName,
			DistributionNamespace: testDistributionNamespace,
			Version:               version,
			Status:                "installed",
			Health:                "",
			PersistenceContext:    nil,
		}
		assert.Equal(t, expected.Name, actual.Name)
		assert.Equal(t, expected.Name, testComponentName)
		assert.Equal(t, expected.DistributionNamespace, actual.DistributionNamespace)
		assert.Equal(t, expected.DistributionNamespace, testDistributionNamespace)
		assert.Equal(t, expected.Status, actual.Status)
		assert.Equal(t, expected.Version, actual.Version)
		assert.Equal(t, expected.Health, actual.Health)
		// map pointers are hard to compare, test each field individually
		assert.Equal(t,
			map[string]any{componentInstallationRepoContextKey: componentInstallationRepoContext{resourceVersion: "42"}}, actual.PersistenceContext)
	})
}

func TestNewComponentInstallationRepo(t *testing.T) {
	t.Run("should return proper repo interface implementation", func(t *testing.T) {
		// given
		mockRepo := newMockComponentRepo(t)

		// when
		actual := NewComponentInstallationRepo(mockRepo)

		// then
		assert.Implements(t, (*domainservice.ComponentInstallationRepository)(nil), actual)
	})
}

func Test_componentInstallationRepo_Update(t *testing.T) {
	t.Run("should patch cr on update", func(t *testing.T) {
		// given
		componentClientMock := newMockComponentRepo(t)
		sut := componentInstallationRepo{
			componentClient: componentClientMock,
		}
		componentInstallation := &ecosystem.ComponentInstallation{
			Name:                  testComponentName,
			DistributionNamespace: testDistributionNamespace,
			Version:               testVersion1,
		}
		patch := []byte("{\"spec\":{\"namespace\":\"k8s\",\"name\":\"my-component\",\"version\":\"1.0.0-1\"}}")
		componentClientMock.EXPECT().Patch(testCtx, testComponentName, types.MergePatchType, patch, metav1.PatchOptions{}).Return(nil, nil)

		// when
		err := sut.Update(testCtx, componentInstallation)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on patch error", func(t *testing.T) {
		// given
		componentClientMock := newMockComponentRepo(t)
		sut := componentInstallationRepo{
			componentClient: componentClientMock,
		}
		componentInstallation := &ecosystem.ComponentInstallation{
			Name:                  testComponentName,
			DistributionNamespace: testDistributionNamespace,
			Version:               testVersion1,
		}
		patch := []byte("{\"spec\":{\"namespace\":\"k8s\",\"name\":\"my-component\",\"version\":\"1.0.0-1\"}}")
		componentClientMock.EXPECT().Patch(testCtx, testComponentName, types.MergePatchType, patch, metav1.PatchOptions{}).Return(nil, assert.AnError)

		// when
		err := sut.Update(testCtx, componentInstallation)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to patch component \"my-component\"")
		assert.IsType(t, err, &domainservice.InternalError{})
	})
}

func Test_componentInstallationRepo_Delete(t *testing.T) {
	t.Run("should delete cr on delete", func(t *testing.T) {
		// given
		componentClientMock := newMockComponentRepo(t)
		sut := componentInstallationRepo{
			componentClient: componentClientMock,
		}

		componentClientMock.EXPECT().Delete(testCtx, testComponentName, metav1.DeleteOptions{}).Return(nil)

		// when
		err := sut.Delete(testCtx, testComponentName)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on delete error", func(t *testing.T) {
		// given
		componentClientMock := newMockComponentRepo(t)
		sut := componentInstallationRepo{
			componentClient: componentClientMock,
		}

		componentClientMock.EXPECT().Delete(testCtx, testComponentName, metav1.DeleteOptions{}).Return(assert.AnError)

		// when
		err := sut.Delete(testCtx, testComponentName)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to delete component CR \"my-component\"")
		assert.IsType(t, err, &domainservice.InternalError{})
	})
}

func Test_componentInstallationRepo_Create(t *testing.T) {
	t.Run("should create cr on create", func(t *testing.T) {
		// given
		componentClientMock := newMockComponentRepo(t)
		sut := componentInstallationRepo{
			componentClient: componentClientMock,
		}
		componentInstallation := &ecosystem.ComponentInstallation{
			Name:                  testComponentName,
			DistributionNamespace: testDistributionNamespace,
			Version:               testVersion1,
		}
		expectedCR := &compV1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: componentInstallation.Name,
				Labels: map[string]string{
					ComponentNameLabelKey:    componentInstallation.Name,
					ComponentVersionLabelKey: componentInstallation.Version.String(),
				},
			},
			Spec: compV1.ComponentSpec{
				Namespace: componentInstallation.DistributionNamespace,
				Name:      componentInstallation.Name,
				Version:   componentInstallation.Version.String(),
			},
		}

		componentClientMock.EXPECT().Create(testCtx, expectedCR, metav1.CreateOptions{}).Return(nil, nil)

		// when
		err := sut.Create(testCtx, componentInstallation)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on patch error", func(t *testing.T) {
		// given
		componentClientMock := newMockComponentRepo(t)
		sut := componentInstallationRepo{
			componentClient: componentClientMock,
		}
		componentInstallation := &ecosystem.ComponentInstallation{
			Name:                  testComponentName,
			DistributionNamespace: testDistributionNamespace,
			Version:               testVersion1,
		}
		expectedCR := &compV1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: componentInstallation.Name,
				Labels: map[string]string{
					ComponentNameLabelKey:    componentInstallation.Name,
					ComponentVersionLabelKey: componentInstallation.Version.String(),
				},
			},
			Spec: compV1.ComponentSpec{
				Namespace: componentInstallation.DistributionNamespace,
				Name:      componentInstallation.Name,
				Version:   componentInstallation.Version.String(),
			},
		}

		componentClientMock.EXPECT().Create(testCtx, expectedCR, metav1.CreateOptions{}).Return(nil, assert.AnError)

		// when
		err := sut.Create(testCtx, componentInstallation)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to create component CR \"my-component\"")
		assert.IsType(t, err, &domainservice.InternalError{})
	})
}