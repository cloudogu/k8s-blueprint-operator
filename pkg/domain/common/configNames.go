package common

type GlobalConfigKey string

type DoguConfigKey struct {
	DoguName SimpleDoguName
	Key      string
}
type SensitiveDoguConfigKey DoguConfigKey
type GlobalConfigValue string
type DoguConfigValue string
type SensitiveDoguConfigValue string
type EncryptedDoguConfigValue string
