package doguregistry

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var myQualifiedTestDoguName = common.QualifiedDoguName{
	Namespace: "testing",
	Name:      "my-dogu",
}

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
		actual, err := sut.GetDogu(myQualifiedTestDoguName, "1.2.3")

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
		actual, err := sut.GetDogu(myQualifiedTestDoguName, "1.2.3")

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
		actual, err := sut.GetDogu(myQualifiedTestDoguName, "1.2.3")

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
			{DoguName: common.QualifiedDoguName{
				Namespace: "testing",
				Name:      "good-dogu",
			}, Version: "0.1.2"},
			{DoguName: common.QualifiedDoguName{
				Namespace: "testing",
				Name:      "not-found",
			}, Version: "1.2.3"},
			{DoguName: common.QualifiedDoguName{
				Namespace: "testing",
				Name:      "other-error",
			}, Version: "2.3.4"},
		}

		expectedDogus := map[common.QualifiedDoguName]*core.Dogu{
			common.QualifiedDoguName{Namespace: "testing", Name: "good-dogu"}:   &expectedDogu,
			common.QualifiedDoguName{Namespace: "testing", Name: "not-found"}:   nil,
			common.QualifiedDoguName{Namespace: "testing", Name: "other-error"}: nil,
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
