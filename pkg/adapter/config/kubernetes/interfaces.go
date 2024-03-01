package kubernetes

import corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

type secretInterface interface {
	corev1.SecretInterface
}
