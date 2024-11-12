package doguregistry

import (
	"context"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var myQualifiedTestDoguName = cescommons.QualifiedDoguName{
	Namespace:  "testing",
	SimpleName: "my-dogu",
}

func TestNewRemote(t *testing.T) {
	// given
	repoMock := newMockRemoteDoguDescriptorRepository(t)

	// when
	actual := NewRemote(repoMock)

	// then
	assert.NotEmpty(t, actual)
}

func TestRemote_GetDogu(t *testing.T) {
	t.Run("should return not found error", func(t *testing.T) {
		// given
		repoMock := newMockRemoteDoguDescriptorRepository(t)
		version, err := core.ParseVersion("1.2.3")
		require.NoError(t, err)
		qDoguVersion := cescommons.QualifiedDoguVersion{
			cescommons.QualifiedDoguName{Namespace: "testing", SimpleName: "my-dogu"},
			version,
		}
		repoMock.EXPECT().Get(context.TODO(), qDoguVersion).Return(nil, cescommons.DoguDescriptorNotFoundError)

		sut := &Remote{repoMock}

		// when
		actual, err := sut.GetDogu(context.TODO(), qDoguVersion)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorContains(t, err, "dogu \"testing/my-dogu\" with version \"1.2.3\" could not be found")
		expectedErr := &domainservice.NotFoundError{}
		assert.ErrorAs(t, err, &expectedErr)
	})
	t.Run("should return internal error", func(t *testing.T) {
		// given
		repoMock := newMockRemoteDoguDescriptorRepository(t)

		version, err := core.ParseVersion("1.2.3")
		require.NoError(t, err)
		qDoguVersion := cescommons.QualifiedDoguVersion{
			cescommons.QualifiedDoguName{Namespace: "testing", SimpleName: "my-dogu"},
			version,
		}
		repoMock.EXPECT().Get(context.TODO(), qDoguVersion).Return(nil, assert.AnError)

		sut := &Remote{repoMock}

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
	t.Run("should return dogu", func(t *testing.T) {
		// given
		expectedDogu := core.Dogu{Name: "testing/my-dogu", Version: "1.2.3"}

		repoMock := newMockRemoteDoguDescriptorRepository(t)
		version, err := core.ParseVersion("1.2.3")
		require.NoError(t, err)
		qDoguVersion := cescommons.QualifiedDoguVersion{
			cescommons.QualifiedDoguName{Namespace: "testing", SimpleName: "my-dogu"},
			version,
		}
		repoMock.EXPECT().Get(context.TODO(), qDoguVersion).Return(&expectedDogu, nil)

		sut := &Remote{repoMock}

		// when
		actual, err := sut.GetDogu(context.TODO(), qDoguVersion)

		// then
		require.NoError(t, err)
		assert.Equal(t, &expectedDogu, actual)
	})
}

func TestRemote_GetDogus(t *testing.T) {
	t.Run("should return collected errors", func(t *testing.T) {
		// given
		repoMock := newMockRemoteDoguDescriptorRepository(t)
		expectedDogu := core.Dogu{Name: "testing/good-dogu", Version: "0.1.2"}
		goodVersion, err := core.ParseVersion("0.1.2")
		require.NoError(t, err)
		NotFoundVersion, err := core.ParseVersion("1.2.3")
		require.NoError(t, err)
		OtherVersion, err := core.ParseVersion("2.3.4")
		require.NoError(t, err)
		qGoodDoguVersion := cescommons.QualifiedDoguVersion{
			cescommons.QualifiedDoguName{Namespace: "testing", SimpleName: "good-dogu"},
			goodVersion,
		}
		qNotFoundDoguVersion := cescommons.QualifiedDoguVersion{
			cescommons.QualifiedDoguName{Namespace: "testing", SimpleName: "not-found"},
			NotFoundVersion,
		}
		qOtherErrorDoguVersion := cescommons.QualifiedDoguVersion{
			cescommons.QualifiedDoguName{Namespace: "testing", SimpleName: "other-error"},
			OtherVersion,
		}
		repoMock.EXPECT().Get(context.TODO(), qGoodDoguVersion).Return(&expectedDogu, nil)
		notFoundError := fmt.Errorf("404 not found")
		repoMock.EXPECT().Get(context.TODO(), qNotFoundDoguVersion).Return(nil, notFoundError)
		repoMock.EXPECT().Get(context.TODO(), qOtherErrorDoguVersion).Return(nil, assert.AnError)

		sut := &Remote{repoMock}
		dogusToLoad := []cescommons.QualifiedDoguVersion{
			qGoodDoguVersion,
			qOtherErrorDoguVersion,
			qNotFoundDoguVersion,
		}

		expectedDogus := map[cescommons.QualifiedDoguName]*core.Dogu{
			cescommons.QualifiedDoguName{Namespace: "testing", SimpleName: "good-dogu"}:   &expectedDogu,
			cescommons.QualifiedDoguName{Namespace: "testing", SimpleName: "not-found"}:   nil,
			cescommons.QualifiedDoguName{Namespace: "testing", SimpleName: "other-error"}: nil,
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
