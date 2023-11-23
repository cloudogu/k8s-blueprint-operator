package main

import (
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type controllerManager interface {
	manager.Manager
}
