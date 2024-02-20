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
		assert.ErrorContains(t, err, "failed to delete encrypted config key \"key\" for dogu \"redmine\"")
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestEtcdSensitiveDoguConfigRepository_GetAllByKey(t *testing.T) {

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
		assert.ErrorContains(t, err, "failed to set encrypted config key \"key\" with value \"value\" for dogu \"redmine\"")
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestEtcdSensitiveDoguConfigRepository_SaveAll(t *testing.T) {

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
