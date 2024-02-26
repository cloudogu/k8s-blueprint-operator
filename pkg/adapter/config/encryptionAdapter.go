package config

import (
	"context"
	"errors"
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
		return "", domainservice.NewNotFoundError(err, "could not get public key")
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
	var encryptionErrors []error
	pubkeys, err := loadPubKeysFromMap(p.registry, entries)
	if err != nil {
		return encryptedEntries, domainservice.NewNotFoundError(err, "could not get public key(s)")
	}
	for configKey, configValue := range entries {
		doguname := string(configKey.DoguName)
		pubkey := pubkeys[doguname]
		encryptedValue, err := pubkey.Encrypt(string(configValue))
		if err != nil {
			encryptionErrors = append(encryptionErrors, domainservice.NewInternalError(err, "could not encrypt value for dogu %v", doguname))
		}
		encryptedEntries[configKey] = common.EncryptedDoguConfigValue(encryptedValue)
	}
	return encryptedEntries, errors.Join(encryptionErrors...)
}

func loadPubKeysFromMap(registry etcdRegistry, entries map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue) (map[string]*keys.PublicKey, error) {
	pubkeys := map[string]*keys.PublicKey{}
	var getPubKeyErrors []error
	for configKey := range entries {
		doguname := string(configKey.DoguName)
		_, pubkeyKnown := pubkeys[doguname]
		if !pubkeyKnown {
			pubkeyFromRegistry, err := resource.GetPublicKey(registry, doguname)
			if err != nil {
				getPubKeyErrors = append(getPubKeyErrors, domainservice.NewNotFoundError(err, "could not get public key for dogu %v", doguname))
				continue
			} else {
				pubkeys[doguname] = pubkeyFromRegistry
			}
		}
	}
	return pubkeys, errors.Join(getPubKeyErrors...)
}

func (p PublicKeyConfigEncryptionAdapter) Decrypt(
	ctx context.Context,
	name common.SimpleDoguName,
	encryptedValue common.EncryptedDoguConfigValue,
) (common.SensitiveDoguConfigValue, error) {
	privateKeySecret, err := p.secrets.Get(ctx, string(name)+"-private", metav1.GetOptions{})
	if err != nil {
		return "", domainservice.NewNotFoundError(err, "could not get private key")
	}
	privateKey := privateKeySecret.Data["private.pem"]
	keyPair, err := getKeyPairFromPrivateKey(privateKey, p.registry)
	if err != nil {
		return "", domainservice.NewInternalError(err, "could not get key pair for dogu %v", string(name))
	}
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
	var decryptionErrors []error
	keypairs, err := getKeyPairsFromMap(ctx, entries, p.secrets, p.registry)
	if err != nil {
		return decryptedEntries, domainservice.NewNotFoundError(err, "could not get key pair(s)")
	}
	for configKey, configValue := range entries {
		doguname := string(configKey.DoguName)
		keyPair := keypairs[doguname]
		decryptedValue, err := keyPair.Private().Decrypt(string(configValue))
		if err != nil {
			decryptionErrors = append(decryptionErrors, domainservice.NewInternalError(err, "could not decrypt encrypted value for dogu %v", doguname))
		}
		decryptedEntries[configKey] = common.SensitiveDoguConfigValue(decryptedValue)
	}
	return decryptedEntries, errors.Join(decryptionErrors...)
}

func getKeyPairsFromMap(ctx context.Context, entries map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue, secrets secret, registry etcdRegistry) (map[string]*keys.KeyPair, error) {
	var getKeyPairErrors []error
	keypairs := map[string]*keys.KeyPair{}
	for configKey := range entries {
		doguname := string(configKey.DoguName)
		_, privateKnown := keypairs[doguname]
		if !privateKnown {
			privateKeySecret, err := secrets.Get(ctx, doguname+"-private", metav1.GetOptions{})
			if err != nil {
				getKeyPairErrors = append(getKeyPairErrors, domainservice.NewNotFoundError(err, "could not get private key for dogu %v", doguname))
				continue
			} else {
				privateKey := privateKeySecret.Data["private.pem"]
				keyPairFromRegistry, err := getKeyPairFromPrivateKey(privateKey, registry)
				if err != nil {
					getKeyPairErrors = append(getKeyPairErrors, domainservice.NewNotFoundError(err, "could not get key pair for dogu %v", doguname))
					continue
				} else {
					keypairs[doguname] = keyPairFromRegistry
				}
			}
		}
	}
	return keypairs, errors.Join(getKeyPairErrors...)
}
