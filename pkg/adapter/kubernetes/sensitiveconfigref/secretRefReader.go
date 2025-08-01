package sensitiveconfigref

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"iter"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"maps"
	"slices"
)

type SecretRefReader struct {
	secretClient k8sv1.SecretInterface
}

func NewSecretRefReader(secretClient secretClient) *SecretRefReader {
	return &SecretRefReader{
		secretClient: secretClient,
	}
}

func (reader *SecretRefReader) ExistAll(ctx context.Context, refs []domain.SensitiveValueRef) (bool, error) {
	secretsByName, secretsError := reader.loadNeededSecrets(ctx, slices.Values(refs))

	var notFoundErr *domainservice.NotFoundError

	if errors.As(secretsError, &notFoundErr) {
		// if any secret was not found, we can just return false because at least one does not exist
		return false, nil
	} else if secretsError != nil {
		// if it was another error, there could be connection problems, which we need to report to the caller
		return false, secretsError
	}

	for _, ref := range refs {
		if !reader.existKeyInSecret(secretsByName[ref.SecretName], ref.SecretKey) {
			return false, nil
		}
	}

	return true, nil
}

func (reader *SecretRefReader) GetValues(ctx context.Context, refs map[common.SensitiveDoguConfigKey]domain.SensitiveValueRef) (map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue, error) {
	secretsByName, secretErrors := reader.loadNeededSecrets(ctx, maps.Values(refs))
	sensitiveConfig, keyErrors := reader.loadKeysFromSecrets(refs, secretsByName)

	// combine errors so that the user gets info about not found secrets and missing keys in existing secrets
	err := errors.Join(secretErrors, keyErrors)
	if err != nil {
		err = fmt.Errorf("could not load sensitive config via references: %w", err)
		return nil, err
	}

	return sensitiveConfig, nil
}

func (reader *SecretRefReader) loadKeysFromSecrets(
	refs map[common.SensitiveDoguConfigKey]domain.SensitiveValueRef,
	secretsByName map[string]*v1.Secret,
) (map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue, error) {
	var errs []error
	loadedConfig := map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{}

	for configKey, ref := range refs {
		secret, found := secretsByName[ref.SecretName]
		if !found {
			// no error here, because we already have an error for missing secrets in the loadNeededSecrets function
			// we want error messages for missing keys too, even if a secret does not exist
			continue
		}
		sensitiveConfigValue, err := reader.loadKeyFromSecret(secret, ref.SecretKey)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		loadedConfig[configKey] = sensitiveConfigValue
	}
	return loadedConfig, errors.Join(errs...)
}

func (reader *SecretRefReader) loadKeyFromSecret(secret *v1.Secret, key string) (common.SensitiveDoguConfigValue, error) {
	valueBytes, exists := secret.Data[key]
	if !exists {
		return "", domainservice.NewNotFoundError(
			nil,
			"referenced secret key %q does not exist", key,
		)
	}
	return common.SensitiveDoguConfigValue(valueBytes), nil
}

func (reader *SecretRefReader) existKeyInSecret(secret *v1.Secret, key string) bool {
	_, exists := secret.Data[key]
	return exists
}

func (reader *SecretRefReader) loadNeededSecrets(
	ctx context.Context,
	refs iter.Seq[domain.SensitiveValueRef],
) (map[string]*v1.Secret, error) {
	secretsByName := map[string]*v1.Secret{}
	var errs []error

	for ref := range refs {
		_, alreadyLoaded := secretsByName[ref.SecretName]
		if alreadyLoaded {
			continue
		}
		secret, err := reader.loadSecret(ctx, ref.SecretName)
		if err != nil {
			errs = append(errs, err)
		}
		// also save nil entries, so that we do not try to load this secret again
		secretsByName[ref.SecretName] = secret
	}
	// delete nil entries
	maps.DeleteFunc(secretsByName, func(s string, secret *v1.Secret) bool {
		return secret == nil
	})
	return secretsByName, errors.Join(errs...)
}

func (reader *SecretRefReader) loadSecret(ctx context.Context, name string) (*v1.Secret, error) {
	secret, err := reader.secretClient.Get(ctx, name, metav1.GetOptions{})
	if secret == nil || err != nil {
		return nil, domainservice.NewNotFoundError(
			err, "referenced secret %q does not exist", name,
		)
	}
	return secret, nil
}
