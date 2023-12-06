package controller

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/blueprint"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type eventRecorder interface {
	record.EventRecorder
}

type ecosystemClientSet interface {
	blueprint.Interface
}

// used for mocks

//nolint:unused
//goland:noinspection GoUnusedType
type controllerManager interface {
	manager.Manager
}
