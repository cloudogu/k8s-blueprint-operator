package v1

import (
	client "github.com/cloudogu/k8s-blueprint-lib/client"
	"k8s.io/client-go/tools/record"
)

type eventRecorder interface {
	record.EventRecorder
}

type blueprintInterface interface {
	client.BlueprintInterface
}
