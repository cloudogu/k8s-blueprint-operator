package etcd

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/client/v2"
	"testing"
)

var testCtx = context.Background()

const testSimpleDoguNameRedmine = common.SimpleDoguName("redmine")

var etcdNotFoundError = client.Error{Code: client.ErrorCodeKeyNotFound}

func TestEtcdDoguConfigRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		etcdMock := newMockEtcdStore(t)
		configurationContextMock := newMockConfigurationContext(t)
		sut := EtcdDoguConfigRepository{etcdStore: etcdMock}

		key := common.DoguConfigKey{Key: "key", DoguName: testSimpleDoguNameRedmine}
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
		sut := EtcdDoguConfigRepository{etcdStore: etcdMock}

		key := common.DoguConfigKey{Key: "key", DoguName: testSimpleDoguNameRedmine}
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
		sut := EtcdDoguConfigRepository{etcdStore: etcdMock}

		key := common.DoguConfigKey{Key: "key", DoguName: testSimpleDoguNameRedmine}
		etcdMock.EXPECT().DoguConfig(string(testSimpleDoguNameRedmine)).Return(configurationContextMock)
		configurationContextMock.EXPECT().Delete(key.Key).Return(assert.AnError)

		// when
		err := sut.Delete(testCtx, key)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to delete config key \"key\" for dogu \"redmine\"")
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestEtcdDoguConfigRepository_GetAllByKey(t *testing.T) {

}

func TestEtcdDoguConfigRepository_Save(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		etcdMock := newMockEtcdStore(t)
		configurationContextMock := newMockConfigurationContext(t)
		sut := EtcdDoguConfigRepository{etcdStore: etcdMock}
		entry := &ecosystem.DoguConfigEntry{
			Key:                common.DoguConfigKey{Key: "key", DoguName: testSimpleDoguNameRedmine},
			Value:              "value",
			PersistenceContext: nil,
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
		sut := EtcdDoguConfigRepository{etcdStore: etcdMock}

		entry := &ecosystem.DoguConfigEntry{
			Key:                common.DoguConfigKey{Key: "key", DoguName: testSimpleDoguNameRedmine},
			Value:              "value",
			PersistenceContext: nil,
		}
		etcdMock.EXPECT().DoguConfig(string(testSimpleDoguNameRedmine)).Return(configurationContextMock)
		configurationContextMock.EXPECT().Set(entry.Key.Key, string(entry.Value)).Return(assert.AnError)

		// when
		err := sut.Save(testCtx, entry)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to set config key \"key\" with value \"value\" for dogu \"redmine\"")
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestEtcdDoguConfigRepository_SaveAll(t *testing.T) {

}

func TestNewEtcdDoguConfigRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		etcdMock := newMockEtcdStore(t)

		// when
		repository := NewEtcdDoguConfigRepository(etcdMock)

		// then
		assert.Equal(t, etcdMock, repository.etcdStore)
	})
}
