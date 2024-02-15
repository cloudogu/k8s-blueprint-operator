package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	globalConfigKey1     = common.GlobalConfigKey("key1")
	globalConfigKeyEmpty = common.GlobalConfigKey("")
)

func TestGlobalConfig_validate(t *testing.T) {
	t.Run("empty config is ok", func(t *testing.T) {
		config := GlobalConfig{}
		err := config.validate()
		assert.NoError(t, err)
	})
	t.Run("config is ok", func(t *testing.T) {
		config := GlobalConfig{
			Present: map[common.GlobalConfigKey]common.GlobalConfigValue{
				"my/key1": "", //empty values are ok
				"my/key2": "test",
			},
			Absent: []common.GlobalConfigKey{
				"key3",
			},
		}

		err := config.validate()

		assert.NoError(t, err)
	})
	t.Run("no empty present keys", func(t *testing.T) {
		config := GlobalConfig{
			Present: map[common.GlobalConfigKey]common.GlobalConfigValue{
				"": "",
			},
		}

		err := config.validate()

		assert.ErrorContains(t, err, "key for present global config should not be empty")
	})
	t.Run("no empty absent keys", func(t *testing.T) {
		config := GlobalConfig{
			Absent: []common.GlobalConfigKey{""},
		}

		err := config.validate()

		assert.ErrorContains(t, err, "key for absent global config should not be empty")
	})
	t.Run("not present and absent at the same time", func(t *testing.T) {
		config := GlobalConfig{
			Present: map[common.GlobalConfigKey]common.GlobalConfigValue{
				"my/key1": "test",
			},
			Absent: []common.GlobalConfigKey{
				"my/key1",
			},
		}

		err := config.validate()

		assert.ErrorContains(t, err, "config key \"my/key1\" cannot be present and absent at the same time")
	})

	t.Run("combine errors", func(t *testing.T) {
		config := GlobalConfig{
			Present: map[common.GlobalConfigKey]common.GlobalConfigValue{
				"":        "",
				"my/key1": "test",
			},
			Absent: []common.GlobalConfigKey{
				"my/key1",
			},
		}

		err := config.validate()

		assert.ErrorContains(t, err, "key for present global config should not be empty")
		assert.ErrorContains(t, err, "config key \"my/key1\" cannot be present and absent at the same time")
	})
}
