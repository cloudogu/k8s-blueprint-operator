package config

import (
	"context"
	"errors"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
)

const (
	privateKey                  = "private.pem"
	fmtDoguPrivateKeySecretName = "%s-private"
)

type PublicKeyConfigEncryptionAdapter struct {
	secrets   secret
	registry  etcdRegistry
	namespace string
}

func NewPublicKeyConfigEncryptionAdapter(secretClient secret, registry etcdRegistry, namespace string) *PublicKeyConfigEncryptionAdapter {
	return &PublicKeyConfigEncryptionAdapter{secrets: secretClient, registry: registry, namespace: namespace}
}

func (p PublicKeyConfigEncryptionAdapter) Encrypt(
	_ context.Context,
	name common.SimpleDoguName,
	value common.SensitiveDoguConfigValue,
) (common.EncryptedDoguConfigValue, error) {
	//TODO: The encryption got removed to update the dogu operator, which does not contain the encryption functions anymore.
	// This code is obsolet anyways as we will not encrypt config anymore but removing this adapter completely is a later step in the refactoring.
	return common.EncryptedDoguConfigValue(value), nil
}

func (p PublicKeyConfigEncryptionAdapter) EncryptAll(
	_ context.Context,
	entries map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue,
) (map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue, error) {
	//TODO: The encryption got removed to update the dogu operator, which does not contain the encryption functions anymore.
	// This code is obsolet anyways as we will not encrypt config anymore but removing this adapter completely is a later step in the refactoring.
	encryptedEntries := map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue{}
	var encryptionErrors []error
	for configKey, configValue := range entries {
		encryptedEntries[configKey] = common.EncryptedDoguConfigValue(configValue)
	}
	return encryptedEntries, errors.Join(encryptionErrors...)
}

func (p PublicKeyConfigEncryptionAdapter) Decrypt(
	ctx context.Context,
	name common.SimpleDoguName,
	encryptedValue common.EncryptedDoguConfigValue,
) (common.SensitiveDoguConfigValue, error) {
	//TODO: The encryption got removed to update the dogu operator, which does not contain the encryption functions anymore.
	// This code is obsolet anyways as we will not encrypt config anymore but removing this adapter completely is a later step in the refactoring.
	return common.SensitiveDoguConfigValue(encryptedValue), nil
}

func (p PublicKeyConfigEncryptionAdapter) DecryptAll(ctx context.Context, entries map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue) (map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue, error) {
	//TODO: The encryption got removed to update the dogu operator, which does not contain the encryption functions anymore.
	// This code is obsolet anyways as we will not encrypt config anymore but removing this adapter completely is a later step in the refactoring.
	decryptedEntries := map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{}
	for configKey, configValue := range entries {
		decryptedEntries[configKey] = common.SensitiveDoguConfigValue(configValue)
	}
	return decryptedEntries, nil
}
