package dogucr

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func parseDoguCR(cr *v1.Dogu) (*ecosystem.DoguInstallation, error) {
	if cr == nil {
		return nil, &domainservice.InternalError{
			WrappedError: nil,
			Message:      "cannot parse dogu CR as it is nil",
		}
	}
	// parse dogu fields
	version, versionErr := core.ParseVersion(cr.Spec.Version)
	doguName, nameErr := common.QualifiedDoguNameFromString(cr.Spec.Name)
	err := errors.Join(versionErr, nameErr)
	if err != nil {
		return nil, &domainservice.InternalError{
			WrappedError: err,
			Message:      "cannot load dogu CR as it cannot be parsed correctly",
		}
	}
	// parse persistence context
	persistenceContext := make(map[string]interface{}, 1)
	persistenceContext[doguInstallationRepoContextKey] = doguInstallationRepoContext{
		resourceVersion: cr.GetResourceVersion(),
	}
	return &ecosystem.DoguInstallation{
		Name:               doguName,
		Version:            version,
		Status:             cr.Status.Status,
		Health:             ecosystem.HealthStatus(cr.Status.Health),
		UpgradeConfig:      ecosystem.UpgradeConfig{AllowNamespaceSwitch: cr.Spec.UpgradeConfig.AllowNamespaceSwitch},
		PersistenceContext: persistenceContext,
	}, nil
}

func toDoguCR(dogu *ecosystem.DoguInstallation) *v1.Dogu {
	return &v1.Dogu{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: string(dogu.Name.Name),
			Labels: map[string]string{
				"app":       "ces",
				"dogu.name": string(dogu.Name.Name),
			},
		},
		Spec: v1.DoguSpec{
			Name:    dogu.Name.String(),
			Version: dogu.Version.Raw,
			Resources: v1.DoguResources{
				DataVolumeSize: "",
			},
			SupportMode: false,
			UpgradeConfig: v1.UpgradeConfig{
				AllowNamespaceSwitch: dogu.UpgradeConfig.AllowNamespaceSwitch,
				ForceUpgrade:         false,
			},
			AdditionalIngressAnnotations: nil,
		},
		Status: v1.DoguStatus{},
	}
}

type doguCRPatch struct {
	Spec doguSpecPatch `json:"spec"`
}

type doguSpecPatch struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	//Resources     doguResourcesPatch `json:"resources"`
	SupportMode   bool               `json:"supportMode"`
	UpgradeConfig upgradeConfigPatch `json:"upgradeConfig"`
	//AdditionalIngressAnnotations ingressAnnotationsPatch `json:"additionalIngressAnnotations"`
}

//type ingressAnnotationsPatch map[string]string

type upgradeConfigPatch struct {
	AllowNamespaceSwitch bool `json:"allowNamespaceSwitch"`
	ForceUpgrade         bool `json:"forceUpgrade"`
}

// DoguResources defines the physical resources used by the dogu.
type doguResourcesPatch struct {
	DataVolumeSize string `json:"dataVolumeSize"`
}

func toDoguCRPatch(dogu *ecosystem.DoguInstallation) *doguCRPatch {
	return &doguCRPatch{
		Spec: doguSpecPatch{
			Name:    dogu.Name.String(),
			Version: dogu.Version.Raw,
			//Resources: doguResourcesPatch{
			//	DataVolumeSize: "",
			//},
			// always set this to false as a dogu cannot start in support mode
			SupportMode: false,
			UpgradeConfig: upgradeConfigPatch{
				AllowNamespaceSwitch: dogu.UpgradeConfig.AllowNamespaceSwitch,
				// this is a useful default as long as blueprints itself have no forceUpgrade flag implemented
				ForceUpgrade: false,
			},
		},
	}
}

func toDoguCRPatchBytes(dogu *ecosystem.DoguInstallation) ([]byte, error) {
	crPatch := toDoguCRPatch(dogu)
	patch, err := json.Marshal(crPatch)
	if err != nil {
		return []byte{}, &domainservice.InternalError{
			WrappedError: err,
			Message:      fmt.Sprintf("cannot patch dogu CR for dogu %q", dogu.Name),
		}
	}
	return patch, nil
}
