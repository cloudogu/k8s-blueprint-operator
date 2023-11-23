package main

import (
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type eventRecorder interface {
	record.EventRecorder
}

type controllerManager interface {
	manager.Manager
}
