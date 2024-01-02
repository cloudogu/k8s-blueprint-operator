package domain

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

type DoguHealthResult struct {
	UnhealthyDogus []UnhealthyDogu
}

type UnhealthyDogu struct {
	Namespace string
	Name      string
	Version   core.Version
	Health    ecosystem.HealthStatus
}
