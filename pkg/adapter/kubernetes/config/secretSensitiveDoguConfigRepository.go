package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/retry"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

const (
	doguSecretFormat       = "%s-secrets"
	etcdKeyDelimiter       = "/"
	doguSecretKeyDelimiter = "."
)

type SecretSensitiveDoguConfigRepository struct {
	client secretInterface
}

func NewSecretSensitiveDoguConfigRepository(secretClient secretInterface) *SecretSensitiveDoguConfigRepository {
	return &SecretSensitiveDoguConfigRepository{
		client: secretClient,
	}
}

// SaveAllForNotInstalledDogus creates or updates a secret for given sensitive dogu entries
// These entries can belong to different dogus.
func (repo *SecretSensitiveDoguConfigRepository) SaveAllForNotInstalledDogus(ctx context.Context, entries []*ecosystem.SensitiveDoguConfigEntry) error {
	groupByName := util.GroupBy(entries, func(entry *ecosystem.SensitiveDoguConfigEntry) common.SimpleDoguName {
		return entry.Key.DoguName
	})

	var errs []error
	for doguName, doguEntries := range groupByName {
		errs = append(errs, repo.SaveAllForNotInstalledDogu(ctx, doguName, doguEntries))
	}

	return errors.Join(errs...)
}

// SaveForNotInstalledDogu creates or updates a secret for the dogu containing the config entry.
// In further processing the dogu-operator uses the secret to encrypt configuration for the dogu.
func (repo *SecretSensitiveDoguConfigRepository) SaveForNotInstalledDogu(ctx context.Context, entry *ecosystem.SensitiveDoguConfigEntry) error {
	return repo.SaveAllForNotInstalledDogu(ctx, entry.Key.DoguName, []*ecosystem.SensitiveDoguConfigEntry{entry})
}

func (repo *SecretSensitiveDoguConfigRepository) SaveAllForNotInstalledDogu(ctx context.Context, doguName common.SimpleDoguName, entries []*ecosystem.SensitiveDoguConfigEntry) error {
	secretName, err := repo.checkAndCreateDoguSecret(ctx, doguName)
	if err != nil {
		return domainservice.NewInternalError(err, "failed to get or create dogu secret %q", secretName)
	}

	err = repo.updateSecretWithEntries(ctx, secretName, entries)
	if err != nil {
		return domainservice.NewInternalError(err, "failed to update dogu secret %q", secretName)
	}

	return nil
}

func (repo *SecretSensitiveDoguConfigRepository) updateSecretWithEntries(ctx context.Context, secretName string, entries []*ecosystem.SensitiveDoguConfigEntry) error {
	if len(entries) == 0 {
		return nil
	}

	return retry.OnConflict(func() error {
		doguSecret, getErr := repo.client.Get(ctx, secretName, metav1.GetOptions{})
		if getErr != nil {
			return getGetError(secretName, getErr)
		}

		if doguSecret.StringData == nil {
			doguSecret.StringData = map[string]string{}
		}

		for _, entry := range entries {
			key, value := createKeyValueEntry(entry)
			doguSecret.StringData[key] = value
		}

		_, updateErr := repo.client.Update(ctx, doguSecret, metav1.UpdateOptions{})
		return updateErr
	})
}

func (repo *SecretSensitiveDoguConfigRepository) checkAndCreateDoguSecret(ctx context.Context, doguName common.SimpleDoguName) (string, error) {
	secretName := getDoguSecretName(string(doguName))
	_, err := repo.client.Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			createErr := repo.createDoguSecret(ctx, doguName)
			if createErr != nil {
				return "", createErr
			}
		} else {
			return "", getGetError(secretName, err)
		}
	}

	return secretName, nil
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
	key = strings.ReplaceAll(entry.Key.Key, etcdKeyDelimiter, doguSecretKeyDelimiter)
	value = string(entry.Value)
	return
}

func getGetError(secretName string, err error) error {
	return fmt.Errorf("failed to get dogu secret %q: %w", secretName, err)
}

func getDoguSecretName(doguName string) string {
	return fmt.Sprintf(doguSecretFormat, doguName)
}
