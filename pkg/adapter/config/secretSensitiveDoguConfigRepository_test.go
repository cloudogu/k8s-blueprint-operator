package config

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"testing"
)

var notFoundErr = apierrors.NewNotFound(schema.GroupResource{}, "")

func TestSecretSensitiveDoguConfigRepository_SaveForNotInstalledDogu(t *testing.T) {
	t.Run("should update secret if it does not exist", func(t *testing.T) {
		// given
		secretMock := newMockSecretInterface(t)
		secretMock.EXPECT().Get(testCtx, string(testSimpleDoguNameRedmine+"-secrets"), metav1.GetOptions{}).Return(nil, notFoundErr).Times(1)
		expectedEmptySecret := &v1.Secret{ObjectMeta: metav1.ObjectMeta{
			Name: string(testSimpleDoguNameRedmine + "-secrets"),
		}}
		secretMock.EXPECT().Create(testCtx, expectedEmptySecret, metav1.CreateOptions{}).Return(nil, nil)
		secretMock.EXPECT().Get(testCtx, string(testSimpleDoguNameRedmine+"-secrets"), metav1.GetOptions{}).Return(expectedEmptySecret, nil).Times(1)
		expectedUpdateSecret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: string(testSimpleDoguNameRedmine + "-secrets"),
			},
			StringData: map[string]string{"key.path": "value"},
		}
		secretMock.EXPECT().Update(testCtx, expectedUpdateSecret, metav1.UpdateOptions{}).Return(nil, nil)

		sut := &SecretSensitiveDoguConfigRepository{
			client: secretMock,
		}

		entry := &ecosystem.DoguConfigEntry{
			Key: common.DoguConfigKey{
				DoguName: testSimpleDoguNameRedmine,
				Key:      "key/path",
			},
			Value: common.DoguConfigValue("value"),
		}

		// when
		err := sut.SaveForNotInstalledDogu(testCtx, entry)

		// then
		require.NoError(t, err)
	})

	t.Run("should retry on conflict error", func(t *testing.T) {
		// given
		secretMock := newMockSecretInterface(t)
		expectedEmptySecret := &v1.Secret{ObjectMeta: metav1.ObjectMeta{
			Name: string(testSimpleDoguNameRedmine + "-secrets"),
		}}
		secretMock.EXPECT().Get(testCtx, string(testSimpleDoguNameRedmine+"-secrets"), metav1.GetOptions{}).Return(expectedEmptySecret, nil).Times(3)
		expectedUpdateSecret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: string(testSimpleDoguNameRedmine + "-secrets"),
			},
			StringData: map[string]string{"key.path": "value"},
		}
		secretMock.EXPECT().Update(testCtx, expectedUpdateSecret, metav1.UpdateOptions{}).Return(nil, &apierrors.StatusError{ErrStatus: metav1.Status{Reason: metav1.StatusReasonConflict}}).Times(1)
		secretMock.EXPECT().Update(testCtx, expectedUpdateSecret, metav1.UpdateOptions{}).Return(nil, nil).Times(1)

		sut := &SecretSensitiveDoguConfigRepository{
			client: secretMock,
		}

		entry := &ecosystem.DoguConfigEntry{
			Key: common.DoguConfigKey{
				DoguName: testSimpleDoguNameRedmine,
				Key:      "key/path",
			},
			Value: common.DoguConfigValue("value"),
		}

		// when
		err := sut.SaveForNotInstalledDogu(testCtx, entry)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on create secret error", func(t *testing.T) {
		// given
		secretMock := newMockSecretInterface(t)
		secretMock.EXPECT().Get(testCtx, string(testSimpleDoguNameRedmine+"-secrets"), metav1.GetOptions{}).Return(nil, notFoundErr).Times(1)
		expectedEmptySecret := &v1.Secret{ObjectMeta: metav1.ObjectMeta{
			Name: string(testSimpleDoguNameRedmine + "-secrets"),
		}}
		secretMock.EXPECT().Create(testCtx, expectedEmptySecret, metav1.CreateOptions{}).Return(nil, assert.AnError)

		sut := &SecretSensitiveDoguConfigRepository{
			client: secretMock,
		}

		entry := &ecosystem.DoguConfigEntry{
			Key: common.DoguConfigKey{
				DoguName: testSimpleDoguNameRedmine,
				Key:      "key/path",
			},
			Value: common.DoguConfigValue("value"),
		}

		// when
		err := sut.SaveForNotInstalledDogu(testCtx, entry)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return get error on other get failure", func(t *testing.T) {
		// given
		secretMock := newMockSecretInterface(t)
		secretMock.EXPECT().Get(testCtx, string(testSimpleDoguNameRedmine+"-secrets"), metav1.GetOptions{}).Return(nil, assert.AnError).Times(1)

		sut := &SecretSensitiveDoguConfigRepository{
			client: secretMock,
		}

		entry := &ecosystem.DoguConfigEntry{
			Key: common.DoguConfigKey{
				DoguName: testSimpleDoguNameRedmine,
				Key:      "key/path",
			},
			Value: common.DoguConfigValue("value"),
		}

		// when
		err := sut.SaveForNotInstalledDogu(testCtx, entry)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get dogu secret \"redmine-secrets\"")
	})

	t.Run("should return get error on other get failure on on conflict procedure", func(t *testing.T) {
		// given
		secretMock := newMockSecretInterface(t)
		secretMock.EXPECT().Get(testCtx, string(testSimpleDoguNameRedmine+"-secrets"), metav1.GetOptions{}).Return(nil, nil).Times(1)
		secretMock.EXPECT().Get(testCtx, string(testSimpleDoguNameRedmine+"-secrets"), metav1.GetOptions{}).Return(nil, assert.AnError).Times(1)

		sut := &SecretSensitiveDoguConfigRepository{
			client: secretMock,
		}

		entry := &ecosystem.DoguConfigEntry{
			Key: common.DoguConfigKey{
				DoguName: testSimpleDoguNameRedmine,
				Key:      "key/path",
			},
			Value: common.DoguConfigValue("value"),
		}

		// when
		err := sut.SaveForNotInstalledDogu(testCtx, entry)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get dogu secret \"redmine-secrets\"")
	})
}

func TestNewSecretSensitiveDoguConfigRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		secretMock := newMockSecretInterface(t)

		// when
		repository := NewSecretSensitiveDoguConfigRepository(secretMock)

		// then
		require.NotNil(t, repository)
		assert.Equal(t, secretMock, repository.client)
	})
}
