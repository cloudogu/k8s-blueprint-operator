package sensitiveconfigref

import k8sv1 "k8s.io/client-go/kubernetes/typed/core/v1"

type secretClient interface {
	k8sv1.SecretInterface
}
