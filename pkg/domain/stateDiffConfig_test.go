package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	dogu1Key1          = common.DoguConfigKey{DoguName: "dogu1", Key: "key1"}
	dogu1Key2          = common.DoguConfigKey{DoguName: "dogu1", Key: "key2"}
	dogu1Key3          = common.DoguConfigKey{DoguName: "dogu1", Key: "key3"}
	dogu1Key4          = common.DoguConfigKey{DoguName: "dogu1", Key: "key4"}
	sensitiveDogu1Key1 = common.SensitiveDoguConfigKey{DoguConfigKey: dogu1Key1}
	sensitiveDogu1Key2 = common.SensitiveDoguConfigKey{DoguConfigKey: dogu1Key2}
	sensitiveDogu1Key3 = common.SensitiveDoguConfigKey{DoguConfigKey: dogu1Key3}
	sensitiveDogu1Key4 = common.SensitiveDoguConfigKey{DoguConfigKey: dogu1Key4}
)

func Test_determineConfigDiff(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		config := Config{}
		dogusConfigDiffs, globalConfigDiff := determineConfigDiff(
			config,
			map[common.GlobalConfigKey]ecosystem.GlobalConfigEntry{},
			map[common.DoguConfigKey]ecosystem.DoguConfigEntry{},
			map[common.SensitiveDoguConfigKey]ecosystem.SensitiveDoguConfigEntry{},
			[]common.SimpleDoguName{},
		)

		assert.Equal(t, map[common.SimpleDoguName]CombinedDoguConfigDiff{}, dogusConfigDiffs)
		assert.Equal(t, GlobalConfigDiffs(nil), globalConfigDiff)
	})
	t.Run("all actions global config", func(t *testing.T) {
		//given ecosystem config
		actualGlobalConfig := map[common.GlobalConfigKey]ecosystem.GlobalConfigEntry{
			"key1": {Key: "key1", Value: "value1"}, // for action none
			"key2": {Key: "key2", Value: "value2"}, // for action set
			"key3": {Key: "key3", Value: "value3"}, // for action delete
			// key4 is absent -> action none
		}

		//given blueprint config
		config := Config{
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
		dogusConfigDiffs, globalConfigDiff := determineConfigDiff(
			config,
			actualGlobalConfig,
			map[common.DoguConfigKey]ecosystem.DoguConfigEntry{},
			map[common.SensitiveDoguConfigKey]ecosystem.SensitiveDoguConfigEntry{},
			[]common.SimpleDoguName{},
		)

		//then
		assert.Equal(t, map[common.SimpleDoguName]CombinedDoguConfigDiff{}, dogusConfigDiffs)
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
			Action: ConfigActionNone,
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
			Action: ConfigActionSet,
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
			Action: ConfigActionRemove,
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
			Action: ConfigActionNone,
		})
	})
	t.Run("all actions normal dogu config", func(t *testing.T) {
		//given ecosystem config
		actualDoguConfig := map[common.DoguConfigKey]ecosystem.DoguConfigEntry{
			dogu1Key1: {Key: dogu1Key1, Value: "value"}, //action none
			dogu1Key2: {Key: dogu1Key2, Value: "value"}, //action set
			dogu1Key3: {Key: dogu1Key3, Value: "value"}, //action delete
			//dogu1Key4 -> absent, so action none
		}
		installedDogus := []common.SimpleDoguName{"dogu1"}

		//given blueprint config
		config := Config{
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
		dogusConfigDiffs, globalConfigDiff := determineConfigDiff(
			config,
			map[common.GlobalConfigKey]ecosystem.GlobalConfigEntry{},
			actualDoguConfig,
			map[common.SensitiveDoguConfigKey]ecosystem.SensitiveDoguConfigEntry{},
			installedDogus,
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
			Action: ConfigActionNone,
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
			Action: ConfigActionSet,
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
			Action: ConfigActionRemove,
		})
		//domain.DoguConfigEntryDiff{Key:common.DoguConfigKey{DoguName:"dogu1", Key:"key3"},
		//Actual:domain.DoguConfigValueState{Value:"value", Exists:true},
		//Expected:domain.DoguConfigValueState{Value:"", Exists:false},
		//Action:"set"}
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
			Action: ConfigActionNone,
		})
	})
	t.Run("all actions for sensitive dogu config for present dogu", func(t *testing.T) {
		//given ecosystem config
		actualSensitiveConfig := map[common.SensitiveDoguConfigKey]ecosystem.SensitiveDoguConfigEntry{
			sensitiveDogu1Key1: {Key: sensitiveDogu1Key1, Value: "value"}, //action none
			sensitiveDogu1Key2: {Key: sensitiveDogu1Key2, Value: "value"}, //action set
			sensitiveDogu1Key3: {Key: sensitiveDogu1Key3, Value: "value"}, //action delete
			//sensitiveDogu1Key4 absent, so action none
		}
		//dogu1Key4 -> absent, so action none
		installedDogus := []common.SimpleDoguName{"dogu1"}

		//given blueprint config
		config := Config{
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
							sensitiveDogu1Key4,
						},
					},
				},
			},
		}

		//when
		dogusConfigDiffs, globalConfigDiff := determineConfigDiff(
			config,
			map[common.GlobalConfigKey]ecosystem.GlobalConfigEntry{},
			map[common.DoguConfigKey]ecosystem.DoguConfigEntry{},
			actualSensitiveConfig,
			installedDogus,
		)
		//then
		assert.Equal(t, GlobalConfigDiffs(nil), globalConfigDiff)
		require.NotNil(t, dogusConfigDiffs["dogu1"])
		assert.Equal(t, DoguConfigDiffs(nil), dogusConfigDiffs["dogu1"].DoguConfigDiff)
		assert.Equal(t, 4, len(dogusConfigDiffs["dogu1"].SensitiveDoguConfigDiff))
		assert.Contains(t, dogusConfigDiffs["dogu1"].SensitiveDoguConfigDiff, SensitiveDoguConfigEntryDiff{
			Key:                  sensitiveDogu1Key1,
			doguAlreadyInstalled: true,
			Actual: DoguConfigValueState{
				Value:  "value",
				Exists: true,
			},
			Expected: DoguConfigValueState{
				Value:  "value",
				Exists: true,
			},
			Action: ConfigActionNone,
		})
		assert.Contains(t, dogusConfigDiffs["dogu1"].SensitiveDoguConfigDiff, SensitiveDoguConfigEntryDiff{
			Key:                  sensitiveDogu1Key2,
			doguAlreadyInstalled: true,
			Actual: DoguConfigValueState{
				Value:  "value",
				Exists: true,
			},
			Expected: DoguConfigValueState{
				Value:  "updated value",
				Exists: true,
			},
			Action: ConfigActionSet,
		})
		assert.Contains(t, dogusConfigDiffs["dogu1"].SensitiveDoguConfigDiff, SensitiveDoguConfigEntryDiff{
			Key:                  sensitiveDogu1Key3,
			doguAlreadyInstalled: true,
			Actual: DoguConfigValueState{
				Value:  "value",
				Exists: true,
			},
			Expected: DoguConfigValueState{
				Value:  "",
				Exists: false,
			},
			Action: ConfigActionRemove,
		})
		assert.Contains(t, dogusConfigDiffs["dogu1"].SensitiveDoguConfigDiff, SensitiveDoguConfigEntryDiff{
			Key:                  sensitiveDogu1Key4,
			doguAlreadyInstalled: true,
			Actual: DoguConfigValueState{
				Value:  "",
				Exists: false,
			},
			Expected: DoguConfigValueState{
				Value:  "",
				Exists: false,
			},
			Action: ConfigActionNone,
		})
	})
	t.Run("all actions for sensitive dogu config for absent dogu", func(t *testing.T) {
		//given ecosystem config
		actualSensitiveConfig := map[common.SensitiveDoguConfigKey]ecosystem.SensitiveDoguConfigEntry{}
		//dogu1Key4 -> absent, so action none
		var installedDogus []common.SimpleDoguName

		//given blueprint config
		config := Config{
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
		dogusConfigDiffs, _ := determineConfigDiff(
			config,
			map[common.GlobalConfigKey]ecosystem.GlobalConfigEntry{},
			map[common.DoguConfigKey]ecosystem.DoguConfigEntry{},
			actualSensitiveConfig,
			installedDogus,
		)
		//then
		require.NotNil(t, dogusConfigDiffs["dogu1"])
		require.Equal(t, 1, len(dogusConfigDiffs["dogu1"].SensitiveDoguConfigDiff))
		assert.Equal(t, dogusConfigDiffs["dogu1"].SensitiveDoguConfigDiff[0], SensitiveDoguConfigEntryDiff{
			Key:                  sensitiveDogu1Key1,
			doguAlreadyInstalled: false,
			Actual: DoguConfigValueState{
				Value:  "",
				Exists: false,
			},
			Expected: DoguConfigValueState{
				Value:  "value",
				Exists: true,
			},
			Action: ConfigActionSetToEncrypt,
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
