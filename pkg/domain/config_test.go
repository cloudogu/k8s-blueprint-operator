package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestDoguConfig_validate(t *testing.T) {
	t.Run("empty is ok", func(t *testing.T) {
		config := DoguConfig{}
		err := config.validate("dogu1")
		assert.NoError(t, err)
	})
	t.Run("config is ok", func(t *testing.T) {
		config := DoguConfig{
			Present: map[common.DoguConfigKey]common.DoguConfigValue{
				common.DoguConfigKey{DoguName: "dogu1", Key: "my/key1"}: "value1",
			},
			Absent: []common.DoguConfigKey{
				{DoguName: "dogu1", Key: "my/key2"},
			},
		}
		err := config.validate("dogu1")
		assert.NoError(t, err)
	})
	t.Run("not absent and present at the same time", func(t *testing.T) {
		config := DoguConfig{
			Present: map[common.DoguConfigKey]common.DoguConfigValue{
				common.DoguConfigKey{DoguName: "dogu1", Key: "my/key"}: "value1",
			},
			Absent: []common.DoguConfigKey{
				{DoguName: "dogu1", Key: "my/key"},
			},
		}
		err := config.validate("dogu1")
		assert.Error(t, err)
	})
	t.Run("not same key multiple times", func(t *testing.T) {
		config := DoguConfig{
			Absent: []common.DoguConfigKey{
				{DoguName: "dogu1", Key: "my/key"},
				{DoguName: "dogu1", Key: "my/key"},
			},
		}
		err := config.validate("dogu1")
		assert.Error(t, err)
	})
	t.Run("combine errors", func(t *testing.T) {
		config := DoguConfig{
			Present: map[common.DoguConfigKey]common.DoguConfigValue{
				common.DoguConfigKey{DoguName: "dogu1", Key: ""}: "value1",
			},
			Absent: []common.DoguConfigKey{
				{DoguName: "dogu1", Key: ""},
			},
		}
		err := config.validate("dogu1")
		//TODO: test is not needed, if we do not combine errors directly for DoguConfig
		assert.ErrorContains(t, err, "present dogu config key invalid")
		assert.ErrorContains(t, err, "absent dogu config key invalid")
	})
}

func TestSensitiveDoguConfig_validate(t *testing.T) {
	t.Run("empty is ok", func(t *testing.T) {
		config := SensitiveDoguConfig{}
		err := config.validate("")
		assert.NoError(t, err)
	})
	t.Run("config is ok", func(t *testing.T) {
		config := SensitiveDoguConfig{
			Present: map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{
				common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{DoguName: "dogu1", Key: "my/key1"}}: "value1",
			},
			Absent: []common.SensitiveDoguConfigKey{
				{DoguConfigKey: common.DoguConfigKey{DoguName: "dogu1", Key: "my/key2"}},
			},
		}
		err := config.validate("dogu1")
		assert.NoError(t, err)
	})
	t.Run("not absent and present at the same time", func(t *testing.T) {
		config := SensitiveDoguConfig{
			Present: map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{
				common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{DoguName: "dogu1", Key: "my/key"}}: "value1",
			},
			Absent: []common.SensitiveDoguConfigKey{
				{DoguConfigKey: common.DoguConfigKey{DoguName: "dogu1", Key: "my/key"}},
			},
		}
		err := config.validate("dogu1")
		assert.Error(t, err)
	})
	t.Run("not same key multiple times", func(t *testing.T) {
		config := SensitiveDoguConfig{
			Absent: []common.SensitiveDoguConfigKey{
				{DoguConfigKey: common.DoguConfigKey{DoguName: "dogu1", Key: "my/key"}},
				{DoguConfigKey: common.DoguConfigKey{DoguName: "dogu1", Key: "my/key"}},
			},
		}
		err := config.validate("dogu1")
		assert.Error(t, err)
	})
	t.Run("combine errors", func(t *testing.T) {
		config := SensitiveDoguConfig{
			Present: map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{
				common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{DoguName: "dogu1", Key: ""}}: "value1",
			},
			Absent: []common.SensitiveDoguConfigKey{
				{DoguConfigKey: common.DoguConfigKey{DoguName: "dogu1", Key: ""}},
			},
		}
		err := config.validate("dogu1")
		//TODO: test is not needed, if we do not combine errors directly for DoguConfig
		assert.ErrorContains(t, err, "present dogu config key invalid")
		assert.ErrorContains(t, err, "absent dogu config key invalid")
	})
}
