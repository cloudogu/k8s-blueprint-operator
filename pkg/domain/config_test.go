package domain

import (
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	libconfig "github.com/cloudogu/k8s-registry-lib/config"
	"github.com/stretchr/testify/assert"
)

var confgiVal1 = libconfig.Value("value1")

func TestGlobalConfig_validate(t *testing.T) {
	t.Run("empty config is ok", func(t *testing.T) {
		config := GlobalConfigEntries{}
		err := config.validate()
		assert.NoError(t, err)
	})
	t.Run("config is ok", func(t *testing.T) {
		config := GlobalConfigEntries{
			{
				Key:   "my/key1",
				Value: nil, //empty values are ok
			},
			{
				Key:   "my/key2",
				Value: &confgiVal1,
			},
			{
				Key:    "key3",
				Absent: true,
			},
		}

		err := config.validate()

		assert.NoError(t, err)
	})
	t.Run("no empty present keys", func(t *testing.T) {
		config := GlobalConfigEntries{
			{
				Key:   "",
				Value: nil,
			},
		}

		err := config.validate()

		assert.ErrorContains(t, err, "key for global config should not be empty")
	})
	t.Run("no empty absent keys", func(t *testing.T) {
		config := GlobalConfigEntries{
			{
				Key:    "",
				Absent: true,
			},
		}

		err := config.validate()

		assert.ErrorContains(t, err, "key for global config should not be empty")
	})
	t.Run("not present and absent at the same time", func(t *testing.T) {
		config := GlobalConfigEntries{
			{
				Key:   "my/key1",
				Value: &confgiVal1,
			},
			{
				Key:    "my/key1",
				Absent: true,
			},
		}

		err := config.validate()

		assert.ErrorContains(t, err, "duplicate dogu config Key found: my/key1")
	})

	t.Run("combine errors", func(t *testing.T) {
		config := GlobalConfigEntries{
			{
				Key:   "",
				Value: &confgiVal1,
			},
			{
				Key:   "my/key1",
				Value: &confgiVal1,
			},
			{
				Key:    "my/key1",
				Absent: true,
			},
		}

		err := config.validate()

		assert.ErrorContains(t, err, "key for global config should not be empty")
		assert.ErrorContains(t, err, "duplicate dogu config Key found: my/key1")
	})
}

func TestDoguConfig_validate(t *testing.T) {
	t.Run("empty is ok", func(t *testing.T) {
		config := DoguConfigEntries{}
		err := config.validate("dogu1")
		assert.NoError(t, err)
	})
	t.Run("config is ok", func(t *testing.T) {
		config := DoguConfigEntries{
			{
				Key:   "my/key1",
				Value: &confgiVal1,
			},
			{
				Key:    "my/key2",
				Absent: true,
			},
		}
		err := config.validate("dogu1")
		assert.NoError(t, err)
	})
	t.Run("not absent and present at the same time", func(t *testing.T) {
		config := DoguConfigEntries{
			{
				Key:   "my/key1",
				Value: &confgiVal1,
			},
			{
				Key:    "my/key1",
				Absent: true,
			},
		}
		err := config.validate("dogu1")
		assert.ErrorContains(t, err, "duplicate dogu config Key found: my/key1")
	})
	t.Run("not same key multiple times", func(t *testing.T) {
		config := DoguConfigEntries{
			{
				Key:   "my/key1",
				Value: &confgiVal1,
			},
			{
				Key:   "my/key1",
				Value: &confgiVal1,
			},
		}
		err := config.validate("dogu1")
		assert.ErrorContains(t, err, "duplicate dogu config Key found: my/key1")
	})
	t.Run("no empty present keys", func(t *testing.T) {
		config := DoguConfigEntries{
			{
				Key:   "",
				Value: nil,
			},
		}

		err := config.validate("dogu1")

		assert.ErrorContains(t, err, "key for config should not be empty")
	})
	t.Run("no empty absent keys", func(t *testing.T) {
		config := DoguConfigEntries{
			{
				Key:    "",
				Absent: true,
			},
		}

		err := config.validate("dogu1")

		assert.ErrorContains(t, err, "key for config should not be empty")
	})
	t.Run("combine errors", func(t *testing.T) {
		config := DoguConfigEntries{
			{
				Key:   "",
				Value: &confgiVal1,
			},
			{
				Key:   "my/key1",
				Value: &confgiVal1,
			},
			{
				Key:    "my/key1",
				Absent: true,
			},
		}
		err := config.validate("dogu1")
		assert.ErrorContains(t, err, "config for dogu \"dogu1\" is invalid")
		assert.ErrorContains(t, err, "key for config should not be empty")
		assert.ErrorContains(t, err, "duplicate dogu config Key found: my/key1")
	})
	t.Run("No absent with value", func(t *testing.T) {
		config := DoguConfigEntries{
			{
				Key:    "my/key1",
				Value:  &confgiVal1,
				Absent: true,
			},
		}
		err := config.validate("dogu1")
		assert.ErrorContains(t, err, "absent entries cannot have value or secretRef")
	})
	t.Run("No absent with secret", func(t *testing.T) {
		config := DoguConfigEntries{
			{
				Key:       "my/key1",
				SecretRef: &SensitiveValueRef{},
				Absent:    true,
			},
		}
		err := config.validate("dogu1")
		assert.ErrorContains(t, err, "absent entries cannot have value or secretRef")
	})
	t.Run("No secret and value", func(t *testing.T) {
		config := DoguConfigEntries{
			{
				Key:       "my/key1",
				SecretRef: &SensitiveValueRef{},
				Value:     &confgiVal1,
			},
		}
		err := config.validate("dogu1")
		assert.ErrorContains(t, err, "config entries can have either a value or a secretRef")
	})
	t.Run("No secret without sensitive", func(t *testing.T) {
		config := DoguConfigEntries{
			{
				Key:       "my/key1",
				SecretRef: &SensitiveValueRef{},
			},
		}
		err := config.validate("dogu1")
		assert.ErrorContains(t, err, "config entries with secret references have to be sensitive")
	})
	t.Run("secret with sensitive allowed", func(t *testing.T) {
		config := DoguConfigEntries{
			{
				Key:       "my/key1",
				SecretRef: &SensitiveValueRef{},
				Sensitive: true,
			},
		}
		err := config.validate("dogu1")
		assert.NoError(t, err)
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
	t.Run("fail with multiple errors", func(t *testing.T) {
		// given
		sut := Config{
			Dogus: map[cescommons.SimpleName]DoguConfigEntries{
				"dogu1": {
					ConfigEntry{
						Key:   "",
						Value: &confgiVal1,
					},
				},
				"dogu2": {
					ConfigEntry{
						Key:   "myKey1",
						Value: &confgiVal1,
					},
					ConfigEntry{
						Key:    "myKey1",
						Absent: true,
					},
				},
			},
			Global: GlobalConfigEntries{
				ConfigEntry{
					Key:    "myKey",
					Value:  &confgiVal1,
					Absent: true,
				},
			},
		}

		// when
		err := sut.validate()

		// then
		assert.ErrorContains(t, err, "config for dogu \"dogu1\" is invalid")
		assert.ErrorContains(t, err, "key for config should not be empty")
		assert.ErrorContains(t, err, "config for dogu \"dogu2\" is invalid")
		assert.ErrorContains(t, err, "duplicate dogu config Key found: myKey1")
		assert.ErrorContains(t, err, "global config is invalid")
		assert.ErrorContains(t, err, "absent entries cannot have value")
	})
}

func TestGlobalConfig_GetGlobalConfigKeys(t *testing.T) {
	var (
		globalKey1 = common.GlobalConfigKey("key1")
		globalKey2 = common.GlobalConfigKey("key2")
	)
	config := GlobalConfigEntries{
		{
			Key:   globalKey1,
			Value: &confgiVal1,
		},
		{
			Key:    globalKey2,
			Absent: true,
		},
	}

	keys := config.GetGlobalConfigKeys()

	assert.ElementsMatch(t, keys, []common.GlobalConfigKey{globalKey1, globalKey2})
}

func TestConfig_GetDoguConfigKeys(t *testing.T) {
	var (
		nginx       = cescommons.SimpleName("nginx")
		postfix     = cescommons.SimpleName("postfix")
		nginxKey1   = common.DoguConfigKey{DoguName: nginx, Key: "key1"}
		nginxKey2   = common.DoguConfigKey{DoguName: nginx, Key: "key2"}
		postfixKey1 = common.DoguConfigKey{DoguName: postfix, Key: "key1"}
		postfixKey2 = common.DoguConfigKey{DoguName: postfix, Key: "key2"}
		postfixKey3 = common.DoguConfigKey{DoguName: postfix, Key: "key3"}
	)
	config := Config{
		Dogus: map[cescommons.SimpleName]DoguConfigEntries{
			nginx: {
				{
					Key:   nginxKey1.Key,
					Value: &confgiVal1,
				},
				{
					Key:    nginxKey2.Key,
					Absent: true,
				},
			},
			postfix: {
				{
					Key:   postfixKey1.Key,
					Value: &confgiVal1,
				},
				{
					Key:    postfixKey2.Key,
					Absent: true,
				},
				{
					Key:       postfixKey3.Key,
					Absent:    true,
					Sensitive: true,
				},
			},
		},
	}

	keys := config.GetDoguConfigKeys()

	assert.ElementsMatch(t, keys, []common.DoguConfigKey{nginxKey1, nginxKey2, postfixKey1, postfixKey2})
}

func TestConfig_GetSensitiveDoguConfigKeys(t *testing.T) {
	var (
		nginx       = cescommons.SimpleName("nginx")
		postfix     = cescommons.SimpleName("postfix")
		nginxKey1   = common.DoguConfigKey{DoguName: nginx, Key: "key1"}
		nginxKey2   = common.DoguConfigKey{DoguName: nginx, Key: "key2"}
		postfixKey1 = common.DoguConfigKey{DoguName: postfix, Key: "key1"}
		postfixKey2 = common.DoguConfigKey{DoguName: postfix, Key: "key2"}
		postfixKey3 = common.DoguConfigKey{DoguName: postfix, Key: "key3"}
	)
	config := Config{
		Dogus: map[cescommons.SimpleName]DoguConfigEntries{
			nginx: {
				{
					Key:       nginxKey1.Key,
					Value:     &confgiVal1,
					Sensitive: true,
				},
				{
					Key:       nginxKey2.Key,
					Absent:    true,
					Sensitive: true,
				},
			},
			postfix: {
				{
					Key:       postfixKey1.Key,
					Value:     &confgiVal1,
					Sensitive: true,
				},
				{
					Key:       postfixKey2.Key,
					Absent:    true,
					Sensitive: true,
				},
				{
					Key:    postfixKey3.Key,
					Absent: true,
				},
			},
		},
	}

	keys := config.GetSensitiveDoguConfigKeys()

	assert.ElementsMatch(t, keys, []common.DoguConfigKey{nginxKey1, nginxKey2, postfixKey1, postfixKey2})
}

func TestConfig_GetDogusWithChangedConfig(t *testing.T) {
	doguConfig := DoguConfigEntries{
		{
			Key:   dogu1Key1.Key,
			Value: &confgiVal1,
		},
		{
			Key:    dogu1Key2.Key,
			Absent: true,
		},
	}
	configEntryPresent := ConfigEntry{
		Key:   dogu1Key1.Key,
		Value: &confgiVal1,
	}
	configEntryAbsent := ConfigEntry{
		Key:    dogu1Key2.Key,
		Absent: true,
	}

	type args struct {
		doguConfig      DoguConfigEntries
		withDogu2Change bool
	}

	var emptyResult []cescommons.SimpleName
	var tests = []struct {
		name string
		args args
		want []cescommons.SimpleName
	}{
		{
			name: "should get multiple Dogus",
			args: args{doguConfig: DoguConfigEntries{configEntryPresent, configEntryAbsent}, withDogu2Change: true},
			want: []cescommons.SimpleName{dogu1, dogu2},
		},
		{
			name: "should get Dogus with changed present and absent config",
			args: args{doguConfig: DoguConfigEntries{configEntryPresent, configEntryAbsent}},
			want: []cescommons.SimpleName{dogu1},
		},
		{
			name: "should get Dogus with changed present config",
			args: args{doguConfig: DoguConfigEntries{configEntryPresent}},
			want: []cescommons.SimpleName{dogu1},
		},
		{
			name: "should get Dogus with changed absent config",
			args: args{doguConfig: DoguConfigEntries{configEntryAbsent}},
			want: []cescommons.SimpleName{dogu1},
		},
		{
			name: "should not get Dogus with no config changes",
			args: args{doguConfig: DoguConfigEntries{}},
			want: emptyResult,
		},
		{
			name: "should not get Dogus with sensitive config changes",
			args: args{doguConfig: DoguConfigEntries{ConfigEntry{Key: "mykey", Value: &confgiVal1, Sensitive: true}}},
			want: emptyResult,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				Dogus: map[cescommons.SimpleName]DoguConfigEntries{
					dogu1: tt.args.doguConfig,
				},
				Global: GlobalConfigEntries{},
			}

			if tt.args.withDogu2Change {
				config.Dogus[dogu2] = doguConfig
			}

			changedDogus := config.GetDogusWithChangedConfig()
			assert.Len(t, changedDogus, len(tt.want))
			for _, doguName := range tt.want {
				assert.Contains(t, changedDogus, doguName)
			}
		})
	}
}

func TestConfig_GetDogusWithChangedSensitiveConfig(t *testing.T) {
	presentConfig := ConfigEntry{
		Key:       dogu1Key1.Key,
		Sensitive: true,
		SecretRef: &SensitiveValueRef{
			SecretName: "mySecret",
			SecretKey:  "myKey",
		},
	}
	absentConfig := ConfigEntry{
		Key:       dogu1Key1.Key,
		Sensitive: true,
		Absent:    true,
	}

	type args struct {
		doguConfig      DoguConfigEntries
		withDogu2Change bool
	}

	var emptyResult []cescommons.SimpleName
	var tests = []struct {
		name string
		args args
		want []cescommons.SimpleName
	}{
		{
			name: "should get multiple Dogus",
			args: args{doguConfig: DoguConfigEntries{presentConfig, absentConfig}, withDogu2Change: true},
			want: []cescommons.SimpleName{dogu1, dogu2},
		},
		{
			name: "should get Dogus with changed present and absent config",
			args: args{doguConfig: DoguConfigEntries{presentConfig, absentConfig}},
			want: []cescommons.SimpleName{dogu1},
		},
		{
			name: "should get Dogus with changed present config",
			args: args{doguConfig: DoguConfigEntries{presentConfig}},
			want: []cescommons.SimpleName{dogu1},
		},
		{
			name: "should get Dogus with changed absent config",
			args: args{doguConfig: DoguConfigEntries{absentConfig}},
			want: []cescommons.SimpleName{dogu1},
		},
		{
			name: "should not get Dogus with no config changes",
			args: args{doguConfig: DoguConfigEntries{}},
			want: emptyResult,
		},
		{
			name: "should not get Dogus with no sensitive config changes",
			args: args{doguConfig: DoguConfigEntries{ConfigEntry{Key: "mykey", Value: &confgiVal1}}},
			want: emptyResult,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				Dogus: map[cescommons.SimpleName]DoguConfigEntries{
					dogu1: tt.args.doguConfig,
				},
				Global: GlobalConfigEntries{},
			}

			if tt.args.withDogu2Change {
				config.Dogus[dogu2] = tt.args.doguConfig
			}

			changedDogus := config.GetDogusWithChangedSensitiveConfig()
			assert.Len(t, changedDogus, len(tt.want))
			for _, doguName := range tt.want {
				assert.Contains(t, changedDogus, doguName)
			}
		})
	}
}
