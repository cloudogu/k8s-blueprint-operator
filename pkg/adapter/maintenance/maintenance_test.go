package maintenance

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSwitch_GetLock(t *testing.T) {
	t.Run("should fail to check if maintenance mode is active", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Exists("maintenance").Return(false, assert.AnError)

		sut := &Switch{globalConfig: globalConfigMock}

		// when
		_, err := sut.GetLock()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		internalError := &domainservice.InternalError{}
		assert.ErrorAs(t, err, &internalError)
		assert.ErrorContains(t, err, "failed to check if maintenance mode registry key exists")
	})
	t.Run("should return lock directly if not active", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Exists("maintenance").Return(false, nil)

		sut := &Switch{globalConfig: globalConfigMock}

		// when
		actual, err := sut.GetLock()

		// then
		require.NoError(t, err)
		assert.Equal(t, lock{isActive: false, isOurs: false}, actual)
	})
	t.Run("should fail to get current maintenance mode", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Exists("maintenance").Return(true, nil)
		globalConfigMock.EXPECT().Get("maintenance").Return("", assert.AnError)

		sut := &Switch{globalConfig: globalConfigMock}

		// when
		_, err := sut.GetLock()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		internalError := &domainservice.InternalError{}
		assert.ErrorAs(t, err, &internalError)
		assert.ErrorContains(t, err, "failed to get maintenance mode from configuration registry")
	})
	t.Run("should fail to unmarshal maintenance mode information", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Exists("maintenance").Return(true, nil)
		globalConfigMock.EXPECT().Get("maintenance").Return("{{[[invalid, \"json\"", nil)

		sut := &Switch{globalConfig: globalConfigMock}

		// when
		_, err := sut.GetLock()

		// then
		require.Error(t, err)
		internalError := &domainservice.InternalError{}
		assert.ErrorAs(t, err, &internalError)
		assert.ErrorContains(t, err, "failed to unmarshal json of maintenance mode object")
	})
	t.Run("should return active lock that isn't ours for non-existent holder", func(t *testing.T) {
		// given
		maintenanceJson := `{
								"title": "myTitle",
								"text": "myText"
							}`
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Exists("maintenance").Return(true, nil)
		globalConfigMock.EXPECT().Get("maintenance").Return(maintenanceJson, nil)

		sut := &Switch{globalConfig: globalConfigMock}

		// when
		actual, err := sut.GetLock()

		// then
		require.NoError(t, err)
		assert.Equal(t, lock{isActive: true, isOurs: false}, actual)
	})
	t.Run("should return active lock that isn't ours for holder other than us", func(t *testing.T) {
		// given
		maintenanceJson := `{
								"title": "myTitle",
								"text": "myText",
								"holder": "not us"
							}`
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Exists("maintenance").Return(true, nil)
		globalConfigMock.EXPECT().Get("maintenance").Return(maintenanceJson, nil)

		sut := &Switch{globalConfig: globalConfigMock}

		// when
		actual, err := sut.GetLock()

		// then
		require.NoError(t, err)
		assert.Equal(t, lock{isActive: true, isOurs: false}, actual)
	})
	t.Run("should return active lock that is ours if holder matches", func(t *testing.T) {
		// given
		maintenanceJson := `{
								"title": "myTitle",
								"text": "myText",
								"holder": "k8s-blueprint-operator"
							}`
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Exists("maintenance").Return(true, nil)
		globalConfigMock.EXPECT().Get("maintenance").Return(maintenanceJson, nil)

		sut := &Switch{globalConfig: globalConfigMock}

		// when
		actual, err := sut.GetLock()

		// then
		require.NoError(t, err)
		assert.Equal(t, lock{isActive: true, isOurs: true}, actual)
	})
}

func TestSwitch_Activate(t *testing.T) {
	t.Run("should fail to activate maintenance mode", func(t *testing.T) {
		// given
		expectedJson := `{"title":"myTitle","text":"myText","holder":"k8s-blueprint-operator"}`
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Set("maintenance", expectedJson).Return(assert.AnError)

		sut := &Switch{globalConfig: globalConfigMock}

		// when
		err := sut.Activate(domainservice.MaintenancePageModel{
			Title: "myTitle",
			Text:  "myText",
		})

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		internalError := &domainservice.InternalError{}
		assert.ErrorAs(t, err, &internalError)
		assert.ErrorContains(t, err, "failed to set maintenance mode registry key")
	})
	t.Run("should succeed to activate maintenance mode", func(t *testing.T) {
		// given
		expectedJson := `{"title":"myTitle","text":"myText","holder":"k8s-blueprint-operator"}`
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Set("maintenance", expectedJson).Return(nil)

		sut := &Switch{globalConfig: globalConfigMock}

		// when
		err := sut.Activate(domainservice.MaintenancePageModel{
			Title: "myTitle",
			Text:  "myText",
		})

		// then
		require.NoError(t, err)
	})
}

func TestSwitch_Deactivate(t *testing.T) {
	t.Run("should fail to deactivate maintenance mode", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Delete("maintenance").Return(assert.AnError)

		sut := &Switch{globalConfig: globalConfigMock}

		// when
		err := sut.Deactivate()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		internalError := &domainservice.InternalError{}
		assert.ErrorAs(t, err, &internalError)
		assert.ErrorContains(t, err, "failed to delete maintenance mode registry key")
	})
	t.Run("should succeed to deactivate maintenance mode", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Delete("maintenance").Return(nil)

		sut := &Switch{globalConfig: globalConfigMock}

		// when
		err := sut.Deactivate()

		// then
		require.NoError(t, err)
	})
}

func TestNewSwitch(t *testing.T) {
	globalConfigMock := newMockGlobalConfig(t)
	actual := NewSwitch(globalConfigMock)
	assert.NotEmpty(t, actual)
}

func Test_lock_IsActive(t *testing.T) {
	type fields struct {
		isActive bool
		isOurs   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{fields: fields{isActive: false, isOurs: false}, want: false},
		{fields: fields{isActive: false, isOurs: true}, want: false},
		{fields: fields{isActive: true, isOurs: true}, want: true},
		{fields: fields{isActive: true, isOurs: false}, want: true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("isActive: %t, isOurs: %t", tt.fields.isActive, tt.fields.isOurs), func(t *testing.T) {
			l := lock{
				isActive: tt.fields.isActive,
				isOurs:   tt.fields.isOurs,
			}
			assert.Equalf(t, tt.want, l.IsActive(), "IsActive()")
		})
	}
}

func Test_lock_IsOurs(t *testing.T) {
	type fields struct {
		isActive bool
		isOurs   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{fields: fields{isActive: false, isOurs: false}, want: false},
		{fields: fields{isActive: false, isOurs: true}, want: true},
		{fields: fields{isActive: true, isOurs: true}, want: true},
		{fields: fields{isActive: true, isOurs: false}, want: false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("isActive: %t, isOurs: %t", tt.fields.isActive, tt.fields.isOurs), func(t *testing.T) {
			l := lock{
				isActive: tt.fields.isActive,
				isOurs:   tt.fields.isOurs,
			}
			assert.Equalf(t, tt.want, l.IsOurs(), "IsOurs()")
		})
	}
}
