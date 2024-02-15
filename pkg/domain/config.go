package domain

import "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"

type Config struct {
	Dogus  map[common.SimpleDoguName]CombinedDoguConfig
	Global GlobalConfig
}

type CombinedDoguConfig struct {
	DoguName        common.SimpleDoguName
	Config          DoguConfig
	SensitiveConfig SensitiveDoguConfig
}

type DoguConfig struct {
	Present map[common.DoguConfigKey]common.DoguConfigValue
	Absent  []common.DoguConfigKey
}

type SensitiveDoguConfig struct {
	Present map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue
	Absent  []common.SensitiveDoguConfigKey
}

type GlobalConfig struct {
	Present map[common.GlobalConfigKey]common.GlobalConfigValue
	Absent  []common.GlobalConfigKey
}
