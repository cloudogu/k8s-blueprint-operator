package maintenance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

func TestMode_Activate(t *testing.T) {
	t.Run("should fail to check if active and ours", func(t *testing.T) {
		// given
		lockMock := newMockLock(t)
		lockMock.EXPECT().isActiveAndOurs().Return(false, false, assert.AnError)

		switcherMock := newMockSwitcher(t)

		sut := &Mode{lock: lockMock, switcher: switcherMock}

		content := domainservice.MaintenancePageModel{
			Title: "myTitle",
			Text:  "myText",
		}

		// when
		err := sut.Activate(content)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		internalError := &domainservice.InternalError{}
		assert.ErrorAs(t, err, &internalError)
		assert.ErrorContains(t, err, "failed to check if maintenance mode is already active and ours")
	})
	t.Run("should fail if maintenance mode is active and not ours", func(t *testing.T) {
		// given
		lockMock := newMockLock(t)
		lockMock.EXPECT().isActiveAndOurs().Return(true, false, nil)

		switcherMock := newMockSwitcher(t)

		sut := &Mode{lock: lockMock, switcher: switcherMock}

		content := domainservice.MaintenancePageModel{
			Title: "myTitle",
			Text:  "myText",
		}

		// when
		err := sut.Activate(content)

		// then
		require.Error(t, err)
		conflictError := &domainservice.ConflictError{}
		assert.ErrorAs(t, err, &conflictError)
		assert.ErrorContains(t, err, "cannot activate maintenance mode as it was already activated by another party")
	})
	t.Run("should fail to activate maintenance mode", func(t *testing.T) {
		// given
		content := domainservice.MaintenancePageModel{
			Title: "myTitle",
			Text:  "myText",
		}

		lockMock := newMockLock(t)
		lockMock.EXPECT().isActiveAndOurs().Return(true, true, nil)

		switcherMock := newMockSwitcher(t)
		switcherMock.EXPECT().activate(content).Return(assert.AnError)

		sut := &Mode{lock: lockMock, switcher: switcherMock}

		// when
		err := sut.Activate(content)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		internalError := &domainservice.InternalError{}
		assert.ErrorAs(t, err, &internalError)
		assert.ErrorContains(t, err, "failed to activate maintenance mode")
	})
	t.Run("should succeed to activate maintenance mode", func(t *testing.T) {
		// given
		content := domainservice.MaintenancePageModel{
			Title: "myTitle",
			Text:  "myText",
		}

		lockMock := newMockLock(t)
		lockMock.EXPECT().isActiveAndOurs().Return(false, false, nil)

		switcherMock := newMockSwitcher(t)
		switcherMock.EXPECT().activate(content).Return(nil)

		sut := &Mode{lock: lockMock, switcher: switcherMock}

		// when
		err := sut.Activate(content)

		// then
		require.NoError(t, err)
	})
}

func TestMode_Deactivate(t *testing.T) {
	t.Run("should fail to check if active and ours", func(t *testing.T) {
		// given
		lockMock := newMockLock(t)
		lockMock.EXPECT().isActiveAndOurs().Return(false, false, assert.AnError)

		switcherMock := newMockSwitcher(t)

		sut := &Mode{lock: lockMock, switcher: switcherMock}

		// when
		err := sut.Deactivate()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		internalError := &domainservice.InternalError{}
		assert.ErrorAs(t, err, &internalError)
		assert.ErrorContains(t, err, "failed to check if maintenance mode is already active")
	})
	t.Run("should do nothing if maintenance mode is not active", func(t *testing.T) {
		// given
		lockMock := newMockLock(t)
		lockMock.EXPECT().isActiveAndOurs().Return(false, false, nil)

		switcherMock := newMockSwitcher(t)

		sut := &Mode{lock: lockMock, switcher: switcherMock}

		// when
		err := sut.Deactivate()

		// then
		require.NoError(t, err)
	})
	t.Run("should fail if maintenance mode was not activated by us", func(t *testing.T) {
		// given
		lockMock := newMockLock(t)
		lockMock.EXPECT().isActiveAndOurs().Return(true, false, nil)

		switcherMock := newMockSwitcher(t)

		sut := &Mode{lock: lockMock, switcher: switcherMock}

		// when
		err := sut.Deactivate()

		// then
		require.Error(t, err)
		conflictError := &domainservice.ConflictError{}
		assert.ErrorAs(t, err, &conflictError)
		assert.ErrorContains(t, err, "cannot deactivate maintenance mode as it was activated by another party")
	})
	t.Run("should fail to deactivate maintenance mode", func(t *testing.T) {
		// given
		lockMock := newMockLock(t)
		lockMock.EXPECT().isActiveAndOurs().Return(true, true, nil)

		switcherMock := newMockSwitcher(t)
		switcherMock.EXPECT().deactivate().Return(assert.AnError)

		sut := &Mode{lock: lockMock, switcher: switcherMock}

		// when
		err := sut.Deactivate()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		internalError := &domainservice.InternalError{}
		assert.ErrorAs(t, err, &internalError)
		assert.ErrorContains(t, err, "failed to deactivate maintenance mode")
	})
	t.Run("should succeed to deactivate maintenance mode", func(t *testing.T) {
		// given
		lockMock := newMockLock(t)
		lockMock.EXPECT().isActiveAndOurs().Return(true, true, nil)

		switcherMock := newMockSwitcher(t)
		switcherMock.EXPECT().deactivate().Return(nil)

		sut := &Mode{lock: lockMock, switcher: switcherMock}

		// when
		err := sut.Deactivate()

		// then
		require.NoError(t, err)
	})
}

func TestNew(t *testing.T) {
	globalConfigMock := newMockGlobalConfig(t)
	actual := New(globalConfigMock)
	assert.NotEmpty(t, actual)
}
