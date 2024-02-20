package config

type SecretEtcdSensitiveDoguConfigRepository struct {
	*EtcdSensitiveDoguConfigRepository
	*SecretSensitiveDoguConfigRepository
}

func NewCombinedSecretEtcdSensitiveDoguConfigRepository(etcdRepo *EtcdSensitiveDoguConfigRepository, secretRepo *SecretSensitiveDoguConfigRepository) *SecretEtcdSensitiveDoguConfigRepository {
	return &SecretEtcdSensitiveDoguConfigRepository{
		EtcdSensitiveDoguConfigRepository:   etcdRepo,
		SecretSensitiveDoguConfigRepository: secretRepo,
	}
}
