package config

import (
	"context"
	"github.com/cloudogu/cesapp-lib/keys"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/cloudogu/k8s-dogu-operator/controllers/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	ctx context.Context,
	name common.SimpleDoguName,
	value common.SensitiveDoguConfigValue,
) (common.EncryptedDoguConfigValue, error) {
	pubkey, err := resource.GetPublicKey(p.registry, string(name))
	if err != nil {
		return "", &domainservice.NotFoundError{
			WrappedError: err,
			Message:      "could not get public key",
		}
	}
	encryptedValue, err := pubkey.Encrypt(string(value))
	if err != nil {
		return "", domainservice.NewInternalError(err, "could not encrypt value")
	}
	return common.EncryptedDoguConfigValue(encryptedValue), nil
}

func (p PublicKeyConfigEncryptionAdapter) EncryptAll(
	ctx context.Context,
	entries map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue,
) (map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue, error) {
	encryptedEntries := map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue{}
	// TODO: only load public key once for every dogu
	for configKey, configValue := range entries {
		pubkey, err := resource.GetPublicKey(p.registry, string(configKey.DoguName))
		if err != nil {
			return map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue{}, &domainservice.NotFoundError{
				WrappedError: err,
				Message:      "could not get public key for dogu " + string(configKey.DoguName),
			}
		}
		encryptedValue, err := pubkey.Encrypt(string(configValue))
		if err != nil {
			return map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue{}, domainservice.NewInternalError(err, "could not encrypt value")
		}
		encryptedEntries[configKey] = common.EncryptedDoguConfigValue(encryptedValue)
	}
	return encryptedEntries, nil
}

func (p PublicKeyConfigEncryptionAdapter) Decrypt(
	ctx context.Context,
	name common.SimpleDoguName,
	encryptedValue common.EncryptedDoguConfigValue) (common.SensitiveDoguConfigValue, error) {
	privateKeySecret, err := p.secrets.Get(ctx, string(name)+"-private", metav1.GetOptions{})
	if err != nil {
		return "", domainservice.NewNotFoundError(err, "could not get private key")
	}
	privateKey := privateKeySecret.Data["private.pem"]
	keyPair, err := getKeyPairFromPrivateKey(privateKey, p.registry)
	decryptedValue, err := keyPair.Private().Decrypt(string(encryptedValue))
	if err != nil {
		return "", domainservice.NewInternalError(err, "could not decrypt encrypted value")
	}
	return common.SensitiveDoguConfigValue(decryptedValue), nil
}

func getKeyPairFromPrivateKey(privateKey []byte, registry etcdRegistry) (*keys.KeyPair, error) {
	keyProvider, err := resource.GetKeyProvider(registry)
	if err != nil {
		return &keys.KeyPair{}, domainservice.NewNotFoundError(err, "could not get key provider")
	}
	keyPair, err := keyProvider.FromPrivateKey(privateKey)
	if err != nil {
		return &keys.KeyPair{}, domainservice.NewInternalError(err, "could not get key pair")
	}
	return keyPair, nil
}

func (p PublicKeyConfigEncryptionAdapter) DecryptAll(ctx context.Context, entries map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue) (map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue, error) {
	decryptedEntries := map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{}
	// TODO: only load private key once for every dogu
	for configKey, configValue := range entries {
		privateKeySecret, err := p.secrets.Get(ctx, string(configKey.DoguName)+"-private", metav1.GetOptions{})
		if err != nil {
			return map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{}, domainservice.NewNotFoundError(err, "could not get private key")
		}
		privateKey := privateKeySecret.Data["private.pem"]
		keyPair, err := getKeyPairFromPrivateKey(privateKey, p.registry)
		decryptedValue, err := keyPair.Private().Decrypt(string(configValue))
		if err != nil {
			return map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{}, domainservice.NewInternalError(err, "could not decrypt encrypted value")
		}
		decryptedEntries[configKey] = common.SensitiveDoguConfigValue(decryptedValue)
	}
	return decryptedEntries, nil
}
