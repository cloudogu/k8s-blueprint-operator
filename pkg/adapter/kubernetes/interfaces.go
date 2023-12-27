package kubernetes

import "k8s.io/client-go/tools/record"

type eventRecorder interface {
	record.EventRecorder
}
