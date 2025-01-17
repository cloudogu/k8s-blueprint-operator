package componentcr

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
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

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
)

var testCtx = context.Background()

const (
	testComponentNameRaw = "my-component"
	testNamespace        = "k8s"
)

var testComponentName = common.QualifiedComponentName{
	Namespace:  testNamespace,
	SimpleName: testComponentNameRaw,
}

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
				Name:      string(testComponentName.SimpleName),
				Namespace: testNamespace,
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
					Name:            string(testComponentName.SimpleName),
					ResourceVersion: "42",
				},
				Spec: compV1.ComponentSpec{
					Name:      string(testComponentName.SimpleName),
					Namespace: string(testComponentName.Namespace),
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

		expected := map[common.SimpleComponentName]*ecosystem.ComponentInstallation{}
		version, _ := semver.NewVersion("1.2.3-4")
		expected[testComponentName.SimpleName] = &ecosystem.ComponentInstallation{
			Name:               testComponentName,
			ExpectedVersion:    version,
			Status:             "installed",
			Health:             "",
			PersistenceContext: nil,
		}
		assert.Equal(t, expected[testComponentName.SimpleName].Name, actual[testComponentName.SimpleName].Name)
		assert.Equal(t, expected[testComponentName.SimpleName].Status, actual[testComponentName.SimpleName].Status)
		assert.Equal(t, expected[testComponentName.SimpleName].ExpectedVersion, actual[testComponentName.SimpleName].ExpectedVersion)
		assert.Equal(t, expected[testComponentName.SimpleName].Health, actual[testComponentName.SimpleName].Health)
		// map pointers are hard to compare, test each field individually
		assert.Equal(t,
			map[string]any{componentInstallationRepoContextKey: componentInstallationRepoContext{resourceVersion: "42"}},
			actual[testComponentName.SimpleName].PersistenceContext)
	})
}

func Test_componentInstallationRepo_GetByName(t *testing.T) {
	t.Run("should return error when k8s client fails generically", func(t *testing.T) {
		// given
		mockRepo := newMockComponentRepo(t)
		mockRepo.EXPECT().Get(testCtx, string(testComponentName.SimpleName), metav1.GetOptions{}).Return(nil, assert.AnError)
		sut := componentInstallationRepo{componentClient: mockRepo}

		// when
		_, err := sut.GetByName(testCtx, testComponentName.SimpleName)

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
				Name:      string(testComponentName.SimpleName),
				Namespace: testNamespace,
				Version:   "a-b.c:d@1.2@parse-fail-here",
			},
		}
		mockRepo.EXPECT().Get(testCtx, string(testComponentName.SimpleName), metav1.GetOptions{}).Return(erroneousComponent, nil)

		// when
		_, err := sut.GetByName(testCtx, testComponentName.SimpleName)

		// then
		require.Error(t, err)
		assert.IsType(t, err, &domainservice.InternalError{})
		assert.ErrorContains(t, err, "cannot load component CR as it cannot be parsed correctly")
	})
	t.Run("should return InternalError when resource is nil", func(t *testing.T) {
		// given
		mockRepo := newMockComponentRepo(t)
		sut := componentInstallationRepo{componentClient: mockRepo}
		mockRepo.EXPECT().Get(testCtx, string(testComponentName.SimpleName), metav1.GetOptions{}).Return(nil, nil)

		// when
		_, err := sut.GetByName(testCtx, testComponentName.SimpleName)

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
			string(testComponentName.SimpleName))
		mockRepo.EXPECT().Get(testCtx, string(testComponentName.SimpleName), metav1.GetOptions{}).Return(nil, errNotFound)

		// when
		_, err := sut.GetByName(testCtx, testComponentName.SimpleName)

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
				Name:            string(testComponentName.SimpleName),
				ResourceVersion: "42",
			},
			Spec: compV1.ComponentSpec{
				Name:      string(testComponentName.SimpleName),
				Namespace: testNamespace,
				Version:   "1.2.3-4",
			},
			Status: compV1.ComponentStatus{
				Status: compV1.ComponentStatusInstalled,
				Health: compV1.PendingHealthStatus,
			},
		}
		mockRepo.EXPECT().Get(testCtx, string(testComponentName.SimpleName), metav1.GetOptions{}).Return(result, nil)

		// when
		actual, err := sut.GetByName(testCtx, testComponentName.SimpleName)

		// then
		require.NoError(t, err)

		version, _ := semver.NewVersion("1.2.3-4")
		expected := ecosystem.ComponentInstallation{
			Name:               testComponentName,
			ExpectedVersion:    version,
			Status:             "installed",
			Health:             "",
			PersistenceContext: nil,
		}
		assert.Equal(t, expected.Name, actual.Name)
		assert.Equal(t, expected.Name, testComponentName)
		assert.Equal(t, expected.Status, actual.Status)
		assert.Equal(t, expected.ExpectedVersion, actual.ExpectedVersion)
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
			Name:            testComponentName,
			ExpectedVersion: testVersion1,
			DeployConfig: map[string]interface{}{
				"deployNamespace": "longhorn-system",
				"overwriteConfig": map[string]interface{}{"key": "value"},
			},
		}
		patch := []byte("{\"spec\":{\"namespace\":\"k8s\",\"name\":\"my-component\",\"version\":\"1.0.0-1\",\"deployNamespace\":\"longhorn-system\",\"valuesYamlOverwrite\":\"key: value\\n\"}}")
		componentClientMock.EXPECT().Patch(testCtx, string(testComponentName.SimpleName), types.MergePatchType, patch, metav1.PatchOptions{}).Return(nil, nil)

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
			Name:            testComponentName,
			ExpectedVersion: testVersion1,
		}
		patch := []byte("{\"spec\":{\"namespace\":\"k8s\",\"name\":\"my-component\",\"version\":\"1.0.0-1\",\"deployNamespace\":null,\"valuesYamlOverwrite\":null}}")
		componentClientMock.EXPECT().Patch(testCtx, string(testComponentName.SimpleName), types.MergePatchType, patch, metav1.PatchOptions{}).Return(nil, assert.AnError)

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

		componentClientMock.EXPECT().Delete(testCtx, string(testComponentName.SimpleName), metav1.DeleteOptions{}).Return(nil)

		// when
		err := sut.Delete(testCtx, testComponentName.SimpleName)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on delete error", func(t *testing.T) {
		// given
		componentClientMock := newMockComponentRepo(t)
		sut := componentInstallationRepo{
			componentClient: componentClientMock,
		}

		componentClientMock.EXPECT().Delete(testCtx, string(testComponentName.SimpleName), metav1.DeleteOptions{}).Return(assert.AnError)

		// when
		err := sut.Delete(testCtx, testComponentName.SimpleName)

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
			Name:            testComponentName,
			ExpectedVersion: testVersion1,
		}
		expectedCR := &compV1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: string(componentInstallation.Name.SimpleName),
				Labels: map[string]string{
					ComponentNameLabelKey:    string(componentInstallation.Name.SimpleName),
					ComponentVersionLabelKey: componentInstallation.ExpectedVersion.String(),
				},
			},
			Spec: compV1.ComponentSpec{
				Namespace: string(componentInstallation.Name.Namespace),
				Name:      string(componentInstallation.Name.SimpleName),
				Version:   componentInstallation.ExpectedVersion.String(),
			},
		}

		componentClientMock.EXPECT().Create(testCtx, expectedCR, metav1.CreateOptions{}).Return(nil, nil)

		// when
		err := sut.Create(testCtx, componentInstallation)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on convert component installation error", func(t *testing.T) {
		// given
		componentClientMock := newMockComponentRepo(t)
		sut := componentInstallationRepo{
			componentClient: componentClientMock,
		}
		componentInstallation := &ecosystem.ComponentInstallation{
			Name:            testComponentName,
			ExpectedVersion: testVersion1,
			DeployConfig: map[string]interface{}{
				"deployNamespace": map[string]string{"no": "string"},
			},
		}

		// when
		err := sut.Create(testCtx, componentInstallation)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to convert component installation \"k8s/my-component\"")
		assert.IsType(t, err, &domainservice.InternalError{})
	})

	t.Run("should return error on patch error", func(t *testing.T) {
		// given
		componentClientMock := newMockComponentRepo(t)
		sut := componentInstallationRepo{
			componentClient: componentClientMock,
		}
		componentInstallation := &ecosystem.ComponentInstallation{
			Name:            testComponentName,
			ExpectedVersion: testVersion1,
		}
		expectedCR := &compV1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name: string(componentInstallation.Name.SimpleName),
				Labels: map[string]string{
					ComponentNameLabelKey:    string(componentInstallation.Name.SimpleName),
					ComponentVersionLabelKey: componentInstallation.ExpectedVersion.String(),
				},
			},
			Spec: compV1.ComponentSpec{
				Namespace: string(componentInstallation.Name.Namespace),
				Name:      string(componentInstallation.Name.SimpleName),
				Version:   componentInstallation.ExpectedVersion.String(),
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
