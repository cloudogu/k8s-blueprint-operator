package common

import (
	"errors"
	"fmt"
	"strings"
)

var K8sK8sLonghornName = QualifiedComponentName{Namespace: "k8s", Name: "k8s-longhorn"}

type QualifiedComponentName struct {
	Namespace ComponentNamespace
	Name      SimpleComponentName
}

type ComponentNamespace string
type SimpleComponentName string

func NewQualifiedComponentName(namespace ComponentNamespace, simpleName SimpleComponentName) (QualifiedComponentName, error) {
	componentName := QualifiedComponentName{Namespace: namespace, Name: simpleName}
	err := componentName.Validate()
	if err != nil {
		return QualifiedComponentName{}, err
	}
	return QualifiedComponentName{Namespace: namespace, Name: simpleName}, nil
}

func (componentName QualifiedComponentName) Validate() error {
	var errorList []error
	if componentName.Namespace == "" {
		errorList = append(errorList, fmt.Errorf("namespace of component %q must not be empty", componentName.Name))
	}
	if componentName.Name == "" {
		errorList = append(errorList, fmt.Errorf("component name must not be empty: '%s/%s'", componentName.Namespace, componentName.Name))
	}
	return errors.Join(errorList...)
}

// String returns the component name with namespace, e.g. k8s/k8s-dogu-operator
func (componentName QualifiedComponentName) String() string {
	return fmt.Sprintf("%s/%s", componentName.Namespace, componentName.Name)
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
