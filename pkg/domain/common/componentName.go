package common

import (
	"errors"
	"fmt"
	"strings"
)

type QualifiedComponentName struct {
	Namespace  ComponentNamespace
	SimpleName SimpleComponentName
}

type ComponentNamespace string
type SimpleComponentName string

func NewQualifiedComponentName(namespace ComponentNamespace, simpleName SimpleComponentName) (QualifiedComponentName, error) {
	componentName := QualifiedComponentName{Namespace: namespace, SimpleName: simpleName}
	err := componentName.Validate()
	if err != nil {
		return QualifiedComponentName{}, err
	}
	return QualifiedComponentName{Namespace: namespace, SimpleName: simpleName}, nil
}

func (componentName QualifiedComponentName) Validate() error {
	var errorList []error
	if componentName.Namespace == "" {
		errorList = append(errorList, fmt.Errorf("namespace of component %q must not be empty", componentName.SimpleName))
	}
	if componentName.SimpleName == "" {
		errorList = append(errorList, fmt.Errorf("component name must not be empty: '%s/%s'", componentName.Namespace, componentName.SimpleName))
	}
	return errors.Join(errorList...)
}

// String returns the component name with namespace, e.g. k8s/k8s-dogu-operator
func (componentName QualifiedComponentName) String() string {
	return fmt.Sprintf("%s/%s", componentName.Namespace, componentName.SimpleName)
}

// QualifiedComponentNameFromString converts a qualified component as a string, e.g. "k8s/k8s-dogu-operator", to a dedicated QualifiedComponentName or raises an error if this is not possible.
func QualifiedComponentNameFromString(qualifiedName string) (QualifiedComponentName, error) {
	splitName := strings.Split(qualifiedName, "/")
	if len(splitName) != 2 {
		return QualifiedComponentName{}, fmt.Errorf("component name needs to be in the form 'namespace/component' but is '%s'", qualifiedName)
	}
	return NewQualifiedComponentName(
		ComponentNamespace(splitName[0]),
		SimpleComponentName(splitName[1]),
	)
}
