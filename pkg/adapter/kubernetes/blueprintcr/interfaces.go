package blueprintcr

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/kubernetes"
	"k8s.io/client-go/tools/record"
)

type eventRecorder interface {
	record.EventRecorder
}

type blueprintInterface interface {
	kubernetes.BlueprintInterface
}
