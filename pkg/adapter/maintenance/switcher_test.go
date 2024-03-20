package maintenance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

func Test_defaultSwitcher_Activate(t *testing.T) {
	t.Run("should fail to activate maintenance mode", func(t *testing.T) {
		// given
		expectedJson := `{"title":"myTitle","text":"myText","holder":"k8s-blueprint-operator"}`
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Set("maintenance", expectedJson).Return(assert.AnError)

		sut := &defaultSwitcher{globalConfig: globalConfigMock}

		// when
		err := sut.activate(domainservice.MaintenancePageModel{
			Title: "myTitle",
			Text:  "myText",
		})

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to set maintenance mode registry key")
	})
	t.Run("should succeed to activate maintenance mode", func(t *testing.T) {
		// given
		expectedJson := `{"title":"myTitle","text":"myText","holder":"k8s-blueprint-operator"}`
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Set("maintenance", expectedJson).Return(nil)

		sut := &defaultSwitcher{globalConfig: globalConfigMock}

		// when
		err := sut.activate(domainservice.MaintenancePageModel{
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

		sut := &defaultSwitcher{globalConfig: globalConfigMock}

		// when
		err := sut.deactivate()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to delete maintenance mode registry key")
	})
	t.Run("should succeed to deactivate maintenance mode", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfig(t)
		globalConfigMock.EXPECT().Delete("maintenance").Return(nil)

		sut := &defaultSwitcher{globalConfig: globalConfigMock}

		// when
		err := sut.deactivate()

		// then
		require.NoError(t, err)
	})
}
