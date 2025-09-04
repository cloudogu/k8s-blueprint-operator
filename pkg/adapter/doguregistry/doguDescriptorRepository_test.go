package doguregistry

import (
	"context"
	"fmt"
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	cloudoguerrors "github.com/cloudogu/ces-commons-lib/errors"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewRemote(t *testing.T) {
	// given
	remoteRepoMock := newMockRemoteDoguDescriptorRepository(t)
	localRepoMock := newMockLocalDoguDescriptorRepository(t)

	// when
	actual := NewDoguDescriptorRepository(remoteRepoMock, localRepoMock)

	// then
	assert.NotEmpty(t, actual)
}

func TestRemote_GetDogu(t *testing.T) {
	t.Run("should return not found error on remote miss", func(t *testing.T) {
		// given
		remoteRepoMock := newMockRemoteDoguDescriptorRepository(t)
		localRepoMock := newMockLocalDoguDescriptorRepository(t)
		version, err := core.ParseVersion("1.2.3")
		require.NoError(t, err)
		qSimpleDoguVersion := cescommons.SimpleNameVersion{Name: "my-dogu", Version: version}
		// no local hit
		localRepoMock.EXPECT().Get(context.TODO(), qSimpleDoguVersion).Return(nil, assert.AnError)
		qDoguVersion := cescommons.QualifiedVersion{
			Name:    cescommons.QualifiedName{Namespace: "testing", SimpleName: "my-dogu"},
			Version: version,
		}
		remoteRepoMock.EXPECT().Get(context.TODO(), qDoguVersion).Return(nil, cloudoguerrors.NewNotFoundError(cloudoguerrors.Error{}))

		sut := &DoguDescriptorRepository{remoteRepoMock, localRepoMock}

		// when
		actual, err := sut.GetDogu(context.TODO(), qDoguVersion)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorContains(t, err, "dogu \"testing/my-dogu\" with version \"1.2.3\" could not be found")
		expectedErr := &domainservice.NotFoundError{}
		assert.ErrorAs(t, err, &expectedErr)
	})
	t.Run("should return internal error on remote error", func(t *testing.T) {
		// given
		remoteRepoMock := newMockRemoteDoguDescriptorRepository(t)
		localRepoMock := newMockLocalDoguDescriptorRepository(t)

		version, err := core.ParseVersion("1.2.3")
		require.NoError(t, err)
		// no local hit
		qSimpleDoguVersion := cescommons.SimpleNameVersion{Name: "my-dogu", Version: version}
		localRepoMock.EXPECT().Get(context.TODO(), qSimpleDoguVersion).Return(nil, assert.AnError)
		qDoguVersion := cescommons.QualifiedVersion{
			Name:    cescommons.QualifiedName{Namespace: "testing", SimpleName: "my-dogu"},
			Version: version,
		}
		remoteRepoMock.EXPECT().Get(context.TODO(), qDoguVersion).Return(nil, assert.AnError)

		sut := &DoguDescriptorRepository{remoteRepoMock, localRepoMock}

		// when
		actual, err := sut.GetDogu(context.TODO(), qDoguVersion)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorContains(t, err, "failed to get dogu \"testing/my-dogu\" with version \"1.2.3\"")
		assert.ErrorIs(t, err, assert.AnError)
		expectedErr := &domainservice.InternalError{}
		assert.ErrorAs(t, err, &expectedErr)
	})
	t.Run("should return dogu on remote hit", func(t *testing.T) {
		// given
		expectedDogu := core.Dogu{Name: "testing/my-dogu", Version: "1.2.3"}

		remoteRepoMock := newMockRemoteDoguDescriptorRepository(t)
		localRepoMock := newMockLocalDoguDescriptorRepository(t)
		version, err := core.ParseVersion("1.2.3")
		require.NoError(t, err)
		// no local hit
		qSimpleDoguVersion := cescommons.SimpleNameVersion{Name: "my-dogu", Version: version}
		localRepoMock.EXPECT().Get(context.TODO(), qSimpleDoguVersion).Return(nil, assert.AnError)
		qDoguVersion := cescommons.QualifiedVersion{
			Name:    cescommons.QualifiedName{Namespace: "testing", SimpleName: "my-dogu"},
			Version: version,
		}
		remoteRepoMock.EXPECT().Get(context.TODO(), qDoguVersion).Return(&expectedDogu, nil)
		localRepoMock.EXPECT().Add(context.TODO(), qDoguVersion.Name.SimpleName, &expectedDogu).Return(nil)

		sut := &DoguDescriptorRepository{remoteRepoMock, localRepoMock}

		// when
		actual, err := sut.GetDogu(context.TODO(), qDoguVersion)

		// then
		require.NoError(t, err)
		assert.Equal(t, &expectedDogu, actual)
	})
	t.Run("should return no error on local dogu description addition error", func(t *testing.T) {
		// given
		expectedDogu := core.Dogu{Name: "testing/my-dogu", Version: "1.2.3"}

		remoteRepoMock := newMockRemoteDoguDescriptorRepository(t)
		localRepoMock := newMockLocalDoguDescriptorRepository(t)
		version, err := core.ParseVersion("1.2.3")
		require.NoError(t, err)
		// no local hit
		qSimpleDoguVersion := cescommons.SimpleNameVersion{Name: "my-dogu", Version: version}
		localRepoMock.EXPECT().Get(context.TODO(), qSimpleDoguVersion).Return(nil, assert.AnError)
		qDoguVersion := cescommons.QualifiedVersion{
			Name:    cescommons.QualifiedName{Namespace: "testing", SimpleName: "my-dogu"},
			Version: version,
		}
		remoteRepoMock.EXPECT().Get(context.TODO(), qDoguVersion).Return(&expectedDogu, nil)
		localRepoMock.EXPECT().Add(context.TODO(), qDoguVersion.Name.SimpleName, &expectedDogu).Return(assert.AnError)

		sut := &DoguDescriptorRepository{remoteRepoMock, localRepoMock}

		// when
		actual, err := sut.GetDogu(context.TODO(), qDoguVersion)

		// then
		require.NoError(t, err)
		assert.Equal(t, &expectedDogu, actual)
	})
	t.Run("should return dogu on local hit", func(t *testing.T) {
		// given
		expectedDogu := core.Dogu{Name: "testing/my-dogu", Version: "1.2.3"}

		localRepoMock := newMockLocalDoguDescriptorRepository(t)
		version, err := core.ParseVersion("1.2.3")
		require.NoError(t, err)
		// no local hit
		qSimpleDoguVersion := cescommons.SimpleNameVersion{Name: "my-dogu", Version: version}
		localRepoMock.EXPECT().Get(context.TODO(), qSimpleDoguVersion).Return(&expectedDogu, nil)

		sut := &DoguDescriptorRepository{nil, localRepoMock}

		// when
		qDoguVersion := cescommons.QualifiedVersion{
			Name:    cescommons.QualifiedName{Namespace: "testing", SimpleName: "my-dogu"},
			Version: version,
		}
		actual, err := sut.GetDogu(context.TODO(), qDoguVersion)

		// then
		require.NoError(t, err)
		assert.Equal(t, &expectedDogu, actual)
	})
}

func TestRemote_GetDogus(t *testing.T) {
	t.Run("should return collected errors on remote calls", func(t *testing.T) {
		// given
		remoteRepoMock := newMockRemoteDoguDescriptorRepository(t)
		localRepoMock := newMockLocalDoguDescriptorRepository(t)
		expectedDogu := core.Dogu{Name: "testing/good-dogu", Version: "0.1.2"}
		goodVersion, err := core.ParseVersion("0.1.2")
		require.NoError(t, err)
		NotFoundVersion, err := core.ParseVersion("1.2.3")
		require.NoError(t, err)
		OtherVersion, err := core.ParseVersion("2.3.4")
		require.NoError(t, err)
		qGoodDoguVersion := cescommons.QualifiedVersion{
			Name:    cescommons.QualifiedName{Namespace: "testing", SimpleName: "good-dogu"},
			Version: goodVersion,
		}
		qNotFoundDoguVersion := cescommons.QualifiedVersion{
			Name:    cescommons.QualifiedName{Namespace: "testing", SimpleName: "not-found"},
			Version: NotFoundVersion,
		}
		qOtherErrorDoguVersion := cescommons.QualifiedVersion{
			Name:    cescommons.QualifiedName{Namespace: "testing", SimpleName: "other-error"},
			Version: OtherVersion,
		}
		// no local hits
		localRepoMock.EXPECT().Get(context.TODO(), mock.Anything).Return(nil, assert.AnError)
		remoteRepoMock.EXPECT().Get(context.TODO(), qGoodDoguVersion).Return(&expectedDogu, nil)
		notFoundError := fmt.Errorf("404 not found")
		remoteRepoMock.EXPECT().Get(context.TODO(), qNotFoundDoguVersion).Return(nil, notFoundError)
		remoteRepoMock.EXPECT().Get(context.TODO(), qOtherErrorDoguVersion).Return(nil, assert.AnError)
		localRepoMock.EXPECT().Add(context.TODO(), qGoodDoguVersion.Name.SimpleName, &expectedDogu).Return(nil)

		sut := &DoguDescriptorRepository{remoteRepoMock, localRepoMock}
		dogusToLoad := []cescommons.QualifiedVersion{
			qGoodDoguVersion,
			qOtherErrorDoguVersion,
			qNotFoundDoguVersion,
		}

		expectedDogus := map[cescommons.QualifiedName]*core.Dogu{
			cescommons.QualifiedName{Namespace: "testing", SimpleName: "good-dogu"}:   &expectedDogu,
			cescommons.QualifiedName{Namespace: "testing", SimpleName: "not-found"}:   nil,
			cescommons.QualifiedName{Namespace: "testing", SimpleName: "other-error"}: nil,
		}

		// when
		actual, err := sut.GetDogus(context.TODO(), dogusToLoad)

		// then
		require.Error(t, err)

		assert.ErrorContains(t, err, "dogu \"testing/not-found\" with version \"1.2.3\": 404 not found")
		assert.ErrorIs(t, err, notFoundError)

		assert.ErrorContains(t, err, "failed to get dogu \"testing/other-error\" with version \"2.3.4\"")
		assert.ErrorIs(t, err, assert.AnError)
		expectedInternal := &domainservice.InternalError{}
		assert.ErrorAs(t, err, &expectedInternal)

		assert.Equal(t, expectedDogus, actual)
	})
}
