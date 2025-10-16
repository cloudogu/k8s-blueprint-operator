package sensitiveconfigref

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

var testCtx = context.TODO()
var (
	redmineSecret = &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "postgres_credentials",
		},
		Data: map[string][]byte{"username": []byte("user1"), "password": []byte("123456")},
	}
	redmineCredentialsUsernameKey = common.DoguConfigKey{
		DoguName: "redmine",
		Key:      "credentials/username",
	}
	redmineCredentialsPasswordKey = common.DoguConfigKey{
		DoguName: "redmine",
		Key:      "credentials/password",
	}

	ldapSecret = &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ldap_credentials",
		},
		Data: map[string][]byte{"username": []byte("user2"), "password": []byte("789123")},
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

func TestSecretRefReader_GetValues(t *testing.T) {
	t.Run("nothing to load", func(t *testing.T) {
		secretMock := newMockSecretClient(t)
		refReader := NewSecretRefReader(secretMock)

		result, err := refReader.GetValues(testCtx, map[common.DoguConfigKey]domain.SensitiveValueRef{})
		require.NoError(t, err)
		assert.Equal(t, map[common.DoguConfigKey]common.SensitiveDoguConfigValue{}, result)
	})
	t.Run("nothing to load with nil input", func(t *testing.T) {
		secretMock := newMockSecretClient(t)
		refReader := NewSecretRefReader(secretMock)

		result, err := refReader.GetValues(testCtx, nil)
		require.NoError(t, err)
		assert.Equal(t, map[common.DoguConfigKey]common.SensitiveDoguConfigValue{}, result)
	})
	t.Run("load secrets with keys", func(t *testing.T) {
		secretMock := newMockSecretClient(t)
		secretMock.EXPECT().
			Get(testCtx, "postgres_credentials", metav1.GetOptions{}).
			Return(redmineSecret, nil)
		secretMock.EXPECT().
			Get(testCtx, "ldap_credentials", metav1.GetOptions{}).
			Return(ldapSecret, nil)

		refReader := NewSecretRefReader(secretMock)

		result, err := refReader.GetValues(testCtx,
			map[common.DoguConfigKey]domain.SensitiveValueRef{
				redmineCredentialsUsernameKey: {
					SecretName: redmineSecret.Name,
					SecretKey:  "username",
				},
				redmineCredentialsPasswordKey: {
					SecretName: redmineSecret.Name,
					SecretKey:  "password",
				},
				ldapCredentialsUsernameKey: {
					SecretName: ldapSecret.Name,
					SecretKey:  "username",
				},
				ldapCredentialsPasswordKey: {
					SecretName: ldapSecret.Name,
					SecretKey:  "password",
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
	t.Run("one secret and one key missing", func(t *testing.T) {
		secretMock := newMockSecretClient(t)
		secretMock.EXPECT().
			Get(testCtx, "postgres_credentials", metav1.GetOptions{}).
			Return(&corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "postgres_credentials",
				},
				StringData: map[string]string{"username": "user1"},
			}, nil)
		secretMock.EXPECT().
			Get(testCtx, "ldap_credentials", metav1.GetOptions{}).
			Return(nil, k8serrors.NewNotFound(schema.GroupResource{}, "ldap_credentials"))

		refReader := NewSecretRefReader(secretMock)

		_, err := refReader.GetValues(testCtx,
			map[common.DoguConfigKey]domain.SensitiveValueRef{
				redmineCredentialsUsernameKey: {
					SecretName: redmineSecret.Name,
					SecretKey:  "username",
				},
				redmineCredentialsPasswordKey: {
					SecretName: redmineSecret.Name,
					SecretKey:  "password",
				},
				ldapCredentialsUsernameKey: {
					SecretName: ldapSecret.Name,
					SecretKey:  "username",
				},
				ldapCredentialsPasswordKey: {
					SecretName: ldapSecret.Name,
					SecretKey:  "password",
				},
			},
		)
		require.Error(t, err)
		assert.ErrorContains(t, err, "could not load sensitive config via references")
		assert.ErrorContains(t, err, "referenced secret \"ldap_credentials\" does not exist")
		assert.ErrorContains(t, err, "referenced key \"password\" in secret \"postgres_credentials\" does not exist")
	})
}
