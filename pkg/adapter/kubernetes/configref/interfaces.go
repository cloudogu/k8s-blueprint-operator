package configref

import k8sv1 "k8s.io/client-go/kubernetes/typed/core/v1"

type configMapClient interface {
	k8sv1.ConfigMapInterface
}
