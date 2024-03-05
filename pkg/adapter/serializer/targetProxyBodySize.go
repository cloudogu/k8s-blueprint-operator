package serializer

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"k8s.io/apimachinery/pkg/api/resource"
)

func ToDomainProxyBodySize(bodySize string) (*ecosystem.BodySize, error) {
	var quantity *ecosystem.BodySize
	var parse ecosystem.BodySize
	var err error
	if bodySize != "" && bodySize != "<nil>" {
		parse, err = resource.ParseQuantity(bodySize)
		quantity = &parse
	}

	return quantity, err
}
