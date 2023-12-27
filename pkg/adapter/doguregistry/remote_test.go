package doguregistry

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewRemote(t *testing.T) {
	// given
	regMock := newMockCesappLibRemoteRegistry(t)

	// when
	actual := NewRemote(regMock)

	// then
	assert.NotEmpty(t, actual)
}

func TestRemote_GetDogu(t *testing.T) {
	t.Run("should return not found error", func(t *testing.T) {
		// given
		notFoundError := fmt.Errorf("404 not found: %w", assert.AnError)
		regMock := newMockCesappLibRemoteRegistry(t)
		regMock.EXPECT().GetVersion("testing/my-dogu", "1.2.3").Return(nil, notFoundError)

		sut := &Remote{regMock}

		// when
		actual, err := sut.GetDogu("testing/my-dogu", "1.2.3")

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorContains(t, err, "dogu \"testing/my-dogu\" with version \"1.2.3\" could not be found")
		assert.ErrorIs(t, err, assert.AnError)
		expectedErr := &domainservice.NotFoundError{}
		assert.ErrorAs(t, err, &expectedErr)
	})
	t.Run("should return internal error", func(t *testing.T) {
		// given
		regMock := newMockCesappLibRemoteRegistry(t)
		regMock.EXPECT().GetVersion("testing/my-dogu", "1.2.3").Return(nil, assert.AnError)

		sut := &Remote{regMock}

		// when
		actual, err := sut.GetDogu("testing/my-dogu", "1.2.3")

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

		regMock := newMockCesappLibRemoteRegistry(t)
		regMock.EXPECT().GetVersion("testing/my-dogu", "1.2.3").Return(&expectedDogu, nil)

		sut := &Remote{regMock}

		// when
		actual, err := sut.GetDogu("testing/my-dogu", "1.2.3")

		// then
		require.NoError(t, err)
		assert.Equal(t, &expectedDogu, actual)
	})
}

func TestRemote_GetDogus(t *testing.T) {
	t.Run("should return collected errors", func(t *testing.T) {
		// given
		regMock := newMockCesappLibRemoteRegistry(t)
		expectedDogu := core.Dogu{Name: "testing/good-dogu", Version: "0.1.2"}
		regMock.EXPECT().GetVersion("testing/good-dogu", "0.1.2").Return(&expectedDogu, nil)
		notFoundError := fmt.Errorf("404 not found")
		regMock.EXPECT().GetVersion("testing/not-found", "1.2.3").Return(nil, notFoundError)
		regMock.EXPECT().GetVersion("testing/other-error", "2.3.4").Return(nil, assert.AnError)

		sut := &Remote{regMock}
		dogusToLoad := []domainservice.DoguToLoad{
			{QualifiedDoguName: "testing/good-dogu", Version: "0.1.2"},
			{QualifiedDoguName: "testing/not-found", Version: "1.2.3"},
			{QualifiedDoguName: "testing/other-error", Version: "2.3.4"},
		}

		expectedDogus := map[string]*core.Dogu{
			"testing/good-dogu":   &expectedDogu,
			"testing/not-found":   nil,
			"testing/other-error": nil,
		}

		// when
		actual, err := sut.GetDogus(dogusToLoad)

		// then
		require.Error(t, err)

		assert.ErrorContains(t, err, "dogu \"testing/not-found\" with version \"1.2.3\" could not be found")
		assert.ErrorIs(t, err, notFoundError)
		expectedNotFound := &domainservice.NotFoundError{}
		assert.ErrorAs(t, err, &expectedNotFound)

		assert.ErrorContains(t, err, "failed to get dogu \"testing/other-error\" with version \"2.3.4\"")
		assert.ErrorIs(t, err, assert.AnError)
		expectedInternal := &domainservice.InternalError{}
		assert.ErrorAs(t, err, &expectedInternal)

		assert.Equal(t, expectedDogus, actual)
	})
}
