package serializer

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"k8s.io/apimachinery/pkg/api/resource"
)

func ToDomainVolumeSize(volumeSize string) (ecosystem.VolumeSize, error) {
	var quantity ecosystem.VolumeSize
	var err error
	if volumeSize != "" {
		quantity, err = resource.ParseQuantity(volumeSize)
	}

	return quantity, err
}
