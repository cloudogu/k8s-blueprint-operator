package controller

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/api/ecosystem"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type eventRecorder interface {
	record.EventRecorder
}

type ecosystemClientSet interface {
	ecosystem.Interface
}

// used for mocks

//nolint:unused
//goland:noinspection GoUnusedType
type controllerManager interface {
	manager.Manager
}
