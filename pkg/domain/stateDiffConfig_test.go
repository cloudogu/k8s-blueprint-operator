package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bpv2 "github.com/cloudogu/blueprint-lib/v2"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-registry-lib/config"
)

var (
	dogu1              = cescommons.SimpleName("dogu1")
	dogu2              = cescommons.SimpleName("dogu2")
	dogu1Key1          = bpv2.DoguConfigKey{DoguName: dogu1, Key: "key1"}
	dogu1Key2          = bpv2.DoguConfigKey{DoguName: dogu1, Key: "key2"}
	dogu1Key3          = bpv2.DoguConfigKey{DoguName: dogu1, Key: "key3"}
	dogu1Key4          = bpv2.DoguConfigKey{DoguName: dogu1, Key: "key4"}
	sensitiveDogu1Key1 = bpv2.SensitiveDoguConfigKey{DoguName: dogu1, Key: "key1"}
	sensitiveDogu1Key2 = bpv2.SensitiveDoguConfigKey{DoguName: dogu1, Key: "key2"}
	sensitiveDogu1Key3 = bpv2.SensitiveDoguConfigKey{DoguName: dogu1, Key: "key3"}
)

func Test_determineConfigDiff(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		emptyConfig := bpv2.Config{}

		dogusConfigDiffs, sensitiveConfigDiffs, globalConfigDiff := determineConfigDiffs(
			emptyConfig,
			config.CreateGlobalConfig(map[config.Key]config.Value{}),
			map[cescommons.SimpleName]config.DoguConfig{},
			map[cescommons.SimpleName]config.DoguConfig{},
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
		givenConfig := bpv2.Config{
			Global: bpv2.GlobalConfig{
				Present: map[bpv2.GlobalConfigKey]bpv2.GlobalConfigValue{
					"key1": "value1",
					"key2": "value2.2",
				},
				Absent: []bpv2.GlobalConfigKey{
					"key3", "key4",
				},
			},
		}

		//when
		dogusConfigDiffs, sensitiveConfigDiffs, globalConfigDiff := determineConfigDiffs(
			givenConfig,
			globalConfig,
			map[cescommons.SimpleName]config.DoguConfig{},
			map[cescommons.SimpleName]config.DoguConfig{},
		)

		//then
		assert.Equal(t, map[cescommons.SimpleName]DoguConfigDiffs{}, dogusConfigDiffs)
		assert.Equal(t, map[cescommons.SimpleName]SensitiveDoguConfigDiffs{}, sensitiveConfigDiffs)
		assert.Equal(t, 4, len(globalConfigDiff))
		assert.Contains(t, globalConfigDiff, GlobalConfigEntryDiff{
			Key: "key1",
			Actual: GlobalConfigValueState{
				Value:  "value1",
				Exists: true,
			},
			Expected: GlobalConfigValueState{
				Value:  "value1",
				Exists: true,
			},
			NeededAction: ConfigActionNone,
		})
		assert.Contains(t, globalConfigDiff, GlobalConfigEntryDiff{
			Key: "key2",
			Actual: GlobalConfigValueState{
				Value:  "value2",
				Exists: true,
			},
			Expected: GlobalConfigValueState{
				Value:  "value2.2",
				Exists: true,
			},
			NeededAction: ConfigActionSet,
		})
		assert.Contains(t, globalConfigDiff, GlobalConfigEntryDiff{
			Key: "key3",
			Actual: GlobalConfigValueState{
				Value:  "value3",
				Exists: true,
			},
			Expected: GlobalConfigValueState{
				Value:  "",
				Exists: false,
			},
			NeededAction: ConfigActionRemove,
		})
		assert.Contains(t, globalConfigDiff, GlobalConfigEntryDiff{
			Key: "key4",
			Actual: GlobalConfigValueState{
				Value:  "",
				Exists: false,
			},
			Expected: GlobalConfigValueState{
				Value:  "",
				Exists: false,
			},
			NeededAction: ConfigActionNone,
		})
	})
	t.Run("all actions normal dogu config", func(t *testing.T) {
		//given ecosystem config
		globalConfigEntries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(globalConfigEntries)

		doguConfigEntries, _ := config.MapToEntries(map[string]any{
			"key1": "value", //action none
			"key2": "value", //action set
			"key3": "value", //action delete
			//key4 -> absent, so action none
		})
		doguConfig := config.CreateDoguConfig(dogu1, doguConfigEntries)

		//given blueprint config
		givenConfig := bpv2.Config{
			Dogus: map[cescommons.SimpleName]bpv2.CombinedDoguConfig{
				"dogu1": {
					DoguName: "dogu1",
					Config: bpv2.DoguConfig{
						Present: map[bpv2.DoguConfigKey]bpv2.DoguConfigValue{
							dogu1Key1: "value",
							dogu1Key2: "updatedValue",
						},
						Absent: []bpv2.DoguConfigKey{
							dogu1Key3, dogu1Key4,
						},
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
		)
		//then
		assert.Equal(t, GlobalConfigDiffs(nil), globalConfigDiff)
		require.NotNil(t, dogusConfigDiffs["dogu1"])
		assert.Equal(t, SensitiveDoguConfigDiffs(nil), sensitiveConfigDiffs["dogu1"])
		assert.Equal(t, 4, len(dogusConfigDiffs["dogu1"]))
		assert.Contains(t, dogusConfigDiffs["dogu1"], DoguConfigEntryDiff{
			Key: dogu1Key1,
			Actual: DoguConfigValueState{
				Value:  "value",
				Exists: true,
			},
			Expected: DoguConfigValueState{
				Value:  "value",
				Exists: true,
			},
			NeededAction: ConfigActionNone,
		})
		assert.Contains(t, dogusConfigDiffs["dogu1"], DoguConfigEntryDiff{
			Key: dogu1Key2,
			Actual: DoguConfigValueState{
				Value:  "value",
				Exists: true,
			},
			Expected: DoguConfigValueState{
				Value:  "updatedValue",
				Exists: true,
			},
			NeededAction: ConfigActionSet,
		})
		assert.Contains(t, dogusConfigDiffs["dogu1"], DoguConfigEntryDiff{
			Key: dogu1Key3,
			Actual: DoguConfigValueState{
				Value:  "value",
				Exists: true,
			},
			Expected: DoguConfigValueState{
				Value:  "",
				Exists: false,
			},
			NeededAction: ConfigActionRemove,
		})
		//domain.DoguConfigEntryDiff{Key:common.DoguConfigKey{DoguName:"dogu1", Key:"key3"},
		//Actual:domain.DoguConfigValueState{Value:"value", Exists:true},
		//Expected:domain.DoguConfigValueState{Value:"", Exists:false},
		//NeededAction:"set"}
		assert.Contains(t, dogusConfigDiffs["dogu1"], DoguConfigEntryDiff{
			Key: dogu1Key4,
			Actual: DoguConfigValueState{
				Value:  "",
				Exists: false,
			},
			Expected: DoguConfigValueState{
				Value:  "",
				Exists: false,
			},
			NeededAction: ConfigActionNone,
		})
	})
	t.Run("all actions for sensitive dogu config for present dogu", func(t *testing.T) {
		//given ecosystem config
		globalConfigEntries, _ := config.MapToEntries(map[string]any{})
		globalConfig := config.CreateGlobalConfig(globalConfigEntries)

		sensitiveDoguConfigEntries, _ := config.MapToEntries(map[string]any{
			"key1": "value", //action none
			"key2": "value", //action set
			//key3 absent, action none
		})
		sensitiveDoguConfig := config.CreateDoguConfig(dogu1, sensitiveDoguConfigEntries)

		//given blueprint config
		givenConfig := bpv2.Config{
			Dogus: map[cescommons.SimpleName]bpv2.CombinedDoguConfig{
				"dogu1": {
					DoguName: "dogu1",
					SensitiveConfig: bpv2.SensitiveDoguConfig{
						Present: map[bpv2.SensitiveDoguConfigKey]bpv2.SensitiveDoguConfigValue{
							sensitiveDogu1Key1: "value",
							sensitiveDogu1Key2: "updated value",
						},
						Absent: []bpv2.SensitiveDoguConfigKey{
							sensitiveDogu1Key3,
						},
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
		)
		//then
		assert.Equal(t, GlobalConfigDiffs(nil), globalConfigDiff)
		assert.Equal(t, DoguConfigDiffs(nil), dogusConfigDiffs["dogu1"])
		require.NotNil(t, sensitiveConfigDiffs["dogu1"])
		assert.Equal(t, 3, len(sensitiveConfigDiffs["dogu1"]))

		entriesDogu1 := []SensitiveDoguConfigEntryDiff{
			{
				Key: sensitiveDogu1Key1,
				Actual: DoguConfigValueState{
					Value:  "value",
					Exists: true,
				},
				Expected: DoguConfigValueState{
					Value:  "value",
					Exists: true,
				},
				NeededAction: ConfigActionNone,
			},
			{
				Key: sensitiveDogu1Key2,
				Actual: DoguConfigValueState{
					Value:  "value",
					Exists: true,
				},
				Expected: DoguConfigValueState{
					Value:  "updated value",
					Exists: true,
				},
				NeededAction: ConfigActionSet,
			},
			{
				Key: sensitiveDogu1Key3,
				Actual: DoguConfigValueState{
					Value:  "",
					Exists: false,
				},
				Expected: DoguConfigValueState{
					Value:  "",
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
		givenConfig := bpv2.Config{
			Dogus: map[cescommons.SimpleName]bpv2.CombinedDoguConfig{
				"dogu1": {
					DoguName: "dogu1",
					SensitiveConfig: bpv2.SensitiveDoguConfig{
						Present: map[bpv2.SensitiveDoguConfigKey]bpv2.SensitiveDoguConfigValue{
							sensitiveDogu1Key1: "value",
						},
						Absent: []bpv2.SensitiveDoguConfigKey{},
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
		)
		//then
		assert.Equal(t, DoguConfigDiffs(nil), dogusConfigDiffs["dogu1"])

		require.NotNil(t, sensitiveConfigDiffs["dogu1"])
		require.Equal(t, 1, len(sensitiveConfigDiffs["dogu1"]))
		assert.Equal(t, sensitiveConfigDiffs["dogu1"][0], SensitiveDoguConfigEntryDiff{
			Key: sensitiveDogu1Key1,
			Actual: DoguConfigValueState{
				Value:  "",
				Exists: false,
			},
			Expected: DoguConfigValueState{
				Value:  "value",
				Exists: true,
			},
			NeededAction: ConfigActionSet,
		})

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
			expected: ConfigValueState{Value: "", Exists: false},
			actual:   ConfigValueState{Value: "", Exists: false},
			want:     ConfigActionNone,
		},
		{
			name:     "action none, for some reason the values are different",
			expected: ConfigValueState{Value: "1", Exists: false},
			actual:   ConfigValueState{Value: "2", Exists: false},
			want:     ConfigActionNone,
		},
		{
			name:     "action none, equal values",
			expected: ConfigValueState{Value: "1", Exists: true},
			actual:   ConfigValueState{Value: "1", Exists: true},
			want:     ConfigActionNone,
		},
		{
			name:     "set new value",
			expected: ConfigValueState{Value: "", Exists: true},
			actual:   ConfigValueState{Value: "", Exists: false},
			want:     ConfigActionSet,
		},
		{
			name:     "update value",
			expected: ConfigValueState{Value: "1", Exists: true},
			actual:   ConfigValueState{Value: "2", Exists: true},
			want:     ConfigActionSet,
		},
		{
			name:     "remove value",
			expected: ConfigValueState{Value: "", Exists: false},
			actual:   ConfigValueState{Value: "", Exists: true},
			want:     ConfigActionRemove,
		},
		{
			name:     "remove value",
			expected: ConfigValueState{Value: "", Exists: false},
			actual:   ConfigValueState{Value: "value3", Exists: true},
			want:     ConfigActionRemove,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getNeededConfigAction(tt.expected, tt.actual), "getNeededConfigAction(%v, %v)", tt.expected, tt.actual)
		})
	}
}

func Test_censorValues(t *testing.T) {
	tests := []struct {
		name         string
		configByDogu map[cescommons.SimpleName]SensitiveDoguConfigDiffs
		want         map[cescommons.SimpleName]SensitiveDoguConfigDiffs
	}{
		{
			name:         "no diff at all",
			configByDogu: map[cescommons.SimpleName]SensitiveDoguConfigDiffs{},
			want:         map[cescommons.SimpleName]SensitiveDoguConfigDiffs{},
		},
		{
			name: "no diff for dogu",
			configByDogu: map[cescommons.SimpleName]SensitiveDoguConfigDiffs{
				dogu1: nil,
			},
			want: map[cescommons.SimpleName]SensitiveDoguConfigDiffs{
				dogu1: nil,
			},
		},
		{
			name: "censored actual and expected values",
			configByDogu: map[cescommons.SimpleName]SensitiveDoguConfigDiffs{
				dogu1: {DoguConfigEntryDiff{
					Key: dogu1Key1,
					Actual: DoguConfigValueState{
						Value:  "123",
						Exists: true,
					},
					Expected: DoguConfigValueState{
						Value:  "1234",
						Exists: true,
					},
					NeededAction: ConfigActionSet,
				}},
			},
			want: map[cescommons.SimpleName]SensitiveDoguConfigDiffs{
				dogu1: {DoguConfigEntryDiff{
					Key: dogu1Key1,
					Actual: DoguConfigValueState{
						Value:  bpv2.CensorValue,
						Exists: true,
					},
					Expected: DoguConfigValueState{
						Value:  bpv2.CensorValue,
						Exists: true,
					},
					NeededAction: ConfigActionSet,
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, censorValues(tt.configByDogu), "censorValues(%v)", tt.configByDogu)
		})
	}
}
