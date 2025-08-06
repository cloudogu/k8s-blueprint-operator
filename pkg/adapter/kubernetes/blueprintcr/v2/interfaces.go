package v2

import (
	client "github.com/cloudogu/k8s-blueprint-lib/v2/client"
	"k8s.io/client-go/tools/record"
)

type eventRecorder interface {
	record.EventRecorder
}

type blueprintInterface interface {
	client.BlueprintInterface
}
