package maintenance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_defaultLock_isActiveAndOurs(t *testing.T) {
	t.Run("should fail to check if key exists", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Exists("maintenance").Return(false, assert.AnError)

		sut := &defaultLock{globalConfig: globalConfigMock}

		// when
		_, _, err := sut.isActiveAndOurs()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to check if maintenance mode registry key exists")
	})
	t.Run("should both be false if key doesn't exist", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Exists("maintenance").Return(false, nil)

		sut := &defaultLock{globalConfig: globalConfigMock}

		// when
		active, ours, err := sut.isActiveAndOurs()

		// then
		require.NoError(t, err)
		assert.False(t, active)
		assert.False(t, ours)
	})
	t.Run("should fail to get key", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Exists("maintenance").Return(true, nil)
		globalConfigMock.EXPECT().Get("maintenance").Return("", assert.AnError)

		sut := &defaultLock{globalConfig: globalConfigMock}

		// when
		_, _, err := sut.isActiveAndOurs()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get maintenance mode from configuration registry")
	})
	t.Run("should fail to unmarshal json", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Exists("maintenance").Return(true, nil)
		globalConfigMock.EXPECT().Get("maintenance").Return("{invalid", nil)

		sut := &defaultLock{globalConfig: globalConfigMock}

		// when
		_, _, err := sut.isActiveAndOurs()

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to parse json of maintenance mode object")
	})
	t.Run("should be active but not ours if holder is missing", func(t *testing.T) {
		// given
		maintenanceJson := `{
								"title": "Attention",
								"text": "This is your captain speaking"
							}`
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Exists("maintenance").Return(true, nil)
		globalConfigMock.EXPECT().Get("maintenance").Return(maintenanceJson, nil)

		sut := &defaultLock{globalConfig: globalConfigMock}

		// when
		active, ours, err := sut.isActiveAndOurs()

		// then
		require.NoError(t, err)
		assert.True(t, active)
		assert.False(t, ours)
	})
	t.Run("should be active but not ours if holder is not us", func(t *testing.T) {
		// given
		maintenanceJson := `{
								"title": "Attention",
								"text": "This is your captain speaking",
								"holder": "cap'n"
							}`
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Exists("maintenance").Return(true, nil)
		globalConfigMock.EXPECT().Get("maintenance").Return(maintenanceJson, nil)

		sut := &defaultLock{globalConfig: globalConfigMock}

		// when
		active, ours, err := sut.isActiveAndOurs()

		// then
		require.NoError(t, err)
		assert.True(t, active)
		assert.False(t, ours)
	})
	t.Run("should be active and ours if holder is us", func(t *testing.T) {
		// given
		maintenanceJson := `{
								"title": "Attention",
								"text": "This is your captain speaking",
								"holder": "k8s-blueprint-operator"
							}`
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Exists("maintenance").Return(true, nil)
		globalConfigMock.EXPECT().Get("maintenance").Return(maintenanceJson, nil)

		sut := &defaultLock{globalConfig: globalConfigMock}

		// when
		active, ours, err := sut.isActiveAndOurs()

		// then
		require.NoError(t, err)
		assert.True(t, active)
		assert.True(t, ours)
	})
}
