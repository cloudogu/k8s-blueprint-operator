package etcd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

var internalErr = &domainservice.InternalError{}
var notFoundErr = &domainservice.NotFoundError{}
var anotherErr = fmt.Errorf("another error for testing")

func TestEtcdGlobalConfigRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigStore(t)
		sut := &GlobalConfigRepository{configStore: globalConfigMock}
		key := common.GlobalConfigKey("key")
		globalConfigMock.EXPECT().Delete(string(key)).Return(nil)

		// when
		err := sut.Delete(testCtx, key)

		// then
		require.NoError(t, err)
	})

	t.Run("should return nil on not found error", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigStore(t)
		sut := &GlobalConfigRepository{configStore: globalConfigMock}
		key := common.GlobalConfigKey("key")
		globalConfigMock.EXPECT().Delete(string(key)).Return(etcdNotFoundError)

		// when
		err := sut.Delete(testCtx, key)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on delete error", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigStore(t)
		sut := &GlobalConfigRepository{configStore: globalConfigMock}
		key := common.GlobalConfigKey("key")
		globalConfigMock.EXPECT().Delete(string(key)).Return(assert.AnError)

		// when
		err := sut.Delete(testCtx, key)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to delete global config key \"key\" from etcd")
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorAs(t, err, &internalErr)
	})
}

func TestEtcdGlobalConfigRepository_Save(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigStore(t)
		sut := &GlobalConfigRepository{configStore: globalConfigMock}
		entry := &ecosystem.GlobalConfigEntry{
			Key:   common.GlobalConfigKey("key"),
			Value: "value",
		}
		globalConfigMock.EXPECT().Set(string(entry.Key), string(entry.Value)).Return(nil)

		// when
		err := sut.Save(testCtx, entry)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on save error", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigStore(t)
		sut := &GlobalConfigRepository{configStore: globalConfigMock}
		entry := &ecosystem.GlobalConfigEntry{
			Key:   common.GlobalConfigKey("key"),
			Value: "value",
		}
		globalConfigMock.EXPECT().Set(string(entry.Key), string(entry.Value)).Return(assert.AnError)

		// when
		err := sut.Save(testCtx, entry)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to set global config key \"key\" with value \"value\" in etcd")
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorAs(t, err, &internalErr)
	})
}

func TestEtcdGlobalConfigRepository_SaveAll(t *testing.T) {
	t.Run("should fail to save with multiple errors", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigStore(t)
		globalConfigMock.EXPECT().Set("key_provider", "pkcs1v15").Return(assert.AnError)
		globalConfigMock.EXPECT().Set("fqdn", "ces.example.com").Return(anotherErr)
		sut := &GlobalConfigRepository{configStore: globalConfigMock}

		entries := []*ecosystem.GlobalConfigEntry{
			{
				Key:   "key_provider",
				Value: "pkcs1v15",
			},
			{
				Key:   "fqdn",
				Value: "ces.example.com",
			},
		}

		// when
		err := sut.SaveAll(testCtx, entries)

		// then
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorIs(t, err, anotherErr)
		assert.ErrorAs(t, err, &internalErr)
		assert.ErrorContains(t, err, "failed to set given global config entries in etcd")
		assert.ErrorContains(t, err, "failed to set global config key \"key_provider\" with value \"pkcs1v15\" in etcd")
		assert.ErrorContains(t, err, "failed to set global config key \"fqdn\" with value \"ces.example.com\" in etcd")
	})
	t.Run("should save multiple", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigStore(t)
		globalConfigMock.EXPECT().Set("key_provider", "pkcs1v15").Return(nil)
		globalConfigMock.EXPECT().Set("fqdn", "ces.example.com").Return(nil)
		sut := &GlobalConfigRepository{configStore: globalConfigMock}

		entries := []*ecosystem.GlobalConfigEntry{
			{
				Key:   "key_provider",
				Value: "pkcs1v15",
			},
			{
				Key:   "fqdn",
				Value: "ces.example.com",
			},
		}

		// when
		err := sut.SaveAll(testCtx, entries)

		// then
		assert.NoError(t, err)
	})
}

func TestNewEtcdGlobalConfigRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigStore(t)

		// when
		repository := NewGlobalConfigRepository(globalConfigMock)

		// then
		assert.Equal(t, globalConfigMock, repository.configStore)
	})
}

func TestEtcdGlobalConfigRepository_Get(t *testing.T) {
	t.Run("should not find key", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigStore(t)
		globalConfigMock.EXPECT().Get("key_provider").Return("", etcdNotFoundError)

		sut := &GlobalConfigRepository{configStore: globalConfigMock}

		// when
		_, err := sut.Get(testCtx, "key_provider")

		// then
		assert.ErrorIs(t, err, etcdNotFoundError)
		assert.ErrorAs(t, err, &notFoundErr)
		assert.ErrorContains(t, err, "could not find key \"key_provider\" from global config in etcd")
	})
	t.Run("should fail to get value for key", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigStore(t)
		globalConfigMock.EXPECT().Get("key_provider").Return("", assert.AnError)

		sut := &GlobalConfigRepository{configStore: globalConfigMock}

		// when
		_, err := sut.Get(testCtx, "key_provider")

		// then
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorAs(t, err, &internalErr)
		assert.ErrorContains(t, err, "failed to get value for key \"key_provider\" from global config in etcd")
	})
	t.Run("should get value for key", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigStore(t)
		globalConfigMock.EXPECT().Get("key_provider").Return("pkcs1v15", nil)

		sut := &GlobalConfigRepository{configStore: globalConfigMock}

		// when
		actualValue, err := sut.Get(testCtx, "key_provider")

		// then
		assert.NoError(t, err)
		assert.Equal(t, &ecosystem.GlobalConfigEntry{
			Key:   "key_provider",
			Value: "pkcs1v15",
		}, actualValue)
	})
}

func TestEtcdGlobalConfigRepository_GetAllByKey(t *testing.T) {
	t.Run("should fail with multiple errors", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigStore(t)
		globalConfigMock.EXPECT().Get("key_provider").Return("", etcdNotFoundError)
		globalConfigMock.EXPECT().Get("fqdn").Return("", assert.AnError)

		sut := &GlobalConfigRepository{configStore: globalConfigMock}

		keysToRetrieve := []common.GlobalConfigKey{"key_provider", "fqdn"}

		// when
		_, err := sut.GetAllByKey(testCtx, keysToRetrieve)

		// then
		assert.ErrorIs(t, err, etcdNotFoundError)
		assert.ErrorAs(t, err, &notFoundErr)

		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorAs(t, err, &internalErr)
	})
	t.Run("should return entries on error", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigStore(t)
		globalConfigMock.EXPECT().Get("key_provider").Return("", etcdNotFoundError)
		globalConfigMock.EXPECT().Get("fqdn").Return("ces.example.com", nil)

		sut := &GlobalConfigRepository{configStore: globalConfigMock}

		keysToRetrieve := []common.GlobalConfigKey{"key_provider", "fqdn"}

		// when
		actualEntries, err := sut.GetAllByKey(testCtx, keysToRetrieve)

		// then
		assert.ErrorIs(t, err, etcdNotFoundError)
		assert.ErrorAs(t, err, &notFoundErr)

		assert.Equal(t, map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry{
			"fqdn": {
				Key:   "fqdn",
				Value: "ces.example.com",
			},
		}, actualEntries)
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigStore(t)
		globalConfigMock.EXPECT().Get("key_provider").Return("pkcs1v15", nil)
		globalConfigMock.EXPECT().Get("fqdn").Return("ces.example.com", nil)

		sut := &GlobalConfigRepository{configStore: globalConfigMock}

		keysToRetrieve := []common.GlobalConfigKey{"key_provider", "fqdn"}

		// when
		actualEntries, err := sut.GetAllByKey(testCtx, keysToRetrieve)

		// then
		assert.NoError(t, err)
		assert.Equal(t, map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry{
			"fqdn": {
				Key:   "fqdn",
				Value: "ces.example.com",
			},
			"key_provider": {
				Key:   "key_provider",
				Value: "pkcs1v15",
			},
		}, actualEntries)
	})
}

func TestEtcdGlobalConfigRepository_DeleteAllByKeys(t *testing.T) {
	t.Run("success with multiple keys", func(t *testing.T) {
		// given
		globalConfigMock := newMockConfigurationContext(t)
		globalConfigMock.EXPECT().Delete("fqdn").Return(nil)
		globalConfigMock.EXPECT().Delete("certificate").Return(nil)

		sut := &GlobalConfigRepository{configStore: globalConfigMock}

		entries := []common.GlobalConfigKey{"fqdn", "certificate"}

		// when
		err := sut.DeleteAllByKeys(testCtx, entries)

		// then
		assert.NoError(t, err)
	})

	t.Run("should return error on delete error", func(t *testing.T) {
		// given
		globalConfigMock := newMockConfigurationContext(t)
		globalConfigMock.EXPECT().Delete("fqdn").Return(nil)
		globalConfigMock.EXPECT().Delete("certificate").Return(assert.AnError)

		sut := &GlobalConfigRepository{configStore: globalConfigMock}

		entries := []common.GlobalConfigKey{"fqdn", "certificate"}

		// when
		err := sut.DeleteAllByKeys(testCtx, entries)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to delete given global config keys in etcd")
	})
}
