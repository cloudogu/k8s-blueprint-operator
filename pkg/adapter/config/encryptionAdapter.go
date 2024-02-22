package config

import (
	"context"
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
	//TODO implement me
	panic("implement me")
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

	keyProvider, err := resource.GetKeyProvider(p.registry)
	if err != nil {
		return "", domainservice.NewNotFoundError(err, "could not get key provider")
	}
	keyPair, err := keyProvider.FromPrivateKey(privateKey)
	if err != nil {
		return "", domainservice.NewInternalError(err, "could not get key pair")
	}
	decryptedValue, err := keyPair.Private().Decrypt(string(encryptedValue))
	if err != nil {
		return "", domainservice.NewInternalError(err, "could not decrypt encrypted value")
	}
	return common.SensitiveDoguConfigValue(decryptedValue), nil
}

func (p PublicKeyConfigEncryptionAdapter) DecryptAll(ctx context.Context, entries map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue) (map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue, error) {
	panic("implement me")
}
