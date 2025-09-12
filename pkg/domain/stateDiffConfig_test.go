package domain

import (
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	dogu1              = cescommons.SimpleName("dogu1")
	dogu2              = cescommons.SimpleName("dogu2")
	dogu1Key1          = common.DoguConfigKey{DoguName: dogu1, Key: "key1"}
	dogu1Key2          = common.DoguConfigKey{DoguName: dogu1, Key: "key2"}
	dogu1Key3          = common.DoguConfigKey{DoguName: dogu1, Key: "key3"}
	dogu1Key4          = common.DoguConfigKey{DoguName: dogu1, Key: "key4"}
	sensitiveDogu1Key1 = common.SensitiveDoguConfigKey{DoguName: dogu1, Key: "key1"}
	sensitiveDogu1Key2 = common.SensitiveDoguConfigKey{DoguName: dogu1, Key: "key2"}
	sensitiveDogu1Key3 = common.SensitiveDoguConfigKey{DoguName: dogu1, Key: "key3"}
	val1               = "value1"
	val2               = "value2"
	val3               = "value3"
)

func Test_determineConfigDiff(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		dogusConfigDiffs, sensitiveConfigDiffs, globalConfigDiff := determineConfigDiffs(
			nil,
			config.CreateGlobalConfig(map[config.Key]config.Value{}),
			map[cescommons.SimpleName]config.DoguConfig{},
			map[cescommons.SimpleName]config.DoguConfig{},
			map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{},
		)

		assert.Nil(t, dogusConfigDiffs)
		assert.Nil(t, sensitiveConfigDiffs)
		assert.Nil(t, globalConfigDiff)
	})
	t.Run("empty", func(t *testing.T) {
		emptyConfig := Config{}

		dogusConfigDiffs, sensitiveConfigDiffs, globalConfigDiff := determineConfigDiffs(
			&emptyConfig,
			config.CreateGlobalConfig(map[config.Key]config.Value{}),
			map[cescommons.SimpleName]config.DoguConfig{},
			map[cescommons.SimpleName]config.DoguConfig{},
			map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{},
		)

		assert.Equal(t, map[cescommons.SimpleName]DoguConfigDiffs{}, dogusConfigDiffs)
		assert.Equal(t, map[cescommons.SimpleName]SensitiveDoguConfigDiffs{}, sensitiveConfigDiffs)
		assert.Equal(t, GlobalConfigDiffs(nil), globalConfigDiff)
	})
	t.Run("all actions global config", func(t *testing.T) {
		//given ecosystem config

		entries, _ := config.MapToEntries(map[string]any{
			"key1": "value1", // for action none
			"key2": "value2", // for action set
			"key3": "value3", // for action delete
			// key4 is absent -> action none
		})
		globalConfig := config.CreateGlobalConfig(entries)
		//given blueprint config
		givenConfig := Config{
			Global: GlobalConfig{
				Present: map[common.GlobalConfigKey]common.GlobalConfigValue{
					"key1": "value1",
					"key2": "value3",
				},
				Absent: []common.GlobalConfigKey{
					"key3", "key4",
				},
			},
		}

		//when
		dogusConfigDiffs, sensitiveConfigDiffs, globalConfigDiff := determineConfigDiffs(
			&givenConfig,
			globalConfig,
			map[cescommons.SimpleName]config.DoguConfig{},
			map[cescommons.SimpleName]config.DoguConfig{},
			map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{},
		)

		//then
		assert.Equal(t, map[cescommons.SimpleName]DoguConfigDiffs{}, dogusConfigDiffs)
		assert.Equal(t, map[cescommons.SimpleName]SensitiveDoguConfigDiffs{}, sensitiveConfigDiffs)
		assert.Equal(t, 4, len(globalConfigDiff))
		hitKeys := make(map[string]bool)
		for _, diff := range globalConfigDiff {
			if diff.Key == "key1" {
				assert.Empty(t, cmp.Diff(diff, GlobalConfigEntryDiff{
					Key: "key1",
					Actual: GlobalConfigValueState{
						Value:  &val1,
						Exists: true,
					},
					Expected: GlobalConfigValueState{
						Value:  &val1,
						Exists: true,
					},
					NeededAction: ConfigActionNone,
				}))
				hitKeys["key1"] = true
			}
			if diff.Key == "key2" {
				assert.Empty(t, cmp.Diff(diff, GlobalConfigEntryDiff{
					Key: "key2",
					Actual: GlobalConfigValueState{
						Value:  &val2,
						Exists: true,
					},
					Expected: GlobalConfigValueState{
						Value:  &val3,
						Exists: true,
					},
					NeededAction: ConfigActionSet,
				}))
				hitKeys["key2"] = true
			}
			if diff.Key == "key3" {
				assert.Empty(t, cmp.Diff(diff, GlobalConfigEntryDiff{
					Key: "key3",
					Actual: GlobalConfigValueState{
						Value:  &val3,
						Exists: true,
					},
					Expected: GlobalConfigValueState{
						Value:  nil,
						Exists: false,
					},
					NeededAction: ConfigActionRemove,
				}))
				hitKeys["key3"] = true
			}
			if diff.Key == "key4" {
				assert.Empty(t, cmp.Diff(diff, GlobalConfigEntryDiff{
					Key: "key4",
					Actual: GlobalConfigValueState{
						Value:  nil,
						Exists: false,
					},
					Expected: GlobalConfigValueState{
						Value:  nil,
						Exists: false,
					},
					NeededAction: ConfigActionNone,
				}))
				hitKeys["key4"] = true
			}
		}
		assert.Equal(t, 4, len(hitKeys))
	})

	t.Run("all actions normal dogu config", func(t *testing.T) {
		//given ecosystem config
		globalConfigEntries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(globalConfigEntries)

		doguConfigEntries, _ := config.MapToEntries(map[string]any{
			"key1": "value1", //action none
			"key2": "value1", //action set
			"key3": "value1", //action delete
			//key4 -> absent, so action none
		})
		doguConfig := config.CreateDoguConfig(dogu1, doguConfigEntries)

		//given blueprint config
		givenConfig := Config{
			Dogus: map[cescommons.SimpleName]CombinedDoguConfig{
				"dogu1": {
					DoguName: "dogu1",
					Config: DoguConfig{
						Present: map[common.DoguConfigKey]common.DoguConfigValue{
							dogu1Key1: "value1",
							dogu1Key2: "value2",
						},
						Absent: []common.DoguConfigKey{
							dogu1Key3, dogu1Key4,
						},
					},
				},
			},
		}

		//when
		dogusConfigDiffs, sensitiveConfigDiffs, globalConfigDiff := determineConfigDiffs(
			&givenConfig,
			globalConfig,
			map[cescommons.SimpleName]config.DoguConfig{
				dogu1: doguConfig,
			},
			map[cescommons.SimpleName]config.DoguConfig{},
			map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{},
		)
		//then
		assert.Equal(t, GlobalConfigDiffs(nil), globalConfigDiff)
		require.NotNil(t, dogusConfigDiffs["dogu1"])
		assert.Equal(t, SensitiveDoguConfigDiffs(nil), sensitiveConfigDiffs["dogu1"])
		assert.Equal(t, 4, len(dogusConfigDiffs["dogu1"]))
		hitKeys := make(map[common.DoguConfigKey]bool)
		for _, diff := range dogusConfigDiffs["dogu1"] {
			if diff.Key == dogu1Key1 {
				assert.Empty(t, cmp.Diff(diff, DoguConfigEntryDiff{
					Key: dogu1Key1,
					Actual: DoguConfigValueState{
						Value:  &val1,
						Exists: true,
					},
					Expected: DoguConfigValueState{
						Value:  &val1,
						Exists: true,
					},
					NeededAction: ConfigActionNone,
				}))
				hitKeys[dogu1Key1] = true
			}
			if diff.Key == dogu1Key2 {
				assert.Empty(t, cmp.Diff(diff, DoguConfigEntryDiff{
					Key: dogu1Key2,
					Actual: DoguConfigValueState{
						Value:  &val1,
						Exists: true,
					},
					Expected: DoguConfigValueState{
						Value:  &val2,
						Exists: true,
					},
					NeededAction: ConfigActionSet,
				}))
				hitKeys[dogu1Key2] = true
			}
			if diff.Key == dogu1Key3 {
				assert.Empty(t, cmp.Diff(diff, DoguConfigEntryDiff{
					Key: dogu1Key3,
					Actual: DoguConfigValueState{
						Value:  &val1,
						Exists: true,
					},
					Expected: DoguConfigValueState{
						Value:  nil,
						Exists: false,
					},
					NeededAction: ConfigActionRemove,
				}))
				hitKeys[dogu1Key3] = true
			}
			//domain.DoguConfigEntryDiff{Key:common.DoguConfigKey{DoguName:"dogu1", Key:"key3"},
			//Actual:domain.DoguConfigValueState{Value:"value", Exists:true},
			//Expected:domain.DoguConfigValueState{Value:"", Exists:false},
			//NeededAction:"set"}
			if diff.Key == dogu1Key4 {
				assert.Empty(t, cmp.Diff(diff, DoguConfigEntryDiff{
					Key: dogu1Key4,
					Actual: DoguConfigValueState{
						Value:  nil,
						Exists: false,
					},
					Expected: DoguConfigValueState{
						Value:  nil,
						Exists: false,
					},
					NeededAction: ConfigActionNone,
				}))
				hitKeys[dogu1Key4] = true
			}
		}
		assert.Equal(t, 4, len(hitKeys))
	})
	t.Run("all actions for sensitive dogu config for present dogu", func(t *testing.T) {
		//given ecosystem config
		globalConfigEntries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(globalConfigEntries)

		sensitiveDoguConfigEntries, _ := config.MapToEntries(map[string]any{
			"key1": "value1", //action none
			"key2": "value1", //action set
			//key3 absent, action none
		})
		sensitiveDoguConfig := config.CreateDoguConfig(dogu1, sensitiveDoguConfigEntries)

		//given blueprint config
		givenConfig := Config{
			Dogus: map[cescommons.SimpleName]CombinedDoguConfig{
				"dogu1": {
					DoguName: "dogu1",
					SensitiveConfig: SensitiveDoguConfig{
						Present: map[common.SensitiveDoguConfigKey]SensitiveValueRef{
							sensitiveDogu1Key1: {
								SecretName: "mySecret1",
								SecretKey:  "myKey1",
							},
							sensitiveDogu1Key2: {
								SecretName: "mySecret2",
								SecretKey:  "myKey2",
							},
						},
						Absent: []common.SensitiveDoguConfigKey{
							sensitiveDogu1Key3,
						},
					},
				},
			},
		}

		//when
		dogusConfigDiffs, sensitiveConfigDiffs, globalConfigDiff := determineConfigDiffs(
			&givenConfig,
			globalConfig,
			map[cescommons.SimpleName]config.DoguConfig{},
			map[cescommons.SimpleName]config.DoguConfig{
				dogu1: sensitiveDoguConfig,
			},
			//loaded referenced sensitive config
			map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{
				sensitiveDogu1Key1: "value1",
				sensitiveDogu1Key2: "value2",
			},
		)
		//then
		assert.Equal(t, GlobalConfigDiffs(nil), globalConfigDiff)
		assert.Equal(t, DoguConfigDiffs(nil), dogusConfigDiffs["dogu1"])
		require.NotNil(t, sensitiveConfigDiffs["dogu1"])
		assert.Equal(t, 3, len(sensitiveConfigDiffs["dogu1"]))

		entriesDogu1 := SensitiveDoguConfigDiffs{
			{
				Key: sensitiveDogu1Key1,
				Actual: DoguConfigValueState{
					Value:  &val1,
					Exists: true,
				},
				Expected: DoguConfigValueState{
					Value:  &val1,
					Exists: true,
				},
				NeededAction: ConfigActionNone,
			},
			{
				Key: sensitiveDogu1Key2,
				Actual: DoguConfigValueState{
					Value:  &val1,
					Exists: true,
				},
				Expected: DoguConfigValueState{
					Value:  &val2,
					Exists: true,
				},
				NeededAction: ConfigActionSet,
			},
			{
				Key: sensitiveDogu1Key3,
				Actual: DoguConfigValueState{
					Value:  nil,
					Exists: false,
				},
				Expected: DoguConfigValueState{
					Value:  nil,
					Exists: false,
				},
				NeededAction: ConfigActionNone,
			},
		}
		assert.ElementsMatch(t, sensitiveConfigDiffs["dogu1"], entriesDogu1)
	})
	t.Run("all actions for sensitive dogu config for absent dogu", func(t *testing.T) {
		//given ecosystem config
		globalConfigEntries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(globalConfigEntries)

		doguConfigEntries, _ := config.MapToEntries(map[string]any{})
		doguConfig := config.CreateDoguConfig(dogu1, doguConfigEntries)

		sensitiveDoguConfigEntries, _ := config.MapToEntries(map[string]any{})
		sensitiveDoguConfig := config.CreateDoguConfig(dogu1, sensitiveDoguConfigEntries)

		//given blueprint config
		givenConfig := Config{
			Dogus: map[cescommons.SimpleName]CombinedDoguConfig{
				"dogu1": {
					DoguName: "dogu1",
					SensitiveConfig: SensitiveDoguConfig{
						Present: map[common.SensitiveDoguConfigKey]SensitiveValueRef{
							sensitiveDogu1Key1: {
								SecretName: "secret1",
								SecretKey:  "key1",
							},
						},
						Absent: []common.SensitiveDoguConfigKey{},
					},
				},
			},
		}

		//when
		dogusConfigDiffs, sensitiveConfigDiffs, _ := determineConfigDiffs(
			&givenConfig,
			globalConfig,
			map[cescommons.SimpleName]config.DoguConfig{
				dogu1: doguConfig,
			},
			map[cescommons.SimpleName]config.DoguConfig{
				dogu1: sensitiveDoguConfig,
			},
			//loaded referenced sensitive config
			map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{
				sensitiveDogu1Key1: "value1",
			},
		)
		//then
		assert.Equal(t, DoguConfigDiffs(nil), dogusConfigDiffs["dogu1"])

		require.NotNil(t, sensitiveConfigDiffs["dogu1"])
		require.Equal(t, 1, len(sensitiveConfigDiffs["dogu1"]))
		assert.Empty(t, cmp.Diff(sensitiveConfigDiffs["dogu1"][0], SensitiveDoguConfigEntryDiff{
			Key: sensitiveDogu1Key1,
			Actual: DoguConfigValueState{
				Value:  nil,
				Exists: false,
			},
			Expected: DoguConfigValueState{
				Value:  &val1,
				Exists: true,
			},
			NeededAction: ConfigActionSet,
		}))

	})
}

func Test_getNeededConfigAction(t *testing.T) {
	tests := []struct {
		name     string
		expected ConfigValueState
		actual   ConfigValueState
		want     ConfigAction
	}{
		{
			name:     "action none, both do not exist",
			expected: ConfigValueState{Value: nil, Exists: false},
			actual:   ConfigValueState{Value: nil, Exists: false},
			want:     ConfigActionNone,
		},
		{
			name:     "action none, for some reason the values are different",
			expected: ConfigValueState{Value: &val1, Exists: false},
			actual:   ConfigValueState{Value: &val2, Exists: false},
			want:     ConfigActionNone,
		},
		{
			name:     "action none, equal values",
			expected: ConfigValueState{Value: &val1, Exists: true},
			actual:   ConfigValueState{Value: &val1, Exists: true},
			want:     ConfigActionNone,
		},
		{
			name:     "set new value",
			expected: ConfigValueState{Value: nil, Exists: true},
			actual:   ConfigValueState{Value: nil, Exists: false},
			want:     ConfigActionSet,
		},
		{
			name:     "update value",
			expected: ConfigValueState{Value: &val1, Exists: true},
			actual:   ConfigValueState{Value: &val2, Exists: true},
			want:     ConfigActionSet,
		},
		{
			name:     "remove value",
			expected: ConfigValueState{Value: nil, Exists: false},
			actual:   ConfigValueState{Value: nil, Exists: true},
			want:     ConfigActionRemove,
		},
		{
			name:     "remove value",
			expected: ConfigValueState{Value: nil, Exists: false},
			actual:   ConfigValueState{Value: &val3, Exists: true},
			want:     ConfigActionRemove,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getNeededConfigAction(tt.expected, tt.actual), "getNeededConfigAction(%v, %v)", tt.expected, tt.actual)
		})
	}
}
