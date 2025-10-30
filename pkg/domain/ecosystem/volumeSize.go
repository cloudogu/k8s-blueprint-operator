package ecosystem

import "k8s.io/apimachinery/pkg/api/resource"

type VolumeSize = resource.Quantity

func GetQuantityReference(quantityStr string) (*resource.Quantity, error) {
	var quantityValue resource.Quantity
	var err error
	if quantityStr != "" && quantityStr != "<nil>" {
		quantityValue, err = resource.ParseQuantity(quantityStr)
		if err == nil {
			return &quantityValue, nil
		}
	}
	return nil, err
}

func GetNonNilQuantityRef(quantityStr string) (*resource.Quantity, error) {
	quantityPtr, err := GetQuantityReference(quantityStr)
	if quantityPtr == nil {
		quantityPtr = &resource.Quantity{}
	}
	return quantityPtr, err
}

func GetQuantityString(quantity *resource.Quantity) *string {
	if quantity == nil {
		return nil
	}
	quantityStr := quantity.String()
	return &quantityStr
}
