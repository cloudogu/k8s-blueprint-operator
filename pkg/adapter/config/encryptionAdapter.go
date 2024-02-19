package config

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
)

type PublicKeyConfigEncryptionAdapter struct {
}

func NewPublicKeyConfigEncryptionAdapter() *PublicKeyConfigEncryptionAdapter {
	return &PublicKeyConfigEncryptionAdapter{}
}

func (p PublicKeyConfigEncryptionAdapter) Encrypt(
	ctx context.Context,
	name common.SimpleDoguName,
	value common.SensitiveDoguConfigValue,
) (common.EncryptedDoguConfigValue, error) {
	//TODO implement me
	panic("implement me")
}

func (p PublicKeyConfigEncryptionAdapter) EncryptAll(
	ctx context.Context,
	entries map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue,
) (map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue, error) {
	//TODO implement me
	panic("implement me")
}
