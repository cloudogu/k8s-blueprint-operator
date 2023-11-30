package main

import (
	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type eventRecorder interface {
	record.EventRecorder
}

type controllerManager interface {
	manager.Manager
}

// used for mocks

//nolint:unused
//goland:noinspection GoUnusedType
type logSink interface {
	logr.LogSink
}
