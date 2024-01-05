package domainservice

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMaintenanceModeUseCase_Activate(t *testing.T) {
	t.Run("should fail to get lock", func(t *testing.T) {
		// given
		maintenanceMock := NewMockMaintenanceMode(t)
		maintenanceMock.EXPECT().GetLock().Return(nil, assert.AnError)

		sut := &MaintenanceModeUseCase{maintenanceMode: maintenanceMock}

		content := MaintenancePageModel{
			Title: "myTitle",
			Text:  "myText",
		}

		// when
		err := sut.Activate(content)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to check if maintenance mode is already active")
	})
	t.Run("should fail if maintenance mode is active and not ours", func(t *testing.T) {
		// given
		lockMock := NewMockMaintenanceLock(t)
		lockMock.EXPECT().IsActive().Return(true)
		lockMock.EXPECT().IsOurs().Return(false)

		maintenanceMock := NewMockMaintenanceMode(t)
		maintenanceMock.EXPECT().GetLock().Return(lockMock, nil)

		sut := &MaintenanceModeUseCase{maintenanceMode: maintenanceMock}

		content := MaintenancePageModel{
			Title: "myTitle",
			Text:  "myText",
		}

		// when
		err := sut.Activate(content)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "cannot activate maintenance mode as someone else already activated it")
	})
	t.Run("should fail to activate maintenance mode", func(t *testing.T) {
		// given
		content := MaintenancePageModel{
			Title: "myTitle",
			Text:  "myText",
		}

		lockMock := NewMockMaintenanceLock(t)
		lockMock.EXPECT().IsActive().Return(true)
		lockMock.EXPECT().IsOurs().Return(true)

		maintenanceMock := NewMockMaintenanceMode(t)
		maintenanceMock.EXPECT().GetLock().Return(lockMock, nil)
		maintenanceMock.EXPECT().Activate(content).Return(assert.AnError)

		sut := &MaintenanceModeUseCase{maintenanceMode: maintenanceMock}

		// when
		err := sut.Activate(content)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to activate maintenance mode")
	})
	t.Run("should succeed to activate maintenance mode", func(t *testing.T) {
		// given
		content := MaintenancePageModel{
			Title: "myTitle",
			Text:  "myText",
		}

		lockMock := NewMockMaintenanceLock(t)
		lockMock.EXPECT().IsActive().Return(false)

		maintenanceMock := NewMockMaintenanceMode(t)
		maintenanceMock.EXPECT().GetLock().Return(lockMock, nil)
		maintenanceMock.EXPECT().Activate(content).Return(nil)

		sut := &MaintenanceModeUseCase{maintenanceMode: maintenanceMock}

		// when
		err := sut.Activate(content)

		// then
		require.NoError(t, err)
	})
}

func TestMaintenanceModeUseCase_Deactivate(t *testing.T) {
	t.Run("should fail to get lock", func(t *testing.T) {
		// given
		maintenanceMock := NewMockMaintenanceMode(t)
		maintenanceMock.EXPECT().GetLock().Return(nil, assert.AnError)

		sut := &MaintenanceModeUseCase{maintenanceMode: maintenanceMock}

		// when
		err := sut.Deactivate()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to check if maintenance mode is already active")
	})
	t.Run("should do nothing if maintenance mode is not active", func(t *testing.T) {
		// given
		lockMock := NewMockMaintenanceLock(t)
		lockMock.EXPECT().IsActive().Return(false)

		maintenanceMock := NewMockMaintenanceMode(t)
		maintenanceMock.EXPECT().GetLock().Return(lockMock, nil)

		sut := &MaintenanceModeUseCase{maintenanceMode: maintenanceMock}

		// when
		err := sut.Deactivate()

		// then
		require.NoError(t, err)
	})
	t.Run("should fail if maintenance mode was not activated by us", func(t *testing.T) {
		// given
		lockMock := NewMockMaintenanceLock(t)
		lockMock.EXPECT().IsActive().Return(true)
		lockMock.EXPECT().IsOurs().Return(false)

		maintenanceMock := NewMockMaintenanceMode(t)
		maintenanceMock.EXPECT().GetLock().Return(lockMock, nil)

		sut := &MaintenanceModeUseCase{maintenanceMode: maintenanceMock}

		// when
		err := sut.Deactivate()

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "cannot deactivate maintenance mode as it was activated by another application")
	})
	t.Run("should fail to deactivate maintenance mode", func(t *testing.T) {
		// given
		lockMock := NewMockMaintenanceLock(t)
		lockMock.EXPECT().IsActive().Return(true)
		lockMock.EXPECT().IsOurs().Return(true)

		maintenanceMock := NewMockMaintenanceMode(t)
		maintenanceMock.EXPECT().GetLock().Return(lockMock, nil)
		maintenanceMock.EXPECT().Deactivate().Return(assert.AnError)

		sut := &MaintenanceModeUseCase{maintenanceMode: maintenanceMock}

		// when
		err := sut.Deactivate()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to deactivate maintenance mode")
	})
	t.Run("should succeed to deactivate maintenance mode", func(t *testing.T) {
		// given
		lockMock := NewMockMaintenanceLock(t)
		lockMock.EXPECT().IsActive().Return(true)
		lockMock.EXPECT().IsOurs().Return(true)

		maintenanceMock := NewMockMaintenanceMode(t)
		maintenanceMock.EXPECT().GetLock().Return(lockMock, nil)
		maintenanceMock.EXPECT().Deactivate().Return(nil)

		sut := &MaintenanceModeUseCase{maintenanceMode: maintenanceMock}

		// when
		err := sut.Deactivate()

		// then
		require.NoError(t, err)
	})
}

func TestNewMaintenanceModeUseCase(t *testing.T) {
	maintenanceMock := NewMockMaintenanceMode(t)
	actual := NewMaintenanceModeUseCase(maintenanceMock)
	assert.NotEmpty(t, actual)
}
