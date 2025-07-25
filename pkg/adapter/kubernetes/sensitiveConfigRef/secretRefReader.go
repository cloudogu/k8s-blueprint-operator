package sensitiveConfigRef

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"maps"
)

type SecretRefReader struct {
	secretClient k8sv1.SecretInterface
}

func (reader *SecretRefReader) ExistAll(ctx context.Context, refs []domain.SensitiveValueRef) (bool, error) {
	secretsByName, secretErrors := reader.loadNeededSecrets(ctx, refs)
	_, keyErrors := reader.loadKeysFromSecrets(refs, secretsByName)

	err := errors.Join(secretErrors, keyErrors)
	var notFoundErr *domainservice.NotFoundError

	if errors.As(err, &notFoundErr) {
		return true, nil
	}
	return false, err
}

func (reader *SecretRefReader) GetValues(ctx context.Context, refs []map[common.SensitiveDoguConfigKey]domain.SensitiveValueRef) (map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue, error) {
	secretsByName, secretErrors := reader.loadNeededSecrets(ctx, refs)
	sensitiveConfig, keyErrors := reader.loadKeysFromSecrets(refs, secretsByName)

	err := errors.Join(secretErrors, keyErrors)
	if err != nil {
		err = fmt.Errorf("could not ")
	}

	return sensitiveConfig
}

func (reader *SecretRefReader) loadKeysFromSecrets(
	refs []map[common.SensitiveDoguConfigKey]domain.SensitiveValueRef,
	secretsByName map[string]*v1.Secret,
) (map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue, error) {
	var errs []error

	for _, ref := range refs {
		secret, found := secretsByName[ref.SecretName]
		if !found {
			continue
		}
		sensitiveConfigValue, err := reader.loadKeyFromSecret(secret, ref.SecretKey)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
}

func (reader *SecretRefReader) loadKeyFromSecret(secret *v1.Secret, key string) (string, error) {
	valueBytes, exists := secret.Data[key]
	if !exists {
		return "", domainservice.NewNotFoundError(
			nil,
			"referenced secret key does not exist", "secretKey", key,
		)
	}
	//TODO: check if the data is base64 encoded
	//decodeString, err := base64.StdEncoding.DecodeString(secretData)
	return string(valueBytes), nil
}

func (reader *SecretRefReader) loadNeededSecrets(ctx context.Context, refs []domain.SensitiveValueRef) (map[string]*v1.Secret, error) {
	secretsByName := map[string]*v1.Secret{}
	var errs []error

	for _, ref := range refs {
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
			err, "referenced secret does not exist", "secretName", name,
		)
	}
	return secret, nil
}
