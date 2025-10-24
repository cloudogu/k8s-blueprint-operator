package configref

import (
	"context"
	"testing"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var testCtx = context.Background()
var (
	redmineConfigMap = &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "postgres_credentials",
		},
		Data: map[string]string{"username": "user1", "password": "123456"},
	}
	redmineCredentialsUsernameKey = common.DoguConfigKey{
		DoguName: "redmine",
		Key:      "credentials/username",
	}
	redmineCredentialsPasswordKey = common.DoguConfigKey{
		DoguName: "redmine",
		Key:      "credentials/password",
	}

	ldapConfigMap = &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ldap_credentials",
		},
		Data: map[string]string{"username": "user2", "password": "789123"},
	}
	ldapCredentialsUsernameKey = common.DoguConfigKey{
		DoguName: "ldap",
		Key:      "credentials/username",
	}
	ldapCredentialsPasswordKey = common.DoguConfigKey{
		DoguName: "ldap",
		Key:      "credentials/password",
	}
)

func TestConfigMapRefReader_GetValues(t *testing.T) {
	t.Run("nothing to load", func(t *testing.T) {
		configMapMock := newMockConfigMapClient(t)
		refReader := NewConfigMapRefReader(configMapMock)

		result, err := refReader.GetValues(testCtx, map[common.DoguConfigKey]domain.ConfigValueRef{})
		require.NoError(t, err)
		assert.Equal(t, map[common.DoguConfigKey]common.DoguConfigValue{}, result)
	})
	t.Run("nothing to load with nil input", func(t *testing.T) {
		configMapMock := newMockConfigMapClient(t)
		refReader := NewConfigMapRefReader(configMapMock)

		result, err := refReader.GetValues(testCtx, nil)
		require.NoError(t, err)
		assert.Equal(t, map[common.DoguConfigKey]common.DoguConfigValue{}, result)
	})
	t.Run("load config maps with keys", func(t *testing.T) {
		configMapMock := newMockConfigMapClient(t)
		configMapMock.EXPECT().
			Get(testCtx, "postgres_credentials", metav1.GetOptions{}).
			Return(redmineConfigMap, nil)
		configMapMock.EXPECT().
			Get(testCtx, "ldap_credentials", metav1.GetOptions{}).
			Return(ldapConfigMap, nil)

		refReader := NewConfigMapRefReader(configMapMock)

		result, err := refReader.GetValues(testCtx,
			map[common.DoguConfigKey]domain.ConfigValueRef{
				redmineCredentialsUsernameKey: {
					ConfigMapName: redmineConfigMap.Name,
					ConfigMapKey:  "username",
				},
				redmineCredentialsPasswordKey: {
					ConfigMapName: redmineConfigMap.Name,
					ConfigMapKey:  "password",
				},
				ldapCredentialsUsernameKey: {
					ConfigMapName: ldapConfigMap.Name,
					ConfigMapKey:  "username",
				},
				ldapCredentialsPasswordKey: {
					ConfigMapName: ldapConfigMap.Name,
					ConfigMapKey:  "password",
				},
			},
		)
		require.NoError(t, err)
		assert.Equal(t, map[common.DoguConfigKey]common.SensitiveDoguConfigValue{
			redmineCredentialsUsernameKey: "user1",
			redmineCredentialsPasswordKey: "123456",
			ldapCredentialsUsernameKey:    "user2",
			ldapCredentialsPasswordKey:    "789123",
		}, result)
	})
	t.Run("one configmap and one key missing", func(t *testing.T) {
		configMapMock := newMockConfigMapClient(t)
		configMapMock.EXPECT().
			Get(testCtx, "postgres_credentials", metav1.GetOptions{}).
			Return(&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: "postgres_credentials",
				},
				Data: map[string]string{"username": "user1"},
			}, nil)
		configMapMock.EXPECT().
			Get(testCtx, "ldap_credentials", metav1.GetOptions{}).
			Return(nil, k8serrors.NewNotFound(schema.GroupResource{}, "ldap_credentials"))

		refReader := NewConfigMapRefReader(configMapMock)

		_, err := refReader.GetValues(testCtx,
			map[common.DoguConfigKey]domain.ConfigValueRef{
				redmineCredentialsUsernameKey: {
					ConfigMapName: redmineConfigMap.Name,
					ConfigMapKey:  "username",
				},
				redmineCredentialsPasswordKey: {
					ConfigMapName: redmineConfigMap.Name,
					ConfigMapKey:  "password",
				},
				ldapCredentialsUsernameKey: {
					ConfigMapName: ldapConfigMap.Name,
					ConfigMapKey:  "username",
				},
				ldapCredentialsPasswordKey: {
					ConfigMapName: ldapConfigMap.Name,
					ConfigMapKey:  "password",
				},
			},
		)
		require.Error(t, err)
		assert.ErrorContains(t, err, "could not load config via references")
		assert.ErrorContains(t, err, "referenced configMap \"ldap_credentials\" does not exist")
		assert.ErrorContains(t, err, "referenced key \"password\" in configMap \"postgres_credentials\" does not exist")
	})
}
