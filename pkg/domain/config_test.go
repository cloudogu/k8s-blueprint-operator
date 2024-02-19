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
	t.Run("not same key multiple times", func(t *testing.T) {
		config := GlobalConfig{
			Absent: []common.GlobalConfigKey{"my/key", "my/key"},
		}
		err := config.validate()
		assert.ErrorContains(t, err, "absent global config should not contain duplicate keys")
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
		assert.ErrorContains(t, err, "key \"my/key\" of dogu \"dogu1\" cannot be present and absent at the same time")
	})
	t.Run("not same key multiple times", func(t *testing.T) {
		config := DoguConfig{
			Absent: []common.DoguConfigKey{
				{DoguName: "dogu1", Key: "my/key"},
				{DoguName: "dogu1", Key: "my/key"},
			},
		}
		err := config.validate("dogu1")
		assert.ErrorContains(t, err, "absent dogu config should not contain duplicate keys: [key \"my/key\" of dogu \"dogu1\"]")
	})
	t.Run("only one referenced dogu name", func(t *testing.T) {
		config := DoguConfig{
			Present: map[common.DoguConfigKey]common.DoguConfigValue{
				common.DoguConfigKey{DoguName: "dogu1", Key: "test"}: "value1",
			},
			Absent: []common.DoguConfigKey{
				{DoguName: "dogu1", Key: "my/key"},
			},
		}
		err := config.validate("dogu2")
		assert.ErrorContains(t, err, "present key \"test\" of dogu \"dogu1\" does not match superordinate dogu name \"dogu2\"")
		assert.ErrorContains(t, err, "absent key \"my/key\" of dogu \"dogu1\" does not match superordinate dogu name \"dogu2\"")
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
		assert.ErrorContains(t, err, "key \"my/key\" of dogu \"dogu1\" cannot be present and absent at the same time")
	})
	t.Run("not same key multiple times", func(t *testing.T) {
		config := SensitiveDoguConfig{
			Absent: []common.SensitiveDoguConfigKey{
				{DoguConfigKey: common.DoguConfigKey{DoguName: "dogu1", Key: "my/key"}},
				{DoguConfigKey: common.DoguConfigKey{DoguName: "dogu1", Key: "my/key"}},
			},
		}
		err := config.validate("dogu1")
		assert.ErrorContains(t, err, "absent dogu config should not contain duplicate keys: [key \"my/key\" of dogu \"dogu1\"]")
	})
	t.Run("only one referenced dogu name", func(t *testing.T) {
		config := SensitiveDoguConfig{
			Present: map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{
				common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{DoguName: "dogu1", Key: "test"}}: "value1",
			},
			Absent: []common.SensitiveDoguConfigKey{
				{DoguConfigKey: common.DoguConfigKey{DoguName: "dogu1", Key: "my/key"}},
			},
		}
		err := config.validate("dogu2")
		assert.ErrorContains(t, err, "present key \"test\" of dogu \"dogu1\" does not match superordinate dogu name \"dogu2\"")
		assert.ErrorContains(t, err, "absent key \"my/key\" of dogu \"dogu1\" does not match superordinate dogu name \"dogu2\"")
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
		assert.ErrorContains(t, err, "present dogu config key invalid")
		assert.ErrorContains(t, err, "absent dogu config key invalid")
	})
}

func TestConfig_validate(t *testing.T) {
	t.Run("succeed for empty config", func(t *testing.T) {
		// given
		sut := Config{}

		// when
		err := sut.validate()

		// then
		assert.NoError(t, err)
	})
	t.Run("fail if dogu name in dogu config does not match dogu key", func(t *testing.T) {
		// given
		sut := Config{
			Dogus: map[common.SimpleDoguName]CombinedDoguConfig{
				"some-name": {DoguName: "another-name"},
			},
		}

		// when
		err := sut.validate()

		// then
		assert.ErrorContains(t, err, "dogu name \"some-name\" in map and dogu name \"another-name\" in value are not equal")
	})
	t.Run("fail with multiple errors", func(t *testing.T) {
		// given
		sut := Config{
			Dogus: map[common.SimpleDoguName]CombinedDoguConfig{
				"some-name": {
					DoguName: "another-name",
					Config: DoguConfig{
						Absent: []common.DoguConfigKey{{DoguName: ""}},
					},
					SensitiveConfig: SensitiveDoguConfig{
						Absent: []common.SensitiveDoguConfigKey{{common.DoguConfigKey{DoguName: ""}}},
					},
				},
			},
			Global: GlobalConfig{Absent: []common.GlobalConfigKey{""}},
		}

		// when
		err := sut.validate()

		// then
		assert.ErrorContains(t, err, "dogu name \"some-name\" in map and dogu name \"another-name\" in value are not equal")
		assert.ErrorContains(t, err, "config for dogu \"another-name\" invalid")
		assert.ErrorContains(t, err, "sensitive config for dogu \"another-name\" invalid")
		assert.ErrorContains(t, err, "key for absent global config should not be empty")
	})
}