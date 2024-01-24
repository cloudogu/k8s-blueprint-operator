package domain

import "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"

type HealthConfig struct {
	ComponentHealthConfig `yaml:"components" json:"components"`
}

type ComponentHealthConfig struct {
	RequiredComponents []ecosystem.RequiredComponent `yaml:"required" json:"required"`
}
