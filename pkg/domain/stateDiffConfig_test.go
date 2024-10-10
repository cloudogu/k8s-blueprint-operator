package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	dogu1              = common.SimpleDoguName("dogu1")
	dogu2              = common.SimpleDoguName("dogu2")
	dogu1Key1          = common.DoguConfigKey{DoguName: dogu1, Key: "key1"}
	dogu1Key2          = common.DoguConfigKey{DoguName: dogu1, Key: "key2"}
	dogu1Key3          = common.DoguConfigKey{DoguName: dogu1, Key: "key3"}
	dogu1Key4          = common.DoguConfigKey{DoguName: dogu1, Key: "key4"}
	dogu2Key1          = common.DoguConfigKey{DoguName: dogu2, Key: "key1"}
	sensitiveDogu1Key1 = common.SensitiveDoguConfigKey{DoguConfigKey: dogu1Key1}
	sensitiveDogu1Key2 = common.SensitiveDoguConfigKey{DoguConfigKey: dogu1Key2}
	sensitiveDogu1Key3 = common.SensitiveDoguConfigKey{DoguConfigKey: dogu1Key3}
	sensitiveDogu1Key4 = common.SensitiveDoguConfigKey{DoguConfigKey: dogu1Key4}
	sensitiveDogu2Key1 = common.SensitiveDoguConfigKey{DoguConfigKey: dogu2Key1}
)

func Test_determineConfigDiff(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		emptyConfig := Config{}

		dogusConfigDiffs, globalConfigDiff := determineConfigDiffs(
			emptyConfig,
			config.CreateGlobalConfig(map[config.Key]config.Value{}),
			map[common.SimpleDoguName]config.DoguConfig{},
			map[common.SimpleDoguName]config.DoguConfig{},
		)

		assert.Equal(t, map[common.SimpleDoguName]CombinedDoguConfigDiffs{}, dogusConfigDiffs)
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
					"key2": "value2.2",
				},
				Absent: []common.GlobalConfigKey{
					"key3", "key4",
				},
			},
		}

		//when
		dogusConfigDiffs, globalConfigDiff := determineConfigDiffs(
			givenConfig,
			globalConfig,
			map[common.SimpleDoguName]config.DoguConfig{},
			map[common.SimpleDoguName]config.DoguConfig{},
		)

		//then
		assert.Equal(t, map[common.SimpleDoguName]CombinedDoguConfigDiffs{}, dogusConfigDiffs)
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
		givenConfig := Config{
			Dogus: map[common.SimpleDoguName]CombinedDoguConfig{
				"dogu1": {
					DoguName: "dogu1",
					Config: DoguConfig{
						Present: map[common.DoguConfigKey]common.DoguConfigValue{
							dogu1Key1: "value",
							dogu1Key2: "updatedValue",
						},
						Absent: []common.DoguConfigKey{
							dogu1Key3, dogu1Key4,
						},
					},
				},
			},
		}

		//when
		dogusConfigDiffs, globalConfigDiff := determineConfigDiffs(
			givenConfig,
			globalConfig,
			map[common.SimpleDoguName]config.DoguConfig{
				dogu1: doguConfig,
			},
			map[common.SimpleDoguName]config.DoguConfig{},
		)
		//then
		assert.Equal(t, GlobalConfigDiffs(nil), globalConfigDiff)
		require.NotNil(t, dogusConfigDiffs["dogu1"])
		assert.Equal(t, SensitiveDoguConfigDiffs(nil), dogusConfigDiffs["dogu1"].SensitiveDoguConfigDiff)
		assert.Equal(t, 4, len(dogusConfigDiffs["dogu1"].DoguConfigDiff))
		assert.Contains(t, dogusConfigDiffs["dogu1"].DoguConfigDiff, DoguConfigEntryDiff{
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
		assert.Contains(t, dogusConfigDiffs["dogu1"].DoguConfigDiff, DoguConfigEntryDiff{
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
		assert.Contains(t, dogusConfigDiffs["dogu1"].DoguConfigDiff, DoguConfigEntryDiff{
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
		assert.Contains(t, dogusConfigDiffs["dogu1"].DoguConfigDiff, DoguConfigEntryDiff{
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
		givenConfig := Config{
			Dogus: map[common.SimpleDoguName]CombinedDoguConfig{
				"dogu1": {
					DoguName: "dogu1",
					SensitiveConfig: SensitiveDoguConfig{
						Present: map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{
							sensitiveDogu1Key1: "value",
							sensitiveDogu1Key2: "updated value",
						},
						Absent: []common.SensitiveDoguConfigKey{
							sensitiveDogu1Key3,
						},
					},
				},
			},
		}

		//when
		dogusConfigDiffs, globalConfigDiff := determineConfigDiffs(
			givenConfig,
			globalConfig,
			map[common.SimpleDoguName]config.DoguConfig{},
			map[common.SimpleDoguName]config.DoguConfig{
				dogu1: sensitiveDoguConfig,
			},
		)
		//then
		assert.Equal(t, GlobalConfigDiffs(nil), globalConfigDiff)
		require.NotNil(t, dogusConfigDiffs["dogu1"])
		assert.Equal(t, DoguConfigDiffs(nil), dogusConfigDiffs["dogu1"].DoguConfigDiff)
		assert.Equal(t, 3, len(dogusConfigDiffs["dogu1"].SensitiveDoguConfigDiff))

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
		assert.ElementsMatch(t, dogusConfigDiffs["dogu1"].SensitiveDoguConfigDiff, entriesDogu1)
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
			Dogus: map[common.SimpleDoguName]CombinedDoguConfig{
				"dogu1": {
					DoguName: "dogu1",
					SensitiveConfig: SensitiveDoguConfig{
						Present: map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{
							sensitiveDogu1Key1: "value",
						},
						Absent: []common.SensitiveDoguConfigKey{},
					},
				},
			},
		}

		//when
		dogusConfigDiffs, _ := determineConfigDiffs(
			givenConfig,
			globalConfig,
			map[common.SimpleDoguName]config.DoguConfig{
				dogu1: doguConfig,
			},
			map[common.SimpleDoguName]config.DoguConfig{
				dogu1: sensitiveDoguConfig,
			},
		)
		//then
		require.NotNil(t, dogusConfigDiffs["dogu1"])
		require.Equal(t, 1, len(dogusConfigDiffs["dogu1"].SensitiveDoguConfigDiff))
		assert.Equal(t, dogusConfigDiffs["dogu1"].SensitiveDoguConfigDiff[0], SensitiveDoguConfigEntryDiff{
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

func TestCombinedDoguConfigDiff_CensorValues(t *testing.T) {
	t.Run("Not censoring normal dogu config", func(t *testing.T) {
		//given
		configDiff := CombinedDoguConfigDiffs{
			DoguConfigDiff: []DoguConfigEntryDiff{
				{
					Key: common.DoguConfigKey{
						DoguName: "ldap",
						Key:      "logging/root",
					},
					Actual: DoguConfigValueState{
						Value:  "ERROR",
						Exists: false,
					},
					Expected: DoguConfigValueState{
						Value:  "DEBUG",
						Exists: false,
					},
					NeededAction: "Update",
				},
			},
		}

		//when
		result := configDiff.censorValues()

		require.Len(t, result.DoguConfigDiff, 1)

		assert.Equal(t, "ldap", string(result.DoguConfigDiff[0].Key.DoguName))
		assert.Equal(t, "logging/root", string(result.DoguConfigDiff[0].Key.Key))
		assert.Equal(t, "ERROR", result.DoguConfigDiff[0].Actual.Value)
		assert.Equal(t, false, result.DoguConfigDiff[0].Actual.Exists)
		assert.Equal(t, "DEBUG", result.DoguConfigDiff[0].Expected.Value)
		assert.Equal(t, false, result.DoguConfigDiff[0].Expected.Exists)
		assert.Equal(t, "Update", string(result.DoguConfigDiff[0].NeededAction))
	})

	t.Run("Censoring sensitive dogu config", func(t *testing.T) {
		//given
		configDiff := CombinedDoguConfigDiffs{
			SensitiveDoguConfigDiff: []SensitiveDoguConfigEntryDiff{
				{
					Key: common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
						DoguName: "ldap",
						Key:      "logging/root",
					}},
					Actual: DoguConfigValueState{
						Value:  "ERROR",
						Exists: false,
					},
					Expected: DoguConfigValueState{
						Value:  "DEBUG",
						Exists: false,
					},
					NeededAction: "Update",
				},
			},
		}

		//when
		result := configDiff.censorValues()

		require.Len(t, result.SensitiveDoguConfigDiff, 1)

		assert.Equal(t, "ldap", string(result.SensitiveDoguConfigDiff[0].Key.DoguName))
		assert.Equal(t, "logging/root", string(result.SensitiveDoguConfigDiff[0].Key.Key))
		assert.Equal(t, censorValue, result.SensitiveDoguConfigDiff[0].Actual.Value)
		assert.Equal(t, false, result.SensitiveDoguConfigDiff[0].Actual.Exists)
		assert.Equal(t, censorValue, result.SensitiveDoguConfigDiff[0].Expected.Value)
		assert.Equal(t, false, result.SensitiveDoguConfigDiff[0].Expected.Exists)
		assert.Equal(t, "Update", string(result.SensitiveDoguConfigDiff[0].NeededAction))
	})

	t.Run("Not censoring sensitive, but empty dogu config", func(t *testing.T) {
		//given
		configDiff := CombinedDoguConfigDiffs{
			SensitiveDoguConfigDiff: []SensitiveDoguConfigEntryDiff{
				{
					Actual: DoguConfigValueState{
						Value: "",
					},
					Expected: DoguConfigValueState{
						Value: "",
					},
				},
			},
		}

		//when
		result := configDiff.censorValues()

		require.Len(t, result.SensitiveDoguConfigDiff, 1)

		assert.Equal(t, "", result.SensitiveDoguConfigDiff[0].Actual.Value)
		assert.Equal(t, "", result.SensitiveDoguConfigDiff[0].Expected.Value)
	})
}

func TestCombinedDoguConfigDiffs_HasChanges(t *testing.T) {
	type fields struct {
		DoguConfigDiff          DoguConfigDiffs
		SensitiveDoguConfigDiff SensitiveDoguConfigDiffs
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "false on empty input",
			fields: fields{
				DoguConfigDiff:          []DoguConfigEntryDiff{},
				SensitiveDoguConfigDiff: []SensitiveDoguConfigEntryDiff{},
			},
			want: false,
		},
		{
			name: "true on non-empty Dogu config diff",
			fields: fields{
				DoguConfigDiff: []DoguConfigEntryDiff{{
					Key:          common.DoguConfigKey{DoguName: "testdogu", Key: "testkey"},
					Actual:       DoguConfigValueState{Value: "changed", Exists: true},
					Expected:     DoguConfigValueState{"initial", true},
					NeededAction: ConfigActionSet}},
				SensitiveDoguConfigDiff: []SensitiveDoguConfigEntryDiff{},
			},
			want: true,
		},
		{
			name: "true on non-empty sensitive config diff",
			fields: fields{
				DoguConfigDiff: []DoguConfigEntryDiff{},
				SensitiveDoguConfigDiff: []SensitiveDoguConfigEntryDiff{{
					Key:          common.SensitiveDoguConfigKey{common.DoguConfigKey{DoguName: "testdogu", Key: "testkey"}},
					Actual:       DoguConfigValueState{Value: "changed", Exists: true},
					Expected:     DoguConfigValueState{"initial", true},
					NeededAction: ConfigActionSet}},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			combinedDiff := CombinedDoguConfigDiffs{
				DoguConfigDiff:          tt.fields.DoguConfigDiff,
				SensitiveDoguConfigDiff: tt.fields.SensitiveDoguConfigDiff,
			}
			assert.Equalf(t, tt.want, combinedDiff.HasChanges(), "HasChanges()")
		})
	}
}
