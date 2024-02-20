package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCombinedDoguConfigDiff_CensorValues(t *testing.T) {
	t.Run("Not censoring normal dogu config", func(t *testing.T) {
		//given
		configDiff := CombinedDoguConfigDiff{
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
					Action: "Update",
				},
			},
		}

		//when
		configDiff.censorValues()

		require.Len(t, configDiff.DoguConfigDiff, 1)

		assert.Equal(t, "ldap", string(configDiff.DoguConfigDiff[0].Key.DoguName))
		assert.Equal(t, "logging/root", configDiff.DoguConfigDiff[0].Key.Key)
		assert.Equal(t, "ERROR", configDiff.DoguConfigDiff[0].Actual.Value)
		assert.Equal(t, false, configDiff.DoguConfigDiff[0].Actual.Exists)
		assert.Equal(t, "DEBUG", configDiff.DoguConfigDiff[0].Expected.Value)
		assert.Equal(t, false, configDiff.DoguConfigDiff[0].Expected.Exists)
		assert.Equal(t, "Update", string(configDiff.DoguConfigDiff[0].Action))
	})

	t.Run("Censoring sensitive dogu config", func(t *testing.T) {
		//given
		configDiff := CombinedDoguConfigDiff{
			SensitiveDoguConfigDiff: []SensitiveDoguConfigEntryDiff{
				{
					Key: common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
						DoguName: "ldap",
						Key:      "logging/root",
					}},
					Actual: EncryptedDoguConfigValueState{
						Value:  "ERROR",
						Exists: false,
					},
					Expected: EncryptedDoguConfigValueState{
						Value:  "DEBUG",
						Exists: false,
					},
					Action: "Update",
				},
			},
		}

		//when
		configDiff.censorValues()

		require.Len(t, configDiff.SensitiveDoguConfigDiff, 1)

		assert.Equal(t, "ldap", string(configDiff.SensitiveDoguConfigDiff[0].Key.DoguName))
		assert.Equal(t, "logging/root", configDiff.SensitiveDoguConfigDiff[0].Key.Key)
		assert.Equal(t, censorValue, configDiff.SensitiveDoguConfigDiff[0].Actual.Value)
		assert.Equal(t, false, configDiff.SensitiveDoguConfigDiff[0].Actual.Exists)
		assert.Equal(t, censorValue, configDiff.SensitiveDoguConfigDiff[0].Expected.Value)
		assert.Equal(t, false, configDiff.SensitiveDoguConfigDiff[0].Expected.Exists)
		assert.Equal(t, "Update", string(configDiff.SensitiveDoguConfigDiff[0].Action))
	})

	t.Run("Not censoring sensitive, but empty dogu config", func(t *testing.T) {
		//given
		configDiff := CombinedDoguConfigDiff{
			SensitiveDoguConfigDiff: []SensitiveDoguConfigEntryDiff{
				{
					Actual: EncryptedDoguConfigValueState{
						Value: "",
					},
					Expected: EncryptedDoguConfigValueState{
						Value: "",
					},
				},
			},
		}

		//when
		configDiff.censorValues()

		require.Len(t, configDiff.SensitiveDoguConfigDiff, 1)

		assert.Equal(t, "", configDiff.SensitiveDoguConfigDiff[0].Actual.Value)
		assert.Equal(t, "", configDiff.SensitiveDoguConfigDiff[0].Expected.Value)
	})
}
