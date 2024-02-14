package serializer

import (
	"fmt"
	"strings"
)

const componentDistributionNamespaceDelimiter = "/"

// SplitComponentName splits a qualified component name into the distribution namespace and the simple name or raises an
// error if this is not possible.
//
//	"k8s/nginx-static" -> "k8s", "nginx-static"
func SplitComponentName(componentName string) (distributionNamespace string, name string, err error) {
	splitName := strings.Split(componentName, componentDistributionNamespaceDelimiter)
	if len(splitName) != 2 {
		return "", "", fmt.Errorf("component name needs to be in the form 'namespace/component' but is '%s'", componentName)
	}
	return splitName[0], splitName[1], nil
}

// JoinComponentName splits a qualified component name into the distribution namespace and the simple name or raises an
// error if this is not possible.
//
//	"k8s/nginx-static" -> "k8s", "nginx-static"
func JoinComponentName(componentName, distributionNamespace string) (string, error) {
	if distributionNamespace == "" {
		return "", fmt.Errorf("distribution namespace of component %s must not be empty", componentName)
	}

	return strings.Join([]string{distributionNamespace, componentName}, componentDistributionNamespaceDelimiter), nil
}
