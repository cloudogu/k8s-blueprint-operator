package etcd

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEtcdGlobalConfigRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigStore(t)
		sut := EtcdGlobalConfigRepository{configStore: globalConfigMock}
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
		sut := EtcdGlobalConfigRepository{configStore: globalConfigMock}
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
		sut := EtcdGlobalConfigRepository{configStore: globalConfigMock}
		key := common.GlobalConfigKey("key")
		globalConfigMock.EXPECT().Delete(string(key)).Return(assert.AnError)

		// when
		err := sut.Delete(testCtx, key)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to delete global config key \"key\"")
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestEtcdGlobalConfigRepository_GetAll(t *testing.T) {

}

func TestEtcdGlobalConfigRepository_Save(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigStore(t)
		sut := EtcdGlobalConfigRepository{configStore: globalConfigMock}
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
		sut := EtcdGlobalConfigRepository{configStore: globalConfigMock}
		entry := &ecosystem.GlobalConfigEntry{
			Key:   common.GlobalConfigKey("key"),
			Value: "value",
		}
		globalConfigMock.EXPECT().Set(string(entry.Key), string(entry.Value)).Return(assert.AnError)

		// when
		err := sut.Save(testCtx, entry)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to set global config key \"key\" with value \"value\"")
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestEtcdGlobalConfigRepository_SaveAll(t *testing.T) {

}

func TestNewEtcdGlobalConfigRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		globalConfigMock := newMockGlobalConfigStore(t)

		// when
		repository := NewEtcdGlobalConfigRepository(globalConfigMock)

		// then
		assert.Equal(t, globalConfigMock, repository.configStore)
	})
}
