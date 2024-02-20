package k8s

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/retry"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

const doguSecretFormat = "%s-secrets"

type SecretSensitiveDoguConfigRepository struct {
	client secretInterface
}

func NewSecretSensitiveDoguConfigRepository(secretClient secretInterface) *SecretSensitiveDoguConfigRepository {
	return &SecretSensitiveDoguConfigRepository{
		client: secretClient,
	}
}

// SaveForNotInstalledDogu create or update a secret for the dogu containing the config entry.
// In further processing the dogu-operator uses the secret to encrypt configuration for the dogu.
func (repo *SecretSensitiveDoguConfigRepository) SaveForNotInstalledDogu(ctx context.Context, entry *ecosystem.SensitiveDoguConfigEntry) error {
	secretName := getDoguSecretName(string(entry.Key.DoguName))
	_, err := repo.client.Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			createErr := repo.createDoguSecret(ctx, entry.Key.DoguName)
			if createErr != nil {
				return createErr
			}
		} else {
			return getGetError(secretName, err)
		}
	}
	// update config map
	return retry.OnConflict(func() error {
		doguSecret, getErr := repo.client.Get(ctx, secretName, metav1.GetOptions{})
		if getErr != nil {
			return getGetError(secretName, getErr)
		}

		key, value := createKeyValueEntry(entry)
		if doguSecret.StringData == nil {
			doguSecret.StringData = map[string]string{}
		}
		doguSecret.StringData[key] = value

		_, updateErr := repo.client.Update(ctx, doguSecret, metav1.UpdateOptions{})
		return updateErr
	})
}

func (repo *SecretSensitiveDoguConfigRepository) createDoguSecret(ctx context.Context, doguName common.SimpleDoguName) error {
	secret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{
		Name: getDoguSecretName(string(doguName)),
	}}

	_, err := repo.client.Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create secret %q: %w", secret.Name, err)
	}

	return nil
}

func createKeyValueEntry(entry *ecosystem.SensitiveDoguConfigEntry) (key string, value string) {
	key = strings.ReplaceAll(entry.Key.Key, "/", ".")
	value = string(entry.Value)
	return
}

func getGetError(secretName string, err error) error {
	return fmt.Errorf("failed to get dogu secret %q: %w", secretName, err)
}

func getDoguSecretName(doguName string) string {
	return fmt.Sprintf(doguSecretFormat, doguName)
}
