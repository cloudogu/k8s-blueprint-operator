package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	bpv2 "github.com/cloudogu/blueprint-lib/v2"
)

func TestGlobalConfig_validate(t *testing.T) {
	t.Run("empty config is ok", func(t *testing.T) {
		config := bpv2.GlobalConfig{}
		err := config.validate()
		assert.NoError(t, err)
	})
	t.Run("config is ok", func(t *testing.T) {
		config := bpv2.GlobalConfig{
			Present: map[bpv2.GlobalConfigKey]bpv2.GlobalConfigValue{
				"my/key1": "", //empty values are ok
				"my/key2": "test",
			},
			Absent: []bpv2.GlobalConfigKey{
				"key3",
			},
		}

		err := config.validate()

		assert.NoError(t, err)
	})
	t.Run("no empty present keys", func(t *testing.T) {
		config := bpv2.GlobalConfig{
			Present: map[bpv2.GlobalConfigKey]bpv2.GlobalConfigValue{
				"": "",
			},
		}

		err := config.validate()

		assert.ErrorContains(t, err, "key for present global config should not be empty")
	})
	t.Run("no empty absent keys", func(t *testing.T) {
		config := bpv2.GlobalConfig{
			Absent: []bpv2.GlobalConfigKey{""},
		}

		err := config.validate()

		assert.ErrorContains(t, err, "key for absent global config should not be empty")
	})
	t.Run("not present and absent at the same time", func(t *testing.T) {
		config := bpv2.GlobalConfig{
			Present: map[bpv2.GlobalConfigKey]bpv2.GlobalConfigValue{
				"my/key1": "test",
			},
			Absent: []bpv2.GlobalConfigKey{
				"my/key1",
			},
		}

		err := config.validate()

		assert.ErrorContains(t, err, "config key \"my/key1\" cannot be present and absent at the same time")
	})

	t.Run("combine errors", func(t *testing.T) {
		config := bpv2.GlobalConfig{
			Present: map[bpv2.GlobalConfigKey]bpv2.GlobalConfigValue{
				"":        "",
				"my/key1": "test",
			},
			Absent: []bpv2.GlobalConfigKey{
				"my/key1",
			},
		}

		err := config.validate()

		assert.ErrorContains(t, err, "key for present global config should not be empty")
		assert.ErrorContains(t, err, "config key \"my/key1\" cannot be present and absent at the same time")
	})
	t.Run("not same key multiple times", func(t *testing.T) {
		config := bpv2.GlobalConfig{
			Absent: []bpv2.GlobalConfigKey{"my/key", "my/key"},
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
			Present: map[bpv2.DoguConfigKey]bpv2.DoguConfigValue{
				bpv2.DoguConfigKey{DoguName: "dogu1", Key: "my/key1"}: "value1",
			},
			Absent: []bpv2.DoguConfigKey{
				{DoguName: "dogu1", Key: "my/key2"},
			},
		}
		err := config.validate("dogu1")
		assert.NoError(t, err)
	})
	t.Run("not absent and present at the same time", func(t *testing.T) {
		config := DoguConfig{
			Present: map[bpv2.DoguConfigKey]bpv2.DoguConfigValue{
				bpv2.DoguConfigKey{DoguName: "dogu1", Key: "my/key"}: "value1",
			},
			Absent: []bpv2.DoguConfigKey{
				{DoguName: "dogu1", Key: "my/key"},
			},
		}
		err := config.validate("dogu1")
		assert.ErrorContains(t, err, "key \"my/key\" of dogu \"dogu1\" cannot be present and absent at the same time")
	})
	t.Run("not same key multiple times", func(t *testing.T) {
		config := DoguConfig{
			Absent: []bpv2.DoguConfigKey{
				{DoguName: "dogu1", Key: "my/key"},
				{DoguName: "dogu1", Key: "my/key"},
			},
		}
		err := config.validate("dogu1")
		assert.ErrorContains(t, err, "absent dogu config should not contain duplicate keys: [key \"my/key\" of dogu \"dogu1\"]")
	})
	t.Run("only one referenced dogu name", func(t *testing.T) {
		config := DoguConfig{
			Present: map[bpv2.DoguConfigKey]bpv2.DoguConfigValue{
				bpv2.DoguConfigKey{DoguName: "dogu1", Key: "test"}: "value1",
			},
			Absent: []bpv2.DoguConfigKey{
				{DoguName: "dogu1", Key: "my/key"},
			},
		}
		err := config.validate("dogu2")
		assert.ErrorContains(t, err, "present key \"test\" of dogu \"dogu1\" does not match superordinate dogu name \"dogu2\"")
		assert.ErrorContains(t, err, "absent key \"my/key\" of dogu \"dogu1\" does not match superordinate dogu name \"dogu2\"")
	})
	t.Run("combine errors", func(t *testing.T) {
		config := DoguConfig{
			Present: map[bpv2.DoguConfigKey]bpv2.DoguConfigValue{
				bpv2.DoguConfigKey{DoguName: "dogu1", Key: ""}: "value1",
			},
			Absent: []bpv2.DoguConfigKey{
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
			Present: map[bpv2.SensitiveDoguConfigKey]bpv2.SensitiveDoguConfigValue{
				bpv2.SensitiveDoguConfigKey{DoguName: "dogu1", Key: "my/key1"}: "value1",
			},
			Absent: []bpv2.SensitiveDoguConfigKey{
				{DoguName: "dogu1", Key: "my/key2"},
			},
		}
		err := config.validate("dogu1")
		assert.NoError(t, err)
	})
	t.Run("not absent and present at the same time", func(t *testing.T) {
		config := SensitiveDoguConfig{
			Present: map[bpv2.SensitiveDoguConfigKey]bpv2.SensitiveDoguConfigValue{
				bpv2.SensitiveDoguConfigKey{DoguName: "dogu1", Key: "my/key"}: "value1",
			},
			Absent: []bpv2.SensitiveDoguConfigKey{
				{DoguName: "dogu1", Key: "my/key"},
			},
		}
		err := config.validate("dogu1")
		assert.ErrorContains(t, err, "key \"my/key\" of dogu \"dogu1\" cannot be present and absent at the same time")
	})
	t.Run("not same key multiple times", func(t *testing.T) {
		config := SensitiveDoguConfig{
			Absent: []bpv2.SensitiveDoguConfigKey{
				{DoguName: "dogu1", Key: "my/key"},
				{DoguName: "dogu1", Key: "my/key"},
			},
		}
		err := config.validate("dogu1")
		assert.ErrorContains(t, err, "absent dogu config should not contain duplicate keys: [key \"my/key\" of dogu \"dogu1\"]")
	})
	t.Run("only one referenced dogu name", func(t *testing.T) {
		config := SensitiveDoguConfig{
			Present: map[bpv2.SensitiveDoguConfigKey]bpv2.SensitiveDoguConfigValue{
				bpv2.SensitiveDoguConfigKey{DoguName: "dogu1", Key: "test"}: "value1",
			},
			Absent: []bpv2.SensitiveDoguConfigKey{
				{DoguName: "dogu1", Key: "my/key"},
			},
		}
		err := config.validate("dogu2")
		assert.ErrorContains(t, err, "present key \"test\" of dogu \"dogu1\" does not match superordinate dogu name \"dogu2\"")
		assert.ErrorContains(t, err, "absent key \"my/key\" of dogu \"dogu1\" does not match superordinate dogu name \"dogu2\"")
	})
	t.Run("combine errors", func(t *testing.T) {
		config := SensitiveDoguConfig{
			Present: map[bpv2.SensitiveDoguConfigKey]bpv2.SensitiveDoguConfigValue{
				bpv2.SensitiveDoguConfigKey{DoguName: "dogu1", Key: ""}: "value1",
			},
			Absent: []bpv2.SensitiveDoguConfigKey{
				{DoguName: "dogu1", Key: ""},
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
			Dogus: map[cesbpv2s.SimpleName]CombinedDoguConfig{
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
			Dogus: map[cesbpv2s.SimpleName]CombinedDoguConfig{
				"some-name": {
					DoguName: "another-name",
					Config: DoguConfig{
						Absent: []bpv2.DoguConfigKey{{DoguName: ""}},
					},
					SensitiveConfig: SensitiveDoguConfig{
						Absent: []bpv2.SensitiveDoguConfigKey{{DoguName: ""}},
					},
				},
			},
			Global: bpv2.GlobalConfig{Absent: []bpv2.GlobalConfigKey{""}},
		}

		// when
		err := sut.validate()

		// then
		assert.ErrorContains(t, err, "dogu name \"some-name\" in map and dogu name \"another-name\" in value are not equal")
		assert.ErrorContains(t, err, "config for dogu \"another-name\" is invalid")
		assert.ErrorContains(t, err, "key for absent global config should not be empty")
	})
}

func TestGlobalConfig_GetGlobalConfigKeys(t *testing.T) {
	var (
		globalKey1 = bpv2.GlobalConfigKey("key1")
		globalKey2 = bpv2.GlobalConfigKey("key2")
	)
	config := bpv2.GlobalConfig{
		Present: map[bpv2.GlobalConfigKey]bpv2.GlobalConfigValue{
			globalKey1: "value",
		},
		Absent: []bpv2.GlobalConfigKey{
			globalKey2,
		},
	}

	keys := config.GetGlobalConfigKeys()

	assert.ElementsMatch(t, keys, []bpv2.GlobalConfigKey{globalKey1, globalKey2})
}

func TestConfig_GetDoguConfigKeys(t *testing.T) {
	var (
		nginx       = cesbpv2s.SimpleName("nginx")
		postfix     = cesbpv2s.SimpleName("postfix")
		nginxKey1   = bpv2.DoguConfigKey{DoguName: nginx, Key: "key1"}
		nginxKey2   = bpv2.DoguConfigKey{DoguName: nginx, Key: "key2"}
		postfixKey1 = bpv2.DoguConfigKey{DoguName: postfix, Key: "key1"}
		postfixKey2 = bpv2.DoguConfigKey{DoguName: postfix, Key: "key2"}
	)
	config := Config{
		Dogus: map[cesbpv2s.SimpleName]CombinedDoguConfig{
			nginx: {
				DoguName: nginx,
				Config: DoguConfig{
					Present: map[bpv2.DoguConfigKey]bpv2.DoguConfigValue{
						nginxKey1: "value",
					},
					Absent: []bpv2.DoguConfigKey{
						nginxKey2,
					},
				},
				SensitiveConfig: SensitiveDoguConfig{},
			},
			postfix: {
				DoguName: postfix,
				Config: DoguConfig{
					Present: map[bpv2.DoguConfigKey]bpv2.DoguConfigValue{
						postfixKey1: "value",
					},
					Absent: []bpv2.DoguConfigKey{
						postfixKey2,
					},
				},
				SensitiveConfig: SensitiveDoguConfig{},
			},
		},
	}

	keys := config.GetDoguConfigKeys()

	assert.ElementsMatch(t, keys, []bpv2.DoguConfigKey{nginxKey1, nginxKey2, postfixKey1, postfixKey2})
}

func TestConfig_GetSensitiveDoguConfigKeys(t *testing.T) {
	var (
		nginx       = cesbpv2s.SimpleName("nginx")
		postfix     = cesbpv2s.SimpleName("postfix")
		nginxKey1   = bpv2.SensitiveDoguConfigKey{DoguName: nginx, Key: "key1"}
		nginxKey2   = bpv2.SensitiveDoguConfigKey{DoguName: nginx, Key: "key2"}
		postfixKey1 = bpv2.SensitiveDoguConfigKey{DoguName: postfix, Key: "key1"}
		postfixKey2 = bpv2.SensitiveDoguConfigKey{DoguName: postfix, Key: "key2"}
	)
	config := Config{
		Dogus: map[cesbpv2s.SimpleName]CombinedDoguConfig{
			nginx: {
				DoguName: nginx,
				SensitiveConfig: SensitiveDoguConfig{
					Present: map[bpv2.SensitiveDoguConfigKey]bpv2.SensitiveDoguConfigValue{
						nginxKey1: "value",
					},
					Absent: []bpv2.SensitiveDoguConfigKey{
						nginxKey2,
					},
				},
			},
			postfix: {
				DoguName: postfix,
				SensitiveConfig: SensitiveDoguConfig{
					Present: map[bpv2.SensitiveDoguConfigKey]bpv2.SensitiveDoguConfigValue{
						postfixKey1: "value",
					},
					Absent: []bpv2.SensitiveDoguConfigKey{
						postfixKey2,
					},
				},
			},
		},
	}

	keys := config.GetSensitiveDoguConfigKeys()

	assert.ElementsMatch(t, keys, []bpv2.SensitiveDoguConfigKey{nginxKey1, nginxKey2, postfixKey1, postfixKey2})
}

func TestCombinedDoguConfig_validate(t *testing.T) {
	normalConfig := DoguConfig{
		Present: map[bpv2.DoguConfigKey]bpv2.DoguConfigValue{
			bpv2.DoguConfigKey{DoguName: "dogu1", Key: "my/key1"}: "value1",
		},
		Absent: []bpv2.DoguConfigKey{
			{DoguName: "dogu1", Key: "my/key2"},
		},
	}
	sensitiveConfig := SensitiveDoguConfig{
		Present: map[bpv2.SensitiveDoguConfigKey]bpv2.SensitiveDoguConfigValue{
			bpv2.SensitiveDoguConfigKey{DoguName: "dogu1", Key: "my/key1"}: "value1",
		},
		Absent: []bpv2.SensitiveDoguConfigKey{
			{DoguName: "dogu1", Key: "my/key2"},
		},
	}

	config := CombinedDoguConfig{
		DoguName:        "dogu1",
		Config:          normalConfig,
		SensitiveConfig: sensitiveConfig,
	}

	err := config.validate()

	assert.ErrorContains(t, err, "dogu config key key \"my/key1\" of dogu \"dogu1\" cannot be in normal and sensitive configuration at the same time")
}

func TestConfig_GetDogusWithChangedConfig(t *testing.T) {
	presentConfig := map[bpv2.DoguConfigKey]bpv2.DoguConfigValue{
		dogu1Key1: "val",
	}
	AbsentConfig := []bpv2.DoguConfigKey{
		dogu1Key1,
	}
	emptyPresentConfig := map[bpv2.DoguConfigKey]bpv2.DoguConfigValue{}
	var emptyAbsentConfig []bpv2.DoguConfigKey

	type args struct {
		doguConfig      DoguConfig
		withDogu2Change bool
	}

	var emptyResult []cesbpv2s.SimpleName
	var tests = []struct {
		name string
		args args
		want []cesbpv2s.SimpleName
	}{
		{
			name: "should get multiple Dogus",
			args: args{doguConfig: DoguConfig{Present: presentConfig, Absent: AbsentConfig}, withDogu2Change: true},
			want: []cesbpv2s.SimpleName{dogu1, dogu2},
		},
		{
			name: "should get Dogus with changed present and absent config",
			args: args{doguConfig: DoguConfig{Present: presentConfig, Absent: AbsentConfig}},
			want: []cesbpv2s.SimpleName{dogu1},
		},
		{
			name: "should get Dogus with changed present config",
			args: args{doguConfig: DoguConfig{Present: presentConfig, Absent: emptyAbsentConfig}},
			want: []cesbpv2s.SimpleName{dogu1},
		},
		{
			name: "should get Dogus with changed absent config",
			args: args{doguConfig: DoguConfig{Present: emptyPresentConfig, Absent: AbsentConfig}},
			want: []cesbpv2s.SimpleName{dogu1},
		},
		{
			name: "should not get Dogus with no config changes",
			args: args{doguConfig: DoguConfig{Present: emptyPresentConfig, Absent: emptyAbsentConfig}},
			want: emptyResult,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emptyDoguConfig := struct {
				Present map[bpv2.DoguConfigKey]bpv2.DoguConfigValue
				Absent  []bpv2.DoguConfigKey
			}{}
			config := Config{
				Dogus: map[cesbpv2s.SimpleName]CombinedDoguConfig{
					dogu1: {
						DoguName:        dogu1,
						Config:          tt.args.doguConfig,
						SensitiveConfig: emptyDoguConfig,
					},
				},
				Global: bpv2.GlobalConfig{},
			}

			if tt.args.withDogu2Change {
				config.Dogus[dogu2] = CombinedDoguConfig{
					DoguName:        dogu2,
					Config:          tt.args.doguConfig,
					SensitiveConfig: emptyDoguConfig,
				}
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
	presentConfig := map[bpv2.DoguConfigKey]bpv2.DoguConfigValue{
		dogu1Key1: "val",
	}
	AbsentConfig := []bpv2.DoguConfigKey{
		dogu1Key1,
	}
	emptyPresentConfig := map[bpv2.DoguConfigKey]bpv2.DoguConfigValue{}
	var emptyAbsentConfig []bpv2.DoguConfigKey

	type args struct {
		doguConfig      DoguConfig
		withDogu2Change bool
	}

	var emptyResult []cesbpv2s.SimpleName
	var tests = []struct {
		name string
		args args
		want []cesbpv2s.SimpleName
	}{
		{
			name: "should get multiple Dogus",
			args: args{doguConfig: DoguConfig{Present: presentConfig, Absent: AbsentConfig}, withDogu2Change: true},
			want: []cesbpv2s.SimpleName{dogu1, dogu2},
		},
		{
			name: "should get Dogus with changed present and absent config",
			args: args{doguConfig: DoguConfig{Present: presentConfig, Absent: AbsentConfig}},
			want: []cesbpv2s.SimpleName{dogu1},
		},
		{
			name: "should get Dogus with changed present config",
			args: args{doguConfig: DoguConfig{Present: presentConfig, Absent: emptyAbsentConfig}},
			want: []cesbpv2s.SimpleName{dogu1},
		},
		{
			name: "should get Dogus with changed absent config",
			args: args{doguConfig: DoguConfig{Present: emptyPresentConfig, Absent: AbsentConfig}},
			want: []cesbpv2s.SimpleName{dogu1},
		},
		{
			name: "should not get Dogus with no config changes",
			args: args{doguConfig: DoguConfig{Present: emptyPresentConfig, Absent: emptyAbsentConfig}},
			want: emptyResult,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emptyDoguConfig := struct {
				Present map[bpv2.DoguConfigKey]bpv2.DoguConfigValue
				Absent  []bpv2.DoguConfigKey
			}{}
			config := Config{
				Dogus: map[cesbpv2s.SimpleName]CombinedDoguConfig{
					dogu1: {
						DoguName:        dogu1,
						Config:          emptyDoguConfig,
						SensitiveConfig: tt.args.doguConfig,
					},
				},
				Global: bpv2.GlobalConfig{},
			}

			if tt.args.withDogu2Change {
				config.Dogus[dogu2] = CombinedDoguConfig{
					DoguName:        dogu2,
					Config:          emptyDoguConfig,
					SensitiveConfig: tt.args.doguConfig,
				}
			}

			changedDogus := config.GetDogusWithChangedSensitiveConfig()
			assert.Len(t, changedDogus, len(tt.want))
			for _, doguName := range tt.want {
				assert.Contains(t, changedDogus, doguName)
			}
		})
	}
}
