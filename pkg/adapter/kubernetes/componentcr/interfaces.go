package componentcr

import compCli "github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"

// used for mocks

//nolint:unused
//goland:noinspection GoUnusedType
type componentRepo interface {
	compCli.ComponentInterface
}
