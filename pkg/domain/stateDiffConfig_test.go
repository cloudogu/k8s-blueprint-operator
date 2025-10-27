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
	dogu1     = cescommons.SimpleName("dogu1")
	dogu2     = cescommons.SimpleName("dogu2")
	dogu1Key1 = common.DoguConfigKey{DoguName: dogu1, Key: "key1"}
	dogu1Key2 = common.DoguConfigKey{DoguName: dogu1, Key: "key2"}
	dogu1Key3 = common.DoguConfigKey{DoguName: dogu1, Key: "key3"}
	dogu1Key4 = common.DoguConfigKey{DoguName: dogu1, Key: "key4"}
	val1      = config.Value("value1")
	val2      = config.Value("value2")
	val3      = config.Value("value3")
)

func Test_determineConfigDiff(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		emptyConfig := Config{}

		dogusConfigDiffs, sensitiveConfigDiffs, globalConfigDiff := determineConfigDiffs(
			emptyConfig,
			config.CreateGlobalConfig(map[config.Key]config.Value{}),
			map[cescommons.SimpleName]config.DoguConfig{},
			map[cescommons.SimpleName]config.DoguConfig{},
			map[common.DoguConfigKey]common.SensitiveDoguConfigValue{},
			map[common.DoguConfigKey]common.DoguConfigValue{},
			map[common.GlobalConfigKey]common.GlobalConfigValue{},
			map[common.GlobalConfigKey]common.GlobalConfigValue{},
		)

		assert.Nil(t, dogusConfigDiffs)
		assert.Nil(t, sensitiveConfigDiffs)
		assert.Equal(t, GlobalConfigDiffs(nil), globalConfigDiff)
	})
	t.Run("all actions global config", func(t *testing.T) {
		//given ecosystem config

		entries, _ := config.MapToEntries(map[string]any{
			"key1": val1.String(), // for action none
			"key2": val2.String(), // for action set
			"key3": val3.String(), // for action delete
			// key4 is absent -> action none
		})
		globalConfig := config.CreateGlobalConfig(entries)
		//given blueprint config
		givenConfig := Config{
			Global: GlobalConfigEntries{
				{
					Key:   "key1",
					Value: &val1,
				},
				{
					Key:   "key2",
					Value: &val3,
				},
				{
					Key:    "key3",
					Absent: true,
				},
				{
					Key:    "key4",
					Absent: true,
				},
			},
		}

		//when
		dogusConfigDiffs, sensitiveConfigDiffs, globalConfigDiff := determineConfigDiffs(
			givenConfig,
			globalConfig,
			map[cescommons.SimpleName]config.DoguConfig{},
			map[cescommons.SimpleName]config.DoguConfig{},
			map[common.DoguConfigKey]common.SensitiveDoguConfigValue{},
			map[common.DoguConfigKey]common.DoguConfigValue{},
			map[common.GlobalConfigKey]common.GlobalConfigValue{},
			map[common.GlobalConfigKey]common.GlobalConfigValue{},
		)

		//then
		assert.Nil(t, dogusConfigDiffs)
		assert.Nil(t, sensitiveConfigDiffs)
		assert.Equal(t, 2, len(globalConfigDiff)) // only changes
		hitKeys := make(map[string]bool)
		for _, diff := range globalConfigDiff {
			if diff.Key == "key2" {
				assert.Empty(t, cmp.Diff(diff, GlobalConfigEntryDiff{
					Key: "key2",
					Actual: GlobalConfigValueState{
						Value:  (*string)(&val2),
						Exists: true,
					},
					Expected: GlobalConfigValueState{
						Value:  (*string)(&val3),
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
						Value:  (*string)(&val3),
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
		}
		assert.Equal(t, 2, len(hitKeys))
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
			Dogus: map[cescommons.SimpleName]DoguConfigEntries{
				"dogu1": {
					{
						Key:   dogu1Key1.Key,
						Value: &val1,
					},
					{
						Key:   dogu1Key2.Key,
						Value: &val2,
					},
					{
						Key:    dogu1Key3.Key,
						Absent: true,
					},
					{
						Key:    dogu1Key4.Key,
						Absent: true,
					},
				},
			},
		}

		//when
		dogusConfigDiffs, sensitiveConfigDiffs, globalConfigDiff := determineConfigDiffs(
			givenConfig,
			globalConfig,
			map[cescommons.SimpleName]config.DoguConfig{
				dogu1: doguConfig,
			},
			map[cescommons.SimpleName]config.DoguConfig{},
			map[common.DoguConfigKey]common.SensitiveDoguConfigValue{},
			map[common.DoguConfigKey]common.DoguConfigValue{
				dogu1Key2: "value1",
			},
			map[common.GlobalConfigKey]common.GlobalConfigValue{},
			map[common.GlobalConfigKey]common.GlobalConfigValue{},
		)
		//then
		assert.Equal(t, GlobalConfigDiffs(nil), globalConfigDiff)
		require.NotNil(t, dogusConfigDiffs["dogu1"])
		assert.Equal(t, SensitiveDoguConfigDiffs(nil), sensitiveConfigDiffs["dogu1"])
		assert.Equal(t, 2, len(dogusConfigDiffs["dogu1"])) // only changes
		hitKeys := make(map[common.DoguConfigKey]bool)
		for _, diff := range dogusConfigDiffs["dogu1"] {
			if diff.Key == dogu1Key2 {
				assert.Empty(t, cmp.Diff(diff, DoguConfigEntryDiff{
					Key: dogu1Key2,
					Actual: DoguConfigValueState{
						Value:  (*string)(&val1),
						Exists: true,
					},
					Expected: DoguConfigValueState{
						Value:  (*string)(&val2),
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
						Value:  (*string)(&val1),
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
		}
		assert.Equal(t, 1, len(hitKeys))
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
			Dogus: map[cescommons.SimpleName]DoguConfigEntries{
				"dogu1": {
					{
						Key:       dogu1Key1.Key,
						Sensitive: true,
						SecretRef: &SensitiveValueRef{
							SecretName: "mySecret1",
							SecretKey:  "myKey1",
						},
					},
					{
						Key:       dogu1Key2.Key,
						Sensitive: true,
						SecretRef: &SensitiveValueRef{
							SecretName: "mySecret2",
							SecretKey:  "myKey2",
						},
					},
					{
						Key:       dogu1Key3.Key,
						Sensitive: true,
						Absent:    true,
					},
				},
			},
		}

		//when
		dogusConfigDiffs, sensitiveConfigDiffs, globalConfigDiff := determineConfigDiffs(
			givenConfig,
			globalConfig,
			map[cescommons.SimpleName]config.DoguConfig{},
			map[cescommons.SimpleName]config.DoguConfig{
				dogu1: sensitiveDoguConfig,
			},
			//loaded referenced sensitive config
			map[common.DoguConfigKey]common.SensitiveDoguConfigValue{
				dogu1Key1: "value1",
				dogu1Key2: "value2",
			},
			map[common.DoguConfigKey]common.DoguConfigValue{},
			map[common.GlobalConfigKey]common.GlobalConfigValue{},
			map[common.GlobalConfigKey]common.GlobalConfigValue{},
		)
		//then
		assert.Equal(t, GlobalConfigDiffs(nil), globalConfigDiff)
		assert.Equal(t, DoguConfigDiffs(nil), dogusConfigDiffs["dogu1"])
		require.NotNil(t, sensitiveConfigDiffs["dogu1"])
		assert.Equal(t, 1, len(sensitiveConfigDiffs["dogu1"])) // only changes

		entriesDogu1 := SensitiveDoguConfigDiffs{
			{
				Key: dogu1Key2,
				Actual: DoguConfigValueState{
					Value:  (*string)(&val1),
					Exists: true,
				},
				Expected: DoguConfigValueState{
					Value:  (*string)(&val2),
					Exists: true,
				},
				NeededAction: ConfigActionSet,
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
			Dogus: map[cescommons.SimpleName]DoguConfigEntries{
				"dogu1": {
					{
						Key:       dogu1Key1.Key,
						Sensitive: true,
						SecretRef: &SensitiveValueRef{
							SecretName: "mySecret1",
							SecretKey:  "myKey1",
						},
					},
				},
			},
		}

		//when
		dogusConfigDiffs, sensitiveConfigDiffs, _ := determineConfigDiffs(
			givenConfig,
			globalConfig,
			map[cescommons.SimpleName]config.DoguConfig{
				dogu1: doguConfig,
			},
			map[cescommons.SimpleName]config.DoguConfig{
				dogu1: sensitiveDoguConfig,
			},
			//loaded referenced sensitive config
			map[common.DoguConfigKey]common.SensitiveDoguConfigValue{
				dogu1Key1: "value1",
			},
			map[common.DoguConfigKey]common.DoguConfigValue{},
			map[common.GlobalConfigKey]common.GlobalConfigValue{},
			map[common.GlobalConfigKey]common.GlobalConfigValue{},
		)
		//then
		assert.Equal(t, DoguConfigDiffs(nil), dogusConfigDiffs["dogu1"])

		require.NotNil(t, sensitiveConfigDiffs["dogu1"])
		require.Equal(t, 1, len(sensitiveConfigDiffs["dogu1"]))
		assert.Empty(t, cmp.Diff(sensitiveConfigDiffs["dogu1"][0], SensitiveDoguConfigEntryDiff{
			Key: dogu1Key1,
			Actual: DoguConfigValueState{
				Value:  nil,
				Exists: false,
			},
			Expected: DoguConfigValueState{
				Value:  (*string)(&val1),
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
			expected: ConfigValueState{Value: (*string)(&val1), Exists: false},
			actual:   ConfigValueState{Value: (*string)(&val2), Exists: false},
			want:     ConfigActionNone,
		},
		{
			name:     "action none, equal values",
			expected: ConfigValueState{Value: (*string)(&val1), Exists: true},
			actual:   ConfigValueState{Value: (*string)(&val1), Exists: true},
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
			expected: ConfigValueState{Value: (*string)(&val1), Exists: true},
			actual:   ConfigValueState{Value: (*string)(&val2), Exists: true},
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
			actual:   ConfigValueState{Value: (*string)(&val3), Exists: true},
			want:     ConfigActionRemove,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getNeededConfigAction(tt.expected, tt.actual), "getNeededConfigAction(%v, %v)", tt.expected, tt.actual)
		})
	}
}
