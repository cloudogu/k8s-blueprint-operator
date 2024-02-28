package etcd

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEtcdSensitiveDoguConfigRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		etcdMock := newMockEtcdStore(t)
		configurationContextMock := newMockConfigurationContext(t)
		sut := EtcdSensitiveDoguConfigRepository{etcdStore: etcdMock}
		key := common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{Key: "key", DoguName: testSimpleDoguNameRedmine}}
		etcdMock.EXPECT().DoguConfig(string(testSimpleDoguNameRedmine)).Return(configurationContextMock)
		configurationContextMock.EXPECT().Delete(key.Key).Return(nil)

		// when
		err := sut.Delete(testCtx, key)

		// then
		require.NoError(t, err)
	})

	t.Run("should return nil on not found error", func(t *testing.T) {
		// given
		etcdMock := newMockEtcdStore(t)
		configurationContextMock := newMockConfigurationContext(t)
		sut := EtcdSensitiveDoguConfigRepository{etcdStore: etcdMock}
		key := common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{Key: "key", DoguName: testSimpleDoguNameRedmine}}
		etcdMock.EXPECT().DoguConfig(string(testSimpleDoguNameRedmine)).Return(configurationContextMock)
		configurationContextMock.EXPECT().Delete(key.Key).Return(etcdNotFoundError)

		// when
		err := sut.Delete(testCtx, key)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on delete error", func(t *testing.T) {
		// given
		etcdMock := newMockEtcdStore(t)
		configurationContextMock := newMockConfigurationContext(t)
		sut := EtcdSensitiveDoguConfigRepository{etcdStore: etcdMock}
		key := common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{Key: "key", DoguName: testSimpleDoguNameRedmine}}
		etcdMock.EXPECT().DoguConfig(string(testSimpleDoguNameRedmine)).Return(configurationContextMock)
		configurationContextMock.EXPECT().Delete(key.Key).Return(assert.AnError)

		// when
		err := sut.Delete(testCtx, key)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to delete encrypted key \"key\" of dogu \"redmine\" from etcd")
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorAs(t, err, &internalErr)
	})
}

func TestEtcdSensitiveDoguConfigRepository_GetAllByKey(t *testing.T) {
	t.Run("should fail with multiple keys", func(t *testing.T) {
		// given
		ldapConfigMock := newMockConfigurationContext(t)
		ldapConfigMock.EXPECT().Get("container_config/memory_limit").Return("", etcdNotFoundError)
		postfixConfigMock := newMockConfigurationContext(t)
		postfixConfigMock.EXPECT().Get("container_config/swap_limit").Return("", assert.AnError)
		etcdMock := newMockEtcdStore(t)
		etcdMock.EXPECT().DoguConfig("ldap").Return(ldapConfigMock)
		etcdMock.EXPECT().DoguConfig("postfix").Return(postfixConfigMock)
		sut := &EtcdSensitiveDoguConfigRepository{etcdStore: etcdMock}

		keys := []common.SensitiveDoguConfigKey{
			{common.DoguConfigKey{
				DoguName: "ldap",
				Key:      "container_config/memory_limit",
			}},
			{common.DoguConfigKey{
				DoguName: "postfix",
				Key:      "container_config/swap_limit",
			}},
		}

		// when
		_, err := sut.GetAllByKey(testCtx, keys)

		// then
		assert.ErrorIs(t, err, etcdNotFoundError)
		assert.ErrorAs(t, err, &notFoundErr)

		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorAs(t, err, &internalErr)
	})
	t.Run("should succeed with multiple keys", func(t *testing.T) {
		// given
		ldapConfigMock := newMockConfigurationContext(t)
		ldapConfigMock.EXPECT().Get("container_config/memory_limit").Return("1024m", nil)
		postfixConfigMock := newMockConfigurationContext(t)
		postfixConfigMock.EXPECT().Get("container_config/swap_limit").Return("512m", nil)
		etcdMock := newMockEtcdStore(t)
		etcdMock.EXPECT().DoguConfig("ldap").Return(ldapConfigMock)
		etcdMock.EXPECT().DoguConfig("postfix").Return(postfixConfigMock)
		sut := &EtcdSensitiveDoguConfigRepository{etcdStore: etcdMock}

		keys := []common.SensitiveDoguConfigKey{
			{common.DoguConfigKey{
				DoguName: "ldap",
				Key:      "container_config/memory_limit",
			}},
			{common.DoguConfigKey{
				DoguName: "postfix",
				Key:      "container_config/swap_limit",
			}},
		}

		// when
		actualEntries, err := sut.GetAllByKey(testCtx, keys)

		// then
		assert.NoError(t, err)
		expectedEntries := map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry{
			{DoguConfigKey: common.DoguConfigKey{
				DoguName: "ldap",
				Key:      "container_config/memory_limit",
			}}: {
				Key: common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
					DoguName: "ldap",
					Key:      "container_config/memory_limit",
				}},
				Value: "1024m",
			},
			{DoguConfigKey: common.DoguConfigKey{
				DoguName: "postfix",
				Key:      "container_config/swap_limit",
			}}: {
				Key: common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
					DoguName: "postfix",
					Key:      "container_config/swap_limit",
				}},
				Value: "512m",
			},
		}
		assert.Equal(t, expectedEntries, actualEntries)
	})
	t.Run("not found errors should produce the other values", func(t *testing.T) {
		// given
		ldapConfigMock := newMockConfigurationContext(t)
		ldapConfigMock.EXPECT().Get("container_config/memory_limit").Return("1024m", nil)
		ldapConfigMock.EXPECT().Get("password_change/notification_enabled").Return("", etcdNotFoundError)
		postfixConfigMock := newMockConfigurationContext(t)
		postfixConfigMock.EXPECT().Get("container_config/swap_limit").Return("512m", nil)
		etcdMock := newMockEtcdStore(t)
		etcdMock.EXPECT().DoguConfig("ldap").Return(ldapConfigMock)
		etcdMock.EXPECT().DoguConfig("postfix").Return(postfixConfigMock)
		sut := &EtcdSensitiveDoguConfigRepository{etcdStore: etcdMock}

		keys := []common.SensitiveDoguConfigKey{
			{common.DoguConfigKey{
				DoguName: "ldap",
				Key:      "container_config/memory_limit",
			}},
			{common.DoguConfigKey{
				DoguName: "ldap",
				Key:      "password_change/notification_enabled",
			}},
			{common.DoguConfigKey{
				DoguName: "postfix",
				Key:      "container_config/swap_limit",
			}},
		}

		// when
		actualEntries, err := sut.GetAllByKey(testCtx, keys)

		// then
		assert.ErrorIs(t, err, etcdNotFoundError)
		assert.ErrorAs(t, err, &notFoundErr)

		expectedEntries := map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry{
			{DoguConfigKey: common.DoguConfigKey{
				DoguName: "ldap",
				Key:      "container_config/memory_limit",
			}}: {
				Key: common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
					DoguName: "ldap",
					Key:      "container_config/memory_limit",
				}},
				Value: "1024m",
			},
			{DoguConfigKey: common.DoguConfigKey{
				DoguName: "postfix",
				Key:      "container_config/swap_limit",
			}}: {
				Key: common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
					DoguName: "postfix",
					Key:      "container_config/swap_limit",
				}},
				Value: "512m",
			},
		}
		assert.Equal(t, expectedEntries, actualEntries)
	})
}

func TestEtcdSensitiveDoguConfigRepository_Save(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		etcdMock := newMockEtcdStore(t)
		configurationContextMock := newMockConfigurationContext(t)
		sut := EtcdSensitiveDoguConfigRepository{etcdStore: etcdMock}
		entry := &ecosystem.SensitiveDoguConfigEntry{
			Key:   common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{Key: "key", DoguName: testSimpleDoguNameRedmine}},
			Value: "value",
		}
		etcdMock.EXPECT().DoguConfig(string(testSimpleDoguNameRedmine)).Return(configurationContextMock)
		configurationContextMock.EXPECT().Set(entry.Key.Key, string(entry.Value)).Return(nil)

		// when
		err := sut.Save(testCtx, entry)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on save error", func(t *testing.T) {
		// given
		etcdMock := newMockEtcdStore(t)
		configurationContextMock := newMockConfigurationContext(t)
		sut := EtcdSensitiveDoguConfigRepository{etcdStore: etcdMock}

		entry := &ecosystem.SensitiveDoguConfigEntry{
			Key:                common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{Key: "key", DoguName: testSimpleDoguNameRedmine}},
			Value:              "value",
			PersistenceContext: nil,
		}
		etcdMock.EXPECT().DoguConfig(string(testSimpleDoguNameRedmine)).Return(configurationContextMock)
		configurationContextMock.EXPECT().Set(entry.Key.Key, string(entry.Value)).Return(assert.AnError)

		// when
		err := sut.Save(testCtx, entry)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to set encrypted key \"key\" of dogu \"redmine\" with value \"value\" in etcd")
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorAs(t, err, &internalErr)
	})
}

func TestEtcdSensitiveDoguConfigRepository_SaveAll(t *testing.T) {
	t.Run("should fail to save multiple entries", func(t *testing.T) {
		// given
		ldapConfigMock := newMockConfigurationContext(t)
		ldapConfigMock.EXPECT().Set("container_config/memory_limit", "1024m").Return(assert.AnError)
		postfixConfigMock := newMockConfigurationContext(t)
		postfixConfigMock.EXPECT().Set("container_config/swap_limit", "512m").Return(anotherErr)
		etcdMock := newMockEtcdStore(t)
		etcdMock.EXPECT().DoguConfig("ldap").Return(ldapConfigMock)
		etcdMock.EXPECT().DoguConfig("postfix").Return(postfixConfigMock)

		sut := &EtcdSensitiveDoguConfigRepository{etcdStore: etcdMock}

		entries := []*ecosystem.SensitiveDoguConfigEntry{
			{
				Key: common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
					DoguName: "ldap",
					Key:      "container_config/memory_limit",
				}},
				Value: "1024m",
			},
			{
				Key: common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
					DoguName: "postfix",
					Key:      "container_config/swap_limit",
				}},
				Value: "512m",
			},
		}

		// when
		err := sut.SaveAll(testCtx, entries)

		// then
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorIs(t, err, anotherErr)
		assert.ErrorContains(t, err, "failed to set given sensitive dogu config entries in etcd")
	})
	t.Run("should succeed to save multiple entries", func(t *testing.T) {
		// given
		ldapConfigMock := newMockConfigurationContext(t)
		ldapConfigMock.EXPECT().Set("container_config/memory_limit", "1024m").Return(nil)
		postfixConfigMock := newMockConfigurationContext(t)
		postfixConfigMock.EXPECT().Set("container_config/swap_limit", "512m").Return(nil)
		etcdMock := newMockEtcdStore(t)
		etcdMock.EXPECT().DoguConfig("ldap").Return(ldapConfigMock)
		etcdMock.EXPECT().DoguConfig("postfix").Return(postfixConfigMock)

		sut := &EtcdSensitiveDoguConfigRepository{etcdStore: etcdMock}

		entries := []*ecosystem.SensitiveDoguConfigEntry{
			{
				Key: common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
					DoguName: "ldap",
					Key:      "container_config/memory_limit",
				}},
				Value: "1024m",
			},
			{
				Key: common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
					DoguName: "postfix",
					Key:      "container_config/swap_limit",
				}},
				Value: "512m",
			},
		}

		// when
		err := sut.SaveAll(testCtx, entries)

		// then
		assert.NoError(t, err)
	})
}

func TestNewEtcdSensitiveDoguConfigRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		etcdMock := newMockEtcdStore(t)

		// when
		repository := NewEtcdSensitiveDoguConfigRepository(etcdMock)

		// then
		assert.Equal(t, etcdMock, repository.etcdStore)
	})
}

func TestEtcdSensitiveDoguConfigRepository_Get(t *testing.T) {
	t.Run("should not find key", func(t *testing.T) {
		// given
		configurationContextMock := newMockConfigurationContext(t)
		configurationContextMock.EXPECT().Get("container_config/memory_limit").Return("", etcdNotFoundError)
		etcdMock := newMockEtcdStore(t)
		etcdMock.EXPECT().DoguConfig("ldap").Return(configurationContextMock)
		sut := &EtcdSensitiveDoguConfigRepository{etcdStore: etcdMock}

		key := common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
			DoguName: "ldap",
			Key:      "container_config/memory_limit",
		}}

		// when
		_, err := sut.Get(testCtx, key)

		// then
		assert.ErrorIs(t, err, etcdNotFoundError)
		assert.ErrorAs(t, err, &notFoundErr)
		assert.ErrorContains(t, err, "could not find sensitive key \"container_config/memory_limit\" of dogu \"ldap\" in etcd")
	})
	t.Run("should fail to get value for key", func(t *testing.T) {
		// given
		configurationContextMock := newMockConfigurationContext(t)
		configurationContextMock.EXPECT().Get("container_config/swap_limit").Return("", assert.AnError)
		etcdMock := newMockEtcdStore(t)
		etcdMock.EXPECT().DoguConfig("ldap").Return(configurationContextMock)
		sut := &EtcdSensitiveDoguConfigRepository{etcdStore: etcdMock}

		key := common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
			DoguName: "ldap",
			Key:      "container_config/swap_limit",
		}}

		// when
		_, err := sut.Get(testCtx, key)

		// then
		assert.ErrorAs(t, err, &internalErr)
		assert.ErrorContains(t, err, "failed to get sensitive key \"container_config/swap_limit\" of dogu \"ldap\" from etcd")
	})
	t.Run("should succeed to get value for key", func(t *testing.T) {
		// given
		configurationContextMock := newMockConfigurationContext(t)
		configurationContextMock.EXPECT().Get("container_config/swap_limit").Return("512m", nil)
		etcdMock := newMockEtcdStore(t)
		etcdMock.EXPECT().DoguConfig("ldap").Return(configurationContextMock)
		sut := &EtcdSensitiveDoguConfigRepository{etcdStore: etcdMock}

		key := common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
			DoguName: "ldap",
			Key:      "container_config/swap_limit",
		}}

		// when
		actualEntry, err := sut.Get(testCtx, key)

		// then
		assert.NoError(t, err)
		assert.Equal(t, &ecosystem.SensitiveDoguConfigEntry{
			Key: common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
				DoguName: "ldap",
				Key:      "container_config/swap_limit",
			}},
			Value: "512m",
		}, actualEntry)
	})
}

func TestEtcdSensitiveDoguConfigRepository_DeleteAllByKeys(t *testing.T) {
	t.Run("success with multiple keys", func(t *testing.T) {
		// given
		ldapConfigMock := newMockConfigurationContext(t)
		ldapConfigMock.EXPECT().Delete("container_config/memory_limit").Return(nil)
		postfixConfigMock := newMockConfigurationContext(t)
		postfixConfigMock.EXPECT().Delete("container_config/swap_limit").Return(nil)
		etcdMock := newMockEtcdStore(t)
		etcdMock.EXPECT().DoguConfig("ldap").Return(ldapConfigMock)
		etcdMock.EXPECT().DoguConfig("postfix").Return(postfixConfigMock)

		sut := &EtcdSensitiveDoguConfigRepository{etcdStore: etcdMock}

		entries := []common.SensitiveDoguConfigKey{
			{
				DoguConfigKey: common.DoguConfigKey{
					DoguName: "ldap",
					Key:      "container_config/memory_limit",
				},
			},
			{
				DoguConfigKey: common.DoguConfigKey{
					DoguName: "postfix",
					Key:      "container_config/swap_limit",
				},
			},
		}

		// when
		err := sut.DeleteAllByKeys(testCtx, entries)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on delete error", func(t *testing.T) {
		// given
		ldapConfigMock := newMockConfigurationContext(t)
		ldapConfigMock.EXPECT().Delete("container_config/memory_limit").Return(nil)
		postfixConfigMock := newMockConfigurationContext(t)
		postfixConfigMock.EXPECT().Delete("container_config/swap_limit").Return(assert.AnError)
		etcdMock := newMockEtcdStore(t)
		etcdMock.EXPECT().DoguConfig("ldap").Return(ldapConfigMock)
		etcdMock.EXPECT().DoguConfig("postfix").Return(postfixConfigMock)

		sut := &EtcdSensitiveDoguConfigRepository{etcdStore: etcdMock}

		entries := []common.SensitiveDoguConfigKey{
			{
				DoguConfigKey: common.DoguConfigKey{
					DoguName: "ldap",
					Key:      "container_config/memory_limit",
				},
			},
			{
				DoguConfigKey: common.DoguConfigKey{
					DoguName: "postfix",
					Key:      "container_config/swap_limit",
				},
			},
		}

		// when
		err := sut.DeleteAllByKeys(testCtx, entries)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to delete given sensitive dogu config keys in etcd")
	})
}
