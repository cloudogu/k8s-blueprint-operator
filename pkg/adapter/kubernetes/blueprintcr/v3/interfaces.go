package v3

import (
	bpv3client "github.com/cloudogu/k8s-blueprint-lib/v3/client"
	"k8s.io/client-go/tools/record"
)

type eventRecorder interface {
	record.EventRecorder
}

type blueprintInterface interface {
	bpv3client.BlueprintInterface
}

type blueprintMaskInterface interface {
	bpv3client.BlueprintMaskInterface
}
