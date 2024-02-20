package k8s

import (
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type secretInterface interface {
	v1.SecretInterface
}
